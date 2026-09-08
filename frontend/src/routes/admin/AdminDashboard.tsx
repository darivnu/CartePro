import { useState } from 'react'
import { useAdminDashboard } from '../../admin/useAdminDashboard'
import { formatCents } from '@/lib/money'
import { Input } from '@/components/ui/input'
import type { AdminPartner } from '../../types/admin'

function toRfc3339Start(date: string) {
  return date ? `${date}T00:00:00Z` : undefined
}

function toRfc3339End(date: string) {
  return date ? `${date}T23:59:59Z` : undefined
}

function countByRegion(partners: AdminPartner[]) {
  const counts = new Map<string, number>()
  for (const partner of partners) {
    counts.set(partner.Region, (counts.get(partner.Region) ?? 0) + 1)
  }
  return Array.from(counts.entries()).sort((a, b) => b[1] - a[1])
}

export function AdminDashboard() {
  const [fromDate, setFromDate] = useState('')
  const [toDate, setToDate] = useState('')

  const dashboard = useAdminDashboard({
    from: toRfc3339Start(fromDate),
    to: toRfc3339End(toDate),
  })

  return (
    <div className="flex flex-col gap-4">
      <h2 className="text-h2 font-display text-brand-blue">Dashboard</h2>

      <div className="flex gap-2">
        <Input
          type="date"
          value={fromDate}
          onChange={(event) => setFromDate(event.target.value)}
          aria-label="From date"
        />
        <Input
          type="date"
          value={toDate}
          onChange={(event) => setToDate(event.target.value)}
          aria-label="To date"
        />
      </div>

      {dashboard.isLoading && <p className="text-body font-sans text-ink-600">Loading...</p>}
      {dashboard.isError && (
        <p className="text-body font-sans text-danger-600">{dashboard.error.message}</p>
      )}
      {dashboard.isSuccess && (
        <>
          <div className="flex flex-col gap-[2px]">
            <div className="flex items-center justify-between rounded-row bg-surface-050 px-3 py-3 shadow-row-raised">
              <p className="text-body-strong font-sans text-ink-900">Transaction volume</p>
              <p className="text-body-strong font-sans text-success-600">
                {formatCents(dashboard.data.total_transaction_volume)}
              </p>
            </div>
            <div className="flex items-center justify-between rounded-row bg-surface-050 px-3 py-3 shadow-row-raised">
              <p className="text-body-strong font-sans text-ink-900">Total clients</p>
              <p className="text-body-strong font-sans text-ink-900">{dashboard.data.total_clients}</p>
            </div>
            <div className="flex items-center justify-between rounded-row bg-surface-050 px-3 py-3 shadow-row-raised">
              <p className="text-body-strong font-sans text-ink-900">Active partners</p>
              <p className="text-body-strong font-sans text-ink-900">
                {dashboard.data.active_partners.length}
              </p>
            </div>
          </div>

          <h3 className="text-label-caps uppercase font-display text-ink-600">By region</h3>

          {dashboard.data.active_partners.length === 0 && (
            <p className="text-body font-sans text-ink-600">No active partners yet.</p>
          )}
          {dashboard.data.active_partners.length > 0 && (
            <div className="flex flex-col gap-[2px]">
              {countByRegion(dashboard.data.active_partners).map(([region, count]) => (
                <div
                  key={region}
                  className="flex items-center justify-between rounded-row bg-surface-050 px-3 py-3 shadow-row-raised"
                >
                  <p className="text-body-strong font-sans text-ink-900">{region}</p>
                  <p className="text-caption font-sans text-ink-600">{count} partners</p>
                </div>
              ))}
            </div>
          )}
        </>
      )}
    </div>
  )
}
