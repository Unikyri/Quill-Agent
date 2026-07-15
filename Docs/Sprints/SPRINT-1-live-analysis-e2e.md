# Sprint 1 — Live Analysis Diagnosis & End-to-End Verification

**Tier:** P1 · **SRS coverage:** LA-1..6, EE-1..3
**Prerequisite:** Sprint 0 complete. The entity browser and graph must faithfully display
backend output before any observation made in this sprint can be trusted.

**The rule of this sprint (LA-1): diagnosis precedes any rewrite.** The chain is fully
wired: `TipTapEditor.tsx:139-161` (5s debounce → `paragraph_submit`) →
`wsStore` → `ws.Hub.handleParagraphSubmit` (`hub.go:266`) →
`AnalysisService.SubmitParagraph` → per-work sequential queue → `broadcastResult`
(`analysis_service.go:554`). It produces no observable result in practice.
Reimplementing a wired-but-broken path is forbidden — find the break.

---

## Task 1.1 — Reproduce and instrument (LA-1)

**Files:** none changed yet — this task produces a diagnosis note, not code.

**Steps:**

1. Boot the stack, open the editor, open browser devtools (Network → WS frames) and
   backend logs side by side.
2. Type a paragraph, wait out the 5s debounce, and answer these questions **in order** —
   each is a link in the chain, stop at the first broken one:
   - Does a `paragraph_submit` frame leave the browser? (If not: front-side bug.)
   - Is the socket `open` at send time? (`wsStore` — a send on a CONNECTING socket is
     silently dropped; this is prime suspect #1.)
   - Does `hub.go:266` log receipt? (Add a temporary log line if none exists.)
   - Does `SubmitParagraph` enqueue, and does the per-work goroutine pick it up?
   - Do the fan-out services (entity/contradiction/relevance/timeline/plot-hole) return
     or error? (Watch for silent error swallowing in the goroutine.)
   - Does `broadcastResult` fire, and does `SendToUser` find the user's connection?
     (Prime suspect #2: user-connection registry mismatch — the result is computed and
     sent to nobody.)
   - Does a WS frame arrive back in devtools, and does `wsStore`'s message switch have a
     case for its `type`? (Prime suspect #3: message arrives, frontend drops it on the
     switch floor.)
3. Also test the cursor-move case: type in paragraph A, immediately click into paragraph
   B, wait for the debounce. `getParagraphAtCursor` (`TipTapEditor.tsx:189`) reads the
   caret at fire time — if paragraph B gets submitted, LA-3 is confirmed as a real bug.
4. Write the diagnosis in the sprint log: which link broke, with the evidence (frame
   captures / log lines). Only then move to 1.2.

**Expected result:** a written root cause. Not a guess — the specific broken link with
evidence.

---

## Task 1.2 — Fix the confirmed break(s) (LA-2, LA-3)

**Files:** determined by 1.1 — most likely `frontend/src/stores/wsStore.ts`,
`frontend/src/components/editor/TipTapEditor.tsx`, possibly `backend/internal/ws/hub.go`

Apply only fixes justified by the diagnosis. The two designed-in fixes from the SRS apply
regardless of which link broke, because both are real defects visible in the code:

1. **Outbound queue (LA-2):** `wsStore` must queue messages while the socket is not
   `open` and flush on connect. A `paragraph_submit` must never be silently discarded.
   Keep it minimal: an array, flushed in `onopen`, capped (e.g. 50) to avoid unbounded
   growth on a dead connection.
2. **Edited paragraph, not cursor paragraph (LA-3):** capture the target paragraph node
   at **edit time** (TipTap transaction) rather than resolving the caret at debounce
   expiry. The debounce timer carries the captured node/position with it.

**Expected result:** typing a paragraph reliably produces an `analysis_result` frame in
devtools, including when the caret moves after typing and when the page was just loaded
(socket still connecting).

---

## Task 1.3 — Lifecycle visibility (LA-4, LA-5)

**Files:** `backend/internal/services/analysis_service.go`, `backend/internal/ws/protocol.go`,
`frontend/src/stores/wsStore.ts`, editor UI component for the status affordance

**Steps:**

1. **Server terminal message (LA-5):** every submitted paragraph ends in exactly one
   terminal WS message — `analysis_result` on success or an `analysis_failed` (new type in
   `protocol.go`) carrying the paragraph ref and a reason. Audit every error path in the
   analysis goroutine: today, an error mid-fanout likely ends in silence. Remember both
   sides: add the Go constant **and** the `wsStore` switch case together.
2. **Client lifecycle state (LA-4):** the editor shows a discreet per-submission state:
   submitted → analyzing → done / failed. Per `design.md`: a small glyph in the margin or
   status bar, not a spinner overlay. A silent failure must be visually distinct from a
   slow success.

**Expected result:** kill the backend mid-analysis; the editor shows *failed*, not an
eternal *analyzing*.

---

## Task 1.4 — Test coverage on the chain (LA-6)

**Files:** `backend/internal/services/analysis_service_test.go`,
`backend/internal/ws/hub_test.go`, `frontend/src/components/editor/TipTapEditor.test.tsx`
(or colocated `__tests__/`)

CodeGraph reports "no covering tests found" for all three links. Minimum set:

1. `SubmitParagraph`: enqueues; same-work paragraphs process in order; an erroring
   analyzer still yields a terminal result (LA-5 regression guard).
2. `handleParagraphSubmit`: a malformed payload does not panic; a valid one reaches the
   submitter (use a stub submitter — the two-phase `SetSubmitter` init makes this easy).
3. `TipTapEditor` submit path: editing paragraph A then moving the caret to B submits A
   (LA-3 regression guard); messages sent while the socket is closed are queued (LA-2).

**Expected result:** `go test ./internal/services/ ./internal/ws/` and `npx vitest run`
cover the three links; each test fails if its sprint fix is reverted.

---

## Task 1.5 — End-to-end suite (EE-1..3)

**Files:** new `e2e/` at repo root (Playwright — drives the real browser editor and can
assert on WS frames; vitest cannot), `docker-compose` as the environment

**Steps:**

1. **EE-1 (live loop):** boot stack → register/login → create universe (multi-tag genre —
   exercises Sprint 0) → create work (`novel`) → create chapter → type a paragraph
   containing an obvious new entity → assert an `analysis_result` arrives over WS with at
   least one extracted entity, and the entity appears in the browser.
2. **EE-2 (ingestion loop):** upload a small fixture (`.md`, ~10 pages, seeded with known
   entities of several types) → assert `ingestion_progress` reaches a terminal state →
   assert chapters, entities (correct types — exercises the Sprint 0 taxonomy), and graph
   edges exist.
3. **EE-3 (one command):** `make e2e` (or an npm script) boots compose, waits for health,
   runs the suite, tears down. Document it in the README. This command is part of the
   definition of done for **every subsequent sprint**.
4. Needs a real `QWEN_API_KEY`; the suite is not in unit CI. It is the pre-merge gate you
   run locally.

**Expected result:** one command proves the assembled system works. Unit and integration
tests passing is not, and has never been, that proof.

---

## Definition of done

- [ ] Root cause written down with evidence (frames/logs), before any fix was applied.
- [ ] Typing reliably produces a visible result; caret movement and slow-connect no longer lose submissions.
- [ ] Every submission ends in a terminal WS message; the editor renders all lifecycle states, including *failed*.
- [ ] The three chain links have failing-on-revert tests.
- [ ] `make e2e` passes: live loop + ingestion loop, one command.
