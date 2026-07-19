import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'
import { render, act, fireEvent } from '@testing-library/react'
import TipTapEditor from '../TipTapEditor'

// ── Mocks ───────────────────────────────────────────────────────────────

const mockWSStoreSend = vi.fn()
let mockSubmissions: Record<string, unknown> = {}

vi.mock('../../../stores/wsStore', () => ({
  useWSStore: vi.fn((selector: (s: unknown) => unknown) => {
    const state = { send: mockWSStoreSend, submissions: mockSubmissions }
    return selector ? selector(state) : state
  }),
}))

// Capture the onUpdate callback passed to useEditor so we can fire it
let capturedOnUpdate: ((props: { editor: any }) => void) | null = null

// The real useEditor() constructs one Editor instance for the component's
// mounted lifetime (see TipTapEditor's own content-sync effect, which mutates
// the existing instance via editor.commands.setContent rather than expecting
// a new one). Mirror that here so a test can grab the instance once and drive
// its selectionUpdate/transaction listeners directly.
let sharedMockEditor: any = null // eslint-disable-line @typescript-eslint/no-explicit-any

vi.mock('@tiptap/react', () => {
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  const MockEditorContent = ({ editor }: any) => (
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    <div data-testid="editor-content">{editor?.getText?.() ?? ''}</div>
  )

  return {
    useEditor: (options: any) => {
      capturedOnUpdate = options.onUpdate ?? null
      if (!sharedMockEditor) sharedMockEditor = createMockEditor(options.content ?? '')
      return sharedMockEditor
    },
    EditorContent: MockEditorContent,
  }
})

vi.mock('@tiptap/starter-kit', () => ({ default: {} }))
vi.mock('@tiptap/extension-placeholder', () => ({ default: { configure: () => ({}) } }))
vi.mock('@tiptap/extension-highlight', () => ({ default: {} }))
vi.mock('@tiptap/extension-underline', () => ({ default: {} }))
vi.mock('@tiptap/extension-link', () => ({ default: { configure: () => ({}) } }))

// ══════════════════════════════════════════════════════════════════════════
// Helpers
// ══════════════════════════════════════════════════════════════════════════

function createMockEditor(text: string) {
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  const listeners: Record<string, Array<() => void>> = {}
  return {
    state: {
      selection: { from: 5, to: 9, empty: true },
      doc: {
        resolve: () => ({
          parent: { isBlock: true, textContent: text },
          depth: 1,
          before: () => 0,
        }),
        textBetween: () => text,
      },
    },
    getHTML: () => `<p>${text}</p>`,
    getText: () => text,
    isActive: () => false,
    getAttributes: () => ({}),
    chain: () => ({
      focus: () => ({
        toggleBold: () => ({ run: vi.fn() }),
        toggleItalic: () => ({ run: vi.fn() }),
        toggleUnderline: () => ({ run: vi.fn() }),
        toggleHighlight: () => ({ run: vi.fn() }),
        extendMarkRange: () => ({ setLink: () => ({ run: vi.fn() }), unsetLink: () => ({ run: vi.fn() }) }),
      }),
    }),
    commands: { setContent: vi.fn() },
    // Minimal TipTap-style event emitter, enough to exercise the
    // selectionUpdate/transaction subscription added to TipTapEditor.
    on(event: string, handler: () => void) {
      listeners[event] = listeners[event] || []
      listeners[event].push(handler)
    },
    off(event: string, handler: () => void) {
      listeners[event] = (listeners[event] || []).filter((existing) => existing !== handler)
    },
    __emit(event: string) {
      (listeners[event] || []).forEach((handler) => handler())
    },
  }
}

function fireOnUpdate(text: string) {
  if (!capturedOnUpdate) throw new Error('onUpdate not captured — did useEditor fire?')
  capturedOnUpdate({ editor: createMockEditor(text) })
}

// ══════════════════════════════════════════════════════════════════════════
// Tests
// ══════════════════════════════════════════════════════════════════════════

describe('TipTapEditor — paragraph submit on idle', () => {
  beforeEach(() => {
    vi.useFakeTimers()
    mockWSStoreSend.mockClear()
    mockSubmissions = {}
    capturedOnUpdate = null
    sharedMockEditor = null
  })

  afterEach(() => {
    vi.useRealTimers()
  })

  const defaultProps = {
    chapterId: 'ch-1',
    workId: 'w-1',
    universeId: 'u-1',
  }

  it('renders the editor and toolbar', () => {
    const { getByTestId, getByTitle } = render(<TipTapEditor {...defaultProps} />)
    expect(getByTestId('editor-content')).toBeInTheDocument()
    expect(getByTitle('Bold (⌘B)')).toBeInTheDocument()
  })

  it('does NOT send paragraph_submit before 5000ms idle', () => {
    render(<TipTapEditor {...defaultProps} />)

    fireOnUpdate('Hello world')
    vi.advanceTimersByTime(2500)

    expect(mockWSStoreSend).not.toHaveBeenCalled()
  })

  it('sends paragraph_submit with correct payload after 5000ms idle', () => {
    render(<TipTapEditor {...defaultProps} />)

    fireOnUpdate('Hello world')
    vi.advanceTimersByTime(5000)

    expect(mockWSStoreSend).toHaveBeenCalledTimes(1)
    expect(mockWSStoreSend).toHaveBeenCalledWith(expect.objectContaining({
      type: 'paragraph_submit',
      payload: expect.objectContaining({
        submission_id: expect.any(String),
        paragraph_ref: 'ch-1:5',
        work_id: 'w-1',
        chapter_id: 'ch-1',
        universe_id: 'u-1',
        text: 'Hello world',
      }),
    }))
  })

  it('resets the timer on subsequent keystrokes (does not double-fire)', () => {
    render(<TipTapEditor {...defaultProps} />)

    fireOnUpdate('Hello')
    vi.advanceTimersByTime(2500)

    // New keystroke resets timer
    fireOnUpdate('Hello world')
    vi.advanceTimersByTime(2500)

    expect(mockWSStoreSend).not.toHaveBeenCalled()

    vi.advanceTimersByTime(2500)
    expect(mockWSStoreSend).toHaveBeenCalledTimes(1)
    expect(mockWSStoreSend).toHaveBeenCalledWith(
      expect.objectContaining({
        type: 'paragraph_submit',
        payload: expect.objectContaining({ text: 'Hello world' }),
      })
    )
  })

  it('does not resend the same paragraph text twice', () => {
    render(<TipTapEditor {...defaultProps} />)

    fireOnUpdate('Same text')
    vi.advanceTimersByTime(5000)
    expect(mockWSStoreSend).toHaveBeenCalledTimes(1)

    // Fire again with same text, advance timer
    fireOnUpdate('Same text')
    vi.advanceTimersByTime(5000)
    expect(mockWSStoreSend).toHaveBeenCalledTimes(1)
  })

  it('sends again when paragraph text changes', () => {
    render(<TipTapEditor {...defaultProps} />)

    fireOnUpdate('First paragraph')
    vi.advanceTimersByTime(5000)
    expect(mockWSStoreSend).toHaveBeenCalledTimes(1)

    fireOnUpdate('Second paragraph')
    vi.advanceTimersByTime(5000)
    expect(mockWSStoreSend).toHaveBeenCalledTimes(2)
    expect(mockWSStoreSend).toHaveBeenLastCalledWith(
      expect.objectContaining({
        payload: expect.objectContaining({ text: 'Second paragraph' }),
      })
    )
  })

  it('submits the paragraph captured at edit time after the cursor moves', () => {
    render(<TipTapEditor {...defaultProps} />)

    const selection = { from: 5, to: 5, empty: true }
    const editor = createMockEditor('Paragraph A')
    editor.state.selection = selection
    ;(editor.state.doc.resolve as unknown as (position: number) => unknown) = (position: number) => ({
      parent: { isBlock: true, textContent: position < 10 ? 'Paragraph A' : 'Paragraph B' },
      depth: 1,
      start: () => position,
      before: () => 0,
    })
    if (!capturedOnUpdate) throw new Error('onUpdate not captured')
    capturedOnUpdate({ editor })
    // Moving into B after the transaction must not change the captured A.
    selection.from = 20
    vi.advanceTimersByTime(5000)

    expect(mockWSStoreSend).toHaveBeenCalledWith(expect.objectContaining({
      payload: expect.objectContaining({ text: 'Paragraph A', paragraph_ref: 'ch-1:5' }),
    }))
  })

  it('renders a lifecycle status for every submission in the current paragraph set', () => {
    mockSubmissions = {
      'submission-a': { submissionId: 'submission-a', paragraphRef: 'ch-1:5', chapterId: 'ch-1', phase: 'analyzing', updatedAt: 2 },
      'submission-b': { submissionId: 'submission-b', paragraphRef: 'ch-1:20', chapterId: 'ch-1', phase: 'failed', reason: 'Qwen unavailable', updatedAt: 1 },
      'other-chapter': { submissionId: 'other-chapter', paragraphRef: 'ch-2:5', chapterId: 'ch-2', phase: 'done', updatedAt: 3 },
    }
    const { getAllByTestId } = render(<TipTapEditor {...defaultProps} />)
    const statuses = getAllByTestId('analysis-submission-status')
    expect(statuses).toHaveLength(2)
    expect(statuses.map((status) => status.getAttribute('data-paragraph-ref'))).toEqual(['ch-1:5', 'ch-1:20'])
  })
})

describe('TipTapEditor — craft-review trigger', () => {
  beforeEach(() => {
    vi.useFakeTimers()
    mockWSStoreSend.mockClear()
    mockSubmissions = {}
    capturedOnUpdate = null
    sharedMockEditor = null
  })

  afterEach(() => {
    vi.useRealTimers()
  })

  const defaultProps = {
    chapterId: 'ch-1',
    workId: 'w-1',
    universeId: 'u-1',
  }

  it('does not render the craft-review control when onCraftReview is not provided', () => {
    const { queryByText } = render(<TipTapEditor {...defaultProps} />)
    expect(queryByText('Ask for a craft review →')).not.toBeInTheDocument()
  })

  it('shows the full-width labeled button, disabled, with a persistent caption when there is no selection', () => {
    const onCraftReview = vi.fn()
    const { getByText } = render(<TipTapEditor {...defaultProps} onCraftReview={onCraftReview} />)

    const button = getByText('Ask for a craft review →')
    expect(button).toBeDisabled()
    // Persistent caption, not a hover-only tooltip: it must be in the DOM already.
    expect(getByText('(select a passage first)')).toBeInTheDocument()
  })

  it('reactively enables the button and hides the caption once a selection is made (selectionUpdate)', () => {
    const onCraftReview = vi.fn()
    const { getByText, queryByText } = render(<TipTapEditor {...defaultProps} onCraftReview={onCraftReview} />)

    expect(getByText('Ask for a craft review →')).toBeDisabled()

    act(() => {
      sharedMockEditor.state.selection.empty = false
      sharedMockEditor.__emit('selectionUpdate')
    })

    expect(getByText('Ask for a craft review →')).not.toBeDisabled()
    expect(queryByText('(select a passage first)')).not.toBeInTheDocument()
  })

  it('reactively enables the button on a transaction event too (programmatic doc changes)', () => {
    const onCraftReview = vi.fn()
    const { getByText } = render(<TipTapEditor {...defaultProps} onCraftReview={onCraftReview} />)

    act(() => {
      sharedMockEditor.state.selection.empty = false
      sharedMockEditor.__emit('transaction')
    })

    expect(getByText('Ask for a craft review →')).not.toBeDisabled()
  })

  it('re-disables the button and re-shows the caption once the selection is cleared again', () => {
    const onCraftReview = vi.fn()
    const { getByText } = render(<TipTapEditor {...defaultProps} onCraftReview={onCraftReview} />)

    act(() => {
      sharedMockEditor.state.selection.empty = false
      sharedMockEditor.__emit('selectionUpdate')
    })
    expect(getByText('Ask for a craft review →')).not.toBeDisabled()

    act(() => {
      sharedMockEditor.state.selection.empty = true
      sharedMockEditor.__emit('selectionUpdate')
    })
    expect(getByText('Ask for a craft review →')).toBeDisabled()
    expect(getByText('(select a passage first)')).toBeInTheDocument()
  })

  it('shows the reviewing label and stays disabled while a review is in flight, even with a selection', () => {
    const onCraftReview = vi.fn()
    const { getByText, queryByText } = render(
      <TipTapEditor {...defaultProps} onCraftReview={onCraftReview} reviewing />
    )

    act(() => {
      sharedMockEditor.state.selection.empty = false
      sharedMockEditor.__emit('selectionUpdate')
    })

    const button = getByText('Reviewing your passage…')
    expect(button).toBeDisabled()
    // No caption while reviewing — the selection isn't the reason it's disabled.
    expect(queryByText('(select a passage first)')).not.toBeInTheDocument()
  })

  it('sends the selected passage unchanged when the enabled button is clicked', () => {
    const onCraftReview = vi.fn()
    const { getByText } = render(
      <TipTapEditor {...defaultProps} initialContent="Hello world" onCraftReview={onCraftReview} />
    )

    act(() => {
      sharedMockEditor.state.selection.empty = false
      sharedMockEditor.state.selection.from = 5
      sharedMockEditor.state.selection.to = 9
      sharedMockEditor.__emit('selectionUpdate')
    })

    fireEvent.click(getByText('Ask for a craft review →'))

    expect(onCraftReview).toHaveBeenCalledWith({ passage: 'Hello world', from: 5, to: 9 })
  })

  it('does not enter an infinite render loop from the selection listener firing repeatedly', () => {
    const consoleErrorSpy = vi.spyOn(console, 'error').mockImplementation(() => {})
    render(<TipTapEditor {...defaultProps} onCraftReview={vi.fn()} />)

    expect(() => {
      act(() => {
        for (let i = 0; i < 50; i += 1) {
          sharedMockEditor.state.selection.empty = i % 2 === 0
          sharedMockEditor.__emit('transaction')
        }
      })
    }).not.toThrow()

    const maxDepthErrors = consoleErrorSpy.mock.calls.filter((call) =>
      call.some((arg) => typeof arg === 'string' && arg.includes('Maximum update depth'))
    )
    expect(maxDepthErrors).toHaveLength(0)
    consoleErrorSpy.mockRestore()
  })

  it('does not affect the existing formatting toolbar', () => {
    const { getByTitle } = render(<TipTapEditor {...defaultProps} onCraftReview={vi.fn()} />)
    expect(getByTitle('Bold (⌘B)')).toBeInTheDocument()
    expect(getByTitle('Heading 1')).toBeInTheDocument()
  })
})
