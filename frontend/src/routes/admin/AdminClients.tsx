import { useState } from 'react'
import { Link } from 'react-router-dom'
import { useAdminClients } from '../../admin/useAdminClients'
import { formatCents } from '@/lib/money'
import { Button } from '@/components/ui/button'

export function AdminClients() {
  const [page, setPage] = useState(1)

  const clients = useAdminClients({ page })

  return (
    <div className="flex flex-col gap-4">
      <h2 className="text-h2 font-display text-brand-blue">Clients</h2>

      {clients.isLoading && <p className="text-body font-sans text-ink-600">Loading...</p>}
      {clients.isError && (
        <p className="text-body font-sans text-danger-600">{clients.error.message}</p>
      )}
      {clients.isSuccess && clients.data.data.length === 0 && (
        <p className="text-body font-sans text-ink-600">No clients found.</p>
      )}
      {clients.isSuccess && clients.data.data.length > 0 && (
        <div className="flex flex-col gap-[2px]">
          {clients.data.data.map((client) => (
            <Link
              key={client.ID}
              to={`/admin/clients/${client.ID}`}
              className="flex items-center justify-between rounded-row bg-surface-050 px-3 py-3 shadow-row-raised"
            >
              <p className="text-body-strong font-sans text-ink-900">{client.Name}</p>
              <p className="text-caption font-sans text-ink-600">{formatCents(client.Balance)}</p>
            </Link>
          ))}
        </div>
      )}

      {clients.isSuccess && (
        <div className="flex items-center justify-between">
          <Button
            size="sm"
            disabled={page <= 1}
            onClick={() => setPage((current) => current - 1)}
            className="rounded-panel bg-surface-300 text-button uppercase font-display text-graphite-900 shadow-key-secondary hover:bg-surface-300"
          >
            Prev
          </Button>
          <p className="text-caption font-sans text-ink-600">Page {page}</p>
          <Button
            size="sm"
            disabled={clients.data.data.length < clients.data.meta.limit}
            onClick={() => setPage((current) => current + 1)}
            className="rounded-panel bg-surface-300 text-button uppercase font-display text-graphite-900 shadow-key-secondary hover:bg-surface-300"
          >
            Next
          </Button>
        </div>
      )}
    </div>
  )
}
