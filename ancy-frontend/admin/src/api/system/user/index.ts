import request from '@/utils/request'
import type {
  GetRoutersRes,
  LoginParams,
  LoginRes,
  LogoutRes,
  UserAddParams,
  UserAddRes,
  UserDeleteRes,
  UserGetByIdRes,
  UserPageParams,
  UserPageRes,
  UserUpdateRes,
} from './type'

export function reqLogin(params: LoginParams): Promise<LoginRes> {
  return request.post(`/user/login`, params)
}

export function getRouters(): Promise<GetRoutersRes> {
  return request.get(`/user/routers`)
}

export function reqLogout(): Promise<LogoutRes> {
  return request.post(`/user/logout`)
}

export function reqUserAdd(params: UserAddParams): Promise<UserAddRes> {
  return request.post(`/user`, params)
}

export function reqUserGetById(id: number): Promise<UserGetByIdRes> {
  return request.get(`/user/${id}`)
}

export function reqUserUpdate(id: number, params: UserAddParams): Promise<UserUpdateRes> {
  return request.put(`/user/${id}`, params)
}

export function reqUserDelete(id: number): Promise<UserDeleteRes> {
  return request.delete(`/user/${id}`)
}

export function reqUserPage(params: UserPageParams): Promise<UserPageRes> {
  return request.get(`/user/page`, { params })
}
