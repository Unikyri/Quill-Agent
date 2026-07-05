import { create } from 'zustand'
import { api } from '../lib/api'
import { parseVertexRaw, ENTITY_TYPES } from '../lib/graphParse'

interface GraphNode {
  id: string
  type: (typeof ENTITY_TYPES)[number]
  position: { x: number; y: number }
  data: {
    label: string
    description?: string
    relevanceScore?: number
    status?: string
    [key: string]: unknown
  }
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

const ALL_TYPES = ENTITY_TYPES as unknown as string[]

// ponytail: shared node-mapping used by both the initial fetch and refresh,
// avoided duplicating the raw-agtype parsing logic a third time.
function mapNodes(rawNodes: any[]): GraphNode[] {
  const total = rawNodes.length || 1
  return rawNodes.map((n: any, i: number) => {
    const angle = (2 * Math.PI * i) / total
    const radius = Math.max(100, total * 30)
    const raw: string = n.properties?.raw || ''
    const v = parseVertexRaw(raw)
    return {
      id: v.entityId || n.id || String(i),
      type: v.type as GraphNode['type'],
      position: { x: Math.cos(angle) * radius, y: Math.sin(angle) * radius },
      data: { label: v.name, relevanceScore: v.relevanceScore, status: v.status },
    }
  })
}

// ponytail: extract source/target from raw AGE edge strings like
// [:KNOWS {source: 'id1', target: 'id2'}]
function extractEdgeSource(e: any): string {
  const raw: string = e.properties?.raw || e.id || ''
  const m = raw.match(/source:\s*'([^']*)'/)
  return m?.[1] || e.source || ''
}
function extractEdgeTarget(e: any): string {
  const raw: string = e.properties?.raw || e.id || ''
  const m = raw.match(/target:\s*'([^']*)'/)
  return m?.[1] || e.target || ''
}

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
      const { nodes: rawNodes, edges: rawEdges } = await api.getGraph(universeId)
      // Transform backend {id, labels, properties} → frontend {id, type, position, data}
      // ponytail: auto-layout with circle packing; no layout lib needed for hackathon
      const nodes: GraphNode[] = mapNodes(rawNodes as any[])
      const edges: GraphEdge[] = (rawEdges as any[]).map((e: any) => ({
        id: e.id || `${e.source}-${e.target}`,
        source: extractEdgeSource(e),
        target: extractEdgeTarget(e),
        label: e.type || e.label || '',
      }))
      set({ nodes, edges, loading: false })
    } catch (err) {
      set({ error: (err as Error).message, loading: false })
    }
  },

  refresh: async () => {
    const { _universeId } = get()
    if (_universeId) {
      try {
        const { nodes: rawNodes, edges: rawEdges } = await api.getGraph(_universeId)
        const nodes: GraphNode[] = mapNodes(rawNodes as any[])
        const edges: GraphEdge[] = (rawEdges as any[]).map((e: any) => ({
          id: e.id || `${e.source}-${e.target}`,
          source: extractEdgeSource(e),
          target: extractEdgeTarget(e),
          label: e.type || e.label || '',
        }))
        set({ nodes, edges, error: null })
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
