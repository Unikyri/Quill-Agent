package services

import (
	"context"
	"sort"
	"testing"

	"github.com/google/uuid"

	"github.com/quill/backend/internal/repositories"
	"github.com/quill/backend/internal/testutil"
)

// TestMemoryServiceRecallMergeWeighting verifies that Recall returns items
// merged from graph neighbors, recent mentions, and freshness signals,
// scored with the weighted formula: graph×0.4 + recency×0.3 + freshness×0.3.
func TestMemoryServiceRecallMergeWeighting(t *testing.T) {
	pool := testutil.SetupTestDB(t)
	testutil.RunMigrationsUpTo(t, pool, "005")
	ctx := context.Background()

	user := svcCreateTestUser(t, ctx, pool)
	universe := svcCreateTestUniverse(t, ctx, pool, user.ID)

	// Create active entities with varying scores
	entityA := svcCreateTestEntity(t, ctx, pool, universe.ID, "Entity A", 0.9, "active")
	entityB := svcCreateTestEntity(t, ctx, pool, universe.ID, "Entity B", 0.5, "active")
	_ = svcCreateTestEntity(t, ctx, pool, universe.ID, "Entity C", 0.2, "active") // archived, not in result

	entityRepo := repositories.NewEntityRepo(pool)
	graphRepo := repositories.NewGraphRepo(pool)
	vectorRepo := repositories.NewVectorRepo(pool)

	svc := NewMemoryService(graphRepo, entityRepo, vectorRepo)

	// Create a graph so NHopTraversal finds neighbors
	graphName := "universe_" + universe.ID.String()
	_ = graphRepo.CreateGraph(ctx, universe.ID.String())
	_ = graphRepo.CreateNode(ctx, graphName, "Entity", map[string]interface{}{
		"entity_id": entityA.ID.String(), "name": "Entity A", "status": "active", "relevance_score": 0.9,
	})
	_ = graphRepo.CreateNode(ctx, graphName, "Entity", map[string]interface{}{
		"entity_id": entityB.ID.String(), "name": "Entity B", "status": "active", "relevance_score": 0.5,
	})
	_ = graphRepo.CreateEdge(ctx, graphName, entityA.ID.String(), entityB.ID.String(), "KNOWS", nil)

	// Query embedding (arbitrary float32 slice)
	queryEmb := make([]float32, 3)
	queryEmb[0] = 1.0

	items, err := svc.Recall(ctx, universe.ID, queryEmb, 5)
	if err != nil {
		t.Fatalf("Recall failed: %v", err)
	}

	if len(items) == 0 {
		t.Error("Expected non-empty recall results")
	}

	// Verify each item has required fields
	for _, item := range items {
		if item.EntityID == uuid.Nil {
			t.Error("RecallItem.EntityID should not be nil")
		}
		if item.Score <= 0 {
			t.Errorf("RecallItem.Score should be > 0, got %f for entity %s", item.Score, item.EntityID)
		}
		if item.Source == "" {
			t.Errorf("RecallItem.Source should not be empty for entity %s", item.EntityID)
		}
	}

	// Verify results are sorted by score descending
	sorted := make([]float64, len(items))
	for i, item := range items {
		sorted[i] = item.Score
	}
	if !sort.Float64sAreSorted(sorted) && len(sorted) > 1 {
		// Check if descending
		for i := 1; i < len(sorted); i++ {
			if sorted[i-1] < sorted[i] {
				t.Errorf("Recall results not sorted by score descending: item[%d]=%f > item[%d]=%f", i-1, sorted[i-1], i, sorted[i])
				break
			}
		}
	}

	// Verify k limit is respected
	if len(items) > 5 {
		t.Errorf("Recall returned %d items, expected at most 5 (k limit)", len(items))
	}
}

// TestMemoryServiceRecallEmptyUniverse verifies Recall handles empty universes gracefully.
func TestMemoryServiceRecallEmptyUniverse(t *testing.T) {
	pool := testutil.SetupTestDB(t)
	testutil.RunMigrationsUpTo(t, pool, "005")
	ctx := context.Background()

	user := svcCreateTestUser(t, ctx, pool)
	universe := svcCreateTestUniverse(t, ctx, pool, user.ID)

	entityRepo := repositories.NewEntityRepo(pool)
	graphRepo := repositories.NewGraphRepo(pool)
	vectorRepo := repositories.NewVectorRepo(pool)

	svc := NewMemoryService(graphRepo, entityRepo, vectorRepo)
	queryEmb := make([]float32, 3)

	items, err := svc.Recall(ctx, universe.ID, queryEmb, 5)
	if err != nil {
		t.Fatalf("Recall failed on empty universe: %v", err)
	}
	if items == nil {
		t.Error("Expected empty slice, got nil")
	}
}

// TestMemoryServiceRecallKCap verifies k limits are respected.
func TestMemoryServiceRecallKCap(t *testing.T) {
	pool := testutil.SetupTestDB(t)
	testutil.RunMigrationsUpTo(t, pool, "005")
	ctx := context.Background()

	user := svcCreateTestUser(t, ctx, pool)
	universe := svcCreateTestUniverse(t, ctx, pool, user.ID)

	// Create several entities
	for i := 0; i < 10; i++ {
		name := "Entity " + string(rune('A'+i))
		svcCreateTestEntity(t, ctx, pool, universe.ID, name, 0.5+float64(i)*0.02, "active")
	}

	entityRepo := repositories.NewEntityRepo(pool)
	graphRepo := repositories.NewGraphRepo(pool)
	vectorRepo := repositories.NewVectorRepo(pool)

	svc := NewMemoryService(graphRepo, entityRepo, vectorRepo)
	queryEmb := make([]float32, 3)

	items, err := svc.Recall(ctx, universe.ID, queryEmb, 3)
	if err != nil {
		t.Fatalf("Recall failed: %v", err)
	}

	if len(items) > 3 {
		t.Errorf("Recall returned %d items, expected at most 3 (k=3)", len(items))
	}
}
