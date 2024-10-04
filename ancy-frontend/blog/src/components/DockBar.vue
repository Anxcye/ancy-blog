<template>
  <div class="dock-bar">
    <div class="dock-bar-item" v-for="item in items" :key="item.key">
      <router-link :to="item.path" custom v-slot="{ isActive, navigate }">
        <div
          @click="navigate"
          :class="['flex items-center justify-center gap-1', isActive ? 'bg-blue-500' : '']"
        >
          <icon :component="item.icon" />
          {{ item.label }}
        </div>
      </router-link>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useDockStore } from '@/stores/dock'
const items = useDockStore().items
import Icon from '@ant-design/icons-vue'
</script>

<style scoped lang="scss">
.dock-bar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  width: 100%;
  height: 100%;

  .ant-btn {
    border: none;
  }

  .btn-active {
    background-color: #1677ff;
  }

  .btn-inactive {
    background-color: #ff0000;
  }
}
</style>
