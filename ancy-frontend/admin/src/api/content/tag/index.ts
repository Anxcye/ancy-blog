import request from '@/utils/request'
import type { TagAddParams, TagAddRes, TagListRes } from './type'

export function reqTagList(): Promise<TagListRes> {
  return request.get(`/tag/list`)
}
export function reqTagAdd(params: TagAddParams): Promise<TagAddRes> {
  return request.post(`/tag`, params)
}
