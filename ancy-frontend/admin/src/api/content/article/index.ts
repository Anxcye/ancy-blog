import request from '@/utils/request'
import type { ArticlePageParams, ArticlePageRes, ArticleDelete } from './type'

export function reqArticlePage(
  params: ArticlePageParams,
): Promise<ArticlePageRes> {
  const { title, summary, pageNum, pageSize } = params
  return request.get(
    `/article/page?title=${title}&summary=${summary}&pageNum=${pageNum}&pageSize=${pageSize}`,
  )
}

export function reqArticleDelete(id: number): Promise<ArticleDelete> {
  return request.delete(`/article/${id}`)
}
