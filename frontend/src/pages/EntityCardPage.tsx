import { useEffect, useState } from 'react'
import { useParams, useNavigate } from 'react-router-dom'
import { api } from '../lib/api'
import styles from './EntityCardPage.module.css'

interface Entity {
  id: string
  name: string
  type: string
  description: string
  universe_id: string
  metadata?: Record<string, unknown>
  attributes?: Record<string, string>
  related_entities?: { id: string; name: string; type: string }[]
}

export default function EntityCardPage() {
  const { entityId } = useParams<{ entityId: string }>()
  const navigate = useNavigate()
  const [entity, setEntity] = useState<Entity | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  const fetchEntity = () => {
    if (!entityId) return
    setLoading(true)
    setError(null)
    api.getEntity(entityId)
      .then((res: { entity: Entity }) => {
        setEntity(res.entity)
        setLoading(false)
      })
      .catch((err: Error) => {
        setError(err.message || 'Failed to load entity')
        setLoading(false)
      })
  }

  useEffect(() => {
    fetchEntity()
  }, [entityId]) // eslint-disable-line react-hooks/exhaustive-deps

  if (loading) {
    return <p className={styles.loading}>Loading entity…</p>
  }

  if (error) {
    return (
      <div className={styles.error}>
        <p className={styles.errorText}>Failed to load entity: {error}</p>
        <button className={styles.retryBtn} onClick={fetchEntity}>
          Retry
        </button>
      </div>
    )
  }

  return (
    <div className={styles.wrap}>
      <nav className={styles.navbar}>
        <button className={styles.backBtn} onClick={() => navigate(-1)}>
          ← Back
        </button>
        <span className={styles.breadcrumb}>
          Entity / {entity?.name || 'Unknown'}
        </span>
      </nav>

      <div className={styles.columns}>
        <main className={styles.main}>
          <h1 className={styles.heading}>{entity?.name || 'Entity'}</h1>
          <p className={styles.subtitle}>{entity?.type}</p>

          <p className={styles.description}>
            {entity?.description || 'No description provided.'}
          </p>

          {entity?.attributes && Object.keys(entity.attributes).length > 0 && (
            <div className={styles.section}>
              <p className={styles.sectionLabel}>Attributes</p>
              <div className={styles.tags}>
                {Object.entries(entity.attributes).map(([key, val]) => (
                  <span key={key} className={styles.tag}>
                    {key}: {val}
                  </span>
                ))}
              </div>
            </div>
          )}

          {entity?.metadata && Object.keys(entity.metadata).length > 0 && (
            <div className={styles.section}>
              <p className={styles.sectionLabel}>Metadata</p>
              <div className={styles.tags}>
                {Object.entries(entity.metadata).map(([key, val]) => (
                  <span key={key} className={styles.tag}>
                    {key}: {String(val)}
                  </span>
                ))}
              </div>
            </div>
          )}
        </main>

        <aside className={styles.sidebar}>
          <div className={styles.sidebarCard}>
            <p className={styles.sidebarLabel}>Type</p>
            <p className={styles.sidebarValue}>{entity?.type || '-'}</p>
          </div>

          {entity?.related_entities && entity.related_entities.length > 0 && (
            <div className={styles.sidebarCard}>
              <p className={styles.sidebarLabel}>Related Entities</p>
              <div className={styles.relatedList}>
                {entity.related_entities.map((rel) => (
                  <div
                    key={rel.id}
                    className={styles.relatedItem}
                    onClick={() => navigate(`/entity/${rel.id}`)}
                  >
                    <span>
                      {rel.name} <span style={{ color: 'var(--ink-40)', fontSize: 11 }}>({rel.type})</span>
                    </span>
                  </div>
                ))}
              </div>
            </div>
          )}
        </aside>
      </div>
    </div>
  )
}
