import { describe, it, expect, vi, beforeEach } from 'vitest'
import { render, screen, waitFor, act, fireEvent } from '@testing-library/react'
import { MemoryRouter, Routes, Route } from 'react-router-dom'
import KnowledgeGraphPage from '../KnowledgeGraphPage'
import { UniverseContext } from '../../contexts/UniverseContext'
import { useGraphStore } from '../../stores/graphStore'

// CSS module mock
vi.mock('../KnowledgeGraphPage.module.css', () => ({ default: new Proxy({}, { get: (_, k) => k }) }))

// Navigate spy for CTA assertion
const mockNavigate = vi.fn()
vi.mock('react-router-dom', async () => {
  const actual = await vi.importActual<typeof import('react-router-dom')>('react-router-dom')
  return { ...actual, useNavigate: () => mockNavigate }
})

// Mutable box for wsStore graphPings — reassign .value to simulate new pings
const { pingBox } = vi.hoisted(() => {
  const box: { value: Array<Record<string, unknown>> } = { value: [] }
  return { pingBox: box }
})

// Mock api
const mockListEntities = vi.fn()
const mockGetEntityNeighbors = vi.fn()
vi.mock('../../lib/api', () => ({
  api: {
    listEntities: (...args: unknown[]) => mockListEntities(...args),
    getEntityNeighbors: (...args: unknown[]) => mockGetEntityNeighbors(...args),
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
  refetchWorks: vi.fn(),
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
  mockNavigate.mockClear()
  pingBox.value = [] // start each test with no pings
  useGraphStore.setState({
    nodes: [],
    edges: [],
    selectedNodeId: null,
    nodeFilter: { character: true, place: true, object: true, event: true, faction: true, world_rule: true, plot_arc: true },
    loading: false,
    error: null,
    _universeId: null,
    focalNodeId: null,
    breadcrumb: [],
  })
})

describe('KnowledgeGraphPage', () => {
  it('shows loading state initially', () => {
    // Freeze the promise so loading state persists
    mockListEntities.mockReturnValue(new Promise(() => {}))
    renderPage()
    expect(screen.getByTestId('loading-state')).toBeInTheDocument()
  })

  it('shows empty state when graph has zero nodes', async () => {
    mockListEntities.mockResolvedValue({ entities: [] })
    renderPage()

    await waitFor(() => {
      expect(screen.getByText('No knowledge graph yet. Ingest a manuscript to build relationships.')).toBeInTheDocument()
    })
  })

  it('renders graph controls and canvas when nodes exist', async () => {
    mockListEntities.mockResolvedValue({ entities: [{ id: 'n1' }] })
    mockGetEntityNeighbors.mockResolvedValue({
      nodes: [
        { id: 'n1', properties: { raw: '{"id":1,"label":"character","properties":{"entity_id":"n1","name":"Alice"}}' } },
      ],
      edges: [],
    })
    renderPage()

    await waitFor(() => {
      // "Character" is rendered both by the filter bar (GraphControls) and the
      // page's own legend — disambiguate with getAllByText instead of getByText.
      expect(screen.getByText('Character')).toBeInTheDocument()
      expect(screen.getByText('Place')).toBeInTheDocument()
    })
  })

  it('shows error state on API failure', async () => {
    mockListEntities.mockRejectedValue(new Error('Fetch failed'))
    renderPage()

    await waitFor(() => {
      expect(screen.getByTestId('error-state')).toBeInTheDocument()
      expect(screen.getByText('Fetch failed')).toBeInTheDocument()
    })
  })

  it('shows retry button on error', async () => {
    mockListEntities.mockRejectedValue(new Error('Oops'))
    renderPage()

    await waitFor(() => {
      expect(screen.getByText('Retry')).toBeInTheDocument()
    })
  })

  it('calls refresh when WS graph_updated ping arrives via wsStore', async () => {
    mockListEntities.mockResolvedValue({ entities: [{ id: 'n1' }] })
    mockGetEntityNeighbors.mockResolvedValue({
      nodes: [{ id: 'n1', properties: { raw: '{"id":1,"label":"character","properties":{"entity_id":"n1","name":"Alice"}}' } }],
      edges: [],
    })
    renderPage()

    await waitFor(() => {
      expect(screen.getByText('Character')).toBeInTheDocument()
    })
    // fetchGraph called once during load
    expect(mockGetEntityNeighbors).toHaveBeenCalledTimes(1)

    // Simulate WS ping: assign a new array so effect dependency reference changes
    pingBox.value = [{ type: 'graph_updated' }]

    // Trigger re-render without unmounting: produce a new nodes reference
    const { nodes } = useGraphStore.getState()
    act(() => {
      useGraphStore.setState({ nodes: [...nodes] })
    })

    // refresh() keeps the current focal neighborhood fresh.
    await waitFor(() => {
      expect(mockGetEntityNeighbors).toHaveBeenCalledTimes(2)
    })
  })

  it('renders CTA button in empty state that navigates to ingestion', async () => {
    mockListEntities.mockResolvedValue({ entities: [] })
    renderPage()

    await waitFor(() => {
      expect(screen.getByText('No knowledge graph yet. Ingest a manuscript to build relationships.')).toBeInTheDocument()
    })

    const ctaButton = screen.getByText('Go to Ingestion')
    expect(ctaButton).toBeInTheDocument()
    expect(ctaButton.tagName).toBe('BUTTON')

    fireEvent.click(ctaButton)
    expect(mockNavigate).toHaveBeenCalledWith('/universe/uni-1/ingest')
  })
})
