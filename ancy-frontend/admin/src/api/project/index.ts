import request from '@/utils/request'
import type {
  ProjectAddParams,
  ProjectAddRes,
  ProjectDeleteRes,
  ProjectGetByIdRes,
  ProjectPageParams,
  ProjectPageRes,
  ProjectUpdateRes,
} from './type'

export function reqProjectPage(params: ProjectPageParams): Promise<ProjectPageRes> {
  return request.get(`/project/page`, { params })
}

export function reqProjectAdd(params: ProjectAddParams): Promise<ProjectAddRes> {
  return request.post(`/project`, params)
}

export function reqProjectUpdate(id: number, params: ProjectAddParams): Promise<ProjectUpdateRes> {
  return request.put(`/project/${id}`, params)
}

export function reqProjectDelete(id: number): Promise<ProjectDeleteRes> {
  return request.delete(`/project/${id}`)
}

export function reqProjectGetById(id: number): Promise<ProjectGetByIdRes> {
  return request.get(`/project/${id}`)
}
