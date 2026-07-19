import { useEffect, useState } from 'react'
import { Link } from 'react-router-dom'
import { api, type EntityMentionDTO } from '../../lib/api'
import { writePath } from '../../lib/canonicalRoutes'
import PageStatus from '../shared/PageStatus'
import styles from './MentionsTab.module.css'

interface MentionsTabProps {
  entityId: string
  universeId: string
}

// Mentions tab for the Story Graph right panel — every persisted mention of
// the selected entity, each linking back to its source chapter in the
// Editor. The Editor route (`writePath`) is chapter-scoped, not
// paragraph-scoped, so the link opens the chapter rather than a specific
// paragraph — the same fallback ReviewPage already uses for chapter
// references (`chapter_id.slice(0, 8)`, no chapter-title lookup available).
export default function MentionsTab({ entityId, universeId }: MentionsTabProps) {
  const [mentions, setMentions] = useState<EntityMentionDTO[] | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [retry, setRetry] = useState(0)

  useEffect(() => {
    if (!entityId || !universeId) {
      setLoading(false)
      setMentions(null)
      return
    }
    let cancelled = false
    setLoading(true)
    setError(null)
    setMentions(null)
    api.getEntityMentions(entityId, universeId)
      .then((res) => {
        if (cancelled) return
        setMentions(res.mentions || [])
        setLoading(false)
      })
      .catch(() => {
        if (cancelled) return
        setError('Could not load mentions. Retry to try again.')
        setLoading(false)
      })
    return () => { cancelled = true }
  }, [entityId, universeId, retry])

  if (loading) return <PageStatus loading />
  if (error) return <PageStatus error={error} onRetry={() => setRetry((attempt) => attempt + 1)} />
  if (!mentions) return null

  if (mentions.length === 0) {
    return <p className={styles.empty}>No mentions are recorded for this entity yet.</p>
  }

  return (
    <ul className={styles.list}>
      {mentions.map((mention) => (
        <li key={mention.id} className={styles.item}>
          <div className={styles.itemHeader}>
            <Link className={styles.chapterLink} to={writePath(universeId, mention.chapter_id)}>
              Chapter {mention.chapter_id.slice(0, 8)}
            </Link>
            <span className={styles.paragraph}>Paragraph {mention.paragraph_index}</span>
          </div>
          {mention.context_snippet && (
            <p className={styles.snippet}>&ldquo;{mention.context_snippet}&rdquo;</p>
          )}
        </li>
      ))}
    </ul>
  )
}
