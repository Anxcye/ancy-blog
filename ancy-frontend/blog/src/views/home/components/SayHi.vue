<template>
  <div class="h-screen-nh flex flex-col justify-between items-center">
    <div></div>
    <div
      class="w-full items-center flex flex-col justify-center md:flex-row md:max-w-3xl md:mx-auto md:justify-between"
    >
      <div class="flex flex-col items-center md:items-start md:mt-10">
        <div class="text-4xl font-bold">{{ greeting }}</div>
        <div class="text-2xl mt-2">{{ role }}</div>
        <div class="text-xl font-thin mt-4">{{ philosophy }}</div>
        <div class="flex flex-row gap-2 mt-2 md:mt-7">
          <div v-for="item in badge" :key="item.index">
            <a-tooltip color="white">
              <template #title>
                <span class="text-black">{{ item.title }}</span>
              </template>
              <a :href="item.url" target="_blank">
                <a-avatar :size="32" :src="item.img" />
              </a>
            </a-tooltip>
          </div>
        </div>
      </div>
      <div>
        <a-avatar :src="avatar" class="shadow-2xl w-48 h-48 mt-12 md:w-64 md:h-64" />
      </div>
    </div>
    <div>
      <DownOutlined class="text-xl text-gray-500" />
    </div>
  </div>
</template>

<script setup lang="ts">
import { useBaseInfoStore } from '@/stores/baseInfo'
import { DownOutlined } from '@ant-design/icons-vue'

const baseInfoStore = useBaseInfoStore()

const avatar = baseInfoStore.getAvatar()
const greeting = baseInfoStore.getGreeting()
const role = baseInfoStore.getRole()
const philosophy = baseInfoStore.getPhilosophy()
const badge = baseInfoStore.getBadge()
</script>

<style scoped lang="scss"></style>
