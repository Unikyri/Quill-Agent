package services

import (
	"context"
	"math"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/quill/backend/internal/repositories"
)

// RelevanceService manages entity relevance scoring with exponential decay.
//
// ponytail: per-chapter event-driven decay; no background worker needed.
// DecayAll runs on chapter advance; Touch on entity mention; Reactivate on
// manual override or story relevance.
type RelevanceService struct {
	pool             *pgxpool.Pool
	entityRepo       *repositories.EntityRepo
	lambda           float64
	archiveThreshold float64
}

// NewRelevanceService creates a relevance service with the given decay lambda
// and archive threshold. lambda controls the decay rate per chapter advance;
// archiveThreshold is the score below which entities get archived.
func NewRelevanceService(pool *pgxpool.Pool, entityRepo *repositories.EntityRepo, lambda, archiveThreshold float64) *RelevanceService {
	return &RelevanceService{
		pool:             pool,
		entityRepo:       entityRepo,
		lambda:           lambda,
		archiveThreshold: archiveThreshold,
	}
}

// Touch resets the idle counter for a mentioned entity, updating
// last_mentioned_chapter_id and last_mentioned_at.
func (s *RelevanceService) Touch(ctx context.Context, entityID, chapterID uuid.UUID) error {
	return s.entityRepo.TouchBatch(ctx, []uuid.UUID{entityID}, chapterID)
}

// Reactivate sets the entity's score to 0.8 and status to "active".
// Used when a previously archived entity becomes relevant again.
func (s *RelevanceService) Reactivate(ctx context.Context, entityID uuid.UUID) error {
	e, err := s.entityRepo.FindByID(ctx, entityID)
	if err != nil {
		return err
	}

	e.RelevanceScore = 0.8
	e.Status = "active"

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	if err := s.entityRepo.Update(ctx, tx, e); err != nil {
		return err
	}

	return tx.Commit(ctx)
}

// DecayAll applies exponential decay (score *= e^-lambda) to all active entities
// in the universe. Entities whose score drops below archiveThreshold are set to
// status "archived".
//
// ponytail: single-pass strategy — decay and archive in one call; no separate
// archive sweep needed.
func (s *RelevanceService) DecayAll(ctx context.Context, universeID uuid.UUID) error {
	if err := s.entityRepo.DecayAll(ctx, universeID, s.lambda); err != nil {
		return err
	}

	// Archive entities that fell below threshold
	_, err := s.pool.Exec(ctx, `
		UPDATE entities SET status = 'archived', updated_at = NOW()
		WHERE universe_id = $1 AND status = 'active' AND relevance_score <= $2
	`, universeID, s.archiveThreshold)
	if err != nil {
		return err
	}

	return nil
}

// applyDecay computes the decayed relevance score after idle chapters.
// score *= e^(-lambda * idle). Exported for testing.
func applyDecay(score float64, idleChapters float64, lambda float64) float64 {
	return score * math.Exp(-lambda*idleChapters)
}
