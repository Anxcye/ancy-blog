<template>
  <div class="aside-logo" v-if="!logo.hideTitle">
    <img :src="logo.image" alt="logo" />
    <h1>{{ logo.title }}</h1>
  </div>
</template>

<script setup lang="ts">
import { logo } from '@/settings'
</script>

<style scoped lang="scss">
.aside-logo {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 100%;
  height: 64px;
  transition: width 0.3s;

  img {
    width: 32px;
    height: 32px;
  }

  h1 {
    margin-left: 10px;
    font-size: 20px;
    font-weight: 600;
  }
}
</style>
