import { useEffect, useState } from 'react'
import { api } from '../../lib/api'
import type { MemoryStatusEntity } from '../../lib/types'
import RelevanceHistoryChart from '../memory/RelevanceHistoryChart'
import PageStatus from '../shared/PageStatus'
import styles from './RelevanceHistoryTab.module.css'

interface RelevanceHistoryTabProps {
  entityId: string
  universeId: string
}

// Relevance-history tab for the Story Graph right panel — the selected
// entity's own decay/relevance line, sourced from the same memory-status
// endpoint DecayTimeline uses for its multi-entity Memory Lab view.
// `entityId` here is the graph node id, which is the same `entities.id` UUID
// memory-status keys its own `entities[].id` by (both are the Postgres
// entities table's primary key — confirmed via entity_service.go always
// setting the AGE vertex's `entity_id` property to `entity.ID.String()`).
export default function RelevanceHistoryTab({ entityId, universeId }: RelevanceHistoryTabProps) {
  const [entities, setEntities] = useState<MemoryStatusEntity[] | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [retry, setRetry] = useState(0)

  useEffect(() => {
    if (!universeId) {
      setLoading(false)
      setEntities(null)
      return
    }
    let cancelled = false
    setLoading(true)
    setError(null)
    setEntities(null)
    api.getMemoryStatus(universeId)
      .then((res) => {
        if (cancelled) return
        setEntities(res.entities || [])
        setLoading(false)
      })
      .catch(() => {
        if (cancelled) return
        setError('Could not load relevance history. Retry to try again.')
        setLoading(false)
      })
    return () => { cancelled = true }
  }, [universeId, retry])

  if (loading) return <PageStatus loading />
  if (error) return <PageStatus error={error} onRetry={() => setRetry((attempt) => attempt + 1)} />
  if (!entities) return null

  // Chart's own empty-state covers "no entities at all"; folding the
  // "found but zero history points" case into the same empty array keeps
  // one empty-state message instead of a second, near-identical one here.
  const found = entities.find((entity) => entity.id === entityId)
  const selected = found && found.history.length > 0 ? [found] : []

  return (
    <div className={styles.wrap}>
      <RelevanceHistoryChart
        entities={selected}
        compact
        emptyMessage="No relevance history is available for this entity yet."
      />
    </div>
  )
}
