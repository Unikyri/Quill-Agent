import { describe, it, expect, vi, beforeEach } from 'vitest'
import { render, screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import RelationshipsTab from '../RelationshipsTab'

function deferred<T>() {
  let resolve!: (value: T) => void
  const promise = new Promise<T>((res) => { resolve = res })
  return { promise, resolve }
}

vi.mock('../RelationshipsTab.module.css', () => ({ default: new Proxy({}, { get: (_, key) => key }) }))

const mockGetEntityNeighbors = vi.fn()
vi.mock('../../../lib/api', () => ({
  api: {
    getEntityNeighbors: (...args: unknown[]) => mockGetEntityNeighbors(...args),
  },
}))

const limits = { hops: 1, max_hops: 2, node_limit: 96, edge_limit: 160, result_limit: 256 }

function vertexRaw(entityId: string, name: string, type: string) {
  return `{"id":1,"label":"${type}","properties":{"entity_id":"${entityId}","name":"${name}"}}`
}

describe('RelationshipsTab', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  it('shows a loading state while relationships are being fetched', () => {
    mockGetEntityNeighbors.mockReturnValue(new Promise(() => {}))
    render(<RelationshipsTab entityId="e1" universeId="u1" />)
    expect(screen.getByTestId('loading-state')).toBeInTheDocument()
  })

  it('renders each direct relationship with type, name, and relation label', async () => {
    mockGetEntityNeighbors.mockResolvedValue({
      nodes: [
        { id: 'n1', properties: { raw: vertexRaw('e1', 'Alice', 'character') } },
        { id: 'n2', properties: { raw: vertexRaw('e2', 'Bob', 'character') } },
      ],
      edges: [{ id: 'edge1', source: 'e1', target: 'e2', type: 'ally_of' }],
      truncated: false,
      limits,
    })

    render(<RelationshipsTab entityId="e1" universeId="u1" />)

    expect(await screen.findByText('Bob')).toBeInTheDocument()
    expect(screen.getByText('CHARACTER')).toBeInTheDocument()
    expect(screen.getByText('ally of')).toBeInTheDocument()
    expect(mockGetEntityNeighbors).toHaveBeenCalledWith('e1', 'u1', 1)
  })

  it('shows an empty message when the entity has no relationships', async () => {
    mockGetEntityNeighbors.mockResolvedValue({
      nodes: [{ id: 'n1', properties: { raw: vertexRaw('e1', 'Alice', 'character') } }],
      edges: [],
      truncated: false,
      limits,
    })

    render(<RelationshipsTab entityId="e1" universeId="u1" />)

    expect(await screen.findByText('No relationships are recorded for this entity yet.')).toBeInTheDocument()
  })

  it('shows a retryable error when relationships cannot be loaded', async () => {
    mockGetEntityNeighbors
      .mockRejectedValueOnce(new Error('offline'))
      .mockResolvedValueOnce({
        nodes: [
          { id: 'n1', properties: { raw: vertexRaw('e1', 'Alice', 'character') } },
          { id: 'n2', properties: { raw: vertexRaw('e2', 'Bob', 'character') } },
        ],
        edges: [{ id: 'edge1', source: 'e1', target: 'e2', type: 'ally_of' }],
        truncated: false,
        limits,
      })

    render(<RelationshipsTab entityId="e1" universeId="u1" />)

    expect(await screen.findByRole('alert')).toBeInTheDocument()
    const user = userEvent.setup()
    await user.click(screen.getByRole('button', { name: 'Retry' }))

    expect(await screen.findByText('Bob')).toBeInTheDocument()
  })

  it('refetches when entityId changes', async () => {
    mockGetEntityNeighbors
      .mockResolvedValueOnce({
        nodes: [
          { id: 'n1', properties: { raw: vertexRaw('e1', 'Alice', 'character') } },
          { id: 'n2', properties: { raw: vertexRaw('e2', 'Bob', 'character') } },
        ],
        edges: [{ id: 'edge1', source: 'e1', target: 'e2', type: 'ally_of' }],
        truncated: false,
        limits,
      })
      .mockResolvedValueOnce({
        nodes: [{ id: 'n3', properties: { raw: vertexRaw('e3', 'Cara', 'character') } }],
        edges: [],
        truncated: false,
        limits,
      })

    const view = render(<RelationshipsTab entityId="e1" universeId="u1" />)
    expect(await screen.findByText('Bob')).toBeInTheDocument()

    view.rerender(<RelationshipsTab entityId="e3" universeId="u1" />)
    expect(await screen.findByText('No relationships are recorded for this entity yet.')).toBeInTheDocument()
    expect(mockGetEntityNeighbors).toHaveBeenCalledTimes(2)
  })

  it('shows the later entity, not a flash of the earlier one, when the first fetch resolves after the second', async () => {
    type Neighborhood = { nodes: unknown[]; edges: unknown[]; truncated: boolean; limits: typeof limits }
    const first = deferred<Neighborhood>()
    const second = deferred<Neighborhood>()
    mockGetEntityNeighbors.mockReturnValueOnce(first.promise).mockReturnValueOnce(second.promise)

    const view = render(<RelationshipsTab entityId="e1" universeId="u1" />)
    view.rerender(<RelationshipsTab entityId="e3" universeId="u1" />)

    // The later request (e3) resolves first — simulating out-of-order network responses.
    second.resolve({
      nodes: [{ id: 'n3', properties: { raw: vertexRaw('e3', 'Cara', 'character') } }],
      edges: [],
      truncated: false,
      limits,
    })
    await screen.findByText('No relationships are recorded for this entity yet.')

    // The stale, earlier request (e1) resolves after — its result must be discarded.
    first.resolve({
      nodes: [
        { id: 'n1', properties: { raw: vertexRaw('e1', 'Alice', 'character') } },
        { id: 'n2', properties: { raw: vertexRaw('e2', 'Bob', 'character') } },
      ],
      edges: [{ id: 'edge1', source: 'e1', target: 'e2', type: 'ally_of' }],
      truncated: false,
      limits,
    })
    await waitFor(() => expect(mockGetEntityNeighbors).toHaveBeenCalledTimes(2))

    expect(screen.getByText('No relationships are recorded for this entity yet.')).toBeInTheDocument()
    expect(screen.queryByText('Bob')).not.toBeInTheDocument()
  })
})
