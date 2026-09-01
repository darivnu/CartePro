import { Routes, Route } from 'react-router-dom'
import { RoleGuard } from './routes/RoleGuard'
import { ClientHome } from './routes/client/ClientHome'
import { PartnerHome } from './routes/partner/PartnerHome'
import { AdminHome } from './routes/admin/AdminHome'
import { LoginPage } from './routes/LoginPage'
import { RootRedirect } from './routes/RootRedirect'

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
      <Route path="/login"
      element={<LoginPage />}
      />
      <Route path="/" element={<RootRedirect />} />
    </Routes>
  )
}

export default App