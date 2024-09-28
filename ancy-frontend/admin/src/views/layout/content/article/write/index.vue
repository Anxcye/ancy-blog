<template>
  <div>
    Write
    <ElInput v-model="title" placeholder="请输入标题"></ElInput>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'

const title = ref('')
</script>

<style scoped lang="scss"></style>
