<!--
File: App.vue
Purpose: Define root shell and top-level layout for admin workspace.
Module: frontend-admin/app, presentation shell layer.
Related: AppShell component and router-view content pages.
-->
<template>
  <AppShell>
    <RouterView />
  </AppShell>
</template>

<script setup lang="ts">
import AppShell from '@/components/AppShell.vue';
</script>
