# Quill — Sprint 01 Action Plan & Codebase Audit

**Track**: MemoryAgent (Global AI Hackathon Series with Qwen Cloud)
**Deadline**: Jul 9, 2026, 2:00 PM PT (~4 days remaining)
**Repo**: `/home/daikyri/Workspace/Hackathon-QwenCloud`
**Status of this document**: read-only audit — no code changed during analysis.

---

## 1. Rules Summary

### Three core capabilities judges look for
The MemoryAgent track evaluates a persistent-memory agent on three named capabilities (`Docs/memoryagent-track-rules.md:10-14`):

1. **Efficient memory storage and retrieval** — how memories are indexed, stored, and fetched *at scale* (explicitly "not just a flat log of messages").
2. **Timely forgetting of outdated information** — mechanisms to decay, prune, or overwrite stale memories ("a genuine forgetting/decay mechanism, not just unlimited accumulation").
3. **Recalling critical memories within limited context windows** — smart selection/compression so the most relevant memories fit a constrained prompt budget (ranking, summarization, compression, relevance scoring).

### Judging criteria and weights (Stage Two, scored)
| Criterion | Weight | What it rewards |
|---|---|---|
| Innovation & AI Creativity | 30% | Sophisticated Qwen Cloud API use (custom skills, MCP-style integrations); novel memory architectures, custom retrieval components, storage/recall optimizations. |
| **Technical Depth & Engineering** | **30%** | *Core focus.* Architecture quality/modularity/scalability around the memory subsystem; clean non-trivial logic (indexing/embedding pipelines, dedup, conflict resolution old-vs-new); advanced patterns (vector stores, hybrid retrieval, summarization/compression, context-window budgeting). |
| Problem Value & Impact | 25% | Authentic user/business pain point (long-term personalization, cross-session continuity); productization/community potential. |
| Presentation & Documentation | 15% | Memory logic (storage, decay, recall) clearly *visualized* in the demo video; architecture doc explains memory design decisions. |

### Judging process
- **Stage One (pass/fail)**: does the project reasonably fit the MemoryAgent theme and reasonably apply Qwen Cloud APIs/SDKs? Failing this means the submission is never scored.
- **Stage Two (scored)**: the four weighted criteria above.

### Submission requirements (all mandatory unless noted)
1. Project built with Qwen models on Qwen Cloud.
2. **Public, open-source repo** with all source + setup instructions, and a **detectable open-source license file visible in the repo's "About" section** (a README mention is not sufficient — GitHub detects a real license file).
3. Text description of features/functionality.
4. **Proof of Alibaba Cloud deployment** — a *direct link to a code file* demonstrating use of Alibaba Cloud services/APIs.
5. **Architecture diagram** — a clear visual showing Qwen Cloud ↔ backend ↔ database ↔ frontend (memory store included).
6. **Demo video ≤ 3 minutes**, hosted publicly (YouTube/Vimeo/Youku), showing the project actually functioning, no unauthorized third-party media.
7. **Track identification** — must state "Track 1: MemoryAgent."
8. *(Optional bonus)* Blog/social post documenting the build journey.

Additional gates: project new-or-significantly-updated during the submission window; functional/testable/non-malicious; no IP infringement; English (or English translations).

### Prize
Grand Prize Track 1: $7,000 cash + $3,000 cloud credits + blog feature + swag (×1). Also independently eligible: Top 10 Honorable Mention and Top 10 Blog Post Award ($500+$500 each).

---

## 2. Codebase Status

Verified by reading source directly, not inferring from commit messages. Quill is **not a skeleton** — the memory stack is substantially built and wired end-to-end. Findings below separate what genuinely runs from what exists-but-is-dead.

### Genuinely built and wired end-to-end

**Apache AGE knowledge graph (real graph store, not a flat log)**
- One graph per universe, named `universe_<uuid>` — created at `backend/internal/repositories/graph_repo.go:90-97` (`CreateGraph`).
- Cypher injection safety: graph names are UUID-derived (injection-safe by construction); interpolated property values are escaped via `escapeCypherString` (`graph_repo.go:51-55`) because AGE doesn't support parameterized queries inside `$$ $$` blocks. This is a mature choice — a naive implementation would be Cypher-injectable.
- Real traversal primitives: `GetNeighbors` (`graph_repo.go:171-191`), `FullQuery` (`:194-209`), `NHopTraversal` BFS (`:222-237`), `DeleteEdge`, `DropGraph`.
- The demo universe deep-clone even reconstructs the AGE graph with entity-ID remapping by parsing agtype edge rows (`backend/internal/services/demo_service.go:392-460`).

**pgvector vector store**
- `paragraph_embeddings` and `entity_embeddings` tables (migrations `007`, `008`).
- Backs both `MemoryService.Recall`'s freshness signal and the agent's `search_vector_memory` tool.

**Weighted hybrid recall (the standout memory feature)**
- `MemoryService.Recall` (`backend/internal/services/memory_service.go:45-157`) merges three signals into one ranked list: graph neighbours ×0.4 + recency (relevance_score) ×0.3 + vector freshness ×0.3, then sorts descending and truncates to `k`.
- This is a genuine recall-ranking algorithm — exactly the "ranking, relevance scoring" language from the rules — not a bare vector top-k.
- It is wired into the live pipeline: `analysis_service.go:353-366` calls `Recall` and pushes results as a `contextual_recall` WS message.
- It is **visible in the UI**: `frontend/src/components/context-panel/ContextPanel.tsx:84-94` renders recall cards with a confidence percentage. The frontend dispatch handles it at `frontend/src/stores/wsStore.ts:102-103`.

**Tool-calling ReAct agent loop (not single-shot prompting)**
- `QwenService.RunAgentLoop` (`backend/internal/services/qwen_service.go:478-564`): OpenAI-style function calling, feeds tool results back as `role: "tool"` messages, loops until the model stops calling tools or `maxDepth` is exhausted; falls back to a single completion when no tools are supplied.
- `QuillExecutor` (`backend/internal/services/agent_tools.go`) implements two real tools via switch dispatch: `search_vector_memory` (embeds query → pgvector similarity, `:61-93`) and `query_entity_graph` (resolves entity name → ID → AGE neighbours, `:100-146`).

**Multi-model Qwen Cloud usage**
- `qwen-turbo` for entity extraction (`qwen_service.go:225-271`) and relationship analysis (`:320-367`); `qwen-max` for contradiction detection (`:383-449`) and the agent loop; `text-embedding-v3` for embeddings (`:273-318`). Concurrency-bounded via per-model semaphores (`:136-173`).

**Analysis services, all reachable from the composition root**
- `EntityService`, `ContradictionService`, `TimelineService`, `PlotHoleService` follow the repo→service→handler layering and push typed WS messages the frontend actually handles (`wsStore.ts:82-107`).
- `PlotHoleService.Scan` (`backend/internal/services/plot_hole_service.go:43-80`) genuinely consumes `relevance_score` as a signal: it skips background entities (`score <= 0.5`, `:62-66`) so it can distinguish "character quietly written out" from "arc abandoned mid-story." This is real engineering sophistication tied to the product thesis.

**Frontend surface**
- WS dispatch handles `analysis_result`, `contradiction_alert`, `entity_discovered`, `contextual_recall`, `graph_updated` (`wsStore.ts:82-107`).
- `ContextPanel` renders contradictions, discovered entities, recall items, and graph pings as dismissible cards (`ContextPanel.tsx:37-113`).

### Built but functionally dead — the single most important finding

**The forgetting/decay mechanism does not execute in production.**
- `RelevanceService.DecayAll` (`backend/internal/services/relevance_service.go:73-88`) implements exponential decay (`score *= e^(-lambda * idle)`, helper `applyDecay` at `:92-94`) and archives entities below `archiveThreshold`.
- It is constructed in the composition root (`backend/cmd/server/main.go:92`) and covered by tests (`relevance_service_test.go`, `entity_repo_test.go`).
- **But `DecayAll` is never called from any handler, service, or scheduler in the running application.** Grepping every caller across `backend/` returns only the two functions themselves plus test files — no production call site. There is no ticker, cron, or chapter-advance hook invoking it.
- `Touch()` (bump relevance *up* on mention) *is* called live (`analysis_service.go:251`). So relevance only ever increases in the running app; nothing decays it downward. The capability judges name explicitly — "timely forgetting" — is present in code but inert at runtime.

**Even if wired, decay is invisible.**
- Grep for `relevance_score`, `relevanceScore`, `archived`, or `decay` across `frontend/src` returns **zero files**. No entity card shows a relevance bar, no "archived/dormant" badge, no decay indicator. There is nothing for a judge to *see* the forgetting mechanism do.

### Submission-requirement gaps (verified directly)

- **No root `LICENSE` file.** Only `frontend/node_modules/**/LICENSE` exist. `README.md:56-58` claims "MIT" in prose, but GitHub's "About" detection needs an actual license file — this is a Stage-One mechanical risk.
- **No demo video and no video script** found anywhere in the repo.
- **Architecture diagram exists but is buried.** Mermaid `graph TB` at `Docs/SRS.md:49-90` and an ERD at `Docs/SRS.md:159-...`. GitHub renders Mermaid natively, so the requirement is technically met — but it's inside an ~11-day-old planning doc and is not linked from `README.md`.
- **Alibaba Cloud proof is arguably already satisfied at the code level.** `backend/internal/config/config.go:45` defaults `QwenBaseURL` to `https://dashscope-intl.aliyuncs.com/compatible-mode/v1` — a direct code-file link demonstrating use of an Alibaba Cloud API, which is exactly what requirement #4 asks for. What's *not* proven is running infrastructure (ECS/ACK) if the team wants the stronger "deployed on Alibaba Cloud" claim.
- **Stale planning checklist.** `Docs/PRD.md:623-635` still marks license, demo video, ECS deployment, and demo URL as "Pendiente" — written ~Jun 28. Do not trust it; re-verify each item's real current state before relying on the table.

---

## 3. Critical Findings

### Strengths that directly score points

1. **Hybrid, hierarchical memory is real** — relational (Postgres) + vector (pgvector) + graph (AGE), not a message log. This is precisely what Technical Depth rewards ("efficient storage/retrieval at scale, not just a flat log") and is the hardest thing to fake. *Scores: Technical Depth (30%), Innovation (30%).*
2. **`MemoryService.Recall`'s weighted multi-signal ranking** (`memory_service.go:45-157`) is a concrete, deterministic implementation of "smart recall within a limited context window" — and it's surfaced in the ContextPanel UI, so it's demoable, not just theoretical. *Scores: Technical Depth, Presentation.*
3. **Tool-calling ReAct agent loop** (`RunAgentLoop` + `QuillExecutor`) is a legitimate "custom skills / MCP-style" pattern. Most hackathon entries do single-shot prompts; this iterates with tool feedback. *Scores: Innovation (30%).*
4. **`PlotHoleService` using relevance to tell narrative abandonment from graceful exit** is a genuinely novel product idea grounded in the memory model — a strong "authentic user pain point" narrative for a creative-writing IDE. *Scores: Problem Value (25%), Innovation.*
5. **AGE graph-per-universe with Cypher escaping** demonstrates engineering maturity and multi-tenant isolation. *Scores: Technical Depth.*

### Weaknesses that will be penalized

1. **Dead forgetting mechanism — the worst possible gap given this specific track.** The rules name "timely forgetting" as one of three core capabilities, and Technical Depth explicitly asks for "a genuine forgetting/decay mechanism (not just unlimited accumulation)." Right now the code builds it, tests it, and never calls it. A judge who reads the source (Technical Depth graders do) will find that `DecayAll` has no production caller almost immediately. This undercuts the single most track-defining claim the project can make.
2. **Forgetting is invisible even if wired.** Zero frontend references to relevance/archived/decay. Presentation is graded on whether "the memory logic (storage, decay, recall) is clearly visualized in the demo video" — decay currently has nothing to show. This is a double failure: not executing *and* not visible.
3. **Missing root LICENSE file.** A mechanical Stage-One risk that costs ~5 minutes to remove but could sink the whole submission if a checker relies on GitHub's "About" detection.
4. **No demo video and no evidence one is planned.** This is a hard submission requirement, not a scoring nicety — without it the submission may not be judged at all. With 4 days left and no script, this is schedule risk.

### Over-engineering / unnecessary effort (leave alone)

1. **`PlotHoleService.Scan` N+1 chapter-order lookups** (`plot_hole_service.go:41-42, 68-72`) — already flagged with a `ponytail:` comment as acceptable at hackathon scale. Correctly deprioritized; do **not** "optimize" it before the deadline.
2. **Demo universe deep-clone** (`demo_service.go`, 10 tables + AGE graph) — thorough finished infrastructure for stable multi-session demos. Done; do not touch.
3. **Two-phase hub/submitter wiring** (`main.go`, `ws.NewHub` with nil submitter → `SetSubmitter`) — a deliberate minimal solution to a circular dependency, not over-engineering. Do not refactor.
4. **1-hop cap in `MemoryService`/graph traversal** — explicitly commented as a diminishing-returns choice. Going deeper is scope creep with no judging payoff.

### Concrete opportunities (ranked by ROI for ~4 days)

1. **Wire `DecayAll` into a real trigger.** Cheapest hook: call `relevSvc.DecayAll(ctx, universeID)` on chapter creation/advance (in `ChapterService` or the analysis path near the existing `Touch` call, `analysis_service.go:251`). ~10-20 lines, reuses already-tested code. Highest single ROI in the whole project.
2. **Surface relevance + archived state in the UI.** Minimum version: an "Archived/Dormant" badge and a relevance bar on entity cards / Context Panel, from data already returned by `ListByUniverseActive` and the entity endpoints. Turns an invisible backend mechanism into a demo moment.
3. **Add a root `LICENSE` (MIT).** One file, near-zero effort, removes a Stage-One risk.
4. **Link the existing Mermaid diagram from the README.** Reuse `Docs/SRS.md:49-90` under a README "Architecture" heading — do not redraw it.
5. **Script and record the 3-minute demo** ending on decay/archive visibly happening. Presentation (15%) is graded on visualized memory logic; decay is currently the weakest-shown of storage/decay/recall.

---

## 4. Sprint 1 Plan

Phased by impact-on-score vs. effort. Do Critical first (or risk not being judged), then High-impact (moves the two biggest-weighted criteria), then Polish.

### Phase A — Critical (must-do to pass Stage One / avoid disqualification) — Day 1

| Task | Why (criterion) | Effort | Dependencies |
|---|---|---|---|
| Add root `LICENSE` file (MIT) | Submission req #2 — detectable license; Stage One | ~5 min | none |
| Link the existing architecture diagram (Mermaid from `Docs/SRS.md:49-90`) into `README.md` under "Architecture" | Submission req #5; Presentation (15%) | ~15 min | none |
| Decide + execute Alibaba Cloud proof: either cite `config.go:45` DashScope endpoint as the required code-file link, **or** actually deploy to ECS if the team wants the stronger claim | Submission req #4; Stage One | ~30 min (link) or 2-4 hrs (ECS) | none — recommend the link unless team explicitly wants deployed infra |
| Clean `docker compose up -d` smoke test with a real `QWEN_API_KEY` (verify end-to-end: write → analysis → WS push) | Stage One — "functional, testable" | ~30 min | Qwen key available |

### Phase B — High-impact (moves Technical Depth 30% + Presentation 15%) — Day 1-2

| Task | Why (criterion) | Effort | Dependencies |
|---|---|---|---|
| Wire `RelevanceService.DecayAll` into a real runtime trigger (chapter-advance hook, near `analysis_service.go:251`) | Technical Depth — "genuine forgetting mechanism"; directly fixes the #1 weakness | 1-2 hrs (reuses tested code) | none |
| Surface relevance score + active/archived badge on entity UI (Entity Card and/or Context Panel), sourced from existing entity data | Technical Depth + Presentation — makes forgetting "visibly demonstrable" | 2-3 hrs | decay wiring (B1) so the badge actually changes |
| Script + record 3-min demo: live write → contradiction alert → plot hole after N silent chapters → **decay/archive visibly happening** → knowledge-graph update | Presentation (15%); reinforces Problem Value (25%) | 3-4 hrs incl. edit | B1 + B2 (decay must be visible on screen) |

### Phase C — Polish (only if time remains) — Day 3-4

| Task | Why (criterion) | Effort | Dependencies |
|---|---|---|---|
| Blog post documenting the build journey | Optional bonus prize ($500 + $500) | 2-3 hrs | demo done (reuse footage/screens) |
| Update `Docs/PRD.md:623-635` checklist to reflect real current state (drop stale "Pendiente") | Presentation clarity for doc-reading judges | ~30 min | Phase A/B done |
| Decay-timeline mini-viz (sparkline of relevance over chapters) | Technical Depth polish | half-day | B1/B2 — **cut first if squeezed** |

---

## 5. Scope Cuts

Explicit "do NOT do" list — each saves time that should go to Phases A/B.

- **Do not refactor `PlotHoleService`'s N+1 query** (`plot_hole_service.go:41-42`). Already flagged acceptable at hackathon scale; zero judging payoff.
- **Do not build real ECS/ACK infrastructure** unless the team explicitly wants the "deployed on Alibaba Cloud" claim beyond "calls an Alibaba Cloud API." Requirement #4 is satisfiable today via the `config.go:45` DashScope endpoint link. Full deployment is high effort for a requirement that's arguably already met.
- **Do not add new AI features** (more entity types, more agent tools, deeper graph hops). `NHopTraversal`/`MemoryService` are intentionally 1-hop with a documented diminishing-returns rationale. Adding surface area this late is pure risk.
- **Do not redesign the architecture diagram.** The existing `Docs/SRS.md` Mermaid diagrams are accurate and GitHub-renderable. Surface them; don't recreate them as a fancier image.
- **Do not build a background decay-sweep worker/cron.** Event-driven decay on chapter advance is the intended design (documented in the `ponytail:` comment at `relevance_service.go:15-17`). A scheduler is unnecessary complexity for a demo — it also makes decay *harder* to show on cue in the video.
- **Do not touch the two-phase hub wiring or the demo deep-clone.** Both are finished, deliberate, and load-bearing for the demo.

---

### Key file references for follow-up work
- `backend/internal/services/relevance_service.go:73-94` — decay logic to wire.
- `backend/internal/services/analysis_service.go:251` — live `Touch` call; decay hook belongs on the chapter-advance path near here.
- `backend/cmd/server/main.go:92,114` — `relevSvc` construction + injection.
- `backend/internal/services/plot_hole_service.go:62-66` — existing `relevance_score` consumer (pattern to mirror in UI).
- `backend/internal/services/memory_service.go:45-157` — weighted hybrid recall (strength to showcase).
- `backend/internal/services/agent_tools.go`, `qwen_service.go:478-564` — tool-calling agent loop.
- `frontend/src/components/context-panel/ContextPanel.tsx`, `frontend/src/stores/wsStore.ts:82-107` — where relevance/archived UI should land.
- `Docs/SRS.md:49-90` — architecture diagram to link from README.
- `backend/internal/config/config.go:45` — Alibaba Cloud (DashScope) proof link.

---
---

# Part 2 — Technical Depth Audit (Implementation Only)

**Scope note**: This part covers ONLY technical implementation depth — submission logistics (license, deployment, video, README) are covered in Part 1 above and not repeated here. All findings are read-only, evidence-backed by file:line, and net-new relative to Part 1 unless marked "confirmed from prior audit."

## 1. Codebase Map

### Backend layers (all wired from `backend/cmd/server/main.go:66-137`)
- **Repositories** (18 files, pgx SQL): `graph_repo.go` (Apache AGE), `entity_repo.go`, `vector_repo.go` (pgvector), `contradiction_repo.go`, `timeline_repo.go`, `plot_hole_repo.go`, `ingestion_repo.go`, plus CRUD repos for user/universe/work/chapter.
- **Services** (25 files): `AnalysisService` (per-work sequential queue), `EntityService` (3-tier dedup), `ContradictionService` (deterministic + semantic agent loop), `RelevanceService` (decay/touch, built but inert — see §2), `TimelineService`, `PlotHoleService`, `MemoryService` (weighted hybrid recall), `QwenService` (agent loop + multi-model dispatch), `IngestionService` (async doc pipeline — built but frontend-invisible, see §2), `DemoService`.
- **Handlers** (19 files): full REST surface for auth/universe/work/chapter/entity/contradiction/timeline/plot-hole/graph/ingestion/demo/health.
- **WS layer**: `hub.go` (per-user single-conn map, heartbeat, auth handshake), `protocol.go` (message type constants).

### Frontend
- Zustand stores: `authStore`, `universeStore`, `editorStore`, `graphStore`, `wsStore` (single WS connection, dispatches by `msg.type`, `frontend/src/stores/wsStore.ts:79-111`).
- `ContextPanel.tsx` renders contradictions, discovered entities, recall items, graph pings.
- **No ingestion/upload UI anywhere** — `grep -r "ingest|upload" frontend/src/**/*.tsx` returns zero matches. The entire document-ingestion pipeline is backend-only and unreachable from the product.

### What genuinely works end-to-end
- Auth → universe → work → paragraph_submit → entity extraction → contradiction check → graph write → WS push is a real, wired pipeline, not stubs. `handlers/e2e_test.go` exercises register→login→create-universe→create-work against a live Postgres integration DB (skipped without `TEST_DATABASE_URL`, but the test exists and is not vaporware).
- Apache AGE graph-per-universe, pgvector embeddings, and the ReAct tool-calling agent loop (`QwenService.RunAgentLoop` + `QuillExecutor`) are all real, non-trivial implementations — confirmed by reading `qwen_service.go:478-564` and `agent_tools.go` directly.

## 2. Technical Findings — Broken, Incomplete, Fragile

### A. Ingestion progress is silently dead (confirmed bug, not just "invisible")

`backend/internal/services/ingestion_service.go:243-262`:
```go
func (s *IngestionService) emitProgress(jobID uuid.UUID, status string, processed, total int) {
	...
	// ponytail: hub.SendToUser requires userID. Ingestion is system-initiated,
	// so userID is empty (uuid.Nil). The hub stores conns by userID.
	_ = s.hub.SendToUser(uuid.Nil, msg)
}
```
`Hub.GetConn` (`ws/hub.go:104-108`) only ever registers real authenticated userIDs (`Register(userID, conn)` in the auth handshake). `GetConn(uuid.Nil)` always returns `nil`, so `SendToUser` always returns `"user %s not connected"` — an error that's discarded (`_ =`). **Every ingestion_progress event ever emitted is a no-op.** Combined with zero frontend references to ingestion at all, this is a fully-built, fully-wired-on-paper pipeline that is completely invisible and non-functional from the user's perspective. Also note `SaveParagraphEmbedding(ctx, uuid.Nil, ...)` (`ingestion_service.go:144`) tags every ingested paragraph with a placeholder `chapterID = uuid.Nil` — ingested content's embeddings are never associated with a real chapter.

### B. `MemoryService.Recall`'s freshness signal is computed wrong

`backend/internal/services/memory_service.go:110-130`:
```go
for range entities {
	if len(queryEmbedding) == 0 {
		break
	}
	matchedID, distance, err := s.vectorRepo.FindSimilarEntity(ctx, universeID, queryEmbedding, 0.8)
	...
	if item, exists := candidateMap[*matchedID]; exists {
		item.Score += freshScore
	}
	...
}
```
The loop variable is discarded (`for range entities`, not `for _, e := range entities`). This calls the identical single top-1 vector match **once per entity in the universe** and keeps re-adding `freshScore` to the same candidate. Net effect: one entity's freshness contribution is silently multiplied by the universe's entity count instead of being computed per-entity against the query. This doesn't corrupt the live per-paragraph pipeline (`analysis_service.go:353` calls `Recall(ctx, job.UniverseID, nil, 5)` — nil embedding short-circuits on the first iteration) but it does corrupt the WS `recall_request` path and `POST /universes/:id/recall` — the two places a user/judge would actually issue a real semantic-recall query. This is the project's flagship "smart recall ranking" algorithm, and its third signal is broken by a copy-paste loop bug.

### C. Deceased/alive contradiction check is structurally defeated by merge ordering

`analysis_service.go:220-246` calls `extractEntities` (step 1) before `contraSvc.CheckDeterministic` (step 2). `extractEntities` → `EntityService.ResolveOrCreate` → `mergeEntity` (`entity_service.go:189-221`):
```go
// Update status if provided
if newData.Status != "" {
	merged.Status = newData.Status
}
```
This **unconditionally overwrites and persists** the entity's status with whatever the LLM extracted for *this* mention, before returning the resolved entity. Then `CheckDeterministic` (`contradiction_service.go:66-110`) checks:
```go
if e.Status == "deceased" {
```
— on the entity that was **just overwritten** in step 1. The extraction prompt (`qwen_service.go:226-240`) has the model emit a `status` field for essentially every character (example shows `"status": "active"` as the default), so on any post-death re-mention, `newData.Status` will be non-empty ("active"), clobbering the DB's "deceased" status *before* the deceased check runs. **The deterministic deceased/alive contradiction rule — the one concrete rule-based check in the whole system — can structurally almost never fire**, because by the time it inspects `e.Status`, it's already looking at the new (wrong) value, not the old one it's supposed to compare against. Fix requires either snapshotting pre-merge status, or reordering `CheckDeterministic` to run against a freshly-fetched pre-merge entity.

### D. Zero test coverage on the two most "innovation-scoring" pieces

- `backend/internal/services/agent_tools.go` (the `QuillExecutor` ToolExecutor dispatch — the concrete "custom skills / tool-calling" implementation judges are told to reward) has **no `_test.go`**, despite being deliberately split into small mockable interfaces (`vectorSearcher`, `graphQuerier`, `entityLister`) for exactly that purpose.
- `backend/internal/services/entity_service.go` (the 3-tier dedup: exact name → alias → pgvector semantic similarity @ 0.15 threshold → create, `entity_service.go:85-187`) — genuinely the strongest "conflict resolution, old-vs-new" implementation in the codebase — has **no `_test.go`**.
- Also untested: `universe_service.go`, `chapter_service.go`, `work_service.go`. Handlers without tests: `entity.go`, `chapter.go`, `work.go`, `universe.go`, `auth.go`, `demo.go`.
- By contrast, `qwen_service_test.go`, `contradiction_service_test.go`, `timeline_service_test.go`, `plot_hole_service_test.go`, `relevance_service_test.go`, `analysis_service_test.go`, `demo_service_test.go`, and most repos DO have tests — coverage is uneven, not absent, but the gaps are on exactly the pieces judges are told to weight most.

### E. WS type parity — one silent gap, one acceptable

`ws/protocol.go:16-26` defines `TypeIngestionProgress` and `TypeAuthError`. `wsStore.ts`'s `dispatch` switch (`wsStore.ts:81-110`) has no `case` for either — both fall into `default: break` and are silently dropped. `auth_error` is low-risk (the connection closes right after anyway). `ingestion_progress` compounds finding A — even if the backend bug were fixed, the frontend still wouldn't render it.

### F. Agent tool entity resolution is exact-match only

`queryEntityGraph` (`agent_tools.go:100-134`) resolves an entity name via `if ent.Name == args.EntityName` — no alias or fuzzy matching, even though `EntityService.ResolveOrCreate` itself supports alias resolution. The agent's own tool can miss entities that the dedup pipeline would happily have found under an alias.

### Confirmed still true from Part 1 (not re-litigated in depth here)
- `RelevanceService.DecayAll` has no production caller (only test call sites) — forgetting logic is built, tested, and inert at runtime.
- Zero frontend references to relevance score/archived/decay anywhere in `frontend/src`.

### Race conditions — checked, none found in the analysis pipeline
`AnalysisService.SubmitParagraph`/`runWorker` (`analysis_service.go:92-203`) has an apparent TOCTOU on the `exists` check before spawning a worker goroutine, but `runWorker` re-checks under `s.mu` via the `cancels` map (`analysis_service.go:163-167`) before proceeding — double-launch is correctly prevented. `broadcastResult` is always guarded by `s.hub != nil` before being called. No unguarded shared-state access found.

### Over-engineered / leave alone (confirmed from Part 1, still valid)
- `PlotHoleService` N+1 chapter-order lookups — explicitly `ponytail:`-flagged, correct call to leave it.
- Demo universe deep-clone infra — finished, load-bearing for demos, don't touch.
- Two-phase `hub`/`submitter` wiring in `main.go` — deliberate, not over-engineering.
- 1-hop graph traversal cap in `MemoryService`/`NHopTraversal` — documented diminishing-returns choice.

## 3. MemoryAgent Assessment (honest verdict per capability)

**1. Efficient storage/retrieval at scale** — **Real, genuinely built.** Postgres + pgvector + Apache AGE (graph-per-universe) is a real hybrid store, not a flat log. This is the strongest of the three claims and needs no further work beyond what's in Part 1 Sprint A/B.

**2. Timely forgetting/decay** — **Built but inert and invisible.** `DecayAll` exists, is tested, and is never called in production; nothing in the frontend shows relevance/archived state. Confirmed again in this pass — still the single biggest gap given this is a track-defining named capability.

**3. Smart recall within limited context windows** — **Real algorithm, but with a live bug undermining it.** `MemoryService.Recall`'s three-signal weighted ranking (graph×0.4 + recency×0.3 + freshness×0.3) is a genuine, demoable implementation — but the freshness signal is computed incorrectly (finding B above) for any query-driven recall path, which is exactly the path a judge would exercise if they tried the recall feature interactively rather than just watching the passive per-paragraph pipeline.

**Bonus — conflict resolution (explicitly named in Technical Depth criterion)**: `EntityService`'s 3-tier dedup is real and well-designed, but (a) it's untested, and (b) its own status-merge logic actively defeats the deceased/alive contradiction check (finding C) — this is the most concrete, fixable "technical depth" bug in the codebase because it directly touches the exact judging language ("conflict resolution old-vs-new").

## 4. Sprint Plan (4 phases, technical-depth focus only)

### Phase 1 — Fix what's broken (blocks core-flow correctness)

| Task | File/Function | Why | Effort | Deps |
|---|---|---|---|---|
| Fix `MemoryService.Recall` freshness loop — iterate `entities` and match query embedding per-entity (or drop the loop and call `FindSimilarEntity` once, then score directly) | `memory_service.go:110-130` | Flagship recall algorithm produces wrong scores on any real query — core "smart recall" claim | S | none |
| Fix deceased/alive check ordering — snapshot entity status *before* `ResolveOrCreate` merges it, pass old status into `CheckDeterministic`, or reorder to check before merge | `analysis_service.go:220-246`, `entity_service.go:189-221`, `contradiction_service.go:66-110` | The one deterministic contradiction rule can structurally never fire — direct "conflict resolution" judging language | M | none |
| Fix or cut ingestion progress delivery — either broadcast to all connected users of the universe's owner (needs a userID on the job) or drop `emitProgress` entirely and rely on REST job-status polling | `ingestion_service.go:243-262` | Dead WS message, wasted effort, misleading "real-time ingestion" claim in CLAUDE.md | S–M | needs a way to resolve universe owner → userID |
| Add minimal ingestion UI (upload button + progress bar) OR explicitly cut the feature from the demo narrative | `frontend/src` (new), `IngestionHandler.Ingest` | Currently built server-side with zero user-facing value; either finish it cheaply or don't claim it | M (UI) / S (cut) | ingestion WS fix above if keeping it |

### Phase 2 — Strengthen the MemoryAgent core (make it genuine AND visible)

| Task | Why | Effort | Deps |
|---|---|---|---|
| Wire `RelevanceService.DecayAll` into a real trigger (chapter-advance hook near `analysis_service.go:251`) | Track-defining "forgetting" capability is currently inert | S–M | none |
| Surface relevance score + archived/dormant badge in `ContextPanel`/entity cards | Turns backend-only decay into a demoable signal | M | decay wiring above |
| Add alias-aware resolution to `queryEntityGraph` (reuse `EntityRepo.FindByAlias` instead of exact `ent.Name ==`) | Agent tool currently misses entities its own dedup pipeline would find | S | none |

### Phase 3 — Engineering quality (tests, error handling)

| Task | Why | Effort | Deps |
|---|---|---|---|
| Add `agent_tools_test.go` (mock the three interfaces already designed for this) | Core innovation surface has zero coverage | S–M | none |
| Add `entity_service_test.go` covering the 3-tier dedup (exact/alias/semantic) and the status-merge fix from Phase 1 | Strongest conflict-resolution logic, currently untested; also regression-guards the Phase 1 fix | M | Phase 1 status-merge fix |
| Add a regression test proving the deceased/alive contradiction fires correctly post-fix | Directly demonstrates the fixed bug won't regress | S | Phase 1 fix |
| Add tests for `universe_service.go`, `chapter_service.go`, `work_service.go`, or explicitly accept the gap given time constraints | Coverage gap on CRUD services | M (if pursued) | none |

### Phase 4 — Polish (only if time remains)

| Task | Why | Effort | Deps |
|---|---|---|---|
| Add `ingestion_progress`/`auth_error` cases to `wsStore.ts` dispatch (even if just logging) | Silent message drops look sloppy to a code-reading judge | S | Phase 1 ingestion fix if keeping the feature |
| Document the freshness-loop fix and the deceased/alive ordering fix as explicit "before/after" bullets in the demo script or README | Presentation credit for genuinely fixing named judging-criterion bugs | S | Phase 1 |

## 5. Scope Cuts (Part 2)

- **Do not build a full ingestion UI from scratch under time pressure** — cheaper to either wire minimal progress display or explicitly cut the pipeline from the pitch; don't let it eat Phase 1/2 time.
- **Do not attempt full test coverage across all CRUD services/handlers** — prioritize `agent_tools.go` and `entity_service.go` only; those are the ones judges are told to weight (custom tools, conflict resolution).
- **Do not add a background decay-sweep worker/cron** — event-driven decay on chapter-advance remains the right call (confirmed from Part 1's `ponytail:` rationale).
- **Do not deepen graph traversal beyond 1-hop** — no judging payoff, documented diminishing returns.
- **Do not refactor `PlotHoleService`'s N+1** — already correctly deprioritized.
