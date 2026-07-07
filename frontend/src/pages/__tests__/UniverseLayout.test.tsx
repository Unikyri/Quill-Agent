import { describe, it, expect, vi, beforeEach } from 'vitest'
import { render, screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { MemoryRouter, Routes, Route } from 'react-router-dom'
import UniverseLayout from '../UniverseLayout'

// CSS module mock
vi.mock('../UniverseLayout.module.css', () => ({ default: new Proxy({}, { get: (_, k) => k }) }))

// Mock api
const mockGetUniverse = vi.fn()
const mockListWorks = vi.fn()
const mockListEntities = vi.fn()
const mockGetContradictions = vi.fn()
const mockGetPlotHoles = vi.fn()
const mockUpdateUniverse = vi.fn()
const mockDeleteUniverse = vi.fn()
const mockCreateUniverse = vi.fn()

vi.mock('../../lib/api', () => ({
  api: {
    getUniverse: (...args: unknown[]) => mockGetUniverse(...args),
    listWorks: (...args: unknown[]) => mockListWorks(...args),
    listEntities: (...args: unknown[]) => mockListEntities(...args),
    getContradictions: (...args: unknown[]) => mockGetContradictions(...args),
    getPlotHoles: (...args: unknown[]) => mockGetPlotHoles(...args),
    updateUniverse: (...args: unknown[]) => mockUpdateUniverse(...args),
    deleteUniverse: (...args: unknown[]) => mockDeleteUniverse(...args),
    createUniverse: (...args: unknown[]) => mockCreateUniverse(...args),
  },
}))

const mockFetchUniverses = vi.fn()
let universeStoreState = {
  universes: [
    { id: 'uni-1', name: 'Middle Earth', genre: 'Fantasy', format: 'Novel Series' },
    { id: 'uni-2', name: 'Second World', genre: 'Sci-Fi', format: 'Novel' },
  ],
  fetchUniverses: mockFetchUniverses,
}
vi.mock('../../stores/universeStore', () => ({
  useUniverseStore: vi.fn((selector?: (state: typeof universeStoreState) => unknown) =>
    selector ? selector(universeStoreState) : universeStoreState
  ),
}))

const mockLogout = vi.fn()
const authStoreState = {
  user: { id: 'u1', email: 'writer@example.com', display_name: 'Author Name' },
  logout: mockLogout,
}
vi.mock('../../stores/authStore', () => ({
  useAuthStore: vi.fn((selector?: (state: typeof authStoreState) => unknown) =>
    selector ? selector(authStoreState) : authStoreState
  ),
}))

const mockNavigate = vi.fn()
vi.mock('react-router-dom', async () => {
  const actual = await vi.importActual<typeof import('react-router-dom')>('react-router-dom')
  return {
    ...actual,
    useNavigate: () => mockNavigate,
  }
})

// Simple child tab for testing
function WorksTab() {
  return <div>Works Content</div>
}
function GraphTab() {
  return <div>Graph Content</div>
}

function renderLayout(initialRoute = '/universe/uni-1/works') {
  return render(
    <MemoryRouter initialEntries={[initialRoute]}>
      <Routes>
        <Route path="/universe/:universeId" element={<UniverseLayout />}>
          <Route path="works" element={<WorksTab />} />
          <Route path="graph" element={<GraphTab />} />
        </Route>
      </Routes>
    </MemoryRouter>
  )
}

beforeEach(() => {
  vi.clearAllMocks()
  universeStoreState = {
    universes: [
      { id: 'uni-1', name: 'Middle Earth', genre: 'Fantasy', format: 'Novel Series' },
      { id: 'uni-2', name: 'Second World', genre: 'Sci-Fi', format: 'Novel' },
    ],
    fetchUniverses: mockFetchUniverses,
  }
  mockGetUniverse.mockResolvedValue({
    universe: { id: 'uni-1', name: 'Middle Earth', genre: 'Fantasy', format: 'Novel Series' },
  })
  mockListWorks.mockResolvedValue({
    works: [{ id: 'w1', title: 'The Hobbit', type: 'novel', order_index: 1 }],
  })
  mockListEntities.mockResolvedValue({
    entities: [],
    pagination: { page: 1, limit: 1, total: 12, total_pages: 12 },
  })
  mockGetContradictions.mockResolvedValue({ contradictions: [] })
  mockGetPlotHoles.mockResolvedValue({ plot_holes: [] })
  vi.stubGlobal('confirm', vi.fn(() => true))
})

describe('UniverseLayout', () => {
  it('shows loading state while fetching', () => {
    mockGetUniverse.mockReturnValue(new Promise(() => {}))
    mockListWorks.mockReturnValue(new Promise(() => {}))
    renderLayout()
    expect(screen.getByText('Loading universe…')).toBeInTheDocument()
  })

  it('renders universe switcher card and all 9 nested shell nav items after load', async () => {
    renderLayout()

    await waitFor(() => {
      expect(screen.getAllByText('Middle Earth').length).toBeGreaterThanOrEqual(1)
    })

    expect(screen.getByText('Fantasy · 12 entities')).toBeInTheDocument()
    expect(screen.getByText('Panorama')).toBeInTheDocument()
    expect(screen.getByRole('link', { name: /Works & Chapters/ })).toBeInTheDocument()
    expect(screen.getByText('Editor')).toBeInTheDocument()
    expect(screen.getByText('Entities')).toBeInTheDocument()
    expect(screen.getByText('Graph')).toBeInTheDocument()
    expect(screen.getByText('Timeline')).toBeInTheDocument()
    expect(screen.getByText('Contradictions')).toBeInTheDocument()
    expect(screen.getByText('Plot Holes')).toBeInTheDocument()
    expect(screen.getByText('Ingestion')).toBeInTheDocument()
  })

  it('renders default Works tab content', async () => {
    renderLayout('/universe/uni-1/works')

    await waitFor(() => {
      expect(screen.getByText('Works Content')).toBeInTheDocument()
    })
  })

  it('navigates to Graph tab on click', async () => {
    const user = userEvent.setup()
    renderLayout('/universe/uni-1/works')

    await waitFor(() => {
      expect(screen.getByText('Graph')).toBeInTheDocument()
    })

    await user.click(screen.getByText('Graph'))

    await waitFor(() => {
      expect(screen.getByText('Graph Content')).toBeInTheDocument()
    })
    expect(screen.queryByText('Works Content')).not.toBeInTheDocument()
  })

  it('shows error state when API fails', async () => {
    mockGetUniverse.mockRejectedValue(new Error('Not found'))
    mockListWorks.mockRejectedValue(new Error('Not found'))
    renderLayout()

    await waitFor(() => {
      expect(screen.getByText(/Failed to load universe/)).toBeInTheDocument()
      expect(screen.getByText(/Not found/)).toBeInTheDocument()
    })
  })

  it('clicking the universe switcher opens a popover listing universes, and clicking a row navigates there', async () => {
    const user = userEvent.setup()
    renderLayout()

    await waitFor(() => {
      expect(screen.getByRole('button', { name: /Middle Earth/ })).toBeInTheDocument()
    })

    await user.click(screen.getByRole('button', { name: /Middle Earth/ }))

    expect(screen.getByText('Second World')).toBeInTheDocument()
    expect(screen.getByText(/Create New Universe/)).toBeInTheDocument()

    await user.click(screen.getByText('Second World'))
    expect(mockNavigate).toHaveBeenCalledWith('/universe/uni-2')
  })

  it('clicking Edit on a universe row opens the edit modal pre-filled with its name', async () => {
    const user = userEvent.setup()
    renderLayout()

    await waitFor(() => {
      expect(screen.getByRole('button', { name: /Middle Earth/ })).toBeInTheDocument()
    })
    await user.click(screen.getByRole('button', { name: /Middle Earth/ }))
    await user.click(screen.getAllByTitle('Edit Universe')[0])

    expect(screen.getByText('Edit Universe')).toBeInTheDocument()
    expect(screen.getByDisplayValue('Middle Earth')).toBeInTheDocument()
  })

  it('clicking Delete on the current universe confirms, deletes, and navigates to the dashboard', async () => {
    const user = userEvent.setup()
    renderLayout()

    await waitFor(() => {
      expect(screen.getByRole('button', { name: /Middle Earth/ })).toBeInTheDocument()
    })
    await user.click(screen.getByRole('button', { name: /Middle Earth/ }))
    await user.click(screen.getAllByTitle('Delete Universe')[0])

    await waitFor(() => {
      expect(mockDeleteUniverse).toHaveBeenCalledWith('uni-1')
    })
    expect(mockNavigate).toHaveBeenCalledWith('/dashboard')
  })

  it('shows the current user in the sidebar footer and signs out on click', async () => {
    const user = userEvent.setup()
    renderLayout()

    await waitFor(() => {
      expect(screen.getByText('Author Name')).toBeInTheDocument()
    })
    expect(screen.getByText('writer@example.com')).toBeInTheDocument()

    await user.click(screen.getByRole('button', { name: /Sign out/i }))
    expect(mockLogout).toHaveBeenCalled()
  })

  it('collapses the sidebar and reveals a menu toggle in the header', async () => {
    const user = userEvent.setup()
    renderLayout()

    await waitFor(() => {
      expect(screen.getByRole('link', { name: /Works & Chapters/ })).toBeInTheDocument()
    })

    await user.click(screen.getByRole('button', { name: /Hide sidebar/i }))

    expect(screen.queryByRole('link', { name: /Works & Chapters/ })).not.toBeInTheDocument()
    expect(screen.getByRole('button', { name: /Show sidebar/i })).toBeInTheDocument()
  })

  it('shows the active tab title and a recall search stub in the header', async () => {
    renderLayout('/universe/uni-1/works')

    await waitFor(() => {
      expect(screen.getByText('Recall from the universe…')).toBeInTheDocument()
    })
  })
})
