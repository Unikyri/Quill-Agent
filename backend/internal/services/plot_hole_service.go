package services

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/quill/backend/internal/models"
	"github.com/quill/backend/internal/repositories"
)

// PlotHoleService detects narrative arcs that have been stale for too many chapters.
type PlotHoleService struct {
	pool         *pgxpool.Pool
	plotHoleRepo *repositories.PlotHoleRepo
	entityRepo   *repositories.EntityRepo
	chapters     int // inactivity threshold
}

func NewPlotHoleService(pool *pgxpool.Pool, plotHoleRepo *repositories.PlotHoleRepo, entityRepo *repositories.EntityRepo, chapters int) *PlotHoleService {
	return &PlotHoleService{
		pool:         pool,
		plotHoleRepo: plotHoleRepo,
		entityRepo:   entityRepo,
		chapters:     chapters,
	}
}

// Scan checks all active entities and creates plot holes for those whose
// last_mentioned chapter is at least `chapters` behind the current chapter.
//
// ponytail: single-pass scan over active entities; O(n) where n is entity count.
// Gap calculation: current_order - last_mentioned_order.
// ponytail: N+1 query per entity for chapter lookup — fine for hackathon scale;
// batch preload chapters if entity count exceeds 1k.
func (s *PlotHoleService) Scan(ctx context.Context, universeID, currentChapterID uuid.UUID) ([]models.PlotHole, error) {
	// Get current chapter's order_index
	var currentOrder int
	if err := s.pool.QueryRow(ctx, "SELECT order_index FROM chapters WHERE id = $1", currentChapterID).Scan(&currentOrder); err != nil {
		return nil, fmt.Errorf("get current chapter order: %w", err)
	}

	// Get all active entities
	entities, err := s.entityRepo.ListByUniverseActive(ctx, universeID)
	if err != nil {
		return nil, fmt.Errorf("list active entities: %w", err)
	}

	var holes []models.PlotHole
	for _, e := range entities {
		if e.LastMentionedChapterID == nil {
			continue // never mentioned → skip
		}

		// Get the last mentioned chapter's order_index
		var lastOrder int
		if err := s.pool.QueryRow(ctx, "SELECT order_index FROM chapters WHERE id = $1", *e.LastMentionedChapterID).Scan(&lastOrder); err != nil {
			continue // chapter may have been deleted
		}

		gap := currentOrder - lastOrder
		if gap >= s.chapters {
			hole := models.PlotHole{
				ID:                       uuid.New(),
				UniverseID:               universeID,
				Title:                    fmt.Sprintf("Stale arc: %s (gap %d chapters)", e.Name, gap),
				Description:              fmt.Sprintf("Entity '%s' has not been mentioned for %d chapters (last seen in chapter %d, currently at chapter %d)", e.Name, gap, lastOrder, currentOrder),
				RelatedEntityIDs:         []uuid.UUID{e.ID},
				FirstMentionedChapterID:  e.LastMentionedChapterID,
				Status:                   "open",
			}
			if err := s.plotHoleRepo.Create(ctx, &hole); err != nil {
				return nil, fmt.Errorf("create plot hole: %w", err)
			}
			holes = append(holes, hole)
		}
	}

	return holes, nil
}
