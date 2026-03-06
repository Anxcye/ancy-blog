<!-- File: app/components/CommentItem.vue
     Purpose: Single comment rendering + recursive children.
     Module: components, presentation layer. -->
<template>
  <div class="comment-item" :id="`comment-${comment.id}`">
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

        <span class="comment-badge" v-if="comment.isPinned">置顶</span>
        <span class="comment-badge admin-badge" v-if="isAdmin">博主</span>

        <span class="comment-meta">
          <time :datetime="comment.createdAt">{{ formatDate(comment.createdAt) }}</time>
        </span>
      </div>

      <div class="comment-content">
        <span class="reply-target" v-if="isSubReply && replyToNickname">@{{ replyToNickname }}</span>
        <!-- Content should ideally be markdown rendered, but for now we just show it. 
             Security note: ensure backend escapes malicious scripts or render safely. -->
        <p>{{ comment.content }}</p>
      </div>

      <div class="comment-actions">
        <button class="action-btn" @click="$emit('reply', comment)">
          <svg viewBox="0 0 24 24" fill="none" class="icon"><path d="M3 10h10a5 5 0 0 1 5 5v2a5 5 0 0 1-5 5H3m0-12l4-4m-4 4l4 4" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"></path></svg>
          回复
        </button>
      </div>
      
      <!-- Sub-comments Recursive Render -->
      <!-- Wait for CommentList children rendering, or handle it via a flat sorted list in CommentList -->
      <slot name="children"></slot>
      
      <div v-if="isReplying" class="reply-box-wrapper">
         <CommentForm 
            :article-id="articleId" 
            :reply-to="comment" 
            @cancel="$emit('cancelReply', comment.id)" 
            @success="$emit('replySuccess')"
         />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { Comment } from '~/composables/useApi'

const props = defineProps<{
  comment: Comment
  articleId: string
  isReplying?: boolean
  isSubReply?: boolean
  replyToNickname?: string
}>()

const emit = defineEmits<{
  (e: 'reply', parent: Comment): void
  (e: 'cancelReply', id: string): void
  (e: 'replySuccess'): void
}>()

// Use standard datetime formatting
function formatDate(dateStr: string) {
  const d = new Date(dateStr)
  return new Intl.DateTimeFormat('zh-CN', {
    month: '2-digit', day: '2-digit', hour: '2-digit', minute: '2-digit'
  }).format(d)
}

// Temporary heuristic for admin badge, could be adjusted
const isAdmin = computed(() => {
   return props.comment.nickname.toLowerCase() === 'admin'
})
</script>

<style scoped>
.comment-item {
  display: flex;
  gap: 16px;
  margin-top: 24px;
}

.comment-avatar {
  width: 40px;
  height: 40px;
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
  border-radius: 4px;
  background: var(--surface-hover);
  color: var(--text-muted);
  font-weight: 500;
  line-height: 1;
}

.admin-badge {
  background: var(--accent-soft);
  color: var(--accent-text);
  border: 1px solid var(--accent);
}

.comment-content {
  font-size: 14px;
  line-height: 1.6;
  color: var(--text);
  margin-bottom: 8px;
  word-break: break-word;
}

.comment-content p {
  margin: 0 0 8px;
}

.comment-content p:last-child {
  margin: 0;
}

.reply-target {
  color: var(--accent);
  font-weight: 500;
  margin-right: 6px;
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
  animation: slide-down 0.2s var(--ease-out);
}

@keyframes slide-down {
  from { opacity: 0; transform: translateY(-10px); }
  to { opacity: 1; transform: translateY(0); }
}
</style>
