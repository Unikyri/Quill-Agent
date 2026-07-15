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
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"
	"unicode"

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

const (
	ingestionRelationshipTimeout     = 15 * time.Second
	ingestionRelationshipCorpusLimit = 12_000
	ingestionRelationshipEntityLimit = 40
	ingestionFallbackEdgeLimit       = 100
)

// IngestionQwen is the minimal Qwen interface used by IngestionService.
// QwenService satisfies this interface.
type IngestionQwen interface {
	ExtractEntities(ctx context.Context, text, universeContext string) (*ExtractedEntities, error)
	AnalyzeRelationships(ctx context.Context, text string, entityNames []string) ([]map[string]interface{}, error)
	GenerateEmbedding(ctx context.Context, text string) ([]float32, error)
	GenerateEmbeddingBatch(ctx context.Context, texts []string) ([][]float32, error)
}

// relationshipEdgeWriter keeps relationship persistence testable without
// bypassing GraphRepo in production. GraphRepo remains the default writer and
// therefore retains AGE search-path and Cypher identifier protections.
type relationshipEdgeWriter interface {
	CreateEdge(ctx context.Context, graphName, sourceEntityID, targetEntityID, relType string, properties map[string]interface{}) error
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
	// relationshipEdges is a test seam; nil uses graphRepo.
	relationshipEdges relationshipEdgeWriter
	qwenSvc           IngestionQwen
	hub               AnalysisHub

	// Post-ingest bounded analysis (D4) — all nil-safe, wired via
	// SetPostIngestAnalysis. Unset means analysis is silently skipped.
	contraSvc           *ContradictionService
	plotHoleSvc         *PlotHoleService
	analysisBudgetMgr   *ContextBudgetManager
	analysisMaxChapters int
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

// supportedFileTypes are the extensions parseDocument can handle. Checked
// synchronously in Start (before any I/O) so unsupported uploads (legacy
// .doc, unknown formats) get an immediate 400 instead of a garbage job row.
var supportedFileTypes = map[string]bool{"md": true, "txt": true, "docx": true, "pdf": true}

// ErrUnsupportedFileType is returned by Start when filename's extension isn't
// one of supportedFileTypes.
var ErrUnsupportedFileType = errors.New("unsupported file type — only .md, .txt, .docx, and .pdf are supported (legacy .doc? Save as .docx)")

// Start creates an ingestion job and kicks off the async pipeline.
// Returns the job ID immediately; duplicate is true when the same content
// was already ingested into this universe (the existing job's ID is
// returned and no worker is started). The caller should return 202 Accepted
// for new jobs and 200 for duplicates.
func (s *IngestionService) Start(ctx context.Context, universeID uuid.UUID, reader io.Reader, filename string) (uuid.UUID, bool, error) {
	fileType := fileTypeOf(filename)
	if !supportedFileTypes[fileType] {
		return uuid.Nil, false, ErrUnsupportedFileType
	}

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
			// Work title = filename stem. The work row is created here in
			// Start, before the document is parsed in runWorker, so a
			// heading-derived title (the proposal's original idea) isn't
			// available yet — and the first heading is usually just
			// "Chapter 1" anyway, not a useful book title.
			title := strings.TrimSuffix(filename, filepath.Ext(filename))
			work := models.Work{
				ID:         uuid.New(),
				UniverseID: universeID,
				Title:      title,
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

		if err := repo.Create(ctx, jobID, universeID, workID, "pending", filename, fileType, hash); err != nil {
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

// ingestedChapter tracks a persisted chapter and the entities resolved from
// it, collected during runWorker's chunk loop for the post-ingest analysis
// pass (SetPostIngestAnalysis).
type ingestedChapter struct {
	ID       uuid.UUID
	Content  string
	Entities []ResolvedEntity
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

	// Parse the raw upload into plain text. Raw binary must never reach
	// splitChunks/chapters.content — a parse failure or empty/whitespace-only
	// extraction fails the job cleanly instead.
	text, err := parseDocument(filename, content)
	if err != nil || strings.TrimSpace(text) == "" {
		msg := "document contains no text"
		if err != nil {
			msg = err.Error()
		}
		s.updateJobStatus(ctx, jobID, "failed", msg)
		s.emitProgress(jobID, ownerID, "failed", 0, 0)
		// The failed job row (with its error_message) is the durable record of
		// this attempt — it must survive so a reload shows "upload failed: …"
		// instead of nothing. We deliberately do NOT delete the Work here:
		// ingestion_jobs.work_id is `NOT NULL REFERENCES works(id) ON DELETE
		// CASCADE` (migration 012), so deleting the Work would cascade-delete
		// this job row and its error_message. The Work has a meaningful title
		// (the filename stem) and is user-removable via the delete-work button.
		return
	}

	// Split parsed text into chapters (markdown/EN/ES/roman/ALL-CAPS heading
	// cascade, falling back to paragraph-boundary chunks).
	chunks := s.splitChunks(text)
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
	var ingestedChapters []ingestedChapter
	var relationshipCorpus strings.Builder
	var relationshipEntities []ResolvedEntity
	var relationshipChunks [][]ResolvedEntity

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
			mentionText := ch.content
			if len(mentionText) > 120 {
				mentionText = mentionText[:120]
			}
			count, resolved := s.resolveAndBuildGraph(ctx, universeID, extracted, mentionText)
			entitiesTotal += count
			relationshipCorpus.WriteString(truncateIngestionRelationshipCorpus(ch.content, relationshipCorpus.Len()))
			relationshipEntities = append(relationshipEntities, resolved...)
			if len(resolved) > 0 {
				relationshipChunks = append(relationshipChunks, resolved)
			}
			if chapterID != uuid.Nil {
				ingestedChapters = append(ingestedChapters, ingestedChapter{
					ID:       chapterID,
					Content:  ch.content,
					Entities: resolved,
				})
			}
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

	// Relationship extraction is deliberately one bounded, best-effort pass for
	// the whole manuscript. Calling the model per chunk made ingestion serially
	// slow and exceeded the E2E budget on long documents.
	s.enrichRelationships(ctx, universeID, relationshipCorpus.String(), relationshipEntities, relationshipChunks)

	// Bounded post-ingest analysis (contradiction + plot-hole checks) runs
	// before the job is marked completed, so the job honestly reports
	// "running" until analysis ends. Best-effort/enrichment: never flips a
	// completed job to failed. No-ops when SetPostIngestAnalysis wasn't
	// called (nil deps).
	s.runPostIngestAnalysis(ctx, universeID, ingestedChapters, ownerID)

	s.updateJobStatus(ctx, jobID, "completed", "")
	s.updateProgress(ctx, jobID, len(chunks), len(chunks), entitiesTotal)
	s.emitProgress(jobID, ownerID, "completed", len(chunks), len(chunks))
}

func (s *IngestionService) enrichRelationships(ctx context.Context, universeID uuid.UUID, corpus string, entities []ResolvedEntity, chunks [][]ResolvedEntity) {
	relationshipCtx, cancelRelationships := context.WithTimeout(ctx, ingestionRelationshipTimeout)
	persisted, err := s.persistRelationships(relationshipCtx, universeID, corpus, entities)
	cancelRelationships()
	if err != nil {
		log.Printf("[ingestion] analyze relationships: %v", err)
	}
	if persisted > 0 {
		return
	}

	// Model enrichment is optional. When it produces no usable edge (including
	// timeout/error), preserve basic graph connectivity from deterministic
	// canonical co-occurrence within each imported chunk.
	if fallbackEdges, fallbackErr := s.persistCooccurrenceEdges(ctx, universeID, chunks); fallbackErr != nil {
		log.Printf("[ingestion] create co-occurrence fallback edges: %v", fallbackErr)
	} else if fallbackEdges > 0 {
		log.Printf("[ingestion] created %d CO_OCCURS_WITH fallback edges", fallbackEdges)
	}
}

func truncateIngestionRelationshipCorpus(chunk string, used int) string {
	if used >= ingestionRelationshipCorpusLimit {
		return ""
	}
	remaining := ingestionRelationshipCorpusLimit - used
	prefix := ""
	if used > 0 && chunk != "" {
		prefix = "\n\n"
		remaining -= len(prefix)
		if remaining <= 0 {
			return ""
		}
	}
	if len(chunk) > remaining {
		chunk = chunk[:remaining]
	}
	return prefix + chunk
}

// persistRelationships enriches the universe graph from relationships found
// in one ingested chunk. GraphRepo owns identifier validation and AGE session
// safety; an invalid model-supplied relationship is therefore logged and
// skipped without failing the ingestion job.
func (s *IngestionService) persistRelationships(ctx context.Context, universeID uuid.UUID, text string, resolved []ResolvedEntity) (int, error) {
	edgeWriter := s.relationshipEdges
	if edgeWriter == nil {
		edgeWriter = s.graphRepo
	}
	if s.qwenSvc == nil || edgeWriter == nil || len(resolved) == 0 {
		return 0, nil
	}

	entities := make([]models.Entity, 0, min(len(resolved), ingestionRelationshipEntityLimit))
	entityNames := make([]string, 0, min(len(resolved), ingestionRelationshipEntityLimit))
	seenEntityIDs := make(map[uuid.UUID]struct{}, len(resolved))
	for _, item := range resolved {
		if item.Entity.Name == "" || len(entityNames) >= ingestionRelationshipEntityLimit {
			continue
		}
		if _, seen := seenEntityIDs[item.Entity.ID]; seen {
			continue
		}
		seenEntityIDs[item.Entity.ID] = struct{}{}
		entities = append(entities, item.Entity)
		entityNames = append(entityNames, item.Entity.Name)
	}
	if len(entityNames) == 0 {
		return 0, nil
	}

	relationships, err := s.qwenSvc.AnalyzeRelationships(ctx, text, entityNames)
	if err != nil {
		return 0, err
	}

	graphName := "universe_" + universeID.String()
	persisted := 0
	for _, relationship := range relationships {
		sourceName, _ := relationship["source"].(string)
		targetName, _ := relationship["target"].(string)
		relationType, _ := relationship["type"].(string)
		if sourceName == "" || targetName == "" || relationType == "" {
			log.Printf("[ingestion] skip malformed relationship: source=%q target=%q type=%q", sourceName, targetName, relationType)
			continue
		}
		source, sourceFound, sourceDiagnostic := resolveIngestionRelationshipEntity(sourceName, entities)
		if !sourceFound {
			log.Printf("[ingestion] skip relationship source %q: %s", sourceName, sourceDiagnostic)
			continue
		}
		target, targetFound, targetDiagnostic := resolveIngestionRelationshipEntity(targetName, entities)
		if !targetFound {
			log.Printf("[ingestion] skip relationship target %q: %s", targetName, targetDiagnostic)
			continue
		}
		if err := edgeWriter.CreateEdge(ctx, graphName, source.ID.String(), target.ID.String(), relationType, nil); err != nil {
			log.Printf("[ingestion] create graph edge %q->%q (%q): %v", sourceName, targetName, relationType, err)
			continue
		}
		persisted++
	}
	return persisted, nil
}

type ingestionEntityPair struct {
	source models.Entity
	target models.Entity
}

func (s *IngestionService) persistCooccurrenceEdges(ctx context.Context, universeID uuid.UUID, chunks [][]ResolvedEntity) (int, error) {
	edgeWriter := s.relationshipEdges
	if edgeWriter == nil {
		edgeWriter = s.graphRepo
	}
	if edgeWriter == nil {
		return 0, nil
	}

	pairs := make(map[string]ingestionEntityPair)
	for _, chunk := range chunks {
		if len(pairs) >= ingestionFallbackEdgeLimit {
			break
		}
		unique := make(map[uuid.UUID]models.Entity, len(chunk))
		for _, item := range chunk {
			if item.Entity.ID != uuid.Nil {
				unique[item.Entity.ID] = item.Entity
			}
		}
		entities := make([]models.Entity, 0, len(unique))
		for _, entity := range unique {
			entities = append(entities, entity)
		}
		sort.Slice(entities, func(i, j int) bool { return entities[i].ID.String() < entities[j].ID.String() })
		for i := 0; i < len(entities); i++ {
			for j := i + 1; j < len(entities); j++ {
				if len(pairs) >= ingestionFallbackEdgeLimit {
					break
				}
				key := entities[i].ID.String() + ":" + entities[j].ID.String()
				if _, exists := pairs[key]; exists {
					continue
				}
				pairs[key] = ingestionEntityPair{source: entities[i], target: entities[j]}
			}
			if len(pairs) >= ingestionFallbackEdgeLimit {
				break
			}
		}
	}

	keys := make([]string, 0, len(pairs))
	for key := range pairs {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	graphName := "universe_" + universeID.String()
	persisted := 0
	for _, key := range keys {
		pair := pairs[key]
		if err := edgeWriter.CreateEdge(ctx, graphName, pair.source.ID.String(), pair.target.ID.String(), "CO_OCCURS_WITH", nil); err != nil {
			log.Printf("[ingestion] create CO_OCCURS_WITH edge %q->%q: %v", pair.source.Name, pair.target.Name, err)
			continue
		}
		persisted++
	}
	return persisted, nil
}

// resolveIngestionRelationshipEntity maps a model-emitted name to exactly one
// persisted entity. Exact canonical/alias matches win; otherwise a shortened
// name may match a canonical name or alias only when its word prefix is unique.
// Ambiguity is an intentional no-edge outcome rather than a guess.
func resolveIngestionRelationshipEntity(name string, entities []models.Entity) (models.Entity, bool, string) {
	query := strings.TrimSpace(name)
	if query == "" {
		return models.Entity{}, false, "empty entity name"
	}

	exact := matchingRelationshipEntities(query, entities, false)
	if len(exact) == 1 {
		return exact[0], true, ""
	}
	if len(exact) > 1 {
		return models.Entity{}, false, fmt.Sprintf("ambiguous canonical or alias match (%d entities)", len(exact))
	}

	prefix := matchingRelationshipEntities(query, entities, true)
	if len(prefix) == 1 {
		return prefix[0], true, ""
	}
	if len(prefix) > 1 {
		return models.Entity{}, false, fmt.Sprintf("ambiguous prefix match (%d entities)", len(prefix))
	}
	return models.Entity{}, false, "no canonical, alias, or unique prefix match"
}

func matchingRelationshipEntities(query string, entities []models.Entity, allowPrefix bool) []models.Entity {
	matches := make([]models.Entity, 0, 1)
	for _, entity := range entities {
		candidateNames := append([]string{entity.Name}, entity.Aliases...)
		for _, candidate := range candidateNames {
			candidate = strings.TrimSpace(candidate)
			exact := strings.EqualFold(query, candidate)
			prefix := allowPrefix && hasRelationshipWordPrefix(query, candidate)
			if exact || prefix {
				matches = append(matches, entity)
				break
			}
		}
	}
	return matches
}

func hasRelationshipWordPrefix(query, candidate string) bool {
	queryWords := strings.Fields(query)
	candidateWords := strings.Fields(candidate)
	if len(queryWords) == 0 || len(queryWords) >= len(candidateWords) {
		return false
	}
	return strings.EqualFold(strings.Join(queryWords, " "), strings.Join(candidateWords[:len(queryWords)], " "))
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

// maxSaneHeadingMatches guards against a heading pattern matching almost
// every line (a false positive, e.g. a doc that happens to contain many bare
// roman-numeral-looking lines) — treated the same as "no match", falling
// through to the next pattern in the cascade.
const maxSaneHeadingMatches = 500

// headingMatch is a single detected chapter heading: title is the extracted
// heading text, start/end are the byte offsets of the whole matched heading
// line in the source content (used to slice out chapter bodies).
type headingMatch struct {
	start, end int
	title      string
}

// headingPatterns is the priority cascade of heading patterns tried in
// splitChunks, in order — the first pattern class with >= 2 matches (and
// <= maxSaneHeadingMatches) wins. Each has exactly one capture group holding
// the extracted title text.
var headingPatterns = []*regexp.Regexp{
	regexp.MustCompile(`(?m)^#{1,3} (.+)$`), // markdown (also styled DOCX via D1 bonus)
	// ponytail: spelled-out numbers (one, two, three…) added because books like The Expanse use "Chapter One" not "Chapter 1"
	regexp.MustCompile(`(?mi)^[ \t]*(chapter[ \t]+(?:\d+|[ivxlc]+|one|two|three|four|five|six|seven|eight|nine|ten|eleven|twelve|thirteen|fourteen|fifteen|sixteen|seventeen|eighteen|nineteen|twenty|thirty|forty|fifty|sixty|seventy|eighty|ninety|hundred|thousand)\b.*)$`),      // English
	regexp.MustCompile(`(?mi)^[ \t]*(cap[ií]tulo[ \t]+(?:\d+|[ivxlc]+|uno|dos|tres|cuatro|cinco|seis|siete|ocho|nueve|diez|once|doce|trece|catorce|quince|dieciséis|diecisiete|dieciocho|diecinueve|veinte|treinta|cuarenta|cincuenta|sesenta|setenta|ochenta|noventa|cien)\b.*)$`), // Spanish
	regexp.MustCompile(`(?m)^[ \t]*([IVXLC]{1,7}\.?)[ \t]*$`),                // bare roman numeral
	regexp.MustCompile(`(?m)^[ \t]*([A-Z][a-z]+(?:\s+[A-Z][a-z]+)*)[ \t]*$`), // title case heading ("Holden", "The Rocinante")
}

// regexHeadingMatches runs re against content and returns one headingMatch
// per match, using the first capture group as the title.
func regexHeadingMatches(content string, re *regexp.Regexp) []headingMatch {
	locs := re.FindAllStringSubmatchIndex(content, -1)
	matches := make([]headingMatch, 0, len(locs))
	for _, loc := range locs {
		matches = append(matches, headingMatch{
			start: loc[0],
			end:   loc[1],
			title: strings.TrimSpace(content[loc[2]:loc[3]]),
		})
	}
	return matches
}

// isAllCapsHeadingLine reports whether a trimmed line looks like an
// ALL-CAPS chapter heading: short (<= 60 chars), no lowercase letters, and
// at least 3 letters (so pure punctuation/numeral lines don't qualify).
//
// ponytail: 10-char minimum plus rejection of sentence-ending punctuation
// prevents PDF artifacts like "JIM" or "THE END." from becoming fake chapters.
func isAllCapsHeadingLine(line string) bool {
	if len(line) < 10 || len(line) > 60 {
		return false
	}
	last := line[len(line)-1]
	if last == '.' || last == '!' || last == '?' || last == '"' || last == '»' || last == ',' || last == ';' || last == ':' {
		return false
	}
	letters := 0
	for _, r := range line {
		if unicode.IsLower(r) {
			return false
		}
		if unicode.IsUpper(r) {
			letters++
		}
	}
	return letters >= 3
}

// allCapsHeadingMatches scans content line-by-line for ALL-CAPS heading
// candidates — this shape (short lines, no lowercase) isn't expressible as a
// single regex the way the other patterns are.
func allCapsHeadingMatches(content string) []headingMatch {
	var matches []headingMatch
	offset := 0
	for _, line := range strings.Split(content, "\n") {
		trimmed := strings.TrimSpace(line)
		if isAllCapsHeadingLine(trimmed) {
			leading := len(line) - len(strings.TrimLeft(line, " \t"))
			start := offset + leading
			matches = append(matches, headingMatch{start: start, end: start + len(trimmed), title: trimmed})
		}
		offset += len(line) + 1 // +1 for the '\n' consumed by Split
	}
	return matches
}

// splitByHeadings turns a set of detected heading matches into chapter
// chunks: the body of each chunk runs from just after its heading line to
// just before the next heading (or end of content).
func splitByHeadings(content string, matches []headingMatch) []ingestionChunk {
	chunks := make([]ingestionChunk, 0, len(matches))
	for i, m := range matches {
		bodyStart := m.end + 1 // after the newline
		var bodyEnd int
		if i+1 < len(matches) {
			bodyEnd = matches[i+1].start
		} else {
			bodyEnd = len(content)
		}
		if bodyStart > len(content) {
			bodyStart = len(content)
		}
		if bodyEnd > bodyStart {
			body := strings.TrimSpace(content[bodyStart:bodyEnd])
			if body != "" {
				chunks = append(chunks, ingestionChunk{title: m.title, content: body})
			}
		}
	}
	return chunks
}

// splitChunks splits document content into chapters. It tries a priority
// cascade of heading patterns (markdown, English "Chapter N", Spanish
// "Capítulo N", bare roman numerals, then short ALL-CAPS lines) — the first
// pattern class with >= 2 matches wins. No pattern matching falls back to
// splitByParagraphs.
//
// ponytail: simple regex cascade — no AST parser needed for chapter detection.
func (s *IngestionService) splitChunks(content string) []ingestionChunk {
	if strings.TrimSpace(content) == "" {
		return nil
	}

	for _, re := range headingPatterns {
		matches := regexHeadingMatches(content, re)
		if len(matches) >= 2 && len(matches) <= maxSaneHeadingMatches {
			return splitByHeadings(content, matches)
		}
	}

	if matches := allCapsHeadingMatches(content); len(matches) >= 2 && len(matches) <= maxSaneHeadingMatches {
		return splitByHeadings(content, matches)
	}

	return splitByParagraphs(content)
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
// returning the number of entities successfully resolved and the resolved
// entities themselves (needed by the post-ingest analysis pass, D4).
//
// ponytail: reuses EntityService.ResolveOrCreate — same dedup/merge logic.
func (s *IngestionService) resolveAndBuildGraph(ctx context.Context, universeID uuid.UUID, extracted *ExtractedEntities, mentionText string) (int, []ResolvedEntity) {
	if extracted == nil {
		return 0, nil
	}

	allEntities := make([]repositories.ExtractedEntity, 0, len(extracted.All()))
	for _, e := range extracted.All() {
		allEntities = append(allEntities, repositories.ExtractedEntity{
			Type: e.Type, Name: e.Name, Aliases: e.Aliases,
			Description: e.Description, Status: e.Status, Properties: e.Properties,
		})
	}

	var resolved []ResolvedEntity
	for _, ee := range allEntities {
		entity, previousStatus, isNew, err := s.entitySvc.ResolveOrCreate(ctx, universeID, ee)
		if err != nil {
			log.Printf("[ingestion] resolve entity %s: %v", ee.Name, err)
			continue
		}
		resolved = append(resolved, ResolvedEntity{
			Entity:         *entity,
			MentionText:    mentionText,
			IsNew:          isNew,
			PreviousStatus: previousStatus,
		})
	}
	return len(resolved), resolved
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
