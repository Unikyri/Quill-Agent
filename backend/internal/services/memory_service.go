package services

import (
	"context"
	"fmt"
	"sort"

	"github.com/google/uuid"

	"github.com/quill/backend/internal/models"
	"github.com/quill/backend/internal/repositories"
)

// ResolvedEntity is a shared type between MemoryService, ContradictionService,
// and AnalysisService. It pairs an entity with the mention text that triggered
// the resolution.
type ResolvedEntity struct {
	Entity      models.Entity
	MentionText string
	IsNew       bool
	// PreviousStatus is the entity's status as it was in the DB before this
	// mention's data was merged in (see EntityService.ResolveOrCreate).
	// Deterministic contradiction checks must compare against this, not
	// Entity.Status, which already reflects the newly-merged value.
	PreviousStatus string
}

// MemoryService provides contextual recall by merging graph neighbourhood
// traversal, recent mentions, and freshness signals into a ranked list.
//
// ponytail: deterministic ranking (weighted average), no LLM needed.
// Weights: graph×0.4 + recency×0.3 + freshness×0.3.
type MemoryService struct {
	graphRepo  *repositories.GraphRepo
	entityRepo *repositories.EntityRepo
	vectorRepo *repositories.VectorRepo
}

// NewMemoryService creates a memory service with the given repos.
func NewMemoryService(graphRepo *repositories.GraphRepo, entityRepo *repositories.EntityRepo, vectorRepo *repositories.VectorRepo) *MemoryService {
	return &MemoryService{
		graphRepo:  graphRepo,
		entityRepo: entityRepo,
		vectorRepo: vectorRepo,
	}
}

// Recall returns up to k RecallItems for a universe, merging graph neighbours,
// recent mentions, and vector-based freshness signals into a single ranked list.
func (s *MemoryService) Recall(ctx context.Context, universeID uuid.UUID, queryEmbedding []float32, k int) ([]models.RecallItem, error) {
	// ── 1. Recent mentions (recency × 0.3) ──
	entities, err := s.entityRepo.ListByUniverseActive(ctx, universeID)
	if err != nil {
		return nil, fmt.Errorf("recall: list active entities: %w", err)
	}

	graphName := "universe_" + universeID.String()

	// Collect candidates from each source
	candidateMap := make(map[uuid.UUID]*models.RecallItem)
	graphScores := make(map[uuid.UUID]float64)

	// ── 2. Graph neighbours (graph × 0.4) ──
	// For each entity, do a 1-hop traversal to find connected entities.
	// ponytail: 1-hop only — deeper traversal adds latency with diminishing returns.
	for _, e := range entities {
		nodes, _, err := s.graphRepo.NHopTraversal(ctx, graphName, e.ID.String(), 1)
		if err != nil {
			continue // skip entities with no graph presence
		}
		for _, node := range nodes {
			// Parse entity ID from graph node
			nid, parseErr := uuid.Parse(node.ID)
			if parseErr != nil {
				continue
			}
			graphScores[nid] += 0.4
			if _, exists := candidateMap[nid]; !exists {
				candidateMap[nid] = &models.RecallItem{
					EntityID: nid,
					Source:   "graph",
				}
			}
		}
	}

	// ── 3. Recent mentions (recency × 0.3) ──
	maxScore := 0.0
	for _, e := range entities {
		if e.RelevanceScore > maxScore {
			maxScore = e.RelevanceScore
		}
	}
	for _, e := range entities {
		normScore := 0.0
		if maxScore > 0 {
			normScore = e.RelevanceScore / maxScore
		}
		recencyScore := normScore * 0.3
		if item, exists := candidateMap[e.ID]; exists {
			item.Score += recencyScore
			if item.Fact == "" {
				item.Fact = fmt.Sprintf("Recently active entity: %s", e.Name)
			}
		} else {
			candidateMap[e.ID] = &models.RecallItem{
				EntityID: e.ID,
				Fact:     fmt.Sprintf("Recently active entity: %s", e.Name),
				Score:    recencyScore,
				Source:   "mention",
			}
		}
	}

	// ── 4. Freshness (vector similarity × 0.3) ──
	// FindSimilarEntity already scans every entity's embedding in the universe
	// and returns the single closest match to queryEmbedding, so this runs once.
	if len(queryEmbedding) > 0 {
		matchedID, distance, err := s.vectorRepo.FindSimilarEntity(ctx, universeID, queryEmbedding, 0.8)
		if err == nil && matchedID != nil {
			// Convert distance to similarity score: 1 - distance (cosine distance)
			freshScore := (1.0 - distance) * 0.3
			if item, exists := candidateMap[*matchedID]; exists {
				item.Score += freshScore
			} else {
				candidateMap[*matchedID] = &models.RecallItem{
					EntityID: *matchedID,
					Score:    freshScore,
					Source:   "freshness",
				}
			}
		}
	}

	// ── 5. Merge graph scores into candidates ──
	for id, gs := range graphScores {
		if item, exists := candidateMap[id]; exists {
			item.Score += gs
		}
	}

	// ── 6. Sort by score descending and limit to k ──
	items := make([]models.RecallItem, 0, len(candidateMap))
	for _, item := range candidateMap {
		items = append(items, *item)
	}

	sort.Slice(items, func(i, j int) bool {
		if items[i].Score == items[j].Score {
			return items[i].EntityID.String() < items[j].EntityID.String()
		}
		return items[i].Score > items[j].Score
	})

	if k > 0 && len(items) > k {
		items = items[:k]
	}

	return items, nil
}
