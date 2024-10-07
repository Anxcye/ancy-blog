import request from '@/utils/request'
import type { ReadPageParam, ReadPageRes } from './type'

export function reqReadPage(params: ReadPageParam): Promise<ReadPageRes> {
  return request.get('/read/page', { params })
}
