# Sprint 3 — Native DashScope Client & Qwen Capabilities

**Tier:** P2 · **SRS coverage:** DS-1..6, CC-1..4, SO-1..3, RR-1..4 (TK/BA are P3 — noted at the end)
**Prerequisite:** Sprint 2 — the throttle and model tiering from Sprint 2 wrap whatever
client is underneath; this sprint swaps the client under them.

**Why:** the current `QwenService` speaks the OpenAI-compatible endpoint. Moving to the
native DashScope protocol unlocks the capabilities the compatible layer hides or degrades
(explicit context cache, rerank) and is itself an innovation signal. **There is no Go SDK**
(SDKs exist only for Python/Java — `Docs/QwenDocs/api-reference/preparation/02-install-sdk.md`);
the client is hand-written. A Python sidecar is rejected: no second runtime in the
analysis hot path.

---

## Task 3.1 — `DashScopeService`: the native client (DS-1..3, DS-6)

**Files:** new `backend/internal/services/dashscope_service.go` (+ `_test.go`);
`backend/internal/services/qwen_service.go` stays untouched as the fallback

**Steps:**

1. Protocol (`Docs/QwenDocs/api-reference/`):
   - `POST {base}/api/v1/services/aigc/text-generation/generation`
   - Request: `{"model": …, "input": {"messages": […]}, "parameters": {…}}`
   - Response: `{"output": {"choices": […]}, "usage": {"input_tokens": …, "output_tokens": …}}`
   - Headers: `Authorization: Bearer $QWEN_API_KEY`; embeddings + rerank have their own
     endpoints — check the exact paths in `Docs/QwenDocs` before writing them.
2. Define a **shared interface** extracted from what callers actually use of
   `QwenService` today (chat completion, tool-calling loop contract, embeddings, JSON
   response format). Both services implement it. Keep the interface minimal — only
   methods with two real implementations.
3. **Tool calling (DS-3):** the native protocol carries `tools` / `tool_choice` inside
   `parameters` and tool calls inside `output.choices[].message`. Map to/from the
   existing structs so `RunAgentLoop` and `QuillExecutor` (`agent_tools.go`) work
   **unchanged** — the agent loop must not know which wire protocol runs underneath.
4. Config: base URL, timeouts, model IDs from env (DS-6).
5. Tests with `httptest.Server` fixtures: request-shape golden test, response parsing,
   tool-call round trip, error/429 surfacing (the Sprint 2 throttle handles retry).

**Expected result:** `DashScopeService` passes the same behavioural tests as
`QwenService`; the agent loop runs against a mocked native endpoint unchanged.

---

## Task 3.2 — Incremental cutover behind a flag (DS-4, DS-5)

**Files:** `backend/cmd/server/main.go` (composition root), config

**Steps:**

1. `LLM_PROTOCOL=openai|dashscope` selects the implementation at wiring time — one
   `if` in `main.go`, both satisfy the Task 3.1 interface. A regression rolls back with a
   config change, not a redeploy.
2. **Usage surfacing (DS-4):** parse `usage.input_tokens` / `output_tokens` **and the
   cached-token counts** into the token-usage logging (NF-4). Cache effectiveness must be
   measurable before Task 3.3 claims anything.
3. Cut over service by service if anything wobbles: extraction first (highest volume,
   simplest calls), agent loop last.

**Expected result:** the full stack runs green with `LLM_PROTOCOL=dashscope`; flipping
back to `openai` also runs green. `make e2e` passes on both.

---

## Task 3.3 — Explicit context cache (CC-1..4)

**Files:** prompt assembly in `analysis_service.go` / craft-review path; `dashscope_service.go`

**Steps:**

1. **Prefix discipline (CC-1):** restructure prompts as `[stable prefix | variable suffix]`:
   - prefix: system instructions + active skill body + universe lore + entity context
   - suffix: the paragraph under analysis
   The assembly code must be deterministic — same universe state ⇒ byte-identical prefix.
   No map-iteration ordering, no timestamps in the prefix. This is the entire trick:
   cache hits require prefix stability.
2. Mark the prefix for **explicit cache** per `Docs/QwenDocs` (cached input bills at 10%
   of standard input price; guaranteed hit window ~5 minutes — ideal for a writing
   session hitting the same universe repeatedly).
3. **Measure (CC-3):** log cache-hit tokens vs total input tokens per call; expose a
   counter. An unverified optimisation is not an optimisation.
4. **Constraint (CC-4):** never combine cache with Batch on one request — the discounts
   are mutually exclusive. Encode as a guard, not a comment.

**Expected result:** two consecutive analyses in the same universe show cached tokens > 0
on the second call, visible in logs.

---

## Task 3.4 — Structured output for extraction (SO-1..3)

**Files:** `dashscope_service.go` (JSON Schema response format), extraction call site

**Steps:**

1. Entity extraction sends a JSON Schema where `type` is an **`enum` of exactly the 7
   canonical types** from Sprint 0. An out-of-vocabulary type becomes *structurally
   impossible* — the third enforcement layer after prompt criteria and DB CHECK.
2. Relationship extraction likewise: relationship types constrained by enum
   (`ALLY_OF`, `MEMBER_OF`, `LOCATED_IN`, … — take the list from migration `014`'s edge
   types plus any the prompt currently allows). This **composes with**
   `validCypherIdentifier` in `graph_repo.go`; it does not replace it (NF-2).
3. Schema violations: log with the offending payload, drop that one item, never the batch
   (SO-3).

**Expected result:** extraction cannot return an invalid entity/relationship type; a
forced-invalid fixture test shows the violation logged and skipped.

---

## Task 3.5 — Rerank after fusion (RR-1..4)

**Files:** `backend/internal/services/memory_service.go` (`Recall`/`RecallWithPipelines`),
new rerank client method in `dashscope_service.go`, `fuse_rrf_explain.go`

**Steps:**

1. After `fuseRRF` produces the fused top-N, call the rerank endpoint (`gte-rerank` per
   `Docs/QwenDocs`) with the query + each item's text; reorder by rerank score. RRF fuses
   by **rank position** and is semantically blind; the reranker actually reads.
2. **Nil-safe optional wiring (RR-2):** `SetReranker(…)` following the exact pattern of
   `SetConsolidationRepo` / `SetBudgetMgr`. Absent ⇒ recall behaves as today.
3. **Measure (RR-3):** run the eval harness before and after —
   `TEST_DATABASE_URL=… QWEN_API_KEY=… go test ./backend/eval/ -run TestMemoryEval -v`.
   Record Recall@k, MRR, nDCG. A claimed improvement without a number does not go in the
   submission; if rerank does not improve the metrics, ship it disabled and say so.
4. **Explainability (RR-4):** extend `RecallExplain` to expose each item's pre- vs
   post-rerank position, so the Memory Theater's FusionExplorer can show what rerank
   changed (frontend chip: "↑3 rerank").

**Expected result:** eval metrics before/after recorded; explain payload carries rerank
deltas; reranker off ⇒ bit-identical legacy behaviour.

---

## P3 notes (do not build now)

- **Thinking mode (TK-1..2):** enable for contradiction/plot-hole when reaching P3;
  reasoning content is diagnostic only, never persisted as user-visible analysis.
- **Batch (BA-1..2):** only ever for consolidation / corpus-wide stylometry. **Never**
  ingestion or live analysis — the SLA is "usually within 24 hours". BA-2 is a standing
  rule, enforce with a guard if Batch is ever wired.

---

## Definition of done

- [ ] `DashScopeService` speaks the native protocol; agent loop unchanged; behavioural tests green against fixtures.
- [ ] `LLM_PROTOCOL` flag switches implementations; `make e2e` green on **both**.
- [ ] Token usage (incl. cached tokens) logged per call.
- [ ] Second same-universe analysis shows cache hits; hit rate visible in logs/metrics.
- [ ] Extraction constrained by JSON Schema enums (entities + relationships); violations logged, not silently dropped.
- [ ] Rerank wired nil-safe; eval numbers (Recall@k / MRR / nDCG) recorded before/after; explain payload shows rerank deltas.
