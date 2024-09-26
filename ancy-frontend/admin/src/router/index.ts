import { createRouter, createWebHistory } from 'vue-router'
import { routers } from './routers'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: routers,
  scrollBehavior: () => ({ left: 0, top: 0 }),
})

export default router
