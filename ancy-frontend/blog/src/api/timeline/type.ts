import type { PageParam, PageResponse } from '../ApiResponse'

export interface TimelineData {
  id: number
  methodName: string
  operateTime: string
  returnValue: number
  summary: string
}

export interface TimelinePageRes extends PageResponse<TimelineData> {}

export interface TimelinePageParam extends PageParam {}
