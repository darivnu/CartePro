import { apiFetch } from './client';
import type { PartnersResponse, PartnerDetailResponse, ValidateQrPayload, ValidateQrResponse } from '../types/partner';

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
