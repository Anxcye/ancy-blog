<script setup lang="ts">
import request from '@/utils/request'
import { onMounted } from 'vue'

onMounted(() => {
  request.get('/articles/page').then((res) => {
    console.log(res.data)
  })
})
</script>

<template>
  <h1>Hello World</h1>
  <el-button type="primary">Primary Button</el-button>
</template>

<style scoped lang="scss"></style>
