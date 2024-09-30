import request from '@/utils/request'
import type {
  CategoryAddParams,
  CategoryAddRes,
  CategoryDeleteRes,
  CategoryListRes,
  CategoryPageParams,
  CategoryPageRes,
  CategoryUpdateRes,
  ExportToXlsxRes,
} from './type'

export function reqCategoryList(): Promise<CategoryListRes> {
  return request.get(`/category/list`)
}

export function reqCategoryAdd(params: CategoryAddParams): Promise<CategoryAddRes> {
  return request.post(`/category`, params)
}

export function reqCategoryPage(params: CategoryPageParams): Promise<CategoryPageRes> {
  const { name, status, pageNum, pageSize } = params
  return request.get(
    `/category/page?${name ? `name=${name}&` : ''}${
      status ? `status=${status}&` : ''
    }pageNum=${pageNum}&pageSize=${pageSize}`,
  )
}

export function reqCategoryUpdate(
  id: number,
  params: CategoryAddParams,
): Promise<CategoryUpdateRes> {
  return request.put(`/category/${id}`, params)
}

export function reqExportToXlsx(): Promise<ExportToXlsxRes> {
  // return request.get(`/category/export/xlsx`)
  return request.get(`/category/export/xlsx`, {
    responseType: 'blob',
  })
}

export function reqCategoryDelete(id: number): Promise<CategoryDeleteRes> {
  return request.delete(`/category/${id}`)
}
