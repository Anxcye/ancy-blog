import { reqLogin, reqLogout } from '@/api/system/user'
import type { LoginData, LoginParams } from '@/api/system/user/type'

import { getLoginInfo, removeLoginInfo, setLoginInfo } from '@/utils/localStorage/login'
import { defineStore } from 'pinia'
import { ref } from 'vue'
import { useRouteStore } from './route'
import { useBaseInfoStore } from './baseInfo'

export const useUserStore = defineStore('user', () => {
  const userInfo = ref<LoginData | null>(getLoginInfo())

  const userLogin = async (data: LoginParams) => {
    const res = await reqLogin(data)
    setLoginInfo(res.data)
    userInfo.value = res.data

    await useRouteStore().initRoutes()
    await useBaseInfoStore().getBaseInfo()
  }

  const logout = async () => {
    await reqLogout()
    removeLoginInfo()
    userInfo.value = null
    useRouteStore().removeRoutes()
  }

  // getter token
  const getToken = (): string | null => {
    return userInfo.value?.token || null
  }

  return {
    userInfo,
    userLogin,
    logout,
    getToken,
  }
})
