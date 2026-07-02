package handlers

import (
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

func TestNewGraphHandler(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for nil graphRepo")
		}
	}()
	NewGraphHandler(nil, nil, nil)
}
