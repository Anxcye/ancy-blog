<!-- File: app/components/CommentItem.vue
     Purpose: Render a single threaded comment node with recursive replies and reply actions.
     Module: components, presentation layer. -->
<template>
  <div class="comment-item" :class="{ 'is-root': depth === 0, 'is-child': depth > 0 }" :id="`comment-${comment.id}`">
    <div class="comment-avatar">
      <img v-if="comment.avatarUrl" :src="comment.avatarUrl" :alt="comment.nickname" />
      <span v-else>{{ comment.nickname.charAt(0).toUpperCase() }}</span>
    </div>

    <div class="comment-main">
      <div class="comment-header">
        <a v-if="comment.website" :href="comment.website" target="_blank" rel="nofollow noopener" class="comment-author">
          {{ comment.nickname }}
        </a>
        <span v-else class="comment-author">{{ comment.nickname }}</span>

        <span class="comment-badge" v-if="comment.isPinned">{{ t('comment.pinned') }}</span>
        <span class="comment-badge admin-badge" v-if="comment.isAuthor">{{ t('comment.authorBadge') }}</span>

        <span class="comment-meta">
          <time :datetime="comment.createdAt">{{ formatDate(comment.createdAt) }}</time>
        </span>
      </div>

      <div class="comment-content markdown-body">
        <span class="reply-target" v-if="comment.toCommentNickname">@{{ comment.toCommentNickname }}</span>
        <div class="comment-markdown" v-html="renderedContent"></div>
      </div>

      <div class="comment-actions">
        <button class="action-btn" @click="$emit('reply', comment)">
          <svg viewBox="0 0 24 24" fill="none" class="icon"><path d="M3 10h10a5 5 0 0 1 5 5v2a5 5 0 0 1-5 5H3m0-12l4-4m-4 4l4 4" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"></path></svg>
          {{ t('comment.reply') }}
        </button>
      </div>

      <div v-if="isReplying" class="reply-box-wrapper">
        <CommentForm
          :content-type="contentType"
          :content-id="contentId"
          :reply-to="comment"
          :require-approval="requireApproval"
          @cancel="$emit('cancelReply', comment.id)"
          @success="$emit('replySuccess')"
        />
      </div>

      <div v-if="comment.children?.length" class="comment-children">
        <CommentItem
          v-for="child in comment.children"
          :key="child.id"
          :comment="child"
          :content-type="contentType"
          :content-id="contentId"
          :is-replying="replyingToId === child.id"
          :replying-to-id="replyingToId"
          :depth="depth + 1"
          :require-approval="requireApproval"
          @reply="$emit('reply', $event)"
          @cancelReply="$emit('cancelReply', $event)"
          @replySuccess="$emit('replySuccess')"
        />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import type { CommentContentType, CommentThread } from '~/composables/useApi'
import { renderCommentMarkdown } from '~/utils/commentMarkdown'

const props = withDefaults(defineProps<{
  comment: CommentThread
  contentType: CommentContentType
  contentId: string
  isReplying?: boolean
  replyingToId?: string | null
  depth?: number
  requireApproval?: boolean
}>(), {
  depth: 0,
  replyingToId: null,
})

const emit = defineEmits<{
  (e: 'reply', parent: CommentThread): void
  (e: 'cancelReply', id: string): void
  (e: 'replySuccess'): void
}>()

const { t, locale } = useI18n()

function formatDate(dateStr: string) {
  const date = new Date(dateStr)
  return new Intl.DateTimeFormat(locale.value === 'en' ? 'en-US' : 'zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
  }).format(date)
}

const renderedContent = computed(() => renderCommentMarkdown(props.comment.content))
</script>

<style scoped>
.comment-item {
  display: flex;
  gap: 14px;
}

.comment-item.is-root {
  margin-top: 0;
}

.comment-item.is-child {
  margin-top: 16px;
  padding-top: 16px;
  border-top: 1px solid color-mix(in srgb, var(--border) 72%, transparent);
}

.comment-avatar {
  width: 42px;
  height: 42px;
  border-radius: 50%;
  background: var(--surface-hover);
  color: var(--text-muted);
  font-weight: 700;
  display: grid;
  place-items: center;
  overflow: hidden;
  flex-shrink: 0;
  border: 1px solid var(--border);
}

.comment-avatar img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.comment-main {
  flex: 1;
  min-width: 0;
}

.comment-header {
  display: flex;
  align-items: center;
  flex-wrap: wrap;
  gap: 8px;
  margin-bottom: 6px;
}

.comment-author {
  font-weight: 600;
  font-size: 14px;
  color: var(--text);
  text-decoration: none;
  transition: color var(--dur-fast);
}

a.comment-author:hover {
  color: var(--accent);
}

.comment-meta {
  font-size: 12px;
  color: var(--text-subtle);
  margin-left: auto;
}

.comment-badge {
  font-size: 11px;
  padding: 2px 6px;
  border-radius: 999px;
  background: var(--surface-hover);
  color: var(--text-muted);
  font-weight: 600;
  line-height: 1;
}

.admin-badge {
  background: var(--accent-soft);
  color: var(--accent-text);
  border: 1px solid color-mix(in srgb, var(--accent) 50%, transparent);
}

.comment-content {
  font-size: 14px;
  line-height: 1.7;
  color: var(--text);
  margin-bottom: 10px;
  word-break: break-word;
}

.reply-target {
  color: var(--accent);
  font-weight: 600;
  margin-right: 6px;
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

.comment-actions {
  display: flex;
  gap: 16px;
  align-items: center;
}

.action-btn {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  background: none;
  border: none;
  font-size: 12px;
  color: var(--text-subtle);
  cursor: pointer;
  padding: 0;
  transition: color var(--dur-fast);
}

.action-btn:hover {
  color: var(--accent);
}

.icon {
  width: 14px;
  height: 14px;
}

.reply-box-wrapper {
  margin-top: 16px;
}

.comment-children {
  margin-top: 16px;
  padding-left: 16px;
  border-left: 2px solid color-mix(in srgb, var(--accent) 16%, var(--border));
}

@media (max-width: 640px) {
  .comment-item {
    gap: 12px;
  }

  .comment-avatar {
    width: 36px;
    height: 36px;
  }

  .comment-children {
    padding-left: 12px;
  }
}
</style>
