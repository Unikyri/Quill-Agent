import { useEffect, useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { useUniverseStore } from '../stores/universeStore'
import { useAuthStore } from '../stores/authStore'
import { api } from '../lib/api'
import styles from './DashboardPage.module.css'

export default function DashboardPage() {
  const { universes, fetchUniverses } = useUniverseStore()
  const { user, logout } = useAuthStore()
  const navigate = useNavigate()

  const [showNewForm, setShowNewForm] = useState(false)
  const [name, setName] = useState('')
  const [genre, setGenre] = useState('sci-fi')
  const [format, setFormat] = useState('novel')
  const [submitError, setSubmitError] = useState<string | null>(null)

  useEffect(() => { fetchUniverses() }, [])

  const handleCreate = async () => {
    if (!name.trim()) { setSubmitError('Name is required'); return }
    setSubmitError(null)
    try {
      await api.createUniverse({ name: name.trim(), genre, format })
      await fetchUniverses()
      setShowNewForm(false)
      setName('')
      setGenre('sci-fi')
      setFormat('novel')
    } catch (err) {
      setSubmitError((err as Error).message || 'Failed to create universe')
    }
  }

  return (
    <div className={styles.layout}>
      <aside className={styles.sidebar}>
        <h1 className={styles.sidebarHeading}>Quill</h1>
        <p className={styles.sidebarSub}>Writer's Desk</p>

        <div className={styles.sidebarDivider} />

        <div className={styles.userSection}>
          <p className={styles.userName}>{user?.display_name}</p>
          <p className={styles.userEmail}>{user?.email}</p>
        </div>

        <div className={styles.sidebarDivider} />

        <div className={styles.stats}>
          <div className={styles.statItem}>
            <p className={styles.statLabel}>Universes</p>
            <p className={styles.statValue}>{universes.length}</p>
          </div>
        </div>

        <div className={styles.memoryBar}>
          <p className={styles.memoryLabel}>Memory</p>
          <div className={styles.memoryTrack}>
            <div className={styles.memoryFill} style={{ width: '24%' }} />
          </div>
          <p className={styles.memoryPercent}>24 GB</p>
        </div>

        <button className={styles.logoutBtn} onClick={logout}>
          Sign Out
        </button>
      </aside>

      <main className={styles.main}>
        <h2 className={styles.mainHeading}>Your Universes</h2>
        <p className={styles.mainSub}>Worlds waiting for ink</p>

        <div className={styles.headerRow}>
          {!showNewForm ? (
            <button
              className={styles.newBtn}
              onClick={() => setShowNewForm(true)}
            >
              + New Universe
            </button>
          ) : (
            <div className={styles.inlineForm}>
              <input
                className={styles.formInput}
                placeholder="Universe name"
                value={name}
                onChange={(e) => setName(e.target.value)}
                onKeyDown={(e) => e.key === 'Enter' && handleCreate()}
              />
              <select className={styles.formSelect} value={genre} onChange={(e) => setGenre(e.target.value)}>
                <option value="sci-fi">Sci-Fi</option>
                <option value="fantasy">Fantasy</option>
                <option value="mystery">Mystery</option>
                <option value="romance">Romance</option>
                <option value="horror">Horror</option>
                <option value="non-fiction">Non-Fiction</option>
                <option value="thriller">Thriller</option>
                <option value="historical">Historical</option>
                <option value="adventure">Adventure</option>
                <option value="comedy">Comedy</option>
                <option value="drama">Drama</option>
              </select>
              <select className={styles.formSelect} value={format} onChange={(e) => setFormat(e.target.value)}>
                <option value="novel">Novel</option>
                <option value="short-story">Short Story</option>
                <option value="screenplay">Screenplay</option>
                <option value="poetry">Poetry</option>
                <option value="essay">Essay</option>
                <option value="article">Article</option>
                <option value="graphic-novel">Graphic Novel</option>
              </select>
              <button className={styles.formSubmit} onClick={handleCreate}>Create</button>
              <button className={styles.formCancel} onClick={() => { setShowNewForm(false); setSubmitError(null) }}>Cancel</button>
            </div>
          )}
          {submitError && <p className={styles.formError}>{submitError}</p>}
        </div>

        {universes.length === 0 ? (
          <div className={styles.emptyCard}>
            <p>No universes yet. Your first world awaits.</p>
          </div>
        ) : (
          <div className={styles.universeGrid}>
            {universes.map((u) => (
              <div
                key={u.id}
                className={styles.universeCard}
                onClick={() => navigate(`/universe/${u.id}`)}
              >
                <h3 className={styles.cardTitle}>{u.name}</h3>
                <div className={styles.cardMeta}>
                  <span className={styles.cardMetaItem}>{u.genre}</span>
                  <span className={styles.cardMetaItem}>{u.format}</span>
                </div>
                {/* ponytail: random progress for visual interest until real data exists */}
                <div className={styles.cardProgress}>
                  <p className={styles.cardProgressLabel}>Progress</p>
                  <div className={styles.cardProgressTrack}>
                    <div
                      className={styles.cardProgressFill}
                      style={{ width: `${(u.id.length % 40) + 20}%` }}
                    />
                  </div>
                </div>
              </div>
            ))}
          </div>
        )}
      </main>
    </div>
  )
}
