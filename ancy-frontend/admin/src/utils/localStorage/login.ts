import type { LoginRes } from '@/api/user/type'

export const setLoginInfo = (data: LoginRes) => {
  localStorage.setItem('userInfo', JSON.stringify(data))
}

export const getLoginInfo = (): LoginRes => {
  return JSON.parse(localStorage.getItem('userInfo') || '{}')
}
