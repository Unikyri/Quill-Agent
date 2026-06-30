import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'
import { api } from '../api'

const API_BASE = '/api/v1'

describe('api', () => {
  let fetchMock: ReturnType<typeof vi.fn>

  beforeEach(() => {
    fetchMock = vi.fn()
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    globalThis.fetch = fetchMock as any
  })

  afterEach(() => {
    vi.restoreAllMocks()
  })

  function mockFetchResponse(body: unknown, status = 200) {
    fetchMock.mockResolvedValueOnce({
      ok: status >= 200 && status < 300,
      status,
      json: () => Promise.resolve(body),
    })
  }

  // ── New Phase 2a methods ──────────────────────────────────────────────

  describe('getContradictions', () => {
    it('calls GET /universes/:id/contradictions and returns typed shape', async () => {
      const contradictions = [
        { id: 'c1', message: 'Plot hole', severity: 'high', entities: ['e1'] },
      ]
      mockFetchResponse({ contradictions })

      const result = await api.getContradictions('uni-1')

      expect(fetchMock).toHaveBeenCalledTimes(1)
      const [url] = fetchMock.mock.calls[0]
      expect(url).toBe(`${API_BASE}/universes/uni-1/contradictions`)
      expect(result.contradictions).toEqual(contradictions)
    })
  })

  describe('getTimeline', () => {
    it('calls GET /universes/:id/timeline and returns typed shape', async () => {
      const events = [
        { id: 'e1', label: 'Birth', timestamp: '2024-01-01', description: 'Born' },
      ]
      mockFetchResponse({ events })

      const result = await api.getTimeline('uni-1')

      expect(fetchMock).toHaveBeenCalledTimes(1)
      const [url] = fetchMock.mock.calls[0]
      expect(url).toBe(`${API_BASE}/universes/uni-1/timeline`)
      expect(result.events).toEqual(events)
    })
  })

  describe('getPlotHoles', () => {
    it('calls GET /universes/:id/plot-holes and returns typed shape', async () => {
      const plot_holes = [
        { id: 'p1', description: 'Missing link', severity: 'medium' },
      ]
      mockFetchResponse({ plot_holes })

      const result = await api.getPlotHoles('uni-1')

      expect(fetchMock).toHaveBeenCalledTimes(1)
      const [url] = fetchMock.mock.calls[0]
      expect(url).toBe(`${API_BASE}/universes/uni-1/plot-holes`)
      expect(result.plot_holes).toEqual(plot_holes)
    })
  })

  describe('getGraph', () => {
    it('calls GET /universes/:id/graph and returns ReactFlow-shaped nodes and edges', async () => {
      const graph = {
        nodes: [{ id: 'n1', type: 'character', position: { x: 0, y: 0 }, data: { name: 'Alice' } }],
        edges: [{ id: 'e1', source: 'n1', target: 'n2', label: 'knows' }],
      }
      mockFetchResponse(graph)

      const result = await api.getGraph('uni-1')

      expect(fetchMock).toHaveBeenCalledTimes(1)
      const [url] = fetchMock.mock.calls[0]
      expect(url).toBe(`${API_BASE}/universes/uni-1/graph`)
      expect(result.nodes).toEqual(graph.nodes)
      expect(result.edges).toEqual(graph.edges)
    })
  })

  describe('recall', () => {
    it('calls POST /universes/:id/recall with query and k in body, returns typed shape', async () => {
      const items = [
        { id: 'r1', fact: 'Alice is a wizard', score: 0.95, entity_id: 'e1' },
      ]
      mockFetchResponse({ items })

      const result = await api.recall('uni-1', 'wizard origin', 5)

      expect(fetchMock).toHaveBeenCalledTimes(1)
      const [url, init] = fetchMock.mock.calls[0]
      expect(url).toBe(`${API_BASE}/universes/uni-1/recall`)
      expect(init.method).toBe('POST')
      expect(JSON.parse(init.body)).toEqual({ query: 'wizard origin', k: 5 })
      expect(result.items).toEqual(items)
    })
  })
})
