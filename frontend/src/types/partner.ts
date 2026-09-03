export interface Partner {
  id: number
  business_name: string
  category: string
  region: string
  address: string
}

export interface PartnersMeta {
  page: number
  limit: number
  total: number
  total_pages: number
}

export interface PartnersResponse {
  data: Partner[]
  meta: PartnersMeta
}

export interface PartnerDetailResponse {
  partner: Partner
}
