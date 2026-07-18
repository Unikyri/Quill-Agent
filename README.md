# Quill — a memory-aware writing studio

Long-form fiction breaks when a writer has to remember every promise made hundreds of pages ago. **Quill turns a manuscript into a queryable story memory:** it extracts entities and relationships, recalls relevant lore while the writer works, and makes decay, retrieval, and context-budget decisions inspectable.

It is prepared for **Qwen Cloud Global AI Hackathon Series — Track 1: MemoryAgent** and **OpenAI Build Week — Apps for Your Life**. This README distinguishes implemented behavior from submission evidence that still needs to be published.

## Judge quick path

1. Start the compose stack and open [http://localhost:3001](http://localhost:3001).
2. Sign in, then use the guided demo or create a universe and work.
3. Draft or ingest text; inspect live analysis, the relationship graph, and Memory.
4. In Memory, run a recall and inspect its contributing pipelines, relevance lifecycle, and budget outcome.

## What makes it a MemoryAgent

| MemoryAgent requirement | Quill implementation |
| --- | --- |
| Persistent memory | Manuscript paragraphs, entities, and relationships persist in PostgreSQL 16 with pgvector and an Apache AGE graph. |
| Retrieval beyond a flat log | Vector, graph, recency, keyword, and consolidated-memory pipelines are fused with Reciprocal Rank Fusion (RRF). |
| Forgetting | Relevance decays over time; low-scoring entities can archive and reactivate when referenced again. History is exposed by the memory-status API. |
| Limited-context recall | A context-budget manager selects recall results that fit the configured window after reserving response tokens. |
| Inspectable behavior | `recall/explain`, `memory-status`, and the Memory screen expose contributing pipelines and lifecycle data. |

Quill also checks prose as it is submitted: a WebSocket hub sends paragraphs to a sequential per-work analysis queue, which fans out to entity, contradiction, relevance, timeline, and plot-hole work. The sequence prevents paragraphs from the same work being analyzed out of order.

## Qwen Cloud usage

The backend uses the Qwen Cloud OpenAI-compatible endpoint by default. Its committed configuration selects:

- `qwen-turbo` for extraction;
- `qwen-max` for reasoning and craft work;
- `text-embedding-v4` with 1024 dimensions for pgvector;
- `qwen3-rerank` when reranking is configured.

The Qwen service provides chat completions, embeddings, tool calling, throttling/retry behavior, and an optional 429 fallback. Its agent loop can call vector-memory search and entity-graph lookup before returning an analysis result. See [`.env.example`](.env.example) and [`backend/internal/services/qwen_service.go`](backend/internal/services/qwen_service.go).

## Architecture

![Quill architecture](Docs/assets/quill-architecture.svg)

The current SPA is [`frontend/`](frontend/): a React 18/Vite application. Docker Compose serves it on port **3001**; the frontend container listens on port 3000 and calls the Go/Fiber backend on **8080**. Compose also starts PostgreSQL on **5432**. The `web/` directory is an experiment scheduled for removal and is not part of the active application or demo path.

## Run locally

Prerequisites: Docker Compose and a Qwen Cloud API key.

```bash
cp .env.example .env
# Set QWEN_API_KEY in .env; do not commit it.
docker compose up -d
```

Open [http://localhost:3001](http://localhost:3001). The API is available at `http://localhost:8080`; use the frontend application rather than calling it directly for the complete editor and WebSocket flow.

For frontend-only development:

```bash
cd frontend
npm install
npm run dev
```

This serves the current SPA on `:3000` and proxies `/api` (including WebSocket) to `localhost:8080`. For backend-only work, run `cd backend && go run cmd/server/main.go` with a reachable `DATABASE_URL`.

## Verification

```bash
# Go packages; database-backed integration tests need TEST_DATABASE_URL.
cd backend && go test ./...

# Current SPA typecheck, tests, and production build.
cd frontend && npm run test
cd frontend && npm run build

# Optional model-backed memory evaluation: requires PostgreSQL + AGE + QWEN_API_KEY.
cd backend
TEST_DATABASE_URL=postgres://quill:quill_dev_password@localhost:5432/quill?sslmode=disable \
  QWEN_API_KEY=your_key go test ./eval/ -run TestMemoryEval -v
```

The committed [evaluation report](Docs/eval/results.md) records a small, dated evaluation corpus and degraded-mode latency measurements. Treat it as a reproducible project artifact—not a production benchmark or a universal quality claim.

## Three-minute demo route

1. Clone or reset the demo universe through the in-app guided demo.
2. Open **Write**, draft or ingest a passage, and observe analysis feedback.
3. Open **Explore** to inspect entities and their relationships.
4. Open **Memory**, ask a lore question, then show the recall explanation, decay history, and budgeted results.
5. Open **Review** to show a surfaced continuity issue or candidate and the author decision.

## Build Week: Codex and GPT-5.6

The repository documents the product work and verification plans from its sprint sequence, rather than inventing session records. [Sprint 7](Docs/Sprints/SPRINT-7.md) specifies the judge-facing SPA work: a guided writer journey, visible memory proof, accessibility checks, real-data failure handling, and evidence requirements. [The sprint index](Docs/Sprints/README.md) places that work after analysis, ingestion, DashScope, writer-memory, skills, and editor/MCP milestones.

**Build Week blocker:** paste the Codex `/feedback` session ID in the submission and link it from the final project materials. No ID is claimed in this repository. The accompanying explanation must state what GPT-5.6/Codex actually changed and what the team reviewed: the sprint record covers the writer journey, memory visibility, accessibility/failure feedback, and submission verification—not a fabricated blanket claim of authorship.

## Submission evidence checklist

### Qwen Cloud / Devpost — Track 1: MemoryAgent

- **Deadline:** July 20, 2026, **5:00 PM EDT** (**2:00 PM PDT**). Complete the Devpost form before this deadline.
- [x] Public-source license is present: [MIT](LICENSE).
- [x] Architecture diagram is committed: [SVG](Docs/assets/quill-architecture.svg).
- [x] Qwen model/API configuration is visible without secrets: [`.env.example`](.env.example).
- [x] Memory storage, retrieval, decay, and budget behavior are implemented and documented above.
- [ ] **Add a direct, public OSI-licensed repository link** and confirm the license is visible in the repository About section.
- [ ] **Required by the rules:** add a public repository link to a code file proving Alibaba Cloud service/API use for the backend. **Recommended supporting evidence:** add the organizer-requested screenshot of the running backend. Docker Compose is not deployment proof.
- [ ] **Add a public video (≤3 minutes)** on YouTube, Vimeo, or Youku that shows the working flow above.
- [ ] Upload/select the architecture diagram, identify **Track 1: MemoryAgent**, add every teammate, and confirm each entrant is eligible in their country/region.
- [ ] Write the English submission description: **what** Quill does, **who** it helps, and **how** its Qwen-backed memory works; name the Qwen models/services in **Built With**.

### OpenAI Build Week — Apps for Your Life

- [x] The implemented product is framed as a personal creative-writing companion.
- [x] Sprint documents describe the judge-facing product work and verification intent.
- [ ] Add the public demo video and its URL.
- [ ] Add dated Codex/GPT-5.6 contribution evidence, human decisions, and the required Codex session ID to the submission.
- [ ] Verify the final Devpost form against the official rules before submitting.

The detailed working checklist is in [Docs/SUBMISSION-CHECKLIST.md](Docs/SUBMISSION-CHECKLIST.md); the local Qwen track snapshot is [Docs/memoryagent-track-rules.md](Docs/memoryagent-track-rules.md). External rules and the submission form take precedence if they differ.

## Repository map

```text
frontend/       Current judge-facing React/Vite SPA
backend/        Go/Fiber API, services, migrations, and Qwen integration
web/            Experimental UI, not part of the active demo path
Docs/           Product, sprint, evaluation, and submission materials
```

## License

This project is licensed under the [MIT License](LICENSE).
