import { NavLink, Outlet } from 'react-router-dom'
import { useLogout } from '../../auth/useAuth'
import { Wordmark } from '@/components/Wordmark'
import { cn } from '@/lib/utils'

const NAV_ITEMS = [
{ to: '/admin/partners', label: 'Partners' },
{ to: '/admin/clients', label: 'Clients' },
{ to: '/admin/dashboard', label: 'Dashboard' },
]

export function AdminLayout() {
const logout = useLogout()

return (
    <div className="min-h-screen bg-surface-200 p-4">
    <div className="mx-auto flex max-w-sm flex-col gap-6">
        <h1 className="sr-only">Ticket Tout — admin dashboard</h1>
        <div className="flex items-center justify-between">
        <Wordmark />
        <button
            onClick={() => logout.mutate()}
            className="text-label-caps uppercase font-display text-ink-600"
        >
            Log out
        </button>
        </div>

        <nav className="flex gap-4">
        {NAV_ITEMS.map((item) => (
            <NavLink
            key={item.to}
            to={item.to}
            className={({ isActive }) =>
                cn(
                'text-label-caps uppercase font-display text-ink-600',
                isActive && 'text-brand-blue',
                )
            }
            >
            {item.label}
            </NavLink>
        ))}
        </nav>

        <Outlet />
    </div>
    </div>
)
}
