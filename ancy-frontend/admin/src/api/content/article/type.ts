import type { ApiResponse } from '@/api/type'

// article page
export interface ArticlePageData {
  categoryId: number
  categoryName: string
  createBy: number
  createTime: string
  id: number
  isTop: string
  summary: string
  thumbnail: string
  title: string
  viewCount: number
}

export interface ArticlePageRes extends ApiResponse {
  data: {
    rows: ArticlePageData[]
    total: number
  }
}

export interface ArticlePageParams {
  title?: string
  summary?: string
  pageNum: number
  pageSize: number
}

// article delete
export interface ArticleDeleteRes extends ApiResponse {
  data: boolean
}
export interface ArticleAddParams {
  title?: string
  content?: string
  summary?: string
  categoryId?: number
  thumbnail?: string
  isTop?: string
  status?: string
  isComment?: string
  tags?: number[]
}

// article add
export interface ArticleAddRes extends ApiResponse {
  data: boolean
}
