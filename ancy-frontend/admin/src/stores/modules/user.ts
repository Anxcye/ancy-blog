import { reqLogin, reqLogout } from '@/api/user'
import type { LoginData, LoginParams } from '@/api/user/type'

import {
  getLoginInfo,
  removeLoginInfo,
  setLoginInfo,
} from '@/utils/localStorage/login'
import { defineStore } from 'pinia'
import { ref } from 'vue'
import { useRouteStore } from './route'

export const useUserStore = defineStore('user', () => {
  const userInfo = ref<LoginData | null>(getLoginInfo())

  const userLogin = async (data: LoginParams) => {
    const res = await reqLogin(data)
    setLoginInfo(res.data)
    userInfo.value = res.data

    await useRouteStore().initRoutes()
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
