import type { ApiResponse } from '@/api/type'

// tag list
export interface TagListData {
  id: number
  name: string
  remark: string
  createBy: number
  createTime: string
  updateBy: number
  updateTime: string
  deleted: number
}

export interface TagListRes extends ApiResponse {
  data: TagListData[]
}
// tag add
export interface TagAddParams {
  name: string
  remark?: string
}

export interface TagAddRes extends ApiResponse {
  data: number
}
