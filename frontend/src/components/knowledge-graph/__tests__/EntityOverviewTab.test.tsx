import { describe, it, expect, vi, beforeEach } from 'vitest'
import { render, screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import EntityOverviewTab from '../EntityOverviewTab'

function deferred<T>() {
  let resolve!: (value: T) => void
  const promise = new Promise<T>((res) => { resolve = res })
  return { promise, resolve }
}

vi.mock('../EntityOverviewTab.module.css', () => ({ default: new Proxy({}, { get: (_, key) => key }) }))

const mockGetEntity = vi.fn()
vi.mock('../../../lib/api', () => ({
  api: {
    getEntity: (...args: unknown[]) => mockGetEntity(...args),
  },
}))

describe('EntityOverviewTab', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  it('shows a loading state while the entity is being fetched', () => {
    mockGetEntity.mockReturnValue(new Promise(() => {}))
    render(<EntityOverviewTab entityId="e1" />)
    expect(screen.getByTestId('loading-state')).toBeInTheDocument()
  })

  it('renders name, aliases, type badge, description, properties, and relevance', async () => {
    mockGetEntity.mockResolvedValue({
      entity: {
        id: 'e1', universe_id: 'u1', type: 'character', name: 'Alice',
        aliases: ['Al'], description: 'A curious protagonist.',
        properties: { eyes: 'green' }, status: 'active', relevance_score: 0.42,
      },
    })

    render(<EntityOverviewTab entityId="e1" />)

    expect(await screen.findByText('Alice')).toBeInTheDocument()
    expect(screen.getByText('Al')).toBeInTheDocument()
    expect(screen.getByText('CHARACTER')).toBeInTheDocument()
    expect(screen.getByText('A curious protagonist.')).toBeInTheDocument()
    expect(screen.getByText('eyes')).toBeInTheDocument()
    expect(screen.getByText('green')).toBeInTheDocument()
    expect(screen.getByText('42')).toBeInTheDocument()

    // dropped demo chrome — no portrait placeholder or add-property control
    expect(screen.queryByText('Add portrait')).not.toBeInTheDocument()
    expect(screen.queryByTitle('Add property')).not.toBeInTheDocument()
  })

  it('omits optional sections when the entity has no aliases, description, or properties', async () => {
    mockGetEntity.mockResolvedValue({
      entity: { id: 'e2', universe_id: 'u1', type: 'place', name: 'The Keep', status: 'active', relevance_score: 0 },
    })

    render(<EntityOverviewTab entityId="e2" />)

    expect(await screen.findByText('The Keep')).toBeInTheDocument()
    expect(screen.queryByText('eyes')).not.toBeInTheDocument()
  })

  it('shows a retryable error when the entity cannot be loaded', async () => {
    mockGetEntity.mockRejectedValueOnce(new Error('offline')).mockResolvedValueOnce({
      entity: { id: 'e1', universe_id: 'u1', type: 'character', name: 'Alice', status: 'active', relevance_score: 0.5 },
    })

    render(<EntityOverviewTab entityId="e1" />)

    expect(await screen.findByRole('alert')).toBeInTheDocument()
    const user = userEvent.setup()
    await user.click(screen.getByRole('button', { name: 'Retry' }))

    expect(await screen.findByText('Alice')).toBeInTheDocument()
  })

  it('refetches when entityId changes', async () => {
    mockGetEntity
      .mockResolvedValueOnce({ entity: { id: 'e1', universe_id: 'u1', type: 'character', name: 'Alice', status: 'active', relevance_score: 0.5 } })
      .mockResolvedValueOnce({ entity: { id: 'e2', universe_id: 'u1', type: 'place', name: 'The Keep', status: 'active', relevance_score: 0.2 } })

    const view = render(<EntityOverviewTab entityId="e1" />)
    expect(await screen.findByText('Alice')).toBeInTheDocument()

    view.rerender(<EntityOverviewTab entityId="e2" />)
    expect(await screen.findByText('The Keep')).toBeInTheDocument()
    expect(mockGetEntity).toHaveBeenCalledTimes(2)
  })

  it('shows the later entity, not a flash of the earlier one, when the first fetch resolves after the second', async () => {
    const first = deferred<{ entity: { id: string; universe_id: string; type: string; name: string; status: string; relevance_score: number } }>()
    const second = deferred<{ entity: { id: string; universe_id: string; type: string; name: string; status: string; relevance_score: number } }>()
    mockGetEntity.mockReturnValueOnce(first.promise).mockReturnValueOnce(second.promise)

    const view = render(<EntityOverviewTab entityId="e1" />)
    view.rerender(<EntityOverviewTab entityId="e2" />)

    // The later request (e2) resolves first — simulating out-of-order network responses.
    second.resolve({ entity: { id: 'e2', universe_id: 'u1', type: 'place', name: 'The Keep', status: 'active', relevance_score: 0.2 } })
    await screen.findByText('The Keep')

    // The stale, earlier request (e1) resolves after — its result must be discarded.
    first.resolve({ entity: { id: 'e1', universe_id: 'u1', type: 'character', name: 'Alice', status: 'active', relevance_score: 0.5 } })
    await waitFor(() => expect(mockGetEntity).toHaveBeenCalledTimes(2))

    expect(screen.getByText('The Keep')).toBeInTheDocument()
    expect(screen.queryByText('Alice')).not.toBeInTheDocument()
  })
})
