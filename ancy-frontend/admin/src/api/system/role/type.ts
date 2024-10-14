import type { PageResponse } from '@/api/pageResult'
import type { ApiResponse } from '@/api/type'

// role add
export interface RoleAddParams {
  id?: number
  roleName?: string
  roleKey?: string
  roleSort?: number
  menuIds?: number[]
  status?: string
  remark?: string
}

export interface RoleAddRes extends ApiResponse<boolean> {}

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

export interface RoleGetByIdRes extends ApiResponse<RoleGetByIdData> {}

export interface RoleUpdateRes extends ApiResponse<boolean> {}

export interface RoleDeleteRes extends ApiResponse<boolean> {}

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

export interface RoleListRes extends ApiResponse<RoleListData[]> {}

export interface RolePageParams {
  name?: string
  status?: string
  pageNum: number
  pageSize: number
}

export interface RolePageRes extends PageResponse<RoleListData> {}
