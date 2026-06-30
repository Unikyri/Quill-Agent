import { describe, it, expect, vi, beforeEach } from 'vitest'
import { render, screen, waitFor } from '@testing-library/react'
import { MemoryRouter, Routes, Route } from 'react-router-dom'
import WorkPage from '../WorkPage'

const mockNavigate = vi.fn()
vi.mock('react-router-dom', async () => {
  const actual = await vi.importActual<typeof import('react-router-dom')>('react-router-dom')
  return {
    ...actual,
    useNavigate: () => mockNavigate,
  }
})

// Mock api
const mockGetWork = vi.fn()
const mockListChapters = vi.fn()

vi.mock('../../lib/api', () => ({
  api: {
    getWork: (...args: unknown[]) => mockGetWork(...args),
    listChapters: (...args: unknown[]) => mockListChapters(...args),
  },
}))

function renderPage(workId = 'work-123') {
  return render(
    <MemoryRouter initialEntries={[`/work/${workId}`]}>
      <Routes>
        <Route path="/work/:workId" element={<WorkPage />} />
      </Routes>
    </MemoryRouter>
  )
}

describe('WorkPage', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  it('shows loading state initially', () => {
    mockGetWork.mockReturnValue(new Promise(() => {})) // never resolves
    mockListChapters.mockReturnValue(new Promise(() => {}))
    renderPage()
    expect(screen.getByText('Loading…')).toBeInTheDocument()
  })

  it('renders work title and chapters on load', async () => {
    mockGetWork.mockResolvedValue({
      work: { id: 'work-123', title: 'My Novel', type: 'Novel' },
    })
    mockListChapters.mockResolvedValue({
      chapters: [
        { id: 'ch-1', title: 'Chapter 1', order_index: 1, word_count: 500 },
        { id: 'ch-2', title: 'Chapter 2', order_index: 2, word_count: 300 },
      ],
    })

    renderPage()

    await waitFor(() => {
      expect(screen.getByText('My Novel')).toBeInTheDocument()
    })

    expect(screen.getByText('Novel')).toBeInTheDocument()
    expect(screen.getByText('Chapter 1')).toBeInTheDocument()
    expect(screen.getByText('Chapter 2')).toBeInTheDocument()
    expect(screen.getByText('500 words')).toBeInTheDocument()
    expect(screen.getByText('300 words')).toBeInTheDocument()
  })

  it('shows error message on fetch failure', async () => {
    mockGetWork.mockRejectedValue(new Error('Network error'))
    mockListChapters.mockRejectedValue(new Error('Network error'))

    renderPage()

    await waitFor(() => {
      expect(screen.getByText(/Network error/)).toBeInTheDocument()
    })
  })

  it('shows empty state when no chapters', async () => {
    mockGetWork.mockResolvedValue({
      work: { id: 'work-123', title: 'Empty Work', type: 'Novel' },
    })
    mockListChapters.mockResolvedValue({ chapters: [] })

    renderPage()

    await waitFor(() => {
      expect(screen.getByText('No chapters yet.')).toBeInTheDocument()
    })
  })
})
