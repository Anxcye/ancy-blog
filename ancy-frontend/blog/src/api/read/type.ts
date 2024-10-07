import type { PageResponse } from '../ApiResponse'

export interface ReadPageParam {
  source?: string
  content?: string
  author?: string
  pageNum: number
  pageSize: number
}

export interface ReadData {
  addFrom: number
  author: string
  content: string
  createBy: number
  createTime: string
  id: number
  source: string
  updateBy: number
  updateTime: string
}

export interface ReadPageRes extends PageResponse<ReadData> {}
