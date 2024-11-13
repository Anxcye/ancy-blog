import type { PageResponse } from '@/api/pageResult'
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
}

export interface TagListRes extends ApiResponse<TagListData[]> {}

// tag add
export interface TagAddData {
  id?: number
  name?: string
  remark?: string
}

export interface TagAddRes extends ApiResponse<number> {}

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
export interface TagUpdateRes extends ApiResponse<boolean> {}

// tag delete
export interface TagDeleteRes extends ApiResponse<boolean> {}

// tag page
export interface TagPageRes extends PageResponse<TagListData> {}

export interface TagPageParams {
  pageNum: number
  pageSize: number
  name?: string
  remark?: string
}
