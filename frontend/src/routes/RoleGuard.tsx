import type { ReactNode } from 'react'
import { Navigate } from 'react-router-dom'
import { useAuth } from '../auth/useAuth'
import type { Role } from '../types/auth'

interface RoleGuardProps {
  role: Role
  children: ReactNode
}

export function RoleGuard({ role, children }: RoleGuardProps) {
  const { status, user } = useAuth()

  if (status === 'loading') {
    return null
  }

  if (status === 'unauthenticated') {
    return <Navigate to="/login" replace />
  }

  if (user.role !== role) {
    return <Navigate to="/" replace />
  }

  return <>{children}</>
}