<template>
  <div class="content">
    <div class="title">日志</div>
    <div class="summary">{{ total }} 条日志</div>
    <TimelineList :list="noteList" :total="total">
      <template #item="{ item }">
        <div
          class="flex flex-col hover:bg-primary-bg-1 p-2 rounded-lg hover:shadow-md hover:scale-105 transition-all"
        >
          <MdViewer :content="item.content" />
          <div class="text-sm text-gray flex flex-row items-center gap-2">
            <span>
              <icon :component="EyeOutlined" />
              {{ item.viewCount }}
            </span>
            <div>
              <CommentOutlined @click="openComment(item.id)" class="text-sm text-gray" />
            </div>
          </div>
        </div>
      </template>
    </TimelineList>
    <a-modal
      wrap-class-name="full-modal"
      v-model:open="commentOpen"
      :width="800"
      :footer="null"
      :maskStyle="{ backdropFilter: 'blur(10px)' }"
    >
      <template #closeIcon></template>
      <div class="bg-bg-color py-2 rounded-lg">
        <DisplayComment :id="NoteId" type="Note" class="w-full" :pagination="true" />
      </div>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { reqNotePage } from '@/api/note'
import type { NoteData, NotePageParams } from '@/api/note/type'
import { onMounted, onUnmounted, ref } from 'vue'
import MdViewer from '@/components/MdViewer.vue'
import { EyeOutlined, CommentOutlined } from '@ant-design/icons-vue'
import Icon from '@ant-design/icons-vue'
import { handleScroll } from '@/utils/handleScroll'
import DisplayComment from '@/components/comment/DisplayComment.vue'

const noteList = ref<NoteData[]>([])
const params = ref<NotePageParams>({
  pageNum: 1,
  pageSize: 10,
})
const total = ref(0)
const loading = ref(false)
const commentOpen = ref(false)
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

const NoteId = ref(0)
const scroll = () => {
  handleScroll(getNoteList, loading.value, total.value === noteList.value.length)
}

const openComment = (id: number) => {
  NoteId.value = id
  commentOpen.value = true
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

<style scoped lang="scss">
.full-modal {
  .ant-modal {
    max-width: 100% !important;
    top: 0;
    padding-bottom: 0;
    margin: 0;
  }
  .ant-modal-content {
    display: flex;
    flex-direction: column;
    height: calc(100vh);
  }
  .ant-modal-body {
    flex: 1;
  }
}
</style>
