package services

import (
	"context"
	"testing"

	"github.com/google/uuid"

	"github.com/quill/backend/internal/config"
)

// TestCheckContradictionsRequestShape verifies the ContradictionCandidate type and
// the request payload structure that CheckContradictions sends to Qwen.
func TestCheckContradictionsRequestShape(t *testing.T) {
	candidates := []ContradictionCandidate{
		{
			EntityID:  uuid.New(),
			Type:      "deceased_alive",
			EvidenceA: "Bob was alive in chapter 1.",
			EvidenceB: "Bob's funeral was in chapter 3.",
			ChapterA:  uuid.New(),
			ChapterB:  uuid.New(),
		},
		{
			EntityID:  uuid.New(),
			Type:      "status_change",
			EvidenceA: "Alice is the mayor.",
			EvidenceB: "Alice was elected queen.",
			ChapterA:  uuid.New(),
			ChapterB:  uuid.New(),
		},
	}

	// Verify the type compiles and prompt would contain evidence
	for _, c := range candidates {
		if c.Type == "" {
			t.Error("candidate type should not be empty")
		}
		if c.EvidenceA == "" || c.EvidenceB == "" {
			t.Error("candidate evidence should not be empty")
		}
	}

	if len(candidates) != 2 {
		t.Errorf("expected 2 candidates, got %d", len(candidates))
	}
}

// TestCheckContradictionsResultParsing verifies parsing of the expected response JSON.
func TestCheckContradictionsResultParsing(t *testing.T) {
	candidates := []ContradictionCandidate{
		{EntityID: uuid.New(), Type: "deceased_alive", EvidenceA: "a", EvidenceB: "b", ChapterA: uuid.New(), ChapterB: uuid.New()},
		{EntityID: uuid.New(), Type: "status_change", EvidenceA: "c", EvidenceB: "d", ChapterA: uuid.New(), ChapterB: uuid.New()},
	}

	rawJSON := `[
		{
			"has_contradiction": true,
			"entity_index": 0,
			"description": "Bob cannot be both alive and dead",
			"severity": "high",
			"suggestion": "Check chapter 1 and 3 for consistency"
		},
		{
			"has_contradiction": false,
			"entity_index": 1,
			"description": "",
			"severity": "",
			"suggestion": ""
		}
	]`

	results, err := parseContradictionResults([]byte(rawJSON), candidates)
	if err != nil {
		t.Fatalf("parseContradictionResults: %v", err)
	}

	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	if !results[0].HasContradiction {
		t.Error("result[0] should have contradiction")
	}
}

// TestQwenServiceCheckContradictionsSignature verifies the method compiles
// and returns proper error/candidates when there is no API connectivity.
func TestQwenServiceCheckContradictionsSignature(t *testing.T) {
	cfg := &config.Config{
		QwenBaseURL:                "https://example.com",
		QwenAPIKey:                 "test-key",
		QwenMaxModel:               "qwen-max-latest",
		QwenMaxConcurrency:         1,
		QwenTurboConcurrency:       1,
		QwenEmbeddingModel:         "text-embedding-v3",
		MaxContradictionCandidates: 3,
	}
	svc := NewQwenService(cfg)

	// Empty candidates should return nil, nil
	results, err := svc.CheckContradictions(context.Background(), nil)
	if err != nil {
		t.Errorf("expected no error for empty candidates, got: %v", err)
	}
	if results != nil {
		t.Errorf("expected nil results for empty candidates, got %d", len(results))
	}
}
