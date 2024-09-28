import { createRouter, createWebHistory } from 'vue-router'
import { routers } from './routers'
import { useRouteStore } from '@/stores/modules/route'
import { useUserStore } from '@/stores/modules/user'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: routers,
  scrollBehavior: () => ({ left: 0, top: 0 }),
})
export default router
