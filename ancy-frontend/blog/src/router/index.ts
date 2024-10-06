import { createRouter, createWebHistory } from 'vue-router'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      component: () => import('@/views/home/HomeLayout.vue'),
      meta: { group: 'home' },
      children: [
        {
          path: '',
          name: 'home',
          component: () => import('@/views/home/HomeIndex.vue'),
          meta: { group: 'home' },
        },
        {
          path: 'home/:id',
          name: 'home-detail',
          component: () => import('@/views/home/HomeDetail.vue'),
          meta: { group: 'home' },
        },
      ],
    },
    {
      path: '/article',
      component: () => import('@/views/article/ArticleLayout.vue'),
      meta: { group: 'article' },
      children: [
        {
          path: '',
          name: 'article-list',
          component: () => import('@/views/article/ArticleList.vue'),
          meta: { group: 'article' },
        },
        {
          path: ':id',
          name: 'article-detail',
          component: () => import('@/views/article/ArticleDetail.vue'),
          meta: { group: 'article' },
        },
      ],
    },
    {
      path: '/category/:id',
      component: () => import('@/views/article/ArticleList.vue'),
      meta: { group: 'article' },
    },
    {
      path: '/timeline',
      component: () => import('@/views/timeline/index.vue'),
      meta: { group: 'timeline' },
    },
    {
      path: '/note',
      component: () => import('@/views/note/index.vue'),
      meta: { group: 'note' },
    },
    {
      path: '/link',
      component: () => import('@/views/link/index.vue'),
      meta: { group: 'link' },
    },
    {
      path: '/project',
      component: () => import('@/views/project/index.vue'),
      meta: { group: 'more' },
    },
    {
      path: '/project/:id',
      component: () => import('@/views/project/ProjectDetail.vue'),
      meta: { group: 'more' },
    },
    {
      path: '/read',
      component: () => import('@/views/read/index.vue'),
      meta: { group: 'more' },
    },
  ],
})

export default router
