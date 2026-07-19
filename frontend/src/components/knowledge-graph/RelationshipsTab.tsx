import { useEffect, useState } from 'react'
import { api } from '../../lib/api'
import { ENTITY_TYPE_META } from '../../lib/entityTypes'
import { adaptEntityNeighborhood, type StoryGraphNeighborhood } from '../../lib/graphElements'
import PageStatus from '../shared/PageStatus'
import styles from './RelationshipsTab.module.css'

interface RelationshipsTabProps {
  entityId: string
  universeId: string
}

// Relationships tab for the Story Graph right panel — the selected entity's
// direct relationships (same edge data drawn on the canvas), surfaced as an
// accessible, readable list (type / name / relation). Absorbs the a11y
// relationship list that previously lived in KnowledgeGraphPage's
// accessibleSummary, now scoped to one entity instead of the whole map.
export default function RelationshipsTab({ entityId, universeId }: RelationshipsTabProps) {
  const [neighborhood, setNeighborhood] = useState<StoryGraphNeighborhood | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [retry, setRetry] = useState(0)

  useEffect(() => {
    if (!entityId || !universeId) {
      setLoading(false)
      setNeighborhood(null)
      return
    }
    let cancelled = false
    setLoading(true)
    setError(null)
    setNeighborhood(null)
    api.getEntityNeighbors(entityId, universeId, 1)
      .then((res) => {
        if (cancelled) return
        setNeighborhood(adaptEntityNeighborhood(res))
        setLoading(false)
      })
      .catch(() => {
        if (cancelled) return
        setError('Could not load relationships. Retry to try again.')
        setLoading(false)
      })
    return () => { cancelled = true }
  }, [entityId, universeId, retry])

  if (loading) return <PageStatus loading />
  if (error) return <PageStatus error={error} onRetry={() => setRetry((attempt) => attempt + 1)} />
  if (!neighborhood) return null

  const directEdges = neighborhood.edges.filter((edge) => edge.source === entityId || edge.target === entityId)

  if (directEdges.length === 0) {
    return <p className={styles.empty}>No relationships are recorded for this entity yet.</p>
  }

  return (
    <ul className={styles.list}>
      {directEdges.map((edge) => {
        const otherId = edge.source === entityId ? edge.target : edge.source
        const other = neighborhood.nodes.find((node) => node.id === otherId)
        const meta = ENTITY_TYPE_META[other?.type as keyof typeof ENTITY_TYPE_META] || ENTITY_TYPE_META.character
        return (
          <li key={edge.id} className={styles.item}>
            <span className={styles.type} style={{ color: meta.color }}>{meta.label.toUpperCase()}</span>
            <span className={styles.name}>{other?.data.label ?? otherId.slice(0, 8)}</span>
            <span className={styles.relation}>{edge.relationshipType.replace(/[_-]+/g, ' ')}</span>
          </li>
        )
      })}
    </ul>
  )
}
