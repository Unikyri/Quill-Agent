import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'
import { render } from '@testing-library/react'
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

vi.mock('@tiptap/react', () => {
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  const MockEditorContent = ({ editor }: any) => (
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    <div data-testid="editor-content">{editor?.getText?.() ?? ''}</div>
  )

  return {
    useEditor: (options: any) => {
      capturedOnUpdate = options.onUpdate ?? null
      return createMockEditor(options.content ?? '')
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
  return {
    state: {
      selection: { from: 5 },
      doc: {
        resolve: () => ({
          parent: { isBlock: true, textContent: text },
          depth: 1,
          before: () => 0,
        }),
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

    const selection = { from: 5 }
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
