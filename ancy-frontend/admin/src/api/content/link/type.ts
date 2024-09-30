import type { ApiResponse } from '@/api/type'

export interface LinkAddParams {
  id?: number
  name?: string
  logo?: string
  description?: string
  address?: string
  status?: string
}

export interface LinkAddRes extends ApiResponse {
  data: number
}

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

export interface LinkGetByIdRes extends ApiResponse {
  data: LinkGetByIdData
}

export interface LinkUpdateRes extends ApiResponse {
  data: boolean
}

export interface LinkDeleteRes extends ApiResponse {
  data: boolean
}

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

export interface LinkListRes extends ApiResponse {
  data: LinkListData[]
}

export interface LinkPageRes extends ApiResponse {
  data: {
    total: number
    rows: LinkListData[]
  }
}

export interface LinkPageParams {
  name?: string
  status?: string
  pageNum: number
  pageSize: number
}
