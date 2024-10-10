<template>
  <div class="content pb-4">
    <CommentForm :id="id" :postType="type" @submit="submitComment($event)" />
    <div class="text-2xl font-medium my-6">
      {{ total }}条评论 ·
      <span class="text-gray text-lg">其中{{ total - rootTotal }}条回复</span>
    </div>
    <div v-for="item in sortedCommentList" :key="item.id" class="mb-8">
      <CommentItem :comment="item" :postType="type" :articleId="id" />
    </div>
    <a-pagination
      v-if="pagination"
      :total="rootTotal"
      @change="getCommentListPagination"
      v-model:current="params.pageNum"
      :page-size="params.pageSize"
      class="w-full mx-auto"
      hideOnSinglePage
    />
  </div>
</template>

<script setup lang="ts">
import {
  commentArticleTotal,
  commentNoteTotal,
  reqCommentByArticleId,
  reqCommentNote,
} from '@/api/comment'
import type { CommentData } from '@/api/comment/type'
import { handleScroll } from '@/utils/handleScroll'
import { computed, onMounted, onUnmounted, ref } from 'vue'

const props = withDefaults(
  defineProps<{
    id: number
    type: 'Article' | 'Note'
    pagination?: boolean
  }>(),
  {
    pagination: false,
  },
)
const commentList = ref<CommentData[]>([])
const rootTotal = ref(0)
const params = ref({
  id: props.id,
  pageNum: 1,
  pageSize: 5,
})
const total = ref(0)
const loading = ref(false)
const sortedCommentList = computed(() => {
  return [...commentList.value].sort((a, b) => {
    if (a.isTop === '1' && (b.isTop === '0' || !b.isTop)) {
      return -1
    }
    if ((!a.isTop || a.isTop === '0') && b.isTop === '1') {
      return 1
    }
    return new Date(b.createTime).getTime() - new Date(a.createTime).getTime()
  })
})

const getTotal = async () => {
  let res
  if (props.type === 'Article') {
    res = await commentArticleTotal(props.id)
  } else {
    res = await commentNoteTotal(props.id)
  }
  total.value = res.data
}

const getCommentListPagination = async () => {
  let res
  if (props.type === 'Article') {
    res = await reqCommentByArticleId(params.value)
  } else {
    res = await reqCommentNote(params.value)
  }
  commentList.value = res.data.rows
  rootTotal.value = res.data.total
}

const getCommentList = async () => {
  if (loading.value) return
  loading.value = true
  try {
    let res
    if (props.type === 'Article') {
      res = await reqCommentByArticleId(params.value)
    } else {
      res = await reqCommentNote(params.value)
    }
    commentList.value = [...commentList.value, ...res.data.rows]
    rootTotal.value = res.data.total
    params.value.pageNum++
  } finally {
    loading.value = false
  }
}

const submitComment = (comment: CommentData) => {
  commentList.value.unshift(comment)
  rootTotal.value++
  total.value++
}

const scroll = () => {
  handleScroll(getCommentList, loading.value, rootTotal.value === commentList.value.length)
}

onMounted(() => {
  getTotal()
  if (props.pagination) {
    getCommentListPagination()
  } else {
    getCommentList()
    window.addEventListener('scroll', scroll)
  }
})

onUnmounted(() => {
  if (props.pagination) return
  window.removeEventListener('scroll', scroll)
})
</script>

<style scoped lang="scss"></style>
