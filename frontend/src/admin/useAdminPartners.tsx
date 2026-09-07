import { useMutation, useQuery } from '@tanstack/react-query'
import {
getAdminPartners,
approvePartner,
rejectPartner,
updatePartnerStatus,
setMinisterPick,
} from '../api/admin'
import type { GetAdminPartnersParams } from '../api/admin'
import type { PartnerStatus } from '../types/admin'
import { queryClient } from '../api/queryClient'

export function useAdminPartners(params: GetAdminPartnersParams = {}) {
const { status, page = 1, limit = 20 } = params

return useQuery({
    queryKey: ['admin-partners', status, page, limit],
    queryFn: () => getAdminPartners({ status, page, limit }),
})
}

function invalidateAdminPartners() {
queryClient.invalidateQueries({ queryKey: ['admin-partners'] })
}

export function useApprovePartner() {
return useMutation({
    mutationFn: approvePartner,
    onSuccess: invalidateAdminPartners,
})
}

export function useRejectPartner() {
return useMutation({
    mutationFn: ({ id, reason }: { id: number; reason: string }) => rejectPartner(id, reason),
    onSuccess: invalidateAdminPartners,
})
}

export function useUpdatePartnerStatus() {
return useMutation({
    mutationFn: ({ id, status }: { id: number; status: PartnerStatus }) =>
    updatePartnerStatus(id, status),
    onSuccess: invalidateAdminPartners,
})
}

export function useSetMinisterPick() {
return useMutation({
    mutationFn: ({ id, ministerPick }: { id: number; ministerPick: boolean }) =>
    setMinisterPick(id, ministerPick),
    onSuccess: invalidateAdminPartners,
})
}