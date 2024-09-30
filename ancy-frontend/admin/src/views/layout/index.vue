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
      margin-bottom: 0px;
      border-radius: 10px;
      padding: 10px;
      // height: 4rem;
      display: flex;
      flex-direction: column;
      align-items: center;

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
        width: 100%;
        display: flex;
        justify-content: space-between;
        align-items: center;

        .tab {
          width: 100%;
          display: flex;
          justify-content: flex-start;
          margin-top: 10px;
          gap: 10px;
          flex-wrap: wrap;
          flex: 1;

          .el-tag {
            cursor: pointer;
          }
        }

        .function {
          display: flex;
          justify-content: flex-end;
          margin-top: 10px;
          gap: 10px;

          svg {
            height: 1.2rem;
          }
        }
      }
    }

    .main {
      margin: 10px;
      margin-left: 5px;
      margin-top: 0px;
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
