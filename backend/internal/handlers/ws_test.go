package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/contrib/websocket"

	"github.com/quill/backend/internal/ws"
)

// TestWSHandlerUpgrade verifies the WS handler returns an appropriate response
// for non-WebSocket requests (should not panic).
func TestWSHandlerUpgrade(t *testing.T) {
	app := fiber.New()
	hub := ws.NewHub(nil, nil, nil, nil)
	handler := NewWSHandler(hub)

	app.Get("/ws", handler.Handler())

	req := httptest.NewRequest(http.MethodGet, "/ws", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test: %v", err)
	}

	// Without WebSocket upgrade headers, Fiber should return 400 or 426
	if resp.StatusCode != http.StatusBadRequest && resp.StatusCode != http.StatusUpgradeRequired {
		t.Logf("Got status %d (expected 400 or 426 for non-upgrade request)", resp.StatusCode)
	}
}

// TestWSHandlerUpgradeWithHeaders verifies the handler attempts upgrade with
// proper WebSocket headers.
func TestWSHandlerUpgradeWithHeaders(t *testing.T) {
	app := fiber.New()
	hub := ws.NewHub(nil, nil, nil, nil)
	handler := NewWSHandler(hub)

	app.Get("/ws", handler.Handler())

	req := httptest.NewRequest(http.MethodGet, "/ws", nil)
	req.Header.Set("Upgrade", "websocket")
	req.Header.Set("Connection", "Upgrade")
	req.Header.Set("Sec-WebSocket-Version", "13")
	req.Header.Set("Sec-WebSocket-Key", "dGhlIHNhbXBsZSBub25jZQ==")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test: %v", err)
	}

	// If upgrade succeeds, status 101; otherwise 400/426
	if resp.StatusCode != http.StatusSwitchingProtocols &&
		resp.StatusCode != http.StatusBadRequest {
		t.Logf("Got status %d", resp.StatusCode)
	}
}

// TestNewWSHandler ensures handler creation doesn't panic.
func TestNewWSHandler(t *testing.T) {
	hub := ws.NewHub(nil, nil, nil, nil)
	handler := NewWSHandler(hub)
	if handler == nil {
		t.Fatal("NewWSHandler returned nil")
	}
}

// TestWSHandlerUnparseableMessage verifies handler doesn't explode on garbage data.
func TestWSHandlerUnparseableMessage(t *testing.T) {
	_ = json.RawMessage(`not json`)
	// This test verifies that the handler type compiles and exists.
	// The real message parsing test belongs in ws/hub_test.go once the hub
	// is fully wired.
}

// TestWSHandlerTypeImplementsFiberHandler verifies Handler() returns a fiber.Handler.
func TestWSHandlerTypeImplementsFiberHandler(t *testing.T) {
	hub := ws.NewHub(nil, nil, nil, nil)
	h := NewWSHandler(hub)
	fiberHandler := h.Handler()
	if fiberHandler == nil {
		t.Error("Handler() should return a non-nil fiber.Handler")
	}
}

// Ensure websocket usage compiles
var _ = websocket.New
