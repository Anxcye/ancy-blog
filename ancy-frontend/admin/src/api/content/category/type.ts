import type { PageResponse } from '@/api/pageResult'
import type { ApiResponse } from '@/api/type'

// category list
export interface CategoryListData {
  id: number
  name: string
  parentId: number
  description: string
  createBy: number
  createTime: string
  updateBy: number
  updateTime: string
}

export interface CategoryListRes extends ApiResponse<CategoryListData[]> {}

// category add
export interface CategoryAddParams {
  id?: number
  name?: string
  parentId?: number
  description?: string
  status?: string
}

export interface CategoryAddRes extends ApiResponse<number> {}

// category page
export interface CategoryPageParams {
  name?: string
  status?: string
  pageNum: number
  pageSize: number
}

export interface CategoryPageRes extends PageResponse<CategoryListData> {}

export interface CategoryUpdateRes extends ApiResponse<boolean> {}

// export to xlsx
export interface ExportToXlsxRes {}

// category delete
export interface CategoryDeleteRes extends ApiResponse<boolean> {}
