// Shared types used by both api.ts and memory-related components. Extracted
// from a former duplicate in ContextPanel.tsx + api.ts (S3 Memory Theater
// design, obs #265) — ContextPanel keeps its own inline copy per spec
// non-goal, this file exists so no THIRD copy appears.

export type Lifecycle = 'active' | 'decaying' | 'archived' | 'consolidated' | 'reactivated'

export interface MemoryHistoryPoint {
  score: number
  recorded_at: string
}

export interface MemoryStatusEntity {
  id: string
  name: string
  type: string
  relevance_score: number
  status: string
  consolidated: boolean
  lifecycle: Lifecycle
  history: MemoryHistoryPoint[]
}
