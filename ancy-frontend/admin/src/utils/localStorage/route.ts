import type { GetRoutersData } from '@/api/user/type'

export const localSetRoutes = (routes: GetRoutersData[]) => {
  localStorage.setItem('routes', JSON.stringify(routes))
}

export const localGetRoutes = (): GetRoutersData[] => {
  return JSON.parse(localStorage.getItem('routes') || '[]')
}
