<template>
  <div class="overflow-hidden">
    <PageHeader />
    <router-view class="router-view"></router-view>
    <ConfigProvider />
  </div>
</template>

<script setup lang="ts">
import PageHeader from '@/components/PageHeader.vue'
import ConfigProvider from '@/provider/ConfigProvider.vue'
</script>

<style scoped lang="scss">
.page-header {
  top: 0;
  left: 0;
  right: 0;
  z-index: 1000;
}
.router-view {
  margin-top: $ac-header-height;
}
</style>
