import type { ApiResponse } from '@/api/type'

export interface MenuAddParams {
  id?: number
  menuName?: string
  parentId?: number
  orderNum?: number
  path?: string
  component?: string
  isFrame?: number
  menuType?: string
  visible?: string
  status?: string
  perms?: string
  icon?: string
  remark?: string
}

export interface MenuAddRes extends ApiResponse<boolean> {}

export interface MenuUpdateRes extends ApiResponse<boolean> {}

export interface MenuDeleteRes extends ApiResponse<boolean> {}

export interface MenuListData {
  id: number
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
  createBy: number
  createTime: string
  updateBy: number
  updateTime: string
  remark: string
  children: MenuListData[]
}

export interface MenuListRes extends ApiResponse<MenuListData[]> {}

export interface MenuPageParams {
  name?: string
  status?: string
  pageNum?: number
  pageSize?: number
}

export interface MenuPageRes {
  code: number
  msg: string
  data: {
    total: number
    rows: MenuListData[]
  }
}

export interface MenuTreeRes extends ApiResponse<MenuListData[]> {}
