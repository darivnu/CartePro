import { useEffect, useState } from 'react'
import { Link } from 'react-router-dom'
import { QRCodeSVG } from 'qrcode.react'
import { useLogout } from '../../auth/useAuth'
import { Wordmark } from '@/components/Wordmark'
import { SimulationNotice } from '@/components/SimulationNotice'
import { useBalance, useTransactions, useGenerateQrCode } from '../../client/useClientData'
import { formatCents } from '../../lib/money'
import { Button } from '@/components/ui/button'
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'

const dateFormatter = new Intl.DateTimeFormat('fr-FR', {
  day: 'numeric',
  month: 'short',
  hour: '2-digit',
  minute: '2-digit',
})

function useCountdown(expiresAt: string | undefined) {
  const [secondsLeft, setSecondsLeft] = useState(0)

  useEffect(() => {
    if (!expiresAt) {
      return
    }

    const expiresAtMs = new Date(expiresAt).getTime()

    function updateSecondsLeft() {
      const diffMs = expiresAtMs - Date.now()
      setSecondsLeft(Math.max(0, Math.floor(diffMs / 1000)))
    }

    updateSecondsLeft()
    const interval = setInterval(updateSecondsLeft, 1000)

    return () => clearInterval(interval)
  }, [expiresAt])

  return secondsLeft
}


function formatCountdown(totalSeconds: number) {
  const minutes = Math.floor(totalSeconds / 60)
  const seconds = totalSeconds % 60
  return `${minutes}:${seconds.toString().padStart(2, '0')}`
}

export function ClientHome() {
  const logout = useLogout()
  const transactions = useTransactions()
  const qrCode = useGenerateQrCode()
  const [qrOpen, setQrOpen] = useState(false)
  const [balanceAtOpen, setBalanceAtOpen] = useState<number | null>(null)

  const secondsLeft = useCountdown(qrCode.data?.expires_at)
  const isExpired = qrCode.isSuccess && secondsLeft <= 0
  const balance = useBalance({
    refetchInterval: qrOpen && qrCode.isSuccess && !isExpired ? 2000 : 10000,
  })

  useEffect(() => {
    if (!qrOpen || balanceAtOpen === null || !balance.data) {
      return
    }
    if (balance.data.balance !== balanceAtOpen) {
      setQrOpen(false)
      transactions.refetch()
    }
  }, [balance.data, qrOpen, balanceAtOpen, transactions])

  function handleGenerate() {
    setBalanceAtOpen(balance.data?.balance ?? null)
    setQrOpen(true)
    qrCode.mutate()
  }

  return (
    <div className="min-h-screen bg-surface-200 p-4">
      <div className="mx-auto flex max-w-sm flex-col gap-6">
        <h1 className="sr-only">CartePro — client dashboard</h1>
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
          <SimulationNotice tone="dark" />
        </div>

        <Button
          onClick={handleGenerate}
          className="h-11 w-full rounded-panel text-button uppercase font-display text-primary-foreground"
        >
          Generate QR code
        </Button>

        <Button
          asChild
          className="h-11 w-full rounded-panel bg-surface-300 text-button uppercase font-display text-graphite-900 shadow-key-secondary hover:bg-surface-300"
        >
          <Link to="/client/partners">Browse partners</Link>
        </Button>

        <div className="flex flex-col gap-3">
          <h2 className="text-h2 font-display text-brand-blue">Transactions</h2>
          <SimulationNotice />

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
                      {tx.partner ? tx.partner.business_name : tx.type === 'topup' ? 'Top-up' : 'Purchase'}
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

      <Dialog open={qrOpen} onOpenChange={setQrOpen}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle className="text-h2 font-display text-brand-blue">
              Payment QR code
            </DialogTitle>
            <DialogDescription className="text-body font-sans text-ink-600">
              Show this to the partner to pay. It expires in 5 minutes and can
              only be used once.
            </DialogDescription>
          </DialogHeader>

          <SimulationNotice />

          {qrCode.isPending && (
            <p className="text-body font-sans text-ink-600">Generating...</p>
          )}
          {qrCode.isError && (
            <p className="text-body font-sans text-danger-600">
              {qrCode.error.message}
            </p>
          )}
          {qrCode.isSuccess && !isExpired && (
            <div className="flex flex-col items-center gap-3">
              <div role="img" aria-label="Payment QR code">
                <QRCodeSVG value={qrCode.data.qr_payload} size={200} />
              </div>
              <p className="text-status-caps uppercase font-display text-pending-600">
                Expires in {formatCountdown(secondsLeft)}
              </p>
            </div>
          )}
          {qrCode.isSuccess && isExpired && (
            <div className="flex flex-col items-center gap-3">
              <p className="text-status-caps uppercase font-display text-danger-600">
                Expired
              </p>
              <Button
                onClick={handleGenerate}
                className="h-11 rounded-panel text-button uppercase font-display text-primary-foreground"
              >
                Generate new code
              </Button>
            </div>
          )}
        </DialogContent>
      </Dialog>
    </div>
  )
}