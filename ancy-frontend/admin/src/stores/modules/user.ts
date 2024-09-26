import { reqLogin } from '@/api/user'
import type { LoginParams, LoginRes } from '@/api/user/type'
import { getLoginInfo, setLoginInfo } from '@/utils/localStorage/login'
import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useUserStore = defineStore('user', () => {
  const userInfo = ref<LoginRes>(getLoginInfo())

  const userLogin = async (data: LoginParams) => {
    const res = await reqLogin(data)
    setLoginInfo(res)
    userInfo.value = res
  }

  // getter token
  const getToken = (): string => {
    return userInfo.value?.data.token || ''
  }

  return {
    userInfo,
    userLogin,
    getToken,
  }
})
