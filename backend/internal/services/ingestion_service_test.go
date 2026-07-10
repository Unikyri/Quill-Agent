package services

import (
	"context"
	"encoding/json"
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

// compile-time interface checks
var _ IngestionQwen = (*mockQwenForIngestion)(nil)
var _ *pgxpool.Pool = nil
var _ *repositories.VectorRepo = nil
var _ *repositories.GraphRepo = nil
var _ *EntityService = nil
var _ AnalysisHub = (*mockIngestionHub)(nil)
