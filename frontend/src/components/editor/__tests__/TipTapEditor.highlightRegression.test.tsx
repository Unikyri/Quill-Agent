// Regression coverage for the risk flagged in the frontend-greenfield-rebuild
// design: adding a selectionUpdate/transaction listener to TipTapEditor must
// not disturb the EntityHighlight/CandidateHighlight ProseMirror decorations.
// TipTapEditor.test.tsx mocks '@tiptap/react' entirely (so it never exercises
// real ProseMirror decorations); this file uses the REAL tiptap stack in
// jsdom so the highlight plugins actually run.
import { describe, it, expect, vi, beforeEach } from 'vitest'
import { render, waitFor, act } from '@testing-library/react'
import TipTapEditor from '../TipTapEditor'

// eslint-disable-next-line @typescript-eslint/no-explicit-any
let capturedEditor: any = null

vi.mock('../../../stores/wsStore', () => ({
  useWSStore: vi.fn((selector: (s: unknown) => unknown) => {
    const state = { send: vi.fn(), submissions: {} }
    return selector ? selector(state) : state
  }),
}))

// Pass-through wrapper around the REAL @tiptap/react so the test can grab the
// live Editor instance and drive it directly (setTextSelection, insertContent).
vi.mock('@tiptap/react', async () => {
  const actual = await vi.importActual<typeof import('@tiptap/react')>('@tiptap/react')
  return {
    ...actual,
    useEditor: (options: Parameters<typeof actual.useEditor>[0]) => {
      const editor = actual.useEditor(options)
      capturedEditor = editor
      return editor
    },
  }
})

describe('TipTapEditor — highlight decorations survive the selection listener (real TipTap)', () => {
  beforeEach(() => {
    capturedEditor = null
  })

  const defaultProps = {
    chapterId: 'ch-1',
    workId: 'w-1',
    universeId: 'u-1',
    initialContent: '<p>Ann met the Captain today.</p>',
    knownEntities: [{ id: 'e-1', name: 'Ann', type: 'character' }],
    candidateEntities: [{ id: 'c-1', name: 'Captain', type: 'character' }],
  }

  it('renders entity and candidate highlight decorations on mount', async () => {
    const { container } = render(<TipTapEditor {...defaultProps} />)

    await waitFor(() => expect(capturedEditor).toBeTruthy())
    await waitFor(() => {
      expect(container.querySelector('.entity-highlight')).toBeTruthy()
      expect(container.querySelector('.candidate-highlight')).toBeTruthy()
    })
  })

  it('keeps highlight decorations intact across a selection-only transaction (no doc change)', async () => {
    const { container } = render(<TipTapEditor {...defaultProps} />)
    await waitFor(() => expect(capturedEditor).toBeTruthy())
    await waitFor(() => expect(container.querySelector('.entity-highlight')).toBeTruthy())

    // Selection-only change — exercises the new selectionUpdate/transaction
    // listener added to TipTapEditor without touching the document.
    act(() => {
      capturedEditor.commands.setTextSelection({ from: 1, to: 4 })
    })

    expect(container.querySelector('.entity-highlight')).toBeTruthy()
    expect(container.querySelector('.candidate-highlight')).toBeTruthy()
  })

  it('updates highlight decorations correctly after a real document change', async () => {
    const { container } = render(<TipTapEditor {...defaultProps} />)
    await waitFor(() => expect(capturedEditor).toBeTruthy())
    await waitFor(() => expect(container.querySelector('.entity-highlight')).toBeTruthy())

    act(() => {
      capturedEditor.commands.setContent('<p>No known names here.</p>')
    })

    await waitFor(() => {
      expect(container.querySelector('.entity-highlight')).toBeFalsy()
      expect(container.querySelector('.candidate-highlight')).toBeFalsy()
    })

    act(() => {
      capturedEditor.commands.setContent('<p>Ann is back.</p>')
    })

    await waitFor(() => {
      expect(container.querySelector('.entity-highlight')).toBeTruthy()
    })
  })

  it('does not throw or loop when many selection changes fire in a row', async () => {
    const consoleErrorSpy = vi.spyOn(console, 'error').mockImplementation(() => {})
    render(<TipTapEditor {...defaultProps} />)
    await waitFor(() => expect(capturedEditor).toBeTruthy())

    expect(() => {
      act(() => {
        for (let i = 0; i < 20; i += 1) {
          capturedEditor.commands.setTextSelection(i % 2 === 0 ? { from: 1, to: 4 } : 1)
        }
      })
    }).not.toThrow()

    const maxDepthErrors = consoleErrorSpy.mock.calls.filter((call) =>
      call.some((arg) => typeof arg === 'string' && arg.includes('Maximum update depth'))
    )
    expect(maxDepthErrors).toHaveLength(0)
    consoleErrorSpy.mockRestore()
  })
})
