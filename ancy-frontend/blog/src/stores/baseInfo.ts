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

  const getName = () => {
    return baseInfo.value?.name
  }

  const getAddress = () => {
    return baseInfo.value?.address
  }

  const getBadge = () => {
    return baseInfo.value?.badge
  }

  const getFooter = () => {
    return baseInfo.value?.footer
  }

  return {
    baseInfo,
    reqBaseInfo,
    getAvatar,
    getGreeting,
    getRole,
    getPhilosophy,
    getName,
    getAddress,
    getBadge,
    getFooter,
  }
})
