package services

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/quill/backend/internal/models"
	"github.com/quill/backend/internal/repositories"
)

// TimelineService validates temporal consistency of timeline events.
type TimelineService struct {
	pool         *pgxpool.Pool
	timelineRepo *repositories.TimelineRepo
}

func NewTimelineService(pool *pgxpool.Pool, timelineRepo *repositories.TimelineRepo) *TimelineService {
	return &TimelineService{pool: pool, timelineRepo: timelineRepo}
}

// ValidatePosition checks whether the event's chapter is chronologically before or
// at the present chapter. Returns an error if the event chapter is in the future.
func (s *TimelineService) ValidatePosition(ctx context.Context, event models.TimelineEvent, presentChapterID uuid.UUID) error {
	// ponytail: compare chapter order_index; nil chapter_id = always valid
	if event.ChapterID == nil {
		return nil
	}

	var eventOrder, presentOrder int
	err := s.pool.QueryRow(ctx, "SELECT order_index FROM chapters WHERE id = $1", *event.ChapterID).Scan(&eventOrder)
	if err != nil {
		return fmt.Errorf("get event chapter: %w", err)
	}
	err = s.pool.QueryRow(ctx, "SELECT order_index FROM chapters WHERE id = $1", presentChapterID).Scan(&presentOrder)
	if err != nil {
		return fmt.Errorf("get present chapter: %w", err)
	}

	if eventOrder > presentOrder {
		return fmt.Errorf("timeline event chapter (%d) is after present chapter (%d): future event not allowed", eventOrder, presentOrder)
	}
	return nil
}
