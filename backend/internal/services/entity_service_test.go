package services

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/quill/backend/internal/config"
	"github.com/quill/backend/internal/models"
	"github.com/quill/backend/internal/repositories"
	"github.com/quill/backend/internal/testutil"
)

// TestResolveOrCreateNewEntityAppendsHistoryRow proves entity creation (Step
// 4 of ResolveOrCreate, brand-new entity, isNew=true) writes a single
// entity_relevance_history row with the initial score (0.8) and status
// (spec: Relevance history persistence requirement).
func TestResolveOrCreateNewEntityAppendsHistoryRow(t *testing.T) {
	pool := testutil.SetupTestDB(t)
	testutil.RunMigrationsUpTo(t, pool, "019")
	ctx := context.Background()

	// Mock Qwen embeddings endpoint so Step 3 (semantic similarity) doesn't
	// need a real API key.
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := EmbeddingResponse{Data: []struct {
			Embedding []float32 `json:"embedding"`
			Index     int       `json:"index"`
		}{{Embedding: make([]float32, 1024)}}}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	cfg := &config.Config{QwenBaseURL: server.URL, QwenAPIKey: "test-key"}
	qwenSvc := NewQwenService(cfg, nil)

	user := svcCreateTestUser(t, ctx, pool)
	universe := svcCreateTestUniverse(t, ctx, pool, user.ID)

	entityRepo := repositories.NewEntityRepo(pool)
	vectorRepo := repositories.NewVectorRepo(pool)
	entitySvc := NewEntityService(pool, entityRepo, vectorRepo, qwenSvc)

	entity, previousStatus, isNew, err := entitySvc.ResolveOrCreate(ctx, universe.ID, repositories.ExtractedEntity{
		Type: "character", Name: "Brand New Wizard",
	})
	if err != nil {
		t.Fatalf("ResolveOrCreate: %v", err)
	}
	if !isNew {
		t.Fatal("expected a brand-new entity to be created")
	}
	if previousStatus != "" {
		t.Errorf("previousStatus = %q, want empty for brand-new entity", previousStatus)
	}

	historyRepo := repositories.NewEntityRelevanceHistoryRepo(pool)
	points, err := historyRepo.ListRecentByUniverse(ctx, universe.ID, 30)
	if err != nil {
		t.Fatalf("ListRecentByUniverse: %v", err)
	}
	if len(points) != 1 {
		t.Fatalf("len(points) = %d, want 1", len(points))
	}
	if points[0].EntityID != entity.ID {
		t.Errorf("EntityID = %v, want %v", points[0].EntityID, entity.ID)
	}
	if points[0].RelevanceScore != 0.8 {
		t.Errorf("RelevanceScore = %f, want 0.8", points[0].RelevanceScore)
	}
}

// newErrorQwenService returns a QwenService whose embedding endpoint always
// fails, so ResolveOrCreate falls through to creating a new entity.
func newErrorQwenService(t *testing.T) *QwenService {
	t.Helper()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	t.Cleanup(server.Close)
	cfg := &config.Config{QwenBaseURL: server.URL, QwenAPIKey: "test-key"}
	return NewQwenService(cfg, nil)
}

// seedEntity creates an entity row directly via the repository so tests can
// exercise ResolveOrCreate paths without a real Qwen embedding service.
func seedEntity(t *testing.T, ctx context.Context, pool *pgxpool.Pool, universeID uuid.UUID, name, entityType string) *models.Entity {
	t.Helper()
	repo := repositories.NewEntityRepo(pool)
	tx, err := pool.Begin(ctx)
	if err != nil {
		t.Fatalf("begin tx: %v", err)
	}
	defer tx.Rollback(ctx)

	e := &models.Entity{
		ID:             uuid.New(),
		UniverseID:     universeID,
		Type:           entityType,
		Name:           name,
		Status:         "active",
		RelevanceScore: 0.8,
	}
	if err := repo.Create(ctx, tx, e); err != nil {
		t.Fatalf("create entity: %v", err)
	}
	if err := tx.Commit(ctx); err != nil {
		t.Fatalf("commit: %v", err)
	}
	return e
}

// TestResolveOrCreateFuzzyMergeShortQuery proves a short incoming name ("Holden")
// merges into an existing longer name ("James Holden") via fuzzy substring match.
func TestResolveOrCreateFuzzyMergeShortQuery(t *testing.T) {
	pool := testutil.SetupTestDB(t)
	testutil.RunMigrationsUpTo(t, pool, "005")
	ctx := context.Background()

	user := svcCreateTestUser(t, ctx, pool)
	universe := svcCreateTestUniverse(t, ctx, pool, user.ID)
	entityRepo := repositories.NewEntityRepo(pool)
	entitySvc := NewEntityService(pool, entityRepo, nil, nil)

	existing := seedEntity(t, ctx, pool, universe.ID, "James Holden", "character")

	merged, prevStatus, isNew, err := entitySvc.ResolveOrCreate(ctx, universe.ID, repositories.ExtractedEntity{
		Type: "character", Name: "Holden",
	})
	if err != nil {
		t.Fatalf("ResolveOrCreate: %v", err)
	}
	if isNew {
		t.Fatal("expected fuzzy merge, not a new entity")
	}
	if merged.ID != existing.ID {
		t.Errorf("merged.ID = %v, want %v", merged.ID, existing.ID)
	}
	if prevStatus != existing.Status {
		t.Errorf("prevStatus = %q, want %q", prevStatus, existing.Status)
	}

	list, total, err := entityRepo.ListByUniverse(ctx, universe.ID, repositories.EntityFilters{Page: 1, Limit: 100})
	if err != nil {
		t.Fatalf("ListByUniverse: %v", err)
	}
	if total != 1 || len(list) != 1 {
		t.Errorf("want 1 entity, got total=%d len=%d", total, len(list))
	}
}

// TestResolveOrCreateFuzzyMergeLongQuery proves a long incoming name ("James Holden")
// merges into an existing short name ("Holden") via fuzzy substring match.
func TestResolveOrCreateFuzzyMergeLongQuery(t *testing.T) {
	pool := testutil.SetupTestDB(t)
	testutil.RunMigrationsUpTo(t, pool, "005")
	ctx := context.Background()

	user := svcCreateTestUser(t, ctx, pool)
	universe := svcCreateTestUniverse(t, ctx, pool, user.ID)
	entityRepo := repositories.NewEntityRepo(pool)
	entitySvc := NewEntityService(pool, entityRepo, nil, nil)

	existing := seedEntity(t, ctx, pool, universe.ID, "Holden", "character")

	merged, _, isNew, err := entitySvc.ResolveOrCreate(ctx, universe.ID, repositories.ExtractedEntity{
		Type: "character", Name: "James Holden",
	})
	if err != nil {
		t.Fatalf("ResolveOrCreate: %v", err)
	}
	if isNew {
		t.Fatal("expected fuzzy merge, not a new entity")
	}
	if merged.ID != existing.ID {
		t.Errorf("merged.ID = %v, want %v", merged.ID, existing.ID)
	}

	list, total, err := entityRepo.ListByUniverse(ctx, universe.ID, repositories.EntityFilters{Page: 1, Limit: 100})
	if err != nil {
		t.Fatalf("ListByUniverse: %v", err)
	}
	if total != 1 || len(list) != 1 {
		t.Errorf("want 1 entity, got total=%d len=%d", total, len(list))
	}
}

// TestResolveOrCreateFuzzyMergeRespectsType proves fuzzy matching does not merge
// across entity types.
func TestResolveOrCreateFuzzyMergeRespectsType(t *testing.T) {
	pool := testutil.SetupTestDB(t)
	testutil.RunMigrationsUpTo(t, pool, "005")
	ctx := context.Background()

	user := svcCreateTestUser(t, ctx, pool)
	universe := svcCreateTestUniverse(t, ctx, pool, user.ID)
	entityRepo := repositories.NewEntityRepo(pool)
	entitySvc := NewEntityService(pool, entityRepo, nil, newErrorQwenService(t))

	_ = seedEntity(t, ctx, pool, universe.ID, "James Holden", "character")

	// Substring match but a different type should create a new entity.
	_, _, isNew, err := entitySvc.ResolveOrCreate(ctx, universe.ID, repositories.ExtractedEntity{
		Type: "location", Name: "Holden",
	})
	if err != nil {
		t.Fatalf("ResolveOrCreate: %v", err)
	}
	if !isNew {
		t.Fatal("expected a new entity for a different type")
	}

	list, total, err := entityRepo.ListByUniverse(ctx, universe.ID, repositories.EntityFilters{Page: 1, Limit: 100})
	if err != nil {
		t.Fatalf("ListByUniverse: %v", err)
	}
	if total != 2 || len(list) != 2 {
		t.Errorf("want 2 entities, got total=%d len=%d", total, len(list))
	}
}
