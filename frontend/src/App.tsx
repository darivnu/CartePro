import { Routes, Route } from 'react-router-dom'
import { RoleGuard } from './routes/RoleGuard'
import { EmployeeHome } from './routes/employee/EmployeeHome'
import { PartnerHome } from './routes/partner/PartnerHome'
import { AdminHome } from './routes/admin/AdminHome'

function App() {
  return (
    <Routes>
      <Route
        path="/employee"
        element={
          <RoleGuard role="employee">
            <EmployeeHome />
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