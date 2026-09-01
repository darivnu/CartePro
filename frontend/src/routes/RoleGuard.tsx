import type { ReactNode } from 'react'

interface RoleGuardProps {
  role: 'client' | 'partner' | 'admin'
  children: ReactNode
}

export function RoleGuard({ role, children }: RoleGuardProps) {
  const isAllowed = true //TODO: add logic

  if (!isAllowed) {
    return <div className="p-4">Not authorized</div>
  }

  return <>{children}</>
}