export interface ApiResponse<T> {
  code: number
  msg: string
  data: T
}

export interface PageResponse<T>
  extends ApiResponse<{
    total: number
    rows: T[]
  }> {}
