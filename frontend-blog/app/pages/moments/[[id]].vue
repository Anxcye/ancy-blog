<!-- File: app/pages/moments/[[id]].vue
     Purpose: Render the moments feed with a URL-driven detail modal and preserved list state.
     Module: frontend-blog/pages, presentation layer.
     Related: components/MomentDetailModal.vue and composables/useApi.ts. -->
<template>
  <div class="moments-page">
    <div class="container">
      <div class="page-hero">
        <h1 class="page-title">{{ t('moments.title') }}</h1>
        <p class="page-subtitle">{{ t('moments.subtitle') }}</p>
      </div>

      <div v-if="pending" class="moments-feed" aria-hidden="true">
        <article v-for="n in 4" :key="n" class="moment-item moment-skeleton">
          <div class="moment-head">
            <span class="moment-dot skeleton-dot" />
            <div class="skeleton-date"></div>
          </div>
          <div class="skeleton-content">
            <div class="skeleton-line wide" />
            <div class="skeleton-line wide" />
            <div class="skeleton-line mid" />
          </div>
        </article>
      </div>

      <div v-else-if="allMoments.length" class="moments-feed">
        <article
          v-for="(moment, i) in allMoments"
          :key="moment.id"
          class="moment-item"
          :class="{ active: selectedMomentId === moment.id }"
          :style="{ animationDelay: `${i * 50}ms` }"
        >
          <button
            class="moment-trigger"
            type="button"
            @click="openMoment(moment.id)"
            @mousemove="setPointerPosition($event)"
            @mouseleave="resetPointerPosition($event)"
          >
            <div class="moment-head">
              <span class="moment-dot" />
              <time class="moment-date">{{ formatDate(moment.publishedAt || moment.createdAt) }}</time>
            </div>

            <div class="moment-body">
              <div class="moment-content markdown-body" v-html="renderMomentPreview(moment.content)"></div>

              <div class="moment-actions">
                <span class="comment-count">{{ t('moments.commentCount', { n: moment.commentCount || 0 }) }}</span>
              </div>
            </div>
          </button>
        </article>
      </div>

      <div v-else class="empty-state">{{ t('moments.empty') }}</div>

      <InfiniteScrollTrigger
        v-if="allMoments.length > 0"
        :loading="loadingMore"
        :done="!hasMore"
        :done-text="t('moments.noMore')"
        @load="loadMore"
      >
        <template #loading>
          <div class="load-more-state">
            <div class="spinner"></div>
            <span>{{ t('moments.loadingMore') }}</span>
          </div>
        </template>
      </InfiniteScrollTrigger>
    </div>

    <MomentDetailModal
      :open="!!selectedMomentId"
      :moment="selectedMoment"
      :loading="detailPending"
      :comment-enabled="siteSettings?.commentEnabled !== false"
      :require-approval="siteSettings?.commentRequireApproval"
      :previous-moment="previousMoment"
      :next-moment="nextMoment"
      @close="closeMoment"
      @count-change="handleDetailCountChange"
      @prev="openPreviousMoment"
      @next="openNextMoment"
    />
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import InfiniteScrollTrigger from '~/components/InfiniteScrollTrigger.vue'
import MomentDetailModal from '~/components/MomentDetailModal.vue'
import type { Moment } from '~/composables/useApi'
import { renderContentMarkdown } from '~/utils/contentMarkdown'

definePageMeta({
  key: 'moments-feed',
  scrollToTop: false,
})

const { t, locale } = useI18n()
const localePath = useLocalePath()
const route = useRoute()
const router = useRouter()
const { listMoments, getMoment, getSiteSettings } = useApi()

const PAGE_SIZE = 12
const allMoments = useState<Moment[]>('moments-feed-rows', () => [])
const total = useState<number>('moments-feed-total', () => 0)
const page = useState<number>('moments-feed-page', () => 1)
const feedReady = useState<boolean>('moments-feed-ready', () => false)
const loadingMore = ref(false)
const detailPending = ref(false)
const selectedMoment = ref<Moment | null>(null)

const [{ data: siteSettings }, { pending: initialPending }] = await Promise.all([
  useAsyncData('moments-site-settings', getSiteSettings),
  useAsyncData('moments-first-page', async () => {
    if (feedReady.value) {
      return {
        rows: allMoments.value,
        total: total.value,
        page: page.value,
        pageSize: PAGE_SIZE,
      }
    }

    const res = await listMoments({ page: 1, pageSize: PAGE_SIZE })
    allMoments.value = res.rows || []
    total.value = res.total || 0
    page.value = 1
    feedReady.value = true
    return res
  }),
])

const pending = computed(() => initialPending.value && !feedReady.value)
const hasMore = computed(() => allMoments.value.length < total.value)
const selectedMomentId = computed(() => typeof route.params.id === 'string' ? route.params.id : '')
const selectedIndex = computed(() => allMoments.value.findIndex((item) => item.id === selectedMomentId.value))
const previousMoment = computed(() => {
  if (selectedIndex.value <= 0) return null
  return allMoments.value[selectedIndex.value - 1] || null
})
const nextMoment = computed(() => {
  if (selectedIndex.value < 0) return null
  return allMoments.value[selectedIndex.value + 1] || null
})

watch(selectedMomentId, async (id) => {
  if (!id) {
    selectedMoment.value = null
    detailPending.value = false
    return
  }

  const cached = allMoments.value.find((item) => item.id === id)
  if (cached) {
    selectedMoment.value = cached
  }

  detailPending.value = true
  try {
    selectedMoment.value = await getMoment(id)
    syncMoment(selectedMoment.value)
  } catch (err) {
    console.error(err)
  } finally {
    detailPending.value = false
  }
}, { immediate: true })

async function loadMore() {
  if (loadingMore.value || !hasMore.value) return
  loadingMore.value = true
  page.value += 1
  try {
    const res = await listMoments({ page: page.value, pageSize: PAGE_SIZE })
    if (res.rows?.length) {
      allMoments.value = dedupeMoments([...allMoments.value, ...res.rows])
      total.value = res.total || total.value
    }
  } catch (err) {
    page.value -= 1
    console.error(err)
  } finally {
    loadingMore.value = false
  }
}

function openMoment(momentId: string) {
  router.push(localePath(`/moments/${momentId}`))
}

function closeMoment() {
  router.push(localePath('/moments'))
}

function openPreviousMoment() {
  if (previousMoment.value) {
    openMoment(previousMoment.value.id)
  }
}

function openNextMoment() {
  if (nextMoment.value) {
    openMoment(nextMoment.value.id)
  }
}

function handleDetailCountChange(count: number) {
  if (!selectedMoment.value) return
  selectedMoment.value.commentCount = count
  syncMoment(selectedMoment.value)
}

function syncMoment(moment: Moment) {
  const index = allMoments.value.findIndex((item) => item.id === moment.id)
  if (index >= 0) {
    allMoments.value[index] = { ...allMoments.value[index], ...moment }
    return
  }
  allMoments.value = [moment, ...allMoments.value]
}

function dedupeMoments(items: Moment[]) {
  const seen = new Map<string, Moment>()
  for (const item of items) {
    seen.set(item.id, item)
  }
  return Array.from(seen.values())
}

function setPointerPosition(event: MouseEvent) {
  const target = event.currentTarget as HTMLElement | null
  if (!target) return
  const rect = target.getBoundingClientRect()
  const x = (event.clientX - rect.left) / rect.width
  const y = (event.clientY - rect.top) / rect.height
  target.style.setProperty('--pointer-x', `${(x * 100).toFixed(2)}%`)
  target.style.setProperty('--pointer-y', `${(y * 100).toFixed(2)}%`)
  target.style.setProperty('--offset-x', `${((x - 0.5) * 8).toFixed(2)}px`)
  target.style.setProperty('--offset-y', `${((y - 0.5) * 8).toFixed(2)}px`)
}

function resetPointerPosition(event: MouseEvent) {
  const target = event.currentTarget as HTMLElement | null
  if (!target) return
  target.style.setProperty('--pointer-x', '50%')
  target.style.setProperty('--pointer-y', '50%')
  target.style.setProperty('--offset-x', '0px')
  target.style.setProperty('--offset-y', '0px')
}

function formatDate(iso: string): string {
  return new Intl.DateTimeFormat(locale.value === 'en' ? 'en-US' : 'zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
  }).format(new Date(iso))
}

function renderMomentPreview(content: string): string {
  return renderContentMarkdown(content)
}

useSeoMeta({
  title: () => selectedMoment.value
    ? `${t('moments.title')} · ${formatDate(selectedMoment.value.publishedAt || selectedMoment.value.createdAt)}`
    : t('moments.title'),
})
</script>

<style scoped>
.moments-page {
  padding-top: calc(var(--header-h) + 48px);
  padding-bottom: 80px;
}

.page-hero {
  margin-bottom: 44px;
}

.page-title {
  font-size: clamp(1.8rem, 4vw, 2.6rem);
  font-weight: 800;
  letter-spacing: 0.04em;
}

.page-subtitle {
  margin-top: 12px;
  max-width: 560px;
  color: var(--text-subtle);
  line-height: 1.8;
}

.moments-feed {
  display: flex;
  flex-direction: column;
  gap: 22px;
}

.moment-item {
  animation: fade-up 0.42s var(--ease-spring) both;
}

.moment-trigger {
  --pointer-x: 50%;
  --pointer-y: 50%;
  --offset-x: 0px;
  --offset-y: 0px;
  position: relative;
  width: 100%;
  overflow: hidden;
  border: none;
  border-radius: 24px;
  background: transparent;
  padding: 18px 20px;
  text-align: left;
  cursor: pointer;
  transition:
    transform 320ms cubic-bezier(0.22, 1.18, 0.36, 1),
    background 220ms ease;
}

.moment-trigger::before {
  content: '';
  position: absolute;
  inset: 0;
  border-radius: inherit;
  background:
    radial-gradient(circle at var(--pointer-x) var(--pointer-y), color-mix(in srgb, var(--accent) 12%, transparent), transparent 38%),
    color-mix(in srgb, var(--bg-secondary) 36%, transparent);
  opacity: 0;
  transition: opacity 220ms ease;
}

.moment-trigger:hover,
.moment-item.active .moment-trigger {
  transform: translate3d(var(--offset-x), var(--offset-y), 0);
}

.moment-trigger:hover::before,
.moment-item.active .moment-trigger::before {
  opacity: 1;
}

.moment-head,
.moment-body {
  position: relative;
  z-index: 1;
}

.moment-head {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 12px;
}

.moment-dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background: var(--accent);
  box-shadow: 0 0 0 4px var(--accent-soft);
  flex-shrink: 0;
}

.moment-date {
  font-size: 12px;
  letter-spacing: 0.08em;
  text-transform: uppercase;
  color: var(--text-subtle);
}

.moment-body {
  padding-left: 20px;
}

.moment-content {
  font-size: 15px;
  line-height: 1.85;
  color: var(--text);
  word-break: break-word;
}

.moment-content :deep(p),
.moment-content :deep(ul),
.moment-content :deep(ol),
.moment-content :deep(blockquote),
.moment-content :deep(pre) {
  margin: 0 0 12px;
}

.moment-content :deep(p:last-child),
.moment-content :deep(ul:last-child),
.moment-content :deep(ol:last-child),
.moment-content :deep(blockquote:last-child),
.moment-content :deep(pre:last-child) {
  margin-bottom: 0;
}

.moment-content :deep(ul),
.moment-content :deep(ol) {
  padding-left: 20px;
}

.moment-content :deep(blockquote) {
  padding-left: 12px;
  border-left: 2px solid var(--border);
  color: var(--text-muted);
}

.moment-content :deep(pre) {
  overflow-x: auto;
  padding: 12px 14px;
  border-radius: 14px;
  background: color-mix(in srgb, var(--bg-secondary) 80%, white);
}

.moment-content :deep(code) {
  font-family: 'Fira Code', monospace;
  font-size: 0.92em;
}

.moment-content :deep(pre code) {
  background: transparent;
}

.moment-actions {
  display: flex;
  align-items: center;
  gap: 14px;
  margin-top: 14px;
}

.comment-count {
  font-size: 12px;
  color: var(--text-subtle);
}

.moment-skeleton {
  opacity: 0.72;
}

.skeleton-date,
.skeleton-line {
  display: block;
  border-radius: 999px;
  background: linear-gradient(90deg, var(--bg-secondary) 25%, var(--surface-hover) 50%, var(--bg-secondary) 75%);
  background-size: 200% 100%;
  animation: shimmer 1.4s infinite;
}

.skeleton-date {
  width: 120px;
  height: 12px;
}

.skeleton-dot {
  box-shadow: none;
  background: var(--border-strong);
}

.skeleton-content {
  padding-left: 20px;
}

.skeleton-line {
  height: 14px;
  margin-top: 8px;
}

.skeleton-line.wide {
  width: 100%;
}

.skeleton-line.mid {
  width: 70%;
}

.load-more-state,
.empty-state {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 10px;
  color: var(--text-subtle);
}

.empty-state {
  padding: 64px 0;
}

.spinner {
  width: 22px;
  height: 22px;
  border: 2px solid var(--border);
  border-top-color: var(--accent);
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
}

@keyframes fade-up {
  from {
    opacity: 0;
    transform: translateY(16px);
  }
  to {
    opacity: 1;
    transform: none;
  }
}

@keyframes shimmer {
  from {
    background-position: 200% 0;
  }
  to {
    background-position: -200% 0;
  }
}

@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}

@media (max-width: 640px) {
  .moment-trigger {
    padding: 14px 0;
    border-radius: 20px;
  }

  .moment-body,
  .skeleton-content {
    padding-left: 0;
  }
}
</style>
