import { apiFetch } from './client';
import type {Balance, Transaction, TransactionsResponse, QrCode} from '../types/client';

export function getBalance() {
  return apiFetch<Balance>('/clients/me/balance');
}

interface RawTransaction {
  ID: number
  Type: 'debit' | 'topup'
  Amount: number
  CreatedAt: string
  Partner: { ID: number; BusinessName: string } | null
}

interface RawTransactionsResponse {
  data: RawTransaction[]
  meta: TransactionsResponse['meta']
}

function normalizeTransaction(raw: RawTransaction): Transaction {
  const createdAt = raw.CreatedAt && !Number.isNaN(new Date(raw.CreatedAt).getTime())
    ? raw.CreatedAt
    : new Date(0).toISOString()

  return {
    id: raw.ID,
    type: raw.Type,
    amount: raw.Amount,
    created_at: createdAt,
    partner: raw.Partner ? { id: raw.Partner.ID, business_name: raw.Partner.BusinessName } : null,
  }
}

export async function getTransactions(page: number = 1, limit: number = 20): Promise<TransactionsResponse> {
  const raw = await apiFetch<RawTransactionsResponse>(`/clients/me/transactions?page=${page}&limit=${limit}`);
  return {
    data: raw.data.map(normalizeTransaction),
    meta: raw.meta,
  };
}

export function generateQrCode() {
    return apiFetch<QrCode>('/clients/me/qrcode', {
        method: 'POST',
    });
}
