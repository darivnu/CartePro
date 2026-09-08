import { apiFetch } from './client';
import type {
  PartnersResponse,
  PartnerDetailResponse,
  ValidateQrPayload,
  ValidateQrResponse,
  PartnerTransactionsResponse,
  PartnerDashboardResponse,
  PartnerRegistrationPayload,
} from '../types/partner';

export interface GetPartnersParams {
  search?: string;
  category?: string;
  page?: number;
  limit?: number;
}

export function getPartners(params: GetPartnersParams = {}) {
  const { search, category, page = 1, limit = 20 } = params;

  const query = new URLSearchParams();
  if (search) query.set('search', search);
  if (category) query.set('category', category);
  query.set('page', String(page));
  query.set('limit', String(limit));

  return apiFetch<PartnersResponse>(`/partners?${query.toString()}`);
}

export function getPartner(id: string) {
  return apiFetch<PartnerDetailResponse>(`/partners/${id}`);
}

export function collectPayment(payload: ValidateQrPayload) {
  return apiFetch<ValidateQrResponse>('/partners/me/validate', {
    method: 'POST',
    body: payload,
  });
}

export interface GetOwnTransactionsParams {
  page?: number;
  limit?: number;
  from?: string;
  to?: string;
}

export function getOwnTransactions(params: GetOwnTransactionsParams = {}) {
  const { page = 1, limit = 20, from, to } = params;

  const query = new URLSearchParams();
  query.set('page', String(page));
  query.set('limit', String(limit));
  if (from) query.set('from', from);
  if (to) query.set('to', to);

  return apiFetch<PartnerTransactionsResponse>(`/partners/me/transactions?${query.toString()}`);
}

export interface GetOwnDashboardParams {
  from?: string;
  to?: string;
}

export function getOwnDashboard(params: GetOwnDashboardParams = {}) {
  const { from, to } = params;

  const query = new URLSearchParams();
  if (from) query.set('from', from);
  if (to) query.set('to', to);

  return apiFetch<PartnerDashboardResponse>(`/partners/me/dashboard?${query.toString()}`);
}

export function registerPartner(payload: PartnerRegistrationPayload) {
  return apiFetch<void>('/partners/register', { method: 'POST', body: payload })
}
