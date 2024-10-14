import type { ApiResponse } from '../type'

export interface Badge {
  index?: string
  title?: string
  url?: string
  img?: string
  orderNum?: number
}

export interface Footer {
  position?: number
  index?: string
  text?: string
  url?: string
  orderNum?: number
}

export interface SettingData {
  avatar?: string
  greeting?: string
  role?: string
  philosophy?: string
  name?: string
  address?: string
  badge?: Badge[]
  footer?: Footer[]
}

export interface SettingListRes extends ApiResponse<SettingData> {}

export interface SettingUpdateParams {
  avatar?: string
  greeting?: string
  role?: string
  philosophy?: string
  name?: string
  address?: string
  badge?: {
    index?: string
    title?: string
    url?: string
    img?: string
    orderNum?: number
  }[]
  footer?: {
    position?: number
    index?: string
    text?: string
    url?: string
    orderNum?: number
  }[]
}

export interface SettingUpdateRes extends ApiResponse<boolean> {}

export interface SettingDeleteRes extends ApiResponse<boolean> {}
