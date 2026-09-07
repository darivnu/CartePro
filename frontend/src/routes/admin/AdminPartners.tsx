import { useState } from 'react'
import { useAdminPartners, useApprovePartner, useRejectPartner, useUpdatePartnerStatus, useSetMinisterPick } from '../../admin/useAdminPartners'
import { formatCents } from '@/lib/money'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogDescription,
  DialogFooter,
} from '@/components/ui/dialog'
import type { AdminPartner, PartnerStatus } from '../../types/admin'
import { cn } from '@/lib/utils'

const STATUS_FILTERS: { label: string; value: PartnerStatus }[] = [
{ label: 'Pending', value: 'pending' },
{ label: 'Approved', value: 'approved' },
{ label: 'Rejected', value: 'rejected' },
{ label: 'All', value: '' },
]

const STATUS_COLORS: Record<string, string> = {
pending: 'text-pending-600',
approved: 'text-success-600',
rejected: 'text-danger-600',
}

const STATUS_OPTIONS: PartnerStatus[] = ['pending', 'approved', 'rejected', 'suspended']

export function AdminPartners() {
const [status, setStatus] = useState<PartnerStatus>('pending')
const [page, setPage] = useState(1)

const partners = useAdminPartners({ status, page })
const approvePartner = useApprovePartner()
const rejectPartner = useRejectPartner()
const updatePartnerStatus = useUpdatePartnerStatus()
const setMinisterPick = useSetMinisterPick()
const [rejectingPartner, setRejectingPartner] = useState<AdminPartner | null>(null)
const [reason, setReason] = useState('')

function closeRejectDialog() {
    setRejectingPartner(null)
    setReason('')
}

function confirmReject() {
    if (!rejectingPartner) return
    rejectPartner.mutate(
        { id: rejectingPartner.ID, reason },
        { onSuccess: closeRejectDialog },
    )
}

function handleStatusChange(next: PartnerStatus) {
    setStatus(next)
    setPage(1)
}

return (
    <div className="flex flex-col gap-4">
    <h2 className="text-h2 font-display text-brand-blue">Partners</h2>

    <div className="flex gap-2">
        {STATUS_FILTERS.map((filter) => (
        <button
            key={filter.label}
            onClick={() => handleStatusChange(filter.value)}
            className={cn(
            'text-label-caps uppercase font-display text-ink-600',
            status === filter.value && 'text-brand-blue',
            )}
        >
            {filter.label}
        </button>
        ))}
    </div>

    {partners.isLoading && <p className="text-body font-sans text-ink-600">Loading...</p>}
    {partners.isError && (
        <p className="text-body font-sans text-danger-600">{partners.error.message}</p>
    )}
    {partners.isSuccess && partners.data.data.length === 0 && (
        <p className="text-body font-sans text-ink-600">No partners found.</p>
    )}
    {partners.isSuccess && partners.data.data.length > 0 && (
        <div className="flex flex-col gap-[2px]">
        {partners.data.data.map((partner) => (
            <div
            key={partner.ID}
            className="flex items-center justify-between rounded-row bg-surface-050 px-3 py-3 shadow-row-raised"
            >
            <div className="flex flex-col gap-1">
                <p className="text-body-strong font-sans text-ink-900">{partner.BusinessName}</p>
                <p className="text-caption font-sans text-ink-600">
                {partner.Category} · {partner.Region}
                </p>
            </div>
            <div className="flex flex-col items-end gap-1">
                <p
                className={cn(
                    'text-caption uppercase font-display',
                    STATUS_COLORS[partner.Status] ?? 'text-ink-600',
                )}
                >
                {partner.Status}
                </p>
                <p className="text-caption font-sans text-ink-600">{formatCents(partner.Balance)}</p>
                <select
                    value={partner.Status}
                    disabled={updatePartnerStatus.isPending}
                    onChange={(event) =>
                    updatePartnerStatus.mutate({ id: partner.ID, status: event.target.value })
                    }
                    className="rounded-panel border border-ink-600 bg-surface-050 px-2 py-1 text-caption font-sans text-ink-900"
                >
                    {STATUS_OPTIONS.map((option) => (
                    <option key={option} value={option}>
                        {option}
                    </option>
                    ))}
                </select>
                {partner.Status === 'approved' && (
                <label className="flex items-center gap-1 text-caption font-sans text-ink-600">
                    <input
                    type="checkbox"
                    checked={partner.MinisterPick}
                    disabled={setMinisterPick.isPending}
                    onChange={(event) =>
                        setMinisterPick.mutate({ id: partner.ID, ministerPick: event.target.checked })
                    }
                    />
                    Minister's Pick
                </label>
                )}
                {partner.Status === 'pending' && (
                <div className="flex gap-2">
                    <Button
                        size="sm"
                        disabled={approvePartner.isPending}
                        onClick={() => approvePartner.mutate(partner.ID)}
                        className="rounded-panel bg-surface-300 text-button uppercase font-display text-graphite-900 shadow-key-secondary hover:bg-surface-300"
                    >
                        Approve
                    </Button>
                    <Button
                        size="sm"
                        onClick={() => setRejectingPartner(partner)}
                        className="rounded-panel bg-surface-300 text-button uppercase font-display text-danger-600 shadow-key-secondary hover:bg-surface-300"
                    >
                        Reject
                    </Button>
                </div>
                )}
            </div>
            </div>
        ))}
        </div>
    )}

    {partners.isSuccess && (
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
            disabled={partners.data.data.length < partners.data.meta.limit}
            onClick={() => setPage((current) => current + 1)}
            className="rounded-panel bg-surface-300 text-button uppercase font-display text-graphite-900 shadow-key-secondary hover:bg-surface-300"
        >
            Next
        </Button>
        </div>
    )}

    <Dialog open={rejectingPartner !== null} onOpenChange={(open) => !open && closeRejectDialog()}>
        <DialogContent>
        <DialogHeader>
            <DialogTitle>Reject {rejectingPartner?.BusinessName}</DialogTitle>
            <DialogDescription>This reason will be visible to the partner.</DialogDescription>
        </DialogHeader>
        <Input
            value={reason}
            onChange={(event) => setReason(event.target.value)}
            placeholder="Reason for rejection"
        />
        <DialogFooter>
            <Button disabled={reason.trim() === '' || rejectPartner.isPending} onClick={confirmReject}>
            Reject
            </Button>
        </DialogFooter>
        </DialogContent>
    </Dialog>
    </div>
)
}