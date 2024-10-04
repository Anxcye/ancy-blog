import request from '@/utils/request'
import type { CategoryListRes } from './type'

export function reqCategoryList(): Promise<CategoryListRes> {
  return request.get(`/category/list`)
}
