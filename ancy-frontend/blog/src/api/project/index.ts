import request from '@/utils/request'
import type { ProjectDetailRes, ProjectListRes, ProjectPageParams, ProjectPageRes } from './type'

export function reqProjectList(): Promise<ProjectListRes> {
  return request.get(`/project/list`)
}

export function reqProjectDetail(id: number): Promise<ProjectDetailRes> {
  return request.get(`/project/${id}`)
}

export function reqProjectPage(params: ProjectPageParams): Promise<ProjectPageRes> {
  return request.get('/project/page', { params })
}
