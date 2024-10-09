import request from '@/utils/request'
import type { LoginParams, LoginRes } from './type'

export function reqLogin(params: LoginParams): Promise<LoginRes> {
  return request.post(`/user/login`, params)
}
