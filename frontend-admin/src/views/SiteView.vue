<!--
File: SiteView.vue
Purpose: Reserve admin site-orchestration page for nav, slots, and footer operations.
Module: frontend-admin/views/site, presentation layer.
Related: site settings APIs, nav items, footer items, slot mapping modules.
-->
<template>
  <section class="panel">
    <h1>{{ t('site.title') }}</h1>
    <p>{{ t('site.empty') }}</p>
  </section>
</template>

<script setup lang="ts">
import { useI18n } from 'vue-i18n';

const { t } = useI18n();
</script>

<style scoped>
.panel {
  padding: 24px;
  border: 1px solid var(--border);
  border-radius: 12px;
  background: var(--surface);
}

h1 {
  margin: 0 0 10px;
}

p {
  margin: 0;
  color: var(--muted);
}
</style>
