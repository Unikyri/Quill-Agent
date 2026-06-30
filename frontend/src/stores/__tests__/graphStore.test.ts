import { describe, it, expect, vi, beforeEach } from 'vitest'
import { useGraphStore } from '../graphStore'

// Mock api
const mockGetGraph = vi.fn()
vi.mock('../../lib/api', () => ({
  api: {
    getGraph: (...args: unknown[]) => mockGetGraph(...args),
  },
}))

function getStore() {
  return useGraphStore.getState()
}

const mockNodes = [
  { id: 'n1', type: 'character', position: { x: 0, y: 0 }, data: { label: 'Alice' } },
  { id: 'n2', type: 'location', position: { x: 100, y: 0 }, data: { label: 'Castle' } },
  { id: 'n3', type: 'concept', position: { x: 200, y: 0 }, data: { label: 'Magic' } },
]

const mockEdges = [
  { id: 'e1', source: 'n1', target: 'n2', label: 'lives in' },
]

beforeEach(() => {
  vi.clearAllMocks()
  useGraphStore.setState({
    nodes: [],
    edges: [],
    selectedNodeId: null,
    nodeFilter: { character: true, location: true, item: true, event: true, concept: true },
    loading: false,
    error: null,
    _universeId: null,
  })
})

describe('graphStore', () => {
  describe('initial state', () => {
    it('has empty nodes and edges', () => {
      expect(getStore().nodes).toEqual([])
      expect(getStore().edges).toEqual([])
    })

    it('has all type filters enabled', () => {
      const f = getStore().nodeFilter
      expect(f.character).toBe(true)
      expect(f.location).toBe(true)
      expect(f.item).toBe(true)
      expect(f.event).toBe(true)
      expect(f.concept).toBe(true)
    })

    it('has null selectedNodeId', () => {
      expect(getStore().selectedNodeId).toBeNull()
    })

    it('is not loading and has no error', () => {
      expect(getStore().loading).toBe(false)
      expect(getStore().error).toBeNull()
    })
  })

  describe('fetchGraph', () => {
    it('sets loading true and populates nodes/edges on success', async () => {
      mockGetGraph.mockResolvedValue({ nodes: mockNodes, edges: mockEdges })

      const promise = getStore().fetchGraph('uni-1')
      expect(getStore().loading).toBe(true)

      await promise
      expect(getStore().nodes).toEqual(mockNodes)
      expect(getStore().edges).toEqual(mockEdges)
      expect(getStore().loading).toBe(false)
      expect(getStore().error).toBeNull()
      expect(getStore()._universeId).toBe('uni-1')
    })

    it('sets error on failure', async () => {
      mockGetGraph.mockRejectedValue(new Error('Network error'))

      await getStore().fetchGraph('uni-1')
      expect(getStore().loading).toBe(false)
      expect(getStore().error).toBe('Network error')
      expect(getStore().nodes).toEqual([])
    })
  })

  describe('refresh', () => {
    it('refetches using stored universeId', async () => {
      mockGetGraph.mockResolvedValueOnce({ nodes: mockNodes, edges: mockEdges })
      await getStore().fetchGraph('uni-1')
      vi.clearAllMocks()

      const updatedNodes = [{ ...mockNodes[0], data: { label: 'Alice Updated' } }]
      mockGetGraph.mockResolvedValueOnce({ nodes: updatedNodes, edges: [] })

      await getStore().refresh()
      expect(mockGetGraph).toHaveBeenCalledWith('uni-1')
      expect(getStore().nodes).toEqual(updatedNodes)
      expect(getStore().edges).toEqual([])
    })

    it('does nothing if no universeId was set', async () => {
      await getStore().refresh()
      expect(mockGetGraph).not.toHaveBeenCalled()
    })
  })

  describe('selectNode', () => {
    it('sets selectedNodeId', () => {
      getStore().selectNode('n1')
      expect(getStore().selectedNodeId).toBe('n1')
    })

    it('clears selectedNodeId with null', () => {
      getStore().selectNode('n1')
      getStore().selectNode(null)
      expect(getStore().selectedNodeId).toBeNull()
    })
  })

  describe('toggleFilter', () => {
    it('toggles a single type filter off', () => {
      expect(getStore().nodeFilter.character).toBe(true)
      getStore().toggleFilter('character')
      expect(getStore().nodeFilter.character).toBe(false)
    })

    it('toggles back on', () => {
      getStore().toggleFilter('character') // off
      getStore().toggleFilter('character') // on
      expect(getStore().nodeFilter.character).toBe(true)
    })

    it('does not affect other filters', () => {
      getStore().toggleFilter('character')
      expect(getStore().nodeFilter.location).toBe(true)
      expect(getStore().nodeFilter.concept).toBe(true)
    })
  })
})
