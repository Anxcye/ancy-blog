import request from '@/utils/request'
import type { LinkAddParams, LinkAddRes, LinkListRes } from './type'
import type { ArticleGetByIdRes } from '../article/type'

export function reqLinkList(): Promise<LinkListRes> {
  return request.get(`/link/list`)
}

export function reqLinkGetArticle(): Promise<ArticleGetByIdRes> {
  return request.get(`/link/article`)
}

export function linkAdd(params: LinkAddParams): Promise<LinkAddRes> {
  return request.post(`/link`, params)
}
