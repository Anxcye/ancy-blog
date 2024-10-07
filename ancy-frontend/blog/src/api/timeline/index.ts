import request from '@/utils/request'
import type { TimelinePageParam, TimelinePageRes } from './type'

export function reqTimelinePage(params: TimelinePageParam): Promise<TimelinePageRes> {
  return request.get('/timeline/page', { params })
}
