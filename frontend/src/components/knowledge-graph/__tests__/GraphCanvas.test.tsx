import { describe, expect, it, beforeEach, vi } from 'vitest'
import { render, screen } from '@testing-library/react'
import GraphCanvas from '../GraphCanvas'
import { useGraphStore } from '../../../stores/graphStore'

vi.mock('reactflow', () => ({
  default: ({ nodes }: { nodes: Array<{ id: string }> }) => (
    <div data-testid="graph-flow" data-node-ids={nodes.map((node) => node.id).join(',')} />
  ),
  Background: () => null,
  Controls: () => null,
  MiniMap: () => null,
  Handle: () => null,
  Position: { Top: 'top', Bottom: 'bottom' },
}))

beforeEach(() => {
  useGraphStore.setState({
    nodes: [
      { id: 'active', type: 'character', position: { x: 0, y: 0 }, data: { label: 'Active', status: 'active' } },
      { id: 'archived', type: 'object', position: { x: 120, y: 0 }, data: { label: 'Archived', status: 'archived' } },
    ],
    edges: [],
    nodeFilter: { character: true, place: true, object: true, faction: true, event: true, world_rule: true, plot_arc: true },
    showArchived: false,
  })
})

describe('GraphCanvas', () => {
  it('hides archived nodes until the archived toggle is enabled', () => {
    const { rerender } = render(<GraphCanvas />)
    expect(screen.getByTestId('graph-flow')).toHaveAttribute('data-node-ids', 'active')

    useGraphStore.setState({ showArchived: true })
    rerender(<GraphCanvas />)

    expect(screen.getByTestId('graph-flow')).toHaveAttribute('data-node-ids', 'active,archived')
  })
})
