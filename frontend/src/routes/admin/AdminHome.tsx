import { useLogout } from '../../auth/useAuth'
import { Wordmark } from '@/components/Wordmark'

export function AdminHome() {
  const logout = useLogout()

  return (
    <div className="min-h-screen bg-surface-200 p-4">
      <div className="mx-auto flex max-w-sm flex-col gap-6">
        <div className="flex items-center justify-between">
          <Wordmark />
          <button
            onClick={() => logout.mutate()}
            className="text-label-caps uppercase font-display text-ink-600"
          >
            Log out
          </button>
        </div>

        <p className="text-body font-sans text-ink-600">Admin area</p>
      </div>
    </div>
  )
}
