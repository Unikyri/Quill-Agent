import { describe, it, expect, vi, beforeEach } from 'vitest'
import { render, screen, waitFor } from '@testing-library/react'
import { MemoryRouter, Routes, Route } from 'react-router-dom'
import PlotHolesPage from '../PlotHolesPage'
import { UniverseContext } from '../../contexts/UniverseContext'

// CSS module mock
vi.mock('../PlotHolesPage.module.css', () => ({ default: new Proxy({}, { get: (_, k) => k }) }))

// Mock api
const mockGetPlotHoles = vi.fn()
vi.mock('../../lib/api', () => ({
  api: {
    getPlotHoles: (...args: unknown[]) => mockGetPlotHoles(...args),
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
      <MemoryRouter initialEntries={['/universe/uni-1/plot-holes']}>
        <Routes>
          <Route path="/universe/:universeId/plot-holes" element={<PlotHolesPage />} />
        </Routes>
      </MemoryRouter>
    </UniverseContext.Provider>
  )
}

beforeEach(() => {
  vi.clearAllMocks()
})

describe('PlotHolesPage', () => {
  it('shows loading state initially', () => {
    mockGetPlotHoles.mockReturnValue(new Promise(() => {}))
    renderPage()
    expect(screen.getByTestId('loading-state')).toBeInTheDocument()
  })

  it('shows empty state when no plot holes', async () => {
    mockGetPlotHoles.mockResolvedValue({ plot_holes: [] })
    renderPage()

    await waitFor(() => {
      expect(screen.getByText('No Plot Holes')).toBeInTheDocument()
    })
  })

  it('shows error state on API failure', async () => {
    mockGetPlotHoles.mockRejectedValue(new Error('Server error'))
    renderPage()

    await waitFor(() => {
      expect(screen.getByTestId('error-state')).toBeInTheDocument()
      expect(screen.getByText('Server error')).toBeInTheDocument()
    })
  })

  it('renders plot holes sorted by severity (critical → high → medium → low)', async () => {
    mockGetPlotHoles.mockResolvedValue({
      plot_holes: [
        { id: 'ph1', description: 'Low priority issue', severity: 'low' },
        { id: 'ph2', description: 'Critical gap', severity: 'critical' },
        { id: 'ph3', description: 'Medium problem', severity: 'medium' },
        { id: 'ph4', description: 'High urgency', severity: 'high' },
      ],
    })
    renderPage()

    await waitFor(() => {
      expect(screen.getByText('Critical gap')).toBeInTheDocument()
      expect(screen.getByText('High urgency')).toBeInTheDocument()
      expect(screen.getByText('Medium problem')).toBeInTheDocument()
      expect(screen.getByText('Low priority issue')).toBeInTheDocument()
    })

    // Verify severity order: critical → high → medium → low
    const descs = screen
      .getAllByText(/Critical gap|High urgency|Medium problem|Low priority issue/)
      .map((el) => el.textContent)
    expect(descs).toEqual(['Critical gap', 'High urgency', 'Medium problem', 'Low priority issue'])
  })

  it('renders severity badges with correct text', async () => {
    mockGetPlotHoles.mockResolvedValue({
      plot_holes: [
        { id: 'ph1', description: 'A problem', severity: 'high' },
      ],
    })
    renderPage()

    await waitFor(() => {
      expect(screen.getByText('high')).toBeInTheDocument()
    })
  })
})
