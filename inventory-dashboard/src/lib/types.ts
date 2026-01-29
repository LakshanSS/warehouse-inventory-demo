export interface Product {
  id: number
  sku: string
  name: string
  quantity: number
  category?: string
  location?: string
  updated_at: string
}

export interface CreateProductRequest {
  sku: string
  name: string
  quantity?: number
  category?: string
  location?: string
}

export interface UpdateQuantityRequest {
  quantity: number
}
