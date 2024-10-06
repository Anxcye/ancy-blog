<template>
  <div class="app-container" :class="{ scrolled: isScrolled }">
    <div class="page-header">
      <div class="menu">
        <a-button
          shape="circle"
          :icon="h(MenuOutlined)"
          class="flex items-center justify-center"
          @click="menuClick"
        />
      </div>
      <div class="avatar">
        <a-avatar :size="36" :src="baseInfoStore.baseInfo?.avatar" @click="avatarClick" />
      </div>
      <div class="dock-bar">
        <DockBar :isScrolled="isScrolled" :drawerOpen="open" @update:drawerOpen="menuClose" />
      </div>
      <div class="flex items-center justify-center">
        <a-button shape="circle" @click="themeClick">
          <template #icon>
            <img
              :src="themeIcon"
              alt="theme"
              class="w-6 h-6 mx-auto"
              :style="{
                filter: colorStore.isDarkMode
                  ? 'invert(100%) sepia(100%) saturate(0%) hue-rotate(249deg) brightness(118%) contrast(119%)'
                  : 'none',
              }"
            />
          </template>
        </a-button>
      </div>
    </div>
  </div>
  <ConfigProvider />
</template>

<script setup lang="ts">
import { h, onMounted, ref, onUnmounted } from 'vue'
import { MenuOutlined } from '@ant-design/icons-vue'
import { useBaseInfoStore } from '@/stores/baseInfo'
import sunIcon from '@/assets/svg/sun.svg'
import moonIcon from '@/assets/svg/moon.svg'
import DockBar from './DockBar.vue'
import router from '@/router'
import ConfigProvider from '@/provider/ConfigProvider.vue'
import { useColorStore } from '@/stores/color'

const colorStore = useColorStore()
const baseInfoStore = useBaseInfoStore()
const themeIcon = ref(sunIcon)
const isScrolled = ref(false)
const open = ref(false)
const handleScroll = () => {
  isScrolled.value = window.scrollY > 0
}

const menuClick = () => {
  open.value = true
}
const menuClose = () => {
  open.value = false
}
const avatarClick = () => {
  router.push('/')
}
const themeClick = () => {
  themeIcon.value = colorStore.isDarkMode ? sunIcon : moonIcon

  colorStore.toggleDarkMode()
}

onMounted(async () => {
  await baseInfoStore.reqBaseInfo()
  window.addEventListener('scroll', handleScroll)
})

onUnmounted(() => {
  window.removeEventListener('scroll', handleScroll)
})
</script>

<style scoped lang="scss">
.app-container {
  width: 100vw;
  display: flex;
  justify-content: center;
  left: 0;
  top: 0;
  right: 0;
  z-index: 1000;
  background-color: transparent;
  position: fixed;
  border-bottom: none;

  &.scrolled {
    // background-color: rgba(255, 255, 255, 0.11);
    background-color: var(--background-color-1);
    backdrop-filter: blur(10px);
    border-bottom: 1px solid #e8e8e8;
  }

  .page-header {
    width: 100%;
    display: flex;
    align-items: center;
    justify-content: space-between;
    height: $ac-header-height;

    // pc
    @media (min-width: 768px) {
      max-width: 800px;
      .menu {
        display: none;
      }
    }

    // mobile
    @media (max-width: 768px) {
      padding: 0 16px;
      .menu {
        display: block;
      }
      .dock-bar {
        display: none;
      }
    }
  }
}
</style>
