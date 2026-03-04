<!--
File: InteractionView.vue
Purpose: Reserve admin interaction page for comments and link review workflows.
Module: frontend-admin/views/interaction, presentation layer.
Related: comment moderation APIs and friend-link review APIs.
-->
<template>
  <section class="panel">
    <h1>{{ t('interaction.title') }}</h1>
    <p>{{ t('interaction.empty') }}</p>
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
