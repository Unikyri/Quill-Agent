package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/quill/backend/internal/models"
	"github.com/quill/backend/internal/repositories"
)

// analysisJob represents a single paragraph to analyze.
type analysisJob struct {
	WorkID     uuid.UUID
	ChapterID  uuid.UUID
	UniverseID uuid.UUID
	Text       string
	UserID     uuid.UUID
}

// AnalysisResult holds the output of a complete analysis pass.
type AnalysisResult struct {
	WorkID         uuid.UUID
	ChapterID      uuid.UUID
	Entities       []models.EntityBrief
	Contradictions []models.Contradiction
	PlotHoles      []models.PlotHole
}

// AnalysisHub is the minimal WebSocket hub interface used by AnalysisService.
// ws.Hub satisfies this interface via its SendToUser method.
type AnalysisHub interface {
	SendToUser(userID uuid.UUID, msg models.WSMessage) error
}

// AnalysisService runs a per-work sequential analysis queue.
//
// ponytail: one goroutine per work, sequential queue. No worker pool needed
// for hackathon scale. Cancel/Shutdown stop the goroutine.
type AnalysisService struct {
	pool        *pgxpool.Pool
	entitySvc   *EntityService
	contraSvc   *ContradictionService
	relevSvc    *RelevanceService
	timelineSvc *TimelineService
	plotHoleSvc *PlotHoleService
	qwenSvc     *QwenService
	hub         AnalysisHub

	queues  map[uuid.UUID]chan analysisJob
	cancels map[uuid.UUID]context.CancelFunc
	mu      sync.Mutex
}

// NewAnalysisService creates an analysis service. All parameters may be nil
// for testing; Submit will only enqueue. Workers start via runWorker.
func NewAnalysisService(
	pool *pgxpool.Pool,
	entitySvc *EntityService,
	contraSvc *ContradictionService,
	relevSvc *RelevanceService,
	timelineSvc *TimelineService,
	plotHoleSvc *PlotHoleService,
	qwenSvc *QwenService,
	hub AnalysisHub,
) *AnalysisService {
	return &AnalysisService{
		pool:        pool,
		entitySvc:   entitySvc,
		contraSvc:   contraSvc,
		relevSvc:    relevSvc,
		timelineSvc: timelineSvc,
		plotHoleSvc: plotHoleSvc,
		qwenSvc:     qwenSvc,
		hub:         hub,
		queues:      make(map[uuid.UUID]chan analysisJob),
		cancels:     make(map[uuid.UUID]context.CancelFunc),
	}
}

// SubmitParagraph is a convenience wrapper that satisfies the ws.ParagraphSubmitter
// interface. It creates an analysisJob and enqueues it via Submit.
// Starts a worker goroutine for first-time work submissions.
func (s *AnalysisService) SubmitParagraph(ctx context.Context, workID, chapterID, universeID, userID uuid.UUID, text string) error {
	s.mu.Lock()
	_, exists := s.queues[workID]
	s.mu.Unlock()

	job := analysisJob{
		WorkID:     workID,
		ChapterID:  chapterID,
		UniverseID: universeID,
		Text:       text,
		UserID:     userID,
	}

	if err := s.Submit(ctx, job); err != nil {
		return err
	}

	// Start a worker if this is a new work ID
	if !exists {
		go s.runWorker(workID)
	}

	return nil
}

// Submit enqueues an analysis job into the per-work channel.
func (s *AnalysisService) Submit(ctx context.Context, job analysisJob) error {
	s.mu.Lock()
	q, exists := s.queues[job.WorkID]
	if !exists {
		q = make(chan analysisJob, 100)
		s.queues[job.WorkID] = q
	}
	s.mu.Unlock()

	select {
	case q <- job:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// Cancel stops the worker goroutine for the given workID.
func (s *AnalysisService) Cancel(workID uuid.UUID) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if cancel, exists := s.cancels[workID]; exists {
		cancel()
		delete(s.cancels, workID)
	}
	delete(s.queues, workID)
}

// Shutdown cancels all running workers and removes all queues.
func (s *AnalysisService) Shutdown() {
	s.mu.Lock()
	defer s.mu.Unlock()

	for workID, cancel := range s.cancels {
		cancel()
		delete(s.cancels, workID)
	}
	for workID := range s.queues {
		delete(s.queues, workID)
	}
}

// runWorker starts a goroutine that drains the per-work queue sequentially.
func (s *AnalysisService) runWorker(workID uuid.UUID) {
	s.mu.Lock()
	if _, exists := s.cancels[workID]; exists {
		s.mu.Unlock()
		return
	}
	workerCtx, cancel := context.WithCancel(context.Background())
	s.cancels[workID] = cancel
	s.mu.Unlock()

	defer func() {
		s.mu.Lock()
		delete(s.cancels, workID)
		s.mu.Unlock()
	}()

	s.mu.Lock()
	q, exists := s.queues[workID]
	s.mu.Unlock()
	if !exists {
		return
	}

	for {
		select {
		case job, ok := <-q:
			if !ok {
				return
			}
			result, err := s.processJob(workerCtx, job)
			if err != nil {
				log.Printf("[analysis] work %s job failed: %v", workID, err)
				continue
			}
			if s.hub != nil && result != nil {
				s.broadcastResult(job.UserID, *result)
			}
		case <-workerCtx.Done():
			return
		}
	}
}

// processJob runs the full analysis pipeline for a single paragraph.
//
// ponytail: sequential pipeline — core pass then enrichment pass.
func (s *AnalysisService) processJob(ctx context.Context, job analysisJob) (*AnalysisResult, error) {
	if ctx.Err() != nil {
		return nil, fmt.Errorf("analysis context cancelled: %w", ctx.Err())
	}

	result := &AnalysisResult{
		WorkID:    job.WorkID,
		ChapterID: job.ChapterID,
	}

	// ── Core Pass (deterministic, fast) ──

	// 1. Extract entities from paragraph text
	var resolvedEntities []ResolvedEntity
	if s.entitySvc != nil && s.pool != nil {
		entities, err := s.extractEntities(ctx, job.UniverseID, job.Text, job.ChapterID)
		if err != nil {
			log.Printf("[analysis] extract entities: %v", err)
		} else {
			resolvedEntities = entities
			for _, re := range resolvedEntities {
				result.Entities = append(result.Entities, models.EntityBrief{
					ID:   re.Entity.ID,
					Name: re.Entity.Name,
					Type: re.Entity.Type,
				})
			}
		}
	}

	// 2. Deterministic contradiction checks (deceased/alive rules)
	if s.contraSvc != nil && len(resolvedEntities) > 0 {
		deterministic, err := s.contraSvc.CheckDeterministic(ctx, job.UniverseID, job.ChapterID, resolvedEntities)
		if err != nil {
			log.Printf("[analysis] deterministic check: %v", err)
		} else {
			result.Contradictions = append(result.Contradictions, deterministic...)
		}
	}

	// 3. Touch relevance for each mentioned entity
	if s.relevSvc != nil {
		for _, re := range resolvedEntities {
			if err := s.relevSvc.Touch(ctx, re.Entity.ID, job.ChapterID); err != nil {
				log.Printf("[analysis] touch entity %s: %v", re.Entity.ID, err)
			}
		}
	}

	// ── Enrichment Pass (Qwen-Max) ──

	// 4. Semantic contradiction checks via Qwen-Max
	if s.contraSvc != nil && len(resolvedEntities) > 0 {
		semantic, err := s.contraSvc.CheckSemantic(ctx, job.UniverseID, job.ChapterID, job.Text, resolvedEntities)
		if err != nil {
			log.Printf("[analysis] semantic check: %v", err)
		} else {
			result.Contradictions = append(result.Contradictions, semantic...)
		}
	}

	// 5. Scan for plot holes
	if s.plotHoleSvc != nil {
		holes, err := s.plotHoleSvc.Scan(ctx, job.UniverseID, job.ChapterID)
		if err != nil {
			log.Printf("[analysis] plot hole scan: %v", err)
		} else {
			result.PlotHoles = holes
		}
	}

	log.Printf("[analysis] work=%s chapter=%s: %d entities, %d contradictions, %d plot holes",
		job.WorkID, job.ChapterID, len(result.Entities), len(result.Contradictions), len(result.PlotHoles))

	return result, nil
}

// extractEntities resolves or creates entities from paragraph text.
//
// ponytail: best-effort; if QwenService is nil or extraction fails, returns
// empty slice rather than failing the whole job.
func (s *AnalysisService) extractEntities(ctx context.Context, universeID uuid.UUID, text string, chapterID uuid.UUID) ([]ResolvedEntity, error) {
	if s.qwenSvc == nil || s.entitySvc == nil {
		return nil, nil
	}

	extracted, err := s.qwenSvc.ExtractEntities(ctx, text, "")
	if err != nil {
		return nil, fmt.Errorf("qwen extract: %w", err)
	}

	// Collect all extracted entities from all categories
	allEntities := make([]ExtractedEntity, 0)
	if extracted != nil {
		allEntities = append(allEntities, extracted.Characters...)
		allEntities = append(allEntities, extracted.Places...)
		allEntities = append(allEntities, extracted.Events...)
		allEntities = append(allEntities, extracted.Factions...)
		allEntities = append(allEntities, extracted.WorldRules...)
		allEntities = append(allEntities, extracted.PlotDevelopments...)
	}

	// ponytail: use first 120 chars of text as mention context
	mentionText := text
	if len(mentionText) > 120 {
		mentionText = mentionText[:120]
	}

	var resolved []ResolvedEntity
	for _, ee := range allEntities {
		entityData := repositories.ExtractedEntity{
			Type:        ee.Type,
			Name:        ee.Name,
			Aliases:     ee.Aliases,
			Description: ee.Description,
			Status:      ee.Status,
			Properties:  ee.Properties,
		}
		entity, isNew, err := s.entitySvc.ResolveOrCreate(ctx, universeID, entityData)
		if err != nil {
			log.Printf("[analysis] resolve entity %s: %v", ee.Name, err)
			continue
		}
		resolved = append(resolved, ResolvedEntity{
			Entity:      *entity,
			MentionText: mentionText,
			IsNew:       isNew,
		})
	}

	return resolved, nil
}

// broadcastResult pushes the analysis result to the user's WebSocket connection.
func (s *AnalysisService) broadcastResult(userID uuid.UUID, result AnalysisResult) {
	payloadBytes, err := json.Marshal(map[string]interface{}{
		"work_id":        result.WorkID,
		"chapter_id":     result.ChapterID,
		"entities":       result.Entities,
		"contradictions": result.Contradictions,
		"plot_holes":     result.PlotHoles,
	})
	if err != nil {
		log.Printf("[analysis] marshal result: %v", err)
		return
	}

	msg := models.WSMessage{
		Type:    "analysis_result",
		Payload: payloadBytes,
	}

	if err := s.hub.SendToUser(userID, msg); err != nil {
		log.Printf("[analysis] send result to user %s: %v", userID, err)
	}
}
