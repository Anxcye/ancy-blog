export const routers = [
  {
    path: '/login',
    name: 'login',
    component: () => import('@/views/login/index.vue'),
    meta: {
      title: '登录',
    },
  },
  {
    path: '/',
    name: 'home',
    component: () => import('@/views/layout/index.vue'),
    meta: {
      title: '首页',
    },
  },
  {
    path: '/404',
    name: '404',
    component: () => import('@/views/404/index.vue'),
    meta: {
      title: '404',
    },
  },
  // {
  //   path: '/:pathMatch(.*)*',
  //   name: 'any',
  //   redirect: '/404',
  // },
]
