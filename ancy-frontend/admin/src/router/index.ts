import { createRouter, createWebHistory } from 'vue-router'
import { routers } from './routers'
import { useRouteStore } from '@/stores/modules/route'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: routers,
  scrollBehavior: () => ({ left: 0, top: 0 }),
})

router.beforeEach(async (to, from, next) => {
  const routerStore = useRouteStore()
  if (routerStore.routesLoaded) {
    next()
  } else {
    await routerStore.setRoutes()
    next({ ...to, replace: true })
  }
})
export default router
