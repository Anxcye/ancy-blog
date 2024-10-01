import { getRouters } from '@/api/system/user'
import type { GetRoutersData } from '@/api/system/user/type'
import router from '@/router'
import { localGetRoutes, localRemoveRoutes, localSetRoutes } from '@/utils/localStorage/route'
import { defineStore } from 'pinia'
import { computed, ref } from 'vue'

const modules = import.meta.glob('../../views/layout/**/*.vue')

export const useRouteStore = defineStore('route', () => {
  const routes = ref<GetRoutersData[]>(localGetRoutes())
  const routesLoaded = ref(false)

  const routeArray = computed(() => {
    const result: GetRoutersData[] = []
    const queue: GetRoutersData[] = [...routes.value]
    while (queue.length) {
      const node = queue.shift()
      if (!node) continue
      result.push(node)
      if (node.children) {
        queue.push(...node.children)
      }
    }
    return result
  })

  const removeRoutes = () => {
    localRemoveRoutes()
    routes.value = []
    router.replace('/login')
  }

  const initRoutes = async () => {
    const res = await getRouters()
    localSetRoutes(res.data.menus)
    routes.value = res.data.menus
    addRouter()
  }

  const setRoutes = async () => {
    if (!routeArray.value.length) {
      await initRoutes()
    }
    addRouter()
  }

  const addRouter = () => {
    for (const item of routeArray.value) {
      if (item.menuType !== 'C') continue
      router.addRoute('home', {
        path: item.path,
        name: item.path,
        component: modules[`../../views/layout/${item.component}.vue`],
        meta: {
          title: item.menuName,
          icon: item.icon,
        },
        children: item.children?.map((child) => ({
          path: child.path,
          name: child.path,
          component: modules[`../../views/layout/${child.component}.vue`],
          meta: {
            title: child.menuName,
            icon: child.icon,
          },
        })),
      })
    }
    router.addRoute({
      path: '/:pathMatch(.*)*',
      name: 'any',
      redirect: '/404',
    })
    routesLoaded.value = true
  }

  const getRoutes = (key: string): GetRoutersData | undefined => {
    return routeArray.value.find((item) => item.id === parseInt(key))
  }

  const getIdByPath = (path: string): number => {
    return routeArray.value.find((item) => item.path === path.slice(1))?.id || -1
  }

  return {
    routes,
    routesLoaded,
    initRoutes,
    removeRoutes,
    setRoutes,
    getRoutes,
    getIdByPath,
  }
})
