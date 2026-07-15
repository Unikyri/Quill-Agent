# Quill — Product Requirements Document (v2)

**Status:** Draft · supersedes the original PRD
**Track:** Global AI Hackathon with Qwen Cloud — **Track 1: MemoryAgent**
**Companion documents:** [SRS.md](./SRS.md) — technical requirements · [SKILLS.md](./SKILLS.md) — skills catalogue and genre vocabulary

---

## 0. Honest State of the System

This section exists because a PRD that describes the system we *wish* we had is worse than
no PRD at all. Everything below was verified against the code, not assumed.

### 0.1 What is genuinely working

- Document ingestion parses real `.pdf` / `.docx` / `.md` / `.txt` and splits into chapters.
- Entity extraction, embeddings, and Apache AGE graph population run and produce data.
- The memory subsystem (RRF fusion across 5 pipelines, relevance decay, consolidation,
  context budgeting) is implemented and has an evaluation harness (`backend/eval/`).
- The Memory Theater UI renders decay timelines, fusion explanations, and budget bars.

### 0.2 What is broken or unproven

| Issue | Evidence | Priority |
|---|---|---|
| **No entity type standard exists** | Three divergent vocabularies: the extraction prompt (`qwen_service.go:283-288`) produces 6 types with **no `object`**; `EntitiesPage.tsx:22` filters on `object` (never produced) and `worldrule` (misspelling of `world_rule` — matches nothing); `graphParse.ts:24` carries a third copy; `entities.type` has no DB constraint. | **P0** |
| **The entity browser hides most of the data** | `TYPE_FILTERS.slice(0, 4)` (`EntitiesPage.tsx:205`) renders only 4 of 7 filter chips. Observed: ~130 entities extracted, ~30 visible. **Until the front displays what the backend produced, E2E observations are untrustworthy.** | **P0** |
| **The graph is an unconnected ring** | It renders every entity at once; with real data the layout degenerates and shows no relationships. | **P0** |
| **Live paragraph analysis does not work at runtime** | The chain is fully wired (`TipTapEditor.tsx:139-161` → `paragraph_submit` → `hub.go:266` → `AnalysisService.SubmitParagraph` → `broadcastResult`) but produces no observable result in practice. **This is a bug, not missing code — it must be diagnosed, not reimplemented.** | **P1** |
| **Zero test coverage on the live-analysis chain** | CodeGraph reports "no covering tests found" for `SubmitParagraph`, `AnalysisHub`, and `TipTapEditor`. | **P1** |
| **Ingestion takes ~1 hour for a 400-page novel** | ~300 sequential LLM calls at ~10s each. A judge testing the product will close the tab. | **P1** |
| **No end-to-end verification exists** | Unit and integration tests pass. That guarantees nothing about the assembled system. In the project owner's words: *"the tests pass, but that guarantees nothing because they are unit and integration tests, not end-to-end."* | **P1** |
| **Duplicate entities** ("James" vs "James Holden") | Observed in practice. Root cause identified: a write-write race in entity resolution. | **P1** |
| **Ingested chapters render as unformatted plain text** | No paragraph structure, no italics, no scene breaks. | **P3** |
| **Autosave has no retry, no failure indicator, no local recovery** | A 5-second debounce and nothing else. If Quill is the source of truth, this is a data-loss risk. | **P3** |

### 0.3 Suspected causes of the live-analysis failure (to be confirmed by diagnosis)

1. The WebSocket may not be open when `send()` fires; the message is dropped silently with no retry queue.
2. `getParagraphAtCursor` keys off the **cursor position**, not the **edited paragraph** — moving the caret before the 5s debounce elapses submits the wrong text.
3. Nobody has ever verified this path end to end.

---

## 1. Vision

**Quill is a writing IDE with a memory of your world — and of you.**

A novelist writing a 200,000-word saga across two years cannot hold their own world in
their head. What colour were Elara's eyes in chapter 3? Did the magic system already
forbid healing the dead? Is the character who vanished in chapter 12 *deliberately*
written out, or forgotten?

Quill remembers. It reads what the writer writes, extracts the world (characters, places,
objects, rules, events, timelines), holds it in a memory that decays and consolidates like
a human one, and surfaces contradictions and plot holes as they appear.

### 1.1 The gap this version closes

**Quill today remembers the STORY. It does not remember the WRITER.**

Track 1 asks for an agent that *"autonomously accumulates experience, **remembers user
preferences**, and makes **increasingly accurate decisions** across multi-turn,
cross-session interactions."*

Story memory — entities, lore, contradictions — is memory of a *document*. It is not
memory of a *person*. This version introduces **Writer Memory**: a subsystem that learns
how the writer works, what they intend, what they reject, and gets measurably less
annoying over time.

That is the difference between a very good RAG pipeline over a novel and an actual
MemoryAgent.

---

## 2. Scope

### 2.1 In scope

Prose fiction writers working in **novels, novellas, and short stories**.

### 2.2 Explicitly out of scope (and why)

| Dropped | Reason |
|---|---|
| **Screenplays** | A screenplay is not "prose with less narration". It is a rigid typographic format (sluglines, centred character names, dialogue margins, parentheticals, transitions) with its own keyboard semantics. It requires custom TipTap nodes and a separate editor mode — a feature, not a dropdown entry. Deferred to v2. |
| **Graphic novels / manga** | Requires image support, which does not exist. |
| **Essays, articles** | Not world-building. Leftovers from the original PRD's over-broad scope. |
| **AI writing the author's prose** | Deliberate product stance. See §4.2. |

---

## 3. Domain Model

| Concept | Definition | Key properties |
|---|---|---|
| **Universe** | One coherent fictional cosmos (the Cosmere model). All works inside it share a world, its rules, and its register. | **Genre: multi-tag, closed vocabulary.** Cross-genre work belongs in a *different* universe, not a different work. |
| **Work** | One book/story within a universe. | **Format: single-valued, closed vocabulary, validated** — `novel`, `novella`, `short-story`. A universe may hold a novel *and* a novella. |
| **Chapter** | An ordered unit of a work. The unit of editing. | |
| **Entity** | A named element of the fictional world, extracted from the text. **Exactly 7 types, closed vocabulary:** `character`, `place`, `object`, `faction`, `event`, `world_rule`, `plot_arc` — with per-type identification criteria (SRS §2.2). | Has a relevance score that decays over time. `event` entities also feed the timeline — the timeline is a view over them, not a separate store. |
| **Skill** | Read-only Markdown carrying **craft** knowledge (genre conventions, register, period terminology, structural expectations). | Curated by us. **Not user-editable.** Activated at the **universe** level. |
| **Observation** | A verifiable **fact** about how the writer writes. *"4.2 adverbs per 100 words."* | Produced by passive stylometry. **Zero LLM cost.** |
| **Preference** | A **belief** about what the writer **intends**. *"Long sentences are deliberate."* | Only created when an Observation is corroborated by an intent signal. Carries confidence, decays. |

### 3.1 Taxonomy migration (P0 — a dependency, not a nicety)

Current state, verified in `backend/internal/services/universe_service.go:16-40`:

- Genre and Format both live on **Universe**, both **single-valued**.
- `Work.Type` exists but is an **unvalidated free string** — a field with no owner.

Required state:

- **Genre → Universe, multi-tag, closed vocabulary.** A book is horror *and* romance *and*
  sci-fi. Forcing one is wrong. **The 20-tag vocabulary is defined in [SKILLS.md §2](./SKILLS.md).**
- **Format → `Work.Type`, single-valued, closed vocabulary, validated.**
- Formats reduced to `novel`, `novella`, `short-story`.

**Why the vocabulary stays closed:** if the model is allowed to invent tags, three books
later the system holds `dark fantasy`, `grimdark`, `fantasía oscura`, and `low fantasy` as
four distinct axes, and preference conditioning (§5.2) disintegrates — you cannot filter on
an axis that mutates.

**Why this is P0:** Writer Memory conditioning and Skill suggestion both *depend* on genre
tags existing and being multiple. Without this migration, the later tiers cannot be built.

### 3.2 Entity taxonomy (P0 — the standard the AI classifies against)

There is currently **no entity type standard** — the extraction prompt is the only
definition, the frontend carries two divergent copies of it, and the database accepts any
string. The consequence is exactly the failure mode this product exists to prevent:
Excalibur has no category (`object` is not in the prompt's vocabulary), and a misspelled
filter (`worldrule` vs `world_rule`) silently matches nothing.

The fix is one canonical, closed 7-type vocabulary with **per-type identification
criteria** ("if it decides, it is a character; a proper name with no agency is an object"),
enforced in three layers: a `CHECK` constraint in the database, a JSON Schema `enum` in
structured output, and a single shared module in the frontend. Full criteria table and
requirements: **SRS §2.2 (ET)**.

---

## 4. Feature Requirements

### 4.1 Skills — craft knowledge, agent-selected

A **Skill** is a Markdown file: frontmatter (`name`, `description`) plus an instruction body.
The model of reference is the Claude Code / claude.ai skill: *"You will act as an expert in
medieval fantasy, familiar with the terminology, places, and narrative conventions of the
period."*

| Req | Requirement |
|---|---|
| **SK-1** | Skills are **read-only**. The writer selects and activates; they do not edit. A writer is a writer, not an engineer. |
| **SK-2** | Skills are activated at the **universe** level, and may be pre-suggested from the universe's genre tags. Craft does not change between chapters — chapter 3 and chapter 40 of a medieval fantasy are both still medieval fantasy. |
| **SK-3** | **The agent auto-selects** which of the active skills apply to a given passage, by reading each skill's `description` frontmatter. The `description` is not human documentation — it is the agent's selection criterion. |
| **SK-4** | Skill selection is **progressive disclosure**: only the short `description` lines enter the prompt; the model chooses; only the winning skill's full body is then loaded. This is a genuine token optimisation, not a cosmetic one. |
| **SK-5** | The UI shows **which skills fired and why**. This is real information, not decoration. |
| **SK-6** | Per-chapter skill overrides are **out of scope** (speculative; the agent already selects per passage). |
| **SK-7** | Quill ships **15 skills** — 7 editorial roles (developmental editor, line editor, copy editor, proofreader, beta reader, sensitivity reader, literary agent), 7 craft skills, and 1 genre skill carrying 20 reference files. Full catalogue: **[SKILLS.md](./SKILLS.md)**. |
| **SK-8** | There is deliberately **no "Continuity Checker" skill.** Continuity, contradictions, plot holes, and timeline validation are the job of the background memory pipeline. Duplicating that engine as a skill would confuse where the intelligence lives. |

### 4.2 Craft review — on demand, never rewriting

| Req | Requirement |
|---|---|
| **CR-1** | Craft review is **100% on demand**. The writer selects a paragraph, a block, or a whole chapter and asks. If they do not ask, the AI says nothing about craft. |
| **CR-2** | Output is a **margin observation**, never a rewrite of the author's prose. *"Three consecutive adverbs here; it breaks the period register this skill establishes."* |
| **CR-3** | **Quill never writes the novel.** The AI is the editor, not the author. A writer who wants prose generated can ask any chatbot; that is not this product. |
| **CR-4** | No pop-ups. Non-invasive affordances only. An editor that interrupts every paragraph is not a good editor — it is a nuisance. |

**Note on autonomy:** the autonomy Track 1 asks for is already satisfied by the background
memory pipeline (§4.3), which runs unprompted. Requiring the *human* to trigger craft
feedback is a correct product decision, not a concession.

### 4.3 Memory writes — propose, then confirm

| Req | Requirement |
|---|---|
| **ME-1** | **The AI never writes an entity silently. It proposes; the human confirms.** The background pipeline produces *candidates*. |
| **ME-2** | Candidates surface in the editor as a **discreet underline** with an inline affordance. No modal, no pop-up. |
| **ME-3** | Known entities and aliases are highlighted **locally, client-side, with zero LLM calls.** If "James" is already a known entity, highlighting it costs nothing. The AI did its work at extraction time; the editor merely harvests. |
| **ME-4** | Bulk ingestion is gated **by confidence**: high-confidence candidates are auto-accepted; low-confidence candidates go to a **review tray**. Asking a writer to confirm 300 entities one by one loses the user. |
| **ME-5** | The **review tray** is mandatory. It is not an admission of failure — it is the *destination of the low-confidence bucket that confidence-gating creates by design*. Without it, the gate cannot be implemented: you would either discard memory or auto-accept everything. It also feeds the learning signals of §5. |

### 4.4 Editor — Quill is the source of truth

The writer's manuscript **lives in Quill**. They do not draft in Word and paste in.

| Req | Requirement |
|---|---|
| **ED-1** | **Faithful import.** Paragraphs, italics (modern prose italicises thought), scene breaks. Today ingestion pastes plain text. |
| **ED-2** | **Robust autosave.** Retry, visible failure state, local recovery. *Losing a chapter kills the product.* There is no second chance with a writer who lost work. |
| **ED-3** | **Markdown export.** If the writer cannot get their book out, they will not put it in. No serious author hands a manuscript to a database they cannot recover it from. |
| **ED-4** | **DOCX export** is deferred to the roadmap (it is what the publishing world actually uses). **PDF export is rejected** — a read format, not an edit format. |
| **ED-5** | Ingestion progress: one bar, one **live action message** beneath it ("Extracting entities from chapter 12…"), plus an **ETA**. The ETA is hidden until ~10% is processed and is smoothed — an ETA computed from two chunks swings wildly and reads as broken. |

### 4.5 Graph & workspace — navigate, don't dump

| Req | Requirement |
|---|---|
| **GW-1** | **Neighborhood graph.** The graph never renders the whole universe. Default view: the highest-relevance entity plus 2 hops. Click any node to re-centre on it; search jumps to any entity by name or alias. The relevance decay engine *is* the hierarchy — no fixed "characters first" ordering needed. |
| **GW-2** | **Complete entity browser.** All 7 type filters visible, values matching what is stored, total count shown, every entity reachable. |
| **GW-3** | **Collapsible, resizable side panels.** Quill is an IDE for writers; a writer who just wants to write collapses everything and keeps only the editor and live alerts. No panel truncates its own labels. |

Details: SRS §2.3–2.5 (FE, GV, WS).

### 4.6 MCP

| Priority | Server | Rationale |
|---|---|---|
| **Mandatory** | **Quill *as* an MCP server** | The tools already exist and already work (`search_vector_memory`, `query_entity_graph` in `agent_tools.go`). Exposing them turns Quill from "an app you write in" into **the memory substrate of your world, queryable from any client**. A writer in any MCP host can ask *"what colour were Elara's eyes?"* and get an answer from their universe. |
| **P3** | **Notion MCP** (official, `mcp.notion.com/mcp`, OAuth) | Writers keep their **worldbuilding bible** in Notion. Quill can ingest it as established lore. This is a **memory source**, not a convenience integration — it feeds the subsystem the track actually scores. Cost: OAuth per user is real work. |
| **P3** | **Wikipedia MCP** (community) | Research and anachronism checks. No auth. **Requires a fallback** — it is the only external network dependency in the entire system. |
| **Rejected** | **Zotero MCP** | It is for academics (citations, BibTeX, APA). A fantasy novelist has no Zotero library. Including it would be exactly the ornament judges detect. |

**Ecosystem finding, worth stating in the submission:** the MCP ecosystem has servers for
programmers and for academics. It has **nothing for fiction writers** — not even a
dictionary, thesaurus, or etymology server. "The first memory MCP for writers" is a genuine
open-source contribution, and it is the direct answer to the *Problem Value & Impact*
criterion: *"could this grow into an open-source project/community?"*

---

## 5. Writer Memory — the core of this release

### 5.1 Observation is not Preference

This distinction is the load-bearing wall of the whole subsystem.

- **Observation** = *"you use 4.2 adverbs per 100 words."* A **fact**. Cheap (pure text
  analysis, **zero LLM calls**), verifiable, and true whether you like it or not.
- **Preference** = *"you want to keep using them."* A **belief about intent**. Your prose is
  **not evidence** for it.

**Stylometry may only write Observations. Never Preferences.**

An Observation is **promoted** to a Preference only when corroborated by an intent signal.

Without this wall, the system fossilises the writer's vices: it confuses *what he does* with
*what he wants*, congratulates him on his worst habits, and never corrects anything.

### 5.2 The four signals

| # | Signal | Weight | Note |
|---|---|---|---|
| 1 | **Explicit accept/reject** of a craft note | Strongest | Cleanest evidence of intent. Scarce — reviews are on demand. |
| 2 | **Silent dismissal** | **Not a rejection** | Silence is noise, not intent. Treating it as rejection is the fastest way to poison a profile. It is *absence of reinforcement* — **already handled by the existing decay engine.** No new code. An ignored attention line fades on its own without ever inferring anything false. |
| 3 | **Behaviour after a note** | Strong | Flagged three adverbs; two minutes later the paragraph is rewritten without them. He did not tell you — he *showed* you. |
| 4 | **Passive stylometry** | Observations only | Sentence length, adverb density, dialogue/narration ratio, lexical richness, rhythm. **Zero LLM cost.** This is what solves cold start: a writer who uploads three manuscripts has a profile on day one. |

### 5.3 The promotion loop (and the demo money-shot)

1. Stylometry observes: *34-word mean sentence length.* → **Observation.** No conclusion drawn.
2. The agent flags: *"this sentence is very long, consider splitting it."*
3. The writer rejects it. Again. And again.
4. The system **promotes** the Observation: *"long sentences are deliberate"* → **Preference.**
5. **It stops flagging them.**

That is, literally, *increasingly accurate decisions across sessions*.

And the writer can ask: **"Why do you believe that about me?"** The system answers with the
evidence trail — *observed in your manuscripts (34 words/sentence), flagged 5 times,
rejected 5 times, confidence 0.87.* **Auditable memory.** A judge sees that and knows there
is no smoke.

### 5.4 Scope and conditioning

Preferences belong to the **user** and **travel across universes** — a new universe must not
start from zero memory.

But each preference carries an **applicability condition**:

- *"Excellent at describing environments"* → **that is the writer.** A craft capability. It
  travels to the grimdark, the children's story, and the romance alike. It is an asset.
- *"Terse, violent prose, short sentences"* → **that is not the writer. That is the writer
  *in the grimdark*.** In a children's story it is poison.

Recall is **conditioned by the target universe's genre tags**. This is what prevents the
system from dragging slasher memories into a story about cats.

**Known risk (§8):** the model must classify which bucket a preference lands in *at write
time*. Misclassifying "graphic violence" as a universal craft capability would contaminate a
children's book. This is a genuine failure point and is designed for explicitly.

### 5.5 Reuse, not reinvention

The decay engine already exists: `RelevanceService`, `DECAY_LAMBDA`, `ARCHIVE_THRESHOLD`,
`entity_relevance_history`. Today it is pointed at story entities. It will *also* be pointed
at writer preferences. Recall enters as **one more RRF pipeline**, conditioned by genre.

This is not a system bolted on the side. It comes in through the door that is already built.

---

## 6. Qwen Cloud Integration

**Principle: we migrate to use capabilities, not to change the shape of the JSON.**

Qwen Cloud exposes the same models behind four protocols (`openai-chat`,
`openai-responses`, `dashscope`, `anthropic`). Switching protocol *alone* buys nothing —
no accuracy, no token savings. The value is in the capabilities.

### 6.1 Protocol

**Native DashScope HTTP, hand-written client in Go.**

- Verified: DashScope SDKs exist **only for Python and Java**
  (`Docs/QwenDocs/api-reference/preparation/02-install-sdk.md`). The backend is Go 1.22.
  **There is no Go DashScope SDK.**
- Rejected: a Python sidecar running the official SDK — it puts a second runtime and a
  second container in the analysis hot path, to be able to say the word "SDK".
- The native HTTP protocol (`POST /api/v1/services/aigc/text-generation/generation`,
  `input`/`parameters` → `output.choices[]`, `usage.input_tokens`) is spoken directly from Go.

### 6.2 Capabilities adopted

| Capability | Where | Why |
|---|---|---|
| **Context cache (explicit)** — 10% of standard input price | Live analysis, skill reviews | Every analysis call ships the same large prefix (skill body, universe lore, entity context). Today we pay for it in full, every time. |
| **Structured output (JSON Schema)** | Entity extraction | The model cannot return a type outside the enum. Today we ask for JSON and pray. This is part of the cure for "Llanura classified as a character". |
| **Rerank API** | After RRF fusion | RRF fuses **by rank** — it is semantically blind. A reranker reads the query *and* the content. Improvement is measurable with the existing `backend/eval/` harness. |
| **Thinking mode** | Contradictions, plot holes | The hardest reasoning task in the product. |
| **Batch API** — 50% off | **Consolidation only** | See the trap below. |

### 6.3 The Batch API trap

Batch is 50% cheaper and looks perfect for ingestion. **It is not.** Its SLA is *"usually
within 24 hours"*. Ingestion has a **live progress bar** and the writer is watching. The demo
video is **three minutes long**. Batch is only usable where **nobody is waiting**:
consolidation, and corpus-wide stylometry re-analysis.

*(Batch and context cache cannot be combined on the same request.)*

### 6.4 Model tiering

| Model | Limits | Used for |
|---|---|---|
| **qwen-turbo** | 600 RPM / **5M TPM** | **Entity extraction** — a simple task (read text, emit names and types) at massive volume. 5× the TPM headroom, and cheaper. |
| **qwen-max** | 600 RPM / **1M TPM** | Contradictions, plot holes, skill reviews — few calls, hard reasoning. |

---

## 7. Ingestion Performance

### 7.1 Root cause of the one-hour ingest

~300 sequential LLM calls × ~10s each. **This is not a slowness problem. It is a fan-out
problem.** (For contrast: when a chatbot "analyses" a 500-page document instantly, it is not
analysing it — it is *reading* it, in one call. It is not the same operation.)

### 7.2 MAP / REDUCE

The pipeline splits into two halves with different natures:

- **MAP — parallel.** Extracting *mentions* from a chunk is **stateless**: "this paragraph
  mentions James" depends on nothing else. Embarrassingly parallel. **95% of the time is here.**
- **REDUCE — sequential.** Deciding whether "James" is new or is the "James Holden" created
  three paragraphs ago **depends on what is already in the database**. It is **stateful** and
  must stay single-threaded.

**This is the root cause of the observed duplicates.** Two goroutines resolving "James" and
"James Holden" concurrently cannot see each other's writes, so both create. It is a
**write-write race**, and the fix is not "parallelise less" — it is *parallelise the half
that can be*.

### 7.3 Throughput governance

- Concurrency governed by a **token bucket on TPM**, not a goroutine count. Grouping
  paragraphs per call lowers RPM but ships **the same tokens** — TPM is the real ceiling.
- **Reserve ~30% of quota for interactive analysis.** Rate limits are **account-level and
  shared**: an ingestion job must never starve the writer typing in the editor.
- **Smooth ramp-up** + **exponential backoff with jitter**. Qwen has a distinct
  `429-Throttling.BurstRate` error: raising your request rate *too fast* throttles you even
  while you are *under* the limit.
- **Backup-model fallback** on 429.

### 7.4 Target

> **A ~400-page novel ingests in ≤ 5 minutes (p95). Stretch goal: ≤ 3 minutes.**
> (~50 pages ≈ 35–40 seconds — a judge uploading a single chapter sees it finish in under a minute.)

---

## 8. Risks

| Risk | Impact | Mitigation |
|---|---|---|
| Live analysis is broken and the cause is unknown | **Fatal.** The core loop is the product. | P0 diagnosis before any new feature. E2E verification. |
| Preference misclassification (universal vs genre-bound) | Contaminates unrelated universes | Explicit classification step; the review/evidence trail makes it visible and correctable. |
| Ingestion starves interactive analysis (shared account quota) | Writer's editor goes silent during ingest | Quota reservation for interactive traffic. |
| Wikipedia MCP unavailable during the demo | Broken feature on camera | Fallback: review proceeds without it. It is the only external dependency. |
| Notion OAuth complexity | Time sink | It is P3. Sacrificed first. |
| Scope | Everything above is a lot | Strict priority tiers (§9). |

---

## 9. Priorities

| Tier | Items | Rationale |
|---|---|---|
| **P0 — the standard & the ruler** | **Entity taxonomy** (7 types, criteria, 3-layer enforcement) · **Taxonomy migration** (genre tags / work format) · Entity browser fixes · Neighborhood graph · Collapsible workspace | You cannot trust an end-to-end observation while the front cannot display what the backend produced — and you cannot expect the AI to classify correctly against a standard that does not exist. Fix the ruler before measuring. Both taxonomies are dependencies of everything below. |
| **P1 — the floor** | Fix live analysis · Parallelise ingestion (<5 min) · Real E2E verification | If the core is broken, everything else is decoration — and a judge sees it in 30 seconds. E2E runs *after* P0 so its observations mean something. |
| **P2 — the ceiling** | **Writer Memory** · Native DashScope client + context cache + structured output + rerank · Skills with agent auto-selection | This is where 60% of the score lives (Technical Depth 30% + Innovation 30%). Writer Memory is the only thing that makes this a MemoryAgent rather than a RAG pipeline over a novel. |
| **P3 — product credibility & ornaments** | Editor as source of truth (highlighting, entity tray, autosave, Markdown export) · Quill as MCP server · Notion MCP · Wikipedia MCP · Thinking mode · Batch for consolidation | Problem Value & Impact, 25%. Sacrificed first if needed. |

---

## 10. Success Criteria

1. A judge uploads a chapter and sees ingestion complete in **under a minute**, with a
   progress bar, a live action message, and an honest ETA.
2. A judge writes a paragraph in the editor and sees entities highlighted, candidates
   proposed, and a contradiction surfaced **without asking for it**.
3. A judge selects a paragraph, requests a craft review, sees **which skills fired and why**,
   rejects the note — and sees the system **learn** and stop repeating it.
4. A judge asks the system **"why do you believe that about me?"** and gets an evidence trail.
5. A judge queries the universe's memory **from outside Quill**, over MCP.
