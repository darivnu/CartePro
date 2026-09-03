import { useCallback, useState } from 'react'
import { useLogout } from '../../auth/useAuth'
import { Wordmark } from '@/components/Wordmark'
import { SimulationNotice } from '@/components/SimulationNotice'
import { QrScanner } from '@/components/QrScanner'
import { useCollectPayment } from '../../partners/usePartnersData'
import { formatCents, parseEurosToCents } from '../../lib/money'
import { ApiError } from '../../api/client'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'

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

  return (
    <div className="min-h-screen bg-surface-200 p-4">
      <div className="mx-auto flex max-w-sm flex-col gap-6">
        <h1 className="sr-only">Ticket Tout — partner dashboard</h1>
        <div className="flex items-center justify-between">
          <Wordmark />
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
          onClick={() => handleDialogOpenChange(true)}
          className="h-11 w-full rounded-panel text-button uppercase font-display text-primary-foreground"
        >
          Scan client QR
        </Button>
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
