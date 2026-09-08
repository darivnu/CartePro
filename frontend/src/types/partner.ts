export interface Partner {
  ID: number
  BusinessName: string
  Category: string
  Region: string
  Address: string
  MinisterPick: boolean
}

export interface PartnersMeta {
  page: number
  limit: number
  total: number
}

export interface PartnersResponse {
  data: Partner[]
  meta: PartnersMeta
}

export interface PartnerDetailResponse {
  partner: Partner
}

export interface CollectedTransaction {
  ID: number
  ClientID: number
  PartnerID: number
  Amount: number
  Type: string
  CreatedAt: string
  IdempotencyKey: string
}

export interface ValidateQrPayload {
  qr_payload: string
  amount: number
  idempotency_key: string
}

export interface ValidateQrResponse {
  transaction: CollectedTransaction
}

export type TransactionType = 'debit' | 'topup' | 'reversal' | string

export interface PartnerTransaction {
  ID: number
  SenderUserID: number
  ReceiverUserID: number
  SenderName: string
  ReceiverName: string
  Amount: number
  CreatedAt: string
  Comment: string | null
  Type: TransactionType
  Cancelled: boolean
  OriginalTransactionID: number | null
}

export interface PartnerTransactionsMeta {
  page: number
  limit: number
  total: number
  total_pages: number
  from: string
  to: string
}

export interface PartnerTransactionsResponse {
  data: PartnerTransaction[]
  meta: PartnerTransactionsMeta
}

export interface PartnerDashboardDayBucket {
  date: string
  total_received: number
  transaction_count: number
}

export interface PartnerDashboardResponse {
  total_received: number
  transaction_count: number
  by_day: PartnerDashboardDayBucket[]
}

export interface PartnerRegistrationPayload {
  business_name: string
  siret: number
  category: string
  address: string
  region: string
  contact_email: string
  password: string
}
