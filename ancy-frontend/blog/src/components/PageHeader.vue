<template>
  <div
    :class="[
      'w-screen flex justify-center left-0 top-0 right-0 z-10 bg-transparent border-b-0 fixed',
      isScrolled ? 'bg-bg-color-1 backdrop-blur-sm border-b border-gray-bg' : '',
    ]"
  >
    <div
      class="w-screen flex items-center justify-between h-ac-header px-2 max-w-3xl mx-auto md:px-0"
    >
      <div class="block md:hidden">
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
      <div class="hidden md:block">
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

<style scoped lang="scss"></style>
