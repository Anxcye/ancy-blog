import request from '@/utils/request'
import type { SettingGetBaseRes } from './type'

export function reqSettingGetBase(): Promise<SettingGetBaseRes> {
  return request.get(`/setting/base`)
}
