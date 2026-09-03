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
