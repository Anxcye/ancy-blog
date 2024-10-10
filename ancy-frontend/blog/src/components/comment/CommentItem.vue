<template>
  <div class="mb-8">
    <div class="flex flex-row items-end justify-start gap-2">
      <a-avatar
        :src="
          comment.userId === 1
            ? baseInfo.getAvatar()
            : comment.avatar?.length > 0
              ? comment.avatar
              : `https:\/\/cravatar.cn/avatar/${comment.nickname}?d=identicon`
        "
      />
      <div class="flex flex-col flex-1">
        <div class="flex flex-row items-end justify-between">
          <div class="flex flex-row items-end gap-2">
            <div class="font-medium">
              {{ comment.userId === 1 ? baseInfo.getName() : comment.nickname }}
            </div>
            <div class="flex flex-row gap-2 items-end" v-if="comment.toCommentNickname">
              <CaretRightOutlined class="text-lg" />
              <div class="font-medium">{{ comment.toCommentNickname }}</div>
            </div>
            <time-tip :time="comment.createTime" class="text-sm text-gray" />
            <div class="flex flex-row items-center gap-6">
              <div class="flex flex-row items-center gap-1">
                {{ comment.likeCount > 0 ? comment.likeCount : '' }}
                <HeartTwoTone
                  :two-tone-color="visitorStore.checkLiked(comment.id) ? '#eb2f96' : '#cccccc'"
                  @click="handleLike(comment)"
                />
              </div>
              <div class="flex flex-row items-center gap-1">
                {{ total > 0 ? total : '' }}
                <CommentOutlined @click="handleCommentForm(comment.id)" class="text-sm text-gray" />
              </div>
            </div>
          </div>
          <div class="text-sm text-gray">
            <a-tag v-if="comment.isTop === '1'" class="bg-primary text-white">置顶</a-tag>
            <a-tag v-if="comment.status === '1'" color="red">已隐藏</a-tag>
          </div>
        </div>
        <div class="flex flex-row gap-2 justify-between w-full items-end">
          <MdViewer :content="comment.content" class="bg-primary-bg-1 p-2 rounded-xl mt-2" />
          <a-dropdown v-if="adminStore.token" type="text" size="small">
            <Icon :component="MoreOutlined" />
            <template #overlay>
              <div class="flex flex-col gap-2">
                <span>
                  置顶
                  <a-switch
                    v-model:checked="comment.isTop"
                    checked-value="1"
                    un-checked-value="0"
                    @change="handleTopChange(comment)"
                  />
                </span>
                <span>
                  公开
                  <a-switch
                    v-model:checked="comment.status"
                    checked-value="0"
                    un-checked-value="1"
                    @change="handleStatusChange(comment)"
                  />
                </span>
              </div>
            </template>
          </a-dropdown>
        </div>
      </div>
    </div>
    <CommentForm
      v-if="openId === comment.id"
      :postType="postType"
      :parentId="isRoot ? comment.id : comment.parentId"
      :id="articleId"
      :toCommentId="isRoot ? null : comment.id"
      :toCommentNickname="isRoot ? null : comment.nickname"
      class="my-3"
      @submit="handleSubmit($event)"
    />
    <div v-if="total > 0" class="ml-6 border-l-4 border-primary-bg pl-2 md:ml-10">
      <div v-for="child in childrenCommentList" :key="child.id" class="mt-2">
        <CommentItem
          :comment="child"
          :postType="postType"
          :articleId="articleId"
          @childSubmit="handleChildrenSubmit"
        />
      </div>
    </div>
    <a-pagination
      v-if="showMore"
      class="w-full mt-2"
      :total="total"
      v-model:current="getChildrenParams.pageNum"
      @change="loadChildrenComment"
      :page-size="getChildrenParams.pageSize"
      hideOnSinglePage
    />
    <a-button v-if="comment.children?.total > 3" class="w-full mt-2" @click="handleLoadMore">
      {{ showMore ? '收起更多' : '查看更多' }}
    </a-button>
  </div>
</template>

<script setup lang="ts">
import { reqCommentChildrenByParentId, reqCommentLike, reqCommentUpdate } from '@/api/comment'
import type { CommentData, CommentPageParams } from '@/api/comment/type'
import { useAdminStore } from '@/stores/admin'
import {
  MoreOutlined,
  CommentOutlined,
  HeartTwoTone,
  CaretRightOutlined,
} from '@ant-design/icons-vue'
import Icon from '@ant-design/icons-vue'
import { ref, computed } from 'vue'
import { useBaseInfoStore } from '@/stores/baseInfo'
import { useVisitorStore } from '@/stores/visitor'

const emit = defineEmits(['childSubmit'])
const adminStore = useAdminStore()
const baseInfo = useBaseInfoStore()
const visitorStore = useVisitorStore()

const openId = ref(0)
const props = defineProps<{
  comment: CommentData
  articleId: number
  postType: 'Article' | 'Note'
}>()
const statusLoading = ref(false)
const childrenCommentList = ref<CommentData[]>(props.comment.children?.rows || [])
const getChildrenParams = ref<CommentPageParams>({
  id: props.comment.id,
  pageNum: 1,
  pageSize: 10,
})
const showMore = ref(false)
const total = ref(props.comment.children?.total || 0)

const isRoot = computed(() => {
  return props.comment.parentId === -1
})

const loadChildrenComment = async () => {
  const res = await reqCommentChildrenByParentId(getChildrenParams.value)
  childrenCommentList.value = res.data.rows
}

const handleLoadMore = async () => {
  if (showMore.value) {
    showMore.value = false
    childrenCommentList.value = props.comment.children?.rows || []
    return
  }
  await loadChildrenComment()
  showMore.value = true
}

const handleTopChange = async (item: CommentData) => {
  statusLoading.value = true
  await reqCommentUpdate(item.id, { isTop: item.isTop }).finally(() => {
    statusLoading.value = false
  })
}

const handleStatusChange = async (item: CommentData) => {
  statusLoading.value = true
  await reqCommentUpdate(item.id, { status: item.status }).finally(() => {
    statusLoading.value = false
  })
}

const handleSubmit = (comment: CommentData) => {
  if (comment.parentId === props.comment.id) {
    total.value++
    childrenCommentList.value.unshift(comment)
  } else {
    emit('childSubmit', comment)
  }
  openId.value = 0
}

const handleChildrenSubmit = (comment: CommentData) => {
  total.value++
  childrenCommentList.value.unshift(comment)
}

const handleCommentForm = (id: number) => {
  openId.value = openId.value === id ? 0 : id
}

const handleLike = async (item: CommentData) => {
  if (!visitorStore.checkLiked(item.id)) {
    visitorStore.addLiked(item.id)
    item.likeCount = item.likeCount ? item.likeCount + 1 : 1
    await reqCommentLike(item.id, true)
  } else {
    visitorStore.removeLiked(item.id)
    item.likeCount -= 1
    await reqCommentLike(item.id, false)
  }
}
</script>

<style scoped lang="scss"></style>
