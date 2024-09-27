import request from '@/utils/request'
import type { GetRoutersRes, LoginParams, LoginRes } from './type'

export function reqLogin(params: LoginParams): Promise<LoginRes> {
  return request.post(`/user/login`, params)
}

export function getRouters(): Promise<GetRoutersRes> {
  return request.get(`/user/routers`)
}
