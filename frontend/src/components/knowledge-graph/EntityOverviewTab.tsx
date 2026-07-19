import { useEffect, useState } from 'react'
import { api } from '../../lib/api'
import { ENTITY_TYPE_META } from '../../lib/entityTypes'
import PageStatus from '../shared/PageStatus'
import styles from './EntityOverviewTab.module.css'

interface Entity {
  id: string; universe_id: string; type: string; name: string
  aliases?: string[]; description?: string; properties?: Record<string, unknown>
  status: string; relevance_score: number
}

interface EntityOverviewTabProps {
  entityId: string
}

// Overview tab for the Story Graph right panel — the presentational facts
// about a selected entity (type, status, confidence, aliases, description).
// Extracted from EntitiesPage's EntityDetail; drops the non-functional demo
// chrome (portrait placeholder, add-property button) and the relations list
// (that lives in the separate Relationships tab).
export default function EntityOverviewTab({ entityId }: EntityOverviewTabProps) {
  const [entity, setEntity] = useState<Entity | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [retry, setRetry] = useState(0)

  useEffect(() => {
    if (!entityId) {
      setLoading(false)
      setEntity(null)
      return
    }
    let cancelled = false
    setLoading(true)
    setError(null)
    setEntity(null)
    api.getEntity(entityId)
      .then((res: { entity: Entity }) => {
        if (cancelled) return
        setEntity(res.entity)
        setLoading(false)
      })
      .catch(() => {
        if (cancelled) return
        setError('Could not load this entity. Retry to try again.')
        setLoading(false)
      })
    return () => { cancelled = true }
  }, [entityId, retry])

  if (loading) return <PageStatus loading />
  if (error) return <PageStatus error={error} onRetry={() => setRetry((attempt) => attempt + 1)} />
  if (!entity) return null

  const meta = ENTITY_TYPE_META[entity.type as keyof typeof ENTITY_TYPE_META] || ENTITY_TYPE_META.character
  const relevancePct = Math.min(100, Math.round((entity.relevance_score ?? 0) * 100))
  const properties = entity.properties || {}
  const propEntries = Object.entries(properties).filter(([, v]) => v !== null && v !== undefined)

  return (
    <>
      <div className={styles.detailCard}>
        <div className={styles.detailNameRow}>
          <h2 className={styles.detailName}>{entity.name}</h2>
        </div>
        {entity.aliases && entity.aliases.length > 0 && (
          <div className={styles.aliasTags}>
            {entity.aliases.map((a) => (
              <span key={a} className={styles.aliasTag}>{a}</span>
            ))}
          </div>
        )}
        <span
          className={styles.typeBadge}
          style={{ background: `${meta.color}18`, color: meta.color, border: `1px solid ${meta.color}30` }}
        >
          {meta.label.toUpperCase()}
        </span>

        {entity.description && (
          <p className={styles.detailDescription}>{entity.description}</p>
        )}

        {propEntries.length > 0 && (
          <div className={styles.propertiesGrid}>
            {propEntries.slice(0, 6).map(([key, val]) => (
              <div key={key} className={styles.propertyTile}>
                <div className={styles.propertyLabel}>{key}</div>
                <div className={styles.propertyValue}>{String(val)}</div>
              </div>
            ))}
          </div>
        )}
      </div>

      <div className={styles.relevanceCard}>
        <div className={styles.relevanceHeader}>
          <span className={styles.relevanceTitle}>Story relevance</span>
          <span className={styles.relevanceScore}>{relevancePct}</span>
        </div>
        <div className={styles.relevanceBar}>
          <div className={styles.relevanceFill} style={{ width: `${relevancePct}%` }} />
        </div>
        <p className={styles.relevanceNote}>A relative score from recent mentions and memory activity. It is not a quality grade.</p>
      </div>
    </>
  )
}
