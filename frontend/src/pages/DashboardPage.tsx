import { useState, useEffect } from 'react'
import { useNavigate, useLocation } from 'react-router-dom'
import { api } from '../lib/api'
import { useUniverseStore } from '../stores/universeStore'
import styles from './DashboardPage.module.css'

const GENRES = [
  { value: 'fantasy', label: 'Fantasy' },
  { value: 'sci-fi', label: 'Sci-Fi' },
  { value: 'mystery', label: 'Mystery' },
  { value: 'romance', label: 'Romance' },
  { value: 'horror', label: 'Horror' },
  { value: 'non-fiction', label: 'Non-Fiction' },
  { value: 'thriller', label: 'Thriller' },
  { value: 'historical', label: 'Historical' },
  { value: 'adventure', label: 'Adventure' },
  { value: 'comedy', label: 'Comedy' },
  { value: 'drama', label: 'Drama' },
]

const FORMATS = [
  { value: 'novel', label: 'Novel' },
  { value: 'short-story', label: 'Short Story' },
  { value: 'screenplay', label: 'Screenplay' },
  { value: 'poetry', label: 'Poetry' },
  { value: 'essay', label: 'Essay' },
  { value: 'article', label: 'Article' },
  { value: 'graphic-novel', label: 'Graphic Novel' },
]

export default function DashboardPage() {
  const navigate = useNavigate()
  const { universes, fetchUniverses, loading } = useUniverseStore()
  const [showCreate, setShowCreate] = useState(false)
  const [newUniverseName, setNewUniverseName] = useState('')
  const [newUniverseDesc, setNewUniverseDesc] = useState('')
  const [newUniverseGenre, setNewUniverseGenre] = useState('fantasy')
  const [newUniverseFormat, setNewUniverseFormat] = useState('novel')
  const [submitError, setSubmitError] = useState<string | null>(null)
  const location = useLocation()
  const isForcingNew = new URLSearchParams(location.search).get('new') === 'true'

  useEffect(() => {
    fetchUniverses()
  }, [fetchUniverses])

  useEffect(() => {
    if (isForcingNew) {
      setShowCreate(true)
      return
    }
    // Automatically redirect if universes exist and we are not forcing create
    if (!loading && universes.length > 0 && !showCreate) {
      navigate(`/universe/${universes[0].id}`, { replace: true })
    } else if (!loading && universes.length === 0) {
      setShowCreate(true)
    }
  }, [loading, universes, showCreate, navigate, isForcingNew])

  const handleCreate = async (e: React.FormEvent) => {
    e.preventDefault()
    if (!newUniverseName.trim()) return
    setSubmitError(null)
    try {
      const { universe } = await api.createUniverse({
        name: newUniverseName,
        description: newUniverseDesc,
        format: newUniverseFormat,
        genre: newUniverseGenre,
      })
      await fetchUniverses()
      navigate(`/universe/${universe.id}`)
    } catch (err) {
      setSubmitError((err as Error).message || 'Failed to create')
    }
  }

  if (loading || (!showCreate && universes.length > 0 && !isForcingNew)) {
    return (
      <div className={styles.layout}>
        <div className={styles.loading}>Entering your universe...</div>
      </div>
    )
  }

  return (
    <div className={styles.layout}>
      <div className={styles.createCard} style={{ margin: '0 auto', maxWidth: 460, marginTop: '10vh' }}>
        <div className={styles.createHeader} style={{ marginBottom: 24, textAlign: 'center', display: 'flex', flexDirection: 'column', alignItems: 'center', position: 'relative' }}>
          {isForcingNew && (
            <button
              onClick={() => navigate(`/universe/${universes[0]?.id}`)}
              style={{ position: 'absolute', left: 0, top: 0, background: 'none', border: 'none', cursor: 'pointer', color: 'var(--muted)', fontSize: 24 }}
              title="Cancel"
            >
              ×
            </button>
          )}
          <div className={styles.createIcon} style={{ fontSize: 32, color: 'var(--teal)', marginBottom: 12 }}>✧</div>
          <h2 className={styles.createTitle} style={{ fontFamily: 'var(--serif)', fontSize: 28, margin: '0 0 8px' }}>Create your first universe</h2>
          <p className={styles.createSub} style={{ color: 'var(--muted)' }}>Give your new world a name and set its genre.</p>
        </div>
        <form className={styles.createForm} onSubmit={handleCreate} style={{ display: 'flex', flexDirection: 'column', gap: 14 }}>
          <div style={{ display: 'flex', flexDirection: 'column', gap: 4 }}>
            <label style={{ fontSize: 11, fontWeight: 600, color: 'var(--muted-3)', textTransform: 'uppercase', letterSpacing: '0.05em' }}>Name</label>
            <input
              className={styles.createInput}
              style={{ padding: '12px 14px', borderRadius: 'var(--r-md)', border: '1px solid var(--line-strong)', background: 'var(--bg-input)', fontSize: 15, width: '100%' }}
              placeholder="Universe Name (e.g. Cosmere)"
              value={newUniverseName}
              onChange={(e) => setNewUniverseName(e.target.value)}
              autoFocus
            />
          </div>

          <div style={{ display: 'flex', flexDirection: 'column', gap: 4 }}>
            <label style={{ fontSize: 11, fontWeight: 600, color: 'var(--muted-3)', textTransform: 'uppercase', letterSpacing: '0.05em' }}>Description</label>
            <textarea
              className={styles.createInput}
              style={{ padding: '12px 14px', borderRadius: 'var(--r-md)', border: '1px solid var(--line-strong)', background: 'var(--bg-input)', fontSize: 14, resize: 'vertical', width: '100%' }}
              placeholder="Brief description (optional)"
              value={newUniverseDesc}
              onChange={(e) => setNewUniverseDesc(e.target.value)}
              rows={3}
            />
          </div>

          <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: 12 }}>
            <div style={{ display: 'flex', flexDirection: 'column', gap: 4 }}>
              <label style={{ fontSize: 11, fontWeight: 600, color: 'var(--muted-3)', textTransform: 'uppercase', letterSpacing: '0.05em' }}>Genre</label>
              <select
                value={newUniverseGenre}
                onChange={(e) => setNewUniverseGenre(e.target.value)}
                style={{ padding: '12px 14px', borderRadius: 'var(--r-md)', border: '1px solid var(--line-strong)', background: 'var(--bg-input)', color: 'var(--ink)', fontSize: 14, width: '100%' }}
              >
                {GENRES.map((g) => (
                  <option key={g.value} value={g.value}>{g.label}</option>
                ))}
              </select>
            </div>

            <div style={{ display: 'flex', flexDirection: 'column', gap: 4 }}>
              <label style={{ fontSize: 11, fontWeight: 600, color: 'var(--muted-3)', textTransform: 'uppercase', letterSpacing: '0.05em' }}>Format</label>
              <select
                value={newUniverseFormat}
                onChange={(e) => setNewUniverseFormat(e.target.value)}
                style={{ padding: '12px 14px', borderRadius: 'var(--r-md)', border: '1px solid var(--line-strong)', background: 'var(--bg-input)', color: 'var(--ink)', fontSize: 14, width: '100%' }}
              >
                {FORMATS.map((f) => (
                  <option key={f.value} value={f.value}>{f.label}</option>
                ))}
              </select>
            </div>
          </div>

          {submitError && <div className={styles.errorText} style={{ color: 'var(--danger)', fontSize: 13, textAlign: 'center' }}>{submitError}</div>}
          <div className={styles.createActions} style={{ marginTop: 8 }}>
            <button type="submit" className={styles.createBtn} disabled={!newUniverseName.trim()} style={{ width: '100%', padding: '14px', background: 'var(--teal)', color: 'var(--parchment-hi)', border: 'none', borderRadius: 'var(--r-md)', fontWeight: 600, cursor: 'pointer' }}>
              Create Universe
            </button>
          </div>
        </form>
      </div>
    </div>
  )
}
