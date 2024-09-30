import request from '@/utils/request'
import type { CategoryAddParams, CategoryAddRes, CategoryListRes } from './type'

export function reqCategoryList(): Promise<CategoryListRes> {
  return request.get(`/category/list`)
}

export function reqCategoryAdd(
  params: CategoryAddParams,
): Promise<CategoryAddRes> {
  return request.post(`/category`, params)
}
