import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { LoginParams } from '@/api/admin/type'
import { reqLogin } from '@/api/admin'

export const useAdminStore = defineStore('admin', () => {
  const token = ref(localStorage.getItem('token') || '')

  const login = async (params: LoginParams) => {
    const res = await reqLogin(params)
    token.value = res.data.token
    localStorage.setItem('token', res.data.token)
  }

  const logout = () => {
    token.value = ''
    localStorage.removeItem('token')
  }

  return {
    token,
    login,
    logout,
  }
})
