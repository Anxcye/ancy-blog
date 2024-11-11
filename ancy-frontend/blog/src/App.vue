<template>
  <div class="overflow-hidden flex flex-col min-h-screen w-full">
    <TopGradient />
    <PageHeader />
    <router-view class="mt-ac-header w-full flex-1"></router-view>
    <HomeFooter />
  </div>
</template>

<script setup lang="ts">
import { useHead } from '@vueuse/head'
import { useBaseInfoStore } from './stores/baseInfo'
import { computed } from 'vue'
import getMeta from './utils/meta'

const baseInfoStore = useBaseInfoStore()

useHead({
  meta: getMeta(
    computed(() => baseInfoStore.baseInfo?.philosophy || 'Anxcye 空间'),
    computed(
      () =>
        baseInfoStore.baseInfo?.badge?.map((item) => item.title) ?? [
          'AnCy',
          'Anxcye',
          'Blog',
          '技术',
          '生活',
          '博客',
          'Github',
          '开源',
          '前端',
          '后端',
          '全栈',
        ],
    ),
  ),
})
</script>

<style scoped lang="scss"></style>
