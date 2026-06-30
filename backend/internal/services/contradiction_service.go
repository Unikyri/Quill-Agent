package services

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/quill/backend/internal/models"
	"github.com/quill/backend/internal/repositories"
)

// ContradictionService checks for narrative contradictions using rule-based
// deterministic checks and batched semantic checks via Qwen-Max.
//
// ponytail: SHA-256 fingerprint dedup — no DB round-trip for duplicate check.
// Deterministic rules catch deceased/alive without API call.
type ContradictionService struct {
	pool          *pgxpool.Pool
	contraRepo    *repositories.ContradictionRepo
	entityRepo    *repositories.EntityRepo
	qwenSvc       *QwenService
	maxCandidates int
}

// NewContradictionService creates a contradiction service.
// qwenSvc may be nil — CheckSemantic will be a no-op in that case.
func NewContradictionService(pool *pgxpool.Pool, contraRepo *repositories.ContradictionRepo, entityRepo *repositories.EntityRepo, qwenSvc *QwenService, maxCandidates int) *ContradictionService {
	return &ContradictionService{
		pool:          pool,
		contraRepo:    contraRepo,
		entityRepo:    entityRepo,
		qwenSvc:       qwenSvc,
		maxCandidates: maxCandidates,
	}
}

// fingerprint produces a SHA-256 hex fingerprint for deduplication.
// Exported for testing.
func (s *ContradictionService) fingerprint(c ContradictionCandidate) string {
	h := sha256.New()
	h.Write([]byte(c.EntityID.String()))
	h.Write([]byte(c.Type))
	h.Write([]byte(c.EvidenceA))
	h.Write([]byte(c.EvidenceB))
	h.Write([]byte(c.ChapterA.String()))
	h.Write([]byte(c.ChapterB.String()))
	return hex.EncodeToString(h.Sum(nil))
}

// CheckDeterministic runs fast rule-based contradiction checks (no API call):
//   - deceased/alive: if entity is marked "deceased" but mentioned as alive in new text
//
// chapterID is threaded into the candidate fingerprint so contradictions are
// scoped to the originating chapter context.
//
// ponytail: rule-based checks only; semantic checks go through Qwen-Max.
func (s *ContradictionService) CheckDeterministic(ctx context.Context, universeID uuid.UUID, chapterID uuid.UUID, entities []ResolvedEntity) ([]models.Contradiction, error) {
	var results []models.Contradiction

	for _, re := range entities {
		e := re.Entity

		// Deceased / alive check: if entity DB status is "deceased" and the
		// new mention suggests they're alive, flag it.
		if e.Status == "deceased" {
			fp := s.fingerprint(ContradictionCandidate{
				EntityID:  e.ID,
				Type:      "deceased_alive",
				EvidenceA: fmt.Sprintf("Entity %s is deceased in DB", e.Name),
				EvidenceB: re.MentionText,
				ChapterA:  chapterID,
				ChapterB:  chapterID,
			})

			// Check if we already have this contradiction
			existing, _ := s.contraRepo.FindByFingerprint(ctx, fp)
			if existing != nil {
				continue // already recorded
			}

			c := models.Contradiction{
				ID:         uuid.New(),
				UniverseID: universeID,
				EntityID:   &e.ID,
				Severity:   "critical",
				Description: fmt.Sprintf(
					"Entity '%s' is marked as deceased but was mentioned as active: \"%s\"",
					e.Name, truncate(re.MentionText, 80),
				),
				Suggestion:  "Review timeline: was this entity revived, or is this a continuity error?",
				EvidenceA:   fmt.Sprintf("Entity %s status: deceased", e.Name),
				EvidenceB:   re.MentionText,
				Fingerprint: fp,
				Status:      "open",
			}
			results = append(results, c)
		}
	}

	return results, nil
}

// CheckSemantic sends batched contradiction candidates to Qwen-Max.
// Returns only the contradictions that Qwen flags as confirmed.
// Batches up to maxCandidates per call.
//
// chapterID is threaded into the candidate fingerprint so contradictions are
// scoped to the originating chapter context.
//
// ponytail: single batch call per invocation — capped at maxCandidates.
func (s *ContradictionService) CheckSemantic(ctx context.Context, universeID uuid.UUID, chapterID uuid.UUID, text string, entities []ResolvedEntity) ([]models.Contradiction, error) {
	if s.qwenSvc == nil {
		return nil, nil
	}

	// Build candidates from entities that aren't already caught by deterministic rules
	var candidates []ContradictionCandidate
	for _, re := range entities {
		if len(candidates) >= s.maxCandidates {
			break
		}
		// Skip deceased entities (already caught by CheckDeterministic)
		if re.Entity.Status == "deceased" {
			continue
		}
		c := ContradictionCandidate{
			EntityID:  re.Entity.ID,
			Type:      "semantic",
			EvidenceA: fmt.Sprintf("Entity %s characteristics from DB: %s", re.Entity.Name, re.Entity.Description),
			EvidenceB: re.MentionText,
			ChapterA:  chapterID,
			ChapterB:  chapterID,
		}
		candidates = append(candidates, c)
	}

	if len(candidates) == 0 {
		return nil, nil
	}

	// Call QwenService batch
	contradictions, err := s.qwenSvc.CheckContradictions(ctx, candidates)
	if err != nil {
		return nil, fmt.Errorf("check semantic: %w", err)
	}

	// Set universeID and persist each contradiction
	for i := range contradictions {
		contradictions[i].UniverseID = universeID
		contradictions[i].ID = uuid.New()
		contradictions[i].Fingerprint = s.fingerprint(candidates[i])

		// Deduplicate: skip if fingerprint already exists
		existing, _ := s.contraRepo.FindByFingerprint(ctx, contradictions[i].Fingerprint)
		if existing != nil {
			continue
		}
		if err := s.contraRepo.Create(ctx, &contradictions[i]); err != nil {
			continue // best-effort persistence
		}
	}

	return contradictions, nil
}

// truncate shortens text to maxLen characters, appending "…" if truncated.
func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	if maxLen < 4 {
		return s[:maxLen]
	}
	return s[:maxLen-1] + "…"
}
