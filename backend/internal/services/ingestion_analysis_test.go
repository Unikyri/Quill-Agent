package services

import (
	"context"
	"encoding/json"
	"sync"
	"testing"

	"github.com/google/uuid"

	"github.com/quill/backend/internal/models"
	"github.com/quill/backend/internal/ws"
)

type recordingPostIngestContradictionAnalyzer struct {
	mu    sync.Mutex
	calls []uuid.UUID
}

func (a *recordingPostIngestContradictionAnalyzer) CheckSemantic(_ context.Context, universeID, chapterID uuid.UUID, _ string, _ []ResolvedEntity, _ ...func(string, *QwenToolCall)) ([]models.Contradiction, error) {
	a.mu.Lock()
	a.calls = append(a.calls, chapterID)
	a.mu.Unlock()
	return []models.Contradiction{{ID: uuid.New(), UniverseID: universeID}}, nil
}

type recordingPostIngestPlotHoleAnalyzer struct {
	mu        sync.Mutex
	chapterID uuid.UUID
	calls     int
}

func (a *recordingPostIngestPlotHoleAnalyzer) Scan(_ context.Context, _ uuid.UUID, chapterID uuid.UUID) ([]models.PlotHole, error) {
	a.mu.Lock()
	a.chapterID = chapterID
	a.calls++
	a.mu.Unlock()
	return nil, nil
}

// TestSafeAnalyzeRecoversPanic proves the recover seam used by
// runPostIngestAnalysis contains a panic instead of letting it escape and
// crash the process. Post-ingest analysis parses untrusted LLM tool-call JSON
// and type-asserts on it; an unrecovered panic on any goroutine would take
// down every user's WebSocket, not just the ingestion job. If safeAnalyze
// failed to recover, this panic would propagate and fail (crash) the test.
func TestSafeAnalyzeRecoversPanic(t *testing.T) {
	safeAnalyze("test unit", func() {
		panic("boom from untrusted analysis")
	})
	// Reaching this line means the panic was recovered, not propagated.
}

func TestEmitContradictionAlertCarriesUniverseID(t *testing.T) {
	hub := &mockIngestionHub{}
	svc := &IngestionService{hub: hub}
	universeID := uuid.New()
	svc.emitContradictionAlert(uuid.New(), models.Contradiction{ID: uuid.New(), UniverseID: universeID})

	msgs := hub.popMessages()
	if len(msgs) != 1 || msgs[0].Type != ws.TypeContradictionAlert {
		t.Fatalf("expected one contradiction_alert, got %+v", msgs)
	}
	var payload models.ContradictionAlertPayload
	if err := json.Unmarshal(msgs[0].Payload, &payload); err != nil {
		t.Fatalf("unmarshal contradiction_alert: %v", err)
	}
	if payload.UniverseID != universeID {
		t.Errorf("contradiction_alert universe_id = %s, want %s", payload.UniverseID, universeID)
	}
}

func TestRunPostIngestAnalysisDispatchesSelectedResults(t *testing.T) {
	universeID := uuid.New()
	ownerID := uuid.New()
	firstID, secondID, thirdID := uuid.New(), uuid.New(), uuid.New()
	contra := &recordingPostIngestContradictionAnalyzer{}
	plotHoles := &recordingPostIngestPlotHoleAnalyzer{}
	hub := &mockIngestionHub{}
	svc := &IngestionService{hub: hub}
	svc.SetPostIngestAnalysis(contra, plotHoles, NewContextBudgetManager(NewTokenizer(), 1_000, 100), 2)

	svc.runPostIngestAnalysis(context.Background(), universeID, []ingestedChapter{
		{ID: firstID, Content: "first", Entities: []ResolvedEntity{{}}},
		{ID: secondID, Content: "second", Entities: []ResolvedEntity{{}}},
		{ID: thirdID, Content: "third", Entities: []ResolvedEntity{{}}},
	}, ownerID)

	contra.mu.Lock()
	calls := append([]uuid.UUID(nil), contra.calls...)
	contra.mu.Unlock()
	if len(calls) != 2 || !containsChapterID(calls, secondID) || !containsChapterID(calls, thirdID) {
		t.Fatalf("CheckSemantic chapter calls = %v, want only %s and %s", calls, secondID, thirdID)
	}
	plotHoles.mu.Lock()
	plotChapterID, plotCalls := plotHoles.chapterID, plotHoles.calls
	plotHoles.mu.Unlock()
	if plotCalls != 1 || plotChapterID != thirdID {
		t.Fatalf("plot-hole Scan = %d calls for %s, want one call for %s", plotCalls, plotChapterID, thirdID)
	}

	msgs := hub.popMessages()
	if len(msgs) != 2 {
		t.Fatalf("contradiction alerts = %d, want 2", len(msgs))
	}
	for _, msg := range msgs {
		if msg.Type != ws.TypeContradictionAlert {
			t.Errorf("message type = %q, want %q", msg.Type, ws.TypeContradictionAlert)
		}
		var payload models.ContradictionAlertPayload
		if err := json.Unmarshal(msg.Payload, &payload); err != nil {
			t.Fatalf("unmarshal contradiction alert: %v", err)
		}
		if payload.UniverseID != universeID || payload.Contradiction.UniverseID != universeID {
			t.Errorf("alert universe IDs = %+v, want %s", payload, universeID)
		}
	}
}

func containsChapterID(chapters []uuid.UUID, want uuid.UUID) bool {
	for _, chapterID := range chapters {
		if chapterID == want {
			return true
		}
	}
	return false
}
