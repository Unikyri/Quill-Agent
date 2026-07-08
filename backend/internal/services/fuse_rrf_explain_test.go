package services

import "testing"

// TestFuseRRFExplainEmptyPipelines proves fuseRRFExplain matches fuseRRF on
// empty input (spec: "Empty pipelines").
func TestFuseRRFExplainEmptyPipelines(t *testing.T) {
	explained := fuseRRFExplain(nil)
	plain := fuseRRF()

	if len(explained) != 0 || len(plain) != 0 {
		t.Fatalf("expected both empty, got explained=%d plain=%d", len(explained), len(plain))
	}
}

// TestFuseRRFExplainSinglePipelineParity proves score+order parity against
// fuseRRF for a single-pipeline input (spec: "Single-pipeline input").
func TestFuseRRFExplainSinglePipelineParity(t *testing.T) {
	list := []rankedEntry{{id: "a", fact: "fact a", source: "vector"}}

	explained := fuseRRFExplain([]namedPipeline{{Name: "vector", Entries: list}})
	plain := fuseRRF(list)

	if len(explained) != 1 || len(plain) != 1 {
		t.Fatalf("expected 1 item each, got explained=%d plain=%d", len(explained), len(plain))
	}
	if explained[0].RRFScore != plain[0].RRFScore {
		t.Errorf("RRFScore mismatch: explained=%v plain=%v", explained[0].RRFScore, plain[0].RRFScore)
	}
	if explained[0].ID != plain[0].ID {
		t.Errorf("ID mismatch: explained=%v plain=%v", explained[0].ID, plain[0].ID)
	}
}

// TestFuseRRFExplainMultiPipelineOverlapParity proves score+order parity
// against fuseRRF when pipelines overlap at different ranks (spec:
// "Multi-pipeline overlap").
func TestFuseRRFExplainMultiPipelineOverlapParity(t *testing.T) {
	graphList := []rankedEntry{{id: "shared", fact: "shared fact", source: "graph"}}
	recencyList := []rankedEntry{{id: "shared", fact: "shared fact", source: "recency"}}
	vectorList := []rankedEntry{{id: "solo", fact: "solo fact", source: "vector"}}

	explained := fuseRRFExplain([]namedPipeline{
		{Name: "graph", Entries: graphList},
		{Name: "recency", Entries: recencyList},
		{Name: "vector", Entries: vectorList},
	})
	plain := fuseRRF(graphList, recencyList, vectorList)

	if len(explained) != len(plain) {
		t.Fatalf("length mismatch: explained=%d plain=%d", len(explained), len(plain))
	}
	for i := range plain {
		if explained[i].ID != plain[i].ID {
			t.Errorf("order mismatch at %d: explained=%v plain=%v", i, explained[i].ID, plain[i].ID)
		}
		if explained[i].RRFScore != plain[i].RRFScore {
			t.Errorf("RRFScore mismatch at %d: explained=%v plain=%v", i, explained[i].RRFScore, plain[i].RRFScore)
		}
	}
}

// TestFuseRRFExplainContributionLedgerSumsToScore proves each
// ExplainedItem's Contributions sum to its RRFScore, with correct
// Pipeline/Rank/Delta bookkeeping (spec: "Contribution ledger sums to
// score").
func TestFuseRRFExplainContributionLedgerSumsToScore(t *testing.T) {
	// "shared" appears rank-1 in vector, rank-2 in graph, rank-1 in keyword.
	vectorList := []rankedEntry{{id: "shared", fact: "f", source: "vector"}}
	graphList := []rankedEntry{{id: "other", fact: "o", source: "graph"}, {id: "shared", fact: "f", source: "graph"}}
	keywordList := []rankedEntry{{id: "shared", fact: "f", source: "keyword"}}

	explained := fuseRRFExplain([]namedPipeline{
		{Name: "vector", Entries: vectorList},
		{Name: "graph", Entries: graphList},
		{Name: "keyword", Entries: keywordList},
	})

	var shared *ExplainedItem
	for i := range explained {
		if explained[i].ID == "shared" {
			shared = &explained[i]
		}
	}
	if shared == nil {
		t.Fatal("expected 'shared' item in result")
	}
	if len(shared.Contributions) != 3 {
		t.Fatalf("expected 3 contributions, got %d: %+v", len(shared.Contributions), shared.Contributions)
	}

	var sum float64
	wantByPipeline := map[string]struct {
		rank  int
		delta float64
	}{
		"vector":  {1, 1.0 / 61.0},
		"graph":   {2, 1.0 / 62.0},
		"keyword": {1, 1.0 / 61.0},
	}
	for _, c := range shared.Contributions {
		want, ok := wantByPipeline[c.Pipeline]
		if !ok {
			t.Fatalf("unexpected pipeline %q in contributions", c.Pipeline)
		}
		if c.Rank != want.rank {
			t.Errorf("pipeline %s: Rank = %d, want %d", c.Pipeline, c.Rank, want.rank)
		}
		if c.Delta != want.delta {
			t.Errorf("pipeline %s: Delta = %v, want %v", c.Pipeline, c.Delta, want.delta)
		}
		sum += c.Delta
	}
	if sum != shared.RRFScore {
		t.Errorf("sum of Deltas = %v, want RRFScore %v", sum, shared.RRFScore)
	}
}

// TestFuseRRFExplainTieBreakByIDAscending proves tie-break by ID ascending
// is preserved (mirrors TestFuseRRFTieBreakByIDAscending).
func TestFuseRRFExplainTieBreakByIDAscending(t *testing.T) {
	listA := []rankedEntry{{id: "zeta", fact: "z", source: "vector"}}
	listB := []rankedEntry{{id: "alpha", fact: "a", source: "keyword"}}

	explained := fuseRRFExplain([]namedPipeline{
		{Name: "vector", Entries: listA},
		{Name: "keyword", Entries: listB},
	})

	if len(explained) != 2 {
		t.Fatalf("expected 2 items, got %d", len(explained))
	}
	if explained[0].RRFScore != explained[1].RRFScore {
		t.Fatalf("expected tied scores, got %v vs %v", explained[0].RRFScore, explained[1].RRFScore)
	}
	if explained[0].ID != "alpha" || explained[1].ID != "zeta" {
		t.Errorf("expected tie-break alpha before zeta, got order: %q, %q", explained[0].ID, explained[1].ID)
	}
}
