import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom'
import { useAuthStore } from './stores/authStore'
import LoginPage from './pages/LoginPage'
import DashboardPage from './pages/DashboardPage'
import UniverseLayout from './pages/UniverseLayout'
import UniverseWorksTab from './pages/UniverseWorksTab'
import KnowledgeGraphPage from './pages/KnowledgeGraphPage'
import TimelinePage from './pages/TimelinePage'
import ContradictionsPage from './pages/ContradictionsPage'
import PlotHolesPage from './pages/PlotHolesPage'
import EditorPage from './pages/EditorPage'
import WorkPage from './pages/WorkPage'

function ProtectedRoute({ children }: { children: React.ReactNode }) {
  const isAuthenticated = useAuthStore((s) => s.isAuthenticated)
  return isAuthenticated ? <>{children}</> : <Navigate to="/login" />
}

export default function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/login" element={<LoginPage />} />
        <Route path="/dashboard" element={<ProtectedRoute><DashboardPage /></ProtectedRoute>} />
        <Route path="/universe/:universeId" element={<ProtectedRoute><UniverseLayout /></ProtectedRoute>}>
          <Route index element={<Navigate to="works" replace />} />
          <Route path="works" element={<UniverseWorksTab />} />
          <Route path="graph" element={<KnowledgeGraphPage />} />
          <Route path="timeline" element={<TimelinePage />} />
          <Route path="contradictions" element={<ContradictionsPage />} />
          <Route path="plot-holes" element={<PlotHolesPage />} />
        </Route>
        <Route path="/work/:workId" element={<ProtectedRoute><WorkPage /></ProtectedRoute>} />
        <Route path="/editor/:chapterId" element={<ProtectedRoute><EditorPage /></ProtectedRoute>} />
        <Route path="*" element={<Navigate to="/dashboard" />} />
      </Routes>
    </BrowserRouter>
  )
}
