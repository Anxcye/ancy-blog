import type { ApiResponse } from './type'

export interface PageResponse<T>
  extends ApiResponse<{
    rows: T[]
    total: number
  }> {}
