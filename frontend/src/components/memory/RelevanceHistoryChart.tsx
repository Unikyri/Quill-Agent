import type { MemoryHistoryPoint, MemoryStatusEntity } from '../../lib/types'
import styles from './RelevanceHistoryChart.module.css'

// Extracted from DecayTimeline so the same chart can render either the
// multi-entity Memory Lab view (DecayTimeline, unchanged) or a future
// single-entity Story Graph Relevance-history tab (`compact`).
const ARCHIVE_THRESHOLD = 0.15
const VIEW_W = 800
const VIEW_H = 300
const PAD = 30
const INNER_W = VIEW_W - PAD * 2
const INNER_H = VIEW_H - PAD * 2

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

function pointX(index: number, length: number) {
  if (length <= 1) return PAD
  return PAD + (index / (length - 1)) * INNER_W
}

interface Crossing {
  index: number
  kind: 'archive' | 'reactivate'
}

function findCrossings(history: MemoryHistoryPoint[]): Crossing[] {
  const crossings: Crossing[] = []
  for (let index = 1; index < history.length; index++) {
    const previous = history[index - 1].score
    const current = history[index].score
    if (previous > ARCHIVE_THRESHOLD && current <= ARCHIVE_THRESHOLD) crossings.push({ index, kind: 'archive' })
    if (previous <= ARCHIVE_THRESHOLD && current > ARCHIVE_THRESHOLD) crossings.push({ index, kind: 'reactivate' })
  }
  return crossings
}

function EntityLine({ entity }: { entity: MemoryStatusEntity }) {
  const meta = LIFECYCLE_META[entity.lifecycle] || LIFECYCLE_META.active
  if (entity.history.length === 0) return null
  if (entity.history.length === 1) {
    return <circle data-testid={`decay-dot-${entity.id}`} cx={pointX(0, 1)} cy={scoreY(entity.history[0].score)} r={4} fill={meta.color}><title>{`${entity.name} — ${entity.history[0].score.toFixed(2)} (${meta.label})`}</title></circle>
  }
  const points = entity.history.map((point, index) => `${pointX(index, entity.history.length)},${scoreY(point.score)}`).join(' ')
  return (
    <g>
      <polyline data-testid={`decay-polyline-${entity.id}`} points={points} fill="none" stroke={meta.color} strokeWidth={2}><title>{`${entity.name} (${meta.label})`}</title></polyline>
      {findCrossings(entity.history).map((crossing) => (
        <text key={`${crossing.kind}-${crossing.index}`} data-testid={`decay-marker-${entity.id}-${crossing.kind}-${crossing.index}`} x={pointX(crossing.index, entity.history.length)} y={scoreY(entity.history[crossing.index].score) + (crossing.kind === 'archive' ? 14 : -8)} textAnchor="middle" fontSize={12} fill={meta.color} aria-hidden="true">
          {crossing.kind === 'archive' ? '▼' : '▲'}
        </text>
      ))}
    </g>
  )
}

interface RelevanceHistoryChartProps {
  entities: MemoryStatusEntity[]
  /** Hides the per-entity legend — the caller already shows the entity's name (e.g. a single-entity tab). */
  compact?: boolean
  emptyMessage?: string
}

export default function RelevanceHistoryChart({
  entities,
  compact = false,
  emptyMessage = 'No entity lifecycle data yet. Quill shows this after it has tracked story entities.',
}: RelevanceHistoryChartProps) {
  if (entities.length === 0) return <p className={styles.emptyPlaceholder}>{emptyMessage}</p>

  const thresholdY = scoreY(ARCHIVE_THRESHOLD)

  return (
    <>
      <svg data-testid="decay-timeline-svg" className={styles.svg} viewBox={`0 0 ${VIEW_W} ${VIEW_H}`} preserveAspectRatio="xMidYMid meet" aria-label="Entity relevance over time">
        <line data-testid="decay-threshold-line" x1={PAD} y1={thresholdY} x2={VIEW_W - PAD} y2={thresholdY} stroke="var(--muted-3)" strokeWidth={1} strokeDasharray="4 3" />
        {entities.map((entity) => <EntityLine key={entity.id} entity={entity} />)}
      </svg>
      {!compact && (
        <ul className={styles.legend} aria-label="Entity lifecycle summary">
          {entities.map((entity) => {
            const meta = LIFECYCLE_META[entity.lifecycle] || LIFECYCLE_META.active
            return <li key={entity.id} className={styles.legendItem}><span className={styles.legendDot} style={{ background: meta.color }} />{entity.name}: {meta.label}, <span className={styles.relevanceFigure}>{Math.round(entity.relevance_score * 100)}% relevance</span></li>
          })}
        </ul>
      )}
    </>
  )
}
