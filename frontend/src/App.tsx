import { Routes, Route, Navigate } from 'react-router-dom'
import { RoleGuard } from './routes/RoleGuard'
import { ClientHome } from './routes/client/ClientHome'
import { PartnerCatalog } from './routes/client/PartnerCatalog'
import { PartnerDetail } from './routes/client/PartnerDetail'
import { PartnerHome } from './routes/partner/PartnerHome'
import { PartnerDashboard } from './routes/partner/PartnerDashboard'
import { LoginPage } from './routes/LoginPage'
import { TermsOfUse } from './routes/TermsOfUse'
import { AccessibilityStatement } from './routes/AccessibilityStatement'
import { RootRedirect } from './routes/RootRedirect'
import { AdminLayout } from './routes/admin/AdminLayout'
import { AdminPartners } from './routes/admin/AdminPartners'
import { AdminClients } from './routes/admin/AdminClients'
import { AdminClientDetail } from './routes/admin/AdminClientDetail'

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
        path="/partner/dashboard"
        element={
          <RoleGuard role="partner">
            <PartnerDashboard />
          </RoleGuard>
        }
      />
      <Route
          path="/admin"
          element={
            <RoleGuard role="admin">
              <AdminLayout />
            </RoleGuard>
          }
        >
          <Route index element={<Navigate to="/admin/partners" replace />} />
          <Route path="partners" element={<AdminPartners />} />
          <Route path="clients" element={<AdminClients />} />
          <Route path="clients/:id" element={<AdminClientDetail />} />
        </Route>
      <Route path="/partners" element={<PartnerCatalog />} />
      <Route path="/partners/:id" element={<PartnerDetail />} />
      <Route path="/login"
      element={<LoginPage />}
      />
      <Route path="/terms" element={<TermsOfUse />} />
      <Route path="/accessibility" element={<AccessibilityStatement />} />
      <Route path="/" element={<RootRedirect />} />
    </Routes>
  )
}

export default App