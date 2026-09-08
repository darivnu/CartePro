import { useState } from 'react'
import { useParams, Link } from 'react-router-dom'
import { useAdminClientDetail, useCreateTopup, useCancelTransaction } from '../../admin/useAdminClients'
import { formatCents, parseEurosToCents } from '@/lib/money'
import { Input } from '@/components/ui/input'
import { Button } from '@/components/ui/button'
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogDescription,
  DialogFooter,
} from '@/components/ui/dialog'
import type { AdminTransaction } from '../../types/admin'
import { cn } from '@/lib/utils'

export function AdminClientDetail() {
  const { id } = useParams<{ id: string }>()
  const clientId = Number(id)

  const clientDetail = useAdminClientDetail(clientId)
  const createTopup = useCreateTopup()
  const [amount, setAmount] = useState('')
  const [comment, setComment] = useState('')

  const parsedAmount = parseEurosToCents(amount)

  function submitTopup() {
    if (parsedAmount === null) return
    createTopup.mutate(
      { client_id: clientId, amount: parsedAmount, comment },
      {
        onSuccess: () => {
          setAmount('')
          setComment('')
        },
      },
    )
  }

  const cancelTransaction = useCancelTransaction()
  const [cancellingTransaction, setCancellingTransaction] = useState<AdminTransaction | null>(null)
  const [cancelReason, setCancelReason] = useState('')

  function closeCancelDialog() {
    setCancellingTransaction(null)
    setCancelReason('')
  }

  function confirmCancel() {
    if (!cancellingTransaction) return
    cancelTransaction.mutate(
      { id: cancellingTransaction.ID, reason: cancelReason, clientId },
      { onSuccess: closeCancelDialog },
    )
  }

  const reversedTransactionIds = new Set(
    (clientDetail.data?.transactions ?? [])
      .map((transaction) => transaction.OriginalTransactionID)
      .filter((id): id is number => id !== null),
  )

  return (
    <div className="flex flex-col gap-4">
      <Link to="/admin/clients" className="text-label-caps uppercase font-display text-ink-600">
        Back
      </Link>

      {clientDetail.isLoading && <p className="text-body font-sans text-ink-600">Loading...</p>}
      {clientDetail.isError && (
        <p className="text-body font-sans text-danger-600">{clientDetail.error.message}</p>
      )}
      {clientDetail.isSuccess && (
        <>
          <div className="flex flex-col gap-1">
            <h2 className="text-h2 font-display text-brand-blue">{clientDetail.data.client.Name}</h2>
            <p className="text-body-strong font-sans text-ink-900">
              {formatCents(clientDetail.data.client.Balance)}
            </p>
          </div>

          <div className="flex flex-col gap-2">
            <h3 className="text-label-caps uppercase font-display text-ink-600">Top up</h3>
            <Input
              value={amount}
              onChange={(event) => setAmount(event.target.value)}
              placeholder="Amount (€)"
              inputMode="decimal"
            />
            <Input
              value={comment}
              onChange={(event) => setComment(event.target.value)}
              placeholder="Comment"
            />
            <Button
              disabled={parsedAmount === null || createTopup.isPending}
              onClick={submitTopup}
              className="rounded-panel bg-surface-300 text-button uppercase font-display text-graphite-900 shadow-key-secondary hover:bg-surface-300"
            >
              Credit balance
            </Button>
          </div>

          <h3 className="text-label-caps uppercase font-display text-ink-600">Transactions</h3>

          {clientDetail.data.transactions.length === 0 && (
            <p className="text-body font-sans text-ink-600">No transactions yet.</p>
          )}
          {clientDetail.data.transactions.length > 0 && (
            <div className="flex flex-col gap-[2px]">
              {clientDetail.data.transactions.map((transaction) => (
                <div
                  key={transaction.ID}
                  className="flex items-center justify-between rounded-row bg-surface-050 px-3 py-3 shadow-row-raised"
                >
                  <div className="flex flex-col gap-1">
                    <p className="text-body-strong font-sans text-ink-900 uppercase">
                      {transaction.Type}
                    </p>
                    <p className="text-caption font-sans text-ink-600">
                      {transaction.Comment ?? '—'}
                    </p>
                  </div>
                  <div className="flex flex-col items-end gap-1">
                    <p
                      className={cn(
                        'text-body-strong font-sans',
                        transaction.Type === 'debit' ? 'text-danger-600' : 'text-success-600',
                      )}
                    >
                      {transaction.Type === 'debit' ? '-' : '+'}
                      {formatCents(transaction.Amount)}
                    </p>
                    {transaction.OriginalTransactionID === null &&
                      transaction.Type === 'debit' &&
                      !reversedTransactionIds.has(transaction.ID) && (
                      <Button
                        size="sm"
                        onClick={() => setCancellingTransaction(transaction)}
                        className="rounded-panel bg-surface-300 text-button uppercase font-display text-danger-600 shadow-key-secondary hover:bg-surface-300"
                      >
                        Cancel
                      </Button>
                    )}
                  </div>
                </div>
              ))}
            </div>
          )}
        </>
      )}

      <Dialog
        open={cancellingTransaction !== null}
        onOpenChange={(open) => !open && closeCancelDialog()}
      >
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Cancel transaction #{cancellingTransaction?.ID}</DialogTitle>
            <DialogDescription>
              This creates a reversal transaction; the original record stays untouched.
            </DialogDescription>
          </DialogHeader>
          <Input
            value={cancelReason}
            onChange={(event) => setCancelReason(event.target.value)}
            placeholder="Reason for cancellation"
          />
          <DialogFooter>
            <Button
              disabled={cancelReason.trim() === '' || cancelTransaction.isPending}
              onClick={confirmCancel}
            >
              Cancel transaction
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  )
}
