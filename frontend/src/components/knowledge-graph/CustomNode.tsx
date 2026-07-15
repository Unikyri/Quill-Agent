import { memo, type CSSProperties } from 'react'
import { Handle, Position, type NodeProps } from 'reactflow'
import styles from './GraphCanvas.module.css'
import { ENTITY_TYPE_META } from '../../lib/entityTypes'

function CustomNode({ data }: NodeProps) {
  const nodeType = (data.type as string) || 'character'
  const meta = ENTITY_TYPE_META[nodeType as keyof typeof ENTITY_TYPE_META] || ENTITY_TYPE_META.character
  const rawRelevance = data.relevanceScore
  const relevance = typeof rawRelevance === 'number'
    ? Math.min(1, Math.max(0, rawRelevance))
    : 0.5
  const opacity = 0.55 + relevance * 0.45
  const scale = 0.84 + relevance * 0.28

  return (
    <div
      className={styles.customNode}
      style={{ borderColor: meta.color, opacity, '--node-scale': scale } as CSSProperties}
      data-relevance={relevance.toFixed(2)}
    >
      <Handle type="target" position={Position.Top} className={styles.handle} />
      <div className={styles.nodeContent}>
        <span className={`${styles.nodeIcon} glyph`}>{meta.glyph}</span>
        <span className={styles.nodeLabel}>{(data.label as string) || 'Untitled'}</span>
      </div>
      <Handle type="source" position={Position.Bottom} className={styles.handle} />
    </div>
  )
}

export default memo(CustomNode)
