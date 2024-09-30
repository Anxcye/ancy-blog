import type { ApiResponse } from '@/api/type'

// category list
export interface CategoryListData {
  id: number
  name: string
  parentId: number
  description: string
}

export interface CategoryListRes extends ApiResponse {
  data: CategoryListData[]
}

// category add
export interface CategoryAddParams {
  name: string

  parentId?: number
  description?: string
  status?: string
}

export interface CategoryAddRes extends ApiResponse {
  data: number
}
