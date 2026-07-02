package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"github.com/quill/backend/internal/repositories"
)

// ── TimelineHandler tests ──

func TestTimelineHandlerListInvalidID(t *testing.T) {
	app := fiber.New()
	h := NewTimelineHandler(nil, repositories.NewTimelineRepo(nil))
	app.Get("/api/v1/universes/:universe_id/timeline", h.ListByUniverse)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/universes/bad/timeline", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test: %v", err)
	}
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", resp.StatusCode)
	}
}

func TestTimelineHandlerCreate(t *testing.T) {
	app := fiber.New()
	h := NewTimelineHandler(nil, repositories.NewTimelineRepo(nil))
	app.Post("/api/v1/universes/:universe_id/timeline", h.Create)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/universes/"+uuid.New().String()+"/timeline", nil)
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test: %v", err)
	}
	if resp.StatusCode < 400 {
		t.Errorf("expected error on empty body, got %d", resp.StatusCode)
	}
}

func TestNewTimelineHandler(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for nil timelineRepo")
		}
	}()
	NewTimelineHandler(nil, nil)
}
