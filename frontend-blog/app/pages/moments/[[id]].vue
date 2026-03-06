<!-- File: app/pages/moments/[[id]].vue
     Purpose: Moments feed with URL-driven detail modal for comments and deep linking.
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
          <button class="moment-trigger" type="button" @click="openMoment(moment.id)">
            <div class="moment-head">
              <span class="moment-dot" />
              <time class="moment-date">{{ formatDate(moment.publishedAt || moment.createdAt) }}</time>
            </div>

            <div class="moment-body">
              <p class="moment-content">{{ moment.content }}</p>

              <div class="moment-actions">
                <span class="comment-count">{{ t('moments.commentCount', { n: moment.commentCount || 0 }) }}</span>
                <span class="detail-link">{{ t('moments.viewDetail') }}</span>
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
      :require-approval="siteSettings?.commentRequireApproval"
      @close="closeMoment"
      @count-change="handleDetailCountChange"
    />
  </div>
</template>

<script setup lang="ts">
import { useI18n } from 'vue-i18n'
import InfiniteScrollTrigger from '~/components/InfiniteScrollTrigger.vue'
import MomentDetailModal from '~/components/MomentDetailModal.vue'
import type { Moment } from '~/composables/useApi'

const { t, locale } = useI18n()
const localePath = useLocalePath()
const route = useRoute()
const router = useRouter()
const { listMoments, getMoment, getSiteSettings } = useApi()

const PAGE_SIZE = 12
const page = ref(1)
const allMoments = ref<Moment[]>([])
const total = ref(0)
const loadingMore = ref(false)
const detailPending = ref(false)
const selectedMoment = ref<Moment | null>(null)

const [{ data: siteSettings }, { data: initialData, pending }] = await Promise.all([
  useAsyncData('moments-site-settings', getSiteSettings),
  useAsyncData('moments-first-page', () => listMoments({ page: 1, pageSize: PAGE_SIZE }))
])

watch(initialData, (newVal) => {
  if (!newVal) return
  allMoments.value = newVal.rows || []
  total.value = newVal.total || 0
  page.value = 1
}, { immediate: true })

const hasMore = computed(() => allMoments.value.length < total.value)
const selectedMomentId = computed(() => typeof route.params.id === 'string' ? route.params.id : '')

watch([selectedMomentId, allMoments], async ([id]) => {
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
      allMoments.value.push(...res.rows)
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

function handleDetailCountChange(count: number) {
  if (!selectedMoment.value) return
  selectedMoment.value.commentCount = count
  syncMoment(selectedMoment.value)
}

function syncMoment(moment: Moment) {
  const index = allMoments.value.findIndex((item) => item.id === moment.id)
  if (index >= 0) {
    allMoments.value[index] = { ...allMoments.value[index], ...moment }
  }
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
  gap: 30px;
}

.moment-item {
  animation: fade-up 0.42s var(--ease-spring) both;
}

.moment-item.active .moment-trigger {
  background: color-mix(in srgb, var(--accent) 4%, transparent);
}

.moment-trigger {
  width: 100%;
  border: none;
  background: transparent;
  padding: 0;
  text-align: left;
  cursor: pointer;
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
  white-space: pre-wrap;
  word-break: break-word;
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

.detail-link {
  font-size: 13px;
  color: var(--accent);
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
  .moment-body,
  .skeleton-content {
    padding-left: 0;
  }
}
</style>
