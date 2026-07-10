package services

import "testing"

// TestSafeAnalyzeRecoversPanic proves the recover seam used by
// runPostIngestAnalysis contains a panic instead of letting it escape and
// crash the process. Post-ingest analysis parses untrusted LLM tool-call JSON
// and type-asserts on it; an unrecovered panic on any goroutine would take
// down every user's WebSocket, not just the ingestion job. If safeAnalyze
// failed to recover, this panic would propagate and fail (crash) the test.
func TestSafeAnalyzeRecoversPanic(t *testing.T) {
	safeAnalyze("test unit", func() {
		panic("boom from untrusted analysis")
	})
	// Reaching this line means the panic was recovered, not propagated.
}
