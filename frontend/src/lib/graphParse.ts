// ponytail: AGE returns agtype vertices/edges as JSON-ish text (double-quoted
// keys/values); a regex scrape is enough for read-only display, no full
// agtype/JSON parser needed.
function extractStr(raw: string, key: string): string | undefined {
  const m = raw.match(new RegExp(`"${key}"\\s*:\\s*"([^"]*)"`))
  return m?.[1]
}

function extractNum(raw: string, key: string): number | undefined {
  const m = raw.match(new RegExp(`"${key}"\\s*:\\s*(-?[\\d.]+)`))
  return m ? Number(m[1]) : undefined
}

export interface ParsedVertex {
  entityId: string
  name: string
  type: string
  status?: string
  relevanceScore?: number
}

// Entity types the extraction pipeline actually produces (qwen_service.go prompt
// + entity_service.go CreateNode label) — keep in sync with backend.
export const ENTITY_TYPES = ['character', 'place', 'event', 'faction', 'world_rule', 'plot_arc'] as const

export function parseVertexRaw(raw: string): ParsedVertex {
  const entityId = extractStr(raw, 'entity_id') || ''
  return {
    entityId,
    name: extractStr(raw, 'name') || entityId.slice(0, 8),
    type: extractStr(raw, 'label') || 'character',
    status: extractStr(raw, 'status'),
    relevanceScore: extractNum(raw, 'relevance_score'),
  }
}
