import { useState, useEffect, useMemo } from 'react'
import { useParams } from 'react-router-dom'
import { api } from '../lib/api'
import { type Contradiction } from '../components/contradictions/ContradictionList'
import styles from './ContradictionsPage.module.css'

const SEVERITIES = ['all', 'high', 'medium', 'low']

type LocalStatus = 'open' | 'resolved' | 'dismissed'

export default function ContradictionsPage() {
  const { universeId } = useParams<{ universeId: string }>()
  const [contradictions, setContradictions] = useState<Contradiction[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [filter, setFilter] = useState('all')
  const [statusOverride, setStatusOverride] = useState<Record<string, LocalStatus>>({})
  const [actionError, setActionError] = useState<string | null>(null)

  useEffect(() => {
    if (!universeId) return
    setLoading(true); setError(null)
    api.getContradictions(universeId)
      .then(({ contradictions: raw }) => { setContradictions(raw || []); setLoading(false) })
      .catch((err: Error) => { setError(err.message); setLoading(false) })
  }, [universeId])

  const filtered = useMemo(() =>
    contradictions.filter((c) => filter === 'all' || c.severity === filter),
    [contradictions, filter]
  )

  const handleResolve = async (id: string) => {
    if (!window.confirm('Mark as resolved?')) return
    setStatusOverride((prev) => ({ ...prev, [id]: 'resolved' })); setActionError(null)
    try { await api.resolveContradiction(universeId!, id) }
    catch (err) { setStatusOverride((prev) => ({ ...prev, [id]: 'open' })); setActionError((err as Error).message) }
  }

  const handleDismiss = async (id: string) => {
    setStatusOverride((prev) => ({ ...prev, [id]: 'dismissed' })); setActionError(null)
    try { await api.dismissContradiction(universeId!, id) }
    catch (err) { setStatusOverride((prev) => ({ ...prev, [id]: 'open' })); setActionError((err as Error).message) }
  }

  if (loading) return (
    <div className={styles.wrap}>
      {Array.from({ length: 3 }).map((_, i) => (
        <div key={i} className={styles.card} style={{ marginBottom: 10 }}>
          <div className={`skeleton`} style={{ height: 10, width: '20%', borderRadius: 4, marginBottom: 10 }} />
          <div className={`skeleton`} style={{ height: 13, width: '85%', borderRadius: 4, marginBottom: 5 }} />
          <div className={`skeleton`} style={{ height: 13, width: '60%', borderRadius: 4 }} />
        </div>
      ))}
    </div>
  )

  if (error) return (
    <div className={styles.wrap}>
      <div className={styles.emptyState}>
        <span className={`glyph ${styles.emptyGlyph}`}>△</span>
        <p className={styles.emptyTitle}>Could not load</p>
        <p className={styles.emptyText}>{error}</p>
      </div>
    </div>
  )

  if (contradictions.length === 0) return (
    <div className={styles.wrap}>
      <div className={styles.emptyState}>
        <span className={`glyph ${styles.emptyGlyph}`}>△</span>
        <p className={styles.emptyTitle}>No Contradictions</p>
        <p className={styles.emptyText}>
          No contradictions detected yet. AI analysis checks your entities and plot events for inconsistencies as you write.
        </p>
      </div>
    </div>
  )

  return (
    <div className={styles.wrap}>
      {actionError && <p className={styles.resolveError}>Action failed: {actionError}</p>}

      <div className={styles.filterBar}>
        {SEVERITIES.map((s) => (
          <button
            key={s}
            className={`${styles.filterBtn} ${filter === s ? styles.filterBtnActive : ''}`}
            onClick={() => setFilter(s)}
          >
            {s === 'all' ? 'All' : s.charAt(0).toUpperCase() + s.slice(1)}
          </button>
        ))}
      </div>

      <div className={styles.listWrap}>
        {filtered.map((c) => {
          const status: LocalStatus = statusOverride[c.id] ?? (c.status as LocalStatus) ?? 'open'
          const isSettled = status !== 'open'
          const sevClass = styles[`severity${c.severity.charAt(0).toUpperCase() + c.severity.slice(1)}` as keyof typeof styles] || styles.severityLow
          return (
            <div key={c.id} className={`${styles.card} ${isSettled ? styles.cardResolved : ''}`}>
              <div className={styles.cardHeader}>
                <span className={`${styles.severity} ${sevClass}`}>{c.severity.toUpperCase()}</span>
                {status === 'resolved' && <span className={styles.resolvedLabel}>✓ Resolved</span>}
                {status === 'dismissed' && <span className={styles.dismissedLabel}>Dismissed — marked intentional.</span>}
                {status === 'open' && (
                  <div className={styles.actions}>
                    <button className={styles.resolveBtn} onClick={() => handleResolve(c.id)}>Resolve</button>
                    <button className={styles.dismissBtn} onClick={() => handleDismiss(c.id)}>Dismiss</button>
                  </div>
                )}
              </div>

              <p className={styles.cardMessage}>{c.description}</p>

              {(c.evidence_a || c.evidence_b) && (
                <div className={styles.evidenceGrid}>
                  {c.evidence_a && (
                    <div className={styles.evidencePanel}>
                      <p className={styles.evidenceQuote}>&ldquo;{c.evidence_a}&rdquo;</p>
                      {c.evidence_a_chapter_id && <span className={styles.evidenceTag}>Ch. {c.evidence_a_chapter_id.slice(0, 8)}</span>}
                    </div>
                  )}
                  {c.evidence_b && (
                    <div className={styles.evidencePanel}>
                      <p className={styles.evidenceQuote}>&ldquo;{c.evidence_b}&rdquo;</p>
                      {c.evidence_b_chapter_id && <span className={styles.evidenceTag}>Ch. {c.evidence_b_chapter_id.slice(0, 8)}</span>}
                    </div>
                  )}
                </div>
              )}

              {c.suggestion && (
                <div className={styles.suggestionBox}>
                  <div className={styles.suggestionKicker}>Suggestion</div>
                  <div className={styles.suggestionText}>{c.suggestion}</div>
                </div>
              )}
            </div>
          )
        })}
      </div>
    </div>
  )
}
