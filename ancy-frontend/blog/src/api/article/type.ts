import type { ApiResponse } from '../ApiResponse'
import type { TagListData } from '../tag/type'

export interface ArticleListData {
  id: number
  title: string
  summary: string
  categoryName: string
  categoryId: number
  thumbnail: string
  isTop: string
  status: string
  type: number
  orderNum: number
  tags: TagListData[]
  viewCount: number
  createBy: number
  createTime: string
  updateTime: string
}

export interface ArticleGetFrontListRes extends ApiResponse<ArticleListData[]> {}

export interface ArticlePageRes
  extends ApiResponse<{
    total: number
    rows: ArticleListData[]
  }> {}

export interface ArticlePageParams {
  pageNum: number
  pageSize: number
  categoryId: number
}

export interface ArticleDetailData {
  id: number
  title: string
  content: string
  summary: string
  categoryId: number
  categoryName: string
  thumbnail: string
  viewCount: number
  tags: TagListData[]
  isComment: string
  isTop: string
  createBy: number
  createTime: string
  updateBy: number
  updateTime: string
}

export interface ArticleGetByIdRes extends ApiResponse<ArticleDetailData> {}
