import { describe, it, expect, vi } from 'vitest'
import { render, screen } from '@testing-library/react'
import RelevanceHistoryChart from '../RelevanceHistoryChart'
import type { MemoryStatusEntity } from '../../../lib/types'

vi.mock('../RelevanceHistoryChart.module.css', () => ({ default: new Proxy({}, { get: (_, key) => key }) }))

const alice: MemoryStatusEntity = {
  id: 'e1', name: 'Alice', type: 'character', relevance_score: 0.6, status: 'active',
  consolidated: false, lifecycle: 'active',
  history: [
    { score: 0.9, recorded_at: '2026-07-01T00:00:00Z' },
    { score: 0.6, recorded_at: '2026-07-02T00:00:00Z' },
    { score: 0.1, recorded_at: '2026-07-03T00:00:00Z' },
  ],
}

describe('RelevanceHistoryChart', () => {
  it('renders an empty placeholder when there are no entities', () => {
    render(<RelevanceHistoryChart entities={[]} />)
    expect(screen.getByText(/no entity lifecycle data yet/i)).toBeInTheDocument()
    expect(screen.queryByTestId('decay-timeline-svg')).not.toBeInTheDocument()
  })

  it('renders one polyline per multi-point entity plus a threshold line and legend', () => {
    render(<RelevanceHistoryChart entities={[alice]} />)
    expect(screen.getByTestId('decay-timeline-svg')).toBeInTheDocument()
    expect(screen.getByTestId('decay-polyline-e1')).toBeInTheDocument()
    expect(screen.getByTestId('decay-threshold-line')).toBeInTheDocument()
    expect(screen.getByTestId('decay-marker-e1-archive-2')).toBeInTheDocument()
    expect(screen.getByText(/alice: active/i)).toBeInTheDocument()
  })

  it('renders a dot instead of a polyline for a single-point history', () => {
    const carol: MemoryStatusEntity = {
      id: 'e3', name: 'Carol', type: 'character', relevance_score: 0.5, status: 'active',
      consolidated: false, lifecycle: 'active', history: [{ score: 0.5, recorded_at: '2026-07-01T00:00:00Z' }],
    }
    render(<RelevanceHistoryChart entities={[carol]} />)
    expect(screen.getByTestId('decay-dot-e3')).toBeInTheDocument()
    expect(screen.queryByTestId('decay-polyline-e3')).not.toBeInTheDocument()
  })

  it('hides the legend in compact mode, for a single-entity tab', () => {
    render(<RelevanceHistoryChart entities={[alice]} compact />)
    expect(screen.getByTestId('decay-timeline-svg')).toBeInTheDocument()
    expect(screen.queryByText(/alice: active/i)).not.toBeInTheDocument()
  })

  it('accepts a custom empty message', () => {
    render(<RelevanceHistoryChart entities={[]} emptyMessage="Nothing here yet." />)
    expect(screen.getByText('Nothing here yet.')).toBeInTheDocument()
  })
})
