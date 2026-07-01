import { useContext } from 'react'
import { useNavigate } from 'react-router-dom'
import { UniverseContext } from '../contexts/UniverseContext'
import styles from './UniverseWorksTab.module.css'

export default function UniverseWorksTab() {
  const { works } = useContext(UniverseContext)
  const navigate = useNavigate()

  return (
    <div className={styles.wrap}>
      <h2 className={styles.heading}>Works</h2>
      {works.length === 0 ? (
        <p className={styles.empty}>No works yet.</p>
      ) : (
        works.map((w) => (
          <div
            key={w.id}
            className={styles.card}
            onClick={() => navigate(`/work/${w.id}`)}
          >
            <h3 className={styles.cardTitle}>{w.title}</h3>
            <p className={styles.cardType}>{w.type}</p>
          </div>
        ))
      )}
    </div>
  )
}
