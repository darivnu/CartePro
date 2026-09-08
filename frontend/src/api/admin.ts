import { apiFetch } from './client'
import type {
AdminPartnersResponse,
AdminPartnerResponse,
PartnerStatus,
AdminClientsResponse,
AdminClientDetailResponse,
CreateTopupRequest,
CreateTopupResponse,
CancelTransactionResponse,
AdminDashboardResponse,
}
from '../types/admin'

export interface GetAdminPartnersParams {
status?: PartnerStatus
page?: number
limit?: number
}

export function getAdminPartners(params: GetAdminPartnersParams = {}) {
const { status, page = 1, limit = 20 } = params

const query = new URLSearchParams()
if (status) query.set('status', status)
query.set('page', String(page))
query.set('limit', String(limit))

return apiFetch<AdminPartnersResponse>(`/admin/partners?${query.toString()}`)
}

export function approvePartner(id: number) {
return apiFetch<AdminPartnerResponse>(`/admin/partners/${id}/approve`, {
    method: 'POST',
})
}

export function rejectPartner(id: number, reason: string) {
return apiFetch<AdminPartnerResponse>(`/admin/partners/${id}/reject`, {
    method: 'POST',
    body: { reason },
})
}

export function updatePartnerStatus(id: number, status: PartnerStatus) {
return apiFetch<AdminPartnerResponse>(`/admin/partners/${id}/status`, {
    method: 'PATCH',
    body: { status },
})
}

export function setMinisterPick(id: number, ministerPick: boolean) {
return apiFetch<AdminPartnerResponse>(`/admin/partners/${id}/minister-pick`, {
    method: 'PATCH',
    body: { minister_pick: ministerPick },
})
}

export interface GetAdminClientsParams {
page?: number
limit?: number
}

export function getAdminClients(params: GetAdminClientsParams = {}) {
const { page = 1, limit = 20 } = params

const query = new URLSearchParams()
query.set('page', String(page))
query.set('limit', String(limit))

return apiFetch<AdminClientsResponse>(`/admin/clients?${query.toString()}`)
}

export function getAdminClientDetail(id: number) {
return apiFetch<AdminClientDetailResponse>(`/admin/clients/${id}`)
}

export function createTopup(payload: CreateTopupRequest) {
return apiFetch<CreateTopupResponse>('/admin/topups', {
    method: 'POST',
    body: payload,
})
}

export function cancelTransaction(id: number, reason: string) {
return apiFetch<CancelTransactionResponse>(`/admin/transactions/${id}/cancel`, {
    method: 'POST',
    body: { reason },
})
}

export interface GetAdminDashboardParams {
from?: string
to?: string
}

export function getAdminDashboard(params: GetAdminDashboardParams = {}) {
const { from, to } = params

const query = new URLSearchParams()
if (from) query.set('from', from)
if (to) query.set('to', to)

const queryString = query.toString()
return apiFetch<AdminDashboardResponse>(`/admin/dashboard${queryString ? `?${queryString}` : ''}`)
}
