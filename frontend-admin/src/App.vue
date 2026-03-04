<!--
File: App.vue
Purpose: Define root providers and top-level shell for the admin workspace.
Module: frontend-admin/app, application composition layer.
Related: AppShell, app store theme settings, Naive UI providers.
-->
<template>
  <NConfigProvider :theme="naiveTheme" :theme-overrides="themeOverrides" :locale="zhCN" :date-locale="dateZhCN">
    <NLoadingBarProvider>
      <NDialogProvider>
        <NNotificationProvider>
          <NMessageProvider>
            <NGlobalStyle />
            <AppShell>
              <RouterView />
            </AppShell>
          </NMessageProvider>
        </NNotificationProvider>
      </NDialogProvider>
    </NLoadingBarProvider>
  </NConfigProvider>
</template>

<script setup lang="ts">
import { computed } from 'vue';
import { dateZhCN, darkTheme, NConfigProvider, NDialogProvider, NGlobalStyle, NLoadingBarProvider, NMessageProvider, NNotificationProvider, zhCN } from 'naive-ui';

import AppShell from '@/components/AppShell.vue';
import { useAppStore } from '@/stores/app';
import { darkThemeOverrides, lightThemeOverrides } from '@/styles/naive-theme';

const appStore = useAppStore();

const naiveTheme = computed(() => (appStore.themeMode === 'dark' ? darkTheme : null));
const themeOverrides = computed(() => (appStore.themeMode === 'dark' ? darkThemeOverrides : lightThemeOverrides));
</script>
