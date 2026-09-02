import { useLogout } from '../../auth/useAuth'
import { useBalance } from '../../client/useClientData'
import { formatCents } from '../../lib/money'

export function ClientHome() {
  const logout = useLogout()
  const balance = useBalance()

  return (
    <div className="min-h-screen bg-surface-200 p-4">
      <div className="mx-auto flex max-w-sm flex-col gap-6">
        <div className="flex items-center justify-between">
          <p className="text-wordmark uppercase font-display text-graphite-900">
            CartePro
          </p>
          <button
            onClick={() => logout.mutate()}
            className="text-label-caps uppercase font-display text-ink-600"
          >
            Log out
          </button>
        </div>

        <div className="flex flex-col gap-2 rounded-card bg-card-gradient p-4 shadow-card-highlight">
          <p className="text-micro-caps uppercase font-sans text-graphite-400">
            Balance
          </p>
          {balance.isLoading && (
            <p className="text-body font-sans text-graphite-300">Loading...</p>
          )}
          {balance.isError && (
            <p className="text-body font-sans text-danger-600">
              {balance.error.message}
            </p>
          )}
          {balance.isSuccess && (
            <p className="text-balance-xl font-display text-graphite-100">
              {formatCents(balance.data.balance)}
            </p>
          )}
        </div>
      </div>
    </div>
  )
}