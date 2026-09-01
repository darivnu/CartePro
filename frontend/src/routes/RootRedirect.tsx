import { Navigate } from 'react-router-dom'
import { useAuth } from '../auth/useAuth'
import type { Role } from '../types/auth'

const roleHomePaths: Record<Role, string> = {
client: '/client',
partner: '/partner',
admin: '/admin',
}

export function RootRedirect() {
const { status, user } = useAuth()

if (status === 'loading') {
    return null
}

if (status === 'unauthenticated') {
    return <Navigate to="/login" replace />
}

return <Navigate to={roleHomePaths[user.role]} replace />
}
