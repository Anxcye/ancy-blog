import { createRouter, createWebHistory } from 'vue-router'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      component: () => import('@/views/home/index.vue'),
    },
    {
      path: '/article',
      component: () => import('@/views/article/ArticleList.vue'),
      children: [
        {
          path: '/article/:id',
          component: () => import('@/views/article/ArticleDetail.vue'),
        },
      ],
    },
    {
      path: '/timeline',
      component: () => import('@/views/timeline/index.vue'),
    },
  ],
})

export default router
