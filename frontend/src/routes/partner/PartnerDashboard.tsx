import { useState } from 'react'
import { Link } from 'react-router-dom'
import { useLogout } from '../../auth/useAuth'
import { Wordmark } from '@/components/Wordmark'
import { useOwnDashboard } from '../../partners/usePartnersData'
import { formatCents } from '../../lib/money'
import { Input } from '@/components/ui/input'

function toRfc3339Start(date: string) {
  return date ? `${date}T00:00:00Z` : undefined
}

function toRfc3339End(date: string) {
  return date ? `${date}T23:59:59Z` : undefined
}

function formatDayLabel(isoDate: string) {
  const [year, month, day] = isoDate.split('-')
  return `${day}/${month}/${year}`
}

export function PartnerDashboard() {
  const logout = useLogout()
  const [fromDate, setFromDate] = useState('')
  const [toDate, setToDate] = useState('')

  const dashboard = useOwnDashboard({
    from: toRfc3339Start(fromDate),
    to: toRfc3339End(toDate),
  })

  return (
    <div className="min-h-screen bg-surface-200 p-4">
      <div className="mx-auto flex max-w-sm flex-col gap-6">
        <h1 className="sr-only">CartePro — partner financial dashboard</h1>
        <div className="flex items-center justify-between">
          <Wordmark />
          <button
            onClick={() => logout.mutate()}
            className="text-label-caps uppercase font-display text-ink-600"
          >
            Log out
          </button>
        </div>

        <Link to="/partner" className="text-label-caps uppercase font-display text-brand-blue">
          Back
        </Link>

        <div className="flex gap-2">
          <Input
            type="date"
            value={fromDate}
            onChange={(e) => setFromDate(e.target.value)}
            aria-label="From date"
          />
          <Input
            type="date"
            value={toDate}
            onChange={(e) => setToDate(e.target.value)}
            aria-label="To date"
          />
        </div>

        {dashboard.isLoading && (
          <p className="text-body font-sans text-ink-600">Loading…</p>
        )}

        {dashboard.isError && (
          <p className="text-body font-sans text-danger-600">{dashboard.error.message}</p>
        )}

        {dashboard.isSuccess && (
          <>
            <div className="flex flex-col gap-2 rounded-card bg-card-gradient p-4 shadow-card-highlight">
              <p className="text-micro-caps uppercase font-sans text-graphite-400">
                Total received
              </p>
              <p className="text-balance-xl font-display text-graphite-100">
                {formatCents(dashboard.data.total_received)}
              </p>
              <p className="text-caption font-sans text-graphite-400">
                {dashboard.data.transaction_count} transactions
              </p>
            </div>

            <p className="text-h2 font-display text-graphite-900">By day</p>

            {dashboard.data.by_day.length === 0 && (
              <p className="text-body font-sans text-ink-600">No transactions in this range.</p>
            )}

            {dashboard.data.by_day.length > 0 && (
              <div className="flex flex-col gap-[2px]">
                {dashboard.data.by_day.map((bucket) => (
                  <div
                    key={bucket.date}
                    className="flex items-center justify-between rounded-row bg-surface-050 px-3 py-3 shadow-row-raised"
                  >
                    <div className="flex flex-col gap-1">
                      <p className="text-body-strong font-sans text-ink-900">
                        {formatDayLabel(bucket.date)}
                      </p>
                      <p className="text-caption font-sans text-ink-600">
                        {bucket.transaction_count} transactions
                      </p>
                    </div>
                    <p className="text-body-strong font-sans text-success-600">
                      {formatCents(bucket.total_received)}
                    </p>
                  </div>
                ))}
              </div>
            )}
          </>
        )}
      </div>
    </div>
  )
}
