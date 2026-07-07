package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"github.com/quill/backend/internal/repositories"
	"github.com/quill/backend/internal/services"
)

// ── GraphHandler tests ──

func TestGraphHandlerFullGraphInvalidID(t *testing.T) {
	app := fiber.New()
	h := NewGraphHandler(repositories.NewGraphRepo(nil), services.NewMemoryService(nil, nil, nil), repositories.NewEntityRepo(nil))
	app.Get("/api/v1/universes/:universe_id/graph", h.FullGraph)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/universes/bad/graph", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test: %v", err)
	}
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", resp.StatusCode)
	}
}

func TestGraphHandlerNeighborsInvalidID(t *testing.T) {
	app := fiber.New()
	h := NewGraphHandler(repositories.NewGraphRepo(nil), services.NewMemoryService(nil, nil, nil), repositories.NewEntityRepo(nil))
	app.Get("/api/v1/entities/:id/neighbors", h.Neighbors)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/entities/bad/neighbors?universe_id="+uuid.New().String(), nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test: %v", err)
	}
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", resp.StatusCode)
	}
}

func TestGraphHandlerRecall(t *testing.T) {
	app := fiber.New()
	h := NewGraphHandler(repositories.NewGraphRepo(nil), services.NewMemoryService(nil, nil, nil), repositories.NewEntityRepo(nil))
	app.Post("/api/v1/universes/:id/recall", h.Recall)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/universes/"+uuid.New().String()+"/recall", nil)
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test: %v", err)
	}

	if resp.StatusCode < 400 {
		t.Errorf("expected error status, got %d", resp.StatusCode)
	}
}

func TestGraphHandlerMemoryStatusInvalidID(t *testing.T) {
	app := fiber.New()
	h := NewGraphHandler(repositories.NewGraphRepo(nil), services.NewMemoryService(nil, nil, nil), repositories.NewEntityRepo(nil))
	app.Get("/api/v1/universes/:id/memory-status", h.MemoryStatus)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/universes/bad/memory-status", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test: %v", err)
	}
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", resp.StatusCode)
	}
}

func TestGraphHandlerNeighborsMissingGraph(t *testing.T) {
	app := fiber.New()

	stub := &stubGraphQuerier{errorMsg: `graph "universe_123e4567-e89b-12d3-a456-426614174000" does not exist`}
	h := &GraphHandler{
		graphRepo:  stub,
		memorySvc:  services.NewMemoryService(nil, nil, nil),
		entityRepo: repositories.NewEntityRepo(nil),
	}
	app.Get("/api/v1/entities/:id/neighbors", h.Neighbors)

	validID := "123e4567-e89b-12d3-a456-426614174000"
	req := httptest.NewRequest(http.MethodGet, "/api/v1/entities/"+validID+"/neighbors?universe_id="+validID, nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200 for missing graph neighbors, got %d", resp.StatusCode)
	}

	var body map[string]json.RawMessage
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatalf("decode: %v", err)
	}
	nodes, ok := body["nodes"]
	if !ok {
		t.Fatal("response missing 'nodes'")
	}
	edges, ok := body["edges"]
	if !ok {
		t.Fatal("response missing 'edges'")
	}
	if string(nodes) != "[]" || string(edges) != "[]" {
		t.Errorf("expected empty arrays, got nodes=%s edges=%s", nodes, edges)
	}
}

func TestNewGraphHandler(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for nil graphRepo")
		}
	}()
	NewGraphHandler(nil, nil, nil)
}

// ── stub graph querier for testing error paths ──

type stubGraphQuerier struct{ errorMsg string }

func (s *stubGraphQuerier) FullQuery(_ context.Context, _ string) ([]repositories.GraphNode, []repositories.GraphEdge, error) {
	return nil, nil, &stubQuerierErr{msg: s.errorMsg}
}
func (s *stubGraphQuerier) NHopTraversal(_ context.Context, _ string, _ string, _ int) ([]repositories.GraphNode, []repositories.GraphEdge, error) {
	return nil, nil, &stubQuerierErr{msg: s.errorMsg}
}

type stubQuerierErr struct{ msg string }

func (e *stubQuerierErr) Error() string { return e.msg }

func TestGraphHandlerFullGraphMissingGraph(t *testing.T) {
	app := fiber.New()

	stub := &stubGraphQuerier{errorMsg: `graph "universe_123e4567-e89b-12d3-a456-426614174000" does not exist`}
	h := &GraphHandler{
		graphRepo:  stub,
		memorySvc:  services.NewMemoryService(nil, nil, nil),
		entityRepo: repositories.NewEntityRepo(nil),
	}
	app.Get("/api/v1/universes/:universe_id/graph", h.FullGraph)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/universes/123e4567-e89b-12d3-a456-426614174000/graph", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200 for missing graph, got %d", resp.StatusCode)
	}

	var body map[string]json.RawMessage
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatalf("decode: %v", err)
	}
	nodes, ok := body["nodes"]
	if !ok {
		t.Fatal("response missing 'nodes'")
	}
	edges, ok := body["edges"]
	if !ok {
		t.Fatal("response missing 'edges'")
	}
	if string(nodes) != "[]" || string(edges) != "[]" {
		t.Errorf("expected empty arrays, got nodes=%s edges=%s", nodes, edges)
	}
}
