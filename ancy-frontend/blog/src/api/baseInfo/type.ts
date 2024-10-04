import type { ApiResponse } from '../ApiResponse'

export interface Badge {
  img: string
  index: string
  title: string
  url: string
}

export interface SettingGetBaseData {
  avatar: string
  badge: Badge[]
  greeting: string
  philosophy: string
  role: string
}
export interface SettingGetBaseRes extends ApiResponse<SettingGetBaseData> {}
