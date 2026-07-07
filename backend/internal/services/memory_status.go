package services

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/quill/backend/internal/repositories"
)

// defaultRelevanceDeltaEpsilon mirrors config.RelevanceDeltaEpsilon's default
// and is used as a fallback when MemoryService's epsilon was never wired
// (e.g. in tests that construct the service without SetRelevanceDeltaEpsilon).
const defaultRelevanceDeltaEpsilon = 0.01

// historyCap bounds how many recent relevance-history rows are fetched per
// entity for the memory-status endpoint (spec: 30, oldest-first per entity).
const historyCap = 30

// allEntitiesLimit is a pragmatic "no pagination" cap passed to
// EntityRepo.ListByUniverse, which requires a positive Limit. Memory-status
// needs every entity in the universe, not a page.
const allEntitiesLimit = 1_000_000

// MemoryStatusHistoryPoint is one sampled relevance-score datapoint exposed
// over the memory-status endpoint.
type MemoryStatusHistoryPoint struct {
	Score      float64   `json:"score"`
	RecordedAt time.Time `json:"recorded_at"`
}

// MemoryStatusEntity is one entity's lifecycle snapshot in the
// memory-status response.
type MemoryStatusEntity struct {
	ID             uuid.UUID                  `json:"id"`
	Name           string                     `json:"name"`
	Type           string                     `json:"type"`
	RelevanceScore float64                    `json:"relevance_score"`
	Status         string                     `json:"status"`
	Consolidated   bool                       `json:"consolidated"`
	Lifecycle      string                     `json:"lifecycle"`
	History        []MemoryStatusHistoryPoint `json:"history"`
}

// MemoryStatusResponse is the full memory-status endpoint payload.
type MemoryStatusResponse struct {
	ConsolidatedCount int                  `json:"consolidated_count"`
	Entities          []MemoryStatusEntity `json:"entities"`
}

// SetHistoryRepo wires the relevance-history pipeline used by MemoryStatus.
// Optional — nil-safe; entities are returned with empty history when unset.
func (s *MemoryService) SetHistoryRepo(r *repositories.EntityRelevanceHistoryRepo) {
	s.historyRepo = r
}

// SetRelevanceDeltaEpsilon wires the epsilon used to distinguish real decay
// from float noise when deriving the "decaying" lifecycle state. Optional —
// falls back to defaultRelevanceDeltaEpsilon when unset (<=0).
func (s *MemoryService) SetRelevanceDeltaEpsilon(epsilon float64) {
	s.relevanceDeltaEpsilon = epsilon
}

// MemoryStatus composes the entity list, capped relevance history, and
// consolidation set for a universe, and derives each entity's lifecycle
// state server-side (see deriveLifecycle).
func (s *MemoryService) MemoryStatus(ctx context.Context, universeID uuid.UUID) (MemoryStatusResponse, error) {
	entities, _, err := s.entityRepo.ListByUniverse(ctx, universeID, repositories.EntityFilters{Page: 1, Limit: allEntitiesLimit})
	if err != nil {
		return MemoryStatusResponse{}, fmt.Errorf("memory status: list entities: %w", err)
	}

	var historyPoints []repositories.RelevanceHistoryPoint
	if s.historyRepo != nil {
		historyPoints, err = s.historyRepo.ListRecentByUniverse(ctx, universeID, historyCap)
		if err != nil {
			return MemoryStatusResponse{}, fmt.Errorf("memory status: list history: %w", err)
		}
	}
	historyByEntity := make(map[uuid.UUID][]repositories.RelevanceHistoryPoint, len(entities))
	for _, p := range historyPoints {
		historyByEntity[p.EntityID] = append(historyByEntity[p.EntityID], p)
	}

	consolidatedSet := make(map[uuid.UUID]bool)
	if s.consolidationRepo != nil {
		ids, err := s.consolidationRepo.EntityIDsWithConsolidation(ctx, universeID)
		if err != nil {
			return MemoryStatusResponse{}, fmt.Errorf("memory status: list consolidated: %w", err)
		}
		for _, id := range ids {
			consolidatedSet[id] = true
		}
	}

	epsilon := s.relevanceDeltaEpsilon
	if epsilon <= 0 {
		epsilon = defaultRelevanceDeltaEpsilon
	}

	resp := MemoryStatusResponse{
		ConsolidatedCount: len(consolidatedSet),
		Entities:          make([]MemoryStatusEntity, 0, len(entities)),
	}
	for _, e := range entities {
		hist := historyByEntity[e.ID]
		consolidated := consolidatedSet[e.ID]

		points := make([]MemoryStatusHistoryPoint, len(hist))
		for i, h := range hist {
			points[i] = MemoryStatusHistoryPoint{Score: h.RelevanceScore, RecordedAt: h.RecordedAt}
		}

		resp.Entities = append(resp.Entities, MemoryStatusEntity{
			ID:             e.ID,
			Name:           e.Name,
			Type:           e.Type,
			RelevanceScore: e.RelevanceScore,
			Status:         e.Status,
			Consolidated:   consolidated,
			Lifecycle:      deriveLifecycle(e.Status, hist, consolidated, epsilon),
			History:        points,
		})
	}
	return resp, nil
}

// deriveLifecycle computes the server-derived lifecycle state for one
// entity, in priority order: consolidated > reactivated > archived >
// decaying > active. It never stores a new status value — this is a pure
// read-time projection over status/history/consolidated.
func deriveLifecycle(status string, history []repositories.RelevanceHistoryPoint, consolidated bool, epsilon float64) string {
	if status == "archived" {
		if consolidated {
			return "consolidated"
		}
		return "archived"
	}

	if isReactivated(history) {
		return "reactivated"
	}
	if isDecaying(history, epsilon) {
		return "decaying"
	}
	return "active"
}

// isReactivated reports whether the most recent status transition within
// the (capped) history window is archived->active. It scans backward for
// the nearest pair of consecutive rows whose status differs — the "most
// recent transition" — rather than just comparing the last two rows, since
// several same-status rows may follow the flip.
func isReactivated(history []repositories.RelevanceHistoryPoint) bool {
	for i := len(history) - 1; i > 0; i-- {
		if history[i].Status != history[i-1].Status {
			return history[i-1].Status == "archived" && history[i].Status == "active"
		}
	}
	return false
}

// isDecaying reports whether the latest history row's score dropped from
// the previous row's score by more than epsilon. Fewer than 2 rows always
// reads as not-decaying (spec: 0/1-row entities are "active", never
// undefined/error).
func isDecaying(history []repositories.RelevanceHistoryPoint, epsilon float64) bool {
	if len(history) < 2 {
		return false
	}
	last := history[len(history)-1]
	prev := history[len(history)-2]
	return last.RelevanceScore-prev.RelevanceScore < -epsilon
}
