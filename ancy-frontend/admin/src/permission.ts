import { useRouteStore } from '@/stores/modules/route'
import { useUserStore } from '@/stores/modules/user'
import NProgress from 'nprogress'
import 'nprogress/nprogress.css'

import router from '@/router'

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
      NProgress.start()
      next()
    } else {
      NProgress.start()
      await routerStore.setRoutes()
      next({ ...to, replace: true })
    }
  }
})

router.afterEach(() => {
  NProgress.done()
})
