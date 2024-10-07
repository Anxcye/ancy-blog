<template>
  <div class="overflow-hidden flex flex-col min-h-screen w-full">
    <TopGradient />
    <PageHeader />
    <router-view class="mt-ac-header w-full flex-1"></router-view>
    <HomeFooter />
  </div>
</template>

<script setup lang="ts"></script>

<style scoped lang="scss"></style>
