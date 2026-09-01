import { Routes, Route } from 'react-router-dom'
import { RoleGuard } from './routes/RoleGuard'
import { ClientHome } from './routes/client/ClientHome'
import { PartnerHome } from './routes/partner/PartnerHome'
import { AdminHome } from './routes/admin/AdminHome'

function App() {
  return (
    <Routes>
      <Route
        path="/client"
        element={
          <RoleGuard role="client">
            <ClientHome />
          </RoleGuard>
        }
      />
      <Route
        path="/partner"
        element={
          <RoleGuard role="partner">
            <PartnerHome />
          </RoleGuard>
        }
      />
      <Route
        path="/admin"
        element={
          <RoleGuard role="admin">
            <AdminHome />
          </RoleGuard>
        }
      />
    </Routes>
  )
}

export default App