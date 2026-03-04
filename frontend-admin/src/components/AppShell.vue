<!--
File: AppShell.vue
Purpose: Provide enterprise-grade admin chrome with sidebar navigation and responsive mobile drawer.
Module: frontend-admin/components, shell/layout layer.
Related: router routes, app store preferences, i18n labels, authentication state.
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
          <small>Control Panel</small>
        </div>
      </div>
      <NMenu :collapsed="collapsed" :collapsed-width="72" :collapsed-icon-size="20" :options="menuOptions" :value="activeKey" @update:value="handleMenuSelect" />
    </NLayoutSider>

    <NLayout>
      <NLayoutHeader bordered class="header">
        <div class="header-left">
          <NButton v-if="isMobile" quaternary circle @click="drawerOpen = true">
            <template #icon>
              <NIcon><MenuOutline /></NIcon>
            </template>
          </NButton>
          <div>
            <h1>{{ pageTitle }}</h1>
            <p>{{ t('layout.subtitle') }}</p>
          </div>
        </div>

        <NSpace align="center" :size="12">
          <NSelect
            size="small"
            style="width: 110px"
            :value="appStore.locale"
            :options="localeOptions"
            @update:value="handleLocaleChange"
          />
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
    </NLayout>

    <NDrawer v-model:show="drawerOpen" placement="left" :width="260">
      <NDrawerContent :title="t('layout.menu')" closable>
        <NMenu :options="menuOptions" :value="activeKey" @update:value="handleMenuSelect" />
      </NDrawerContent>
    </NDrawer>
  </NLayout>
</template>

<script setup lang="ts">
import { computed, h, onBeforeUnmount, onMounted, ref } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { useI18n } from 'vue-i18n';
import { BookOutline, ChatbubblesOutline, CogOutline, HomeOutline, LogOutOutline, MoonOutline, MenuOutline, NewspaperOutline, SunnyOutline } from '@vicons/ionicons5';
import type { MenuOption, SelectOption } from 'naive-ui';
import { NButton, NDrawer, NDrawerContent, NIcon, NLayout, NLayoutContent, NLayoutHeader, NLayoutSider, NMenu, NSelect, NSpace, NSwitch } from 'naive-ui';

import { useAppStore } from '@/stores/app';

const router = useRouter();
const route = useRoute();
const { t } = useI18n();
const appStore = useAppStore();

const collapsed = ref(false);
const drawerOpen = ref(false);
const isMobile = ref(false);

function renderIcon(icon: typeof HomeOutline) {
  return () => h(NIcon, null, { default: () => h(icon) });
}

const menuOptions = computed<MenuOption[]>(() => [
  { label: t('nav.workbench'), key: 'dashboard', icon: renderIcon(HomeOutline) },
  {
    label: t('nav.content'),
    key: 'content-group',
    icon: renderIcon(NewspaperOutline),
    children: [
      { label: t('articles.tabArticles'), key: 'articles', icon: renderIcon(BookOutline) },
      { label: t('articles.tabMoments'), key: 'moments', icon: renderIcon(ChatbubblesOutline) },
    ],
  },
  { label: t('nav.site'), key: 'site', icon: renderIcon(BookOutline) },
  { label: t('nav.interaction'), key: 'interaction', icon: renderIcon(ChatbubblesOutline) },
  { label: t('nav.system'), key: 'system', icon: renderIcon(CogOutline) },
]);

const routeNameToKey: Record<string, string> = {
  dashboard: 'dashboard',
  articles: 'articles',
  'article-new': 'articles',
  'article-edit': 'articles',
  moments: 'moments',
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

const localeOptions = computed<SelectOption[]>(() => [
  { label: '中文', value: 'zh-CN' },
  { label: 'English', value: 'en-US' },
]);

function syncViewport(): void {
  isMobile.value = window.innerWidth <= 992;
}

function handleMenuSelect(key: string): void {
  const targetMap: Record<string, string> = {
    dashboard: '/dashboard',
    articles: '/content/articles',
    moments: '/content/moments',
    site: '/site',
    interaction: '/interaction',
    system: '/system',
  };
  const target = targetMap[key];
  if (target) {
    drawerOpen.value = false;
    router.push(target);
  }
}

function handleThemeChange(enabled: boolean): void {
  appStore.setThemeMode(enabled ? 'dark' : 'light');
}

function handleLocaleChange(value: string): void {
  appStore.setLocale(value === 'en-US' ? 'en-US' : 'zh-CN');
}

function logout(): void {
  appStore.clearToken();
  drawerOpen.value = false;
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
  background:
    radial-gradient(circle at 14% 18%, rgba(38, 166, 154, 0.12), transparent 34%),
    radial-gradient(circle at 86% 82%, rgba(38, 166, 154, 0.08), transparent 42%),
    #f4f7f6;
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
  background: linear-gradient(135deg, #26a69a, #3ec7bb);
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

.brand-text small {
  color: #7a8791;
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
  font-size: 12px;
  color: #7b8791;
}

.content-area {
  padding: 18px;
}

.content-inner {
  width: min(1260px, 100%);
  margin: 0 auto;
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
