import { describe, it, expect, vi, beforeEach } from 'vitest'
import { render, screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { MemoryRouter } from 'react-router-dom'
import DashboardPage from '../DashboardPage'

// CSS module mock
vi.mock('../DashboardPage.module.css', () => ({ default: new Proxy({}, { get: (_, k) => k }) }))

const mockFetchUniverses = vi.fn()
let universeStoreState = {
  universes: [] as { id: string; name: string; genre: string; format: string }[],
  fetchUniverses: mockFetchUniverses,
  loading: false,
}
vi.mock('../../stores/universeStore', () => ({
  useUniverseStore: vi.fn((selector?: (state: typeof universeStoreState) => unknown) =>
    selector ? selector(universeStoreState) : universeStoreState
  ),
}))

const mockCreateUniverse = vi.fn()
vi.mock('../../lib/api', () => ({
  api: {
    createUniverse: (...args: unknown[]) => mockCreateUniverse(...args),
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

function renderPage(route = '/dashboard') {
  return render(
    <MemoryRouter initialEntries={[route]}>
      <DashboardPage />
    </MemoryRouter>
  )
}

describe('DashboardPage', () => {
  beforeEach(() => {
    vi.clearAllMocks()
    universeStoreState = {
      universes: [],
      fetchUniverses: mockFetchUniverses,
      loading: false,
    }
  })

  it('redirects to the first universe when universes exist', async () => {
    universeStoreState.universes = [{ id: 'uni-1', name: 'World One', genre: 'fantasy', format: 'novel' }]
    renderPage()

    await waitFor(() => {
      expect(mockNavigate).toHaveBeenCalledWith('/universe/uni-1', { replace: true })
    })
  })

  it('shows a create-first-universe form when no universes exist', async () => {
    renderPage()

    expect(await screen.findByText('Create your first universe')).toBeInTheDocument()
    expect(mockNavigate).not.toHaveBeenCalled()
  })

  it('forces the create form via ?new=true even when universes exist', async () => {
    universeStoreState.universes = [{ id: 'uni-1', name: 'World One', genre: 'fantasy', format: 'novel' }]
    renderPage('/dashboard?new=true')

    expect(await screen.findByText('Create your first universe')).toBeInTheDocument()
    expect(mockNavigate).not.toHaveBeenCalledWith('/universe/uni-1', { replace: true })
  })

  it('shows a loading message while universes are being fetched', () => {
    universeStoreState.loading = true
    renderPage()

    expect(screen.getByText('Entering your universe...')).toBeInTheDocument()
  })

  it('submits the create form and navigates to the new universe', async () => {
    mockCreateUniverse.mockResolvedValue({ universe: { id: 'uni-new' } })
    const user = userEvent.setup()
    renderPage()

    const nameInput = await screen.findByPlaceholderText('Universe Name (e.g. Cosmere)')
    await user.type(nameInput, 'New World')
    await user.click(screen.getByRole('button', { name: 'Create Universe' }))

    await waitFor(() => {
      expect(mockCreateUniverse).toHaveBeenCalledWith(
        expect.objectContaining({ name: 'New World' })
      )
    })
    expect(mockNavigate).toHaveBeenCalledWith('/universe/uni-new')
  })
})
