# Final submission checklist

Use this as the final handoff checklist. Do not mark external evidence complete until it is
publicly reachable and reviewed by a human.

## Deadline

**July 20, 2026 — 5:00 PM EDT (2:00 PM PDT).** Submit the Devpost project before the deadline;
saved drafts are not a final submission.

## Qwen Cloud / Devpost final form

- [ ] Confirm every entrant is an eligible adult in a Qwen Cloud-supported, non-restricted
  country/region; appoint the authorized representative for a team submission.
- [ ] Select **Track 1: MemoryAgent** and add every teammate to the Devpost project.
- [ ] Write the English description: **what** Quill does, **who** it helps, and **how** persistent
  memory stores, retrieves, forgets, and budgets story context.
- [ ] Name Qwen Cloud and the Qwen models/services actually used in the submission's **Built With**
  field; do not expose `QWEN_API_KEY`.
- [ ] Provide a public repository URL for judging/testing. Confirm it contains source, assets, setup
  instructions, and the [MIT License](../LICENSE), and that the OSI license is visible in its About
  section.
- [ ] Attach the committed [architecture SVG](assets/quill-architecture.svg), showing frontend,
  Go/Fiber backend, Qwen Cloud, PostgreSQL/pgvector/AGE, and the WebSocket path.
- [ ] **Required by the rules:** add a public repository link to a code file demonstrating Alibaba
  Cloud service/API use for the backend. **Recommended supporting evidence:** add the
  organizer-requested screenshot of the running backend. Docker Compose is not deployment proof.
- [ ] Publish a public, functioning demo video on YouTube, Vimeo, or Youku; keep it under three
  minutes and add its URL. Show the actual application, not mocked runtime data.

## OpenAI Build Week blocker

- [ ] Paste/link the Codex `/feedback` session ID in the final Build Week materials. No ID exists in
  this repository; code and tests cannot generate it.
- [ ] Explain the actual GPT-5.6/Codex contribution and human review using the sprint record:
  writer journey, visible memory proof, accessibility and failure feedback, and submission
  verification. Do not claim a session, model use, or authorship that cannot be evidenced.

## Local judge-proof evidence

- [x] `cd frontend && npm run build` passed on 2026-07-17. The current build emits a separate
  `KnowledgeGraphPage-Dz9yI3_a.js` artifact (585.26 kB; 181.27 kB gzip); the entry bundle is
  272.33 kB (84.93 kB gzip).
- [x] Playwright Sprint 7 proof passed 10/10 on desktop Chromium and Pixel 5: guided real routes,
  failure/retry feedback, keyboard focus/reduced motion, axe WCAG A/AA checks, and lazy graph
  loading.
- [ ] Review the Vite warning for the isolated graph chunk exceeding 500 kB before submission; the
  browser proof confirms lazy loading, but the warning is still present.

## Re-run before submission

```bash
cd frontend
npm run build
npm run test:e2e -- --grep 'Sprint 7'
rg --files dist/assets | rg 'KnowledgeGraphPage-.*\\.js$'
```

The Playwright runner uses Chromium and starts a local preview server. On a new machine, install
the browser first with `npx playwright install chromium`.

## Submission caveats

- Keep API keys, tokens, and private prompt material out of recordings, screenshots, and commits.
- Local browser mocks are confined to Playwright request boundaries; they are not runtime demo
  data and do not establish that a deployed backend is available.
