import { describe, it, expect, vi, beforeEach } from 'vitest'
import { render, screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { MemoryRouter } from 'react-router-dom'
import MentionsTab from '../MentionsTab'

function deferred<T>() {
  let resolve!: (value: T) => void
  const promise = new Promise<T>((res) => { resolve = res })
  return { promise, resolve }
}

vi.mock('../MentionsTab.module.css', () => ({ default: new Proxy({}, { get: (_, key) => key }) }))

const mockGetEntityMentions = vi.fn()
vi.mock('../../../lib/api', () => ({
  api: {
    getEntityMentions: (...args: unknown[]) => mockGetEntityMentions(...args),
  },
}))

function renderTab(entityId = 'e1', universeId = 'u1') {
  return render(
    <MemoryRouter>
      <MentionsTab entityId={entityId} universeId={universeId} />
    </MemoryRouter>
  )
}

describe('MentionsTab', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  it('shows a loading state while mentions are being fetched', () => {
    mockGetEntityMentions.mockReturnValue(new Promise(() => {}))
    renderTab()
    expect(screen.getByTestId('loading-state')).toBeInTheDocument()
  })

  it('renders each mention with chapter, paragraph index, snippet, and a link to the editor', async () => {
    mockGetEntityMentions.mockResolvedValue({
      mentions: [
        {
          id: 'm1', entity_id: 'e1', chapter_id: 'c1234567-aaaa', paragraph_index: 2,
          character_offset: 40, context_snippet: 'she walked into the hall', mention_type: 'explicit',
          created_at: '2024-01-01T00:00:00Z',
        },
      ],
      total: 1,
    })

    renderTab()

    expect(await screen.findByText(/she walked into the hall/)).toBeInTheDocument()
    expect(screen.getByText(/Paragraph 2/)).toBeInTheDocument()
    const link = screen.getByRole('link', { name: /chapter c123456/i })
    expect(link).toHaveAttribute('href', '/universe/u1/write/c1234567-aaaa')
    expect(mockGetEntityMentions).toHaveBeenCalledWith('e1', 'u1')
  })

  it('shows an empty message when the entity has no mentions', async () => {
    mockGetEntityMentions.mockResolvedValue({ mentions: [], total: 0 })
    renderTab()
    expect(await screen.findByText('No mentions are recorded for this entity yet.')).toBeInTheDocument()
  })

  it('shows a retryable error when mentions cannot be loaded', async () => {
    mockGetEntityMentions
      .mockRejectedValueOnce(new Error('offline'))
      .mockResolvedValueOnce({
        mentions: [{
          id: 'm1', entity_id: 'e1', chapter_id: 'c1', paragraph_index: 0,
          character_offset: 0, context_snippet: 'hello', mention_type: 'explicit', created_at: '2024-01-01T00:00:00Z',
        }],
        total: 1,
      })

    renderTab()

    expect(await screen.findByRole('alert')).toBeInTheDocument()
    const user = userEvent.setup()
    await user.click(screen.getByRole('button', { name: 'Retry' }))

    expect(await screen.findByText(/hello/)).toBeInTheDocument()
  })

  it('refetches when entityId changes', async () => {
    mockGetEntityMentions
      .mockResolvedValueOnce({
        mentions: [{
          id: 'm1', entity_id: 'e1', chapter_id: 'c1', paragraph_index: 0,
          character_offset: 0, context_snippet: 'hello', mention_type: 'explicit', created_at: '2024-01-01T00:00:00Z',
        }],
        total: 1,
      })
      .mockResolvedValueOnce({ mentions: [], total: 0 })

    const view = renderTab('e1')
    expect(await screen.findByText(/hello/)).toBeInTheDocument()

    view.rerender(
      <MemoryRouter>
        <MentionsTab entityId="e2" universeId="u1" />
      </MemoryRouter>
    )
    expect(await screen.findByText('No mentions are recorded for this entity yet.')).toBeInTheDocument()
    expect(mockGetEntityMentions).toHaveBeenCalledTimes(2)
  })

  it('shows the later entity, not a flash of the earlier one, when the first fetch resolves after the second', async () => {
    type Mentions = { mentions: Array<Record<string, unknown>>; total: number }
    const first = deferred<Mentions>()
    const second = deferred<Mentions>()
    mockGetEntityMentions.mockReturnValueOnce(first.promise).mockReturnValueOnce(second.promise)

    const view = renderTab('e1')
    view.rerender(
      <MemoryRouter>
        <MentionsTab entityId="e2" universeId="u1" />
      </MemoryRouter>
    )

    // The later request (e2) resolves first — simulating out-of-order network responses.
    second.resolve({ mentions: [], total: 0 })
    await screen.findByText('No mentions are recorded for this entity yet.')

    // The stale, earlier request (e1) resolves after — its result must be discarded.
    first.resolve({
      mentions: [{
        id: 'm1', entity_id: 'e1', chapter_id: 'c1', paragraph_index: 0,
        character_offset: 0, context_snippet: 'hello', mention_type: 'explicit', created_at: '2024-01-01T00:00:00Z',
      }],
      total: 1,
    })
    await waitFor(() => expect(mockGetEntityMentions).toHaveBeenCalledTimes(2))

    expect(screen.getByText('No mentions are recorded for this entity yet.')).toBeInTheDocument()
    expect(screen.queryByText(/hello/)).not.toBeInTheDocument()
  })
})
