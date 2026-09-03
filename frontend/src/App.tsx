import { Routes, Route } from 'react-router-dom'
import { RoleGuard } from './routes/RoleGuard'
import { ClientHome } from './routes/client/ClientHome'
import { PartnerCatalog } from './routes/client/PartnerCatalog'
import { PartnerDetail } from './routes/client/PartnerDetail'
import { PartnerHome } from './routes/partner/PartnerHome'
import { AdminHome } from './routes/admin/AdminHome'
import { LoginPage } from './routes/LoginPage'
import { TermsOfUse } from './routes/TermsOfUse'
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
        path="/client/partners"
        element={
          <RoleGuard role="client">
            <PartnerCatalog />
          </RoleGuard>
        }
      />
      <Route
        path="/client/partners/:id"
        element={
          <RoleGuard role="client">
            <PartnerDetail />
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
      <Route path="/terms" element={<TermsOfUse />} />
      <Route path="/" element={<RootRedirect />} />
    </Routes>
  )
}

export default App