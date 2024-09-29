<template>
  <div class="aside-menu" v-if="treeData.length > 0">
    <el-menu @select="handleSelect" :default-active="initialId.toString()">
      <template v-for="item in treeData" :key="item.id">
        <template v-if="item.children.length === 0">
          <el-menu-item :index="item.id.toString()" class="menu-item">
            <el-icon>
              <component :is="iconComponent(item.icon)" />
            </el-icon>

            <span>{{ item.menuName }}</span>
          </el-menu-item>
        </template>

        <el-sub-menu v-else :index="item.id.toString()" class="menu-menu">
          <template #title>
            <el-icon>
              <component :is="iconComponent(item.icon)" />
            </el-icon>
            <span>{{ item.menuName }}</span>
          </template>

          <template v-for="child in item.children" :key="child.id">
            <el-menu-item :index="child.id.toString()" class="menu-item">
              <el-icon>
                <component :is="iconComponent(child.icon)" />
              </el-icon>
              <span>{{ child.menuName }}</span>
            </el-menu-item>
          </template>
        </el-sub-menu>
      </template>
    </el-menu>
  </div>
</template>

<script setup lang="ts">
import router from '@/router'
import { useRouteStore } from '@/stores/modules/route'
import { useTabStore } from '@/stores/modules/tab'
import * as ElementPlusIconsVue from '@element-plus/icons-vue'

const routeStore = useRouteStore()
const tabStore = useTabStore()
const treeData = routeStore.routes

const handleSelect = (key: string) => {
  const menu = routeStore.getRoutes(key)
  router.replace(`/${menu?.path}`)
}

let initialPath = tabStore.currentTab

let initialId = routeStore.getIdByPath(initialPath)

const iconComponent = (icon: string) => {
  return ElementPlusIconsVue[icon as keyof typeof ElementPlusIconsVue]
}
</script>

<style scoped lang="scss">
.aside-menu {
  width: 100%;
  font-size: 14px;
  font-weight: 600;
  background-color: white;
}
</style>
