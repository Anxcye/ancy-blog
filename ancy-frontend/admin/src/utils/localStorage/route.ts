import type { GetRoutersData } from '@/api/system/user/type'

export const localSetRoutes = (routes: GetRoutersData[]) => {
  localStorage.setItem('routes', JSON.stringify(routes))
}

export const localGetRoutes = (): GetRoutersData[] => {
  return JSON.parse(localStorage.getItem('routes') || '[]')
}

export const localRemoveRoutes = () => {
  localStorage.removeItem('routes')
}
