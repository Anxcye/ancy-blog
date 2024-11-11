<template>
  <div class="flex flex-col justify-between items-center px-2 min-h-screen-nh">
    <div></div>
    <div
      class="w-full items-center flex flex-col justify-center md:flex-row md:max-w-3xl md:mx-auto md:justify-between"
    >
      <div class="flex flex-col items-center md:items-start md:mt-10">
        <div class="text-4xl font-bold">{{ baseInfoStore.getGreeting() }}</div>
        <div class="text-2xl mt-2">{{ baseInfoStore.getRole() }}</div>
        <div class="text-xl font-thin mt-4 text-gray">{{ baseInfoStore.getPhilosophy() }}</div>
        <div class="flex flex-row gap-2 mt-2 md:mt-7">
          <div v-for="item in baseInfoStore.getBadge()" :key="item.index">
            <a-tooltip>
              <template #title>
                <span>{{ item.title }}</span>
              </template>
              <a :href="item.url" target="_blank">
                <a-avatar :size="32" :src="item.img" alt="badge" />
              </a>
            </a-tooltip>
          </div>
        </div>
      </div>
      <div>
        <a-avatar
          :src="baseInfoStore.getAvatar()"
          class="shadow-2xl w-48 h-48 mt-12 md:w-64 md:h-64"
          alt="avatar"
        />
      </div>
    </div>
    <div>
      <DownOutlined class="text-xl text-gray" />
    </div>
  </div>
</template>

<script setup lang="ts">
import { useBaseInfoStore } from '@/stores/baseInfo'
import { DownOutlined } from '@ant-design/icons-vue'

const baseInfoStore = useBaseInfoStore()
</script>

<style scoped lang="scss"></style>
