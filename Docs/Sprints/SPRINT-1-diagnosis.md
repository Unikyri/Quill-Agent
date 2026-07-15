# Sprint 1 — Live Analysis Diagnosis

**Status:** diagnosis completed before implementation
**Date:** 2026-07-15

## Evidence

The normal composition root wires the live-analysis chain correctly: `main.go` creates
the WebSocket hub, creates `AnalysisService`, and calls `hub.SetSubmitter(analysisSvc)`.
The editor emits `paragraph_submit` through `wsStore`, `Hub.handleParagraphSubmit`
delegates to the submitter, and `AnalysisService` owns the per-work queue. A nil
submitter is therefore not the normal failure mode.

The local runtime was initially not reproducible: Docker showed backend and frontend
without PostgreSQL, and the backend logged thirty failed database connection attempts.
This is an environment-health observation, not evidence of a browser-to-server break.

## Confirmed defects

1. **Lost outbound submissions (LA-2).** `wsStore.send` returns when the socket is not
   `open`; `onopen` sends only `auth_init`. A paragraph edited during `connecting` or
   reconnecting is silently discarded.
2. **Wrong paragraph after cursor movement (LA-3).** `TipTapEditor` calls
   `getParagraphAtCursor` when the five-second debounce expires. That helper resolves
   the current selection, so an edit in paragraph A followed by moving the caret to B
   submits B.
3. **No per-submission terminal contract (LA-4/LA-5).** The protocol has no
   `analysis_failed` message or submission reference. Worker errors and duplicate
   submissions are logged/continued without a terminal message. The frontend stores
   success results but has no lifecycle state or UI consumer for an individual
   submission.

## Consequences for the implementation

The fix must add a bounded outbound queue, capture the edited paragraph at transaction
time, carry a client-generated `submission_id` across progress/result/failure messages,
and guarantee one terminal message for every accepted submission. Tests and Playwright
coverage must exercise these contracts.

## Implementation evidence

- `ParagraphSubmitPayload` now carries `submission_id` and `paragraph_ref`; the same
  correlation fields travel through progress, `analysis_result`, and the terminal
  `analysis_failed` payload.
- The client preserves up to 50 outbound messages while a socket is connecting or
  reconnecting. An overflow becomes a visible failed submission instead of a silent
  drop.
- The editor captures the edited block during `onUpdate`; its five-second timer never
  re-resolves the later cursor position.
- The per-work worker emits one terminal result for each accepted job. A worker error
  emits `analysis_failed`; duplicate work returns an empty, correlated result so the
  client does not remain analyzing.
- Focused regression checks: `go test ./internal/services ./internal/ws`, `vitest` for
  `wsStore` and `TipTapEditor`, plus the Playwright suite listed by
  `npx playwright test --list --config=../e2e/playwright.config.ts`.
