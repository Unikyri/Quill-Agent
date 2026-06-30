import { create } from 'zustand'
import { api } from '../lib/api'

interface GraphNode {
  id: string
  type: 'character' | 'location' | 'item' | 'event' | 'concept'
  position: { x: number; y: number }
  data: { label: string; description?: string; [key: string]: unknown }
}

interface GraphEdge {
  id: string
  source: string
  target: string
  label: string
}

// ponytail: global store for graph state; a single store is cheaper than 5 useState slices per page
interface GraphState {
  nodes: GraphNode[]
  edges: GraphEdge[]
  selectedNodeId: string | null
  nodeFilter: Record<string, boolean> // { character: true, location: true, ... }
  loading: boolean
  error: string | null
  _universeId: string | null
  fetchGraph: (universeId: string) => Promise<void>
  refresh: () => Promise<void>
  selectNode: (id: string | null) => void
  toggleFilter: (type: string) => void
}

const ALL_TYPES = ['character', 'location', 'item', 'event', 'concept']

export const useGraphStore = create<GraphState>((set, get) => ({
  nodes: [],
  edges: [],
  selectedNodeId: null,
  nodeFilter: Object.fromEntries(ALL_TYPES.map((t) => [t, true])),
  loading: false,
  error: null,
  _universeId: null,

  fetchGraph: async (universeId) => {
    set({ loading: true, error: null, _universeId: universeId })
    try {
      const { nodes, edges } = await api.getGraph(universeId)
      set({ nodes: nodes as GraphNode[], edges, loading: false })
    } catch (err) {
      set({ error: (err as Error).message, loading: false })
    }
  },

  refresh: async () => {
    const { _universeId } = get()
    if (_universeId) {
      try {
        const { nodes, edges } = await api.getGraph(_universeId)
        set({ nodes: nodes as GraphNode[], edges, error: null })
      } catch (err) {
        set({ error: (err as Error).message })
      }
    }
  },

  selectNode: (id) => set({ selectedNodeId: id }),

  toggleFilter: (type) => {
    const current = get().nodeFilter
    set({ nodeFilter: { ...current, [type]: !current[type] } })
  },
}))
