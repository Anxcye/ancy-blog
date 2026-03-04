<!--
File: AppShell.vue
Purpose: Provide enterprise-grade admin chrome with desktop sider and mobile bottom navigation.
Module: frontend-admin/components, shell/layout layer.
Related: router routes, app store theme preference, authentication state.
-->
<template>
  <div v-if="isLoginPage" class="login-layout">
    <slot />
  </div>

  <NLayout v-else has-sider class="shell-layout">
    <NLayoutSider
      v-if="!isMobile"
      bordered
      collapse-mode="width"
      :collapsed-width="72"
      :width="250"
      :collapsed="collapsed"
      show-trigger
      @collapse="collapsed = true"
      @expand="collapsed = false"
    >
      <div class="brand" :class="{ compact: collapsed }">
        <span class="brand-mark">A</span>
        <div v-if="!collapsed" class="brand-text">
          <strong>Ancy Admin</strong>
          <small>后台管理系统</small>
        </div>
      </div>
      <NMenu :collapsed="collapsed" :collapsed-width="72" :collapsed-icon-size="20" :options="menuOptions" :value="activeKey" @update:value="handleMenuSelect" />
    </NLayoutSider>

    <NLayout>
      <NLayoutHeader bordered class="header">
        <div class="header-left">
          <div>
            <h1>{{ pageTitle }}</h1>
            <p>{{ t('layout.subtitle') }}</p>
          </div>
        </div>

        <NSpace align="center" :size="12">
          <NSwitch :value="appStore.themeMode === 'dark'" @update:value="handleThemeChange">
            <template #checked-icon>
              <NIcon><MoonOutline /></NIcon>
            </template>
            <template #unchecked-icon>
              <NIcon><SunnyOutline /></NIcon>
            </template>
          </NSwitch>
          <NButton tertiary size="small" @click="logout">{{ t('layout.logout') }}</NButton>
        </NSpace>
      </NLayoutHeader>

      <NLayoutContent class="content-area">
        <div class="content-inner">
          <slot />
        </div>
      </NLayoutContent>

      <nav v-if="isMobile" class="mobile-tabbar">
        <NButton
          v-for="item in mobileTabs"
          :key="item.key"
          text
          size="small"
          class="tab-item"
          :class="{ active: activeKey === item.key }"
          @click="handleMenuSelect(item.key)"
        >
          <span>{{ item.label }}</span>
        </NButton>
      </nav>
    </NLayout>
  </NLayout>
</template>

<script setup lang="ts">
import { computed, h, onBeforeUnmount, onMounted, ref } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { useI18n } from 'vue-i18n';
import { BookOutline, ChatbubblesOutline, CogOutline, HomeOutline, MoonOutline, NewspaperOutline, SunnyOutline } from '@vicons/ionicons5';
import type { MenuOption } from 'naive-ui';
import { NButton, NIcon, NLayout, NLayoutContent, NLayoutHeader, NLayoutSider, NMenu, NSpace, NSwitch } from 'naive-ui';

import { useAppStore } from '@/stores/app';

const router = useRouter();
const route = useRoute();
const { t } = useI18n();
const appStore = useAppStore();

const collapsed = ref(false);
const isMobile = ref(false);

function renderIcon(icon: typeof HomeOutline) {
  return () => h(NIcon, null, { default: () => h(icon) });
}

const menuOptions = computed<MenuOption[]>(() => [
  { label: t('nav.workbench'), key: 'dashboard', icon: renderIcon(HomeOutline) },
  { label: t('nav.content'), key: 'content', icon: renderIcon(NewspaperOutline) },
  { label: t('nav.site'), key: 'site', icon: renderIcon(BookOutline) },
  { label: t('nav.interaction'), key: 'interaction', icon: renderIcon(ChatbubblesOutline) },
  { label: t('nav.system'), key: 'system', icon: renderIcon(CogOutline) },
]);

const routeNameToKey: Record<string, string> = {
  dashboard: 'dashboard',
  articles: 'content',
  'article-new': 'content',
  'article-edit': 'content',
  moments: 'content',
  site: 'site',
  interaction: 'interaction',
  system: 'system',
};

const routeTitleMap = computed<Record<string, string>>(() => ({
  dashboard: t('dashboard.title'),
  articles: t('articles.title'),
  'article-new': t('editor.createTitle'),
  'article-edit': t('editor.editTitle'),
  moments: t('moments.title'),
  site: t('site.title'),
  interaction: t('interaction.title'),
  system: t('system.title'),
}));

const activeKey = computed(() => routeNameToKey[String(route.name || '')] || 'dashboard');
const pageTitle = computed(() => routeTitleMap.value[String(route.name || '')] || 'Ancy Admin');
const isLoginPage = computed(() => route.name === 'login');

const mobileTabs = computed(() => [
  { key: 'dashboard', label: t('nav.workbench') },
  { key: 'content', label: t('nav.content') },
  { key: 'site', label: t('nav.site') },
  { key: 'interaction', label: t('nav.interaction') },
  { key: 'system', label: t('nav.system') },
]);

function syncViewport(): void {
  isMobile.value = window.innerWidth <= 992;
}

function handleMenuSelect(key: string): void {
  const targetMap: Record<string, string> = {
    dashboard: '/dashboard',
    content: '/content/articles',
    articles: '/content/articles',
    moments: '/content/moments',
    site: '/site',
    interaction: '/interaction',
    system: '/system',
  };
  const target = targetMap[key];
  if (target) {
    router.push(target);
  }
}

function handleThemeChange(enabled: boolean): void {
  appStore.setThemeMode(enabled ? 'dark' : 'light');
}

function logout(): void {
  appStore.clearToken();
  router.push({ name: 'login' });
}

onMounted(() => {
  syncViewport();
  window.addEventListener('resize', syncViewport);
});

onBeforeUnmount(() => {
  window.removeEventListener('resize', syncViewport);
});
</script>

<style scoped>
.shell-layout {
  min-height: 100vh;
}

.login-layout {
  min-height: 100vh;
  display: grid;
  place-items: center;
  padding: 20px;
  background-color: var(--n-body-color);
}

.brand {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 18px 14px 10px;
}

.brand.compact {
  justify-content: center;
}

.brand-mark {
  width: 36px;
  height: 36px;
  border-radius: 12px;
  background: linear-gradient(135deg, var(--n-primary-color), var(--n-primary-color-hover));
  color: #fff;
  display: grid;
  place-items: center;
  font-size: 16px;
  font-weight: 700;
}

.brand-text {
  display: grid;
  line-height: 1.2;
}

.brand-text strong {
  font-size: 14px;
}

.brand-text small,
.header-left p {
  color: var(--n-text-color-3);
  font-size: 12px;
}

.header {
  height: 70px;
  padding: 0 18px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 12px;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 10px;
}

.header-left h1 {
  margin: 0;
  font-size: 18px;
  font-weight: 650;
  line-height: 1.2;
}

.header-left p {
  margin: 3px 0 0;
}

.content-area {
  padding: 18px;
  overflow-x: hidden;
}

.content-inner {
  width: min(1260px, 100%);
  margin: 0 auto;
  padding-bottom: 84px;
}

.mobile-tabbar {
  position: fixed;
  left: 0;
  right: 0;
  bottom: 0;
  z-index: 30;
  display: grid;
  grid-template-columns: repeat(5, minmax(0, 1fr));
  gap: 0;
  padding: 6px 6px calc(6px + env(safe-area-inset-bottom));
  background: var(--n-card-color);
  border-top: 1px solid var(--n-border-color);
  box-shadow: 0 -8px 24px color-mix(in srgb, var(--n-text-color) 12%, transparent);
  isolation: isolate;
  backdrop-filter: none;
  -webkit-backdrop-filter: none;
}

.mobile-tabbar::before {
  content: '';
  position: absolute;
  inset: 0;
  z-index: -1;
  background: var(--n-card-color);
  opacity: 1;
}

.tab-item {
  display: grid;
  gap: 0;
  place-items: center;
  color: var(--n-text-color-2);
  padding: 8px 2px;
  border-radius: 8px;
  font-size: 14px;
  font-weight: 650;
  height: 52px;
}

.tab-item.active {
  color: var(--n-primary-color);
  background: color-mix(in srgb, var(--n-primary-color) 20%, transparent);
  font-weight: 750;
}

@media (max-width: 992px) {
  .header {
    height: 64px;
    padding: 0 12px;
  }

  .header-left h1 {
    font-size: 16px;
  }

  .header-left p {
    display: none;
  }

  .content-area {
    padding: 10px;
  }
}
</style>
