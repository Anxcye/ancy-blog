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
export interface TagAddData {
  id?: number
  name?: string
  remark?: string
}

export interface TagAddRes extends ApiResponse {
  data: number
}

// tag data
export interface TagData {
  id: number
  name: string
  remark: string
  createBy: number
  createTime: string
  updateBy: number
  updateTime: string
}
// tag update
export interface TagUpdateRes extends ApiResponse {
  data: boolean
}

// tag delete
export interface TagDeleteRes extends ApiResponse {
  data: boolean
}

// tag page
export interface TagPageRes extends ApiResponse {
  data: {
    total: number
    rows: TagListData[]
  }
}

export interface TagPageParams {
  pageNum: number
  pageSize: number
  name?: string
  remark?: string
}
