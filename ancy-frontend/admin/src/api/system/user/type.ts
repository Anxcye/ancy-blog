import type { PageResponse } from '@/api/pageResult'
import type { ApiResponse } from '@/api/type'

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
export interface LoginRes extends ApiResponse<LoginData> {}

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
export interface GetRoutersRes
  extends ApiResponse<{
    menus: GetRoutersData[]
  }> {}

// Response interface
export interface LogoutRes extends ApiResponse<boolean> {}

// User add
export interface UserAddParams {
  id?: number
  userName?: string
  nickName?: string
  password?: string
  status?: string
  email?: string
  roleIds?: number[]
  phonenumber?: string
  sex?: string
  avatar?: string
}

export interface UserAddRes extends ApiResponse<number> {}

// User get by id
export interface UserGetByIdData {
  id: number
  userName: string
  nickName: string
  type: string
  status: string
  roleIds: number[]
  email: string
  sex: string
  avatar: string
  createBy: number
  createTime: string
  updateBy: number
  updateTime: string
}

export interface UserGetByIdRes extends ApiResponse<UserGetByIdData> {}

// User update
export interface UserUpdateRes extends ApiResponse<boolean> {}

// User delete
export interface UserDeleteRes extends ApiResponse<boolean> {}

// user list
export interface UserListData {
  avatar: string
  createBy: number
  createTime: string
  email: string
  id: number
  nickName: string
  roleIds: number[]
  sex: '0' | '1' | '2'
  status: '0' | '1'
  type: '0' | '1'
  updateBy: number
  updateTime: string
  userName: string
}

export interface UserPageRes extends PageResponse<UserListData> {}

export interface UserPageParams {
  pageNum: number
  pageSize: number
  userName?: string
  nickName?: string
  type?: string
  status?: string
  email?: string
  phonenumber?: string
  sex?: string
}
