import type { PageResponse } from '../ApiResponse'

export interface NoteData {
  createTime: string
  deleted: number
  id: number
  isComment: '1' | '0'
  isTop: '1' | '0'
  orderNum: number
  viewCount: 0
}
export interface NotePageRes extends PageResponse<NoteData> {}

export interface NotePageParams {
  pageNum: number
  pageSize: number
}
