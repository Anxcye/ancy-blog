import type { ApiResponse } from '@/api/type'
import type { TagData } from '../tag/type'
import type { PageResponse } from '@/api/pageResult'

// article page
export interface ArticlePageData {
  categoryId: number
  categoryName: string
  createBy: number
  createTime: string
  updateTime: string
  id: number
  isTop: string
  status: '1' | '0'
  summary: string
  thumbnail: string
  title: string
  viewCount: number
}

export interface ArticlePageRes extends PageResponse<ArticlePageData> {}

export interface ArticlePageParams {
  title?: string
  summary?: string
  pageNum: number
  pageSize: number
}

// article delete
export interface ArticleDeleteRes extends ApiResponse<boolean> {}

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
export interface ArticleAddRes extends ApiResponse<number> {}

// article get by id
export interface ArticleGetByIdData {
  id: number
  title: string
  content: string
  summary: string
  categoryId: number
  categoryName: string
  thumbnail: string
  viewCount: number
  tags: TagData[]
  isComment: '0' | '1'
  isTop: '0' | '1'
  createBy: number
  createTime: string
  updateBy: number
  updateTime: string
}

export interface ArticleGetByIdRes extends ApiResponse<ArticleGetByIdData> {}

// article update
export interface ArticleUpdateRes extends ApiResponse<boolean> {}
