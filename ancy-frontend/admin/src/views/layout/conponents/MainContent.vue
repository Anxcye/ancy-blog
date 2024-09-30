<template>
  <div>
    <router-view v-slot="{ Component }">
      <transition name="fade">
        <keep-alive v-if="refresh">
          <component :is="Component" v-if="refresh" :key="tabStore.currentTab?.path" />
        </keep-alive>
      </transition>
    </router-view>
  </div>
</template>

<script setup lang="ts">
import { nextTick, ref, watch } from 'vue'
import { useLayoutStore } from '@/stores/modules/layout'
import { useTabStore } from '@/stores/modules/tab'

const layoutStore = useLayoutStore()
const tabStore = useTabStore()
const refresh = ref(true)

watch(
  () => layoutStore.refresh,
  () => {
    refresh.value = false
    nextTick(() => {
      refresh.value = true
    })
  },
)
</script>

<style scoped lang="scss">
.fade-enter-active {
  transition: all 0.5s;
  transform: translateX(0);
}

.fade-enter-from {
  opacity: 0;
  transform: translateX(-100%);
}
.fade-leave-to {
  opacity: 0;
  transform: translateX(100%);
}
</style>
