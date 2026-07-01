import { useEffect, useState } from 'react'
import { useParams, useNavigate } from 'react-router-dom'
import { api } from '../lib/api'
import styles from './WorkPage.module.css'

interface Work {
  id: string
  title: string
  type: string
  universe_id: string
}

interface Chapter {
  id: string
  title: string
  order_index: number
  word_count: number
}

export default function WorkPage() {
  const { workId } = useParams<{ workId: string }>()
  const navigate = useNavigate()
  const [work, setWork] = useState<Work | null>(null)
  const [chapters, setChapters] = useState<Chapter[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    if (!workId) return
    setLoading(true)
    setError(null)

    Promise.all([
      api.getWork(workId),
      api.listChapters(workId),
    ])
      .then(([{ work }, { chapters }]) => {
        setWork(work)
        setChapters(chapters)
      })
      .catch((err) => {
        setError(err.message || 'Failed to load work')
      })
      .finally(() => setLoading(false))
  }, [workId])

  if (loading) {
    return <p className={styles.loading}>Loading…</p>
  }

  if (error) {
    return <p className={styles.error}>Error: {error}</p>
  }

  return (
    <div className={styles.wrap}>
      <button className={styles.backBtn} onClick={() => navigate(-1)}>
        ← Back
      </button>

      <h1 className={styles.heading}>{work?.title || 'Untitled Work'}</h1>
      {work?.type && <p className={styles.type}>{work.type}</p>}

      <h2 className={styles.sectionHeading}>Chapters</h2>
      {chapters.length === 0 ? (
        <p className={styles.empty}>No chapters yet.</p>
      ) : (
        chapters
          .sort((a, b) => a.order_index - b.order_index)
          .map((ch) => (
            <div
              key={ch.id}
              className={styles.chapterCard}
              onClick={() => navigate(`/editor/${ch.id}`)}
            >
              <h3 className={styles.chapterTitle}>{ch.title}</h3>
              <p className={styles.chapterMeta}>
                {ch.word_count > 0 ? `${ch.word_count} words` : 'Empty'}
              </p>
            </div>
          ))
      )}
    </div>
  )
}
