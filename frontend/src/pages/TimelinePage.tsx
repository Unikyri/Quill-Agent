import { useState, useEffect } from 'react'
import { useParams } from 'react-router-dom'
import { api } from '../lib/api'
import styles from './TimelinePage.module.css'

interface TimelineEvent {
  id: string; label: string; timestamp: string; description?: string
}

export default function TimelinePage() {
  const { universeId } = useParams<{ universeId: string }>()
  const [events, setEvents] = useState<TimelineEvent[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  const fetchTimeline = () => {
    if (!universeId) return
    setLoading(true); setError(null)
    api.getTimeline(universeId)
      .then(({ events: raw }) => {
        const sorted = (raw || [])
          .map((e: TimelineEvent) => ({ ...e, timestamp: e.timestamp || '' }))
          .sort((a: TimelineEvent, b: TimelineEvent) => (a.timestamp || '').localeCompare(b.timestamp || ''))
        setEvents(sorted)
        setLoading(false)
      })
      .catch((err: Error) => { setError(err.message); setLoading(false) })
  }

  useEffect(() => { fetchTimeline() }, [universeId]) // eslint-disable-line react-hooks/exhaustive-deps

  if (loading) return (
    <div className={styles.wrap}>
      {Array.from({ length: 5 }).map((_, i) => (
        <div key={i} className={styles.timelineItem}>
          <div className={styles.timelineDot} style={{ background: 'var(--muted-3)', boxShadow: 'none' }} />
          <div className={styles.timelineContent}>
            <div className={`skeleton`} style={{ height: 9, width: '30%', borderRadius: 4, marginBottom: 6 }} />
            <div className={`skeleton`} style={{ height: 13, width: '70%', borderRadius: 4 }} />
          </div>
        </div>
      ))}
    </div>
  )

  if (error) return (
    <div className={styles.wrap}>
      <div className={styles.emptyState}>
        <span className={`glyph ${styles.emptyGlyph}`}>⌇</span>
        <p className={styles.emptyTitle}>Could not load timeline</p>
        <p className={styles.emptyText}>{error}</p>
      </div>
    </div>
  )

  if (events.length === 0) return (
    <div className={styles.wrap}>
      <div className={styles.emptyState}>
        <span className={`glyph ${styles.emptyGlyph}`}>⌇</span>
        <p className={styles.emptyTitle}>No Timeline Events</p>
        <p className={styles.emptyText}>
          AI extracts timeline events from your chapters automatically during ingestion. Upload a manuscript to get started.
        </p>
      </div>
    </div>
  )

  return (
    <div className={styles.wrap}>
      <div className={styles.timelineList}>
        {events.map((event, i) => (
          <div key={event.id || i} className={styles.timelineItem}>
            <div className={styles.timelineDot} />
            <div className={styles.timelineContent}>
              {event.timestamp && (
                <div className={styles.timelineEra}>{event.timestamp}</div>
              )}
              <div className={styles.timelineLabel}>{event.label}</div>
              {event.description && (
                <div className={styles.timelineDescription}>{event.description}</div>
              )}
            </div>
          </div>
        ))}
      </div>
    </div>
  )
}
