import request from '@/utils/request'
import type { ProjectDetailRes, ProjectListRes } from './type'

export function reqProjectList(): Promise<ProjectListRes> {
  return request.get(`/project/list`)
}

export function reqProjectDetail(id: number): Promise<ProjectDetailRes> {
  return request.get(`/project/${id}`)
}
