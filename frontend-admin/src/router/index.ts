// File: router/index.ts
// Purpose: Configure admin SPA routes and lightweight auth guard.
// Module: frontend-admin/router, navigation layer.
// Related: auth store, login view, dashboard and content views.
import { createRouter, createWebHistory } from 'vue-router';

import { useAppStore } from '@/stores/app';
import DashboardView from '@/views/DashboardView.vue';
import LoginView from '@/views/LoginView.vue';
import ArticlesView from '@/views/ArticlesView.vue';
import SiteView from '@/views/SiteView.vue';
import InteractionView from '@/views/InteractionView.vue';
import SystemView from '@/views/SystemView.vue';
import ArticleEditorView from '@/views/ArticleEditorView.vue';
import MomentsView from '@/views/MomentsView.vue';

const router = createRouter({
  history: createWebHistory(),
  routes: [
    { path: '/login', name: 'login', component: LoginView },
    { path: '/', redirect: '/dashboard' },
    { path: '/dashboard', name: 'dashboard', component: DashboardView, meta: { requiresAuth: true } },
    { path: '/content/articles', name: 'articles', component: ArticlesView, meta: { requiresAuth: true } },
    { path: '/content/moments', name: 'moments', component: MomentsView, meta: { requiresAuth: true } },
    { path: '/content/articles/new', name: 'article-new', component: ArticleEditorView, meta: { requiresAuth: true } },
    { path: '/content/articles/:id/edit', name: 'article-edit', component: ArticleEditorView, meta: { requiresAuth: true } },
    { path: '/site', name: 'site', component: SiteView, meta: { requiresAuth: true } },
    { path: '/interaction', name: 'interaction', component: InteractionView, meta: { requiresAuth: true } },
    { path: '/system', name: 'system', component: SystemView, meta: { requiresAuth: true } },
  ],
});

router.beforeEach((to) => {
  const appStore = useAppStore();

  if (to.meta.requiresAuth && !appStore.isAuthenticated) {
    return { name: 'login' };
  }

  if (to.name === 'login' && appStore.isAuthenticated) {
    return { name: 'dashboard' };
  }

  return true;
});

export default router;
