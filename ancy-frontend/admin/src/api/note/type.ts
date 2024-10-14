import type { PageResponse } from '../pageResult'
import type { ApiResponse } from '../type'

export interface NotePageData {
  content: string
  createBy: number
  createTime: string
  id: number
  isComment: '0' | '1'
  isTop: '0' | '1'
  orderNum: number
  status: '0' | '1'
  updateBy: string | null
  updateTime: string | null
  viewCount: number
}

export interface NotePageRes extends PageResponse<NotePageData> {}

export interface NotePageParams {
  content?: string
  status?: string
  pageNum: number
  pageSize: number
}

export interface NoteAddParams {
  content?: string
  isTop?: '0' | '1'
  status?: '0' | '1'
  orderNum?: number
  viewCount?: number
  isComment?: '0' | '1'
}

export interface NoteAddRes extends ApiResponse<number> {}

export interface NoteData {
  id: number
  content: string
  isTop: '0' | '1'
  status: '0' | '1'
  orderNum: number
  viewCount: number
  isComment: '0' | '1'
  createBy: number
  createTime: string
  updateBy: number
  updateTime: string
}

export interface NoteGetByIdRes extends ApiResponse<NoteData> {}

export interface NoteUpdateRes extends ApiResponse<boolean> {}

export interface NoteDeleteRes extends ApiResponse<boolean> {}
