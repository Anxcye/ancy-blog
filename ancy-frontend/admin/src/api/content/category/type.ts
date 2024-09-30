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
  id?: number
  name?: string
  parentId?: number
  description?: string
  status?: string
}

export interface CategoryAddRes extends ApiResponse {
  data: number
}

// category page
export interface CategoryPageParams {
  name?: string
  status?: string
  pageNum: number
  pageSize: number
}

export interface CategoryPageRes extends ApiResponse {
  data: {
    total: number
    rows: CategoryListData[]
  }
}

export interface CategoryUpdateRes extends ApiResponse {
  data: boolean
}

// export to xlsx
export interface ExportToXlsxRes {}

// category delete
export interface CategoryDeleteRes extends ApiResponse {
  data: boolean
}
