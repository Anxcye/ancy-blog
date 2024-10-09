import type { ApiResponse } from '../ApiResponse'

export interface Badge {
  img: string
  index: string
  title: string
  url: string
}

export interface Footer {
  index: string
  text: string
  url: string
  position: number
}

export interface SettingGetBaseData {
  avatar: string
  badge: Badge[]
  footer: Footer[]
  greeting: string
  philosophy: string
  name: string
  address: string
  role: string
}
export interface SettingGetBaseRes extends ApiResponse<SettingGetBaseData> {}
