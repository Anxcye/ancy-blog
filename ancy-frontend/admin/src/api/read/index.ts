import request from '@/utils/request'
import type {
  ReadAddParams,
  ReadAddRes,
  ReadDeleteRes,
  ReadPageParams,
  ReadPageRes,
  ReadUpdateRes,
} from './type'

export function reqReadPage(params: ReadPageParams): Promise<ReadPageRes> {
  return request.get(`/read/page`, { params })
}

export function reqReadAdd(params: ReadAddParams): Promise<ReadAddRes> {
  return request.post(`/read`, params)
}

export function reqReadUpdate(id: number, params: ReadAddParams): Promise<ReadUpdateRes> {
  return request.put(`/read/${id}`, params)
}

export function reqReadDelete(id: number): Promise<ReadDeleteRes> {
  return request.delete(`/read/${id}`)
}
