# Sprint 0 — Entity Standard & Front Correctness

**Tier:** P0 · **SRS coverage:** TX-1..7, ET-1..5, FE-1..4, GV-1..6, WS-1..3
**Prerequisite:** none — this is the first sprint.
**Why first:** you cannot trust an end-to-end observation while the front cannot display
what the backend produced, and you cannot expect the AI to classify correctly against a
standard that does not exist. Everything later depends on this sprint's two closed
vocabularies (genre tags, entity types).

Work through the tasks in order — each task lists its files, steps, and expected result.

---

## Task 0.1 — Migration `020`: both taxonomies (TX-1..6, ET-2)

**Files:** `backend/migrations/020_taxonomy.up.sql`, `backend/migrations/020_taxonomy.down.sql`

One migration carries both taxonomy changes; they are one conceptual change ("close the
vocabularies") and they roll back together.

**Steps — `020_taxonomy.up.sql`:**

1. Universes:
   ```sql
   ALTER TABLE universes ADD COLUMN genre_tags text[] NOT NULL DEFAULT '{}';
   UPDATE universes SET genre_tags = ARRAY[genre] WHERE genre IS NOT NULL AND genre <> '';
   ALTER TABLE universes DROP COLUMN genre;
   ALTER TABLE universes DROP COLUMN format;
   ```
   Before dropping `format`, copy it down to works (step 2). Genre tag values must be
   normalized to the 20-tag vocabulary in `Docs/SKILLS.md §2`; legacy values that do not
   map (verify against the current `allowedGenres` list in
   `backend/internal/services/universe_service.go:16-40`) map to the closest tag
   (e.g. `sci-fi` → `science-fiction`) in the same UPDATE.
2. Works — format moves here:
   ```sql
   UPDATE works w SET type = u.format FROM universes u
     WHERE w.universe_id = u.id AND u.format IN ('novel','novella','short-story');
   UPDATE works SET type = 'novel' WHERE type IS NULL OR type NOT IN ('novel','novella','short-story');
   ALTER TABLE works ADD CONSTRAINT works_type_check CHECK (type IN ('novel','novella','short-story'));
   ```
3. Entities — the canonical 7-type CHECK (ET-2):
   ```sql
   -- Map any stray existing values first (the prompt never produced 'object', but be safe):
   UPDATE entities SET type = 'world_rule' WHERE type IN ('worldrule', 'rule');
   UPDATE entities SET type = 'place' WHERE type = 'location';
   ALTER TABLE entities ADD CONSTRAINT entities_type_check
     CHECK (type IN ('character','place','object','faction','event','world_rule','plot_arc'));
   ```
   Before writing the mapping UPDATEs, run `SELECT DISTINCT type FROM entities;` against a
   dev DB with real ingested data and cover every value found.

**Steps — `020_taxonomy.down.sql`:** restore `genre` (first tag wins: `genre_tags[1]`),
restore `universes.format` from the most common work type per universe (or `'novel'`),
drop both CHECK constraints, drop `genre_tags`.

**Watch out:** migration `014` seeds the demo universe — check what `type` values it
inserts for entities and what `genre`/`format` it sets, and make sure the 020 data
migration maps them cleanly (the migration runs after 014 in a fresh compose boot).

**Expected result:** `docker compose down -v && docker compose up -d` boots clean; the
`migrations` service applies 020 without error; `\d entities` shows the CHECK constraint.

---

## Task 0.2 — Backend model & service: genre tags + work type (TX-1..4, TX-7 server side)

**Files:** `backend/internal/models/models.go`, `backend/internal/repositories/universe_repo.go`,
`backend/internal/services/universe_service.go`, `backend/internal/services/work_service.go` (or
wherever work create/update validation lives), `backend/internal/handlers/universe.go` (shape only)

**Steps:**

1. `models.Universe`: `Genre string` → `GenreTags []string`; delete `Format`.
   `models.CreateUniverseRequest` likewise.
2. `universe_repo.go`: every query lists columns explicitly (`Create`, `FindByID`,
   `ListByUser`, `Update`, `FindBySessionID`) — replace `genre, format` with `genre_tags`
   in all five, and update the matching `Scan` calls. pgx maps `text[]` ↔ `[]string`
   natively.
3. `universe_service.go`: replace the single-genre validation (`allowedGenres`, lines
   16-40) with validation of **every** tag in `GenreTags` against the 20-tag vocabulary
   from `Docs/SKILLS.md §2`. Define the vocabulary as a package-level `map[string]bool`
   — this is the single backend source of truth for genre tags. Empty `GenreTags` is
   allowed (a universe without tags gets no genre-conditioned behaviour). Remove
   `allowedFormats` from universe validation entirely.
4. Work create/update: validate `Type` against `{novel, novella, short-story}`; reject
   others with `VALIDATION_ERROR`. Default to `novel` when omitted.
5. `DemoService` clone path: verify it copies `genre_tags` (it deep-copies universes —
   check the INSERT it uses).

**Expected result:** `go build ./...` passes; `go test ./...` passes (fix any repo tests
that reference `genre`/`format`); creating a universe with `genre_tags: ["fantasy","romance"]`
succeeds, with `["fantasi"]` returns `VALIDATION_ERROR`.

---

## Task 0.3 — Extraction prompt: add `object`, embed criteria (ET-3)

**Files:** `backend/internal/services/qwen_service.go` (extraction prompt ~lines 283-288 and
the `ExtractionResult` struct ~lines 241-245), `backend/internal/services/entity_service.go`
(wherever the result buckets are iterated)

**Steps:**

1. Add `Objects []ExtractedEntity \`json:"objects"\`` to the extraction result struct.
2. Add the `"objects"` line to the prompt's JSON example: `{"name": "...", "type": "object", ...}`.
3. Embed the identification criteria in the prompt (from SRS §2.2, condensed):
   - character: an agent with its own will that acts — if it decides, it is a character
   - object: a named thing with no will (a sword, a ship). The ship is an object; its pilot is a character
   - place / faction / event / world_rule / plot_arc: one line each
   The criteria are what make classification deterministic instead of vibe-based; this is
   the actual fix for "Excalibur classified as a character".
4. Wire the new bucket through `entity_service.go` wherever `Characters`/`Places`/etc.
   are looped (search for `WorldRules` to find every loop).
5. Grep for any other prompt that enumerates entity types (contradiction, plot-hole,
   consolidation prompts) and align them.

**Expected result:** ingesting a fixture containing "Excalibur, the sword of kings" yields
an entity of type `object`. The DB CHECK from Task 0.1 backstops any prompt drift.

---

## Task 0.4 — Frontend: one canonical entity-type module (ET-4)

**Files:** new `frontend/src/lib/entityTypes.ts`; edit `frontend/src/pages/EntitiesPage.tsx`,
`frontend/src/lib/graphParse.ts`, and wherever `NODE_TYPE_META` lives (imported by
`EntitiesPage.tsx` and `EntityCardPage.tsx` — follow the import)

**Steps:**

1. Create `entityTypes.ts` exporting:
   ```ts
   export const ENTITY_TYPES = ['character','place','object','faction','event','world_rule','plot_arc'] as const
   export type EntityType = typeof ENTITY_TYPES[number]
   export const ENTITY_TYPE_META: Record<EntityType, { label: string; color: string; glyph: string }> = { … }
   ```
   Fold the existing `NODE_TYPE_META` content into `ENTITY_TYPE_META` (keep the palette —
   colors must satisfy `design.md`). One module, one list, one meta map.
2. Delete `TYPE_FILTERS` (`EntitiesPage.tsx:22`) and `ENTITY_TYPES` (`graphParse.ts:24`);
   import from the new module. This kills the `worldrule`/`world_rule` mismatch (FE-2)
   structurally — there is no second copy left to drift.

**Expected result:** `npm run build` passes; grep confirms exactly one file in
`frontend/src` declares the entity-type list.

---

## Task 0.5 — Entity browser: show everything (FE-1..4)

**Files:** `frontend/src/pages/EntitiesPage.tsx`

**Steps:**

1. Remove `.slice(0, 4)` at line 205 — render a chip for `'All'` plus every canonical
   type. If seven chips overflow the rail, wrap them (`flex-wrap`), do not truncate.
2. Total count: `api.listEntities` already returns `pagination` (`entity.go` handler
   returns `total`). Display "N entities" above the list, from `pagination.total` —
   **not** from `entities.length`.
3. Reach every entity: the handler default limit is 50 (`entity.go:34`) and the page
   currently requests `limit: '200'`. Replace with proper pagination: request pages of
   100 with a "Load more" button (or auto-load on scroll), accumulating until
   `entities.length === pagination.total`. Verify the backend actually honors `limit`
   and `page` end to end — this is the prime suspect for "~130 extracted, ~30 visible".
4. Per-type counts on chips (FE-4): derive from a one-shot request or extend the list
   endpoint with a `counts_by_type` field — cheapest correct option wins; a `GROUP BY type`
   count query in the repo is trivial.

**Expected result:** with a universe holding 130 entities, the browser shows "130
entities", all seven filters work, and scrolling/Load-more reaches the last one.

---

## Task 0.6 — Graph: neighborhood navigation (GV-1..6)

**Files:** `frontend/src/stores/graphStore.ts`, `frontend/src/pages/KnowledgeGraphPage.tsx`,
`frontend/src/components/knowledge-graph/GraphCanvas.tsx`; backend
`backend/internal/handlers/entity.go` / `entity_repo.go` (sort support, if absent)

The current behaviour: `graphStore.fetchGraph` calls `api.getGraph(universeId)` (the
**entire** universe) and lays nodes on a circle (`graphStore.ts:54`, `cos/sin * radius`) —
that is the unconnected ring. The fix is to stop fetching the world.

**Steps:**

1. **Focal node default (GV-2):** the backend must answer "highest relevance active
   entity". `EntityFilters` already supports `MinRelevance`; add a `sort=relevance`
   option to the list query (ORDER BY relevance_score DESC) if not present, and call
   `listEntities(universeId, { limit: '1', sort: 'relevance', status: 'active' })`.
2. **Neighborhood fetch (GV-1):** replace the store's full-graph fetch with
   `api.getEntityNeighbors(focalId, universeId)` — the endpoint already exists and takes
   `hops`; use `hops=2`. Keep `parseVertexRaw` for the agtype parsing. If the backend
   neighbors query caps at 1 hop internally, extend the Cypher in `graph_repo.go` to a
   variable-length pattern `-[*1..2]-` — identifiers are not interpolated here, but keep
   everything inside `withAgeTx`/`withAgeConn` as always.
3. **Re-centre on click (GV-3):** node click handler sets the clicked node as the new
   focal and refetches its neighborhood. Keep a small breadcrumb of visited focals so
   the writer can step back.
4. **Search (GV-4):** a search input over `listEntities(universeId, { search })`
   (the repo `Search` filter exists); selecting a result makes it the focal node.
   Alias matching: the repo search currently matches `name` — extend the WHERE to
   `OR aliases` (same pattern as `FindByAlias`).
5. **Layout:** with ≤ a few dozen nodes, a radial layout around the focal node is fine
   and honest — 1-hop ring, 2-hop outer ring. Do not add a layout library; `design.md`
   forbids chart libraries and the node counts no longer justify one.
6. **Relevance emphasis + type filter (GV-5, GV-6):** node size or opacity from
   `relevance_score` (already parsed in `parseVertexRaw`); a type-filter chip row reusing
   `ENTITY_TYPE_META`; archived entities hidden unless a toggle is on (`status` is parsed).

**Expected result:** opening the graph shows one readable neighborhood (focal + 2 hops)
with visible edges; clicking walks the graph; searching "Avasarala" jumps to her
neighborhood. The full-universe ring is gone.

---

## Task 0.7 — Workspace: collapsible panels (WS-1..3)

**Files:** the layout component (`frontend/src/…/UniverseLayout` — follow the router) and
its CSS module; the panels it hosts

**Steps:**

1. Each side panel gets a collapse toggle (a thin edge button with a chevron glyph, per
   `design.md` inline-glyph style). Collapsed state = a narrow rail with re-open affordance.
2. Persist collapsed/width state in `localStorage` keyed by panel id.
3. Resizable (WS-3): a drag handle on the panel edge adjusting a CSS grid column;
   pointer-events based, ~15 lines, no library. Clamp to a min width at which nothing
   truncates (WS-2) — audit current panels at default width and fix any label/icon
   truncation found while you are there.

**Expected result:** a writer can collapse everything and keep only the editor; sizes
survive reload; no panel truncates its own labels at any allowed width.

---

## Definition of done

- [ ] `docker compose down -v && docker compose up -d` boots; migration 020 applies; down.sql restores the prior shape.
- [ ] `go build ./...` and `go test ./...` pass; a repo test asserts an invalid entity type is **rejected by the DB** (ET-2 verification).
- [ ] `npm run build` and `npm run test` pass; a component test asserts **all filter chips render** and the displayed total equals `pagination.total` (FE-1/FE-3 verification).
- [ ] Manual pass: create universe (multi-genre select), create work (format select), ingest the demo fixture, open Entities (all reachable), open Graph (neighborhood, not ring), collapse panels.
- [ ] Exactly one entity-type list exists in the frontend and one genre vocabulary in the backend (grep proof).
