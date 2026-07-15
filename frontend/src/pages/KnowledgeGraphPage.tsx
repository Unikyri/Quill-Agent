import { useEffect, useContext, useRef, useState } from 'react'
import { useParams, useNavigate } from 'react-router-dom'
import { useGraphStore } from '../stores/graphStore'
import { useWSStore } from '../stores/wsStore'
import { UniverseContext } from '../contexts/UniverseContext'
import GraphCanvas from '../components/knowledge-graph/GraphCanvas'
import GraphControls from '../components/knowledge-graph/GraphControls'
import PageStatus from '../components/shared/PageStatus'
import { api } from '../lib/api'
import styles from './KnowledgeGraphPage.module.css'

export default function KnowledgeGraphPage() {
  const { universeId } = useParams<{ universeId: string }>()
  const { universe } = useContext(UniverseContext)
  const navigate = useNavigate()
  const fetchGraph = useGraphStore((s) => s.fetchGraph)
  const refresh = useGraphStore((s) => s.refresh)
  const loading = useGraphStore((s) => s.loading)
  const error = useGraphStore((s) => s.error)
  const nodes = useGraphStore((s) => s.nodes)
  const selectedNodeId = useGraphStore((s) => s.selectedNodeId)
  const focalNodeId = useGraphStore((s) => s.focalNodeId)
  const breadcrumb = useGraphStore((s) => s.breadcrumb)
  const focusNode = useGraphStore((s) => s.focusNode)
  const goBack = useGraphStore((s) => s.goBack)
  const graphPings = useWSStore((s) => s.graphPings)
  const prevPingCount = useRef(graphPings.length)
  const [searchQuery, setSearchQuery] = useState('')
  const [searchResults, setSearchResults] = useState<Array<{ id: string; name: string; type: string }>>([])

  useEffect(() => {
    if (universeId) fetchGraph(universeId)
  }, [universeId]) // eslint-disable-line react-hooks/exhaustive-deps

  useEffect(() => {
    if (graphPings.length > prevPingCount.current) {
      prevPingCount.current = graphPings.length
      refresh()
    }
  }, [graphPings, refresh])

  useEffect(() => {
    if (!universeId || !searchQuery.trim()) {
      setSearchResults([])
      return
    }
    const timer = window.setTimeout(() => {
      api.listEntities(universeId, { search: searchQuery.trim(), limit: '8' })
        .then((response) => setSearchResults(response.entities || []))
        .catch(() => setSearchResults([]))
    }, 180)
    return () => window.clearTimeout(timer)
  }, [searchQuery, universeId])

  // Selected node data
  const selectedNodeData = selectedNodeId
    ? nodes.find((n) => n.id === selectedNodeId)
    : null

  if (loading || error) {
    return (
      <PageStatus
        loading={loading}
        error={error}
        onRetry={() => universeId && fetchGraph(universeId)}
      />
    )
  }

  return (
    <div className={styles.wrap}>
      {/* Graph canvas */}
      <div className={styles.canvasArea}>
        {nodes.length === 0 && !loading ? (
          <div style={{ display: 'flex', flexDirection: 'column', alignItems: 'center', justifyContent: 'center', height: '100%', gap: 12, color: 'var(--muted-2)' }}>
            <span className="glyph" style={{ fontSize: 36, color: 'var(--muted-3)' }}>✳</span>
            <p style={{ fontSize: 14, color: 'var(--muted)' }}>
              No knowledge graph yet. Ingest a manuscript to build relationships.
            </p>
            <button
              className="primary"
              style={{ marginTop: 4 }}
              onClick={() => navigate(`/universe/${universeId}/ingest`)}
            >
              Go to Ingestion
            </button>
          </div>
        ) : (
          <>
            <GraphControls />
            <GraphCanvas />
          </>
        )}
      </div>

      {/* Node info panel — always visible */}
      <div className={`${styles.nodePanel} q-scroll`}>
        {selectedNodeData ? (
          <>
            <div>
              <div className={styles.nodePanelTitle}>Central node</div>
              <div className={styles.nodePanelName}>{selectedNodeData.data?.label || 'Entity'}</div>
            </div>

            <div className={styles.statsGrid}>
              <div className={styles.statTile}>
                <div className={styles.statTileValue}>{String((selectedNodeData.data as { connections?: number })?.connections ?? '—')}</div>
                <div className={styles.statTileLabel}>Links</div>
              </div>
              <div className={styles.statTile}>
              <div className={styles.statTileValue}>2</div>
              <div className={styles.statTileLabel}>Max hops</div>
              </div>
              <div className={styles.statTile}>
                <div className={styles.statTileValue}>{String((selectedNodeData.data as { relevance?: number })?.relevance ?? '—')}</div>
                <div className={styles.statTileLabel}>Relevance</div>
              </div>
            </div>

            <div className={styles.conflictSection}>
              <div className={styles.conflictKicker}>Conflicts</div>
              <div className={styles.conflictItem} style={{ color: 'var(--muted-3)', fontStyle: 'italic' }}>
                No conflicts detected for this node.
              </div>
            </div>
          </>
        ) : (
          <>
            <div>
              <div className={styles.nodePanelTitle}>Knowledge Graph</div>
              <div style={{ fontSize: 12, color: 'var(--muted)', marginBottom: 12, fontStyle: 'italic' }}>
                {loading ? 'Loading graph…' : `${nodes.length} nodes · ${universe?.name || ''}`}
              </div>
            </div>

            {/* Skeleton stats */}
            <div className={styles.statsGrid}>
              {[1,2,3].map((i) => (
                <div key={i} className={styles.statTile}>
                  <div className={`skeleton ${styles.skRow}`} style={{ height: 22, marginBottom: 4 }} />
                  <div className={`skeleton ${styles.skRow}`} style={{ height: 8 }} />
                </div>
              ))}
            </div>

            {/* Skeleton conflicts */}
            <div className={styles.conflictSection}>
              <div className={styles.conflictKicker}>Relations in conflict</div>
              <div className={`skeleton ${styles.skRow}`} style={{ height: 44, marginBottom: 6 }} />
              <p className={styles.emptyNodePanel}>Click a node to see its conflicts</p>
            </div>
          </>
        )}

        <div className={styles.graphSearch}>
          <div className={styles.semanticKicker}>Jump to entity</div>
          <input
            className={styles.semanticInput}
            placeholder="Search a name or alias…"
            value={searchQuery}
            onChange={(e) => setSearchQuery(e.target.value)}
          />
          {searchResults.length > 0 && (
            <div className={styles.searchResults}>
              {searchResults.map((entity) => (
                <button
                  key={entity.id}
                  className={styles.searchResult}
                  onClick={() => { void focusNode(entity.id); setSearchQuery(''); setSearchResults([]) }}
                >
                  <span>{entity.name}</span><small>{entity.type.replace('_', ' ')}</small>
                </button>
              ))}
            </div>
          )}
        </div>

        {breadcrumb.length > 0 && (
          <button className={styles.backButton} onClick={() => void goBack()}>
            ← Back to previous focal
          </button>
        )}

        {focalNodeId && (
          <p className={styles.focalNote}>Showing the focal entity and its two-hop neighborhood.</p>
        )}

        {/* Semantic search stays available for a later recall-specific interaction. */}
        <div className={styles.semanticSection}>
          <div className={styles.semanticKicker}>Semantic Memory</div>
          <input
            className={styles.semanticInput}
            placeholder="How does magic work?"
            value=""
            readOnly
            onKeyDown={(e) => { if (e.key === 'Enter') { /* TODO: query semantic memory */ } }}
          />
        </div>
      </div>
    </div>
  )
}
