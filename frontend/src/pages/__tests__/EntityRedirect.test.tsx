import { describe, it, expect, vi, beforeEach } from 'vitest'
import { fireEvent, render, screen, waitFor } from '@testing-library/react'
import { MemoryRouter, Routes, Route, useLocation } from 'react-router-dom'
import EntityRedirect from '../EntityRedirect'

const mockGetEntity = vi.fn()
vi.mock('../../lib/api', () => ({
  api: { getEntity: (...args: unknown[]) => mockGetEntity(...args) },
}))

function Target() {
  const location = useLocation()
  return (
    <div>
      Canonical explore screen
      <output data-testid="target-location">{location.pathname}{location.search}</output>
    </div>
  )
}

function Dashboard() {
  return <div>Dashboard screen</div>
}

function renderRedirect(entityId = 'ent-1') {
  return render(
    <MemoryRouter initialEntries={[`/entity/${entityId}`]}>
      <Routes>
        <Route path="/entity/:entityId" element={<EntityRedirect />} />
        <Route path="/universe/:universeId/explore/map" element={<Target />} />
        <Route path="/dashboard" element={<Dashboard />} />
      </Routes>
    </MemoryRouter>
  )
}

describe('EntityRedirect', () => {
  beforeEach(() => vi.clearAllMocks())

  it('redirects to the consolidated Story Graph map once the entity resolves', async () => {
    mockGetEntity.mockResolvedValue({ entity: { id: 'ent-1', universe_id: 'uni-2' } })
    renderRedirect('ent-1')

    await waitFor(() => {
      expect(screen.getByText('Canonical explore screen')).toBeInTheDocument()
    })
  })

  it('preserves the resolved entity id as a target query param on the map, not just the universe', async () => {
    mockGetEntity.mockResolvedValue({ entity: { id: 'ent-1', universe_id: 'uni-2' } })
    renderRedirect('ent-1')

    await waitFor(() => {
      expect(screen.getByTestId('target-location')).toHaveTextContent('/universe/uni-2/explore/map?entity=ent-1')
    })
  })

  it('keeps the legacy link visible with an accessible retry when the entity fetch fails', async () => {
    mockGetEntity
      .mockRejectedValueOnce(new Error('not found'))
      .mockResolvedValueOnce({ entity: { id: 'ent-1', universe_id: 'uni-2' } })
    renderRedirect('ent-1')

    expect(await screen.findByRole('alert')).toHaveTextContent('Could not open this entity.')
    expect(screen.queryByText('Dashboard screen')).not.toBeInTheDocument()

    fireEvent.click(screen.getByRole('button', { name: 'Retry' }))

    await waitFor(() => expect(screen.getByText('Canonical explore screen')).toBeInTheDocument())
  })
})
