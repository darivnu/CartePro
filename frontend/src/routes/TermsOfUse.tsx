import { Link } from 'react-router-dom'
import { Wordmark } from '@/components/Wordmark'

export function TermsOfUse() {
  return (
    <div className="min-h-screen bg-surface-200 p-4">
      <div className="mx-auto flex max-w-sm flex-col gap-6 pb-12">
        <div className="flex items-center justify-between">
          <Wordmark />
          <Link
            to="/login"
            className="text-label-caps uppercase font-display text-ink-600"
          >
            Back
          </Link>
        </div>

        <h1 className="text-h1 font-display text-brand-blue">Terms of Use</h1>

        <p className="text-body font-sans text-ink-600">
          Ticket Tout is a simulation. No real monetary value is issued,
          transferred, or held at any point: balances, top-ups and payments
          shown in the application are for demonstration purposes only and do
          not constitute a payment service.
        </p>

        <div className="flex flex-col gap-2">
          <h2 className="text-h2 font-display text-brand-blue">Insufficient balance</h2>
          <p className="text-body font-sans text-ink-600">
            A payment that would exceed your available balance is declined in
            full. No partial payment is taken, your balance is not changed,
            and no transaction is recorded. You may generate a new QR code to
            try again.
          </p>
        </div>

        <div className="flex flex-col gap-2">
          <h2 className="text-h2 font-display text-brand-blue">Unspent balance</h2>
          <p className="text-body font-sans text-ink-600">
            An unspent balance does not expire and is not forfeited at the end
            of any period. It remains available until spent.
          </p>
        </div>

        <div className="flex flex-col gap-2">
          <h2 className="text-h2 font-display text-brand-blue">Transaction finality</h2>
          <p className="text-body font-sans text-ink-600">
            Once validated, a transaction is final. Ticket Tout does not
            currently offer a way to cancel, reverse, or refund a validated
            transaction.
          </p>
        </div>
      </div>
    </div>
  )
}
