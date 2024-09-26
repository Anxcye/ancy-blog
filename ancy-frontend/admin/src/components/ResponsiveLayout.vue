<template>
  <div class="common-layout">
    <el-container>
      <el-aside
        :class="{ 'mobile-aside': isMobile, 'aside-content': true }"
        :style="{ width: asideWidth }"
      >
        <slot name="aside" class="aside-content"></slot>
      </el-aside>
      <el-container class="main-container">
        <el-header>
          <slot name="header"></slot>
        </el-header>
        <el-main class="main-content">
          <slot name="main"></slot>
        </el-main>
      </el-container>
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
    return props.isAsideVisible ? '80%' : '0'
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
.aside-content {
  height: 100vh;
}
.main-container {
  display: flex;
  flex-direction: column;
  height: 100vh;
}
.main-content {
  flex: 1;
}

.common-layout {
  height: 100vh;
}

.el-aside {
  transition: width 0.3s;
}

.el-header {
  display: flex;
  align-items: center;
}

.mobile-aside {
  position: fixed;
  top: 0;
  left: 0;
  height: 100%;
  z-index: 1000;
  overflow-x: hidden;

  .aside-content {
    padding: 20px;
  }

  .close-btn {
    position: absolute;
    top: 10px;
    right: 10px;
  }
}

@media (max-width: 768px) {
  .el-aside {
    width: 0;
  }
}
</style>
