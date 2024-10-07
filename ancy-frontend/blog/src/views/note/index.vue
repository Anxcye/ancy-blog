<template>
  <div class="max-w-3xl mx-auto px-3">
    <TimelineList :list="noteList" :total="total">
      <template #item="{ item }">
        <div
          class="flex flex-col hover:bg-primary-bg-1 p-2 rounded-lg hover:shadow-md hover:scale-105 transition-all"
        >
          <MdViewer :content="item.content" />
          <div class="text-sm text-gray">
            <span>
              <icon :component="EyeOutlined" />
              {{ item.viewCount }}
            </span>
          </div>
        </div>
        <div></div>
      </template>
    </TimelineList>
  </div>
</template>

<script setup lang="ts">
import { reqNotePage } from '@/api/note'
import type { NoteData, NotePageParams } from '@/api/note/type'
import { onMounted, onUnmounted, ref } from 'vue'
import MdViewer from '@/components/MdViewer.vue'
import { EyeOutlined } from '@ant-design/icons-vue'
import Icon from '@ant-design/icons-vue'
import { handleScroll } from '@/utils/handleScroll'

const noteList = ref<NoteData[]>([])
const params = ref<NotePageParams>({
  pageNum: 1,
  pageSize: 10,
})
const total = ref(0)
const loading = ref(false)

const getNoteList = async () => {
  if (loading.value) return
  loading.value = true
  try {
    const res = await reqNotePage(params.value)
    noteList.value = [...noteList.value, ...res.data.rows]
    total.value = res.data.total
    params.value.pageNum++
  } finally {
    loading.value = false
  }
}

const scroll = () => {
  handleScroll(getNoteList, loading.value, total.value === noteList.value.length)
}

onMounted(async () => {
  await getNoteList()
  scroll()
  window.addEventListener('scroll', scroll)
})
onUnmounted(() => {
  window.removeEventListener('scroll', scroll)
})
</script>

<style scoped></style>
