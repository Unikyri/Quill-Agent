## Exploration: Quill Phase 2 — Memory Engine + Frontend Core

### Current State

Phase 0+1 is complete and archived. The project has:

**Backend (Go/Fiber):**
- All 12 database migrations deployed (users, universes, works, chapters, entities, entity_mentions, entity_embeddings, paragraph_embeddings, contradictions, timeline_events, plot_holes, ingestion_jobs) + seed demo
- Full CRUD handlers for auth, universe, work, chapter, entity
- EntityService with 4-step resolution (exact → alias → semantic → create)
- QwenService with semaphore rate limiting (max 3, turbo 5), ExtractEntities, GenerateEmbedding, AnalyzeRelationships
- GraphRepo with basic Cypher operations (bare-bones)
- Health endpoint with real DB/ext/Qwen probes
- Tests: auth_service, health_handler, CRUD repos

**Frontend (React/Vite):**
- Vite + React 18 + TipTap (in package.json but NOT wired — EditorPage uses plain `<textarea>`)
- Pages: LoginPage, DashboardPage, UniversePage, EditorPage
- `components/` directory is EMPTY — no Layout, no Sidebar, no MemoryPanel
- Stores: authStore, universeStore, editorStore (basic)
- `ReactFlow` is in package.json but not used
- Dark mode CSS

**Database** — the schema has ALL the tables Phase 2 needs:
- `contradictions` (severity, fingerprint, evidence_a/b, status)
- `timeline_events` (timeline_position, participants, chapter_id)
- `plot_holes` (related_entity_ids, first_mentioned_chapter_id, status)
- `paragraph_embeddings` (embedding vector(1024), chapter_id, paragraph_node_id)
- `entity_embeddings` (embedding vector(1024), entity_id)
- `entity_mentions` (paragraph_node_id for TipTap)

**Key existing code:**
- `EntityService.ResolveOrCreate()` — transactional entity resolution (exists)
- `VectorRepo.FindSimilarParagraphs()` — semantic search (exists)
- `EntityRepo.CreateMention()`, `CountMentions()`, `GetMaxMentionsInUniverse()` — mention tracking (exists)

### Affected Areas

- `backend/internal/services/qwen_service.go` — Add `DetectContradictionsBatch()` method (uses Qwen-Max)
- `backend/internal/services/chapter_service.go` — Add analysis trigger on save
- `backend/cmd/server/main.go` — Wire new services, WebSocket, new routes
- `backend/internal/services/analysis_service.go` — **NEW**: Sequential queue per chapter with context cancellation, core/enrichment separation
- `backend/internal/services/contradiction_service.go` — **NEW**: Fingerprint dedup, deterministic checks, batched LLM contradictions
- `backend/internal/services/relevance_service.go` — **NEW**: Exponential decay, touch/reactivation, decay_all
- `backend/internal/services/memory_service.go` — **NEW**: Contextual recall combining graph facts + recent mentions
- `backend/internal/services/timeline_service.go` — **NEW**: Temporal consistency validation
- `backend/internal/services/plothole_service.go` — **NEW**: Stale arc scanning
- `backend/internal/handlers/memory.go` — **NEW**: Contradiction CRUD, timeline, plot holes, graph endpoints
- `backend/internal/handlers/ws.go` — **NEW**: WebSocket hub with auth, paragraph_submit, events
- `backend/internal/repositories/graph_repo.go` — Expand: richer relationship ops, full graph query, neighbor traversal
- `backend/internal/config/config.go` — Add decay/timeline/plothole constants
- `backend/internal/models/models.go` — Add request/response types for memory/graph endpoints
- `frontend/src/pages/EditorPage.tsx` — Replace textarea with TipTap, add entity highlighting, debounce → WS
- `frontend/src/components/` — **NEW**: MemoryPanel, EntityTooltip, WebSocketStatus
- `frontend/src/stores/editorStore.ts` — Add TipTap instance, analysis state
- `frontend/src/stores/memoryStore.ts` — **NEW**: Contextual memories, contradictions, plot holes, graph updates
- `frontend/src/lib/api.ts` — Add memory/graph/timeline endpoints
- `frontend/src/lib/ws.ts` — **NEW**: WebSocket client with exponential reconnect
- `frontend/src/pages/KnowledgeGraphPage.tsx` — **NEW**: React Flow Knowledge Graph
- `frontend/src/pages/TimelinePage.tsx` — **NEW**: Timeline visualization
- `frontend/src/App.tsx` — Add /graph/:universeId and /timeline/:universeId routes
- `backend/migrations/014_seed_demo_saga.up.sql` — **NEW**: "Echoes of Eternity" seed data

### Approaches

1. **Full Phase 2 as single SDD cycle** — Build all memory services + frontend core + demo saga in one pass
   - Pros: Faster iteration, single orchestration, all pieces land together
   - Cons: Large scope (~3,000+ lines across 30+ files), harder to review, risk of context overflow
   - Effort: High

2. **Split Phase 2 into sub-phases: 2a (Backend Memory) → 2b (Frontend Core) → 2c (Demo Saga)**
   - Pros: Manageable PRs, each sub-phase independently verifiable, easier to parallelize
   - Cons: More orchestration overhead, dependency chain (frontend depends on WS events from backend)
   - Effort: Medium (recommended)

3. **Minimum-viable Phase 2** — Focus on contradictions + contextual recall + TipTap editor only. Defer timeline, plot holes, KG visualization
   - Pros: Smallest scope, highest demo impact per line of code
   - Cons: Leaves functionality gaps visible in demo
   - Effort: Low

### Recommendation

**Approach 2: Split into sub-phases 2a → 2b → 2c**.

Rationale:
- Phase 2 is the biggest scope of the hackathon — trying to build it all in one SDD cycle risks incomplete delivery
- Backend memory engine (2a) can be built and tested independently with Go tests + mock WS
- Frontend core (2b) depends on the WS protocol and event types defined in 2a
- Demo saga (2c) needs the entity model from 2a and the seed structure
- Each sub-phase produces a shippable increment

### Sub-phase Breakdown

**Phase 2a: Backend Memory Engine**
- RelevanceService (decay algorithm + touch + decay_all)
- ContradictionService (fingerprint dedup, deterministic checks, batched Qwen-Max)
- MemoryService (contextual recall combining graph facts + vector + mentions)
- TimelineService (temporal position validation)
- PlotHoleService (stale arc scanning)
- AnalysisService (sequential per-chapter queue, core/enrichment separation)
- WebSocket hub (auth_init, paragraph_submit → all WS event types)
- GraphRepo expansions (full graph query, neighbor traversal, edge ops)
- New endpoints: contradictions CRUD, timeline, plot holes, graph data
- Config: decay constants, archive threshold, max contradiction candidates
- Tests: RelevanceService TDD, ContradictionService TDD, WS flow

**Phase 2b: Frontend Core**
- TipTap editor (wired with StarterKit, Placeholder, Highlight extensions)
- Paragraph debounce 5s + WS paragraph_submit
- Entity highlighting in editor text (purple=characters, green=places, gold=events)
- MemoryPanel component (contextual recall, contradiction alerts, plot holes)
- WebSocket client (exponential reconnect, auth_init flow)
- memoryStore (Zustand): contextual memories, contradictions, discovered entities, graph updates
- KnowledgeGraphPage (React Flow with colored/labeled nodes)
- TimelinePage (horizontal timeline with events)
- Routes: /graph/:universeId, /timeline/:universeId
- WS status indicator (green/yellow/red)

**Phase 2c: Demo Saga + Polish**
- Generate "Echoes of Eternity" seed content via Qwen
- Migration with 2-3 chapters of content, 20+ entities, planted contradictions
- Entity detail views / character cards
- Knowledge Graph layout tuning
- Video demo prep

### Key Technical Decisions Needed

1. **WebSocket library**: Use `gofiber/contrib/websocket` (Fiber's own, based on gorilla/websocket) vs raw gorilla/websocket. Decision: Use `gofiber/contrib/websocket` for consistency with Fiber.

2. **Analysis pipeline**: Single goroutine per chapter with `context.WithCancel` for enrichment steps. Core transaction (entity extraction + DB writes) is NOT cancelable once started.

3. **Decay constants**: `DECAY_LAMBDA=0.1`, `ARCHIVE_THRESHOLD=0.15`, `REACTIVATION_SCORE=0.8` — align with SRS Section 6.3, make configurable via env vars.

4. **Contradiction max candidates**: Hard limit of `MAX_CONTRADICTION_CANDIDATES=3` per paragraph to control Qwen-Max cost/latency.

5. **Paragraph addressing**: Use TipTap's stable `paragraph_node_id` (not index) — already in the schema.

6. **Fingerprint approach**: SHA-256 of `entity_id + evidence_a_chapter_id + evidence_b_chapter_id` — prevents duplicate contradiction alerts.

### Risks

- **Qwen API key not available**: Entire analysis pipeline depends on Qwen calls. Without API key, Phase 2a cannot be verified end-to-end. Mitigation: Mock QwenService for unit tests; real integration tests only with key.
- **AGE + pgvector stability**: The Docker Compose setup needs testing with full memory operations (AGE graph updates + pgvector embeddings in same transaction). Phase 0+1 confirmed basic coexistence, but heavy concurrent operations may expose issues.
- **Phase 2 scope is ~3,000+ lines**: Sub-phasing mitigates this, but each sub-phase still needs careful PR budgeting. 2a alone is ~1,500 lines across ~15 files.
- **WebSocket complexity**: Client reconnection, message ordering, and race conditions in the analysis queue need careful handling. Mitigation: WS tested in isolation first.
- **Demo saga generation needs Qwen**: Generating coherent "Echoes of Eternity" content with planted contradictions requires multiple Qwen calls. If Qwen is unavailable, fallback to manually crafted seed data.
- **Deadline pressure**: Jul 9 deadline. Phase 2 must be delivered by ~Jul 5-6 to leave time for deploy + video. Each sub-phase must be ~1.5 days max.

### Ready for Proposal
Yes — proceed to SDD Proposal phase for **Phase 2a: Backend Memory Engine**. The orchestrator should present this exploration to the user, confirm the sub-phase approach, then launch `sdd-propose` for Phase 2a.
