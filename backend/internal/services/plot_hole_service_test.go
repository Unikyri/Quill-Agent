package services

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/quill/backend/internal/models"
	"github.com/quill/backend/internal/repositories"
	"github.com/quill/backend/internal/testutil"
)

func TestPlotHoleServiceScanDetectsStaleArc(t *testing.T) {
	pool := testutil.SetupTestDB(t)
	testutil.RunMigrationsUpTo(t, pool, "011")
	ctx := context.Background()

	user := svcCreateTestUser(t, ctx, pool)
	universe := svcCreateTestUniverse(t, ctx, pool, user.ID)
	chapters := svcCreateChapters(t, ctx, pool, universe, 10)
	ch1 := chapters[0]   // order_index=1
	ch10 := chapters[9]  // order_index=10

	// Create entity with last_mentioned at chapter 1
	entity := svcCreateTestEntity(t, ctx, pool, universe.ID, "Forgotten Hero", 0.3, "active")
	// Update last_mentioned_chapter_id to ch1
	if _, err := pool.Exec(ctx, "UPDATE entities SET last_mentioned_chapter_id = $1, last_mentioned_at = NOW() WHERE id = $2",
		ch1.ID, entity.ID); err != nil {
		t.Fatalf("set last_mentioned: %v", err)
	}

	svc := NewPlotHoleService(pool, repositories.NewPlotHoleRepo(pool), repositories.NewEntityRepo(pool), 8)

	holes, err := svc.Scan(ctx, universe.ID, ch10.ID)
	if err != nil {
		t.Fatalf("Scan failed: %v", err)
	}

	if len(holes) == 0 {
		t.Error("expected at least 1 plot hole for stale entity (gap=9 ≥ 8)")
	}
	for _, h := range holes {
		if h.Status != "open" {
			t.Errorf("plot hole %s status = %s, want open", h.Title, h.Status)
		}
	}
}

// TestPlotHoleServiceScanNoGap ensures entities within threshold don't trigger.
func TestPlotHoleServiceScanNoGap(t *testing.T) {
	pool := testutil.SetupTestDB(t)
	testutil.RunMigrationsUpTo(t, pool, "011")
	ctx := context.Background()

	user := svcCreateTestUser(t, ctx, pool)
	universe := svcCreateTestUniverse(t, ctx, pool, user.ID)
	chapters := svcCreateChapters(t, ctx, pool, universe, 5)
	ch1 := chapters[0]
	ch5 := chapters[4]

	entity := svcCreateTestEntity(t, ctx, pool, universe.ID, "Current Hero", 0.5, "active")
	if _, err := pool.Exec(ctx, "UPDATE entities SET last_mentioned_chapter_id = $1, last_mentioned_at = NOW() WHERE id = $2",
		ch1.ID, entity.ID); err != nil {
		t.Fatalf("set last_mentioned: %v", err)
	}

	svc := NewPlotHoleService(pool, repositories.NewPlotHoleRepo(pool), repositories.NewEntityRepo(pool), 8)

	holes, err := svc.Scan(ctx, universe.ID, ch5.ID)
	if err != nil {
		t.Fatalf("Scan failed: %v", err)
	}

	if len(holes) != 0 {
		t.Errorf("expected 0 plot holes (gap=4 < 8), got %d", len(holes))
	}
}

// TestPlotHoleServiceScanMixed has both stale and current entities.
func TestPlotHoleServiceScanMixed(t *testing.T) {
	pool := testutil.SetupTestDB(t)
	testutil.RunMigrationsUpTo(t, pool, "011")
	ctx := context.Background()

	user := svcCreateTestUser(t, ctx, pool)
	universe := svcCreateTestUniverse(t, ctx, pool, user.ID)
	chapters := svcCreateChapters(t, ctx, pool, universe, 10)
	ch1 := chapters[0]
	ch8 := chapters[7]  // order_index=8
	ch10 := chapters[9] // order_index=10

	// Stale entity: last mentioned at chapter 1, scan at chapter 10 → gap 9
	stale := svcCreateTestEntity(t, ctx, pool, universe.ID, "Stale Arc", 0.3, "active")
	pool.Exec(ctx, "UPDATE entities SET last_mentioned_chapter_id = $1, last_mentioned_at = NOW() WHERE id = $2", ch1.ID, stale.ID)

	// Recent entity: last mentioned at chapter 8, scan at chapter 10 → gap 2
	recent := svcCreateTestEntity(t, ctx, pool, universe.ID, "Recent Arc", 0.5, "active")
	pool.Exec(ctx, "UPDATE entities SET last_mentioned_chapter_id = $1, last_mentioned_at = NOW() WHERE id = $2", ch8.ID, recent.ID)

	svc := NewPlotHoleService(pool, repositories.NewPlotHoleRepo(pool), repositories.NewEntityRepo(pool), 8)

	holes, err := svc.Scan(ctx, universe.ID, ch10.ID)
	if err != nil {
		t.Fatalf("Scan failed: %v", err)
	}

	if len(holes) != 1 {
		t.Errorf("expected 1 plot hole, got %d", len(holes))
	}
}

// TestPlotHoleServiceScanNoLastMentioned skips entities without last_mentioned data.
func TestPlotHoleServiceScanNoLastMentioned(t *testing.T) {
	pool := testutil.SetupTestDB(t)
	testutil.RunMigrationsUpTo(t, pool, "011")
	ctx := context.Background()

	user := svcCreateTestUser(t, ctx, pool)
	universe := svcCreateTestUniverse(t, ctx, pool, user.ID)
	chapters := svcCreateChapters(t, ctx, pool, universe, 3)
	ch3 := chapters[2]

	// Entity with NULL last_mentioned — never mentioned
	svcCreateTestEntity(t, ctx, pool, universe.ID, "Never Mentioned", 0.2, "active")

	svc := NewPlotHoleService(pool, repositories.NewPlotHoleRepo(pool), repositories.NewEntityRepo(pool), 8)

	holes, err := svc.Scan(ctx, universe.ID, ch3.ID)
	if err != nil {
		t.Fatalf("Scan failed: %v", err)
	}

	if len(holes) != 0 {
		t.Errorf("expected 0 plot holes (null last_mentioned should be skipped), got %d", len(holes))
	}
}

// svcCreateChapters creates a work + N chapters and returns the chapters.
func svcCreateChapters(t *testing.T, ctx context.Context, pool *pgxpool.Pool, universe models.Universe, n int) []models.Chapter {
	t.Helper()
	w := models.Work{ID: uuid.New(), UniverseID: universe.ID, Title: "Test Work", Type: "book", OrderIndex: 1, Status: "draft"}
	if _, err := pool.Exec(ctx, "INSERT INTO works (id, universe_id, title, type, order_index, status) VALUES ($1,$2,$3,$4,$5,$6)",
		w.ID, w.UniverseID, w.Title, w.Type, w.OrderIndex, w.Status); err != nil {
		t.Fatalf("create work: %v", err)
	}

	chapters := make([]models.Chapter, n)
	for i := 0; i < n; i++ {
		ch := models.Chapter{ID: uuid.New(), WorkID: w.ID, Title: "Ch", OrderIndex: i + 1, Status: "draft"}
		if _, err := pool.Exec(ctx, "INSERT INTO chapters (id, work_id, title, order_index, content, raw_text, word_count, status) VALUES ($1,$2,$3,$4,'','',0,$5)",
			ch.ID, ch.WorkID, ch.Title, ch.OrderIndex, ch.Status); err != nil {
			t.Fatalf("create chapter: %v", err)
		}
		chapters[i] = ch
	}
	return chapters
}
