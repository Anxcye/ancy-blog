import { useRouteStore } from '@/stores/modules/route'
import { useUserStore } from '@/stores/modules/user'
import NProgress from 'nprogress'
import 'nprogress/nprogress.css'

import router from '@/router'
import { useTabStore } from './stores/modules/tab'
import { useBrowserStore } from './stores/modules/browser'
import { ref } from 'vue'

const whiteList = ['/login']

router.beforeEach(async (to, from, next) => {
  const routerStore = useRouteStore()
  const userStore = useUserStore()
  const tabStore = useTabStore()
  const browserStore = useBrowserStore()
  const setFavicon = ref(false)

  if (!setFavicon.value) {
    browserStore.setIcon(to.meta.icon as string)
    setFavicon.value = true
  }

  const go = (param?: any) => {
    NProgress.start()
    if (to.path !== '/') {
      tabStore.addHistoryTab(to.path)
    }
    tabStore.currentTab = to.path
    browserStore.setTitle(to.meta.title as string)
    next(param)
  }

  if (whiteList.includes(to.path)) {
    go()
  } else {
    if (!userStore.getToken()) {
      go({ path: '/login' })
    }

    if (routerStore.routesLoaded) {
      go()
    } else {
      NProgress.start()
      await routerStore.setRoutes()
      go({ ...to, replace: true })
    }
  }
})

router.afterEach(() => {
  NProgress.done()
})
