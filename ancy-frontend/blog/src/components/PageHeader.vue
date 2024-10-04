<template>
  <div class="page-header">
    <div class="menu"><a-button shape="circle" :icon="h(MenuOutlined)" /></div>
    <div class="avatar">
      <a-avatar :size="48" :src="baseInfoStore.baseInfo?.avatar" />
    </div>
    <div class="dock-bar">
      <DockBar />
    </div>
    <div class="theme">
      <a-button shape="circle" @click="toggleTheme">
        <template #icon>
          <img :src="themeIcon" alt="theme" style="width: 24px; height: 24px" />
        </template>
      </a-button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { h, onMounted, ref } from 'vue'
import { MenuOutlined } from '@ant-design/icons-vue'
import { useBaseInfoStore } from '@/stores/baseInfo'
import sunIcon from '@/assets/svg/sun.svg'
import moonIcon from '@/assets/svg/moon.svg'
import DockBar from './DockBar.vue'
const baseInfoStore = useBaseInfoStore()
const themeIcon = ref(sunIcon)

const toggleTheme = () => {
  themeIcon.value = themeIcon.value === sunIcon ? moonIcon : sunIcon
}

onMounted(async () => {
  await baseInfoStore.reqBaseInfo()
})
</script>

<style scoped lang="scss">
.page-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  height: 64px;
  padding: 0 24px;
  background-color: #fff;

  // pc
  @media (min-width: 768px) {
    .menu {
      display: none;
    }
  }

  // mobile
  @media (max-width: 768px) {
    .menu {
      display: block;
    }
    .dock-bar {
      display: none;
    }
  }
}
</style>
