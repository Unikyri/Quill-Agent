import { useState, useMemo } from 'react'
import { api } from '../../lib/api'
import styles from './ContradictionList.module.css'

export interface Contradiction {
  id: string
  message: string
  severity: string
  entities: string[]
}

interface ContradictionListProps {
  universeId: string
  contradictions: Contradiction[]
}

const SEVERITY_CLASS: Record<string, string> = {
  low: styles.severityLow,
  medium: styles.severityMedium,
  high: styles.severityHigh,
}

export default function ContradictionList({ universeId, contradictions }: ContradictionListProps) {
  const [filter, setFilter] = useState<string>('all')
  const [resolved, setResolved] = useState<Set<string>>(new Set())
  const [resolveError, setResolveError] = useState<string | null>(null)

  const filtered = useMemo(() => {
    return contradictions.filter((c) => {
      if (filter !== 'all' && c.severity !== filter) return false
      return true
    })
  }, [contradictions, filter])

  const handleResolve = async (id: string) => {
    // ponytail: native confirm() per product decision — zero deps, zero UI complexity
    if (!window.confirm('Mark this contradiction as resolved?')) return

    // Optimistic local mark
    setResolved((prev) => new Set(prev).add(id))
    setResolveError(null)

    try {
      await api.resolveContradiction(universeId, id)
    } catch (err) {
      // Roll back optimistic update and surface the error — no silent tolerance
      setResolved((prev) => {
        const next = new Set(prev)
        next.delete(id)
        return next
      })
      setResolveError((err as Error).message)
    }
  }

  const severities = ['all', 'low', 'medium', 'high']

  return (
    <div>
      {resolveError && <p className={styles.resolveError}>Failed to resolve: {resolveError}</p>}
      <div className={styles.filterBar}>
        {severities.map((s) => (
          <button
            key={s}
            className={`${styles.filterBtn} ${filter === s ? styles.filterBtnActive : ''}`}
            onClick={() => setFilter(s)}
          >
            {s === 'all' ? 'All' : s}
          </button>
        ))}
      </div>

      <div className={styles.listWrap}>
        {filtered.map((c) => {
          const isResolved = resolved.has(c.id)
          return (
            <div key={c.id} className={`${styles.card} ${isResolved ? styles.cardResolved : ''}`}>
              <div className={styles.cardHeader}>
                <span className={`${styles.severity} ${SEVERITY_CLASS[c.severity] || styles.severityLow}`}>
                  {c.severity}
                </span>
                {isResolved ? (
                  <span className={styles.resolvedLabel}>✓ Resolved</span>
                ) : (
                  <button className={styles.resolveBtn} onClick={() => handleResolve(c.id)}>
                    Resolve
                  </button>
                )}
              </div>
              <p className={styles.cardMessage}>{c.message}</p>
              {c.entities && c.entities.length > 0 && (
                <div className={styles.cardEntities}>
                  {c.entities.map((e) => (
                    <span key={e} className={styles.entityTag}>{e}</span>
                  ))}
                </div>
              )}
            </div>
          )
        })}
      </div>
    </div>
  )
}
