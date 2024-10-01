import type { ApiResponse } from '@/api/type'

// role add
export interface RoleAddParams {
  roleName?: string
  roleKey?: string
  roleSort?: number
  menuIds?: number[]
  status?: string
  remark?: string
}

export interface RoleAddRes extends ApiResponse {
  data: boolean
}

// role/get by id
export interface RoleGetByIdData {
  id: number
  roleName: string
  roleKey: string
  roleSort: number
  menuIds: number[]
  status: string
  createBy: number
  createTime: string
  updateBy: number
  updateTime: string
  remark: string
}

export interface RoleGetByIdRes extends ApiResponse {
  data: RoleGetByIdData
}

export interface RoleUpdateRes extends ApiResponse {
  data: boolean
}

export interface RoleDeleteRes extends ApiResponse {
  data: boolean
}

export interface RoleListData {
  id: number
  roleName: string
  roleKey: string
  roleSort: number
  menuIds: number[]
  status: '0' | '1'
  createBy: number
  createTime: string
  updateBy: number
  updateTime: string
  remark: string
}

export interface RoleListRes extends ApiResponse {
  data: RoleListData[]
}

export interface RolePageParams {
  name: string
  status: string
  pageNum: number
  pageSize: number
}

export interface RolePageRes extends ApiResponse {
  data: {
    total: number
    rows: RoleListData[]
  }
}
