import request from '@/utils/request'
import type {
  ArticleGetByIdRes,
  ArticleGetFrontListRes,
  ArticleListRes,
  ArticlePageParams,
  ArticlePageRes,
} from './type'

export function reqArticleGetFrontList(): Promise<ArticleGetFrontListRes> {
  return request.get(`/article/front`)
}

export function reqArticlePage(params: ArticlePageParams): Promise<ArticlePageRes> {
  return request.get('/article/page', { params })
}

export function reqArticleHot(): Promise<ArticlePageRes> {
  return request.get(`/article/hot`)
}

export function reqArticleGetById(id: number): Promise<ArticleGetByIdRes> {
  return request.get(`/article/${id}`)
}

export function reqArticleRecent(): Promise<ArticleListRes> {
  return request.get(`/article/recent`)
}
