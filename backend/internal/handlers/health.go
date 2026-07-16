package handlers

import (
	"context"
	"net/http"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/quill/backend/internal/config"
	"github.com/quill/backend/internal/services"
)

type HealthHandler struct {
	pool        *pgxpool.Pool
	qwenSvc     services.LLMHealthChecker
	qwenTimeout time.Duration
	start       time.Time
}

func NewHealthHandler(pool *pgxpool.Pool, qwenSvc services.LLMHealthChecker, cfg *config.Config) *HealthHandler {
	return &HealthHandler{
		pool:        pool,
		qwenSvc:     qwenSvc,
		qwenTimeout: cfg.QwenHealthTimeout,
		start:       time.Now(),
	}
}

func (h *HealthHandler) Check(c *fiber.Ctx) error {
	ctx := c.Context()

	dbStatus := probeDB(ctx, h.pool)
	ageStatus := probeExtension(ctx, h.pool, "age")
	pgvectorStatus := probeExtension(ctx, h.pool, "vector")
	qwenStatus := probeQwen(ctx, h.qwenSvc, h.qwenTimeout)

	var stat syscall.Statfs_t
	diskFreeMB := int64(0)
	if err := syscall.Statfs("/", &stat); err == nil {
		diskFreeMB = int64(stat.Bavail) * int64(stat.Bsize) / (1024 * 1024)
	}

	status := "healthy"
	if dbStatus != "connected" || ageStatus != "available" || pgvectorStatus != "available" {
		status = "unhealthy"
	} else if qwenStatus != "reachable" {
		status = "degraded"
	}

	code := http.StatusOK
	if status != "healthy" {
		code = http.StatusServiceUnavailable
	}

	return c.Status(code).JSON(fiber.Map{
		"status":         status,
		"db":             dbStatus,
		"age":            ageStatus,
		"pgvector":       pgvectorStatus,
		"qwen_api":       qwenStatus,
		"disk_free_mb":   diskFreeMB,
		"uptime_seconds": int64(time.Since(h.start).Seconds()),
	})
}

func probeDB(ctx context.Context, pool *pgxpool.Pool) string {
	if err := pool.Ping(ctx); err != nil {
		return "disconnected"
	}
	return "connected"
}

func probeExtension(ctx context.Context, pool *pgxpool.Pool, extname string) string {
	var exists bool
	if err := pool.QueryRow(ctx, "SELECT EXISTS (SELECT 1 FROM pg_extension WHERE extname = $1)", extname).Scan(&exists); err != nil {
		return "unavailable"
	}
	if !exists {
		return "unavailable"
	}
	return "available"
}

func probeQwen(ctx context.Context, svc services.LLMHealthChecker, timeout time.Duration) string {
	if svc == nil {
		return "not_configured"
	}

	checkCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	if err := svc.HealthCheck(checkCtx); err != nil {
		return "unreachable"
	}
	return "reachable"
}
