package eval

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/quill/backend/internal/config"
	"github.com/quill/backend/internal/models"
	"github.com/quill/backend/internal/services"
)

type nativeRerankMetricRow struct {
	Query                           string
	BeforeRecallAt5, AfterRecallAt5 float64
	BeforeMRR, AfterMRR             float64
	BeforeNDCGAt5, AfterNDCGAt5     float64
}

// TestMemoryEvalDashScopeRerank is intentionally opt-in. It compares the
// existing RRF ordering with the native DashScope reranker on the same gold
// queries and writes a reproducible before/after section to results.md.
//
// The default test suite remains provider-neutral: without an explicit
// LLM_PROTOCOL=dashscope, a Qwen key, a migrated AGE database, and a reachable
// native endpoint this test is skipped rather than turning CI red.
func TestMemoryEvalDashScopeRerank(t *testing.T) {
	if !strings.EqualFold(os.Getenv("LLM_PROTOCOL"), "dashscope") {
		t.Skip("LLM_PROTOCOL=dashscope required for native rerank evaluation")
	}
	if os.Getenv("QWEN_API_KEY") == "" {
		t.Skip("QWEN_API_KEY required for native rerank evaluation")
	}

	cfg, err := config.Load()
	if err != nil {
		t.Skipf("DashScope config unavailable: %v", err)
	}
	native := services.NewDashScopeService(cfg, nil)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if _, err := native.Rerank(ctx, "native rerank health check", []string{"health check"}, 1); err != nil {
		t.Skipf("DashScope rerank API unavailable: %v", err)
	}
	compatible := services.NewQwenService(cfg, nil)
	if err := compatible.HealthCheck(ctx); err != nil {
		t.Skipf("Qwen-compatible embedding API unavailable: %v", err)
	}

	fx := setupSagaEval(t)
	ctx = context.Background()
	rows := make([]nativeRerankMetricRow, 0, len(fx.gold.Queries))

	// Baseline and reranked calls run against the same fixture and embeddings;
	// only the optional MemoryService reranker changes between passes.
	fx.svc.SetReranker(nil)
	for _, query := range fx.gold.Queries {
		embedding, err := fx.qwen.GenerateEmbedding(ctx, query.Query)
		if err != nil {
			t.Skipf("embedding API unavailable during rerank evaluation: %v", err)
		}
		items, err := fx.svc.RecallWithQuery(ctx, fx.universeID, embedding, query.Query, 5)
		if err != nil {
			t.Fatalf("baseline recall %s: %v", query.ID, err)
		}
		before := rerankMetricRow(items, query, fx.paragraphToEntities)

		// Keep the reranker on the same service so graph/vector/recency inputs
		// are identical and the comparison isolates post-RRF native reranking.
		fx.svc.SetReranker(native)
		rankedItems, err := fx.svc.RecallWithQuery(ctx, fx.universeID, embedding, query.Query, 5)
		if err != nil {
			t.Fatalf("reranked recall %s: %v", query.ID, err)
		}
		after := rerankMetricRow(rankedItems, query, fx.paragraphToEntities)
		rows = append(rows, nativeRerankMetricRow{
			Query:           query.Query,
			BeforeRecallAt5: before.RecallAt5, AfterRecallAt5: after.RecallAt5,
			BeforeMRR: before.MRR, AfterMRR: after.MRR,
			BeforeNDCGAt5: before.NDCGAt5, AfterNDCGAt5: after.NDCGAt5,
		})
		// Reset before the next baseline pass; this keeps every row paired.
		fx.svc.SetReranker(nil)
	}

	writeNativeRerankReport(t, "../../Docs/eval/results.md", rows)
}

func rerankMetricRow(items []models.RecallItem, query GoldQuery, paragraphToEntities map[string][]uuid.UUID) QueryReport {
	retrieved := retrievedEntityIDs(items, paragraphToEntities)
	relevant := relevantSet(query)
	return QueryReport{
		RecallAt5: recallAtK(retrieved, relevant, 5),
		MRR:       mrr(retrieved, relevant),
		NDCGAt5:   ndcgAtK(retrieved, relevant, 5),
	}
}

func writeNativeRerankReport(t *testing.T, path string, rows []nativeRerankMetricRow) {
	t.Helper()
	var beforeRecall, afterRecall, beforeMRR, afterMRR, beforeNDCG, afterNDCG float64
	for _, row := range rows {
		beforeRecall += row.BeforeRecallAt5
		afterRecall += row.AfterRecallAt5
		beforeMRR += row.BeforeMRR
		afterMRR += row.AfterMRR
		beforeNDCG += row.BeforeNDCGAt5
		afterNDCG += row.AfterNDCGAt5
	}
	if len(rows) > 0 {
		count := float64(len(rows))
		beforeRecall, afterRecall = beforeRecall/count, afterRecall/count
		beforeMRR, afterMRR = beforeMRR/count, afterMRR/count
		beforeNDCG, afterNDCG = beforeNDCG/count, afterNDCG/count
	}

	var body strings.Builder
	body.WriteString("Generated from TestMemoryEvalDashScopeRerank with LLM_PROTOCOL=dashscope.\n\n")
	body.WriteString("| Mode | recall@5 | MRR | nDCG@5 |\n")
	body.WriteString("|------|----------|-----|--------|\n")
	fmt.Fprintf(&body, "| RRF baseline | %.3f | %.3f | %.3f |\n", beforeRecall, beforeMRR, beforeNDCG)
	fmt.Fprintf(&body, "| Native rerank | %.3f | %.3f | %.3f |\n", afterRecall, afterMRR, afterNDCG)
	fmt.Fprintf(&body, "| Delta | %+.3f | %+.3f | %+.3f |\n", afterRecall-beforeRecall, afterMRR-beforeMRR, afterNDCG-beforeNDCG)
	body.WriteString("\n| Query | RRF recall@5 | Rerank recall@5 | RRF MRR | Rerank MRR | RRF nDCG@5 | Rerank nDCG@5 |\n")
	body.WriteString("|-------|--------------|-----------------|---------|------------|-------------|----------------|\n")
	for _, row := range rows {
		fmt.Fprintf(&body, "| %s | %.3f | %.3f | %.3f | %.3f | %.3f | %.3f |\n", row.Query, row.BeforeRecallAt5, row.AfterRecallAt5, row.BeforeMRR, row.AfterMRR, row.BeforeNDCGAt5, row.AfterNDCGAt5)
	}

	existing := ""
	if data, err := os.ReadFile(path); err == nil {
		existing = string(data)
	}
	writeReportFile(t, path, replaceSection(existing, "## Native DashScope Rerank Evaluation", body.String()))
}
