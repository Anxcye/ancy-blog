<template>
  <div class="flex justify-between items-center gap-2">
    <div class="dock-bar-item" v-for="item in items" :key="item.key">
      <router-link :to="item.path" custom v-slot="{ navigate }">
        <a-dropdown>
          <div
            @click="navigate"
            :class="[
              'flex items-center justify-center gap-1',
              isActive(item.group) ? 'text-blue-500' : '',
            ]"
          >
            <icon :component="item.icon" v-if="isActive(item.group)" />
            {{ item.label }}
          </div>
          <template #overlay v-if="item.children">
            <a-menu>
              <a-menu-item v-for="child in item.children" :key="child.key">
                <router-link :to="child.path">{{ child.label }}</router-link>
              </a-menu-item>
            </a-menu>
          </template>
        </a-dropdown>
      </router-link>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useDockStore } from '@/stores/dock'
import Icon from '@ant-design/icons-vue'
import router from '@/router'

const dockStore = useDockStore()
const items = dockStore.items

const isActive = (group: string) => {
  console.log('router', router.currentRoute.value.meta.group)
  // console.log('group', group)
  return router.currentRoute.value.meta.group === group
}
</script>

<style scoped lang="scss"></style>
