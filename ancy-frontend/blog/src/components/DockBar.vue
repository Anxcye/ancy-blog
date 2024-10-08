<template>
  <div
    :class="[
      'flex justify-between items-center gap-3 rounded-full px-3 py-1 transition-all duration-300',
      props.isScrolled ? '' : 'shadow border border-gray-bg',
    ]"
  >
    <div class="dock-bar-item" v-for="item in items" :key="item.key">
      <a-dropdown>
        <div
          @click="navigate(item.path)"
          :class="[
            'flex items-center justify-center gap-1 text-sm cursor-pointer ',
            isActive(item.group) ? 'text-primary' : 'text-gray',
          ]"
        >
          <icon :component="item.icon" v-if="isActive(item.group)" style="font-size: 12px" />
          {{ getActiveLabel(item) }}
        </div>
        <template #overlay v-if="item.children">
          <a-menu>
            <a-menu-item v-for="child in item.children" :key="child.key">
              <router-link :to="child.path" class="hover:text-primary">
                {{ child.label }}
              </router-link>
            </a-menu-item>
          </a-menu>
        </template>
      </a-dropdown>
    </div>
  </div>

  <div>
    <a-drawer
      height="60vh"
      placement="bottom"
      class="rounded-t-2xl overflow-hidden"
      :open="props.drawerOpen"
      @close="menuClose"
      :closable="false"
      :maskStyle="{ background: 'rgba(0, 0, 0, 0.1)' }"
    >
      <div v-for="item in items" :key="item.key" class="my-4">
        <div
          @click="navigate(item.path)"
          :class="[
            'gap-2 font-medium text-lg w-full flex flex-row items-center',
            isActive(item.group) ? 'text-primary' : 'text-gray',
          ]"
        >
          <icon :component="item.icon" />
          <div class="">{{ item.label }}</div>
        </div>
        <div class="grid grid-cols-2 gap-2">
          <div v-for="child in item.children" :key="child.key">
            <router-link :to="child.path" custom v-slot="{ isExactActive }">
              <div
                @click="navigate(child.path)"
                :class="[
                  isExactActive ? 'text-lg' : 'text-sm',
                  isExactActive ? 'text-primary' : 'text-gray',
                ]"
              >
                <!-- :style="{
                  color: isExactActive ? getPrimaryColor() : '',
                }" -->
                {{ child.label }}
              </div>
            </router-link>
          </div>
        </div>
      </div>
    </a-drawer>
  </div>
</template>

<script setup lang="ts">
import { useDockStore } from '@/stores/dock'
import Icon from '@ant-design/icons-vue'
import router from '@/router'
const dockStore = useDockStore()
const items = dockStore.items

const isActive = (group: string) => router.currentRoute.value.meta.group === group

const props = defineProps<{
  isScrolled: boolean
  drawerOpen: boolean
}>()

const emit = defineEmits(['update:drawerOpen'])

const menuClose = () => {
  emit('update:drawerOpen', false)
}

const navigate = (path?: string) => {
  if (path) {
    router.push(path)
    if (props.drawerOpen) {
      emit('update:drawerOpen', false)
    }
  }
}

const getActiveLabel = (item: any) => {
  if (item.children) {
    const activeChild = item.children.find((child: any) =>
      router.currentRoute.value.path.startsWith(child.path),
    )
    return activeChild ? activeChild.label : item.label
  }
  return item.label
}
</script>

<style scoped lang="scss"></style>
