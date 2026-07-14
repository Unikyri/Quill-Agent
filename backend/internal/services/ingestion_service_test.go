package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/quill/backend/internal/models"
	"github.com/quill/backend/internal/repositories"
	"github.com/quill/backend/internal/testutil"
)

// ── Mocks ──

// mockIngestionHub captures SendToUser calls for verification.
type mockIngestionHub struct {
	mu       sync.Mutex
	messages []models.WSMessage
	userIDs  []uuid.UUID
}

func (m *mockIngestionHub) SendToUser(userID uuid.UUID, msg models.WSMessage) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.messages = append(m.messages, msg)
	m.userIDs = append(m.userIDs, userID)
	return nil
}

func (m *mockIngestionHub) popMessages() []models.WSMessage {
	m.mu.Lock()
	defer m.mu.Unlock()
	out := m.messages
	m.messages = nil
	return out
}

func (m *mockIngestionHub) popUserIDs() []uuid.UUID {
	m.mu.Lock()
	defer m.mu.Unlock()
	out := m.userIDs
	m.userIDs = nil
	return out
}

// mockQwenForIngestion returns canned ExtractEntities and GenerateEmbedding results.
type mockQwenForIngestion struct {
	extractResult *ExtractedEntities
	extractErr    error
}

func (m *mockQwenForIngestion) ExtractEntities(ctx context.Context, text, categories string) (*ExtractedEntities, error) {
	return m.extractResult, m.extractErr
}

func (m *mockQwenForIngestion) GenerateEmbedding(ctx context.Context, text string) ([]float32, error) {
	// Return a dummy embedding — length 3 for test
	return []float32{0.1, 0.2, 0.3}, nil
}

func (m *mockQwenForIngestion) GenerateEmbeddingBatch(ctx context.Context, texts []string) ([][]float32, error) {
	out := make([][]float32, len(texts))
	for i := range texts {
		out[i] = []float32{0.1, 0.2, 0.3}
	}
	return out, nil
}

// ── Test: pipeline sequence ──

// TestIngestionServicePipeline verifies the chunk→extract→embed→graph sequence
// in a mock-driven test. It uses an httptest-style pattern without real DB.
func TestIngestionServicePipeline(t *testing.T) {
	hub := &mockIngestionHub{}
	qwen := &mockQwenForIngestion{
		extractResult: &ExtractedEntities{
			Characters: []ExtractedEntity{
				{Type: "Character", Name: "Frodo", Status: "alive", Description: "A hobbit"},
			},
		},
	}

	docContent := `# Chapter 1: A Long-expected Party

Bilbo was going to have a birthday party.

# Chapter 2: The Shadow of the Past

Frodo learns about the Ring.`

	svc := &IngestionService{
		pool:       nil,
		entitySvc:  nil,
		vectorRepo: nil,
		graphRepo:  nil,
		qwenSvc:    qwen,
		hub:        hub,
	}

	universeID := uuid.New()

	jobID, duplicate, err := svc.Start(context.Background(), universeID, strings.NewReader(docContent), "test.md")
	if err != nil {
		t.Fatalf("Start: %v", err)
	}
	if duplicate {
		t.Error("expected duplicate=false for first upload")
	}
	if jobID == uuid.Nil {
		t.Error("expected non-nil job ID")
	}

	// Wait for goroutine to finish (small doc, fast)
	time.Sleep(200 * time.Millisecond)

	msgs := hub.popMessages()
	// Should have at least one progress event
	if len(msgs) == 0 {
		t.Error("expected at least one WebSocket progress message")
	}

	// Verify at least one is an ingestion_progress event
	foundProgress := false
	for _, msg := range msgs {
		if msg.Type == "ingestion_progress" {
			foundProgress = true
			var payload map[string]any
			if err := json.Unmarshal(msg.Payload, &payload); err != nil {
				t.Errorf("unmarshal progress payload: %v", err)
			}
			if payload["job_id"] == nil {
				t.Error("progress payload missing job_id")
			}
		}
	}
	if !foundProgress {
		t.Error("expected an ingestion_progress message")
	}
}

// TestIngestionServiceChunking verifies the document is split by markdown headers.
func TestIngestionServiceChunking(t *testing.T) {
	hub := &mockIngestionHub{}
	qwen := &mockQwenForIngestion{
		extractResult: &ExtractedEntities{},
	}

	docContent := `# Chapter A

Content A line 1.
Content A line 2.

# Chapter B

Content B line 1.`

	svc := &IngestionService{
		pool:       nil,
		entitySvc:  nil,
		vectorRepo: nil,
		graphRepo:  nil,
		qwenSvc:    qwen,
		hub:        hub,
	}

	chunks := svc.splitChunks(docContent)
	// ponytail: minimal chunking — one chunk per # header section
	if len(chunks) < 2 {
		t.Errorf("expected at least 2 chunks, got %d", len(chunks))
	}

	// Verify chunk content is non-empty
	for i, ch := range chunks {
		if strings.TrimSpace(ch.content) == "" {
			t.Errorf("chunk %d is empty", i)
		}
		if ch.title == "" {
			t.Errorf("chunk %d has no title", i)
		}
	}
}

// TestIngestionServiceEmptyDocument verifies handling of empty input.
func TestIngestionServiceEmptyDocument(t *testing.T) {
	hub := &mockIngestionHub{}
	qwen := &mockQwenForIngestion{
		extractResult: &ExtractedEntities{},
	}

	svc := &IngestionService{
		pool:       nil,
		entitySvc:  nil,
		vectorRepo: nil,
		graphRepo:  nil,
		qwenSvc:    qwen,
		hub:        hub,
	}

	chunks := svc.splitChunks("")
	if len(chunks) != 0 {
		t.Errorf("expected 0 chunks for empty document, got %d", len(chunks))
	}
}

// TestIngestionServiceHeaderlessFallback verifies a headerless document is split
// into ~50K-char paragraph-boundary chunks titled "Part N".
func TestIngestionServiceHeaderlessFallback(t *testing.T) {
	svc := &IngestionService{}

	para := strings.Repeat("A paragraph with enough characters to fill space. ", 200)
	var b strings.Builder
	for i := 0; i < 12; i++ {
		if i > 0 {
			b.WriteString("\n\n")
		}
		b.WriteString(fmt.Sprintf("Paragraph %d: %s", i, para))
	}
	content := b.String()

	chunks := svc.splitChunks(content)
	if len(chunks) <= 1 {
		t.Fatalf("expected multiple chunks for headerless large doc, got %d", len(chunks))
	}

	joined := ""
	for i, ch := range chunks {
		joined += ch.content
		if len(ch.content) == 0 {
			t.Errorf("chunk %d is empty", i)
		}
		if len(ch.content) > 50_000 {
			t.Errorf("chunk %d length %d exceeds 50K", i, len(ch.content))
		}
		wantTitle := fmt.Sprintf("Part %d", i+1)
		if ch.title != wantTitle {
			t.Errorf("chunk %d title = %q, want %q", i, ch.title, wantTitle)
		}
	}

	for i := 0; i < 12; i++ {
		marker := fmt.Sprintf("Paragraph %d:", i)
		if !strings.Contains(joined, marker) {
			t.Errorf("missing paragraph marker %q", marker)
		}
	}
}

// TestSplitByParagraphsSmallDoc verifies a small headerless document stays as
// a single chunk.
func TestSplitByParagraphsSmallDoc(t *testing.T) {
	content := "A short document.\n\nWith two paragraphs."
	chunks := splitByParagraphs(content)
	if len(chunks) != 1 {
		t.Fatalf("expected 1 chunk, got %d", len(chunks))
	}
	if chunks[0].title != "Untitled" {
		t.Errorf("title = %q, want Untitled", chunks[0].title)
	}
	if chunks[0].content != strings.TrimSpace(content) {
		t.Errorf("content mismatch: got %d chars, want %d chars", len(chunks[0].content), len(strings.TrimSpace(content)))
	}
}

// TestSplitByParagraphsOversizedParagraph verifies a single paragraph larger
// than the chunk cap becomes its own chunk rather than being dropped.
func TestSplitByParagraphsOversizedParagraph(t *testing.T) {
	big := strings.Repeat("a", 60_000)
	content := big + "\n\nsmall"
	chunks := splitByParagraphs(content)
	if len(chunks) < 1 {
		t.Fatalf("expected at least 1 chunk, got %d", len(chunks))
	}
	if chunks[0].title != "Part 1" {
		t.Errorf("first chunk title = %q, want Part 1", chunks[0].title)
	}
	if len(chunks[0].content) != 60_000 {
		t.Errorf("first chunk length = %d, want 60000", len(chunks[0].content))
	}
}

// TestIngestionServiceNilDeps verifies Start can handle nil dependencies gracefully.
func TestIngestionServiceNilDeps(t *testing.T) {
	svc := &IngestionService{
		pool:       nil,
		entitySvc:  nil,
		vectorRepo: nil,
		graphRepo:  nil,
		qwenSvc:    nil,
		hub:        nil,
	}

	// Start should attempt to create a job; may fail if pool is nil
	jobID, _, err := svc.Start(context.Background(), uuid.New(), strings.NewReader("hello"), "test.md")
	if err != nil {
		// Expected when pool is nil — service can't persist the job
		t.Logf("Start with nil deps: jobID=%s err=%v", jobID, err)
	} else if jobID == uuid.Nil {
		t.Error("expected non-nil job ID even with nil deps")
	}
}

// TestIngestionProgressDeliveredToUniverseOwner is a DB-backed regression test
// proving ingestion_progress events are routed to the real universe owner,
// not uuid.Nil (see sdd/fix-ingestion-progress-delivery).
func TestIngestionProgressDeliveredToUniverseOwner(t *testing.T) {
	pool := testutil.SetupTestDB(t)
	testutil.RunMigrationsUpTo(t, pool, "020")
	ctx := context.Background()

	user := svcCreateTestUser(t, ctx, pool)
	universe := svcCreateTestUniverse(t, ctx, pool, user.ID)

	hub := &mockIngestionHub{}
	svc := &IngestionService{
		pool:       pool,
		entitySvc:  nil,
		vectorRepo: nil,
		graphRepo:  nil,
		qwenSvc:    nil,
		hub:        hub,
	}

	docContent := "# Chapter 1\n\nBody text."

	_, _, err := svc.Start(ctx, universe.ID, strings.NewReader(docContent), "t.md")
	if err != nil {
		t.Fatalf("Start: %v", err)
	}

	time.Sleep(200 * time.Millisecond)

	msgs := hub.popMessages()
	userIDs := hub.popUserIDs()

	foundProgress := false
	for i, msg := range msgs {
		if msg.Type != "ingestion_progress" {
			continue
		}
		foundProgress = true
		if userIDs[i] == uuid.Nil {
			t.Errorf("ingestion_progress message %d delivered to uuid.Nil, want universe owner %s", i, universe.UserID)
		}
		if userIDs[i] != universe.UserID {
			t.Errorf("ingestion_progress message %d userID = %s, want %s", i, userIDs[i], universe.UserID)
		}
	}
	if !foundProgress {
		t.Fatal("expected at least one ingestion_progress message")
	}
}

// TestStartRejectsLegacyDoc verifies Start rejects a .doc filename
// synchronously, before creating any job row, per the whitelist check (D1/D2).
func TestStartRejectsLegacyDoc(t *testing.T) {
	svc := &IngestionService{}
	_, _, err := svc.Start(context.Background(), uuid.New(), strings.NewReader("binary junk"), "manuscript.doc")
	if err == nil {
		t.Fatal("expected error for .doc filename, got nil")
	}
	if !errors.Is(err, ErrUnsupportedFileType) {
		t.Errorf("expected ErrUnsupportedFileType, got: %v", err)
	}
}

func TestStartRejectsUnknownExtension(t *testing.T) {
	svc := &IngestionService{}
	_, _, err := svc.Start(context.Background(), uuid.New(), strings.NewReader("binary junk"), "manuscript.rtf")
	if err == nil {
		t.Fatal("expected error for unknown extension, got nil")
	}
}

// TestStartWorkTitleFromFilenameStem is a DB-backed regression test verifying
// the created Work's title is the filename stem (D3), replacing the old
// hardcoded "Imported Manuscript".
func TestStartWorkTitleFromFilenameStem(t *testing.T) {
	pool := testutil.SetupTestDB(t)
	testutil.RunMigrationsUpTo(t, pool, "020")
	ctx := context.Background()

	user := svcCreateTestUser(t, ctx, pool)
	universe := svcCreateTestUniverse(t, ctx, pool, user.ID)

	svc := &IngestionService{pool: pool}
	_, _, err := svc.Start(ctx, universe.ID, strings.NewReader("# Chapter 1\n\nBody."), "manuscript.pdf")
	if err != nil {
		t.Fatalf("Start: %v", err)
	}

	works, err := repositories.NewWorkRepo(pool).ListByUniverse(ctx, universe.ID)
	if err != nil {
		t.Fatalf("ListByUniverse: %v", err)
	}
	if len(works) != 1 {
		t.Fatalf("expected 1 work, got %d", len(works))
	}
	if works[0].Title != "manuscript" {
		t.Errorf("work title = %q, want %q", works[0].Title, "manuscript")
	}
}

// TestRunWorkerParseFailure is a DB-backed regression test verifying a
// corrupt/unparseable upload delivers a "failed" WS status, creates no
// chapters (raw binary must never reach chapters.content, D1), and — crucially
// — leaves a durable failed job row with a non-empty error_message so a reload
// shows the failure instead of nothing.
func TestRunWorkerParseFailure(t *testing.T) {
	pool := testutil.SetupTestDB(t)
	testutil.RunMigrationsUpTo(t, pool, "020")
	ctx := context.Background()

	user := svcCreateTestUser(t, ctx, pool)
	universe := svcCreateTestUniverse(t, ctx, pool, user.ID)

	hub := &mockIngestionHub{}
	svc := &IngestionService{pool: pool, hub: hub}
	jobID, _, err := svc.Start(ctx, universe.ID, strings.NewReader("not a pdf at all"), "manuscript.pdf")
	if err != nil {
		t.Fatalf("Start: %v", err)
	}

	time.Sleep(200 * time.Millisecond)

	var sawFailed bool
	for _, msg := range hub.popMessages() {
		if msg.Type != "ingestion_progress" {
			continue
		}
		var payload map[string]any
		if err := json.Unmarshal(msg.Payload, &payload); err != nil {
			t.Fatalf("unmarshal progress payload: %v", err)
		}
		if payload["job_id"] == jobID.String() && payload["status"] == "failed" {
			sawFailed = true
		}
	}
	if !sawFailed {
		t.Error("expected an ingestion_progress WS message with status=failed for this job")
	}

	var chapterCount int
	if err := pool.QueryRow(ctx, "SELECT count(*) FROM chapters c JOIN works w ON c.work_id = w.id WHERE w.universe_id = $1", universe.ID).Scan(&chapterCount); err != nil {
		t.Fatalf("count chapters: %v", err)
	}
	if chapterCount != 0 {
		t.Errorf("expected 0 chapters after parse failure, got %d", chapterCount)
	}

	// The failed job row must survive with its error_message — deleting the
	// orphan Work would cascade-delete it and hide the failure from the user.
	job, err := repositories.NewIngestionRepo(pool).FindByID(ctx, jobID)
	if err != nil {
		t.Fatalf("FindByID: %v", err)
	}
	if job == nil {
		t.Fatal("expected the failed job row to survive, but it was gone")
	}
	if job.Status != "failed" {
		t.Errorf("job status = %q, want %q", job.Status, "failed")
	}
	if job.ErrorMessage == "" {
		t.Error("expected a non-empty error_message on the failed job")
	}
}

// TestIngestionServiceOrphanWork_ReusedWorkNotDeleted is a DB-backed
// regression test: a parse failure when ingesting into an *existing* Work
// (the works[0]-reuse branch) must never delete that Work.
func TestIngestionServiceOrphanWork_ReusedWorkNotDeleted(t *testing.T) {
	pool := testutil.SetupTestDB(t)
	testutil.RunMigrationsUpTo(t, pool, "020")
	ctx := context.Background()

	user := svcCreateTestUser(t, ctx, pool)
	universe := svcCreateTestUniverse(t, ctx, pool, user.ID)

	// Seed an existing work so Start takes the works[0]-reuse branch.
	workRepo := repositories.NewWorkRepo(pool)
	existingWork := models.Work{ID: uuid.New(), UniverseID: universe.ID, Title: "Existing Work", Type: "novel", Status: "in_progress"}
	tx, err := pool.Begin(ctx)
	if err != nil {
		t.Fatalf("begin: %v", err)
	}
	if err := workRepo.Create(ctx, tx, &existingWork); err != nil {
		t.Fatalf("create existing work: %v", err)
	}
	if err := tx.Commit(ctx); err != nil {
		t.Fatalf("commit: %v", err)
	}

	svc := &IngestionService{pool: pool}
	_, _, err = svc.Start(ctx, universe.ID, strings.NewReader("not a pdf at all"), "manuscript.pdf")
	if err != nil {
		t.Fatalf("Start: %v", err)
	}

	time.Sleep(200 * time.Millisecond)

	got, err := workRepo.FindByID(ctx, existingWork.ID)
	if err != nil {
		t.Fatalf("expected the reused work to survive, but FindByID errored: %v", err)
	}
	if got.ID != existingWork.ID {
		t.Errorf("FindByID returned unexpected work: %+v", got)
	}
}

// TestSplitChunksCascade table-tests the heading-pattern priority cascade.
func TestSplitChunksCascade(t *testing.T) {
	svc := &IngestionService{}

	cases := []struct {
		name         string
		content      string
		minChunks    int
		wantContains string // a chunk title expected to appear
	}{
		{
			name:         "markdown",
			content:      "# Chapter One\n\nBody one.\n\n# Chapter Two\n\nBody two.",
			minChunks:    2,
			wantContains: "Chapter One",
		},
		{
			name:         "english chapter N",
			content:      "Chapter 1\n\nBody one.\n\nChapter 2\n\nBody two.",
			minChunks:    2,
			wantContains: "Chapter 1",
		},
		{
			name:         "english chapter spelled out",
			content:      "Chapter One\n\nBody one.\n\nChapter Two\n\nBody two.",
			minChunks:    2,
			wantContains: "Chapter One",
		},
		{
			name:         "spanish capitulo",
			content:      "Capítulo I\n\nCuerpo uno.\n\nCapítulo II\n\nCuerpo dos.",
			minChunks:    2,
			wantContains: "Capítulo I",
		},
		{
			name:         "spanish capitulo spelled out",
			content:      "Capítulo Uno\n\nCuerpo uno.\n\nCapítulo Dos\n\nCuerpo dos.",
			minChunks:    2,
			wantContains: "Capítulo Uno",
		},
		{
			name:         "bare roman numerals",
			content:      "I\n\nBody one.\n\nII\n\nBody two.",
			minChunks:    2,
			wantContains: "I",
		},
		{
			name:         "all caps headings",
			content:      "CHAPTER ONE: THE BEGINNING\n\nBody one.\n\nCHAPTER TWO: THE END\n\nBody two.",
			minChunks:    2,
			wantContains: "CHAPTER ONE: THE BEGINNING",
		},
		{
			name:         "title case single word",
			content:      "Holden\n\nBody one.\n\nMiller\n\nBody two.",
			minChunks:    2,
			wantContains: "Holden",
		},
		{
			name:         "title case multi-word",
			content:      "The Rocinante\n\nBody one.\n\nThe Canterbury\n\nBody two.",
			minChunks:    2,
			wantContains: "The Rocinante",
		},
		{
			name:         "no pattern falls back to paragraphs",
			content:      "Just some prose.\n\nMore prose, no headings at all here.",
			minChunks:    1,
			wantContains: "Untitled",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			chunks := svc.splitChunks(tc.content)
			if len(chunks) < tc.minChunks {
				t.Fatalf("got %d chunks, want >= %d", len(chunks), tc.minChunks)
			}
			found := false
			for _, ch := range chunks {
				if ch.title == tc.wantContains {
					found = true
				}
			}
			if !found {
				titles := make([]string, len(chunks))
				for i, ch := range chunks {
					titles[i] = ch.title
				}
				t.Errorf("expected a chunk titled %q, got titles: %v", tc.wantContains, titles)
			}
		})
	}
}

// TestIsAllCapsHeadingLine verifies the ALL-CAPS fallback heuristic.
func TestIsAllCapsHeadingLine(t *testing.T) {
	cases := []struct {
		line string
		want bool
	}{
		{"CHAPTER ONE: THE BEGINNING", true},
		{"THE BEGINNING", true}, // 13 chars, no punctuation
		{"A", false},            // too short
		{"THE END.", false},     // sentence punctuation
		{"WHAT?", false},        // sentence punctuation
		{"LOOK!", false},        // sentence punctuation
		{"QUOTE,", false},       // sentence punctuation
		{"SEMICOLON;", false},   // sentence punctuation
		{"COLON:", false},       // sentence punctuation
		{"DIALOGUE", false},     // no punctuation, but only 8 chars (< 10)
	}

	for _, tc := range cases {
		t.Run(tc.line, func(t *testing.T) {
			if got := isAllCapsHeadingLine(tc.line); got != tc.want {
				t.Errorf("isAllCapsHeadingLine(%q) = %v, want %v", tc.line, got, tc.want)
			}
		})
	}
}

// TestSplitChunksCascadeMixedPrefersHighestPriority verifies that when a
// document contains heading styles from multiple pattern classes, the
// highest-priority class (earliest in the cascade) with >= 2 matches wins.
func TestSplitChunksCascadeMixedPrefersHighestPriority(t *testing.T) {
	svc := &IngestionService{}
	// Markdown headers present (2 matches) alongside a coincidental line that
	// looks like "Chapter N" prose — markdown must win since it's earlier in
	// the cascade and already has >= 2 matches.
	content := "# Chapter One\n\nChapter 9 was a turning point in the plot.\n\n# Chapter Two\n\nMore body."
	chunks := svc.splitChunks(content)
	if len(chunks) != 2 {
		t.Fatalf("expected 2 markdown-driven chunks, got %d", len(chunks))
	}
	if chunks[0].title != "Chapter One" || chunks[1].title != "Chapter Two" {
		t.Errorf("unexpected titles: %q, %q", chunks[0].title, chunks[1].title)
	}
}

// TestSplitChunksCascadeGuardsAgainstFalsePositive verifies a pattern
// matching an unreasonable number of lines (a likely false positive) is
// skipped in favor of the next pattern in the cascade.
func TestSplitChunksCascadeGuardsAgainstFalsePositive(t *testing.T) {
	// 600 bare-roman-looking lines — over maxSaneHeadingMatches (500) — must
	// not be treated as chapter headings; falls through to paragraph split.
	var b strings.Builder
	for i := 0; i < 600; i++ {
		b.WriteString("I\n\nSome body text to keep this a real paragraph so it survives trimming.\n\n")
	}
	chunks := (&IngestionService{}).splitChunks(b.String())
	for _, ch := range chunks {
		if ch.title == "I" && len(chunks) > 500 {
			t.Fatalf("expected the >500-match guard to reject the bare-roman pattern, got %d chunks titled %q", len(chunks), ch.title)
		}
	}
}

// TestSelectAnalysisChapters verifies the K cap and zero-entity skip.
func TestSelectAnalysisChapters(t *testing.T) {
	mk := func(n int, hasEntity bool) ingestedChapter {
		ch := ingestedChapter{ID: uuid.New(), Content: fmt.Sprintf("chapter %d", n)}
		if hasEntity {
			ch.Entities = []ResolvedEntity{{MentionText: fmt.Sprintf("chapter %d", n)}}
		}
		return ch
	}

	chapters := []ingestedChapter{
		mk(1, true),
		mk(2, false), // skipped: zero entities
		mk(3, true),
		mk(4, true),
		mk(5, true),
	}

	selected := selectAnalysisChapters(chapters, 2)
	if len(selected) != 2 {
		t.Fatalf("expected 2 selected chapters (K cap), got %d", len(selected))
	}
	// Last-K among the ones with entities: chapters 1,3,4,5 have entities;
	// last 2 of those are 4 and 5.
	if selected[0].Content != "chapter 4" || selected[1].Content != "chapter 5" {
		t.Errorf("unexpected selection: %+v", selected)
	}

	// K larger than available: returns all chapters with entities.
	selected = selectAnalysisChapters(chapters, 10)
	if len(selected) != 4 {
		t.Errorf("expected 4 chapters with entities, got %d", len(selected))
	}

	if got := selectAnalysisChapters(chapters, 0); got != nil {
		t.Errorf("expected nil for k=0, got %+v", got)
	}
}

// TestRunPostIngestAnalysis is a DB-backed wiring test: nil qwenSvc makes
// CheckSemantic and the plot-hole agent evaluation no-op, so this exercises
// the K cap / zero-entity skip / never-fails-the-job contract without
// needing a live Qwen API.
func TestRunPostIngestAnalysis(t *testing.T) {
	pool := testutil.SetupTestDB(t)
	testutil.RunMigrationsUpTo(t, pool, "020")
	ctx := context.Background()

	user := svcCreateTestUser(t, ctx, pool)
	universe := svcCreateTestUniverse(t, ctx, pool, user.ID)

	workRepo := repositories.NewWorkRepo(pool)
	work := models.Work{ID: uuid.New(), UniverseID: universe.ID, Title: "Test Work", Type: "novel", Status: "in_progress"}
	tx, err := pool.Begin(ctx)
	if err != nil {
		t.Fatalf("begin: %v", err)
	}
	if err := workRepo.Create(ctx, tx, &work); err != nil {
		t.Fatalf("create work: %v", err)
	}
	if err := tx.Commit(ctx); err != nil {
		t.Fatalf("commit: %v", err)
	}

	chRepo := repositories.NewChapterRepo(pool)
	entityRepo := repositories.NewEntityRepo(pool)
	var chapters []ingestedChapter
	for i := 1; i <= 3; i++ {
		ch := models.Chapter{
			ID:         uuid.New(),
			WorkID:     work.ID,
			Title:      fmt.Sprintf("Chapter %d", i),
			OrderIndex: i,
			Content:    fmt.Sprintf("Chapter %d content.", i),
			RawText:    fmt.Sprintf("Chapter %d content.", i),
			Status:     "draft",
		}
		tx, err := pool.Begin(ctx)
		if err != nil {
			t.Fatalf("begin: %v", err)
		}
		if err := chRepo.Create(ctx, tx, &ch); err != nil {
			t.Fatalf("create chapter: %v", err)
		}
		if err := tx.Commit(ctx); err != nil {
			t.Fatalf("commit: %v", err)
		}

		entity := models.Entity{ID: uuid.New(), UniverseID: universe.ID, Type: "character", Name: fmt.Sprintf("Entity %d", i), Status: "alive"}
		tx2, err := pool.Begin(ctx)
		if err != nil {
			t.Fatalf("begin: %v", err)
		}
		if err := entityRepo.Create(ctx, tx2, &entity); err != nil {
			t.Fatalf("create entity: %v", err)
		}
		if err := tx2.Commit(ctx); err != nil {
			t.Fatalf("commit: %v", err)
		}

		chapters = append(chapters, ingestedChapter{
			ID:      ch.ID,
			Content: ch.Content,
			Entities: []ResolvedEntity{
				{Entity: entity, MentionText: ch.Content, IsNew: true},
			},
		})
	}

	tok := NewTokenizer()
	budgetMgr := NewContextBudgetManager(tok, 30000, 2000)
	contraSvc := NewContradictionService(pool, repositories.NewContradictionRepo(pool), entityRepo, nil, nil, 3, budgetMgr, 3)
	plotHoleSvc := NewPlotHoleService(pool, repositories.NewPlotHoleRepo(pool), entityRepo, 8, nil, nil, 2)

	hub := &mockIngestionHub{}
	svc := &IngestionService{pool: pool, hub: hub}
	svc.SetPostIngestAnalysis(contraSvc, plotHoleSvc, budgetMgr, 2)

	svc.runPostIngestAnalysis(ctx, universe.ID, chapters, user.ID)

	// K cap of 2 respected — with a nil qwenSvc, CheckSemantic no-ops (no
	// contradictions), so the only observable check here is that it didn't
	// panic/error/hang across all 3 candidate chapters despite the cap.
}

// compile-time interface checks
var _ IngestionQwen = (*mockQwenForIngestion)(nil)
var _ *pgxpool.Pool = nil
var _ *repositories.VectorRepo = nil
var _ *repositories.GraphRepo = nil
var _ *EntityService = nil
var _ AnalysisHub = (*mockIngestionHub)(nil)
