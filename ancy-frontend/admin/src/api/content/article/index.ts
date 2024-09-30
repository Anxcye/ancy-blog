import request from '@/utils/request'
import type {
  ArticlePageParams,
  ArticlePageRes,
  ArticleDeleteRes,
  ArticleAddParams,
  ArticleAddRes,
  ArticleGetByIdRes,
  ArticleUpdateRes,
} from './type'

export function reqArticlePage(params: ArticlePageParams): Promise<ArticlePageRes> {
  const { title, summary, pageNum, pageSize } = params
  return request.get(
    `/article/page?${title ? `title=${title}&` : ''}${
      summary ? `summary=${summary}&` : ''
    }pageNum=${pageNum}&pageSize=${pageSize}`,
  )
}

export function reqArticleDelete(id: number): Promise<ArticleDeleteRes> {
  return request.delete(`/article/${id}`)
}

export function reqArticleAdd(params: ArticleAddParams): Promise<ArticleAddRes> {
  return request.post(`/article`, params)
}

export function reqArticleGetById(id: number): Promise<ArticleGetByIdRes> {
  return request.get(`/article/${id}`)
}

export function reqArticleUpdate(id: number, params: ArticleAddParams): Promise<ArticleUpdateRes> {
  return request.put(`/article/${id}`, params)
}
