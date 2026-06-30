import { useEffect, useState } from 'react'
import { useParams, useNavigate, Outlet, NavLink } from 'react-router-dom'
import { UniverseContext, type UniverseContextValue } from '../contexts/UniverseContext'
import { api } from '../lib/api'

// ponytail: re-fetches universe on mount instead of sharing universeStore;
// universeStore is already populated by DashboardPage, but self-contained layout handles direct nav
export default function UniverseLayout() {
  const { universeId } = useParams<{ universeId: string }>()
  const navigate = useNavigate()
  const [ctx, setCtx] = useState<UniverseContextValue>({ universe: null, works: [] })
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    if (!universeId) return
    setLoading(true)
    setError(null)
    Promise.all([api.getUniverse(universeId), api.listWorks(universeId)])
      .then(([{ universe }, { works }]) => {
        setCtx({ universe, works })
        setLoading(false)
      })
      .catch((err) => {
        setError((err as Error).message)
        setLoading(false)
      })
  }, [universeId])

  if (loading) {
    return (
      <div style={{ padding: 24 }}>
        <p style={{ color: '#888' }}>Loading universe…</p>
      </div>
    )
  }

  if (error) {
    return (
      <div style={{ padding: 24 }}>
        <p className="error">Failed to load universe: {error}</p>
        <button className="primary" onClick={() => navigate('/dashboard')}>
          Back to Dashboard
        </button>
      </div>
    )
  }

  const tabs = [
    { to: 'works', label: 'Works' },
    { to: 'graph', label: 'Graph' },
    { to: 'timeline', label: 'Timeline' },
    { to: 'contradictions', label: 'Contradictions' },
    { to: 'plot-holes', label: 'Plot-holes' },
  ]

  return (
    <UniverseContext.Provider value={ctx}>
      <div style={{ padding: 24 }}>
        <button
          onClick={() => navigate('/dashboard')}
          style={{ background: 'transparent', color: '#6c5ce7', marginBottom: 16 }}
        >
          ← Back to Dashboard
        </button>

        <h1>{ctx.universe?.name || 'Universe'}</h1>
        <p style={{ color: '#888', marginBottom: 16 }}>
          {ctx.universe?.genre} · {ctx.universe?.format}
        </p>

        {/* Tab bar */}
        <nav style={{ display: 'flex', gap: 4, marginBottom: 24, borderBottom: '1px solid #333', paddingBottom: 8 }}>
          {tabs.map((tab) => (
            <NavLink
              key={tab.to}
              to={tab.to}
              end
              style={({ isActive }) => ({
                padding: '6px 16px',
                borderRadius: '6px 6px 0 0',
                color: isActive ? '#6c5ce7' : '#888',
                background: isActive ? '#16213e' : 'transparent',
                textDecoration: 'none',
                fontSize: 14,
              })}
            >
              {tab.label}
            </NavLink>
          ))}
        </nav>

        <Outlet />
      </div>
    </UniverseContext.Provider>
  )
}
