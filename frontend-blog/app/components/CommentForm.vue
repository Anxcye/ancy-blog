<!-- File: app/components/CommentForm.vue
     Purpose: A reusable form for submitting root comments or replying to a specific comment.
     Module: components, presentation layer. -->
<template>
  <div class="comment-form-wrap" :class="{ 'is-reply': !!replyTo }">
    <button
      v-if="!replyTo && collapsed"
      type="button"
      class="composer-trigger"
      @click="openComposer"
    >
      <span class="composer-trigger-title">{{ t('comment.openComposer') }}</span>
      <span class="composer-trigger-hint">{{ t('comment.openComposerHint') }}</span>
    </button>

    <div v-if="replyTo" class="reply-header">
      <span>{{ t('comment.replyingTo', { name: replyTo.nickname }) }}</span>
      <button class="cancel-reply-btn" @click="$emit('cancel', replyTo.id)">{{ t('comment.cancelReply') }}</button>
    </div>

    <form v-if="!collapsed || !!replyTo" class="comment-form" @submit.prevent="handleSubmit">
      <div class="form-row meta-row">
        <label class="input-wrap">
          <input v-model="form.nickname" type="text" :placeholder="`${t('comment.nickname')} *`" required :disabled="submitting" />
        </label>
        <label class="input-wrap">
          <input v-model="form.email" type="email" :placeholder="t('comment.email')" :disabled="submitting" />
        </label>
        <label class="input-wrap">
          <input v-model="form.website" type="url" :placeholder="t('comment.website')" :disabled="submitting" />
        </label>
      </div>

      <div class="form-row meta-row">
        <label class="input-wrap">
          <input v-model="form.avatarUrl" type="url" :placeholder="t('comment.avatarUrl')" :disabled="submitting" />
        </label>
      </div>

      <div class="editor-shell">
        <div class="editor-toolbar">
          <div class="editor-tabs">
            <button type="button" class="editor-tab" :class="{ active: editorMode === 'write' }" @click="editorMode = 'write'">
              {{ t('comment.write') }}
            </button>
            <button type="button" class="editor-tab" :class="{ active: editorMode === 'preview' }" @click="editorMode = 'preview'">
              {{ t('comment.preview') }}
            </button>
          </div>
        </div>

        <div v-if="editorMode === 'write'" class="content-row">
          <textarea
            v-model="form.content"
            :placeholder="t('comment.content')"
            required
            rows="4"
            :disabled="submitting"
          ></textarea>
        </div>

        <div v-else class="comment-preview markdown-body" :class="{ empty: !form.content.trim() }">
          <div v-if="form.content.trim()" v-html="previewHtml"></div>
          <p v-else>{{ t('comment.previewEmpty') }}</p>
        </div>

        <p class="editor-help">
          {{ t('comment.markdownHelpPrefix') }}
          <span> </span>
          <a href="https://docs.github.com/en/get-started/writing-on-github/getting-started-with-writing-and-formatting-on-github/basic-writing-and-formatting-syntax" target="_blank" rel="noopener noreferrer nofollow">
            {{ t('comment.markdownHelpLink') }}
          </a>
          {{ t('comment.markdownHelpSuffix') }}
        </p>
      </div>

      <div class="form-actions">
        <span class="form-hint">{{ reviewHint }}</span>
        <button type="submit" class="submit-btn" :disabled="submitting || !isValid">
          <template v-if="submitting">
            <span class="spinner"></span> {{ t('comment.submitting') }}
          </template>
          <template v-else>{{ t('article.leaveComment') }}</template>
        </button>
      </div>
    </form>

    <p v-if="successMessage" class="submit-feedback">
      {{ successMessage }}
    </p>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { useApi, type Comment } from '~/composables/useApi'
import { useStorage } from '@vueuse/core'
import { renderCommentMarkdown } from '~/utils/commentMarkdown'

const props = defineProps<{
  articleId: string
  replyTo?: Comment
  requireApproval?: boolean
}>()

const emit = defineEmits<{
  (e: 'success'): void
  (e: 'cancel', commentId: string): void
}>()

const { createComment } = useApi()
const { t } = useI18n()
const toast = useToast()

const submitting = ref(false)
const editorMode = ref<'write' | 'preview'>('write')
const collapsed = ref(!props.replyTo)
const successMessage = ref('')

// Persist user info in local storage (like typical blog comments)
const storedReviewer = useStorage('blog-commenter-info', {
  nickname: '',
  email: '',
  website: '',
  avatarUrl: ''
})

const form = reactive({
  nickname: '',
  email: '',
  website: '',
  avatarUrl: '',
  content: '',
})

// Hydrate form from storage on mount
onMounted(() => {
  form.nickname = storedReviewer.value.nickname
  form.email = storedReviewer.value.email
  form.website = storedReviewer.value.website
  form.avatarUrl = storedReviewer.value.avatarUrl
})

watch(() => form.content, () => {
  successMessage.value = ''
})

const isValid = computed(() => {
  return form.nickname.trim().length > 0 && form.content.trim().length > 0
})

const reviewHint = computed(() => (
  props.requireApproval ? t('comment.reviewRequired') : t('comment.reviewNotRequired')
))

const previewHtml = computed(() => renderCommentMarkdown(form.content))

async function handleSubmit() {
  if (!isValid.value || submitting.value) return

  storedReviewer.value = {
    nickname: form.nickname.trim(),
    email: form.email.trim(),
    website: form.website.trim(),
    avatarUrl: form.avatarUrl.trim()
  }

  submitting.value = true
  try {
    let parentId = undefined
    let rootId = undefined
    if (props.replyTo) {
      parentId = props.replyTo.id
      rootId = props.replyTo.rootId || props.replyTo.id
    }
    await createComment({
      articleId: props.articleId,
      parentId,
      rootId,
      toCommentId: props.replyTo ? props.replyTo.id : undefined,
      content: form.content.trim(),
      nickname: form.nickname.trim(),
      email: form.email.trim() || undefined,
      website: form.website.trim() || undefined,
      avatarUrl: form.avatarUrl.trim() || undefined,
    })

    toast.add({
      title: t('comment.submitSuccessTitle'),
      description: reviewHint.value,
      color: 'green'
    })

    form.content = ''
    editorMode.value = 'write'
    successMessage.value = reviewHint.value
    if (!props.replyTo) {
      collapsed.value = true
    }
    emit('success')
  } catch (err: any) {
    if (!err.fatal) {
      toast.add({
        title: t('comment.submitFailedTitle'),
        description: err.message || '服务异常',
        color: 'red'
      })
    }
    console.error(err)
  } finally {
    submitting.value = false
  }
}

function openComposer() {
  successMessage.value = ''
  collapsed.value = false
}
</script>

<style scoped>
.comment-form-wrap {
  border-radius: var(--radius-lg);
  padding: 0;
  transition: all var(--dur-base);
}

.comment-form {
  animation: composer-pop 0.42s cubic-bezier(0.22, 1.18, 0.36, 1);
  transform-origin: top center;
}

.composer-trigger {
  width: 100%;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 10px;
  margin: 0 auto;
  padding: 22px 18px;
  border: 1px solid var(--border);
  border-radius: var(--radius-md);
  background: color-mix(in srgb, var(--bg-secondary) 42%, transparent);
  color: var(--text);
  text-align: center;
  cursor: pointer;
  box-shadow: 0 8px 20px color-mix(in srgb, var(--accent) 7%, transparent);
  transition:
    border-color 180ms ease,
    transform 340ms cubic-bezier(0.22, 1.18, 0.36, 1),
    box-shadow 340ms cubic-bezier(0.22, 1.18, 0.36, 1),
    background 180ms ease;
}

.composer-trigger:hover {
  border-color: color-mix(in srgb, var(--accent) 40%, var(--border));
  background: color-mix(in srgb, var(--accent) 4%, var(--bg-secondary));
  transform: translateY(-4px) scale(1.01);
  box-shadow: 0 18px 36px color-mix(in srgb, var(--accent) 14%, transparent);
}

.composer-trigger:active {
  transform: translateY(-1px) scale(0.985);
  box-shadow: 0 10px 22px color-mix(in srgb, var(--accent) 10%, transparent);
}

.composer-trigger-title {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  min-width: 108px;
  padding: 10px 20px;
  border-radius: 999px;
  background: var(--accent);
  color: #fff;
  font-size: 15px;
  font-weight: 600;
  box-shadow: 0 10px 22px color-mix(in srgb, var(--accent) 24%, transparent);
  transition:
    transform 340ms cubic-bezier(0.22, 1.18, 0.36, 1),
    box-shadow 340ms cubic-bezier(0.22, 1.18, 0.36, 1);
}

.composer-trigger:hover .composer-trigger-title {
  transform: translateY(-1px) scale(1.035);
  box-shadow: 0 16px 28px color-mix(in srgb, var(--accent) 30%, transparent);
}

.composer-trigger-hint {
  font-size: 12px;
  color: var(--text-subtle);
  transition: transform 340ms cubic-bezier(0.22, 1.18, 0.36, 1), opacity 180ms ease;
}

.composer-trigger:hover .composer-trigger-hint {
  transform: translateY(1px);
}

.comment-form-wrap.is-reply {
  box-shadow: none;
}

.submit-feedback {
  margin-top: 14px;
  font-size: 13px;
  line-height: 1.6;
  color: var(--accent-text);
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

.input-wrap {
  flex: 1;
}

input,
textarea {
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

input:focus,
textarea:focus {
  outline: none;
  border-color: var(--accent);
  box-shadow: 0 0 0 2px var(--accent-soft);
}

input:disabled,
textarea:disabled {
  opacity: 0.7;
  cursor: not-allowed;
  background: var(--bg-secondary);
}

.editor-shell {
  margin-bottom: 16px;
  border: 1px solid var(--border);
  border-radius: var(--radius-md);
  overflow: hidden;
  background: color-mix(in srgb, var(--bg-secondary) 60%, transparent);
  animation: editor-shell-rise 0.5s cubic-bezier(0.22, 1.18, 0.36, 1);
}

.editor-toolbar {
  display: flex;
  align-items: center;
  padding: 8px 10px 0;
  background: color-mix(in srgb, var(--bg-secondary) 82%, transparent);
}

.editor-tabs {
  display: inline-flex;
  gap: 8px;
}

.editor-tab {
  padding: 6px 12px;
  border: 1px solid transparent;
  border-bottom: none;
  border-radius: 10px 10px 0 0;
  background: transparent;
  color: var(--text-subtle);
  font-size: 12px;
  cursor: pointer;
  transition: all var(--dur-fast);
}

.editor-tab.active {
  border-color: var(--border);
  background: var(--bg-primary);
  color: var(--text);
}

.content-row {
  margin: 0;
}

textarea {
  resize: vertical;
  min-height: 120px;
  border: none;
  border-top: 1px solid var(--border);
  border-radius: 0;
  background: transparent;
}

textarea:focus {
  box-shadow: none;
}

.comment-preview {
  min-height: 152px;
  padding: 14px 16px;
  border-top: 1px solid var(--border);
  background: var(--bg-primary);
}

.comment-preview.empty {
  display: flex;
  align-items: center;
  color: var(--text-subtle);
}

.editor-help {
  margin: 0;
  padding: 0 14px 12px;
  font-size: 12px;
  color: var(--text-subtle);
}

.editor-help a {
  color: var(--accent-text);
  text-decoration: underline;
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

.markdown-body :deep(p),
.markdown-body :deep(ul),
.markdown-body :deep(ol),
.markdown-body :deep(blockquote),
.markdown-body :deep(pre) {
  margin: 0 0 10px;
}

.markdown-body :deep(p:last-child),
.markdown-body :deep(ul:last-child),
.markdown-body :deep(ol:last-child),
.markdown-body :deep(blockquote:last-child),
.markdown-body :deep(pre:last-child) {
  margin-bottom: 0;
}

.markdown-body :deep(ul),
.markdown-body :deep(ol) {
  padding-left: 20px;
}

.markdown-body :deep(blockquote) {
  padding-left: 12px;
  border-left: 2px solid var(--border);
  color: var(--text-muted);
}

.markdown-body :deep(pre) {
  overflow-x: auto;
  padding: 12px 14px;
  border-radius: var(--radius-md);
  background: var(--bg-secondary);
  border: 1px solid var(--border);
}

.markdown-body :deep(code) {
  font-family: 'Fira Code', monospace;
  font-size: 0.92em;
  padding: 0.15em 0.35em;
  border-radius: 4px;
  background: var(--bg-secondary);
}

.markdown-body :deep(pre code) {
  padding: 0;
  background: transparent;
}

.markdown-body :deep(a) {
  color: var(--accent-text);
  text-decoration: underline;
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
  transform: translateY(-2px) scale(1.01);
}

.submit-btn:active:not(:disabled) {
  transform: translateY(0) scale(0.985);
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

@keyframes composer-pop {
  0% {
    opacity: 0;
    transform: translateY(10px) scale(0.985);
  }
  60% {
    opacity: 1;
    transform: translateY(-4px) scale(1.01);
  }
  100% {
    opacity: 1;
    transform: translateY(0) scale(1);
  }
}

@keyframes editor-shell-rise {
  0% {
    opacity: 0;
    transform: translateY(12px) scale(0.985);
  }
  65% {
    opacity: 1;
    transform: translateY(-3px) scale(1.005);
  }
  100% {
    opacity: 1;
    transform: translateY(0) scale(1);
  }
}

@media (max-width: 640px) {
  .meta-row {
    flex-direction: column;
    gap: 12px;
  }

  .reply-header,
  .form-actions {
    flex-direction: column;
    align-items: flex-start;
    gap: 10px;
  }
}
</style>
