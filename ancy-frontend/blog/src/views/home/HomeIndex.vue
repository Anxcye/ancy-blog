<template>
  <div
    class="absolute top-0 left-0 w-full h-1/4 bg-gradient-to-b from-[var(--primary-background)] to-transparent"
  ></div>
  <SayHi />

  <div v-for="i in 1000" :key="i">
    {{ i }}FrontPageodfahsdjhfaldskjfhalskdjfhalskdjfhalskdjhflaskdjhf
  </div>
</template>

<script setup lang="ts">
import SayHi from './components/SayHi.vue'
</script>

<style scoped></style>
