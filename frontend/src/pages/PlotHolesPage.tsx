import { useState, useEffect } from 'react'
import { useParams } from 'react-router-dom'
import { api } from '../lib/api'
import { type PlotHole } from '../components/plot-holes/PlotHoleList'
import styles from './PlotHolesPage.module.css'

export default function PlotHolesPage() {
  const { universeId } = useParams<{ universeId: string }>()
  const [plotHoles, setPlotHoles] = useState<PlotHole[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    if (!universeId) return
    setLoading(true); setError(null)
    api.getPlotHoles(universeId)
      .then(({ plot_holes }) => { setPlotHoles(plot_holes || []); setLoading(false) })
      .catch((err: Error) => { setError(err.message); setLoading(false) })
  }, [universeId])

  if (loading) return (
    <div className={styles.wrap}>
      {Array.from({ length: 3 }).map((_, i) => (
        <div key={i} className={styles.card} style={{ marginBottom: 10 }}>
          <div className={`skeleton`} style={{ height: 10, width: '18%', borderRadius: 4, marginBottom: 10 }} />
          <div className={`skeleton`} style={{ height: 13, width: '80%', borderRadius: 4, marginBottom: 5 }} />
          <div className={`skeleton`} style={{ height: 13, width: '50%', borderRadius: 4 }} />
        </div>
      ))}
    </div>
  )

  if (error) return (
    <div className={styles.wrap}>
      <div className={styles.emptyState}>
        <span className={`glyph ${styles.emptyGlyph}`}>◠</span>
        <p className={styles.emptyTitle}>Could not load</p>
        <p className={styles.emptyText}>{error}</p>
      </div>
    </div>
  )

  if (plotHoles.length === 0) return (
    <div className={styles.wrap}>
      <div className={styles.emptyState}>
        <span className={`glyph ${styles.emptyGlyph}`}>◠</span>
        <p className={styles.emptyTitle}>No Plot Holes</p>
        <p className={styles.emptyText}>
          No plot holes detected. AI analysis scans your works for narrative gaps, inconsistencies, and unresolved threads.
        </p>
      </div>
    </div>
  )

  return (
    <div className={styles.wrap}>
      <div className={styles.listWrap}>
        {plotHoles.map((ph) => (
          <div key={ph.id} className={styles.card}>
            <div className={styles.cardHeader}>
              <span className={styles.cardKicker}>Plot Hole</span>
              <div className={styles.actions}>
                <button className={styles.resolveBtn}>Resolve</button>
                <button className={styles.dismissBtn}>Dismiss</button>
              </div>
            </div>
            <p className={styles.cardDescription}>{ph.description}</p>
            {ph.first_mentioned_chapter_id && (
              <div className={styles.cardMeta}>Chapter {ph.first_mentioned_chapter_id.slice(0, 8)}</div>
            )}
          </div>
        ))}
      </div>
    </div>
  )
}
