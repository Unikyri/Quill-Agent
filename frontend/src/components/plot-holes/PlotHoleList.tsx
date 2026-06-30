import { useMemo } from 'react'
import styles from './PlotHoleList.module.css'

export interface PlotHole {
  id: string
  description: string
  severity: string
}

interface PlotHoleListProps {
  plotHoles: PlotHole[]
}

const SEVERITY_ORDER: Record<string, number> = {
  critical: 0,
  high: 1,
  medium: 2,
  low: 3,
}

const SEVERITY_CLASS: Record<string, string> = {
  critical: styles.severityCritical,
  high: styles.severityHigh,
  medium: styles.severityMedium,
  low: styles.severityLow,
}

export default function PlotHoleList({ plotHoles }: PlotHoleListProps) {
  const sorted = useMemo(() => {
    return [...plotHoles].sort((a, b) => {
      const orderA = SEVERITY_ORDER[a.severity] ?? 99
      const orderB = SEVERITY_ORDER[b.severity] ?? 99
      return orderA - orderB
    })
  }, [plotHoles])

  return (
    <div className={styles.listWrap}>
      {sorted.map((ph) => (
        <div key={ph.id} className={styles.card}>
          <div className={styles.cardHeader}>
            <span className={`${styles.severity} ${SEVERITY_CLASS[ph.severity] || styles.severityLow}`}>
              {ph.severity}
            </span>
          </div>
          <p className={styles.cardDesc}>{ph.description}</p>
        </div>
      ))}
    </div>
  )
}
