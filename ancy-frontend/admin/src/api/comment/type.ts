import type { PageResponse } from '../pageResult'
import type { ApiResponse } from '../type'

export interface CommentListData {
  articleId: number
  avatar: string
  children: {
    rows: []
    total: number
  }
  type: '0' | '1'
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
  toCommentId: number | null
  toCommentNickname: string | null
  ua: string
  updateBy: number | null
  updateTime: string | null
  userId: number
}

export interface CommentPageRes extends PageResponse<CommentListData> {}

export interface CommentPageParams {
  articleId?: number
  email?: string
  nickname?: string
  content?: string
  status?: string
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
  status?: '0' | '1'
  likeCount?: number
  isTop?: '0' | '1'
  toCommentNickname?: string
  toCommentId?: number
}

export interface CommentAddRes extends ApiResponse<number> {}

export interface CommentUpdateRes extends ApiResponse<boolean> {}

export interface CommentDeleteRes extends ApiResponse<boolean> {}
