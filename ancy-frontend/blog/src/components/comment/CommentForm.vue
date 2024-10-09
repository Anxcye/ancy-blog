<template>
  <div>
    <a-form :model="commentForm" class="flex flex-col gap-2" ref="formRef" :rules="rules">
      <a-tabs class="border border-gray-bg rounded-xl p-2">
        <a-tab-pane key="1" tab="编辑">
          <div class="flex flex-col gap-2 border border-gray-bg rounded-md">
            <a-form-item name="content">
              <a-textarea
                v-model:value="commentForm.content"
                placeholder="说点什么吧"
                :bordered="false"
                :rows="4"
                name="content"
                :maxlength="1000"
              />
            </a-form-item>
            <div class="text-xs text-gray px-2 pb-2">
              支持
              <a href="https://guides.github.com/features/mastering-markdown/" target="_blank">
                Markdown
              </a>
            </div>
          </div>
        </a-tab-pane>
        <a-tab-pane key="2" tab="预览" force-render>
          <MdViewer :content="commentForm.content" />
        </a-tab-pane>
      </a-tabs>
      <div class="flex flex-col gap-2 w-full md:flex-row">
        <div
          class="mx-auto bg-gray-bg rounded-md p-1 flex-1 w-full text-center"
          v-if="adminStore.token"
        >
          <span>已管理员登录，将以管理员身份提交，若失败请检查权限</span>
        </div>
        <div class="flex flex-col gap-2 w-full md:flex-row" v-else>
          <a-form-item name="nickname" class="flex-1">
            <a-input v-model:value="commentForm.nickname" placeholder="昵称 *" />
          </a-form-item>
          <a-form-item name="email" class="flex-1">
            <a-input v-model:value="commentForm.email" placeholder="邮箱 *" />
          </a-form-item>
          <a-form-item name="avatar" class="flex-1">
            <a-input v-model:value="commentForm.avatar" placeholder="头像 url" />
          </a-form-item>
        </div>

        <a-button type="primary" @click="handleSubmit">Biu~</a-button>
      </div>
    </a-form>
  </div>
</template>

<script setup lang="ts">
import { reqCommentAdd } from '@/api/comment'
import type { CommentAddParams } from '@/api/comment/type'
import { message, type FormInstance } from 'ant-design-vue'
import { ref, watch } from 'vue'
import { useVisitorStore } from '@/stores/visitor'
import { useAdminStore } from '@/stores/admin'

const adminStore = useAdminStore()
const visitorStore = useVisitorStore()
const props = defineProps<{
  id: number
  postType: 'Article' | 'Note'
  parentId?: number
  toCommentId?: number
  toCommentNickname?: string
}>()
const emits = defineEmits(['submit'])

const commentForm = ref<CommentAddParams>({
  nickname: visitorStore.getUserInfo().nickname,
  email: visitorStore.getUserInfo().email,
  avatar: visitorStore.getUserInfo().avatar,
})

const formRef = ref<FormInstance>()

const rules = {
  content: [{ required: true, message: '写点什么吧' }],
  nickname: [{ required: true, message: '起个昵称吧' }],
  email: [{ required: true, message: '邮箱必须写哦' }],
}

watch(
  () => commentForm.value,
  () => {
    visitorStore.setUserInfo({
      nickname: commentForm.value.nickname,
      email: commentForm.value.email,
      avatar: commentForm.value.avatar,
    })
  },
  { deep: true },
)

const handleSubmit = async () => {
  if (!adminStore.token) {
    await formRef.value?.validate()
  } else {
    commentForm.value.nickname = ''
    commentForm.value.email = ''
    commentForm.value.avatar = ''
  }
  const res = await reqCommentAdd({
    ...commentForm.value,
    parentId: props.parentId,
    articleId: props.id,
    type: props.postType === 'Article' ? '0' : '1',
    toCommentId: props.toCommentId,
    toCommentNickname: props.toCommentNickname,
  })
  message.success('提交成功')
  emits('submit', {
    ...commentForm.value,
    parentId: props.parentId,
    id: res.data,
    createTime: new Date().toISOString(),
    userId: adminStore.token ? 1 : -1,
    toCommentNickname: props.toCommentNickname,
    toCommentId: props.toCommentId,
  })
}
</script>

<style scoped lang="scss"></style>
