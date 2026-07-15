import { create } from 'zustand'
import { api } from '../lib/api'
import { parseVertexRaw } from '../lib/graphParse'
import { ENTITY_TYPES, type EntityType } from '../lib/entityTypes'

interface GraphNode {
  id: string
  type: EntityType
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
  showArchived: boolean
  loading: boolean
  error: string | null
  requestVersion: number
  _universeId: string | null
  focalNodeId: string | null
  breadcrumb: string[]
  fetchGraph: (universeId: string) => Promise<void>
  refresh: () => Promise<void>
  focusNode: (id: string) => Promise<void>
  goBack: () => Promise<void>
  selectNode: (id: string | null) => void
  toggleFilter: (type: string) => void
  toggleArchived: () => void
}

const ALL_TYPES = ENTITY_TYPES as unknown as string[]

function isCurrentRequest(
  get: () => GraphState,
  requestVersion: number,
  universeId: string,
  focalNodeId?: string,
) {
  const state = get()
  return state.requestVersion === requestVersion
    && state._universeId === universeId
    && (focalNodeId === undefined || state.focalNodeId === focalNodeId)
}

// ponytail: shared node-mapping used by both the initial fetch and refresh,
// avoided duplicating the raw-agtype parsing logic a third time.
function mapNodes(rawNodes: any[], focalNodeId: string): GraphNode[] {
  const neighbors = rawNodes.filter((node: any) => (parseVertexRaw(String(node.properties?.raw || '')).entityId || node.id) !== focalNodeId)
  const total = neighbors.length || 1
  let neighborIndex = 0
  return rawNodes.map((n: any, i: number) => {
    const raw: string = n.properties?.raw || ''
    const v = parseVertexRaw(raw)
    const id = v.entityId || n.id || String(i)
    if (id === focalNodeId) {
      return {
        id,
        type: v.type as GraphNode['type'],
        position: { x: 0, y: 0 },
        data: { label: v.name, relevanceScore: v.relevanceScore, status: v.status },
      }
    }
    const angle = (2 * Math.PI * neighborIndex++) / total
    const radius = 220
    return {
      id,
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
  showArchived: false,
  loading: false,
  error: null,
  requestVersion: 0,
  _universeId: null,
  focalNodeId: null,
  breadcrumb: [],

  fetchGraph: async (universeId) => {
    const requestVersion = get().requestVersion + 1
    set({
      loading: true,
      error: null,
      requestVersion,
      _universeId: universeId,
      focalNodeId: null,
      selectedNodeId: null,
      breadcrumb: [],
    })
    try {
      let { entities } = await api.listEntities(universeId, { limit: '1', status: 'active' })
      if (!isCurrentRequest(get, requestVersion, universeId)) return
      if (entities.length === 0) {
        ({ entities } = await api.listEntities(universeId, { limit: '1', status: 'archived' }))
        if (!isCurrentRequest(get, requestVersion, universeId)) return
      }
      const focalNodeId = entities[0]?.id
      if (!focalNodeId) {
        if (isCurrentRequest(get, requestVersion, universeId)) {
          set({ nodes: [], edges: [], focalNodeId: null, loading: false })
        }
        return
      }
      const { nodes: rawNodes, edges: rawEdges } = await api.getEntityNeighbors(focalNodeId, universeId, 2)
      if (!isCurrentRequest(get, requestVersion, universeId)) return
      const nodes: GraphNode[] = mapNodes(rawNodes as any[], focalNodeId)
      const edges: GraphEdge[] = (rawEdges as any[]).map((e: any) => ({
        id: e.id || `${e.source}-${e.target}`,
        source: extractEdgeSource(e),
        target: extractEdgeTarget(e),
        label: e.type || e.label || '',
      }))
      set({ nodes, edges, focalNodeId, selectedNodeId: focalNodeId, loading: false })
    } catch (err) {
      if (isCurrentRequest(get, requestVersion, universeId)) {
        set({ error: (err as Error).message, loading: false })
      }
    }
  },

  refresh: async () => {
    const { _universeId, focalNodeId } = get()
    if (_universeId && focalNodeId) {
      const requestVersion = get().requestVersion + 1
      set({ requestVersion })
      try {
        const { nodes: rawNodes, edges: rawEdges } = await api.getEntityNeighbors(focalNodeId, _universeId, 2)
        if (!isCurrentRequest(get, requestVersion, _universeId, focalNodeId)) return
        const nodes: GraphNode[] = mapNodes(rawNodes as any[], focalNodeId)
        const edges: GraphEdge[] = (rawEdges as any[]).map((e: any) => ({
          id: e.id || `${e.source}-${e.target}`,
          source: extractEdgeSource(e),
          target: extractEdgeTarget(e),
          label: e.type || e.label || '',
        }))
        set({ nodes, edges, error: null })
      } catch (err) {
        if (isCurrentRequest(get, requestVersion, _universeId, focalNodeId)) {
          set({ error: (err as Error).message })
        }
      }
    }
  },

  focusNode: async (id) => {
    const { _universeId, focalNodeId, breadcrumb } = get()
    if (!_universeId || id === focalNodeId) return
    const requestVersion = get().requestVersion + 1
    set({ loading: true, error: null, requestVersion })
    try {
      const { nodes: rawNodes, edges: rawEdges } = await api.getEntityNeighbors(id, _universeId, 2)
      if (!isCurrentRequest(get, requestVersion, _universeId)) return
      const nodes = mapNodes(rawNodes as any[], id)
      const edges: GraphEdge[] = (rawEdges as any[]).map((e: any) => ({
        id: e.id || `${e.source}-${e.target}`,
        source: extractEdgeSource(e),
        target: extractEdgeTarget(e),
        label: e.type || e.label || '',
      }))
      set({ nodes, edges, focalNodeId: id, selectedNodeId: id, breadcrumb: focalNodeId ? [...breadcrumb, focalNodeId] : breadcrumb, loading: false })
    } catch (err) {
      if (isCurrentRequest(get, requestVersion, _universeId)) {
        set({ error: (err as Error).message, loading: false })
      }
    }
  },

  goBack: async () => {
    const { breadcrumb } = get()
    const previous = breadcrumb[breadcrumb.length - 1]
    if (!previous) return
    set({ breadcrumb: breadcrumb.slice(0, -1) })
    await get().focusNode(previous)
    set({ breadcrumb: breadcrumb.slice(0, -1) })
  },

  selectNode: (id) => set({ selectedNodeId: id }),

  toggleFilter: (type) => {
    const current = get().nodeFilter
    set({ nodeFilter: { ...current, [type]: !current[type] } })
  },

  toggleArchived: () => set((state) => ({ showArchived: !state.showArchived })),
}))
