<template>
  <div class="common-layout">
    <el-container>
      <el-aside
        :class="{ 'mobile-aside': isMobile, 'aside-content': true }"
        :style="{ width: asideWidth }"
      >
        <slot name="aside"></slot>
      </el-aside>
      <slot name="main"></slot>
    </el-container>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted, computed } from 'vue'

const isMobile = ref(false)
const emit = defineEmits(['update:isAsideVisible'])

const props = defineProps({
  isAsideVisible: {
    type: Boolean,
    default: true,
  },
})

const asideWidth = computed(() => {
  if (isMobile.value) {
    return props.isAsideVisible ? '200px' : '0'
  }
  return '200px'
})

const checkMobile = () => {
  isMobile.value = window.innerWidth <= 768
  if (isMobile.value) {
    emit('update:isAsideVisible', false)
  }
}

onMounted(() => {
  checkMobile()
  window.addEventListener('resize', checkMobile)
})

onUnmounted(() => {
  window.removeEventListener('resize', checkMobile)
})
</script>

<style scoped lang="scss">
.common-layout {
  height: 100vh;

  .el-container {
    height: 100%;
  }

  .el-aside {
    transition: width 0.3s;

    &.aside-content {
      height: 100vh;
    }

    &.mobile-aside {
      position: fixed;
      top: 0;
      left: 0;
      height: 100%;
      z-index: 1000;
      overflow-x: hidden;
    }
  }

  @media (max-width: 768px) {
    .el-aside {
      width: 0;
    }
  }
}
</style>
