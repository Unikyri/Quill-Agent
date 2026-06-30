import { useEffect, useState } from 'react'
import { useParams, useNavigate } from 'react-router-dom'
import { api } from '../lib/api'

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
    return <div style={{ padding: 24, color: '#888' }}>Loading…</div>
  }

  if (error) {
    return <div style={{ padding: 24, color: '#ff6b6b' }}>Error: {error}</div>
  }

  return (
    <div style={{ padding: 24 }}>
      <button
        onClick={() => navigate(-1)}
        style={{ background: 'transparent', color: '#6c5ce7', marginBottom: 16 }}
      >
        ← Back
      </button>

      <h1>{work?.title || 'Untitled Work'}</h1>
      {work?.type && <p style={{ color: '#888', marginBottom: 24 }}>{work.type}</p>}

      <h2 style={{ marginBottom: 16 }}>Chapters</h2>
      {chapters.length === 0 ? (
        <div className="card"><p>No chapters yet.</p></div>
      ) : (
        chapters
          .sort((a, b) => a.order_index - b.order_index)
          .map((ch) => (
            <div
              key={ch.id}
              className="card"
              style={{ cursor: 'pointer' }}
              onClick={() => navigate(`/editor/${ch.id}`)}
            >
              <h3>{ch.title}</h3>
              <p style={{ color: '#888', fontSize: 13 }}>
                {ch.word_count > 0 ? `${ch.word_count} words` : 'Empty'}
              </p>
            </div>
          ))
      )}
    </div>
  )
}
