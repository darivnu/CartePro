import { useCallback, useState } from 'react'
import { Link } from 'react-router-dom'
import { useAuth, useLogout } from '../../auth/useAuth'
import { Wordmark } from '@/components/Wordmark'
import { SimulationNotice } from '@/components/SimulationNotice'
import { QrScanner } from '@/components/QrScanner'
import { useCollectPayment, useOwnTransactions, useOwnPartner } from '../../partners/usePartnersData'
import type { PartnerStatus } from '../../types/partner'
import { formatCents, parseEurosToCents } from '../../lib/money'
import { ApiError } from '../../api/client'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { cn } from '@/lib/utils'
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'

const TRANSACTION_TYPE_LABELS: Record<string, string> = {
  debit: 'Payment',
  topup: 'Top-up',
  reversal: 'Reversal',
}

const STATUS_COLORS: Record<PartnerStatus, string> = {
  pending: 'text-pending-600',
  approved: 'text-success-600',
  rejected: 'text-danger-600',
}

function toRfc3339Start(date: string) {
  return date ? `${date}T00:00:00Z` : undefined
}

function toRfc3339End(date: string) {
  return date ? `${date}T23:59:59Z` : undefined
}


function getCollectErrorMessage(error: Error | null): string {
  if (error instanceof ApiError) {
    switch (error.status) {
      case 400:
        return 'This QR code is invalid or expired. Ask the client to generate a new one.'
      case 409:
        return 'This QR code has already been used or expired. Ask the client to generate a new one.'
      case 402:
        return "This client doesn't have enough balance for this amount."
      case 403:
        return "Your partner account isn't approved yet."
      case 401:
        return 'Your session has expired. Please log in again.'
      case 404:
        return 'Client not found. Please try scanning again.'
    }
  }
  return 'Something went wrong. Please try again.'
}

export function PartnerHome() {
  const logout = useLogout()
  const collectPayment = useCollectPayment()

  const [dialogOpen, setDialogOpen] = useState(false)
  const [step, setStep] = useState<'amount' | 'scan'>('amount')
  const [amountInput, setAmountInput] = useState('')
  const [amountError, setAmountError] = useState<string | null>(null)
  const [cents, setCents] = useState<number | null>(null)
  const [idempotencyKey, setIdempotencyKey] = useState<string | null>(null)
  const [cameraError, setCameraError] = useState<string | null>(null)
  const { user } = useAuth()
  const [fromDate, setFromDate] = useState('')
  const [toDate, setToDate] = useState('')
  const [page, setPage] = useState(1)

  const ownPartner = useOwnPartner()
  const status = ownPartner.data?.partner.Status
  const isApproved = status === 'approved'

  const transactions = useOwnTransactions({
    page,
    from: toRfc3339Start(fromDate),
    to: toRfc3339End(toDate),
  })

  function handleDialogOpenChange(open: boolean) {
    setDialogOpen(open)
    if (open) {
      setStep('amount')
      setAmountInput('')
      setAmountError(null)
      setCents(null)
      setIdempotencyKey(null)
      setCameraError(null)
      collectPayment.reset()
    }
  }

  function handleConfirmAmount() {
    const parsed = parseEurosToCents(amountInput)
    if (parsed === null) {
      setAmountError('Enter a valid amount above 0.')
      return
    }
    setAmountError(null)
    setCents(parsed)
    setIdempotencyKey(crypto.randomUUID())
    setStep('scan')
  }

  const handleScan = useCallback(
    (qrPayload: string) => {
      if (cents === null || idempotencyKey === null) {
        return
      }
      collectPayment.mutate({
        qr_payload: qrPayload,
        amount: cents,
        idempotency_key: idempotencyKey,
      })
    },
    [cents, idempotencyKey, collectPayment.mutate],
  )

  const handleCameraError = useCallback((message: string) => {
    setCameraError(message)
  }, [])

  function handleScanAgain() {
    setCameraError(null)
    collectPayment.reset()
  }

  function handleCollectAnother() {
    setStep('amount')
    setAmountInput('')
    setAmountError(null)
    setCents(null)
    setIdempotencyKey(null)
    setCameraError(null)
    collectPayment.reset()
  }

  if (!user) {
    return null
  }


  return (
    <div className="min-h-screen bg-surface-200 p-4">
      <div className="mx-auto flex max-w-sm flex-col gap-6">
        <h1 className="sr-only">CartePro — partner dashboard</h1>
        <div className="flex items-center justify-between">
          <div className="flex items-center gap-2">
            <Wordmark />
            {status && (
              <p
                className={cn(
                  'text-status-caps uppercase font-display',
                  STATUS_COLORS[status],
                )}
              >
                {status}
              </p>
            )}
          </div>
          <button
            onClick={() => logout.mutate()}
            className="text-label-caps uppercase font-display text-ink-600"
          >
            Log out
          </button>
        </div>

        <div className="flex flex-col gap-2 rounded-card bg-card-gradient p-4 shadow-card-highlight">
          <p className="text-micro-caps uppercase font-sans text-graphite-400">
            Collect a payment
          </p>
          <p className="text-body font-sans text-graphite-100">
            Scan a client's QR code to charge their card for a purchase.
          </p>
          <SimulationNotice tone="dark" />
        </div>

        <Button
          disabled={!isApproved}
          onClick={() => handleDialogOpenChange(true)}
          className="h-11 w-full rounded-panel text-button uppercase font-display text-primary-foreground"
        >
          Scan client QR
        </Button>

        <div className="flex items-center justify-between">
          <p className="text-h2 font-display text-graphite-900">Transaction history</p>
          <Link
            to="/partner/dashboard"
            className="text-label-caps uppercase font-display text-brand-blue"
          >
            Dashboard
          </Link>
        </div>

        <div className="flex gap-2">
          <Input
            type="date"
            value={fromDate}
            onChange={(e) => {
              setFromDate(e.target.value)
              setPage(1)
            }}
            aria-label="From date"
          />
          <Input
            type="date"
            value={toDate}
            onChange={(e) => {
              setToDate(e.target.value)
              setPage(1)
            }}
            aria-label="To date"
          />
        </div>

        {transactions.isLoading && (
          <p className="text-body font-sans text-ink-600">Loading…</p>
        )}

        {transactions.isError && (
          <p className="text-body font-sans text-danger-600">{transactions.error.message}</p>
        )}

        {transactions.isSuccess && transactions.data.data.length === 0 && (
          <p className="text-body font-sans text-ink-600">No transactions found.</p>
        )}

        {transactions.isSuccess && transactions.data.data.length > 0 && (
          <div className="flex flex-col gap-[2px]">
            {transactions.data.data.map((tx) => {
              const isReceived = tx.ReceiverUserID === Number(user.id)
              const clientName = isReceived ? tx.SenderName : tx.ReceiverName
              return (
                <div
                  key={tx.ID}
                  className="flex items-center justify-between rounded-row bg-surface-050 px-3 py-3 shadow-row-raised"
                >
                  <div className="flex flex-col gap-1">
                    <p className="text-body-strong font-sans text-ink-900">
                      {clientName}
                    </p>
                    <p className="text-caption font-sans text-ink-600">
                      {TRANSACTION_TYPE_LABELS[tx.Type] ?? tx.Type} · {new Date(tx.CreatedAt).toLocaleDateString('fr-FR')}
                    </p>
                  </div>
                  <div className="flex flex-col items-end gap-1">
                    <p
                      className={cn(
                        'text-body-strong font-sans',
                        isReceived ? 'text-success-600' : 'text-danger-600',
                      )}
                    >
                      {isReceived ? '+' : '-'}
                      {formatCents(tx.Amount)}
                    </p>
                    {tx.Cancelled && (
                      <p className="text-caption uppercase font-display text-danger-600">
                        Reversed
                      </p>
                    )}
                  </div>
                </div>
              )
            })}
          </div>
        )}

        {transactions.isSuccess && (
          <div className="flex items-center justify-between">
            <Button
              size="sm"
              disabled={page <= 1}
              onClick={() => setPage((c) => c - 1)}
              className="rounded-panel bg-surface-300 text-button uppercase font-display text-graphite-900 shadow-key-secondary hover:bg-surface-300"
            >
              Prev
            </Button>
            <p className="text-caption font-sans text-ink-600">
              Page {page} of {transactions.data.meta.total_pages}
            </p>
            <Button
              size="sm"
              disabled={page >= transactions.data.meta.total_pages}
              onClick={() => setPage((c) => c + 1)}
              className="rounded-panel bg-surface-300 text-button uppercase font-display text-graphite-900 shadow-key-secondary hover:bg-surface-300"
            >
              Next
            </Button>
          </div>
        )}
      </div>

      <Dialog open={dialogOpen} onOpenChange={handleDialogOpenChange}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle className="text-h2 font-display text-brand-blue">
              Collect payment
            </DialogTitle>
            <DialogDescription className="text-body font-sans text-ink-600">
              {step === 'amount'
                ? "Enter the amount to charge, then scan the client's QR code."
                : "Point the camera at the client's QR code."}
            </DialogDescription>
          </DialogHeader>

          <SimulationNotice />

          {step === 'amount' && (
            <div className="flex flex-col gap-3">
              <Input
                type="text"
                inputMode="decimal"
                placeholder="0,00"
                value={amountInput}
                onChange={(e) => setAmountInput(e.target.value)}
                aria-label="Amount in euros"
              />
              {amountError && (
                <p className="text-body font-sans text-danger-600">{amountError}</p>
              )}
              <Button
                onClick={handleConfirmAmount}
                className="h-11 w-full rounded-panel text-button uppercase font-display text-primary-foreground"
              >
                Continue to scan
              </Button>
            </div>
          )}

          {step === 'scan' && (
            <div className="flex flex-col items-center gap-3">
              {collectPayment.isIdle && !cameraError && (
                <>
                  <QrScanner onScan={handleScan} onError={handleCameraError} />
                  <p className="text-caption font-sans text-ink-600">
                    Charging {cents !== null ? formatCents(cents) : ''}
                  </p>
                </>
              )}

              {collectPayment.isIdle && cameraError && (
                <>
                  <p className="text-body font-sans text-danger-600">{cameraError}</p>
                  <Button
                    onClick={() => setCameraError(null)}
                    className="h-11 rounded-panel text-button uppercase font-display text-primary-foreground"
                  >
                    Try again
                  </Button>
                </>
              )}

              {collectPayment.isPending && (
                <p className="text-body font-sans text-ink-600">Validating…</p>
              )}

              {collectPayment.isSuccess && (
                <>
                  <p className="text-status-caps uppercase font-display text-success-600">
                    Payment collected
                  </p>
                  <p className="text-balance-xl font-display text-graphite-100">
                    {formatCents(collectPayment.data.transaction.Amount)}
                  </p>
                  <div className="flex w-full flex-col gap-2">
                    <Button
                      onClick={handleCollectAnother}
                      className="h-11 w-full rounded-panel text-button uppercase font-display text-primary-foreground"
                    >
                      Collect another payment
                    </Button>
                    <Button
                      onClick={() => handleDialogOpenChange(false)}
                      className="h-11 w-full rounded-panel bg-surface-300 text-button uppercase font-display text-graphite-900 shadow-key-secondary hover:bg-surface-300"
                    >
                      Done
                    </Button>
                  </div>
                </>
              )}

              {collectPayment.isError && (
                <>
                  <p className="text-body font-sans text-danger-600">
                    {getCollectErrorMessage(collectPayment.error)}
                  </p>
                  <div className="flex w-full flex-col gap-2">
                    <Button
                      onClick={handleScanAgain}
                      className="h-11 w-full rounded-panel text-button uppercase font-display text-primary-foreground"
                    >
                      Scan again
                    </Button>
                    <Button
                      onClick={() => handleDialogOpenChange(false)}
                      className="h-11 w-full rounded-panel bg-surface-300 text-button uppercase font-display text-graphite-900 shadow-key-secondary hover:bg-surface-300"
                    >
                      Cancel
                    </Button>
                  </div>
                </>
              )}
            </div>
          )}
        </DialogContent>
      </Dialog>
    </div>
  )
}
