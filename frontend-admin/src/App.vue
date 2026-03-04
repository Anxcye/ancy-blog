<!--
File: App.vue
Purpose: Define root providers and top-level shell for the admin workspace.
Module: frontend-admin/app, application composition layer.
Related: AppShell, i18n runtime locale, app store theme settings, Naive UI providers.
-->
<template>
  <NConfigProvider :theme="naiveTheme" :theme-overrides="themeOverrides" :locale="naiveLocale" :date-locale="naiveDateLocale">
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
import { computed, watch } from 'vue';
import { dateEnUS, dateZhCN, darkTheme, enUS, NConfigProvider, NDialogProvider, NGlobalStyle, NLoadingBarProvider, NMessageProvider, NNotificationProvider, zhCN } from 'naive-ui';

import AppShell from '@/components/AppShell.vue';
import i18n from '@/i18n';
import { useAppStore } from '@/stores/app';
import { darkThemeOverrides, lightThemeOverrides } from '@/styles/naive-theme';

const appStore = useAppStore();

watch(
  () => appStore.locale,
  (value) => {
    i18n.global.locale.value = value;
  },
  { immediate: true },
);

const naiveTheme = computed(() => (appStore.themeMode === 'dark' ? darkTheme : null));
const themeOverrides = computed(() => (appStore.themeMode === 'dark' ? darkThemeOverrides : lightThemeOverrides));
const naiveLocale = computed(() => (appStore.locale === 'en-US' ? enUS : zhCN));
const naiveDateLocale = computed(() => (appStore.locale === 'en-US' ? dateEnUS : dateZhCN));
</script>
