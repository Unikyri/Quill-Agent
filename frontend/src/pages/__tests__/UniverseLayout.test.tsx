import { describe, it, expect, vi, beforeEach } from 'vitest'
import { render, screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { MemoryRouter, Routes, Route } from 'react-router-dom'
import UniverseLayout from '../UniverseLayout'

// Mock api
const mockGetUniverse = vi.fn()
const mockListWorks = vi.fn()

vi.mock('../../lib/api', () => ({
  api: {
    getUniverse: (...args: unknown[]) => mockGetUniverse(...args),
    listWorks: (...args: unknown[]) => mockListWorks(...args),
  },
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
  mockGetUniverse.mockResolvedValue({
    universe: { id: 'uni-1', name: 'Middle Earth', genre: 'Fantasy', format: 'Novel Series' },
  })
  mockListWorks.mockResolvedValue({
    works: [{ id: 'w1', title: 'The Hobbit', type: 'novel', order_index: 1 }],
  })
})

describe('UniverseLayout', () => {
  it('shows loading state while fetching', () => {
    mockGetUniverse.mockReturnValue(new Promise(() => {}))
    mockListWorks.mockReturnValue(new Promise(() => {}))
    renderLayout()
    expect(screen.getByText('Loading universe…')).toBeInTheDocument()
  })

  it('renders universe name and 5 tabs after load', async () => {
    renderLayout()

    await waitFor(() => {
      expect(screen.getByText('Middle Earth')).toBeInTheDocument()
    })

    expect(screen.getByText('Fantasy · Novel Series')).toBeInTheDocument()
    expect(screen.getByText('Works')).toBeInTheDocument()
    expect(screen.getByText('Graph')).toBeInTheDocument()
    expect(screen.getByText('Timeline')).toBeInTheDocument()
    expect(screen.getByText('Contradictions')).toBeInTheDocument()
    expect(screen.getByText('Plot-holes')).toBeInTheDocument()
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

    // Graph tab content should render
    await waitFor(() => {
      expect(screen.getByText('Graph Content')).toBeInTheDocument()
    })
    // Works tab content should be gone
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

  it('has back to dashboard button', async () => {
    renderLayout()

    await waitFor(() => {
      expect(screen.getByText('← Back to Dashboard')).toBeInTheDocument()
    })
  })
})
