import type { PageResponse } from '../pageResult'
import type { ApiResponse } from '../type'

export interface ReadListData {
  addFrom: 0 | 1
  author: string
  content: string
  createBy: number
  createTime: string
  id: number
  source: string
  updateBy: number | null
  updateTime: string | null
}

export interface ReadPageRes extends PageResponse<ReadListData> {}

export interface ReadPageParams {
  source?: string
  content?: string
  author?: string
  pageNum: number
  pageSize: number
}

export interface ReadAddParams {
  source?: string
  content?: string
  author?: string
  addFrom?: number
}

export interface ReadAddRes extends ApiResponse<number> {}

export interface ReadUpdateRes extends ApiResponse<boolean> {}

export interface ReadDeleteRes extends ApiResponse<boolean> {}
