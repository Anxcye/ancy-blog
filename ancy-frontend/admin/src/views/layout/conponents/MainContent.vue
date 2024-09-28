<template>
  <div>
    <router-view v-slot="{ Component }">
      <transition name="fade">
        <component :is="Component" v-if="refresh" />
      </transition>
    </router-view>
  </div>
</template>

<script setup lang="ts">
import { nextTick, ref, watch } from 'vue'
import { useLayoutStore } from '@/stores/modules/layout'

const layoutStore = useLayoutStore()
const refresh = ref(false)

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
