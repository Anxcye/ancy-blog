<!--
File: SystemView.vue
Purpose: Reserve admin system page for integrations and runtime configuration.
Module: frontend-admin/views/system, presentation layer.
Related: integration providers APIs, translation job controls, security settings.
-->
<template>
  <section class="panel">
    <h1>{{ t('system.title') }}</h1>
    <p>{{ t('system.empty') }}</p>
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
