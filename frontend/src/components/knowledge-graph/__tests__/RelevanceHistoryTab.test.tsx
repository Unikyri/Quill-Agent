import { describe, it, expect, vi, beforeEach } from 'vitest'
import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import RelevanceHistoryTab from '../RelevanceHistoryTab'

vi.mock('../RelevanceHistoryTab.module.css', () => ({ default: new Proxy({}, { get: (_, key) => key }) }))

const mockGetMemoryStatus = vi.fn()
vi.mock('../../../lib/api', () => ({
  api: {
    getMemoryStatus: (...args: unknown[]) => mockGetMemoryStatus(...args),
  },
}))

function entity(overrides: Partial<Record<string, unknown>> = {}) {
  return {
    id: 'e1', name: 'Alice', type: 'character', relevance_score: 0.6, status: 'active',
    consolidated: false, lifecycle: 'active',
    history: [
      { score: 0.9, recorded_at: '2024-01-01T00:00:00Z' },
      { score: 0.6, recorded_at: '2024-01-02T00:00:00Z' },
    ],
    ...overrides,
  }
}

describe('RelevanceHistoryTab', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  it('shows a loading state while memory status is being fetched', () => {
    mockGetMemoryStatus.mockReturnValue(new Promise(() => {}))
    render(<RelevanceHistoryTab entityId="e1" universeId="u1" />)
    expect(screen.getByTestId('loading-state')).toBeInTheDocument()
  })

  it('renders the chart for only the selected entity', async () => {
    mockGetMemoryStatus.mockResolvedValue({
      consolidated_count: 0,
      entities: [entity({ id: 'e1' }), entity({ id: 'e2', name: 'Bob' })],
    })

    render(<RelevanceHistoryTab entityId="e1" universeId="u1" />)

    expect(await screen.findByTestId('decay-timeline-svg')).toBeInTheDocument()
    expect(screen.getByTestId('decay-polyline-e1')).toBeInTheDocument()
    expect(screen.queryByTestId('decay-polyline-e2')).not.toBeInTheDocument()
    expect(mockGetMemoryStatus).toHaveBeenCalledWith('u1')
  })

  it('shows an empty message when the entity has no history', async () => {
    mockGetMemoryStatus.mockResolvedValue({
      consolidated_count: 0,
      entities: [entity({ id: 'e1', history: [] })],
    })

    render(<RelevanceHistoryTab entityId="e1" universeId="u1" />)

    expect(await screen.findByText(/No relevance history/i)).toBeInTheDocument()
  })

  it('shows an empty message when the entity is not present in memory status', async () => {
    mockGetMemoryStatus.mockResolvedValue({ consolidated_count: 0, entities: [] })

    render(<RelevanceHistoryTab entityId="e1" universeId="u1" />)

    expect(await screen.findByText(/No relevance history/i)).toBeInTheDocument()
  })

  it('shows a retryable error when memory status cannot be loaded', async () => {
    mockGetMemoryStatus
      .mockRejectedValueOnce(new Error('offline'))
      .mockResolvedValueOnce({ consolidated_count: 0, entities: [entity({ id: 'e1' })] })

    render(<RelevanceHistoryTab entityId="e1" universeId="u1" />)

    expect(await screen.findByRole('alert')).toBeInTheDocument()
    const user = userEvent.setup()
    await user.click(screen.getByRole('button', { name: 'Retry' }))

    expect(await screen.findByTestId('decay-polyline-e1')).toBeInTheDocument()
  })
})
