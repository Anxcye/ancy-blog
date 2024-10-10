import request from '@/utils/request'
import type {
  CommentAddParams,
  CommentAddRes,
  CommentLikeRes,
  CommentPageParams,
  CommentUpdateRes,
  CommentListRes,
  CommentTotalRes,
} from './type'

export function reqCommentByArticleId(params: CommentPageParams): Promise<CommentListRes> {
  const { id, pageNum, pageSize } = params
  return request.get(`/comment/article/${id}?pageNum=${pageNum}&pageSize=${pageSize}`)
}

export function reqCommentNote(params: CommentPageParams): Promise<CommentListRes> {
  const { id, pageNum, pageSize } = params
  return request.get(`/comment/note/${id}?pageNum=${pageNum}&pageSize=${pageSize}`)
}

export function reqCommentAdd(params: CommentAddParams): Promise<CommentAddRes> {
  return request.post(`/comment`, params)
}

export function reqCommentLike(id: number, increase: boolean): Promise<CommentLikeRes> {
  return request.put(`/comment/${id}/like?increase=${increase}`)
}

export function reqCommentChildrenByParentId(params: CommentPageParams): Promise<CommentListRes> {
  const { id, pageNum, pageSize } = params
  return request.get(`/comment/${id}/children?pageNum=${pageNum}&pageSize=${pageSize}`)
}

export function reqCommentUpdate(id: number, params: CommentAddParams): Promise<CommentUpdateRes> {
  return request.put(`/comment/admin/${id}`, params)
}

export function commentArticleTotal(id: number): Promise<CommentTotalRes> {
  return request.get(`/comment/article/${id}/total`)
}
export function commentNoteTotal(id: number): Promise<CommentTotalRes> {
  return request.get(`/comment/note/${id}/total`)
}
