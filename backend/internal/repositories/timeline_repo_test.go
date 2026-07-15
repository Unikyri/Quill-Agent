package repositories

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/quill/backend/internal/models"
	"github.com/quill/backend/internal/testutil"
)

func setupTimelineFixtures(t *testing.T, pool *pgxpool.Pool) (models.Universe, models.Chapter, models.Entity) {
	t.Helper()
	ctx := context.Background()
	user := createTestUser(t, ctx, pool)

	universe := models.Universe{ID: uuid.New(), UserID: user.ID, Name: "Timeline Universe", GenreTags: []string{"fantasy"}}
	work := models.Work{ID: uuid.New(), UniverseID: universe.ID, Title: "TL Work", Type: "novel", OrderIndex: 1, Status: "in_progress"}
	chapter := models.Chapter{ID: uuid.New(), WorkID: work.ID, Title: "Ch1", OrderIndex: 5, Content: "content", RawText: "text", WordCount: 10, Status: "draft"}
	entity := models.Entity{ID: uuid.New(), UniverseID: universe.ID, Type: "character", Name: "TL Entity", Status: "active", RelevanceScore: 0.8}

	tx, err := pool.Begin(ctx)
	if err != nil {
		t.Fatalf("begin tx: %v", err)
	}
	defer tx.Rollback(ctx)

	if _, err := tx.Exec(ctx, "INSERT INTO universes (id, user_id, name, format) VALUES ($1,$2,$3,$4)",
		universe.ID, universe.UserID, universe.Name, "novel"); err != nil {
		t.Fatalf("insert universe: %v", err)
	}
	if _, err := tx.Exec(ctx, "INSERT INTO works (id, universe_id, title, type, order_index, status) VALUES ($1,$2,$3,$4,$5,$6)",
		work.ID, work.UniverseID, work.Title, work.Type, work.OrderIndex, work.Status); err != nil {
		t.Fatalf("insert work: %v", err)
	}
	if _, err := tx.Exec(ctx, "INSERT INTO chapters (id, work_id, title, order_index, content, raw_text, word_count, status) VALUES ($1,$2,$3,$4,$5,$6,$7,$8)",
		chapter.ID, chapter.WorkID, chapter.Title, chapter.OrderIndex, chapter.Content, chapter.RawText, chapter.WordCount, chapter.Status); err != nil {
		t.Fatalf("insert chapter: %v", err)
	}
	if _, err := tx.Exec(ctx, "INSERT INTO entities (id, universe_id, type, name, description, status, relevance_score) VALUES ($1,$2,$3,$4,$5,$6,$7)",
		entity.ID, entity.UniverseID, entity.Type, entity.Name, "", entity.Status, entity.RelevanceScore); err != nil {
		t.Fatalf("insert entity: %v", err)
	}
	if err := tx.Commit(ctx); err != nil {
		t.Fatalf("commit: %v", err)
	}
	return universe, chapter, entity
}

func TestTimelineRepoCreate(t *testing.T) {
	pool := testutil.SetupTestDB(t)
	testutil.RunMigrationsUpTo(t, pool, "012")
	universe, chapter, entity := setupTimelineFixtures(t, pool)

	ctx := context.Background()
	repo := NewTimelineRepo(pool)

	pos := 3.5
	evt := &models.TimelineEvent{
		ID:               uuid.New(),
		UniverseID:       universe.ID,
		EventEntityID:    &entity.ID,
		Title:            "Battle of Test",
		Description:      "A great battle",
		TimelinePosition: &pos,
		TimelineLabel:    "Chapter 3",
		ChapterID:        &chapter.ID,
		Participants:     []uuid.UUID{entity.ID},
	}
	if err := repo.Create(ctx, evt); err != nil {
		t.Fatalf("Create: %v", err)
	}

	found, err := repo.GetByID(ctx, evt.ID)
	if err != nil {
		t.Fatalf("GetByID: %v", err)
	}
	if found.Title != evt.Title {
		t.Errorf("Title = %q, want %q", found.Title, evt.Title)
	}
}

func TestTimelineRepoListByUniverseOrdered(t *testing.T) {
	pool := testutil.SetupTestDB(t)
	testutil.RunMigrationsUpTo(t, pool, "012")
	universe, _, entity := setupTimelineFixtures(t, pool)

	ctx := context.Background()
	repo := NewTimelineRepo(pool)

	pos1 := 1.0
	pos3 := 3.0
	pos2 := 2.0

	events := []*models.TimelineEvent{
		{ID: uuid.New(), UniverseID: universe.ID, EventEntityID: &entity.ID, Title: "Event A", TimelinePosition: &pos3},
		{ID: uuid.New(), UniverseID: universe.ID, EventEntityID: &entity.ID, Title: "Event B", TimelinePosition: &pos1},
		{ID: uuid.New(), UniverseID: universe.ID, EventEntityID: &entity.ID, Title: "Event C", TimelinePosition: &pos2},
	}
	for _, e := range events {
		if err := repo.Create(ctx, e); err != nil {
			t.Fatalf("Create %s: %v", e.Title, err)
		}
	}

	list, err := repo.ListByUniverse(ctx, universe.ID)
	if err != nil {
		t.Fatalf("ListByUniverse: %v", err)
	}
	if len(list) != 3 {
		t.Fatalf("len = %d, want 3", len(list))
	}
	// Must be ordered by timeline_position ascending
	if list[0].Title != "Event B" { // pos=1.0
		t.Errorf("first = %q, want Event B (pos 1.0)", list[0].Title)
	}
	if list[1].Title != "Event C" { // pos=2.0
		t.Errorf("second = %q, want Event C (pos 2.0)", list[1].Title)
	}
	if list[2].Title != "Event A" { // pos=3.0
		t.Errorf("third = %q, want Event A (pos 3.0)", list[2].Title)
	}
}

func TestTimelineRepoGetByIDNotFound(t *testing.T) {
	pool := testutil.SetupTestDB(t)
	testutil.RunMigrationsUpTo(t, pool, "012")

	ctx := context.Background()
	repo := NewTimelineRepo(pool)

	_, err := repo.GetByID(ctx, uuid.New())
	if err == nil {
		t.Error("expected error for nonexistent timeline event")
	}
}
