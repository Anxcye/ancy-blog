import type { ApiResponse } from '../ApiResponse'

export interface LinkListData {
  id: number
  name: string
  logo: string
  description: string
  address: string
  status: string
  createBy: number
  createTime: string
  updateBy: number
  updateTime: string
}

export interface LinkListRes extends ApiResponse<LinkListData[]> {}

export interface LinkAddParams {
  name: string
  logo: string
  description: string
  address: string
}

export interface LinkAddRes extends ApiResponse<number> {}
