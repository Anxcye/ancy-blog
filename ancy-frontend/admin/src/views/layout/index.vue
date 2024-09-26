<template>
  <div class="layout-container">
    <ResponsiveLayout :is-aside-visible="isAsideVisible">
      <template #aside>
        <div class="aside">Aside</div>
      </template>
      <template #header>
        <div class="header">
          Header
          <el-button @click="toggleAside" class="menu-button">菜单</el-button>
        </div>
      </template>
      <template #main>
        <div class="main">Main</div>
      </template>
    </ResponsiveLayout>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'

const isAsideVisible = ref(false)

const toggleAside = () => {
  isAsideVisible.value = !isAsideVisible.value
}
</script>

<style scoped>
.layout-container {
  height: 100%;
  width: 100%;
}

.aside {
  background-color: #f0f0f0;
  height: 100%;
  width: 100%;
}

.header {
  background-color: #d3dce6;
  height: 100%;
  width: 100%;
  .menu-button {
    float: right;
  }
}

@media (min-width: 768px) {
  .menu-button {
    display: none;
  }
}

.main {
  background-color: #e9eef3;
  height: 100%;
  width: 100%;
}
</style>
