import request from '@/utils/request'
import type {
  ArticlePageParams,
  ArticlePageRes,
  ArticleDeleteRes,
  ArticleAddParams,
  ArticleAddRes,
} from './type'

export function reqArticlePage(
  params: ArticlePageParams,
): Promise<ArticlePageRes> {
  const { title, summary, pageNum, pageSize } = params
  return request.get(
    `/article/page?title=${title}&summary=${summary}&pageNum=${pageNum}&pageSize=${pageSize}`,
  )
}

export function reqArticleDelete(id: number): Promise<ArticleDeleteRes> {
  return request.delete(`/article/${id}`)
}

export function articleAdd(params: ArticleAddParams): Promise<ArticleAddRes> {
  return request.post(`/article`, params)
}
