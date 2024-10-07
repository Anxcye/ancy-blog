<template>
  <div class="content">
    <div class="title">回溯</div>
    <div class="summary">{{ total }} 条记录</div>
    <TimelineList :list="timelineList" timeField="operateTime" :total="total">
      <template #item="{ item }">
        <div
          class="flex flex-col hover:bg-primary-bg-1 p-2 rounded-lg hover:shadow-md hover:scale-105 transition-all"
        >
          <router-link :to="getInfo(item).path">
            <div class="flex flex-row items-center gap-2 justify-between">
              <div
                class="text-lg font-medium flex flex-row items-center gap-2 justify-between w-full"
              >
                {{ getInfo(item).summary }}
                <Icon :component="getInfo(item).icon" />
              </div>
            </div>
            <div class="text-lg text-gray">{{ item.summary }}</div>
          </router-link>
        </div>
      </template>
    </TimelineList>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted, type Component } from 'vue'
import { reqTimelinePage } from '@/api/timeline/index'
import type { TimelineData, TimelinePageParam } from '@/api/timeline/type'
import { handleScroll } from '@/utils/handleScroll'
import Icon from '@ant-design/icons-vue'
import { BookOutlined, CalendarOutlined, ToolOutlined } from '@ant-design/icons-vue'

const timelineList = ref<TimelineData[]>([])
const total = ref(0)
const loading = ref(false)
const params = ref<TimelinePageParam>({
  pageNum: 1,
  pageSize: 10,
})
interface TimelineItem {
  summary: string
  icon: Component
  path: string
}

const getInfo = (item: TimelineData): TimelineItem => {
  switch (item.methodName) {
    case 'addArticle':
      return {
        icon: BookOutlined,
        summary: '发布了一篇文章',
        path: '/article/' + item.returnValue,
      }
    case 'addNote':
      return {
        icon: CalendarOutlined,
        summary: '添加了一篇日志',
        path: '/note',
      }
    case 'addProject':
      return {
        icon: ToolOutlined,
        summary: '记下了一个项目',
        path: '/project/' + item.returnValue,
      }
    default:
      return {
        icon: BookOutlined,
        summary: '未知操作',
        path: '/',
      }
  }
}

const getTimelineList = async () => {
  if (loading.value) return
  loading.value = true
  try {
    const res = await reqTimelinePage(params.value)
    timelineList.value = [...timelineList.value, ...res.data.rows]
    total.value = res.data.total
    params.value.pageNum++
  } finally {
    loading.value = false
  }
}

const scroll = () => {
  handleScroll(getTimelineList, loading.value, total.value === timelineList.value.length)
}

onMounted(async () => {
  await getTimelineList()
  scroll()
  window.addEventListener('scroll', scroll)
})
onUnmounted(() => {
  window.removeEventListener('scroll', scroll)
})
</script>

<style scoped></style>
