<!--
File: DashboardView.vue
Purpose: Provide initial admin overview page and pending work entry point.
Module: frontend-admin/views/dashboard, presentation layer.
Related: translation jobs, content management modules, metrics widgets.
-->
<template>
  <section class="panel">
    <h1>{{ t('dashboard.title') }}</h1>
    <p>{{ t('dashboard.subtitle') }}</p>
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
