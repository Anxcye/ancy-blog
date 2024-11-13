import type { PageResponse } from '../pageResult'
import type { ApiResponse } from '../type'

export interface ProjectListData {
  beginDate: string
  createBy: number
  createTime: string
  displayUrl: string
  id: number
  isTop: '0' | '1'
  orderNum: number
  srcUrl: string
  status: '0' | '1'
  summary: string
  thumbnail: string
  title: string
  type: string
  updateBy: string | null
  updateTime: string | null
}

export interface ProjectPageRes extends PageResponse<ProjectListData> {}

export interface ProjectPageParams {
  title?: string
  summary?: string
  status?: string
  type?: string
  pageNum: number
  pageSize: number
}

export interface ProjectAddParams {
  id?: number
  title?: string
  content?: string
  summary?: string
  thumbnail?: string
  isTop?: '0' | '1'
  status?: '0' | '1'
  type?: string
  srcUrl?: string
  displayUrl?: string
  orderNum?: number
  beginDate?: string
}

export interface ProjectAddRes extends ApiResponse<number> {}

export interface ProjectUpdateRes extends ApiResponse<boolean> {}

export interface ProjectDeleteRes extends ApiResponse<boolean> {}

export interface ProjectDetailData {
  id: number
  title: string
  content: string
  summary: string
  thumbnail: string
  isTop: '0' | '1'
  status: '0' | '1'
  type: string
  srcUrl: string
  displayUrl: string
  orderNum: number
  beginDate: string
  createBy: number
  createTime: string
  updateBy: number
  updateTime: string
}

export interface ProjectGetByIdRes extends ApiResponse<ProjectDetailData> {}
