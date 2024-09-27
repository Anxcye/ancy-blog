<template>
  <div class="layout-container">
    <ResponsiveLayout :is-aside-visible="isAsideVisible">
      <template #aside>
        <div class="aside">
          <el-button @click="toggleAside" class="menu-button">菜单</el-button>
          <AsideLogo />
          <el-scrollbar class="scrollbar">
            <AsideMenu />
          </el-scrollbar>
        </div>
      </template>
      <template #main>
        <div class="main-container">
          <div class="header">
            Header
            <el-button @click="toggleAside" class="menu-button">菜单</el-button>
          </div>
          <div class="main">
            <router-view />
          </div>
        </div>
      </template>
    </ResponsiveLayout>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import AsideLogo from './conponents/AsideLogo.vue'
import AsideMenu from './conponents/AsideMenu.vue'
import { useRouteStore } from '@/stores/modules/route'
const isAsideVisible = ref(false)

const toggleAside = () => {
  isAsideVisible.value = !isAsideVisible.value
}

const routeStore = useRouteStore()

routeStore.setRoutes()
</script>

<style scoped>
.layout-container {
  height: 100%;
  width: 100%;

  .aside {
    background-color: #f0f0f0;
    height: 100%;
    width: 100%;
    display: flex;
    flex-direction: column;

    .scrollbar {
      flex: 1;
    }
  }

  .main-container {
    display: flex;
    flex-direction: column;
    height: 100%;
    width: 100%;

    .header {
      background-color: #d3dce6;
      /* height: 100%; */
      /* width: 100%; */
      .menu-button {
        float: left;
      }
    }

    .main {
      background-color: #e9eef3;
      /* height: 100%; */
      /* width: 100%; */
      flex: 1;
      overflow: auto;
    }
  }
}

@media (min-width: 768px) {
  .menu-button {
    display: none;
  }
}
</style>
