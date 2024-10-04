import { reqSettingGetBase } from '@/api/baseInfo'
import type { SettingGetBaseData } from '@/api/baseInfo/type'
import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useBaseInfoStore = defineStore('baseInfo', () => {
  const baseInfo = ref<SettingGetBaseData>()

  const reqBaseInfo = async () => {
    const res = await reqSettingGetBase()
    baseInfo.value = res.data
  }

  return { baseInfo, reqBaseInfo }
})
