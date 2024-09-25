// Parameter interface
export interface LoginParams {
  userName?: string
  nickName?: string
  email?: string
  password?: string
}

// Response interface
export interface LoginRes {
  code: number
  msg: string
  data: {
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
  }
}
