<!-- File: app/components/CommentList.vue
     Purpose: Fetch and render a threaded vertical list of comments for an article.
     Module: components, presentation layer. -->
<template>
  <div class="comment-section">
    <h3 class="section-title">
      {{ t('article.commentsCount', { n: totalComments }) }}
    </h3>

    <CommentForm
      class="root-form"
      :content-type="contentType"
      :content-id="contentId"
      :require-approval="requireApproval"
      @success="handleRootSuccess"
    />

    <div v-if="pending" class="comments-list skeleton-list" aria-hidden="true">
      <article v-for="card in skeletonCards" :key="card" class="thread-card skeleton-card">
        <div class="skeleton-comment-head">
          <span class="skeleton-avatar"></span>
          <div class="skeleton-lines">
            <span class="skeleton-line short"></span>
            <span class="skeleton-line tiny"></span>
          </div>
        </div>
        <div class="skeleton-body">
          <span class="skeleton-line wide"></span>
          <span class="skeleton-line wide"></span>
          <span class="skeleton-line mid"></span>
        </div>
        <div class="skeleton-replies">
          <span class="skeleton-line short"></span>
          <span class="skeleton-line mid"></span>
        </div>
      </article>
    </div>

    <div v-else-if="commentThreads.length" class="comments-list">
      <article
        v-for="comment in commentThreads"
        :key="comment.id"
        class="thread-card"
      >
        <CommentItem
          :comment="comment"
          :content-type="contentType"
          :content-id="contentId"
          :is-replying="replyingToId === comment.id"
          :replying-to-id="replyingToId"
          :require-approval="requireApproval"
          @reply="setReplyTarget"
          @cancelReply="clearReplyTarget"
          @replySuccess="handleReplySuccess"
        />
      </article>
    </div>

    <div v-else class="empty-state">
      {{ t('article.noComments') }}
    </div>

    <InfiniteScrollTrigger
      v-if="commentThreads.length > 0"
      :loading="loadingMore"
      :done="!hasMore"
      :done-text="t('comment.noMore')"
      @load="loadMore"
    >
      <template #loading>
        <div class="load-more-state">
          <div class="spinner"></div>
          <span>{{ t('comment.loadingMore') }}</span>
        </div>
      </template>
    </InfiniteScrollTrigger>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { useApi, type CommentContentType, type CommentThread } from '~/composables/useApi'
import CommentItem from './CommentItem.vue'
import CommentForm from './CommentForm.vue'
import InfiniteScrollTrigger from './InfiniteScrollTrigger.vue'

const props = defineProps<{
  contentType: CommentContentType
  contentId: string
  requireApproval?: boolean
}>()

const emit = defineEmits<{
  (e: 'countChange', total: number): void
}>()

const { t } = useI18n()
const { listComments, getCommentTotal } = useApi()

const page = ref(1)
const commentThreads = ref<CommentThread[]>([])
const totalComments = ref(0)
const threadTotal = ref(0)
const loadingMore = ref(false)
const pending = ref(true)
const replyingToId = ref<string | null>(null)
const skeletonCards = [1, 2, 3]

const hasMore = computed(() => commentThreads.value.length < threadTotal.value)

function setReplyTarget(target: CommentThread) {
  replyingToId.value = target.id
}

function clearReplyTarget() {
  replyingToId.value = null
}

async function fetchTotal() {
  try {
    totalComments.value = await getCommentTotal(props.contentType, props.contentId)
    emit('countChange', totalComments.value)
  } catch (err) {
    console.error('Failed to count comments', err)
  }
}

async function fetchComments(nextPage: number, append = false) {
  try {
    const res = await listComments(props.contentType, props.contentId, { page: nextPage, pageSize: 12 })
    const rows = res.rows || []
    commentThreads.value = append ? dedupeThreads([...commentThreads.value, ...rows]) : rows
    threadTotal.value = res.total || 0
  } catch (err) {
    console.error('Failed to list comments', err)
    if (append) {
      page.value = Math.max(1, page.value - 1)
    }
  } finally {
    pending.value = false
    loadingMore.value = false
  }
}

onMounted(async () => {
  pending.value = true
  await Promise.all([
    fetchTotal(),
    fetchComments(1),
  ])
})

async function loadMore() {
  if (loadingMore.value || !hasMore.value) return
  loadingMore.value = true
  page.value += 1
  await fetchComments(page.value, true)
}

function handleRootSuccess() {
  page.value = 1
  pending.value = true
  Promise.all([fetchTotal(), fetchComments(1)])
}

function handleReplySuccess() {
  clearReplyTarget()
  page.value = 1
  pending.value = true
  Promise.all([fetchTotal(), fetchComments(1)])
}

function dedupeThreads(items: CommentThread[]) {
  const seen = new Map<string, CommentThread>()
  for (const item of items) {
    seen.set(item.id, item)
  }
  return Array.from(seen.values())
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
  gap: 18px;
}

.thread-card {
  padding: 18px;
  border-radius: var(--radius-lg);
}

.skeleton-card {
  position: relative;
  overflow: hidden;
}

.skeleton-card::after {
  content: '';
  position: absolute;
  inset: 0;
  background: linear-gradient(90deg, transparent, color-mix(in srgb, #fff 34%, transparent), transparent);
  transform: translateX(-100%);
  animation: skeleton-sheen 1.35s ease-in-out infinite;
}

.skeleton-comment-head {
  display: flex;
  align-items: center;
  gap: 14px;
  margin-bottom: 18px;
}

.skeleton-avatar {
  width: 42px;
  height: 42px;
  border-radius: 50%;
  background: color-mix(in srgb, var(--text) 8%, var(--bg-secondary));
  flex-shrink: 0;
}

.skeleton-lines {
  display: grid;
  gap: 8px;
  width: 100%;
}

.skeleton-body,
.skeleton-replies {
  display: grid;
  gap: 10px;
}

.skeleton-replies {
  margin-top: 18px;
  padding-left: 14px;
  border-left: 2px solid color-mix(in srgb, var(--accent) 16%, var(--border));
}

.skeleton-line {
  display: block;
  height: 12px;
  border-radius: 999px;
  background: color-mix(in srgb, var(--text) 8%, var(--bg-secondary));
}

.skeleton-line.tiny {
  width: 26%;
}

.skeleton-line.short {
  width: 34%;
}

.skeleton-line.mid {
  width: 62%;
}

.skeleton-line.wide {
  width: 100%;
}

.load-more-state,
.empty-state {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 10px;
  color: var(--text-subtle);
}

.empty-state {
  width: 100%;
  padding: 48px 0;
}

.spinner {
  width: 22px;
  height: 22px;
  border: 2px solid var(--border);
  border-top-color: var(--accent);
  border-radius: 50%;
  animation: comment-spin 0.8s linear infinite;
}

@keyframes comment-spin {
  to {
    transform: rotate(360deg);
  }
}

@keyframes skeleton-sheen {
  to {
    transform: translateX(100%);
  }
}
</style>
