import request from '@/utils/request'
import type { SettingListRes, SettingUpdateParams, SettingUpdateRes } from './type'

export function reqSettingList(): Promise<SettingListRes> {
  return request.get(`/setting`)
}

export function reqSettingUpdate(params: SettingUpdateParams): Promise<SettingUpdateRes> {
  return request.put(`/setting`, params)
}
