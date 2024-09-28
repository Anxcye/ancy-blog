import request from '@/utils/request'
import type { GetRoutersRes, LoginParams, LoginRes, LogoutRes } from './type'

export function reqLogin(params: LoginParams): Promise<LoginRes> {
  return request.post(`/user/login`, params)
}

export function getRouters(): Promise<GetRoutersRes> {
  return request.get(`/user/routers`)
}

export function reqLogout(): Promise<LogoutRes> {
  return request.post(`/user/logout`)
}
