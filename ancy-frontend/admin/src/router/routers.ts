export const routers = [
  {
    path: '/login',
    name: 'login',
    component: () => import('@/views/login/index.vue'),
  },
  {
    path: '/',
    name: 'home',
    component: () => import('@/views/layout/index.vue'),

    children: [
      {
        path: '/ancy',
        name: 'ancy',
        component: () => import('@/views/layout/conponents/AncyLogo.vue'),
        meta: {
          title: 'Ancy',
        },
      },
    ],
  },
  {
    path: '/404',
    name: '404',
    component: () => import('@/views/404/index.vue'),
  },
  // {
  //   path: '/:pathMatch(.*)*',
  //   name: 'any',
  //   redirect: '/404',
  // },
]
