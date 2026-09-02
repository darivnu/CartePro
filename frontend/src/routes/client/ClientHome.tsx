import { useLogout } from '../../auth/useAuth'
import { useBalance, useTransactions } from '../../client/useClientData'
import { formatCents } from '../../lib/money'

const dateFormatter = new Intl.DateTimeFormat('fr-FR', {
  day: 'numeric',
  month: 'short',
  hour: '2-digit',
  minute: '2-digit',
})

export function ClientHome() {
  const logout = useLogout()
  const balance = useBalance()
  const transactions = useTransactions()

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

        <div className="flex flex-col gap-3">
          <h2 className="text-h2 font-display text-graphite-900">Transactions</h2>

          {transactions.isLoading && (
            <p className="text-body font-sans text-ink-600">Loading...</p>
          )}
          {transactions.isError && (
            <p className="text-body font-sans text-danger-600">
              {transactions.error.message}
            </p>
          )}
          {transactions.isSuccess && transactions.data.data.length === 0 && (
            <p className="text-body font-sans text-ink-600">No transactions yet.</p>
          )}
          {transactions.isSuccess && (
            <div className="flex flex-col gap-[2px]">
              {transactions.data.data.map((tx) => (
                <div
                  key={tx.id}
                  className="flex items-center justify-between rounded-row bg-surface-050 px-3 py-3 shadow-row-raised"
                >
                  <div className="flex flex-col gap-1">
                    <p className="text-body-strong font-sans text-ink-900">
                      {tx.partner ? tx.partner.business_name : 'Top-up'}
                    </p>
                    <p className="text-caption font-sans text-ink-600">
                      {dateFormatter.format(new Date(tx.created_at))}
                    </p>
                  </div>
                  <p className="text-amount-row font-display text-ink-900">
                    {tx.type === 'debit' ? '-' : '+'}
                    {formatCents(tx.amount)}
                  </p>
                </div>
              ))}
            </div>
          )}
        </div>
      </div>
    </div>
  )
}
