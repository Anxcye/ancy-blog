import { createRouter, createWebHistory } from 'vue-router'
import { routers } from './routers'
import { useRouteStore } from '@/stores/modules/route'
import { useUserStore } from '@/stores/modules/user'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: routers,
  scrollBehavior: () => ({ left: 0, top: 0 }),
})

const whiteList = ['/login']

router.beforeEach(async (to, from, next) => {
  const routerStore = useRouteStore()
  const userStore = useUserStore()

  if (whiteList.includes(to.path)) {
    next()
  } else {
    if (!userStore.getToken()) {
      next({ path: '/login' })
    }

    if (routerStore.routesLoaded) {
      next()
    } else {
      await routerStore.setRoutes()
      next({ ...to, replace: true })
    }
  }
})

export default router
