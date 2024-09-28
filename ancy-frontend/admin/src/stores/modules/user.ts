import { reqLogin } from '@/api/user'
import type { LoginData, LoginParams } from '@/api/user/type'
import { getLoginInfo, setLoginInfo } from '@/utils/localStorage/login'
import { defineStore } from 'pinia'
import { ref } from 'vue'
import { useRouteStore } from './route'

export const useUserStore = defineStore('user', () => {
  const userInfo = ref<LoginData>(getLoginInfo())

  const userLogin = async (data: LoginParams) => {
    const res = await reqLogin(data)
    setLoginInfo(res.data)
    userInfo.value = res.data
    await useRouteStore().initRoutes()
  }

  // getter token
  const getToken = (): string => {
    return userInfo.value?.token || ''
  }

  return {
    userInfo,
    userLogin,
    getToken,
  }
})
