// Parameter interface
export interface LoginParams {
  userName?: string
  nickName?: string
  email?: string
  password?: string
}

export interface LoginData {
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

// Response interface
export interface LoginRes {
  code: number
  msg: string
  data: LoginData
}

export interface GetRoutersData {
  id: number

  children: GetRoutersData[]
  menuName: string
  parentId: number
  orderNum: number
  path: string
  component: string
  isFrame: number
  menuType: string
  visible: string
  status: string
  perms: string
  icon: string
  remark: string

  createBy: number
  createTime: string
  updateBy: number
  updateTime: string
}

// Response interface
export interface GetRoutersRes {
  code: number
  msg: string

  data: {
    menus: GetRoutersData[]
  }
}
// Response interface
export interface LogoutRes {
  code: number
  msg: string
  data: boolean
}
