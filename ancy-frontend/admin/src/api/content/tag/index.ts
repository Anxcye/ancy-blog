import request from '@/utils/request'
import type {
  TagAddParams,
  TagAddRes,
  TagDeleteRes,
  TagListRes,
  TagPageParams,
  TagPageRes,
  TagUpdateRes,
} from './type'

export function reqTagList(): Promise<TagListRes> {
  return request.get(`/tag/list`)
}
export function reqTagAdd(params: TagAddParams): Promise<TagAddRes> {
  return request.post(`/tag`, params)
}
// Parameter interface
export interface TagUpdateParams {
  name?: string

  remark?: string
}

export function tagUpdate(id: number, params: TagUpdateParams): Promise<TagUpdateRes> {
  return request.put(`/tag/${id}`, params)
}
export function tagDelete(id: number): Promise<TagDeleteRes> {
  return request.delete(`/tag/${id}`)
}

export function tagPage(params: TagPageParams): Promise<TagPageRes> {
  const { pageNum, pageSize, name, remark } = params
  return request.get(
    `/tag/page?${name ? `name=${name}&` : ''}${remark ? `remark=${remark}&` : ''}pageNum=${pageNum}&pageSize=${pageSize}`,
  )
}
