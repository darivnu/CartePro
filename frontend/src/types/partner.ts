export interface Partner {
  ID: number
  BusinessName: string
  Category: string
  Region: string
  Address: string
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

