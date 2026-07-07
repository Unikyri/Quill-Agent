package services

import (
	"context"
	"testing"
	"time"

	"github.com/quill/backend/internal/repositories"
	"github.com/quill/backend/internal/testutil"
)

// TestDeriveLifecycle is a pure table test covering all 5 lifecycle states,
// the epsilon boundary for "decaying", and the archived->active
// flip-within-window rule for "reactivated" — no DB required.
func TestDeriveLifecycle(t *testing.T) {
	now := time.Now()
	pt := func(score float64, status string, offsetSeconds int) repositories.RelevanceHistoryPoint {
		return repositories.RelevanceHistoryPoint{
			RelevanceScore: score,
			Status:         status,
			RecordedAt:     now.Add(time.Duration(offsetSeconds) * time.Second),
		}
	}

	const epsilon = 0.01

	tests := []struct {
		name         string
		status       string
		history      []repositories.RelevanceHistoryPoint
		consolidated bool
		want         string
	}{
		{
			name:   "active default with no history",
			status: "active",
			want:   "active",
		},
		{
			name:    "active with a single history row stays active",
			status:  "active",
			history: []repositories.RelevanceHistoryPoint{pt(0.8, "active", 0)},
			want:    "active",
		},
		{
			name:   "decaying when latest score drops beyond epsilon",
			status: "active",
			history: []repositories.RelevanceHistoryPoint{
				pt(0.8, "active", 0),
				pt(0.5, "active", 1),
			},
			want: "decaying",
		},
		{
			name:   "not decaying when drop is within epsilon (float noise)",
			status: "active",
			history: []repositories.RelevanceHistoryPoint{
				pt(0.800, "active", 0),
				pt(0.795, "active", 1), // delta -0.005, |delta| < epsilon(0.01)
			},
			want: "active",
		},
		{
			name:   "decaying just past epsilon boundary",
			status: "active",
			history: []repositories.RelevanceHistoryPoint{
				pt(0.80, "active", 0),
				pt(0.7899, "active", 1), // delta -0.0101, strictly less than -epsilon
			},
			want: "decaying",
		},
		{
			name:         "archived without consolidation",
			status:       "archived",
			consolidated: false,
			want:         "archived",
		},
		{
			name:         "consolidated takes priority over plain archived",
			status:       "archived",
			consolidated: true,
			want:         "consolidated",
		},
		{
			name:   "reactivated when most recent transition is archived->active",
			status: "active",
			history: []repositories.RelevanceHistoryPoint{
				pt(0.5, "active", 0),
				pt(0.1, "archived", 1),
				pt(0.8, "active", 2),
			},
			want: "reactivated",
		},
		{
			name:   "reactivated flag persists while no newer transition supersedes it",
			status: "active",
			history: []repositories.RelevanceHistoryPoint{
				pt(0.1, "archived", 0),
				pt(0.8, "active", 1),
				pt(0.82, "active", 2), // score increase after reactivation, no new transition
			},
			want: "reactivated",
		},
		{
			name:   "reactivated re-detected after an intervening active->archived->active cycle",
			status: "active",
			history: []repositories.RelevanceHistoryPoint{
				pt(0.5, "active", 0),
				pt(0.1, "archived", 1), // older transition: active->archived
				pt(0.6, "active", 2),   // most recent transition: archived->active
			},
			want: "reactivated",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := deriveLifecycle(tt.status, tt.history, tt.consolidated, epsilon)
			if got != tt.want {
				t.Errorf("deriveLifecycle(%q, %v, %v) = %q, want %q", tt.status, tt.history, tt.consolidated, got, tt.want)
			}
		})
	}
}

// TestMemoryServiceMemoryStatusZeroHistoryAndCapping is an integration test
// (requires TEST_DATABASE_URL) covering the spec's "zero-history entity" and
// "large universe capping" scenarios end to end through MemoryStatus.
func TestMemoryServiceMemoryStatusZeroHistoryAndCapping(t *testing.T) {
	pool := testutil.SetupTestDB(t)
	testutil.RunMigrationsUpTo(t, pool, "019")
	ctx := context.Background()

	user := svcCreateTestUser(t, ctx, pool)
	universe := svcCreateTestUniverse(t, ctx, pool, user.ID)

	freshEntity := svcCreateTestEntity(t, ctx, pool, universe.ID, "Fresh Entity", 0.8, "active")
	busyEntity := svcCreateTestEntity(t, ctx, pool, universe.ID, "Busy Entity", 0.3, "active")

	historyRepo := repositories.NewEntityRelevanceHistoryRepo(pool)
	// 35 rows for busyEntity — only the most recent 30 should come back.
	for i := 0; i < 35; i++ {
		if err := historyRepo.AppendOne(ctx, busyEntity.ID); err != nil {
			t.Fatalf("append history row %d: %v", i, err)
		}
	}

	entityRepo := repositories.NewEntityRepo(pool)
	svc := NewMemoryService(repositories.NewGraphRepo(pool), entityRepo, repositories.NewVectorRepo(pool))
	svc.SetHistoryRepo(historyRepo)

	status, err := svc.MemoryStatus(ctx, universe.ID)
	if err != nil {
		t.Fatalf("MemoryStatus: %v", err)
	}

	byID := make(map[string]MemoryStatusEntity, len(status.Entities))
	for _, e := range status.Entities {
		byID[e.ID.String()] = e
	}

	fresh, ok := byID[freshEntity.ID.String()]
	if !ok {
		t.Fatalf("expected fresh entity in response, got: %+v", status.Entities)
	}
	if fresh.History == nil || len(fresh.History) != 0 {
		t.Errorf("expected zero-history entity to have empty (non-nil) history slice, got %v", fresh.History)
	}
	if fresh.Lifecycle != "active" {
		t.Errorf("expected zero-history entity lifecycle 'active', got %q", fresh.Lifecycle)
	}

	busy, ok := byID[busyEntity.ID.String()]
	if !ok {
		t.Fatalf("expected busy entity in response, got: %+v", status.Entities)
	}
	if len(busy.History) != 30 {
		t.Errorf("expected capped history of 30 rows, got %d", len(busy.History))
	}
}
