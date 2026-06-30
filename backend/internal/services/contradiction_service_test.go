package services

import (
	"context"
	"testing"

	"github.com/google/uuid"

	"github.com/quill/backend/internal/config"
	"github.com/quill/backend/internal/models"
	"github.com/quill/backend/internal/repositories"
	"github.com/quill/backend/internal/testutil"
)

// TestContradictionFingerprintDeterminism verifies SHA-256 fingerprint is
// deterministic — same inputs always produce the same hash. Pure unit test.
func TestContradictionFingerprintDeterminism(t *testing.T) {
	entityA := uuid.New()
	entityB := uuid.New()

	candidates := []ContradictionCandidate{
		{EntityID: entityA, Type: "deceased_alive", EvidenceA: "Bob alive ch1", EvidenceB: "Bob dead ch3"},
		{EntityID: entityB, Type: "status_change", EvidenceA: "Alice mayor", EvidenceB: "Alice queen"},
	}

	// Create service with nil dependencies — fingerprint is pure, needs no DB
	svc := NewContradictionService(nil, nil, nil, nil, 3)

	fp1 := svc.fingerprint(candidates[0])
	fp2 := svc.fingerprint(candidates[0])

	if fp1 == "" {
		t.Error("fingerprint should not be empty")
	}
	if fp1 != fp2 {
		t.Errorf("same input produced different fingerprints: %s vs %s", fp1, fp2)
	}

	// Different inputs should produce different fingerprints
	fp3 := svc.fingerprint(candidates[1])
	if fp3 == "" {
		t.Error("fingerprint for candidate[1] should not be empty")
	}
	if fp3 == fp1 {
		t.Error("different inputs should produce different fingerprints")
	}
}

// TestContradictionFingerprintChaptersIncluded verifies that ChapterA/ChapterB
// affect the fingerprint — two candidates identical except for chapter fields
// must produce different fingerprints.
func TestContradictionFingerprintChaptersIncluded(t *testing.T) {
	svc := NewContradictionService(nil, nil, nil, nil, 3)
	entityID := uuid.New()
	chA := uuid.New()
	chB := uuid.New()
	chOther := uuid.New()

	c1 := ContradictionCandidate{
		EntityID:  entityID,
		Type:      "semantic",
		EvidenceA: "evidence A",
		EvidenceB: "evidence B",
		ChapterA:  chA,
		ChapterB:  chB,
	}
	c2 := ContradictionCandidate{
		EntityID:  entityID,
		Type:      "semantic",
		EvidenceA: "evidence A",
		EvidenceB: "evidence B",
		ChapterA:  chOther, // different chapter
		ChapterB:  chB,
	}

	fp1 := svc.fingerprint(c1)
	fp2 := svc.fingerprint(c2)

	if fp1 == fp2 {
		t.Error("fingerprints should differ when ChapterA differs — chapter fields must be included in hash")
	}

	// Triangulate: ChapterB differs
	c3 := ContradictionCandidate{
		EntityID:  entityID,
		Type:      "semantic",
		EvidenceA: "evidence A",
		EvidenceB: "evidence B",
		ChapterA:  chA,
		ChapterB:  chOther,
	}
	fp3 := svc.fingerprint(c3)
	if fp1 == fp3 {
		t.Error("fingerprints should differ when ChapterB differs — chapter fields must be included in hash")
	}
	if fp3 == fp2 {
		t.Error("fingerprints with different chapters should all differ")
	}
}

// TestContradictionFingerprintFormat verifies the fingerprint is a valid hex string
// (64 chars for SHA-256).
func TestContradictionFingerprintFormat(t *testing.T) {
	cfg := &config.Config{MaxContradictionCandidates: 3}
	svc := NewContradictionService(nil, nil, nil, nil, cfg.MaxContradictionCandidates)

	c := ContradictionCandidate{
		EntityID:  uuid.New(),
		Type:      "semantic",
		EvidenceA: "test evidence A",
		EvidenceB: "test evidence B",
	}

	fp := svc.fingerprint(c)
	if len(fp) != 64 {
		t.Errorf("SHA-256 fingerprint length = %d, want 64 hex chars", len(fp))
	}
	// Verify all characters are hex
	for _, ch := range fp {
		if !((ch >= '0' && ch <= '9') || (ch >= 'a' && ch <= 'f')) {
			t.Errorf("fingerprint contains non-hex character: %c", ch)
			break
		}
	}
}

// TestContradictionCheckDeterministicDeceasedAlive verifies that the
// deterministic rule catches deceased/alive contradictions without calling Qwen API.
func TestContradictionCheckDeterministicDeceasedAlive(t *testing.T) {
	pool := testutil.SetupTestDB(t)
	testutil.RunMigrationsUpTo(t, pool, "005")
	ctx := context.Background()

	user := svcCreateTestUser(t, ctx, pool)
	universe := svcCreateTestUniverse(t, ctx, pool, user.ID)

	// Create a "deceased" entity
	deceasedID := uuid.New()
	// Insert entity with "deceased" status via pool directly
	if _, err := pool.Exec(ctx,
		"INSERT INTO entities (id, universe_id, type, name, description, status, relevance_score) VALUES ($1,$2,'character','Dead Bob','','deceased',0.8)",
		deceasedID, universe.ID); err != nil {
		t.Fatalf("create deceased entity: %v", err)
	}

	entityRepo := repositories.NewEntityRepo(pool)
	contraRepo := repositories.NewContradictionRepo(pool)
	cfg := config.Config{MaxContradictionCandidates: 3}
	svc := NewContradictionService(pool, contraRepo, entityRepo, nil, cfg.MaxContradictionCandidates) // nil qwenSvc

	// Pass an entity marked as "alive" but the DB says "deceased"
	entities := []ResolvedEntity{
		{
			Entity:     models.Entity{ID: deceasedID, UniverseID: universe.ID, Type: "character", Name: "Dead Bob", Status: "deceased"},
			MentionText: "Bob walked into the room",
			IsNew:       false,
		},
	}

	chapterID := uuid.New()
	contradictions, err := svc.CheckDeterministic(ctx, universe.ID, chapterID, entities)
	if err != nil {
		t.Fatalf("CheckDeterministic: %v", err)
	}

	// Should detect that "Dead Bob" is mentioned as alive but DB says deceased
	if len(contradictions) == 0 {
		t.Error("Expected at least one contradiction for deceased entity mentioned as alive")
	}
	if len(contradictions) > 0 && contradictions[0].Severity == "" {
		t.Error("Contradiction severity should not be empty")
	}
	if len(contradictions) > 0 && contradictions[0].Severity != "critical" {
		t.Errorf("Contradiction severity for deceased_alive should be 'critical', got '%s'", contradictions[0].Severity)
	}
}

// TestContradictionCheckDeterministicNoIssues verifies deterministic check
// returns empty when there are no issues.
func TestContradictionCheckDeterministicNoIssues(t *testing.T) {
	pool := testutil.SetupTestDB(t)
	testutil.RunMigrationsUpTo(t, pool, "005")
	ctx := context.Background()

	user := svcCreateTestUser(t, ctx, pool)
	universe := svcCreateTestUniverse(t, ctx, pool, user.ID)

	activeEntity := svcCreateTestEntity(t, ctx, pool, universe.ID, "Alive Alice", 0.8, "active")

	entityRepo := repositories.NewEntityRepo(pool)
	contraRepo := repositories.NewContradictionRepo(pool)
	cfg := config.Config{MaxContradictionCandidates: 3}
	svc := NewContradictionService(pool, contraRepo, entityRepo, nil, cfg.MaxContradictionCandidates)

	entities := []ResolvedEntity{
		{Entity: activeEntity, MentionText: "Alice walked to the store", IsNew: false},
	}

	chapterID := uuid.New()
	contradictions, err := svc.CheckDeterministic(ctx, universe.ID, chapterID, entities)
	if err != nil {
		t.Fatalf("CheckDeterministic: %v", err)
	}

	if len(contradictions) != 0 {
		t.Errorf("Expected 0 contradictions for active entity, got %d", len(contradictions))
	}
}

// TestContradictionCheckDeterministicChapterThreaded verifies that ChapterA/ChapterB
// are populated on ContradictionCandidate structs when CheckDeterministic is called
// with a chapterID. This ensures the fingerprint embeds chapter context.
//
// RED: CheckDeterministic currently takes 3 params (ctx, universeID, entities).
// This test adds a 4th param (chapterID) — won't compile until production code
// is updated.
func TestContradictionCheckDeterministicChapterThreaded(t *testing.T) {
	pool := testutil.SetupTestDB(t)
	testutil.RunMigrationsUpTo(t, pool, "005")
	ctx := context.Background()

	user := svcCreateTestUser(t, ctx, pool)
	universe := svcCreateTestUniverse(t, ctx, pool, user.ID)

	deceasedID := uuid.New()
	if _, err := pool.Exec(ctx,
		"INSERT INTO entities (id, universe_id, type, name, description, status, relevance_score) VALUES ($1,$2,'character','Ghost Bob','','deceased',0.8)",
		deceasedID, universe.ID); err != nil {
		t.Fatalf("create deceased entity: %v", err)
	}

	entityRepo := repositories.NewEntityRepo(pool)
	contraRepo := repositories.NewContradictionRepo(pool)
	cfg := config.Config{MaxContradictionCandidates: 3}
	svc := NewContradictionService(pool, contraRepo, entityRepo, nil, cfg.MaxContradictionCandidates)

	chapterID := uuid.New()

	entities := []ResolvedEntity{
		{
			Entity:      models.Entity{ID: deceasedID, UniverseID: universe.ID, Type: "character", Name: "Ghost Bob", Status: "deceased"},
			MentionText: "Bob walked into the room",
			IsNew:       false,
		},
	}

	// RED: this line won't compile — CheckDeterministic currently only takes 3 params
	contradictions, err := svc.CheckDeterministic(ctx, universe.ID, chapterID, entities)
	if err != nil {
		t.Fatalf("CheckDeterministic: %v", err)
	}

	if len(contradictions) == 0 {
		t.Fatal("Expected at least one contradiction for deceased entity")
	}

	// Verify the fingerprint embeds the chapter — reconstruct expected fingerprint
	expectedFP := svc.fingerprint(ContradictionCandidate{
		EntityID:  deceasedID,
		Type:      "deceased_alive",
		EvidenceA: "Entity Ghost Bob is deceased in DB",
		EvidenceB: "Bob walked into the room",
		ChapterA:  chapterID,
		ChapterB:  chapterID,
	})

	if contradictions[0].Fingerprint != expectedFP {
		t.Errorf("Fingerprint mismatch — chapterID not threaded into candidate?\n  got:  %s\n  want: %s",
			contradictions[0].Fingerprint, expectedFP)
	}
}

// TestContradictionCheckSemanticSignature verifies the CheckSemantic method
// compiles and routes to QwenService.CheckContradictions.
func TestContradictionCheckSemanticSignature(t *testing.T) {
	pool := testutil.SetupTestDB(t)
	testutil.RunMigrationsUpTo(t, pool, "005")
	ctx := context.Background()

	contraRepo := repositories.NewContradictionRepo(pool)
	entityRepo := repositories.NewEntityRepo(pool)

	// Create QwenService with dummy config — CheckSemantic will call it
	cfgQwen := config.Config{
		QwenBaseURL:                "https://example.com",
		QwenAPIKey:                 "test-key",
		QwenMaxModel:               "qwen-max-latest",
		QwenMaxConcurrency:         1,
		QwenTurboConcurrency:       1,
		QwenEmbeddingModel:         "text-embedding-v3",
		MaxContradictionCandidates: 3,
	}
	qwenSvc := NewQwenService(&cfgQwen)

	cfg := config.Config{MaxContradictionCandidates: 3}
	svc := NewContradictionService(pool, contraRepo, entityRepo, qwenSvc, cfg.MaxContradictionCandidates)

	entities := []ResolvedEntity{
		{Entity: models.Entity{ID: uuid.New(), Type: "character", Name: "Test"}, MentionText: "Test text", IsNew: false},
	}

	// Should not panic, will try to call Qwen and fail gracefully
	_, err := svc.CheckSemantic(ctx, uuid.New(), uuid.New(), "some text", entities)
	// Expected to fail because Qwen API is unreachable — that's OK, test verifies compilation
	if err == nil {
		t.Log("Unexpected: CheckSemantic succeeded with dummy Qwen endpoint")
	}
}
