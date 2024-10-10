import type { ApiResponse, PageResponse } from '../ApiResponse'

export interface CommentData {
  articleId: number
  avatar: string
  children: {
    rows: CommentData[]
    total: number
  }
  content: string
  createBy: number
  createTime: string
  deleted: number
  email: string
  id: number
  ip: string
  isTop: '0' | '1'
  likeCount: number
  nickname: string
  parentId: number
  status: '0' | '1'
  toCommentId: number
  toCommentNickname: string
  ua: string
  updateBy: number
  updateTime: string
  userId: number
}

export interface CommentListRes extends PageResponse<CommentData> {}

export interface CommentPageParams {
  id: number
  pageNum: number
  pageSize: number
}
export interface CommentAddParams {
  type?: string
  articleId?: number
  parentId?: number
  avatar?: string
  nickname?: string
  email?: string
  content?: string
  status?: string
  likeCount?: number
  isTop?: string
  toCommentNickname?: string
  toCommentId?: number
}

export interface CommentAddRes extends ApiResponse<number> {}

export interface CommentLikeRes extends ApiResponse<boolean> {}

export interface CommentUpdateRes extends ApiResponse<boolean> {}

export interface CommentTotalRes extends ApiResponse<number> {}
