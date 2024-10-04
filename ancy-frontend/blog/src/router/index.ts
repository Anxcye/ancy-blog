import { createRouter, createWebHistory } from 'vue-router'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      component: () => import('../views/layout/index.vue'),
      children: [
        {
          path: '/article/:id',
          component: () => import('../views/ArticleDetail/index.vue'),
        },
      ],
    },
  ],
})

export default router
