export const ENTITY_TYPES = [
  'character',
  'place',
  'object',
  'faction',
  'event',
  'world_rule',
  'plot_arc',
] as const

export type EntityType = (typeof ENTITY_TYPES)[number]

export const ENTITY_TYPE_META: Record<EntityType, { label: string; color: string; glyph: string }> = {
  character: { color: 'var(--node-character)', glyph: '●', label: 'Character' },
  place: { color: 'var(--node-place)', glyph: '◆', label: 'Place' },
  object: { color: 'var(--gold)', glyph: '◇', label: 'Object' },
  faction: { color: 'var(--node-faction)', glyph: '■', label: 'Faction' },
  event: { color: 'var(--node-event)', glyph: '▲', label: 'Event' },
  world_rule: { color: 'var(--node-worldrule)', glyph: '◈', label: 'World Rule' },
  plot_arc: { color: 'var(--node-plotarc)', glyph: '◉', label: 'Plot Arc' },
}
