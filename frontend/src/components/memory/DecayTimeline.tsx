import { useCallback, useEffect, useState } from 'react'
import { api } from '../../lib/api'
import type { MemoryHistoryPoint, MemoryStatusEntity } from '../../lib/types'
import styles from './DecayTimeline.module.css'

interface DecayTimelineProps {
  universeId: string
}

// Mirrors backend RelevanceService's ARCHIVE_THRESHOLD default (config.go:66).
// Not part of the memory-status payload, so it's hardcoded here per design.
const ARCHIVE_THRESHOLD = 0.15

const VIEW_W = 800
const VIEW_H = 300
const PAD = 30
const INNER_W = VIEW_W - PAD * 2
const INNER_H = VIEW_H - PAD * 2

// Reuses ContextPanel's LIFECYCLE_META color convention without importing
// from ContextPanel.tsx (which must not be modified per spec).
const LIFECYCLE_META: Record<string, { color: string; label: string }> = {
  active: { color: 'var(--success-2)', label: 'active' },
  decaying: { color: 'var(--gold-ink)', label: 'decaying' },
  archived: { color: 'var(--muted-3)', label: 'archived' },
  consolidated: { color: 'var(--node-event)', label: 'consolidated' },
  reactivated: { color: 'var(--teal)', label: 'reactivated' },
}

function scoreY(score: number) {
  return PAD + (1 - score) * INNER_H
}

function pointX(i: number, len: number) {
  if (len <= 1) return PAD
  return PAD + (i / (len - 1)) * INNER_W
}

interface Crossing {
  index: number
  kind: 'archive' | 'reactivate'
}

function findCrossings(history: MemoryHistoryPoint[]): Crossing[] {
  const crossings: Crossing[] = []
  for (let i = 1; i < history.length; i++) {
    const prev = history[i - 1].score
    const cur = history[i].score
    if (prev > ARCHIVE_THRESHOLD && cur <= ARCHIVE_THRESHOLD) {
      crossings.push({ index: i, kind: 'archive' })
    } else if (prev <= ARCHIVE_THRESHOLD && cur > ARCHIVE_THRESHOLD) {
      crossings.push({ index: i, kind: 'reactivate' })
    }
  }
  return crossings
}

function EntityLine({ entity }: { entity: MemoryStatusEntity }) {
  const meta = LIFECYCLE_META[entity.lifecycle] || LIFECYCLE_META.active
  const { history } = entity

  if (history.length === 0) return null

  if (history.length === 1) {
    const x = pointX(0, 1)
    const y = scoreY(history[0].score)
    return (
      <circle
        data-testid={`decay-dot-${entity.id}`}
        cx={x}
        cy={y}
        r={4}
        fill={meta.color}
      >
        <title>{`${entity.name} — ${history[0].score.toFixed(2)} (${meta.label})`}</title>
      </circle>
    )
  }

  const points = history.map((h, i) => `${pointX(i, history.length)},${scoreY(h.score)}`).join(' ')
  const crossings = findCrossings(history)

  return (
    <g>
      <polyline
        data-testid={`decay-polyline-${entity.id}`}
        points={points}
        fill="none"
        stroke={meta.color}
        strokeWidth={2}
      >
        <title>{`${entity.name} (${meta.label})`}</title>
      </polyline>
      {crossings.map((c) => (
        <text
          key={c.index}
          data-testid={`decay-marker-${entity.id}-${c.kind}-${c.index}`}
          x={pointX(c.index, history.length)}
          y={scoreY(history[c.index].score) + (c.kind === 'archive' ? 14 : -8)}
          textAnchor="middle"
          fontSize={12}
          fill={meta.color}
        >
          {c.kind === 'archive' ? '▼' : '▲'}
        </text>
      ))}
    </g>
  )
}

export default function DecayTimeline({ universeId }: DecayTimelineProps) {
  const [entities, setEntities] = useState<MemoryStatusEntity[]>([])
  const [running, setRunning] = useState(false)

  const fetchStatus = useCallback(() => {
    return api.getMemoryStatus(universeId)
      .then((res) => setEntities(res.entities))
      .catch(() => { /* keep last-known entities on transient fetch failure */ })
  }, [universeId])

  useEffect(() => {
    fetchStatus()
  }, [fetchStatus])

  const handleAdvanceChapter = async () => {
    setRunning(true)
    try {
      await api.runDecay(universeId)
      await fetchStatus()
    } finally {
      setRunning(false)
    }
  }

  const thresholdY = scoreY(ARCHIVE_THRESHOLD)

  return (
    <div className={styles.wrap}>
      <div className={styles.header}>
        <span className={styles.kicker}>Decay Timeline</span>
        <button className={styles.advanceBtn} onClick={handleAdvanceChapter} disabled={running}>
          Advance chapter → run decay
        </button>
      </div>

      {entities.length === 0 ? (
        <p className={styles.emptyPlaceholder}>No memory data yet — entities appear here once tracked</p>
      ) : (
        <>
          <svg
            data-testid="decay-timeline-svg"
            className={styles.svg}
            viewBox={`0 0 ${VIEW_W} ${VIEW_H}`}
            preserveAspectRatio="xMidYMid meet"
          >
            <line
              data-testid="decay-threshold-line"
              x1={PAD}
              y1={thresholdY}
              x2={VIEW_W - PAD}
              y2={thresholdY}
              stroke="var(--muted-3)"
              strokeWidth={1}
              strokeDasharray="4 3"
            />
            {entities.map((e) => (
              <EntityLine key={e.id} entity={e} />
            ))}
          </svg>
          <div className={styles.legend}>
            {entities.map((e) => {
              const meta = LIFECYCLE_META[e.lifecycle] || LIFECYCLE_META.active
              return (
                <span key={e.id} className={styles.legendItem}>
                  <span className={styles.legendDot} style={{ background: meta.color }} />
                  {e.name}
                </span>
              )
            })}
          </div>
        </>
      )}
    </div>
  )
}
