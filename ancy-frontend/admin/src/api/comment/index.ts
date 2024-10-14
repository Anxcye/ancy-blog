import request from '@/utils/request'
import type {
  CommentAddParams,
  CommentAddRes,
  CommentDeleteRes,
  CommentPageParams,
  CommentPageRes,
  CommentUpdateRes,
} from './type'

export function reqCommentPage(params: CommentPageParams): Promise<CommentPageRes> {
  return request.get(`/comment/page`, { params })
}

export function reqCommentAdd(params: CommentAddParams): Promise<CommentAddRes> {
  return request.post(`/comment`, params)
}

export function reqCommentUpdate(id: number, params: CommentAddParams): Promise<CommentUpdateRes> {
  return request.put(`/comment/${id}`, params)
}

export function reqCommentDelete(id: number): Promise<CommentDeleteRes> {
  return request.delete(`/comment/${id}`)
}
