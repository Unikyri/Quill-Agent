package services

import (
	"context"
	"fmt"
	"log"
	"sort"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/quill/backend/internal/models"
	"github.com/quill/backend/internal/repositories"
)

var (
	allowedGenreTags = map[string]struct{}{
		"fantasy": {}, "epic-fantasy": {}, "urban-fantasy": {}, "romantasy": {},
		"science-fiction": {}, "space-opera": {}, "dystopian": {}, "horror": {},
		"gothic": {}, "paranormal": {}, "romance": {}, "mystery": {},
		"cozy-mystery": {}, "thriller": {}, "crime": {}, "historical": {},
		"literary": {}, "adventure": {}, "young-adult": {}, "coming-of-age": {},
	}
)

func validateUniverseEnums(input models.CreateUniverseRequest) error {
	for _, tag := range input.GenreTags {
		if _, ok := allowedGenreTags[tag]; !ok {
			return fmt.Errorf("invalid genre tag %q: must be one of %s", tag, joinKeys(allowedGenreTags))
		}
	}
	return nil
}

func joinKeys(m map[string]struct{}) string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return strings.Join(keys, ", ")
}

type UniverseService struct {
	pool         *pgxpool.Pool
	universeRepo *repositories.UniverseRepo
	graphRepo    *repositories.GraphRepo
}

func NewUniverseService(pool *pgxpool.Pool, universeRepo *repositories.UniverseRepo, graphRepo *repositories.GraphRepo) *UniverseService {
	return &UniverseService{
		pool:         pool,
		universeRepo: universeRepo,
		graphRepo:    graphRepo,
	}
}

func (s *UniverseService) Create(ctx context.Context, userID uuid.UUID, input models.CreateUniverseRequest) (*models.Universe, error) {
	if input.Name == "" {
		return nil, fmt.Errorf("universe name is required")
	}
	if err := validateUniverseEnums(input); err != nil {
		return nil, err
	}

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	u := &models.Universe{
		ID:          uuid.New(),
		UserID:      userID,
		Name:        input.Name,
		Description: input.Description,
		GenreTags:   input.GenreTags,
	}

	if err := s.universeRepo.Create(ctx, tx, u); err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("commit transaction: %w", err)
	}

	// Create AGE graph for the new universe
	if s.graphRepo != nil {
		if err := s.graphRepo.CreateGraph(ctx, u.ID.String()); err != nil {
			log.Printf("[universe] create AGE graph for %s: %v", u.ID, err)
			// non-fatal — the graph can be created later
		}
	}

	return u, nil
}

func (s *UniverseService) GetByID(ctx context.Context, id uuid.UUID) (*models.Universe, error) {
	return s.universeRepo.FindByID(ctx, id)
}

func (s *UniverseService) ListByUser(ctx context.Context, userID uuid.UUID, page, limit int) ([]models.Universe, int, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	return s.universeRepo.ListByUser(ctx, userID, page, limit)
}

func (s *UniverseService) Update(ctx context.Context, id uuid.UUID, input models.CreateUniverseRequest) (*models.Universe, error) {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	u, err := s.universeRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if err := validateUniverseEnums(input); err != nil {
		return nil, err
	}

	if input.Name != "" {
		u.Name = input.Name
	}
	if input.Description != "" {
		u.Description = input.Description
	}
	if input.GenreTags != nil {
		u.GenreTags = input.GenreTags
	}

	if err := s.universeRepo.Update(ctx, tx, u); err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("commit transaction: %w", err)
	}

	return u, nil
}

func (s *UniverseService) Delete(ctx context.Context, id uuid.UUID) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	if err := s.universeRepo.Delete(ctx, tx, id); err != nil {
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return err
	}

	// Best-effort: drop the universe's AGE graph so it doesn't leak into
	// ag_catalog forever. The relational delete already committed.
	if s.graphRepo != nil {
		if err := s.graphRepo.DropGraph(ctx, "universe_"+id.String()); err != nil {
			log.Printf("[universe] drop graph for %s: %v", id, err)
		}
	}
	return nil
}
