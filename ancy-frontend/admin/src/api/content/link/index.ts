import request from '@/utils/request'
import type {
  LinkAddParams,
  LinkAddRes,
  LinkDeleteRes,
  LinkGetByIdRes,
  LinkListRes,
  LinkPageParams,
  LinkPageRes,
  LinkUpdateArticleRes,
  LinkUpdateRes,
} from './type'
import type { ArticleAddParams, ArticleGetByIdRes } from '../article/type'

export function reqLinkAdd(params: LinkAddParams): Promise<LinkAddRes> {
  return request.post(`/links`, params)
}

export function reqLinkGetById(id: number): Promise<LinkGetByIdRes> {
  return request.get(`/links/${id}`)
}

export function reqLinkUpdate(id: number, params: LinkAddParams): Promise<LinkUpdateRes> {
  return request.put(`/links/${id}`, params)
}

export function reqLinkDelete(id: number): Promise<LinkDeleteRes> {
  return request.delete(`/links/${id}`)
}

export function reqLinkList(): Promise<LinkListRes> {
  return request.get(`/links/list`)
}

export function reqLinkPage(params: LinkPageParams): Promise<LinkPageRes> {
  return request.get(`/links/page`, { params })
}

export function reqLinkUpdateArticle(params: ArticleAddParams): Promise<LinkUpdateArticleRes> {
  return request.put(`/links/article`, params)
}

export function reqLinkGetArticle(): Promise<ArticleGetByIdRes> {
  return request.get(`/links/article`)
}
