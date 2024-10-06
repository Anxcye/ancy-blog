<template>
  <div class="max-w-3xl mx-auto px-3">
    <TimelineList :list="projectList" timeField="beginDate">
      <template #item="{ item }">
        <div
          class="flex flex-col hover:bg-primary-bg-1 p-2 rounded-lg hover:shadow-md hover:scale-105 transition-all"
        >
          <router-link :to="`/project/${item.id}`" class="flex flex-row items-center gap-2 pb-2">
            <a-avatar :src="item.thumbnail" shape="square" :size="60" />
            <div class="flex flex-col w-full">
              <div class="flex flex-row items-center gap-2 justify-between">
                <div class="text-lg font-medium">
                  {{ item.title }}
                </div>
                <a-tag v-if="item.isTop === '1'" class="bg-primary text-white">置顶</a-tag>
              </div>
              <div class="text-sm text-gray">{{ item.summary }}</div>
              <div class="flex flex-row items-center gap-2 justify-between">
                <div>
                  <a-tag v-if="item.type === '0'" class="bg-green-500 text-white">Active</a-tag>
                  <a-tag v-if="item.type === '1'">Archived</a-tag>
                </div>
              </div>
            </div>
          </router-link>
          <div class="flex flex-row items-center gap-2 justify-end">
            <a v-if="item.srcUrl" :href="item.srcUrl" target="_blank">开源地址</a>
            <a v-if="item.displayUrl" :href="item.displayUrl" target="_blank">展示地址</a>
            <router-link :to="`/project/${item.id}`">查看详情</router-link>
          </div>
        </div>
      </template>
    </TimelineList>
  </div>
</template>

<script setup lang="ts">
import { reqProjectList } from '@/api/project'
import type { ProjectData } from '@/api/project/type'
import { onMounted, ref } from 'vue'

const projectList = ref<ProjectData[]>([])

const getProjectList = async () => {
  const res = await reqProjectList()
  projectList.value = res.data
}

onMounted(async () => {
  await getProjectList()
})
</script>

<style scoped></style>
