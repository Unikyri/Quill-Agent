package services

import (
	"sort"

	"github.com/google/uuid"
)

// namedPipeline pairs a pipeline name with its ranked entries. fuseRRFExplain
// takes an ORDERED SLICE of these (not a map) because Go map iteration order
// is randomized — fuseRRF fuses positionally, so a map input would make
// tie-break ordering non-deterministic across calls. Callers must supply
// pipelines in the SAME order fuseRRF receives its positional arguments.
type namedPipeline struct {
	Name    string
	Entries []rankedEntry
}

// RRFContribution records one pipeline appearance's contribution to an
// ExplainedItem's RRFScore.
type RRFContribution struct {
	Pipeline string  `json:"pipeline"`
	Rank     int     `json:"rank"`  // 1-indexed within its pipeline
	Delta    float64 `json:"delta"` // 1/(rrfK+rank)
}

// ExplainedItem is a fused, ranked recall result carrying a full per-pipeline
// contribution ledger, for the /recall/explain endpoint.
type ExplainedItem struct {
	ID                 string            `json:"id"`
	EntityID           uuid.UUID         `json:"entity_id"`
	Fact               string            `json:"fact"`
	RRFScore           float64           `json:"rrf_score"`
	Contributions      []RRFContribution `json:"contributions"`
	FitInBudget        bool              `json:"fit_in_budget"`
	PreRerankPosition  int               `json:"pre_rerank_position,omitempty"`
	PostRerankPosition int               `json:"post_rerank_position,omitempty"`
	RerankDelta        int               `json:"rerank_delta,omitempty"`
	RerankScore        float64           `json:"rerank_score,omitempty"`
}

// fuseRRFExplain mirrors fuseRRF's dedup/score/sort logic exactly (same
// dedup-by-id, same first-non-empty EntityID/Fact wins, same
// RRFScore-desc/ID-asc comparator) but additionally records one
// RRFContribution per pipeline appearance, so RRFScore == sum(Contributions
// Deltas) for every item and the final score/ordering matches fuseRRF
// byte-for-byte given the same pipelines in the same order.
func fuseRRFExplain(pipelines []namedPipeline) []ExplainedItem {
	byID := make(map[string]*ExplainedItem)
	order := make([]string, 0)

	for _, pipeline := range pipelines {
		for i, e := range pipeline.Entries {
			item, exists := byID[e.id]
			if !exists {
				item = &ExplainedItem{ID: e.id}
				byID[e.id] = item
				order = append(order, e.id)
			}
			delta := 1.0 / float64(rrfK+i+1)
			item.RRFScore += delta
			item.Contributions = append(item.Contributions, RRFContribution{
				Pipeline: pipeline.Name,
				Rank:     i + 1,
				Delta:    delta,
			})
			if item.EntityID == uuid.Nil && e.entityID != uuid.Nil {
				item.EntityID = e.entityID
			}
			if item.Fact == "" && e.fact != "" {
				item.Fact = e.fact
			}
		}
	}

	result := make([]ExplainedItem, 0, len(order))
	for _, id := range order {
		result = append(result, *byID[id])
	}
	sort.Slice(result, func(i, j int) bool {
		if result[i].RRFScore == result[j].RRFScore {
			return result[i].ID < result[j].ID
		}
		return result[i].RRFScore > result[j].RRFScore
	})
	return result
}
