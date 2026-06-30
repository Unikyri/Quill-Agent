import { describe, it, expect } from 'vitest'
import { render } from '@testing-library/react'
import { ReactFlowProvider, Position } from 'reactflow'
import CustomNode from '../CustomNode'

// ReactFlow nodes need a provider in tests
function renderNode(type: string, label: string) {
  return render(
    <ReactFlowProvider>
      <CustomNode
        id="test-1"
        type="custom"
        data={{ type, label }}
        xPos={0}
        yPos={0}
        zIndex={0}
        selected={false}
        dragging={false}
        isConnectable={true}
        sourcePosition={Position.Bottom}
        targetPosition={Position.Top}
      />
    </ReactFlowProvider>
  )
}

describe('CustomNode', () => {
  it('renders character node with purple border and person icon', () => {
    const { container } = renderNode('character', 'Alice')
    const node = container.firstElementChild as HTMLElement
    expect(node.style.borderColor.toLowerCase()).toBe('rgb(108, 92, 231)')
    expect(container.textContent).toContain('👤')
    expect(container.textContent).toContain('Alice')
  })

  it('renders location node with green border and pin icon', () => {
    const { container } = renderNode('location', 'Castle')
    const node = container.firstElementChild as HTMLElement
    expect(node.style.borderColor.toLowerCase()).toBe('rgb(0, 184, 148)')
    expect(container.textContent).toContain('📍')
    expect(container.textContent).toContain('Castle')
  })

  it('renders item node with yellow border and crystal icon', () => {
    const { container } = renderNode('item', 'Sword')
    const node = container.firstElementChild as HTMLElement
    expect(node.style.borderColor.toLowerCase()).toBe('rgb(253, 203, 110)')
    expect(container.textContent).toContain('🔮')
  })

  it('renders event node with red border and lightning icon', () => {
    const { container } = renderNode('event', 'Battle')
    const node = container.firstElementChild as HTMLElement
    expect(node.style.borderColor.toLowerCase()).toBe('rgb(225, 112, 85)')
    expect(container.textContent).toContain('⚡')
  })

  it('renders concept node with blue border and idea icon', () => {
    const { container } = renderNode('concept', 'Magic System')
    const node = container.firstElementChild as HTMLElement
    expect(node.style.borderColor.toLowerCase()).toBe('rgb(116, 185, 255)')
    expect(container.textContent).toContain('💡')
  })

  it('falls back to concept style for unknown type', () => {
    const { container } = renderNode('unknown_type', 'Mystery')
    const node = container.firstElementChild as HTMLElement
    expect(node.style.borderColor.toLowerCase()).toBe('rgb(116, 185, 255)')
    expect(container.textContent).toContain('💡')
  })

  it('shows "Untitled" when label is empty', () => {
    const { container } = renderNode('character', '')
    expect(container.textContent).toContain('Untitled')
  })
})
