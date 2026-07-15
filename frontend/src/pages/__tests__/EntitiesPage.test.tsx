import { describe, it, expect, vi, beforeEach } from 'vitest'
import { fireEvent, render, screen, waitFor } from '@testing-library/react'
import { MemoryRouter, Route, Routes } from 'react-router-dom'
import EntitiesPage from '../EntitiesPage'

vi.mock('../EntitiesPage.module.css', () => ({ default: new Proxy({}, { get: (_, key) => key }) }))

const mockListEntities = vi.fn()
vi.mock('../../lib/api', () => ({
  api: {
    listEntities: (...args: unknown[]) => mockListEntities(...args),
    getEntity: vi.fn(),
    getEntityNeighbors: vi.fn(),
  },
}))

const counts = { character: 2, place: 0, object: 1, faction: 0, event: 0, world_rule: 0, plot_arc: 0 }

function renderPage() {
  return render(
    <MemoryRouter initialEntries={['/universe/uni-1/entities']}>
      <Routes>
        <Route path="/universe/:universeId/entities" element={<EntitiesPage />} />
        <Route path="/universe/:universeId/entities/:entityId" element={<EntitiesPage />} />
      </Routes>
    </MemoryRouter>,
  )
}

describe('EntitiesPage', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  it('requests the selected type on every paginated page', async () => {
    mockListEntities
      .mockResolvedValueOnce({ entities: [], counts_by_type: counts, pagination: { total: 3 } })
      .mockResolvedValueOnce({
        entities: Array.from({ length: 100 }, (_, index) => ({ id: `object-${index}`, name: `Object ${index}`, type: 'object' })),
        counts_by_type: counts,
        pagination: { total: 101 },
      })
      .mockResolvedValueOnce({
        entities: [{ id: 'object-100', name: 'Object 100', type: 'object' }],
        counts_by_type: counts,
        pagination: { total: 101 },
      })

    renderPage()
    await waitFor(() => expect(mockListEntities).toHaveBeenCalledWith('uni-1', { limit: '100', page: '1' }))

    fireEvent.click(screen.getByRole('button', { name: 'Objects (1)' }))
    await waitFor(() => {
      expect(mockListEntities).toHaveBeenLastCalledWith('uni-1', { limit: '100', page: '1', type: 'object' })
    })

    fireEvent.click(await screen.findByRole('button', { name: 'Load more (100 of 101)' }))
    await waitFor(() => {
      expect(mockListEntities).toHaveBeenLastCalledWith('uni-1', { limit: '100', page: '2', type: 'object' })
    })
  })

  it('requests search terms from the server instead of filtering the loaded page locally', async () => {
    mockListEntities
      .mockResolvedValueOnce({ entities: [], counts_by_type: counts, pagination: { total: 3 } })
      .mockResolvedValueOnce({
        entities: [{ id: 'character-1', name: 'Filip', type: 'character' }],
        counts_by_type: counts,
        pagination: { total: 1 },
      })

    renderPage()
    await waitFor(() => expect(mockListEntities).toHaveBeenCalledTimes(1))

    fireEvent.change(screen.getByPlaceholderText('Search entity or alias…'), { target: { value: 'Fil' } })
    await waitFor(() => {
      expect(mockListEntities).toHaveBeenLastCalledWith('uni-1', { limit: '100', page: '1', search: 'Fil' })
    })
    expect(await screen.findByText('Filip')).toBeInTheDocument()
  })

  it('discards an older pagination response after the query changes', async () => {
    let resolveStalePage: (value: unknown) => void = () => {}
    mockListEntities
      .mockResolvedValueOnce({ entities: [], counts_by_type: counts, pagination: { total: 3 } })
      .mockResolvedValueOnce({
        entities: Array.from({ length: 100 }, (_, index) => ({ id: `object-${index}`, name: `Object ${index}`, type: 'object' })),
        counts_by_type: counts,
        pagination: { total: 101 },
      })
      .mockImplementationOnce(() => new Promise((resolve) => { resolveStalePage = resolve }))
      .mockResolvedValueOnce({
        entities: [{ id: 'character-1', name: 'Character result', type: 'character' }],
        counts_by_type: counts,
        pagination: { total: 2 },
      })

    renderPage()
    await waitFor(() => expect(mockListEntities).toHaveBeenCalledTimes(1))

    fireEvent.click(screen.getByRole('button', { name: 'Objects (1)' }))
    await screen.findByRole('button', { name: 'Load more (100 of 101)' })
    fireEvent.click(screen.getByRole('button', { name: 'Load more (100 of 101)' }))
    await waitFor(() => {
      expect(mockListEntities).toHaveBeenLastCalledWith('uni-1', { limit: '100', page: '2', type: 'object' })
    })

    fireEvent.click(screen.getByRole('button', { name: 'Characters (2)' }))
    await screen.findByText('Character result')

    resolveStalePage({
      entities: [{ id: 'object-100', name: 'Stale object', type: 'object' }],
      counts_by_type: counts,
      pagination: { total: 101 },
    })

    await waitFor(() => {
      expect(screen.queryByText('Stale object')).not.toBeInTheDocument()
      expect(screen.getByText('Character result')).toBeInTheDocument()
    })
  })
})
