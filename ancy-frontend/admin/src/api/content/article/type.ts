import type { Response } from '@/api/type'

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

export interface ArticlePageRes extends Response {
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

export interface ArticleDelete extends Response {
  data: boolean
}
