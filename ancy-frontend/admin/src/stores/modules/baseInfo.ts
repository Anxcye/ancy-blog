import { reqSettingList } from '@/api/setting'
import type { SettingData } from '@/api/setting/type'
import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useBaseInfoStore = defineStore('baseInfo', () => {
  const baseInfo = ref<SettingData>(JSON.parse(localStorage.getItem('baseInfo') || '{}'))

  const getBaseInfo = async () => {
    const res = await reqSettingList()

    localStorage.setItem('baseInfo', JSON.stringify(res.data))
    baseInfo.value = res.data
  }
  return { baseInfo, getBaseInfo }
})
