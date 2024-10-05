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

  const getAvatar = () => {
    return baseInfo.value?.avatar
  }

  const getGreeting = () => {
    return baseInfo.value?.greeting
  }

  const getRole = () => {
    return baseInfo.value?.role
  }

  const getPhilosophy = () => {
    return baseInfo.value?.philosophy
  }

  const getBadge = () => {
    return baseInfo.value?.badge
  }

  return {
    baseInfo,
    reqBaseInfo,
    getAvatar,
    getGreeting,
    getRole,
    getPhilosophy,
    getBadge,
  }
})
