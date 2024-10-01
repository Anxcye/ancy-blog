<template>
  <div class="layout-container">
    <ResponsiveLayout :is-aside-visible="isAsideVisible">
      <template #aside>
        <div class="aside">
          <el-button @click="toggleAside" class="menu-button" type="primary">关闭</el-button>
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
            <div class="title">
              <el-button @click="toggleAside" class="menu-button" type="primary">菜单</el-button>

              <div class="header-content">Header</div>
              <!-- 右侧用户信息  -->
              <UserInfo />
            </div>
            <div class="controller">
              <div class="tab">
                <el-tag
                  v-for="tag in tabStore.historyTabs"
                  :key="tag.path"
                  closable
                  :type="tag.path === tabStore.currentTab ? 'primary' : 'info'"
                  @click="handleTabClick(tag)"
                  @close="removeTab(tag)"
                >
                  {{ tag.meta.title }}
                </el-tag>
              </div>
              <div class="function">
                <FullScreen @click="toggleFullscreen" />
                <Refresh @click="refresh" />
              </div>
            </div>
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
import AsideLogo from './conponents/AsideLogo.vue'
import AsideMenu from './conponents/AsideMenu.vue'
import UserInfo from './conponents/UserInfo.vue'
import MainContent from './conponents/MainContent.vue'
import { Refresh, FullScreen } from '@element-plus/icons-vue'
import { useLayoutStore } from '@/stores/modules/layout'
import { useTabStore } from '@/stores/modules/tab'

const layoutStore = useLayoutStore()
const tabStore = useTabStore()

const isAsideVisible = ref(false)

const toggleAside = () => {
  isAsideVisible.value = !isAsideVisible.value
}

const refresh = () => {
  layoutStore.setRefresh()
}

const toggleFullscreen = () => {
  const full = document.fullscreenElement
  if (full) {
    document.exitFullscreen()
  } else {
    document.documentElement.requestFullscreen()
  }
}

import { ref } from 'vue'
import router from '@/router'
import type { RouteLocationNormalized } from 'vue-router'

const removeTab = (targetName: RouteLocationNormalized) => {
  tabStore.removeHistoryTab(targetName.path)
}

const handleTabClick = (tab: RouteLocationNormalized) => {
  router.push(tab.path)
  tabStore.currentTab = tab.path
}
</script>

<style scoped lang="scss">
.layout-container {
  .aside {
    display: flex;
    flex-direction: column;
    height: calc(100% - 20px);
    padding: 10px;
    margin: 10px;
    margin-right: 5px;
    background-color: $ac-background;
    border: 1px solid $ac-border-color;
    border-radius: 20px;

    .scrollbar {
      flex: 1;
      overflow: hidden;
      font-size: 14px;
      font-weight: 600;
      background-color: white;
      border-radius: 10px;
    }

    .menu-button {
      font-weight: 600;
    }
  }

  .main-container {
    display: flex;
    flex-direction: column;
    width: 100%;
    height: 100%;

    .header {
      // height: 4rem;
      display: flex;
      flex-direction: column;
      align-items: center;
      padding: 10px;
      margin: 10px;
      margin-bottom: 0;
      margin-left: 5px;
      border-radius: 10px;

      .title {
        display: flex;
        align-items: center;
        width: 100%;

        .menu-button {
          font-weight: 600;
        }

        .header-content {
          flex: 1;
          text-align: center;
        }
      }

      .controller {
        display: flex;
        align-items: center;
        justify-content: space-between;
        width: 100%;

        .tab {
          display: flex;
          flex: 1;
          flex-wrap: wrap;
          gap: 10px;
          justify-content: flex-start;
          width: 100%;
          margin-top: 10px;

          .el-tag {
            cursor: pointer;
          }
        }

        .function {
          display: flex;
          gap: 10px;
          justify-content: flex-end;
          margin-top: 10px;

          svg {
            height: 1.2rem;
          }
        }
      }
    }

    .main {
      flex: 1;
      padding: 10px;
      margin: 10px;
      margin-top: 0;
      margin-left: 5px;
      overflow: auto;
      background-color: $ac-background;
      border: 1px solid $ac-border-color;
      border-radius: 10px;
    }
  }
}

@media (width >= 768px) {
  .menu-button {
    display: none;
  }
}
</style>
