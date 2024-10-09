<template>
  <div class="content">
    <CommentForm :id="id" :postType="type" @submit="submitComment($event)" />
    <div class="text-2xl font-medium mb-4">{{ total }}条评论</div>
    <div v-for="item in sortedCommentList" :key="item.id" class="mb-8">
      <CommentItem :comment="item" :postType="type" :articleId="id" />
    </div>
  </div>
</template>

<script setup lang="ts">
import { reqCommentByArticleId, reqCommentNote } from '@/api/comment'
import type { CommentData } from '@/api/comment/type'
import { computed, onMounted, ref } from 'vue'

const props = defineProps<{
  id: number
  type: 'Article' | 'Note'
}>()
const commentList = ref<CommentData[]>([])
const total = ref(0)
const params = ref({
  id: props.id,
  pageNum: 1,
  pageSize: 100,
})

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

const getCommentList = async () => {
  let res
  if (props.type === 'Article') {
    res = await reqCommentByArticleId(params.value)
  } else {
    res = await reqCommentNote(params.value)
  }
  commentList.value = res.data.rows
  total.value = res.data.total
}

const submitComment = (comment: CommentData) => {
  commentList.value.unshift(comment)
  total.value++
}

onMounted(() => {
  getCommentList()
})
</script>

<style scoped lang="scss"></style>
