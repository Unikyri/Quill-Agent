# Phase 6: Bug Fixes & UI Completion Implementation Plan

This plan addresses the critical bugs and omissions discovered during the QA testing phase. The focus is on resolving the backend database deadlock (which cascaded into breaking the AI Memory system) and completing the frontend navigation structure.

## 1. Backend Fixes (Database Deadlock & WebSocket)

### 1.1 Resolving the `conn busy` Deadlock in `demo_service.go`
**The Problem:**
During `CloneUniverse`, a main transaction (`tx`) is started to clone the universe, works, chapters, entities, etc. However, when cloning the AGE Graph, the `cloneGraph` function calls `s.graphRepo.CreateNode(ctx, ...)` and `s.graphRepo.CreateEdge(ctx, ...)`. These repository methods are currently executing queries against the global connection pool (`pool.Exec`) rather than participating in the active transaction (`tx`). This causes a transaction deadlock: the `tx` holds locks on rows, while the pool is exhausted waiting for the `tx` to finish, resulting in `insert work: conn busy`.

**The Solution (Unit-of-Work Pattern):**
- **Refactor `GraphRepo`:** Modify `CreateGraph`, `CreateNode`, and `CreateEdge` to accept a transaction interface (e.g., `pgx.Tx`) or modify them to use the active transaction if injected into the context.
- **Update `demo_service.go`:** 
  - Change `s.graphRepo.CreateNode(ctx, ...)` to `s.graphRepo.CreateNodeTx(ctx, tx, ...)` (or equivalent implementation).
  - Change `s.graphRepo.CreateEdge(ctx, ...)` to `s.graphRepo.CreateEdgeTx(ctx, tx, ...)`.
- This ensures all graph creation queries participate in the same atomic transaction, releasing the lock correctly upon `tx.Commit(ctx)`.

### 1.2 Restoring the AI Memory Pipeline (WebSocket)
**The Problem:**
The WebSocket handlers (`Chapter Editor` and `Context Panel`) remained in a disconnected/unresponsive state (`🟡`), resulting in the AI not extracting entities or detecting contradictions.
**The Solution:**
- Fixing the `demo_service` deadlock (above) will likely free up the connection pool, allowing the AnalysisQueue's worker goroutines to successfully query the database.
- **Verify WebSocket Reconnection:** Ensure the frontend `ws/client.ts` implements exponential backoff reconnection. Check if the WebSocket URL is correctly prefixed with `ws://` in production vs development environments.
- **Verify Error Handling in Pipeline:** In `backend/internal/services/analysis.go` (or `AnalysisQueue`), ensure that if a database query fails, the context is properly cancelled and a `ws_send("analysis_error", ...)` message is emitted so the frontend doesn't hang indefinitely waiting for insights.

## 2. Frontend Fixes (Missing UI Creation Flows)

### 2.1 Add "New Universe" Button
- **File:** `frontend/src/pages/UniversesPage.tsx` (or the equivalent Dashboard/Home list component).
- **Action:** Add a "Create Universe" button.
- **Flow:** Clicking opens a Modal with a form (Name, Genre, Format). Submitting makes a `POST /api/v1/universes` call and redirects to the new Universe layout `navigate('/universe/:id')`.

### 2.2 Add "New Work" Button
- **File:** `frontend/src/pages/UniverseWorksTab.tsx`
- **Action:** Add a "Create Work" button in the header of the tab.
- **Flow:** Clicking opens a Modal with a form (Title, Type, Synopsis). Submitting makes a `POST /api/v1/universes/:id/works` call and updates the list.

### 2.3 Add "New Chapter" Button
- **File:** `frontend/src/pages/WorkPage.tsx`
- **Action:** Add a "New Chapter" button at the bottom or top of the chapters list.
- **Flow:** Clicking opens a prompt for Chapter Title. Submitting makes a `POST /api/v1/works/:id/chapters` call and redirects to the editor `navigate('/editor/:chapter_id')`.

### 2.4 Add "Back to Work" Navigation
- **File:** `frontend/src/pages/ChapterEditor.tsx` (or `EditorLayout.tsx`)
- **Action:** Add a "← Back to Work" link/button in the top left header next to the Chapter Title.
- **Flow:** Uses `navigate(-1)` or explicitly navigates back to `/universe/:id/works/:id`.

## Summary
By injecting the active `pgx.Tx` into the graph creation methods, we will eliminate the `conn busy` deadlock, which in turn will unblock the entire WebSocket AI pipeline. Simultaneously, completing the 4 missing CRUD buttons in the React UI will make the app fully usable beyond the pre-loaded demo.
