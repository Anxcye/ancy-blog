<!-- File: app/components/CommentList.vue
     Purpose: Fetch and render a nested list of comments for an article.
     Module: components, presentation layer. -->
<template>
  <div class="comment-section">
    <h3 class="section-title">
      共 {{ total }} 条评论
    </h3>

    <!-- Root Submission Form -->
    <CommentForm
      class="root-form"
      :article-id="articleId"
      @success="handleRootSuccess"
    />

    <!-- Comments Tree -->
    <div class="comments-list" v-if="comments.length">
      <!-- Only render root comments (no parentId) at top level -->
      <template v-for="comment in rootComments" :key="comment.id">
        <CommentItem
          :comment="comment"
          :article-id="articleId"
          :is-replying="replyingToId === comment.id"
          @reply="setReplyTarget"
          @cancelReply="clearReplyTarget"
          @replySuccess="handleReplySuccess"
        >
          <template #children>
            <div class="sub-comments" v-if="getChildren(comment.id).length">
              <CommentItem
                v-for="sub in getChildren(comment.id)"
                :key="sub.id"
                :comment="sub"
                :article-id="articleId"
                :is-replying="replyingToId === sub.id"
                :is-sub-reply="true"
                :reply-to-nickname="getCommenterName(sub.toCommentId)"
                @reply="setReplyTarget"
                @cancelReply="clearReplyTarget"
                @replySuccess="handleReplySuccess"
              />
            </div>
          </template>
        </CommentItem>
      </template>
    </div>

    <div v-if="pending" class="loading-state">
      <div class="spinner"></div> 加载中...
    </div>

    <div v-else-if="!comments.length" class="empty-state">
      —— 快来抢沙发吧 ——
    </div>

    <!-- Infinite Scroll Trigger (optional, assuming we load pages of root comments) -->
    <InfiniteScrollTrigger
      v-if="comments.length > 0 && hasMore"
      :loading="loadingMore"
      :done="!hasMore"
      done-text="没有更多评论了"
      @load="loadMore"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { useApi, type Comment } from '~/composables/useApi'
import CommentItem from './CommentItem.vue'
import CommentForm from './CommentForm.vue'
import InfiniteScrollTrigger from './InfiniteScrollTrigger.vue'

const props = defineProps<{
  articleId: string
}>()

const { listComments, getCommentTotal } = useApi()

const page = ref(1)
const comments = ref<Comment[]>([])
const total = ref(0)
const loadingMore = ref(false)
const pending = ref(true)
const hasMore = computed(() => comments.value.length < total.value)

// Reply state
const replyingToId = ref<string | null>(null)

function setReplyTarget(target: Comment) {
  replyingToId.value = target.id
}

function clearReplyTarget() {
  replyingToId.value = null
}

async function fetchTotal() {
  try {
    const res = await getCommentTotal(props.articleId)
    total.value = res?.total || 0
  } catch (err) {
    console.error('Failed to count comments', err)
  }
}

async function fetchComments(p: number, append = false) {
  try {
    const res = await listComments(props.articleId, { page: p, pageSize: 20 })
    if (res.rows) {
      if (append) {
        comments.value.push(...res.rows)
      } else {
        comments.value = res.rows
      }
    }
    // Update total just in case
    total.value = res.total || 0
  } catch (err) {
    console.error('Failed to list comments', err)
  } finally {
    pending.value = false
    loadingMore.value = false
  }
}

// Initial fetch
onMounted(async () => {
    pending.value = true
    await Promise.all([
        fetchTotal(),
        fetchComments(1)
    ])
})

async function loadMore() {
  if (loadingMore.value || !hasMore.value) return
  loadingMore.value = true
  page.value++
  await fetchComments(page.value, true)
}

function handleRootSuccess() {
  // New comment might be queued for review, or immediately approved
  // Refresh to show if immediately visible
  page.value = 1
  pending.value = true
  Promise.all([fetchTotal(), fetchComments(1)])
}

function handleReplySuccess() {
  clearReplyTarget()
  // Same deal, refresh current fetched comments
  page.value = 1
  pending.value = true
  Promise.all([fetchTotal(), fetchComments(1)])
}

// Data sorting / nesting logic
const rootComments = computed(() => {
  return comments.value.filter(c => !c.parentId)
})

function getChildren(parentId: string) {
  // Return all comments whose parentId matches the root comment's id
  // Note: in ancy-blog, multi-level replies point to the SAME parentId to flatten the tree visually
  return comments.value.filter(c => c.parentId === parentId).sort((a,b) => new Date(a.createdAt).getTime() - new Date(b.createdAt).getTime())
}

function getCommenterName(id?: string) {
  if (!id) return ''
  const t = comments.value.find(c => c.id === id)
  return t ? t.nickname : ''
}
</script>

<style scoped>
.comment-section {
  margin-top: 64px;
}

.section-title {
  font-size: 1.25rem;
  font-weight: 700;
  margin-bottom: 24px;
  display: flex;
  align-items: center;
  gap: 12px;
}

.section-title::before {
  content: '';
  display: block;
  width: 4px;
  height: 18px;
  background: var(--accent);
  border-radius: 2px;
}

.root-form {
  margin-bottom: 40px;
}

.comments-list {
  display: flex;
  flex-direction: column;
}

.sub-comments {
  margin-top: 16px;
  padding: 16px 0 16px 20px;
  border-left: 2px solid var(--border);
  background: var(--bg-secondary);
  border-radius: 0 var(--radius-md) var(--radius-md) 0;
}

.sub-comments .comment-item {
  margin-top: 16px;
}

.sub-comments .comment-item:first-child {
  margin-top: 0;
}

.loading-state, .empty-state {
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 48px 0;
  color: var(--text-subtle);
  font-size: 14px;
  gap: 12px;
}

.spinner {
  width: 18px;
  height: 18px;
  border: 2px solid var(--border);
  border-top-color: var(--accent);
  border-radius: 50%;
  animation: s-spin 0.6s linear infinite;
}

@keyframes s-spin {
  to { transform: rotate(360deg); }
}

@media (max-width: 640px) {
  .sub-comments {
    padding-left: 12px;
    margin-left: -20px; /* pull it out a bit on mobile to save space */
  }
}
</style>
