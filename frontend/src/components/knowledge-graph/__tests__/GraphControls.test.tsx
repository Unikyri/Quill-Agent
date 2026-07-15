import { describe, expect, it, beforeEach } from 'vitest'
import { fireEvent, render, screen } from '@testing-library/react'
import GraphControls from '../GraphControls'
import { ENTITY_TYPE_META, ENTITY_TYPES } from '../../../lib/entityTypes'
import { useGraphStore } from '../../../stores/graphStore'

beforeEach(() => {
  useGraphStore.setState({
    nodeFilter: Object.fromEntries(ENTITY_TYPES.map((type) => [type, true])),
    showArchived: false,
  })
})

describe('GraphControls', () => {
  it('renders every canonical entity type and exposes the archived toggle', () => {
    render(<GraphControls />)

    for (const type of ENTITY_TYPES) {
      expect(screen.getByRole('checkbox', { name: `Toggle ${ENTITY_TYPE_META[type].label} entities` })).toBeInTheDocument()
    }

    const archivedToggle = screen.getByRole('checkbox', { name: 'Show archived entities' })
    expect(archivedToggle).not.toBeChecked()
    fireEvent.click(archivedToggle)
    expect(archivedToggle).toBeChecked()
  })
})
