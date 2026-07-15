# Sprint 6 — Editor as Source of Truth & MCP

**Tier:** P3 · **SRS coverage:** EH-1..4, EC-1..6, ID-1..6, MC-1..6
**Prerequisite:** Sprint 1 (live analysis works — candidates ride on it), Sprint 4/5 for
the tray-decisions-as-learning-signals wiring.
**Sacrifice order:** this sprint is cut first if capacity runs out. Within it, build in
the order below — autosave durability and Markdown export before highlighting cosmetics;
Quill-as-MCP-server before any external MCP.

---

## Task 6.1 — Durability first: autosave and recovery (ID-2..4)

**Files:** `frontend/src/components/editor/TipTapEditor.tsx` (save path),
`frontend/src/stores/editorStore.ts`

**Steps:**

1. Retry on save failure with backoff (the current 5s debounce has no retry).
2. Honest save state in the editor chrome: *saving / saved / failed*. A silent save
   failure is a data-loss event wearing a green checkmark.
3. Local recovery: persist unsaved content to `localStorage` (keyed by chapter) on every
   debounce; on load, if local is newer than server, offer restore. *Losing a chapter
   kills the product — there is no second chance with a writer who lost work.*

**Expected result:** kill the backend, keep typing, reload the page → the editor offers
the unsaved content back, and the UI never claimed it was saved.

---

## Task 6.2 — Import fidelity and Markdown export (ID-1, ID-5, ID-6)

**Files:** `backend/internal/services/ingestion_service.go` (chapter content storage),
editor content loading; new export endpoint + frontend affordance

**Steps:**

1. **Import (ID-1):** preserve paragraph structure, italics, and scene breaks from
   ingested `.md`/`.txt` into the stored chapter content, loading into TipTap as
   formatted content — not the current plain-text dump. Store as the editor's native
   rich format (TipTap JSON or HTML — whichever the editor already persists).
2. **Export (ID-5):** `GET /chapters/:id/export.md` and `GET /works/:id/export.md`
   (chapters concatenated with headers). TipTap content → Markdown. Without export the
   product is a trap, and no serious author will adopt a trap.
3. DOCX stays roadmap; PDF stays rejected (ID-6). Do not build them.

**Expected result:** ingest a Markdown fixture with italics and scene breaks → the editor
shows them; export round-trips the work back to Markdown a human can diff against the
original.

---

## Task 6.3 — Entity highlighting, zero LLM (EH-1..4)

**Files:** TipTap extension (new file under `frontend/src/components/editor/`),
entity data from `universeStore`/API

**Steps:**

1. Client-side match of known entity names + aliases against the visible text — a TipTap
   decoration plugin. **Zero LLM calls**: the AI did its work at extraction time; this is
   a string match against the known-entity set.
2. Stable colour per entity (hash entity id → the type's palette family from
   `ENTITY_TYPE_META`); an alias highlights identically to its canonical name.
3. Performance (EH-3): build one matcher (regex alternation or Aho-Corasick if needed)
   per entity-set change, decorate only visible ranges; typing latency must not
   measurably degrade — test with a 500-entity universe.
4. Click a highlight → entity context (route to the entity card or a popover with the
   Sprint 0 detail panel).

**Expected result:** known entities glow in the manuscript as the writer types, at zero
marginal cost, and clicking one opens its card.

---

## Task 6.4 — Entity candidates and the review tray (EC-1..6)

**Files:** extraction result handling (`analysis_service.go` / ingestion REDUCE),
new candidate status or table, editor underline UI, tray component

**Steps:**

1. Extraction produces **candidates**, not committed entities: add a
   `candidate` status (or confidence column) to the entity flow. The AI proposes; the
   human confirms.
2. **Confidence gate (EC-3, EC-5):** extraction returns confidence per entity (add to the
   structured-output schema). Above threshold (configurable) → auto-accept; below → tray.
3. **Live candidates (EC-2):** in the editor, a candidate is a discreet underline with an
   inline accept/dismiss affordance. No modals, no pop-ups.
4. **The tray (EC-4):** a panel listing low-confidence candidates from bulk ingestion —
   name, type, evidence quote, accept/merge/dismiss. The tray is the *destination* of the
   low-confidence bucket the gate creates; without it, confidence-gating would be
   decorative.
5. **Learning signals (EC-6):** tray and inline decisions are recorded as feedback events
   (Sprint 4 infrastructure) — corrections are training data, not just cleanup.

**Expected result:** ingesting a noisy fixture yields auto-accepted confident entities and
a tray of uncertain ones; accepting from the tray commits the entity into browser + graph.

---

## Task 6.5 — Quill as an MCP server (MC-1..4)

**Files:** new `backend/internal/mcp/` (server over the MCP Go SDK or a minimal
streamable-HTTP JSON-RPC implementation), wired in `main.go`

**Steps:**

1. Expose the memory subsystem as MCP tools, **reusing the existing implementations** —
   `QuillExecutor`'s `search_vector_memory` and `query_entity_graph`
   (`agent_tools.go`) plus `MemoryService.Recall`. **No new retrieval logic** (MC-2);
   these tools already run in production against the internal agent — this task is a
   transport, not a feature.
2. Tool surface: `search_memory(universe, query)`, `query_entities(universe, name)`,
   `recall(universe, query, k)` — thin adapters mapping MCP tool calls onto the existing
   functions.
3. **Auth and tenancy (MC-3, MC-4):** authenticate no weaker than the HTTP API (bearer
   token); every tool call scoped to a universe the authenticated user owns. A
   cross-tenant read is a security failure, not a bug.
4. Prove it from a real client: connect Claude Code (or any MCP host) and ask *"what
   colour were Elara's eyes?"* against the demo universe. That interaction **is** Success
   Criterion 5 and the demo clip.

**Expected result:** an external MCP client lists the tools, calls them with a token, gets
universe-scoped answers, and is refused on another user's universe.

---

## Task 6.6 — External MCPs, fallback mandatory (MC-5, MC-6) — build last

**Steps:**

1. **Notion MCP** (memory *source*): ingest a worldbuilding page as established lore via
   the official server (`mcp.notion.com/mcp`, OAuth). OAuth per user is the cost driver —
   timebox it; if OAuth eats the budget, ship Quill-as-server only.
2. **Wikipedia MCP** (research): available to craft review as a lookup tool. **Must
   degrade gracefully** — it is the only external network dependency in the system; if
   unreachable, the review completes without it. It must not be *able* to fail a review
   (timeout + absent-tool behaviour, tested by pointing it at a dead endpoint).

**Expected result:** with Wikipedia unreachable, reviews still complete; a Notion page
lands as universe lore (or the feature was consciously cut, documented in the submission).

---

## Definition of done

- [ ] Autosave retries, shows honest state, and recovers locally; the kill-reload test passes.
- [ ] Import preserves structure/italics/scene breaks; Markdown export exists per chapter and per work.
- [ ] Zero-LLM entity highlighting with stable per-entity colour and click-through; no measurable typing latency at 500 entities.
- [ ] Candidates flow: confidence gate, inline underline, review tray, decisions recorded as learning signals.
- [ ] Quill speaks MCP: existing tools exposed, token-authenticated, universe-scoped, demonstrated from an external client.
- [ ] External MCPs behind fallbacks — or consciously cut.
- [ ] `make e2e` green.
