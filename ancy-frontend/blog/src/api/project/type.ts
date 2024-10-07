import type { ApiResponse, PageResponse } from '../ApiResponse'
import type { ArticleDetailData } from '../article/type'

export interface ProjectData {
  id: number
  title: string
  articleId: number
  summary: string
  thumbnail: string
  isTop: string
  status: string
  type: string
  srcUrl: string
  displayUrl: string
  orderNum: number
  beginDate: string
  createBy: number
  createTime: string
  updateBy: number
  updateTime: string
  deleted: number
}

export interface ProjectListRes extends ApiResponse<ProjectData[]> {}

export interface ProjectDetailData extends ProjectData {
  articleDetailVo: ArticleDetailData
}
export interface ProjectDetailRes extends ApiResponse<ProjectDetailData> {}

export interface ProjectPageParams {
  pageNum: number
  pageSize: number
}

export interface ProjectPageRes extends PageResponse<ProjectData> {}
