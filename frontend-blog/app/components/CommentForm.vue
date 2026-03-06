<!-- File: app/components/CommentForm.vue
     Purpose: A reusable form for submitting root comments or replying to a specific comment.
     Module: components, presentation layer. -->
<template>
  <div class="comment-form-wrap" :class="{ 'is-reply': !!replyTo }">
    <div v-if="replyTo" class="reply-header">
      <span>正在回复 <strong>{{ replyTo.nickname }}</strong></span>
      <button class="cancel-reply-btn" @click="$emit('cancel', replyTo.id)">取消回复</button>
    </div>

    <form class="comment-form" @submit.prevent="handleSubmit">
      <div class="form-row meta-row">
        <label class="input-wrap">
          <input v-model="form.nickname" type="text" placeholder="昵称 *" required :disabled="submitting" />
        </label>
        <label class="input-wrap">
          <input v-model="form.email" type="email" placeholder="邮箱" :disabled="submitting" />
        </label>
        <label class="input-wrap">
          <input v-model="form.website" type="url" placeholder="网址" :disabled="submitting" />
        </label>
      </div>

      <div class="form-row content-row">
        <textarea
          v-model="form.content"
          placeholder="说点什么吧... 支持 Markdown 语法"
          required
          rows="4"
          :disabled="submitting"
        ></textarea>
      </div>

      <div class="form-actions">
        <span class="form-hint">提交后需等待审核后才会显示</span>
        <button type="submit" class="submit-btn" :disabled="submitting || !isValid">
          <template v-if="submitting">
            <span class="spinner"></span> 提交中...
          </template>
          <template v-else>发表评论</template>
        </button>
      </div>
    </form>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { useApi, type Comment } from '~/composables/useApi'
import { useStorage } from '@vueuse/core'

const props = defineProps<{
  articleId: string
  replyTo?: Comment
}>()

const emit = defineEmits<{
  (e: 'success'): void
  (e: 'cancel', commentId: string): void
}>()

const { createComment } = useApi()
const toast = useToast()

const submitting = ref(false)

// Persist user info in local storage (like typical blog comments)
const storedReviewer = useStorage('blog-commenter-info', {
  nickname: '',
  email: '',
  website: ''
})

const form = reactive({
  nickname: '',
  email: '',
  website: '',
  content: '',
})

// Hydrate form from storage on mount
onMounted(() => {
  form.nickname = storedReviewer.value.nickname
  form.email = storedReviewer.value.email
  form.website = storedReviewer.value.website
})

const isValid = computed(() => {
  return form.nickname.trim().length > 0 && form.content.trim().length > 0
})

async function handleSubmit() {
  if (!isValid.value || submitting.value) return

  // Save info for future
  storedReviewer.value = {
    nickname: form.nickname.trim(),
    email: form.email.trim(),
    website: form.website.trim()
  }

  submitting.value = true
  try {
    let parentId = undefined;
    if (props.replyTo) {
         // In ancy-blog backend, parent_id is the root comment ID, to_comment_id is the direct replied ID
         parentId = props.replyTo.parentId || props.replyTo.id
    }
    await createComment({
      articleId: props.articleId,
      parentId: parentId, 
      toCommentId: props.replyTo ? props.replyTo.id : undefined,
      content: form.content.trim(),
      nickname: form.nickname.trim(),
      email: form.email.trim() || undefined,
      website: form.website.trim() || undefined,
    })

    toast.add({
      title: '评论已成功提交！',
      description: '您的评论可能需要等待管理员审核。',
      color: 'green'
    })

    form.content = ''
    emit('success')
  } catch (err: any) {
    if (!err.fatal) {
       toast.add({
          title: '提交失败',
          description: err.message || '服务异常',
          color: 'red'
       })
    }
    console.error(err)
  } finally {
    submitting.value = false
  }
}
</script>

<style scoped>
.comment-form-wrap {
  background: var(--surface);
  border: 1px solid var(--border);
  border-radius: var(--radius-lg);
  padding: 20px;
  box-shadow: var(--shadow-sm);
  transition: all var(--dur-base);
}

.comment-form-wrap.is-reply {
  background: var(--bg-secondary);
  border-color: var(--accent-soft);
  box-shadow: none;
}

.reply-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 13px;
  color: var(--text-muted);
  margin-bottom: 16px;
  padding-bottom: 12px;
  border-bottom: 1px solid var(--border);
}

.reply-header strong {
  color: var(--text);
  font-weight: 600;
}

.cancel-reply-btn {
  color: var(--accent);
  background: none;
  border: none;
  font-size: 13px;
  cursor: pointer;
}

.cancel-reply-btn:hover {
  text-decoration: underline;
}

.form-row {
  display: flex;
  gap: 16px;
  margin-bottom: 16px;
}

@media (max-width: 640px) {
  .meta-row {
    flex-direction: column;
    gap: 12px;
  }
}

.input-wrap {
  flex: 1;
}

input, textarea {
  width: 100%;
  padding: 12px 16px;
  background: var(--bg-primary);
  border: 1px solid var(--border);
  border-radius: var(--radius-md);
  color: var(--text);
  font-family: inherit;
  font-size: 14px;
  transition: border-color var(--dur-fast), box-shadow var(--dur-fast);
}

textarea {
  resize: vertical;
  min-height: 100px;
}

input:focus, textarea:focus {
  outline: none;
  border-color: var(--accent);
  box-shadow: 0 0 0 2px var(--accent-soft);
}

input:disabled, textarea:disabled {
  opacity: 0.7;
  cursor: not-allowed;
  background: var(--bg-secondary);
}

.form-actions {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.form-hint {
  font-size: 12px;
  color: var(--text-subtle);
}

.submit-btn {
  background: var(--accent);
  color: #fff;
  border: none;
  padding: 10px 24px;
  border-radius: var(--radius-md);
  font-size: 14px;
  font-weight: 600;
  cursor: pointer;
  transition: opacity var(--dur-fast), transform var(--dur-fast);
  display: inline-flex;
  align-items: center;
  gap: 8px;
}

.submit-btn:hover:not(:disabled) {
  opacity: 0.9;
  transform: translateY(-1px);
}

.submit-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.spinner {
  width: 14px;
  height: 14px;
  border: 2px solid rgba(255,255,255,0.3);
  border-top-color: #fff;
  border-radius: 50%;
  animation: spinner 0.6s linear infinite;
}

@keyframes spinner {
  to { transform: rotate(360deg); }
}
</style>
