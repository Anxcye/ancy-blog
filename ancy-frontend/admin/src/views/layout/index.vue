<template>
  <div class="layout-container">
    <ResponsiveLayout :is-aside-visible="isAsideVisible">
      <template #aside>
        <div class="aside">
          <el-button @click="toggleAside" class="menu-button" type="primary">
            关闭
          </el-button>
          <AsideLogo style="overflow: hidden" />
          <el-scrollbar class="scrollbar">
            <div class="aside-menu">
              <AsideMenu />
            </div>
          </el-scrollbar>
        </div>
      </template>
      <template #main>
        <div class="main-container">
          <div class="header">
            <el-button @click="toggleAside" class="menu-button" type="primary">
              菜单
            </el-button>
            <div class="header-content">Header</div>
          </div>
          <div class="main">
            <MainContent />
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
import MainContent from './conponents/MainContent.vue'
const isAsideVisible = ref(false)

const toggleAside = () => {
  isAsideVisible.value = !isAsideVisible.value
}
</script>

<style scoped lang="scss">
.layout-container {
  .aside {
    margin: 10px;
    margin-right: 5px;
    height: calc(100% - 20px);

    border-radius: 20px;
    background-color: $ac-background;
    display: flex;
    flex-direction: column;
    padding: 10px;

    border: 1px solid $ac-border-color;

    .scrollbar {
      flex: 1;
      border-radius: 10px;
      overflow: hidden;
      font-size: 14px;
      font-weight: 600;
      background-color: white;
    }

    .menu-button {
      font-weight: 600;
    }
  }

  .main-container {
    display: flex;
    flex-direction: column;
    height: 100%;
    width: 100%;

    .header {
      margin: 10px;
      margin-left: 5px;
      margin-bottom: 5px;
      border-radius: 10px;
      padding: 10px;

      background-color: $ac-background;
      height: 4rem;
      display: flex;
      align-items: center;
      // justify-content: space-between;
      border: 1px solid $ac-border-color;

      .menu-button {
        font-weight: 600;
      }

      .header-content {
        flex: 1;
        text-align: center;
      }
    }

    .main {
      margin: 10px;
      margin-left: 5px;
      margin-top: 5px;
      border-radius: 10px;
      padding: 10px;

      background-color: $ac-background;
      border: 1px solid $ac-border-color;
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
