export interface Balance {
    balance: number
    currency: string
    updated_at: string
}

export type TransactionType = 'debit' | 'topup'

export interface TransactionPartner {
    business_name: string
}

export interface Transaction {
    id: number
    type: TransactionType
    amount: number
    created_at: string
    partner: TransactionPartner | null
}

export interface TransactionsMeta {
    page: number
    limit: number
    total: number
    total_pages: number
}

export interface TransactionsResponse {
    data: Transaction[]
    meta: TransactionsMeta
}

export interface QrCode {
    token: string
    qr_payload: string
    expires_at: string
}
