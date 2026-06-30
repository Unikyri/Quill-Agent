import { describe, it, expect, vi } from 'vitest'
import { render, screen } from '@testing-library/react'
import Toolbar from '../Toolbar'

// Mock TipTap editor interface
function createMockEditor(activeMarks: string[] = []) {
  const chain = {
    focus: vi.fn().mockReturnThis(),
    toggleBold: vi.fn().mockReturnThis(),
    toggleItalic: vi.fn().mockReturnThis(),
    toggleUnderline: vi.fn().mockReturnThis(),
    toggleHighlight: vi.fn().mockReturnThis(),
    extendMarkRange: vi.fn().mockReturnThis(),
    setLink: vi.fn().mockReturnThis(),
    unsetLink: vi.fn().mockReturnThis(),
    run: vi.fn(),
  }

  return {
    isActive: (mark: string) => activeMarks.includes(mark),
    getAttributes: vi.fn().mockReturnValue({}),
    chain: () => chain,
    state: { selection: { from: 0, to: 0 }, doc: { resolve: vi.fn() } },
  }
}

describe('Toolbar', () => {
  it('renders nothing when editor is null', () => {
    const { container } = render(<Toolbar editor={null} />)
    expect(container.innerHTML).toBe('')
  })

  it('renders formatting buttons', () => {
    const mockEditor = createMockEditor()
    render(<Toolbar editor={mockEditor as any} />)

    expect(screen.getByTitle('Bold (Ctrl+B)')).toBeInTheDocument()
    expect(screen.getByTitle('Italic (Ctrl+I)')).toBeInTheDocument()
    expect(screen.getByTitle('Underline (Ctrl+U)')).toBeInTheDocument()
    expect(screen.getByTitle('Highlight')).toBeInTheDocument()
    expect(screen.getByTitle('Add link')).toBeInTheDocument()
  })

  it('shows active state for bold when editor.isActive("bold") is true', () => {
    const mockEditor = createMockEditor(['bold'])
    render(<Toolbar editor={mockEditor as any} />)

    const boldBtn = screen.getByTitle('Bold (Ctrl+B)')
    expect(boldBtn.className).toContain('buttonActive')
  })

  it('does not show active state when mark is inactive', () => {
    const mockEditor = createMockEditor([])
    render(<Toolbar editor={mockEditor as any} />)

    const boldBtn = screen.getByTitle('Bold (Ctrl+B)')
    expect(boldBtn.className).not.toContain('buttonActive')
  })
})
