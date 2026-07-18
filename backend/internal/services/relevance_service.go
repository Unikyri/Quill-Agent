package services

import (
	"context"
	"log"
	"math"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/quill/backend/internal/models"
	"github.com/quill/backend/internal/repositories"
)

const mentionRelevanceBump = 0.15

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
	consolidationSvc Consolidator
	historyRepo      *repositories.EntityRelevanceHistoryRepo
}

// NewRelevanceService creates a relevance service with the given decay lambda
// and archive threshold. lambda controls the decay rate per chapter advance;
// archiveThreshold is the score below which entities get archived.
// consolidationSvc may be nil — deconsolidation on reactivate will be skipped.
// historyRepo defaults to a pool-backed instance; override via
// SetHistoryRepo in tests if needed.
func NewRelevanceService(pool *pgxpool.Pool, entityRepo *repositories.EntityRepo, lambda, archiveThreshold float64, consolidationSvc Consolidator) *RelevanceService {
	return &RelevanceService{
		pool:             pool,
		entityRepo:       entityRepo,
		lambda:           lambda,
		archiveThreshold: archiveThreshold,
		consolidationSvc: consolidationSvc,
		historyRepo:      repositories.NewEntityRelevanceHistoryRepo(pool),
	}
}

// Touch records a real entity mention. It bumps active relevance, restores an
// archived entity to active memory, snapshots the mutation for the timeline,
// and synchronizes AGE's denormalized node best-effort.
func (s *RelevanceService) Touch(ctx context.Context, entityID, chapterID uuid.UUID) error {
	before, err := s.entityRepo.FindByID(ctx, entityID)
	if err != nil {
		return err
	}
	updated, err := s.entityRepo.ReinforceMention(ctx, entityID, chapterID, mentionRelevanceBump)
	if err != nil {
		return err
	}
	if updated == nil {
		return nil
	}
	if s.historyRepo != nil {
		if err := s.historyRepo.AppendOne(ctx, updated.ID); err != nil {
			log.Printf("[relevance] append history for mentioned entity %s: %v", updated.ID, err)
		}
	}
	s.syncGraphState(ctx, updated)
	if before.Status == "archived" && s.consolidationSvc != nil {
		if err := s.consolidationSvc.DeconsolidateEntity(ctx, updated.ID); err != nil {
			log.Printf("[relevance] deconsolidate mentioned entity %s: %v", updated.ID, err)
		}
	}
	return nil
}

// Reactivate sets the entity's score to 0.8 and status to "active".
// Used when a previously archived entity becomes relevant again.
// Also triggers deconsolidation: removes the consolidated memory row
// since the entity is no longer fully archived.
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

	if err := tx.Commit(ctx); err != nil {
		return err
	}

	// spec: reactivation writes one entity_relevance_history row
	if s.historyRepo != nil {
		if err := s.historyRepo.AppendOne(ctx, entityID); err != nil {
			log.Printf("[relevance] append history for reactivated entity %s: %v", entityID, err)
		}
	}
	s.syncGraphState(ctx, e)

	// spec: after reactivation, deconsolidate (nil-safe)
	if s.consolidationSvc != nil {
		if err := s.consolidationSvc.DeconsolidateEntity(ctx, entityID); err != nil {
			log.Printf("[relevance] deconsolidate entity %s: %v", entityID, err)
		}
	}

	return nil
}

// DecayAll applies exponential decay (score *= e^-lambda) to all active entities
// in the universe. Entities whose score drops below archiveThreshold are set to
// status "archived". After archiving, newly-archived entities are consolidated
// asynchronously (fire-and-forget, errors logged).
//
// ponytail: single-pass strategy — decay and archive in one call; no separate
// archive sweep needed.
func (s *RelevanceService) DecayAll(ctx context.Context, universeID uuid.UUID) error {
	return s.DecayExcept(ctx, universeID, nil)
}

// DecayForChapter advances relevance after a completed chapter. The entities
// mentioned in that chapter are explicitly excluded: a fresh mention must not
// lose relevance in the exact event that records it.
func (s *RelevanceService) DecayForChapter(ctx context.Context, universeID, chapterID uuid.UUID) error {
	mentioned, err := s.entityRepo.ListMentionedEntityIDs(ctx, chapterID)
	if err != nil {
		return err
	}
	return s.DecayExcept(ctx, universeID, mentioned)
}

// DecayExcept is the chapter-aware decay primitive. DecayAll remains for an
// explicit manual maintenance sweep; production chapter flows call this with
// the mentions from the completed chapter.
func (s *RelevanceService) DecayExcept(ctx context.Context, universeID uuid.UUID, mentionedIDs []uuid.UUID) error {
	if err := s.entityRepo.DecayExcept(ctx, universeID, mentionedIDs, s.lambda); err != nil {
		return err
	}

	// Identify entities about to be archived BEFORE the UPDATE
	newlyArchivedIDs, err := s.entityRepo.FindNewlyArchivable(ctx, universeID, s.archiveThreshold)
	if err != nil {
		log.Printf("[relevance] find newly archivable: %v", err)
	}

	// Archive entities that fell below threshold
	_, err = s.pool.Exec(ctx, `
		UPDATE entities SET status = 'archived', updated_at = NOW()
		WHERE universe_id = $1 AND status = 'active' AND relevance_score <= $2
	`, universeID, s.archiveThreshold)
	if err != nil {
		return err
	}

	// spec: piggyback one set-based snapshot INSERT...SELECT after the bulk
	// decay+archive UPDATEs, capturing every entity's post-decay score/status.
	if s.historyRepo != nil {
		if err := s.historyRepo.AppendSnapshot(ctx, universeID); err != nil {
			log.Printf("[relevance] append decay snapshot for universe %s: %v", universeID, err)
		}
	}

	// Fire-and-forget consolidation goroutines for newly-archived entities
	if s.consolidationSvc != nil && len(newlyArchivedIDs) > 0 {
		for _, entityID := range newlyArchivedIDs {
			go func(eid uuid.UUID) {
				defer func() {
					if r := recover(); r != nil {
						log.Printf("[relevance] panic consolidating entity %s: %v", eid, r)
					}
				}()
				if err := s.consolidationSvc.ConsolidateEntity(context.Background(), eid, universeID); err != nil {
					log.Printf("[relevance] consolidate entity %s: %v", eid, err)
				}
			}(entityID)
		}
	}

	states, err := s.entityRepo.ListRelevanceStates(ctx, universeID)
	if err != nil {
		log.Printf("[relevance] list graph relevance states for universe %s: %v", universeID, err)
		return nil
	}
	for index := range states {
		s.syncGraphState(ctx, &states[index])
	}

	return nil
}

func (s *RelevanceService) syncGraphState(ctx context.Context, entity *models.Entity) {
	if s.pool == nil || entity == nil {
		return
	}
	graphRepo := repositories.NewGraphRepo(s.pool)
	if err := graphRepo.UpdateNodeState(ctx, "universe_"+entity.UniverseID.String(), entity.ID.String(), entity.RelevanceScore, entity.Status); err != nil {
		// SQL remains the source of truth. Graph synchronization is presentation
		// enrichment and must never reject a successful writing operation.
		log.Printf("[relevance] sync graph node %s: %v", entity.ID, err)
	}
}

// applyDecay computes the decayed relevance score after idle chapters.
// score *= e^(-lambda * idle). Exported for testing.
func applyDecay(score float64, idleChapters float64, lambda float64) float64 {
	return score * math.Exp(-lambda*idleChapters)
}
