import { apiFetch } from './client';
import type {Balance, TransactionsResponse, QrCode} from '../types/client';

export function getBalance() {
  return apiFetch<Balance>('/clients/me/balance');
}

export function getTransactions(page: number = 1, limit: number = 20) {
  return apiFetch<TransactionsResponse>(`/clients/me/transactions?page=${page}&limit=${limit}`);
}

export function generateQrCode() {
    return apiFetch<QrCode>('/clients/me/qrcode', {
        method: 'POST',
    });
}
