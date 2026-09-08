import { Link, useLocation, useParams } from 'react-router-dom'
import { usePartner } from '../../partners/usePartnersData'
import { Wordmark } from '@/components/Wordmark'

export function PartnerDetail() {
  const { id } = useParams()
  const partner = usePartner(id ?? '')
  const location = useLocation()
  const partnersBasePath = location.pathname.startsWith('/client')
    ? '/client/partners'
    : '/partners'

  return (
    <div className="min-h-screen bg-surface-200 p-4">
      <div className="mx-auto flex max-w-sm flex-col gap-6">
        <div className="flex items-center justify-between">
          <Wordmark />
          <Link
            to={partnersBasePath}
            className="text-label-caps uppercase font-display text-ink-600"
          >
            Back
          </Link>
        </div>

        {partner.isLoading && (
          <p className="text-body font-sans text-ink-600">Loading...</p>
        )}
        {partner.isError && (
          <p className="text-body font-sans text-danger-600">
            {partner.error.message}
          </p>
        )}
        {partner.isSuccess && (
          <div className="flex flex-col gap-6">
            <div className="flex flex-col gap-1">
              <h1 className="text-h1 font-display text-ink-900">
                {partner.data.partner.BusinessName}
              </h1>
              {partner.data.partner.MinisterPick && (
                <p className="text-caption uppercase font-display text-gold-500">
                  ⭐ Minister Pick
                </p>
              )}
            </div>

            <div className="flex flex-col gap-4">
              <div className="flex flex-col gap-1">
                <p className="text-label-caps uppercase font-display text-ink-600">
                  Category
                </p>
                <p className="text-body font-sans text-ink-900">
                  {partner.data.partner.Category}
                </p>
              </div>

              <div className="flex flex-col gap-1">
                <p className="text-label-caps uppercase font-display text-ink-600">
                  Region
                </p>
                <p className="text-body font-sans text-ink-900">
                  {partner.data.partner.Region}
                </p>
              </div>

              <div className="flex flex-col gap-1">
                <p className="text-label-caps uppercase font-display text-ink-600">
                  Address
                </p>
                <p className="text-body font-sans text-ink-900">
                  {partner.data.partner.Address}
                </p>
              </div>
            </div>
          </div>
        )}
      </div>
    </div>
  )
}
