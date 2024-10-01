import request from '@/utils/request'
import type {
  RoleAddParams,
  RoleAddRes,
  RoleDeleteRes,
  RoleGetByIdRes,
  RoleListRes,
  RolePageParams,
  RolePageRes,
  RoleUpdateRes,
} from './type'

export function reqRoleAdd(params: RoleAddParams): Promise<RoleAddRes> {
  return request.post(`/role`, params)
}

export function reqRoleGetById(id: number): Promise<RoleGetByIdRes> {
  return request.get(`/role/${id}`)
}

export function reqRoleUpdate(id: number, params: RoleAddParams): Promise<RoleUpdateRes> {
  return request.put(`/role/${id}`, params)
}

export function reqRoleDelete(id: number): Promise<RoleDeleteRes> {
  return request.delete(`/role/${id}`)
}

export function reqRoleList(): Promise<RoleListRes> {
  return request.get(`/role/list`)
}

export function reqRolePage(params: RolePageParams): Promise<RolePageRes> {
  return request.get(`/role/page`, { params })
}
