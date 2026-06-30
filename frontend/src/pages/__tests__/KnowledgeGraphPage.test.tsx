import { describe, it, expect, vi, beforeEach } from 'vitest'
import { render, screen, waitFor, act } from '@testing-library/react'
import { MemoryRouter, Routes, Route } from 'react-router-dom'
import KnowledgeGraphPage from '../KnowledgeGraphPage'
import { UniverseContext } from '../../contexts/UniverseContext'
import { useGraphStore } from '../../stores/graphStore'

// Mutable box for wsStore graphPings — reassign .value to simulate new pings
const { pingBox } = vi.hoisted(() => {
  const box: { value: Array<Record<string, unknown>> } = { value: [] }
  return { pingBox: box }
})

// Mock api
const mockGetGraph = vi.fn()
vi.mock('../../lib/api', () => ({
  api: {
    getGraph: (...args: unknown[]) => mockGetGraph(...args),
  },
}))

// Mock wsStore — reads pingBox.value so reassigning it gives a new reference.
// Must apply the Zustand selector so `useWSStore(s => s.graphPings)` returns the array, not the wrapper object.
vi.mock('../../stores/wsStore', () => ({
  useWSStore: (selector: unknown) => {
    const state = { graphPings: pingBox.value }
    return typeof selector === 'function' ? (selector as (s: typeof state) => unknown)(state) : state
  },
}))

const defaultContext = {
  universe: { id: 'uni-1', name: 'Test Universe', genre: 'Fantasy', format: 'Novel' },
  works: [],
}

function renderPage() {
  return render(
    <UniverseContext.Provider value={defaultContext}>
      <MemoryRouter initialEntries={['/universe/uni-1/graph']}>
        <Routes>
          <Route path="/universe/:universeId/graph" element={<KnowledgeGraphPage />} />
        </Routes>
      </MemoryRouter>
    </UniverseContext.Provider>
  )
}

beforeEach(() => {
  vi.clearAllMocks()
  pingBox.value = [] // start each test with no pings
  useGraphStore.setState({
    nodes: [],
    edges: [],
    selectedNodeId: null,
    nodeFilter: { character: true, location: true, item: true, event: true, concept: true },
    loading: false,
    error: null,
    _universeId: null,
  })
})

describe('KnowledgeGraphPage', () => {
  it('shows loading state initially', () => {
    // Freeze the promise so loading state persists
    mockGetGraph.mockReturnValue(new Promise(() => {}))
    renderPage()
    expect(screen.getByTestId('loading-state')).toBeInTheDocument()
  })

  it('shows empty state when graph has zero nodes', async () => {
    mockGetGraph.mockResolvedValue({ nodes: [], edges: [] })
    renderPage()

    await waitFor(() => {
      expect(screen.getByText('No Knowledge Graph')).toBeInTheDocument()
    })
  })

  it('renders graph controls and canvas when nodes exist', async () => {
    mockGetGraph.mockResolvedValue({
      nodes: [
        { id: 'n1', type: 'character', position: { x: 0, y: 0 }, data: { label: 'Alice' } },
      ],
      edges: [],
    })
    renderPage()

    await waitFor(() => {
      // Filter bar with checkboxes should be visible
      expect(screen.getByText('character')).toBeInTheDocument()
      expect(screen.getByText('location')).toBeInTheDocument()
    })
  })

  it('shows error state on API failure', async () => {
    mockGetGraph.mockRejectedValue(new Error('Fetch failed'))
    renderPage()

    await waitFor(() => {
      expect(screen.getByTestId('error-state')).toBeInTheDocument()
      expect(screen.getByText('Fetch failed')).toBeInTheDocument()
    })
  })

  it('shows retry button on error', async () => {
    mockGetGraph.mockRejectedValue(new Error('Oops'))
    renderPage()

    await waitFor(() => {
      expect(screen.getByText('Retry')).toBeInTheDocument()
    })
  })

  it('calls refresh when WS graph_updated ping arrives via wsStore', async () => {
    mockGetGraph.mockResolvedValue({
      nodes: [{ id: 'n1', type: 'character', position: { x: 0, y: 0 }, data: { label: 'Alice' } }],
      edges: [],
    })
    renderPage()

    await waitFor(() => {
      expect(screen.getByText('character')).toBeInTheDocument()
    })
    // fetchGraph called once during load
    expect(mockGetGraph).toHaveBeenCalledTimes(1)

    // Simulate WS ping: assign a new array so effect dependency reference changes
    pingBox.value = [{ type: 'graph_updated' }]

    // Trigger re-render without unmounting: produce a new nodes reference
    const { nodes } = useGraphStore.getState()
    act(() => {
      useGraphStore.setState({ nodes: [...nodes] })
    })

    // refresh() calls api.getGraph internally — should now be called twice
    await waitFor(() => {
      expect(mockGetGraph).toHaveBeenCalledTimes(2)
    })
  })
})
