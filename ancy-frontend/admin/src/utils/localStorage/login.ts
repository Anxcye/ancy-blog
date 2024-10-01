import type { LoginData } from '@/api/system/user/type'

export const setLoginInfo = (data: LoginData) => {
  localStorage.setItem('userInfo', JSON.stringify(data))
}

export const getLoginInfo = (): LoginData => {
  return JSON.parse(localStorage.getItem('userInfo') || '{}')
}

export const removeLoginInfo = () => {
  localStorage.removeItem('userInfo')
}
