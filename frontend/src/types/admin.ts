export type PartnerStatus = 'pending' | 'approved' | 'rejected' | string

export interface AdminPartner {
ID: number
BusinessName: string
Siret: number
Category: string
Address: string
Region: string
Balance: number
Status: PartnerStatus
RejectReason: string | null
MinisterPick: boolean
CreatedAt: string
}

export interface AdminPartnersMeta {
page: number
limit: number
total: number
}

export interface AdminPartnersResponse {
data: AdminPartner[]
meta: AdminPartnersMeta
}

export interface AdminPartnerResponse {
partner: AdminPartner
}

export interface AdminClient {
ID: number
UserID: number
Name: string
Balance: number
CreatedAt: string
}

export interface AdminClientsMeta {
page: number
limit: number
total: number
}

export interface AdminClientsResponse {
data: AdminClient[]
meta: AdminClientsMeta
}

export type TransactionType = 'debit' | 'topup' | 'reversal' | string

export interface AdminTransaction {
ID: number
SenderUserID: number
ReceiverUserID: number
Amount: number
Comment: string | null
Type: TransactionType
CreatedAt: string
OriginalTransactionID: number | null
}

export interface AdminClientDetailResponse {
client: AdminClient
transactions: AdminTransaction[]
}

export interface CreateTopupRequest {
client_id: number
amount: number
comment: string
}

export interface CreateTopupResponse {
transaction: AdminTransaction
}

export interface CancelTransactionResponse {
transaction: AdminTransaction
}