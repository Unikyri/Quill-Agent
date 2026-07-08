import { describe, it, expect, beforeEach, vi } from 'vitest'
import { render, screen, waitFor } from '@testing-library/react'
import DecayTimeline from '../DecayTimeline'
import { api } from '../../../lib/api'

vi.mock('../../../lib/api', () => ({
  api: {
    getMemoryStatus: vi.fn(),
    runDecay: vi.fn(),
  },
}))

const getMemoryStatus = api.getMemoryStatus as ReturnType<typeof vi.fn>
const runDecay = api.runDecay as ReturnType<typeof vi.fn>

beforeEach(() => {
  vi.clearAllMocks()
})

describe('DecayTimeline', () => {
  it('renders one polyline per entity plus a threshold line', async () => {
    getMemoryStatus.mockResolvedValue({
      consolidated_count: 0,
      entities: [
        {
          id: 'e1', name: 'Alice', type: 'character', relevance_score: 0.6, status: 'active',
          consolidated: false, lifecycle: 'active',
          history: [
            { score: 0.9, recorded_at: '2026-07-01T00:00:00Z' },
            { score: 0.6, recorded_at: '2026-07-02T00:00:00Z' },
            { score: 0.3, recorded_at: '2026-07-03T00:00:00Z' },
          ],
        },
        {
          id: 'e2', name: 'Bob', type: 'character', relevance_score: 0.8, status: 'active',
          consolidated: false, lifecycle: 'decaying',
          history: [
            { score: 0.8, recorded_at: '2026-07-01T00:00:00Z' },
            { score: 0.7, recorded_at: '2026-07-02T00:00:00Z' },
          ],
        },
      ],
    })

    render(<DecayTimeline universeId="u1" />)

    await waitFor(() => expect(screen.getByTestId('decay-timeline-svg')).toBeInTheDocument())
    expect(screen.getByTestId('decay-polyline-e1')).toBeInTheDocument()
    expect(screen.getByTestId('decay-polyline-e2')).toBeInTheDocument()
    expect(screen.getByTestId('decay-threshold-line')).toBeInTheDocument()
  })

  it('renders empty state without crashing when there are no entities', async () => {
    getMemoryStatus.mockResolvedValue({ consolidated_count: 0, entities: [] })

    render(<DecayTimeline universeId="u1" />)

    await waitFor(() => expect(screen.getByText(/no memory data yet/i)).toBeInTheDocument())
    expect(screen.queryByTestId('decay-timeline-svg')).not.toBeInTheDocument()
  })

  it('renders a dot instead of a polyline for a single-point history', async () => {
    getMemoryStatus.mockResolvedValue({
      consolidated_count: 0,
      entities: [
        {
          id: 'e3', name: 'Carol', type: 'character', relevance_score: 0.5, status: 'active',
          consolidated: false, lifecycle: 'active',
          history: [{ score: 0.5, recorded_at: '2026-07-01T00:00:00Z' }],
        },
      ],
    })

    render(<DecayTimeline universeId="u1" />)

    await waitFor(() => expect(screen.getByTestId('decay-timeline-svg')).toBeInTheDocument())
    expect(screen.queryByTestId('decay-polyline-e3')).not.toBeInTheDocument()
    expect(screen.getByTestId('decay-dot-e3')).toBeInTheDocument()
  })

  it('does not render crossing markers when the entity never crosses the threshold', async () => {
    getMemoryStatus.mockResolvedValue({
      consolidated_count: 0,
      entities: [
        {
          id: 'e4', name: 'Dave', type: 'character', relevance_score: 0.9, status: 'active',
          consolidated: false, lifecycle: 'active',
          history: [
            { score: 0.9, recorded_at: '2026-07-01T00:00:00Z' },
            { score: 0.8, recorded_at: '2026-07-02T00:00:00Z' },
          ],
        },
      ],
    })

    render(<DecayTimeline universeId="u1" />)

    await waitFor(() => expect(screen.getByTestId('decay-timeline-svg')).toBeInTheDocument())
    expect(screen.queryByTestId(/decay-marker-e4-/)).not.toBeInTheDocument()
  })

  it('renders a crossing marker where an entity drops below the threshold', async () => {
    getMemoryStatus.mockResolvedValue({
      consolidated_count: 0,
      entities: [
        {
          id: 'e5', name: 'Eve', type: 'character', relevance_score: 0.1, status: 'archived',
          consolidated: false, lifecycle: 'archived',
          history: [
            { score: 0.2, recorded_at: '2026-07-01T00:00:00Z' },
            { score: 0.1, recorded_at: '2026-07-02T00:00:00Z' },
          ],
        },
      ],
    })

    render(<DecayTimeline universeId="u1" />)

    await waitFor(() => expect(screen.getByTestId('decay-timeline-svg')).toBeInTheDocument())
    expect(screen.getByTestId('decay-marker-e5-archive-1')).toBeInTheDocument()
  })

  it('runs decay and refetches memory-status when the advance-chapter button is clicked', async () => {
    getMemoryStatus.mockResolvedValue({
      consolidated_count: 0,
      entities: [
        {
          id: 'e1', name: 'Alice', type: 'character', relevance_score: 0.6, status: 'active',
          consolidated: false, lifecycle: 'active',
          history: [
            { score: 0.9, recorded_at: '2026-07-01T00:00:00Z' },
            { score: 0.6, recorded_at: '2026-07-02T00:00:00Z' },
          ],
        },
      ],
    })
    runDecay.mockResolvedValue({ ok: true })

    render(<DecayTimeline universeId="u1" />)
    await waitFor(() => expect(getMemoryStatus).toHaveBeenCalledTimes(1))

    screen.getByRole('button', { name: /advance chapter/i }).click()

    await waitFor(() => expect(runDecay).toHaveBeenCalledWith('u1'))
    await waitFor(() => expect(getMemoryStatus).toHaveBeenCalledTimes(2))
  })
})
