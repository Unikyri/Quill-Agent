import { useEffect, useState, useMemo } from 'react'
import { useParams, useNavigate } from 'react-router-dom'
import { api } from '../lib/api'
import { parseVertexRaw } from '../lib/graphParse'
import { NODE_TYPE_META } from '../components/knowledge-graph/nodeTypeMeta'
import styles from './EntitiesPage.module.css'

interface EntitySummary {
  id: string; name: string; type: string
}

interface Entity {
  id: string; universe_id: string; type: string; name: string
  aliases?: string[]; description?: string; properties?: Record<string, unknown>
  status: string; relevance_score: number
}

interface RelatedEntity {
  id: string; name: string; type: string; relation?: string
}

const TYPE_FILTERS = ['All', 'character', 'place', 'object', 'faction', 'event', 'worldrule']

function getInitial(name: string) { return (name || '?').charAt(0).toUpperCase() }

// ── Entity detail panel ──────────────────────────────────────────────────────
function EntityDetail({ entityId, universeId }: { entityId: string; universeId: string }) {
  const navigate = useNavigate()
  const [entity, setEntity] = useState<Entity | null>(null)
  const [related, setRelated] = useState<RelatedEntity[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    if (!entityId) return
    setLoading(true); setError(null)
    api.getEntity(entityId)
      .then((res: { entity: Entity }) => {
        setEntity(res.entity)
        setLoading(false)
        return api.getEntityNeighbors(entityId, res.entity.universe_id).then((g) => {
          const others = g.nodes
            .map((n) => parseVertexRaw(String(n.properties?.raw || '')))
            .filter((v) => v.entityId && v.entityId !== entityId)
          setRelated(others.map((v) => ({ id: v.entityId, name: v.name, type: v.type })))
        })
      })
      .catch((err: Error) => { setError(err.message); setLoading(false) })
  }, [entityId])

  if (loading) return <p className={styles.loadingText}>Loading entity…</p>
  if (error) return <p className={styles.loadingText} style={{ color: 'var(--danger)' }}>Error: {error}</p>
  if (!entity) return null

  const meta = NODE_TYPE_META[entity.type] || NODE_TYPE_META.character
  const relevancePct = Math.min(100, Math.round((entity.relevance_score ?? 0) * 100))
  const properties = entity.properties || {}
  const propEntries = Object.entries(properties).filter(([, v]) => v !== null && v !== undefined)

  return (
    <>
      {/* Main detail card */}
      <div className={styles.detailCard}>
        <div className={styles.detailHeaderRow}>
          <div className={styles.portraitBox}>
            <span className={styles.portraitGlyph}>⊕</span>
            <span>Add portrait</span>
          </div>
          <div className={styles.detailHeaderInfo}>
            <div className={styles.detailNameRow}>
              <h2 className={styles.detailName}>{entity.name}</h2>
              <span className={`glyph ${styles.editIcon}`} title="Edit name">✎</span>
            </div>
            {entity.aliases && entity.aliases.length > 0 && (
              <div className={styles.aliasTags}>
                {entity.aliases.map((a) => (
                  <span key={a} className={styles.aliasTag}>{a}</span>
                ))}
              </div>
            )}
            <span
              className={styles.typeBadge}
              style={{ background: `${meta.color}18`, color: meta.color, border: `1px solid ${meta.color}30` }}
            >
              {meta.label.toUpperCase()}
            </span>
          </div>
        </div>

        {entity.description && (
          <p className={styles.detailDescription}>{entity.description}</p>
        )}

        {propEntries.length > 0 && (
          <div className={styles.propertiesGrid}>
            {propEntries.slice(0, 6).map(([key, val]) => (
              <div key={key} className={styles.propertyTile}>
                <div className={styles.propertyLabel}>{key}</div>
                <div className={styles.propertyValue}>{String(val)}</div>
              </div>
            ))}
            <button className={styles.addPropertyBtn} title="Add property">+</button>
          </div>
        )}
      </div>

      {/* Relevance */}
      <div className={styles.relevanceCard}>
        <div className={styles.relevanceHeader}>
          <span className={styles.relevanceTitle}>Relevance</span>
          <span className={styles.relevanceScore}>{relevancePct}</span>
        </div>
        <div className={styles.relevanceBar}>
          <div className={styles.relevanceFill} style={{ width: `${relevancePct}%`, background: meta.color }} />
        </div>
        <p className={styles.relevanceNote}>How central this entity is to the universe, according to AI.</p>
      </div>

      {/* Relations */}
      {related.length > 0 && (
        <div className={styles.relationsCard}>
          <div className={styles.relationsHeader}>
            <span className={styles.relationsTitle}>Direct Relations</span>
            <span
              className={styles.relationsLink}
              role="button" tabIndex={0}
              onClick={() => navigate(`/universe/${universeId}/graph`)}
            >
              View graph →
            </span>
          </div>
          <div className={styles.relationChips}>
            {related.map((r) => (
              <span
                key={r.id}
                className={styles.relationChip}
                role="button" tabIndex={0}
                onClick={() => navigate(`/universe/${universeId}/entities/${r.id}`)}
                onKeyDown={(e) => { if (e.key === 'Enter' || e.key === ' ') navigate(`/universe/${universeId}/entities/${r.id}`) }}
              >
                {r.name}
                <span style={{ color: 'var(--muted-3)', marginLeft: 4, fontSize: 10 }}>
                  ({NODE_TYPE_META[r.type]?.label || r.type})
                </span>
              </span>
            ))}
          </div>
        </div>
      )}
    </>
  )
}

// ── Main split-pane component ────────────────────────────────────────────────
export default function EntitiesPage() {
  const { universeId, entityId: paramEntityId } = useParams<{ universeId: string; entityId?: string }>()
  const navigate = useNavigate()

  const [entities, setEntities] = useState<EntitySummary[]>([])
  const [loading, setLoading] = useState(true)
  const [search, setSearch] = useState('')
  const [filter, setFilter] = useState('All')
  const [selectedId, setSelectedId] = useState<string | null>(paramEntityId || null)

  useEffect(() => {
    if (!universeId) return
    setLoading(true)
    api.listEntities(universeId, { limit: '200' })
      .then((res) => { setEntities(res.entities || []); setLoading(false) })
      .catch(() => setLoading(false))
  }, [universeId])

  // Sync URL param
  useEffect(() => {
    if (paramEntityId && paramEntityId !== selectedId) setSelectedId(paramEntityId)
  }, [paramEntityId]) // eslint-disable-line react-hooks/exhaustive-deps

  const filtered = useMemo(() => {
    return entities.filter((e) => {
      if (filter !== 'All' && e.type !== filter) return false
      if (search && !e.name.toLowerCase().includes(search.toLowerCase())) return false
      return true
    })
  }, [entities, filter, search])

  const handleSelect = (id: string) => {
    setSelectedId(id)
    navigate(`/universe/${universeId}/entities/${id}`, { replace: true })
  }

  return (
    <div className={styles.wrap}>
      {/* List rail */}
      <div className={`${styles.listRail} q-scroll`}>
        <div className={styles.searchBar}>
          <input
            className={styles.searchInput}
            placeholder="Search entity or alias…"
            value={search}
            onChange={(e) => setSearch(e.target.value)}
          />
        </div>

        <div className={styles.filterChips}>
          {TYPE_FILTERS.slice(0, 4).map((f) => (
            <button
              key={f}
              className={`${styles.filterChip} ${filter === f ? styles.filterChipActive : ''}`}
              onClick={() => setFilter(f)}
            >
              {f === 'All' ? 'All' : f.charAt(0).toUpperCase() + f.slice(1) + 's'}
            </button>
          ))}
        </div>

        <div className={styles.entityList}>
          {loading ? (
            Array.from({ length: 6 }).map((_, i) => (
              <div key={i} className={styles.entityListItem}>
                <div className={`skeleton ${styles.entityAvatar}`} style={{ background: 'var(--surface-sunken)' }} />
                <div style={{ flex: 1 }}>
                  <div className={`skeleton`} style={{ height: 11, width: '65%', borderRadius: 4, marginBottom: 4 }} />
                  <div className={`skeleton`} style={{ height: 8, width: '40%', borderRadius: 4 }} />
                </div>
              </div>
            ))
          ) : filtered.length === 0 ? (
            <p style={{ padding: 16, fontSize: 12.5, color: 'var(--muted)', textAlign: 'center', fontStyle: 'italic' }}>
              No entities found.
            </p>
          ) : (
            filtered.map((e) => {
              const meta = NODE_TYPE_META[e.type] || NODE_TYPE_META.character
              return (
                <div
                  key={e.id}
                  className={`${styles.entityListItem} ${selectedId === e.id ? styles.entityListItemActive : ''}`}
                  role="button" tabIndex={0}
                  onClick={() => handleSelect(e.id)}
                  onKeyDown={(ev) => { if (ev.key === 'Enter' || ev.key === ' ') handleSelect(e.id) }}
                >
                  <div className={styles.entityAvatar} style={{ background: meta.color }}>
                    {getInitial(e.name)}
                  </div>
                  <div className={styles.entityItemInfo}>
                    <div className={styles.entityItemName}>{e.name}</div>
                    <div className={styles.entityItemType} style={{ color: meta.color }}>
                      {meta.label.toUpperCase()}
                    </div>
                  </div>
                </div>
              )
            })
          )}
        </div>
      </div>

      {/* Detail panel */}
      <div className={`${styles.detailPanel} q-scroll`}>
        {selectedId && universeId ? (
          <EntityDetail entityId={selectedId} universeId={universeId} />
        ) : (
          <div className={styles.emptyDetail}>
            <span className={`glyph ${styles.emptyGlyph}`}>○</span>
            <p className={styles.emptyText}>Select an entity to view its details</p>
          </div>
        )}
      </div>
    </div>
  )
}
