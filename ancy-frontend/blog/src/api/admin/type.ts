import type { ApiResponse } from '../ApiResponse'

export interface LoginParams {
  userName?: string
  nickName?: string
  email?: string
  password?: string
}

export interface LoginRes
  extends ApiResponse<{
    token: string
    permissions: string[]
    role: string[]
    userInfoVo: {
      id: number
      nickName: string
      avatar: string
      sex: string
      email: string
    }
  }> {}
