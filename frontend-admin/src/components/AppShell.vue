<!--
File: AppShell.vue
Purpose: Provide common admin chrome including header and content container.
Module: frontend-admin/components, shell/layout layer.
Related: router routes, app store, localized menu labels.
-->
<template>
  <div class="shell">
    <header class="header">
      <div class="brand">Ancy Admin</div>
      <nav class="nav">
        <RouterLink to="/dashboard">{{ t('nav.workbench') }}</RouterLink>
        <RouterLink to="/content/articles">{{ t('nav.content') }}</RouterLink>
        <RouterLink to="/site">{{ t('nav.site') }}</RouterLink>
        <RouterLink to="/interaction">{{ t('nav.interaction') }}</RouterLink>
        <RouterLink to="/system">{{ t('nav.system') }}</RouterLink>
      </nav>
    </header>
    <main class="content">
      <slot />
    </main>
    <nav class="mobile-tabbar">
      <RouterLink to="/dashboard">{{ t('nav.workbench') }}</RouterLink>
      <RouterLink to="/content/articles">{{ t('nav.content') }}</RouterLink>
      <RouterLink to="/site">{{ t('nav.site') }}</RouterLink>
      <RouterLink to="/interaction">{{ t('nav.interaction') }}</RouterLink>
      <RouterLink to="/system">{{ t('nav.system') }}</RouterLink>
    </nav>
  </div>
</template>

<script setup lang="ts">
import { useI18n } from 'vue-i18n';

const { t } = useI18n();
</script>

<style scoped>
.shell {
  min-height: 100vh;
}

.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 24px;
  border-bottom: 1px solid var(--border);
  background: var(--surface);
  position: sticky;
  top: 0;
  z-index: 10;
}

.brand {
  font-weight: 700;
  letter-spacing: 0.02em;
}

.nav {
  display: flex;
  gap: 16px;
}

.nav a {
  text-decoration: none;
  color: var(--muted);
}

.nav a.router-link-active {
  color: var(--accent);
}

.content {
  width: min(1100px, 100%);
  margin: 24px auto;
  padding: 0 16px 88px;
}

.mobile-tabbar {
  display: none;
}

@media (max-width: 900px) {
  .header {
    padding: 10px 14px;
  }

  .brand {
    font-size: 14px;
  }

  .nav {
    display: none;
  }

  .content {
    margin: 14px auto;
    padding: 0 12px 88px;
  }

  .mobile-tabbar {
    display: grid;
    grid-template-columns: repeat(5, minmax(0, 1fr));
    gap: 4px;
    position: fixed;
    bottom: 0;
    left: 0;
    right: 0;
    z-index: 20;
    padding: 8px 6px calc(8px + env(safe-area-inset-bottom));
    border-top: 1px solid var(--border);
    background: var(--surface);
  }

  .mobile-tabbar a {
    text-align: center;
    text-decoration: none;
    font-size: 12px;
    line-height: 1.4;
    padding: 6px 4px;
    border-radius: 8px;
    color: var(--muted);
  }

  .mobile-tabbar a.router-link-active {
    color: var(--accent-hover);
    background: var(--accent-soft);
    font-weight: 600;
  }
}
</style>
