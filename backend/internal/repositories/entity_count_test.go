package repositories

import (
	"context"
	"testing"

	"github.com/google/uuid"

	"github.com/quill/backend/internal/testutil"
)

func TestEntityRepoCountByTypeIncludesArchivedEntities(t *testing.T) {
	pool := testutil.SetupTestDB(t)
	testutil.RunMigrationsUpTo(t, pool, "021")
	ctx := context.Background()

	userID := uuid.New()
	universeID := uuid.New()
	if _, err := pool.Exec(ctx, "INSERT INTO users (id, email, password_hash, display_name) VALUES ($1, $2, $3, $4)", userID, "count@example.com", "hash", "Count User"); err != nil {
		t.Fatalf("insert user: %v", err)
	}
	if _, err := pool.Exec(ctx, "INSERT INTO universes (id, user_id, name, genre_tags) VALUES ($1, $2, $3, $4)", universeID, userID, "Count Universe", []string{"fantasy"}); err != nil {
		t.Fatalf("insert universe: %v", err)
	}

	fixtures := []struct {
		entityType string
		status     string
	}{
		{entityType: "character", status: "active"},
		{entityType: "character", status: "archived"},
		{entityType: "object", status: "active"},
	}
	for index, fixture := range fixtures {
		if _, err := pool.Exec(ctx, `
			INSERT INTO entities (id, universe_id, type, name, status, relevance_score)
			VALUES ($1, $2, $3, $4, $5, $6)
		`, uuid.New(), universeID, fixture.entityType, "Entity "+string(rune('A'+index)), fixture.status, 0.5); err != nil {
			t.Fatalf("insert entity: %v", err)
		}
	}

	counts, err := NewEntityRepo(pool).CountByType(ctx, universeID)
	if err != nil {
		t.Fatalf("CountByType: %v", err)
	}
	if counts["character"] != 2 {
		t.Errorf("character count = %d, want 2", counts["character"])
	}
	if counts["object"] != 1 {
		t.Errorf("object count = %d, want 1", counts["object"])
	}
	if _, exists := counts["place"]; exists {
		t.Errorf("place should not be returned when its count is zero: %v", counts)
	}
}

func TestEntityTypeCheckRejectsNonCanonicalValue(t *testing.T) {
	pool := testutil.SetupTestDB(t)
	testutil.RunMigrationsUpTo(t, pool, "021")
	ctx := context.Background()

	userID := uuid.New()
	universeID := uuid.New()
	if _, err := pool.Exec(ctx, "INSERT INTO users (id, email, password_hash, display_name) VALUES ($1, $2, $3, $4)", userID, "taxonomy@example.com", "hash", "Taxonomy User"); err != nil {
		t.Fatalf("insert user: %v", err)
	}
	if _, err := pool.Exec(ctx, "INSERT INTO universes (id, user_id, name, genre_tags) VALUES ($1, $2, $3, $4)", universeID, userID, "Taxonomy Universe", []string{"fantasy"}); err != nil {
		t.Fatalf("insert universe: %v", err)
	}

	_, err := pool.Exec(ctx, `
		INSERT INTO entities (id, universe_id, type, name, status, relevance_score)
		VALUES ($1, $2, $3, $4, $5, $6)
	`, uuid.New(), universeID, "invalid_type", "Invalid Entity", "active", 0.5)
	if err == nil {
		t.Fatal("expected entities_type_check to reject a non-canonical type")
	}
}
