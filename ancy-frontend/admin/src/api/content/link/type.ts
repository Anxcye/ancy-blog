import type { PageResponse } from '@/api/pageResult'
import type { ApiResponse } from '@/api/type'

export interface LinkAddParams {
  id?: number
  name?: string
  logo?: string
  description?: string
  address?: string
  status?: string
}

export interface LinkAddRes extends ApiResponse<number> {}

export interface LinkGetByIdData {
  id: number
  name: string
  logo: string
  description: string
  address: string
  status: string
  createBy: number
  createTime: string
}

export interface LinkGetByIdRes extends ApiResponse<LinkGetByIdData> {}

export interface LinkUpdateRes extends ApiResponse<boolean> {}

export interface LinkDeleteRes extends ApiResponse<boolean> {}

export interface LinkListData {
  id: number
  name: string
  logo: string
  description: string
  address: string
  status: string
  createBy: number
  createTime: Record<string, unknown>
  updateBy: number
  updateTime: Record<string, unknown>
  deleted: number
}

export interface LinkListRes extends ApiResponse<LinkListData[]> {}

export interface LinkPageRes extends PageResponse<LinkListData> {}

export interface LinkPageParams {
  name?: string
  status?: string
  pageNum: number
  pageSize: number
}

export interface LinkUpdateArticleRes extends ApiResponse<boolean> {}
