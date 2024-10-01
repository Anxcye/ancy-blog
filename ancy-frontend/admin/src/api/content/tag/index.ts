import request from '@/utils/request'
import type {
  TagAddData,
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
export function reqTagAdd(params: TagAddData): Promise<TagAddRes> {
  return request.post(`/tag`, params)
}

export function reqTagUpdate(id: number, params: TagAddData): Promise<TagUpdateRes> {
  return request.put(`/tag/${id}`, params)
}
export function reqTagDelete(id: number): Promise<TagDeleteRes> {
  return request.delete(`/tag/${id}`)
}

export function reqTagPage(params: TagPageParams): Promise<TagPageRes> {
  return request.get(`/tag/page`, { params })
}
