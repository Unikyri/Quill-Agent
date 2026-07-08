import { describe, it, expect, beforeEach, vi } from 'vitest'
import { render, screen, waitFor, fireEvent } from '@testing-library/react'
import FusionExplorer from '../FusionExplorer'
import { api } from '../../../lib/api'

vi.mock('../../../lib/api', () => ({
  api: {
    recallExplain: vi.fn(),
  },
}))

const recallExplain = api.recallExplain as ReturnType<typeof vi.fn>

beforeEach(() => {
  vi.clearAllMocks()
})

function renderAndSearch(query = 'dragon') {
  render(<FusionExplorer universeId="u1" />)
  fireEvent.change(screen.getByRole('textbox'), { target: { value: query } })
  fireEvent.click(screen.getByRole('button', { name: /explain/i }))
}

describe('FusionExplorer', () => {
  it('renders one column per RRF pipeline, the fused list, and per-item contributions', async () => {
    recallExplain.mockResolvedValue({
      query: 'dragon',
      pipeline_sizes: { vector: 2, graph: 1, recency: 1, keyword: 0, consolidated: 0 },
      items: [
        {
          id: 'i1',
          entity_id: 'e1',
          fact: 'The dragon guards the tower',
          rrf_score: 0.9,
          contributions: [
            { pipeline: 'vector', rank: 1, delta: 0.5 },
            { pipeline: 'graph', rank: 2, delta: 0.4 },
          ],
          fit_in_budget: true,
        },
        {
          id: 'i2',
          entity_id: 'e2',
          fact: 'A lone knight rides north',
          rrf_score: 0.3,
          contributions: [{ pipeline: 'recency', rank: 1, delta: 0.3 }],
          fit_in_budget: false,
        },
      ],
      budget: {
        max_context_tokens: 1000,
        available: 400,
        entities_tokens: 200,
        vector_tokens: 300,
        tools_tokens: 100,
        used_percent: 60,
      },
    })

    renderAndSearch()

    await waitFor(() => expect(screen.getByTestId('fused-item-i1')).toBeInTheDocument())

    // one column per pipeline
    for (const pipeline of ['vector', 'graph', 'recency', 'keyword', 'consolidated']) {
      expect(screen.getByTestId(`pipeline-column-${pipeline}`)).toBeInTheDocument()
    }

    // fused result list
    expect(screen.getByTestId('fused-item-i1')).toBeInTheDocument()
    expect(screen.getByTestId('fused-item-i2')).toBeInTheDocument()
    expect(screen.getByText(/the dragon guards the tower/i)).toBeInTheDocument()

    // per-item contribution breakdown
    expect(screen.getByTestId('contribution-i1-vector')).toHaveTextContent(/vector/i)
    expect(screen.getByTestId('contribution-i1-vector')).toHaveTextContent('1')
    expect(screen.getByTestId('contribution-i1-graph')).toBeInTheDocument()

    // fit-in-budget indicator
    expect(screen.getByTestId('fit-in-budget-i1')).toHaveTextContent(/fit/i)
    expect(screen.getByTestId('fit-in-budget-i2')).toHaveTextContent(/dropped|not/i)
  })

  it('renders an empty state without crashing when recallExplain returns zero items', async () => {
    recallExplain.mockResolvedValue({
      query: 'dragon',
      pipeline_sizes: { vector: 0, graph: 0, recency: 0, keyword: 0, consolidated: 0 },
      items: [],
      budget: {
        max_context_tokens: 1000,
        available: 1000,
        entities_tokens: 0,
        vector_tokens: 0,
        tools_tokens: 0,
        used_percent: 0,
      },
    })

    renderAndSearch()

    await waitFor(() => expect(screen.getByText(/no results/i)).toBeInTheDocument())
  })

  it('renders a single contribution for an item with only one pipeline hit', async () => {
    recallExplain.mockResolvedValue({
      query: 'dragon',
      pipeline_sizes: { vector: 1, graph: 0, recency: 0, keyword: 0, consolidated: 0 },
      items: [
        {
          id: 'i3',
          entity_id: 'e3',
          fact: 'Only vector found this',
          rrf_score: 0.2,
          contributions: [{ pipeline: 'vector', rank: 1, delta: 0.2 }],
          fit_in_budget: true,
        },
      ],
      budget: {
        max_context_tokens: 1000,
        available: 900,
        entities_tokens: 50,
        vector_tokens: 50,
        tools_tokens: 0,
        used_percent: 10,
      },
    })

    renderAndSearch()

    await waitFor(() => expect(screen.getByTestId('fused-item-i3')).toBeInTheDocument())
    expect(screen.getByTestId('contribution-i3-vector')).toBeInTheDocument()
    expect(screen.queryByTestId('contribution-i3-graph')).not.toBeInTheDocument()
  })

  it('shows a loading state while the request is in flight and an error state on failure', async () => {
    let reject!: (err: Error) => void
    recallExplain.mockReturnValue(new Promise((_resolve, rej) => { reject = rej }))

    renderAndSearch()

    expect(screen.getByText(/loading|explaining/i)).toBeInTheDocument()

    reject(new Error('boom'))

    await waitFor(() => expect(screen.getByText(/error|failed/i)).toBeInTheDocument())
  })
})
