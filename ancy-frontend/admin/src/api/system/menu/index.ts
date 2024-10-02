import request from '@/utils/request'
import type {
  MenuAddParams,
  MenuAddRes,
  MenuDeleteRes,
  MenuListRes,
  MenuPageParams,
  MenuPageRes,
  MenuTreeRes,
  MenuUpdateRes,
} from './type'

export function reqMenuAdd(params: MenuAddParams): Promise<MenuAddRes> {
  return request.post(`/menus`, params)
}

export function reqMenuUpdate(id: number, params: MenuAddParams): Promise<MenuUpdateRes> {
  return request.put(`/menus/${id}`, params)
}

export function reqMenuDelete(id: number): Promise<MenuDeleteRes> {
  return request.delete(`/menus/${id}`)
}

export function reqMenuList(): Promise<MenuListRes> {
  return request.get(`/menus/list`)
}

export function reqMenuPage(params: MenuPageParams): Promise<MenuPageRes> {
  return request.get(`/menus/page`, { params })
}

export function reqMenuTree(): Promise<MenuTreeRes> {
  return request.get(`/menus/tree`)
}

export function reqMenuListByRoleId(roleId: number): Promise<MenuListRes> {
  return request.get(`/menus/list/role/${roleId}`)
}
