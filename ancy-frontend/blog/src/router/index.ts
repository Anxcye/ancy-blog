import { useBrowserStore } from '@/stores/browser'
import { useColorStore } from '@/stores/color'
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
          name: '首页',
          component: () => import('@/views/home/HomeIndex.vue'),
          meta: { group: 'home' },
        },
        {
          path: 'home/:id',
          name: '首页详情',
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
          name: '文章列表',
          component: () => import('@/views/article/ArticleList.vue'),
          meta: { group: 'article' },
        },
        {
          path: ':id',
          name: '文章详情',
          component: () => import('@/views/article/ArticleDetail.vue'),
          meta: { group: 'article' },
        },
      ],
    },
    {
      path: '/category/:id',
      name: '分类',
      component: () => import('@/views/article/ArticleList.vue'),
      meta: { group: 'article' },
    },
    {
      path: '/timeline',
      name: '时间线',
      component: () => import('@/views/timeline/index.vue'),
      meta: { group: 'timeline' },
    },
    {
      path: '/note',
      name: '日志',
      component: () => import('@/views/note/index.vue'),
      meta: { group: 'note' },
    },
    {
      path: '/link',
      name: '友链',
      component: () => import('@/views/link/index.vue'),
      meta: { group: 'link' },
    },
    {
      path: '/project',
      name: '项目',
      component: () => import('@/views/project/index.vue'),
      meta: { group: 'more' },
    },
    {
      path: '/project/:id',
      name: '项目详情',
      component: () => import('@/views/project/ProjectDetail.vue'),
      meta: { group: 'more' },
    },
    {
      path: '/read',
      name: '阅读',
      component: () => import('@/views/read/index.vue'),
      meta: { group: 'more' },
    },
  ],
})

router.beforeEach((to, from, next) => {
  const colorStore = useColorStore()
  const browserStore = useBrowserStore()

  colorStore.changeColor()

  browserStore.setTitle(to.name as string)

  next()
})

export default router
