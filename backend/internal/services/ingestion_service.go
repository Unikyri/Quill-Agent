package services

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"regexp"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/quill/backend/internal/models"
	"github.com/quill/backend/internal/repositories"
)

// ── Ingestion service ──

// ingestionChunk represents a parsed section of the uploaded document.
type ingestionChunk struct {
	title   string
	content string
}

// IngestionQwen is the minimal Qwen interface used by IngestionService.
// QwenService satisfies this interface.
type IngestionQwen interface {
	ExtractEntities(ctx context.Context, text, universeContext string) (*ExtractedEntities, error)
	GenerateEmbedding(ctx context.Context, text string) ([]float32, error)
	GenerateEmbeddingBatch(ctx context.Context, texts []string) ([][]float32, error)
}

// IngestionService processes document uploads asynchronously:
// file → chunk by headers → extract entities → embed → graph.
//
// ponytail: one goroutine per job, sequential chunk processing. No worker pool
// needed for hackathon scale. Cancel via ctx.Done().
type IngestionService struct {
	pool       *pgxpool.Pool
	entitySvc  *EntityService
	vectorRepo *repositories.VectorRepo
	graphRepo  *repositories.GraphRepo
	qwenSvc    IngestionQwen
	hub        AnalysisHub
}

// NewIngestionService creates an IngestionService. All parameters may be nil
// for testing; Start will create a job ID but the worker will be a no-op.
func NewIngestionService(
	pool *pgxpool.Pool,
	entitySvc *EntityService,
	vectorRepo *repositories.VectorRepo,
	graphRepo *repositories.GraphRepo,
	qwenSvc IngestionQwen,
	hub AnalysisHub,
) *IngestionService {
	return &IngestionService{
		pool:       pool,
		entitySvc:  entitySvc,
		vectorRepo: vectorRepo,
		graphRepo:  graphRepo,
		qwenSvc:    qwenSvc,
		hub:        hub,
	}
}

// Start creates an ingestion job and kicks off the async pipeline.
// Returns the job ID immediately; duplicate is true when the same content
// was already ingested into this universe (the existing job's ID is
// returned and no worker is started). The caller should return 202 Accepted
// for new jobs and 200 for duplicates.
func (s *IngestionService) Start(ctx context.Context, universeID uuid.UUID, reader io.Reader, filename string) (uuid.UUID, bool, error) {
	jobID := uuid.New()

	// ponytail: read the full content synchronously before spawning the
	// goroutine. The handler's file.Close() runs as soon as Start returns, so
	// passing the io.Reader to a goroutine would read from a closed handle.
	content, err := io.ReadAll(reader)
	if err != nil {
		return uuid.Nil, false, fmt.Errorf("read uploaded file: %w", err)
	}

	sum := sha256.Sum256(content)
	hash := hex.EncodeToString(sum[:])

	var workID uuid.UUID
	if s.pool != nil {
		repo := repositories.NewIngestionRepo(s.pool)
		existing, err := repo.FindByContentHash(ctx, universeID, hash)
		if err != nil {
			return uuid.Nil, false, fmt.Errorf("check duplicate ingestion: %w", err)
		}
		if existing != nil {
			return existing.ID, true, nil
		}

		workRepo := repositories.NewWorkRepo(s.pool)
		works, err := workRepo.ListByUniverse(ctx, universeID)
		if err != nil {
			return uuid.Nil, false, fmt.Errorf("resolve work: %w", err)
		}
		if len(works) > 0 {
			// ponytail: ingest into the first work. A future UI should let users
			// pick which work to target when a universe has more than one.
			workID = works[0].ID
		} else {
			tx, err := s.pool.Begin(ctx)
			if err != nil {
				return uuid.Nil, false, fmt.Errorf("begin transaction: %w", err)
			}
			orderIdx, err := workRepo.GetMaxOrderIndex(ctx, universeID)
			if err != nil {
				_ = tx.Rollback(ctx)
				return uuid.Nil, false, fmt.Errorf("get max order index: %w", err)
			}
			work := models.Work{
				ID:         uuid.New(),
				UniverseID: universeID,
				Title:      "Imported Manuscript",
				Type:       "novel",
				Status:     "in_progress",
				OrderIndex: orderIdx + 1,
			}
			if err := workRepo.Create(ctx, tx, &work); err != nil {
				_ = tx.Rollback(ctx)
				return uuid.Nil, false, fmt.Errorf("create default work: %w", err)
			}
			if err := tx.Commit(ctx); err != nil {
				return uuid.Nil, false, fmt.Errorf("commit transaction: %w", err)
			}
			workID = work.ID
		}

		if err := repo.Create(ctx, jobID, universeID, workID, "pending", filename, hash); err != nil {
			// Unique violation: another upload of the same content won the
			// race between our FindByContentHash and this insert.
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) && pgErr.Code == "23505" {
				if existing, ferr := repo.FindByContentHash(ctx, universeID, hash); ferr == nil && existing != nil {
					return existing.ID, true, nil
				}
			}
			return uuid.Nil, false, fmt.Errorf("create ingestion job: %w", err)
		}
	}

	go s.runWorker(jobID, universeID, workID, content, filename)

	return jobID, false, nil
}

// ListJobs returns the recent ingestion jobs for a universe.
func (s *IngestionService) ListJobs(ctx context.Context, universeID uuid.UUID) ([]models.IngestionJob, error) {
	if s.pool == nil {
		return []models.IngestionJob{}, nil
	}
	return repositories.NewIngestionRepo(s.pool).ListByUniverse(ctx, universeID)
}

// runWorker processes the document in a background goroutine.
//
// ponytail: synchronous per-chunk — no parallel chunk extraction to avoid
// overwhelming the Qwen API rate limit.
func (s *IngestionService) runWorker(jobID, universeID, workID uuid.UUID, content []byte, filename string) {
	ctx := context.Background()

	s.updateJobStatus(ctx, jobID, "running", "")

	// Resolve the universe owner once per job — this never changes during a
	// job's lifetime, so N+1 identical lookups per emit would be wasteful.
	// Failure (deleted universe, or pool==nil in unit tests) degrades to
	// best-effort: ownerID stays uuid.Nil and progress simply won't be routed.
	var ownerID uuid.UUID
	if s.pool != nil {
		if u, err := repositories.NewUniverseRepo(s.pool).FindByID(ctx, universeID); err != nil {
			log.Printf("[ingestion] resolve universe owner %s: %v (progress events will not be delivered)", universeID, err)
		} else {
			ownerID = u.UserID
		}
	}

	// Split content by markdown headers
	chunks := s.splitChunks(string(content))
	if len(chunks) == 0 {
		s.updateJobStatus(ctx, jobID, "completed", "")
		return
	}

	// One real chapters row per chunk, under the imported work, so paragraph
	// embeddings get a valid chapter FK and the document survives reloads.
	var chRepo *repositories.ChapterRepo
	baseOrder := 0
	if s.pool != nil {
		chRepo = repositories.NewChapterRepo(s.pool)
		if bo, err := chRepo.GetMaxOrderIndex(ctx, workID); err != nil {
			log.Printf("[ingestion] get max chapter order for work %s: %v", workID, err)
		} else {
			baseOrder = bo
		}
	}

	entitiesTotal := 0
	s.updateProgress(ctx, jobID, len(chunks), 0, entitiesTotal)
	s.emitProgress(jobID, ownerID, "running", 0, len(chunks))

	anySucceeded := false
	var lastErr error

	for i, ch := range chunks {
		select {
		case <-ctx.Done():
			s.updateJobStatus(ctx, jobID, "failed", "cancelled")
			return
		default:
		}

		chapterID := uuid.Nil
		if chRepo != nil {
			chapter := models.Chapter{
				ID:         uuid.New(),
				WorkID:     workID,
				Title:      ch.title,
				OrderIndex: baseOrder + i + 1,
				Content:    ch.content,
				RawText:    ch.content,
				WordCount:  chRepo.CountWords(ch.content),
				Status:     "draft",
			}
			if err := s.createChapter(ctx, chRepo, &chapter); err != nil {
				// Without a valid chapter FK there is nothing to persist for
				// this chunk — skip it entirely.
				log.Printf("[ingestion] create chapter chunk %d: %v", i, err)
				s.updateProgress(ctx, jobID, len(chunks), i+1, entitiesTotal)
				s.emitProgress(jobID, ownerID, "running", i+1, len(chunks))
				continue
			}
			chapterID = chapter.ID
		}

		// Embed the chunk's paragraphs in batches of 10 (DashScope
		// text-embedding-v3 batch limit). Batch failure → log + skip slice,
		// same best-effort semantics as the old per-paragraph calls.
		if s.qwenSvc != nil && s.vectorRepo != nil {
			var texts []string
			var indexes []int
			for pIdx, p := range strings.Split(ch.content, "\n\n") {
				p = strings.TrimSpace(p)
				if p == "" {
					continue
				}
				const maxEmbedChars = 30_000
				if len(p) > maxEmbedChars {
					log.Printf("[ingestion] skip embedding oversized paragraph chunk %d para %d (%d chars)", i, pIdx, len(p))
					continue
				}
				texts = append(texts, p)
				indexes = append(indexes, pIdx)
			}

			const embedBatchSize = 10
			for start := 0; start < len(texts); start += embedBatchSize {
				end := start + embedBatchSize
				if end > len(texts) {
					end = len(texts)
				}
				embeddings, err := s.qwenSvc.GenerateEmbeddingBatch(ctx, texts[start:end])
				if err != nil {
					log.Printf("[ingestion] embed batch chunk %d paras %d-%d: %v", i, start, end-1, err)
					continue
				}
				if len(embeddings) != end-start {
					log.Printf("[ingestion] embed batch chunk %d paras %d-%d: got %d embeddings for %d texts", i, start, end-1, len(embeddings), end-start)
				}
				for j, emb := range embeddings {
					if emb == nil {
						continue
					}
					if err := s.vectorRepo.SaveParagraphEmbedding(ctx, chapterID, indexes[start+j], ch.title, texts[start+j], emb); err != nil {
						log.Printf("[ingestion] save paragraph embedding chunk %d para %d: %v", i, indexes[start+j], err)
					}
				}
			}
		}

		// Extract entities from chunk
		if s.qwenSvc != nil && s.entitySvc != nil && s.pool != nil {
			extracted, err := s.qwenSvc.ExtractEntities(ctx, ch.content, "")
			if err != nil {
				log.Printf("[ingestion] extract entities chunk %d: %v", i, err)
				lastErr = err
				s.updateProgress(ctx, jobID, len(chunks), i+1, entitiesTotal)
				s.emitProgress(jobID, ownerID, "running", i+1, len(chunks))
				continue
			}
			anySucceeded = true
			entitiesTotal += s.resolveAndBuildGraph(ctx, universeID, extracted)
		}

		s.updateProgress(ctx, jobID, len(chunks), i+1, entitiesTotal)
		s.emitProgress(jobID, ownerID, "running", i+1, len(chunks))
	}

	if !anySucceeded && lastErr != nil {
		s.updateJobStatus(ctx, jobID, "failed", fmt.Sprintf("entity extraction failed for all %d chunks", len(chunks)))
		s.updateProgress(ctx, jobID, len(chunks), len(chunks), entitiesTotal)
		s.emitProgress(jobID, ownerID, "failed", len(chunks), len(chunks))
		return
	}

	s.updateJobStatus(ctx, jobID, "completed", "")
	s.updateProgress(ctx, jobID, len(chunks), len(chunks), entitiesTotal)
	s.emitProgress(jobID, ownerID, "completed", len(chunks), len(chunks))
}

// createChapter wraps ChapterRepo.Create (which requires a transaction) in a
// short single-statement transaction.
func (s *IngestionService) createChapter(ctx context.Context, chRepo *repositories.ChapterRepo, ch *models.Chapter) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)
	if err := chRepo.Create(ctx, tx, ch); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

// splitChunks splits document content by markdown headers (# Chapter N).
//
// ponytail: simple regex split — no AST parser needed for markdown chapters.
func (s *IngestionService) splitChunks(content string) []ingestionChunk {
	if strings.TrimSpace(content) == "" {
		return nil
	}

	headerRe := regexp.MustCompile(`(?m)^# (.+)$`)
	locs := headerRe.FindAllStringSubmatchIndex(content, -1)

	if len(locs) == 0 {
		return splitByParagraphs(content)
	}

	chunks := make([]ingestionChunk, 0, len(locs))
	for i, loc := range locs {
		title := content[loc[2]:loc[3]]
		bodyStart := loc[1] + 1 // after the newline
		var bodyEnd int
		if i+1 < len(locs) {
			bodyEnd = locs[i+1][0]
		} else {
			bodyEnd = len(content)
		}
		body := strings.TrimSpace(content[bodyStart:bodyEnd])
		if body != "" {
			chunks = append(chunks, ingestionChunk{title: title, content: body})
		}
	}
	return chunks
}

// splitByParagraphs splits content at paragraph boundaries when there are no
// markdown headers. Each chunk targets ~50K chars so entity extraction gets
// manageable text and progress is granular.
//
// ponytail: greedy paragraph fill — splits at \n\n boundaries, no tokenizer.
func splitByParagraphs(content string) []ingestionChunk {
	const maxChunkSize = 50_000
	content = strings.TrimSpace(content)
	if content == "" {
		return nil
	}
	if len(content) <= maxChunkSize {
		return []ingestionChunk{{title: "Untitled", content: content}}
	}

	paragraphs := strings.Split(content, "\n\n")
	chunks := make([]ingestionChunk, 0, len(paragraphs)/3+1)
	var buf strings.Builder
	part := 1
	for _, p := range paragraphs {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		if buf.Len() > 0 && buf.Len()+len(p) > maxChunkSize {
			chunks = append(chunks, ingestionChunk{
				title:   fmt.Sprintf("Part %d", part),
				content: buf.String(),
			})
			buf.Reset()
			part++
		}
		if buf.Len() > 0 {
			buf.WriteString("\n\n")
		}
		buf.WriteString(p)
	}
	if buf.Len() > 0 {
		chunks = append(chunks, ingestionChunk{
			title:   fmt.Sprintf("Part %d", part),
			content: buf.String(),
		})
	}
	return chunks
}

// resolveAndBuildGraph resolves or creates entities and builds graph nodes,
// returning the number of entities successfully resolved.
//
// ponytail: reuses EntityService.ResolveOrCreate — same dedup/merge logic.
func (s *IngestionService) resolveAndBuildGraph(ctx context.Context, universeID uuid.UUID, extracted *ExtractedEntities) int {
	if extracted == nil {
		return 0
	}

	allEntities := make([]repositories.ExtractedEntity, 0)
	for _, e := range extracted.Characters {
		allEntities = append(allEntities, repositories.ExtractedEntity{
			Type: e.Type, Name: e.Name, Aliases: e.Aliases,
			Description: e.Description, Status: e.Status, Properties: e.Properties,
		})
	}
	for _, e := range extracted.Places {
		allEntities = append(allEntities, repositories.ExtractedEntity{
			Type: e.Type, Name: e.Name, Aliases: e.Aliases,
			Description: e.Description, Status: e.Status, Properties: e.Properties,
		})
	}
	for _, e := range extracted.Events {
		allEntities = append(allEntities, repositories.ExtractedEntity{
			Type: e.Type, Name: e.Name, Aliases: e.Aliases,
			Description: e.Description, Status: e.Status, Properties: e.Properties,
		})
	}
	for _, e := range extracted.Factions {
		allEntities = append(allEntities, repositories.ExtractedEntity{
			Type: e.Type, Name: e.Name, Aliases: e.Aliases,
			Description: e.Description, Status: e.Status, Properties: e.Properties,
		})
	}
	for _, e := range extracted.WorldRules {
		allEntities = append(allEntities, repositories.ExtractedEntity{
			Type: e.Type, Name: e.Name, Aliases: e.Aliases,
			Description: e.Description, Status: e.Status, Properties: e.Properties,
		})
	}
	for _, e := range extracted.PlotDevelopments {
		allEntities = append(allEntities, repositories.ExtractedEntity{
			Type: e.Type, Name: e.Name, Aliases: e.Aliases,
			Description: e.Description, Status: e.Status, Properties: e.Properties,
		})
	}

	resolved := 0
	for _, ee := range allEntities {
		if _, _, _, err := s.entitySvc.ResolveOrCreate(ctx, universeID, ee); err != nil {
			log.Printf("[ingestion] resolve entity %s: %v", ee.Name, err)
			continue
		}
		resolved++
	}
	return resolved
}

// emitProgress sends an ingestion_progress WebSocket event to the resolved
// universe owner.
func (s *IngestionService) emitProgress(jobID, userID uuid.UUID, status string, processed, total int) {
	if s.hub == nil {
		return
	}
	payload, _ := json.Marshal(map[string]any{
		"job_id":             jobID.String(),
		"status":             status,
		"chapters_processed": processed,
		"total_chapters":     total,
	})
	msg := models.WSMessage{
		Type:    "ingestion_progress",
		Payload: payload,
	}
	// ponytail: userID is the universe owner resolved once in runWorker.
	// Delivery remains best-effort — SendToUser's error is discarded because
	// an offline/missing WS connection is expected and non-fatal.
	_ = s.hub.SendToUser(userID, msg)
}

// updateProgress persists the chapter/entity counters, mirroring what
// emitProgress reports over WS. Best-effort like updateJobStatus.
func (s *IngestionService) updateProgress(ctx context.Context, jobID uuid.UUID, totalDetected, processed, entities int) {
	if s.pool == nil {
		return
	}
	repo := repositories.NewIngestionRepo(s.pool)
	if err := repo.UpdateProgress(ctx, jobID, totalDetected, processed, entities); err != nil {
		log.Printf("[ingestion] update progress job %s: %v", jobID, err)
	}
}

// updateJobStatus persists a status change to the ingestion_jobs table.
func (s *IngestionService) updateJobStatus(ctx context.Context, jobID uuid.UUID, status, errMsg string) {
	if s.pool == nil {
		return
	}
	repo := repositories.NewIngestionRepo(s.pool)
	if err := repo.UpdateStatus(ctx, jobID, status, errMsg); err != nil {
		log.Printf("[ingestion] update status %s: %v", status, err)
	}
}
