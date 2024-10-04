<template>
  <div class="dock-bar">
    <div class="dock-bar-item" v-for="item in items" :key="item.key">
      <!-- <a-button :icon="item.icon" :href="item.path" :class="setClass(item.key)">
        {{ item.label }}
      </a-button> -->
      <router-link
        :to="item.path"
        exact-active-class="btn-active"
        :class="{ 'btn-inactive': router.currentRoute.value.path !== item.path }"
      >
        {{ item.label }}
      </router-link>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import router from '@/router'
import { useDockStore } from '@/stores/dock'

const currentKey = ref('home1')
const items = useDockStore().items

const setClass = (key: string) => {
  return currentKey.value === key ? 'btn-active' : 'btn-inactive'
}
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
