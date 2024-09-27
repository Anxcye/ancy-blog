<template>
  <div class="aside-menu" v-if="treeData.length > 0">
    <el-menu @select="handleSelect">
      <template v-for="item in treeData" :key="item.id">
        <template v-if="item.children.length === 0">
          <el-menu-item :index="item.id.toString()">
            <el-icon>
              <component :is="iconComponent(item.icon)" />
            </el-icon>

            <span>{{ item.menuName }}</span>
          </el-menu-item>
        </template>

        <el-sub-menu v-else :index="item.id.toString()">
          <template #title>
            <el-icon>
              <component :is="iconComponent(item.icon)" />
            </el-icon>
            <span>{{ item.menuName }}</span>
          </template>

          <template v-for="child in item.children" :key="child.id">
            <el-menu-item :index="child.id.toString()">
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
import * as ElementPlusIconsVue from '@element-plus/icons-vue'
// import { useRouter } from 'vue-router'

const routeStore = useRouteStore()
// const router = useRouter()
const treeData = routeStore.routes

const handleSelect = (key: string) => {
  // router.push(key)
  const menu = routeStore.getRoutes(key)
  router.replace(`/${menu?.path}`)
}

const iconComponent = (icon: string) => {
  return ElementPlusIconsVue[icon as keyof typeof ElementPlusIconsVue]
}
</script>

<style scoped lang="scss">
.aside-menu {
  height: 100%;
  width: 200px;
  // background-color: #fff;
}
</style>
