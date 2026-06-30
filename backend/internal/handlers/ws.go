package handlers

import (
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"

	"github.com/quill/backend/internal/ws"
)

// WSHandler bridges Fiber HTTP requests to the WebSocket hub.
//
// ponytail: thin adapter — just upgrades the connection and delegates to hub.Handle.
type WSHandler struct {
	hub *ws.Hub
}

// NewWSHandler creates a WebSocket handler backed by the given hub.
func NewWSHandler(hub *ws.Hub) *WSHandler {
	return &WSHandler{hub: hub}
}

// Upgrade is a Fiber-compatible WebSocket upgrade handler.
// Use with: app.Get("/api/v1/ws", websocket.New(handler.Upgrade))
func (h *WSHandler) Upgrade(c *websocket.Conn) {
	h.hub.Handle(c)
}

// Handler returns a Fiber-compatible route handler using websocket.New.
// Use: app.Get("/api/v1/ws", handler.Handler())
func (h *WSHandler) Handler() fiber.Handler {
	return websocket.New(func(c *websocket.Conn) {
		h.hub.Handle(c)
	})
}
