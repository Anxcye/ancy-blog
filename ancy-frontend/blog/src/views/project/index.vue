<template>
  <div class="content">
    <div class="title">项目</div>
    <div class="summary">{{ total }} 个项目</div>
    <TimelineList :list="projectList" timeField="beginDate" :total="total">
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
import { reqProjectPage } from '@/api/project'
import type { ProjectData, ProjectPageParams } from '@/api/project/type'
import { handleScroll } from '@/utils/handleScroll'
import { onMounted, onUnmounted, ref } from 'vue'

const projectList = ref<ProjectData[]>([])
const total = ref(0)
const params = ref<ProjectPageParams>({
  pageNum: 1,
  pageSize: 10,
})
const loading = ref(false)

const getProjectList = async () => {
  if (loading.value) return
  loading.value = true
  try {
    const res = await reqProjectPage(params.value)
    projectList.value = [...projectList.value, ...res.data.rows]
    total.value = res.data.total
    params.value.pageNum++
  } finally {
    loading.value = false
  }
}

const scroll = () => {
  handleScroll(getProjectList, loading.value, total.value === projectList.value.length)
}

onMounted(async () => {
  await getProjectList()
  scroll()
  window.addEventListener('scroll', scroll)
})
onUnmounted(() => {
  window.removeEventListener('scroll', scroll)
})
</script>

<style scoped></style>
