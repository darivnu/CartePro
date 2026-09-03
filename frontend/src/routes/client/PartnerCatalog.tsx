import { useState } from 'react'
import { Link } from 'react-router-dom'
import { usePartners } from '../../partners/usePartnersData'
import { Input } from '@/components/ui/input'
import { Button } from '@/components/ui/button'

export function PartnerCatalog() {
  const [search, setSearch] = useState('')
  const [category, setCategory] = useState('')
  const [page, setPage] = useState(1)

  const partners = usePartners({ search, category, page })

  function handleSearchChange(event: React.ChangeEvent<HTMLInputElement>) {
    setSearch(event.target.value)
    setPage(1)
  }

  function handleCategoryChange(event: React.ChangeEvent<HTMLInputElement>) {
    setCategory(event.target.value)
    setPage(1)
  }

  return (
    <div className="min-h-screen bg-surface-200 p-4">
      <div className="mx-auto flex max-w-sm flex-col gap-6">
        <div className="flex items-center justify-between">
          <p className="text-wordmark uppercase font-display text-graphite-900">
            CartePro
          </p>
          <Link
            to="/client"
            className="text-label-caps uppercase font-display text-ink-600"
          >
            Back
          </Link>
        </div>

        <h1 className="text-h1 font-display text-graphite-900">Partners</h1>

        <div className="flex flex-col gap-2">
          <Input
            value={search}
            onChange={handleSearchChange}
            placeholder="Search partners"
          />
          <Input
            value={category}
            onChange={handleCategoryChange}
            placeholder="Category"
          />
        </div>

        {partners.isLoading && (
          <p className="text-body font-sans text-ink-600">Loading...</p>
        )}
        {partners.isError && (
          <p className="text-body font-sans text-danger-600">
            {partners.error.message}
          </p>
        )}
        {partners.isSuccess && partners.data.data.length === 0 && (
          <p className="text-body font-sans text-ink-600">No partners found.</p>
        )}
        {partners.isSuccess && partners.data.data.length > 0 && (
          <div className="flex flex-col gap-[2px]">
            {partners.data.data.map((partner) => (
              <Link
                key={partner.id}
                to={`/client/partners/${partner.id}`}
                className="flex items-center justify-between rounded-row bg-surface-050 px-3 py-3 shadow-row-raised"
              >
                <div className="flex flex-col gap-1">
                  <p className="text-body-strong font-sans text-ink-900">
                    {partner.business_name}
                  </p>
                  <p className="text-caption font-sans text-ink-600">
                    {partner.category} · {partner.region}
                  </p>
                </div>
              </Link>
            ))}
          </div>
        )}

        {partners.isSuccess && partners.data.meta.total_pages > 1 && (
          <div className="flex items-center justify-between">
            <Button
              size="sm"
              disabled={page <= 1}
              onClick={() => setPage((current) => current - 1)}
              className="rounded-panel bg-surface-300 text-button uppercase font-display text-graphite-900 shadow-key-secondary hover:bg-surface-300"
            >
              Prev
            </Button>
            <p className="text-caption font-sans text-ink-600">
              Page {page} of {partners.data.meta.total_pages}
            </p>
            <Button
              size="sm"
              disabled={page >= partners.data.meta.total_pages}
              onClick={() => setPage((current) => current + 1)}
              className="rounded-panel bg-surface-300 text-button uppercase font-display text-graphite-900 shadow-key-secondary hover:bg-surface-300"
            >
              Next
            </Button>
          </div>
        )}
      </div>
    </div>
  )
}
