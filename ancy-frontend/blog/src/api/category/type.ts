import type { ApiResponse } from '../ApiResponse'

export interface CategoryListData {
  id: number
  name: string
  parentId: number
  description: string
}

export interface CategoryListRes extends ApiResponse<CategoryListData[]> {}
