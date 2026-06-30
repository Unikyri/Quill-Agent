import { useContext } from 'react'
import { useNavigate } from 'react-router-dom'
import { UniverseContext } from '../contexts/UniverseContext'

export default function UniverseWorksTab() {
  const { works } = useContext(UniverseContext)
  const navigate = useNavigate()

  return (
    <div>
      <h2 style={{ marginBottom: 16 }}>Works</h2>
      {works.length === 0 ? (
        <div className="card"><p>No works yet.</p></div>
      ) : (
        works.map((w) => (
          <div key={w.id} className="card" style={{ cursor: 'pointer' }} onClick={() => navigate(`/work/${w.id}`)}>
            <h3>{w.title}</h3>
            <p style={{ color: '#888' }}>{w.type}</p>
          </div>
        ))
      )}
    </div>
  )
}
