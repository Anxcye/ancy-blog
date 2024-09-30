import type { ApiResponse } from '@/api/type'
import type { TagData } from '../tag/type'

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
  isComment: string
  createBy: number
  createTime: string
  updateBy: number
  updateTime: string
}

export interface ArticleGetByIdRes extends ApiResponse {
  code: number
  msg: string
  data: ArticleGetByIdData
}

// article update
export interface ArticleUpdateRes extends ApiResponse {
  data: boolean
}
