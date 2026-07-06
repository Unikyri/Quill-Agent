import { useWSStore, type WSStatus } from '../../stores/wsStore'
import { NODE_TYPE_META } from '../knowledge-graph/nodeTypeMeta'
import styles from './ContextPanel.module.css'

interface ContextPanelProps {
  status: WSStatus
}

export default function ContextPanel({ status }: ContextPanelProps) {
  const contradictions = useWSStore((s) => s.contradictions)
  const discoveredEntities = useWSStore((s) => s.discoveredEntities)
  const recallItems = useWSStore((s) => s.recallItems)
  const graphPings = useWSStore((s) => s.graphPings)

  const dismissContradiction = (id: string) => {
    useWSStore.setState((s) => ({ contradictions: s.contradictions.filter((c) => c.id !== id) }))
  }

  const statusClass =
    status === 'open' ? styles.statusOpen
    : status === 'reconnecting' ? styles.statusReconnecting
    : styles.statusClosed

  const isConnected = status === 'open'

  return (
    <div className={styles.panel}>
      <div className={styles.panelHeader}>
        <h3 className={styles.panelTitle}>
          Live Analysis
          {isConnected && <span className={styles.liveIndicator}>● live</span>}
        </h3>
        <span className={`glyph ${styles.statusIndicator} ${statusClass}`} title={`WS: ${status}`}>●</span>
      </div>

      <div className={`${styles.panelBody} q-scroll`}>

        {/* Graph pings */}
        {graphPings.map((_g, i) => (
          <div key={`graph-${i}`} className={styles.graphPing}>
            <span className={`glyph ${styles.graphPingIcon}`}>✳</span>
            <span className={styles.graphPingText}>Knowledge graph updated</span>
            <button
              className={styles.graphPingDismiss}
              onClick={() => useWSStore.setState((s) => ({
                graphPings: s.graphPings.filter((_, idx) => idx !== i),
              }))}
            >✕</button>
          </div>
        ))}

        {/* ── Entities in this paragraph ─ always rendered ─────────── */}
        <div className={styles.section}>
          <div className={styles.sectionHeader}>
            <span className={styles.sectionKicker}>Entities in this paragraph</span>
          </div>
          <div className={styles.sectionBody}>
            {discoveredEntities.length === 0 ? (
              <div className={styles.entityChips}>
                {/* Skeleton chips when idle/loading */}
                <div className={`skeleton ${styles.chipSkeleton}`} />
                <div className={`skeleton ${styles.chipSkeleton}`} style={{ width: 52 }} />
                <div className={`skeleton ${styles.chipSkeleton}`} style={{ width: 78 }} />
              </div>
            ) : (
              <div className={styles.entityChips}>
                {discoveredEntities.map((e) => {
                  const meta = NODE_TYPE_META[e.type || ''] || NODE_TYPE_META.character
                  return (
                    <span
                      key={e.id || e.name}
                      className={styles.entityChip}
                      style={{
                        background: `${meta.color}18`,
                        borderColor: `${meta.color}30`,
                        color: meta.color,
                      }}
                    >
                      <span className={styles.entityChipDot} style={{ background: meta.color }} />
                      {e.name || 'Entity'}
                    </span>
                  )
                })}
              </div>
            )}
          </div>
        </div>

        {/* ── Contradiction detected ─ always rendered ─────────────── */}
        <div className={styles.section}>
          <div className={styles.sectionHeader}>
            <span className={styles.sectionKicker}>⚠ Contradiction detected</span>
          </div>
          <div className={styles.sectionBody}>
            {contradictions.length === 0 ? (
              <>
                <div className={`skeleton ${styles.skRow}`} style={{ width: '90%', height: 40, marginBottom: 8 }} />
                <div className={`skeleton ${styles.skRow}`} style={{ width: '75%', height: 32 }} />
                <p className={styles.emptyPlaceholder} style={{ marginTop: 8 }}>
                  AI contradiction analysis will appear here
                </p>
              </>
            ) : (
              contradictions.map((c) => (
                <div key={c.id || String(Math.random())} className={styles.contradictionCard} style={{ marginBottom: 8 }}>
                  <div className={styles.contradictionKicker}>
                    Contradiction
                    {c.severity && (
                      <span className={styles.severityBadge}>{c.severity.toUpperCase()}</span>
                    )}
                  </div>
                  <p className={styles.contradictionText}>{c.message || String(c)}</p>
                  {(c as any).suggestion && (
                    <div className={styles.suggestionBox}>
                      <div className={styles.suggestionKicker}>Suggestion</div>
                      <div className={styles.suggestionText}>{(c as any).suggestion}</div>
                    </div>
                  )}
                  <div className={styles.contradictionActions}>
                    <button className={styles.resolveBtn}>Resolve</button>
                    <button className={styles.dismissBtn} onClick={() => dismissContradiction(c.id || '')}>
                      Dismiss
                    </button>
                  </div>
                </div>
              ))
            )}
          </div>
        </div>

        {/* ── Relevant memory ─ always rendered ───────────────────── */}
        <div className={styles.section}>
          <div className={styles.sectionHeader}>
            <span className={styles.sectionKicker}>Relevant memory</span>
          </div>
          <div className={styles.sectionBody}>
            {recallItems.length === 0 ? (
              <>
                <div className={`skeleton ${styles.skRow}`} style={{ height: 48, marginBottom: 6 }} />
                <p className={styles.emptyPlaceholder}>Semantic memory appears as you write</p>
              </>
            ) : (
              recallItems.map((r) => (
                <div key={r.id || String(Math.random())} style={{ marginBottom: 8 }}>
                  <p className={styles.memoryQuote}>"{r.fact}"</p>
                  <div className={styles.memorySource}>
                    <span className={styles.memorySrc}>Relevant memory</span>
                    {r.score && (
                      <span className={styles.memoryScore}>
                        {(r.score * 100).toFixed(0)}%
                      </span>
                    )}
                  </div>
                </div>
              ))
            )}
          </div>
        </div>

      </div>
    </div>
  )
}
