<!-- File: app/pages/timeline.vue
     Purpose: Present a grouped chronological stream of articles and moments with infinite loading.
     Module: frontend-blog/pages, presentation layer.
     Related: app/components/MomentDetailModal.vue, app/components/InfiniteScrollTrigger.vue, composables/useApi.ts. -->
<template>
  <div class="timeline-page">
    <div class="container">
      <div class="page-hero">
        <span class="hero-eyebrow">{{ t('timeline.eyebrow') }}</span>
        <div class="hero-main">
          <div class="hero-copy">
            <h1 class="page-title">{{ t('timeline.title') }}</h1>
            <p class="page-subtitle">{{ t('timeline.subtitle') }}</p>
          </div>
          <div class="hero-stats">
            <span class="hero-stat">{{ t('timeline.total', { n: total || allItems.length }) }}</span>
            <span class="hero-stat muted">{{ t('timeline.mixedLabel') }}</span>
          </div>
        </div>
      </div>

      <div v-if="pending" class="timeline-groups timeline-groups--loading" aria-hidden="true">
        <section v-for="group in 3" :key="group" class="timeline-group timeline-group--loading">
          <div class="group-head group-head--loading">
            <div class="skeleton-pill skeleton-month"></div>
            <div class="skeleton-line skeleton-count"></div>
          </div>
          <div class="group-stream">
            <article v-for="card in 4" :key="card" class="tl-row tl-row--loading">
              <div class="skeleton-line skeleton-date"></div>
              <div class="skeleton-line skeleton-text"></div>
              <div class="skeleton-line skeleton-tail"></div>
            </article>
          </div>
        </section>
      </div>

      <div v-else-if="groupedItems.length" class="timeline-groups">
        <section
          v-for="(group, groupIndex) in groupedItems"
          :key="group.key"
          class="timeline-group"
          :style="{ animationDelay: `${groupIndex * 60}ms` }"
        >
          <button
            class="group-head"
            type="button"
            :aria-expanded="isGroupExpanded(group.key)"
            @click="toggleGroup(group.key)"
          >
            <span class="group-orbit" aria-hidden="true"></span>
            <div class="group-copy">
              <span class="group-kicker">{{ group.year }}</span>
              <h2 class="group-title">{{ group.label }}</h2>
            </div>
            <div class="group-meta">
              <span class="group-total">{{ t('timeline.groupCount', { n: group.items.length }) }}</span>
              <span class="group-toggle" aria-hidden="true">{{ isGroupExpanded(group.key) ? '−' : '+' }}</span>
            </div>
          </button>

          <div v-if="isGroupExpanded(group.key)" class="group-stream">
            <template v-for="(item, itemIndex) in group.items" :key="`${group.key}-${item.id}`">
              <NuxtLink
                v-if="item.contentType === 'article' && item.slug"
                :to="localePath(`/articles/${item.slug}`)"
                class="tl-row tl-row--article"
                :style="{ animationDelay: `${itemIndex * 45}ms` }"
              >
                <time class="tl-date">{{ formatItemDate(item.publishedAt) }}</time>
                <span class="tl-text">{{ getItemLabel(item) }}</span>
                <span class="tl-tail">{{ getItemTail(item) }}</span>
              </NuxtLink>

              <button
                v-else
                type="button"
                class="tl-row tl-row--moment"
                :style="{ animationDelay: `${itemIndex * 45}ms` }"
                @click="openMoment(item.id)"
              >
                <time class="tl-date">{{ formatItemDate(item.publishedAt) }}</time>
                <span class="tl-text">{{ getItemLabel(item) }}</span>
                <span class="tl-tail">{{ getItemTail(item) }}</span>
              </button>
            </template>
          </div>
        </section>
      </div>

      <div v-else class="empty-state">{{ t('timeline.empty') }}</div>

      <InfiniteScrollTrigger
        v-if="allItems.length > 0"
        :loading="loadingMore"
        :done="!hasMore"
        :done-text="t('timeline.noMore')"
        @load="loadMore"
      >
        <template #loading>
          <div class="load-more-state">
            <div class="spinner"></div>
            <span>{{ t('timeline.loadingMore') }}</span>
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
import { computed, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import InfiniteScrollTrigger from '~/components/InfiniteScrollTrigger.vue'
import MomentDetailModal from '~/components/MomentDetailModal.vue'
import type { Moment, TimelineItem } from '~/composables/useApi'

const { t, locale } = useI18n()
const localePath = useLocalePath()
const { listTimeline, getMoment, getSiteSettings } = useApi()

const PAGE_SIZE = 18
const page = ref(1)
const allItems = ref<TimelineItem[]>([])
const total = ref(0)
const loadingMore = ref(false)
const detailPending = ref(false)
const selectedMomentId = ref('')
const selectedMoment = ref<Moment | null>(null)
const expandedGroups = ref<string[]>([])

const [{ data: siteSettings }, { data: initialData, pending }] = await Promise.all([
  useAsyncData('timeline-site-settings', getSiteSettings),
  useAsyncData(
    'timeline-first-page',
    () => listTimeline({ page: 1, pageSize: PAGE_SIZE }),
    { getCachedData: () => undefined },
  ),
])

watch(initialData, (value) => {
  if (!value) return
  allItems.value = value.rows || []
  total.value = value.total || 0
  page.value = 1
}, { immediate: true })

const hasMore = computed(() => allItems.value.length < total.value)

const groupedItems = computed(() => {
  const groups = new Map<string, { key: string; year: string; label: string; items: TimelineItem[] }>()

  for (const item of allItems.value) {
    const date = new Date(item.publishedAt)
    const year = `${date.getFullYear()}`
    const month = `${date.getMonth() + 1}`.padStart(2, '0')
    const key = `${year}-${month}`
    if (!groups.has(key)) {
      const label = new Intl.DateTimeFormat(locale.value === 'en' ? 'en-US' : 'zh-CN', {
        month: locale.value === 'en' ? 'long' : 'numeric',
      }).format(date)
      groups.set(key, {
        key,
        year,
        label,
        items: [],
      })
    }
    groups.get(key)!.items.push(item)
  }

  return Array.from(groups.values())
})

watch(groupedItems, (groups) => {
  if (!groups.length) return
  const known = new Set(expandedGroups.value)
  if (!expandedGroups.value.length) {
    expandedGroups.value = [groups[0].key]
    return
  }
  expandedGroups.value = expandedGroups.value.filter((key) => groups.some((group) => group.key === key))
  for (const group of groups) {
    if (known.has(group.key)) continue
    // Keep newly loaded older months collapsed by default.
  }
}, { immediate: true })

const timelineMoments = computed(() => {
  return allItems.value
    .filter((item) => item.contentType === 'moment')
    .map((item) => toMomentStub(item))
})

const selectedMomentIndex = computed(() => timelineMoments.value.findIndex((item) => item.id === selectedMomentId.value))

const previousMoment = computed(() => {
  if (selectedMomentIndex.value <= 0) return null
  return timelineMoments.value[selectedMomentIndex.value - 1] || null
})

const nextMoment = computed(() => {
  if (selectedMomentIndex.value < 0) return null
  return timelineMoments.value[selectedMomentIndex.value + 1] || null
})

async function loadMore() {
  if (loadingMore.value || !hasMore.value) return
  loadingMore.value = true
  page.value += 1
  try {
    const res = await listTimeline({ page: page.value, pageSize: PAGE_SIZE })
    allItems.value = dedupeTimelineItems([...allItems.value, ...(res.rows || [])])
    total.value = res.total || total.value
  } catch (err) {
    page.value -= 1
    console.error(err)
  } finally {
    loadingMore.value = false
  }
}

watch(selectedMomentId, async (id) => {
  if (!id) {
    selectedMoment.value = null
    detailPending.value = false
    return
  }

  const cached = timelineMoments.value.find((item) => item.id === id)
  if (cached) {
    selectedMoment.value = cached
  }

  detailPending.value = true
  try {
    selectedMoment.value = await getMoment(id)
  } catch (err) {
    console.error(err)
  } finally {
    detailPending.value = false
  }
}, { immediate: true })

function openMoment(id: string) {
  selectedMomentId.value = id
}

function closeMoment() {
  selectedMomentId.value = ''
}

function isGroupExpanded(key: string) {
  return expandedGroups.value.includes(key)
}

function toggleGroup(key: string) {
  expandedGroups.value = isGroupExpanded(key)
    ? expandedGroups.value.filter((item) => item !== key)
    : [...expandedGroups.value, key]
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
}

function toMomentStub(item: TimelineItem): Moment {
  const cached = selectedMoment.value?.id === item.id ? selectedMoment.value : null
  return {
    id: item.id,
    content: item.content || '',
    status: cached?.status || 'published',
    allowComment: cached?.allowComment ?? true,
    commentCount: cached?.commentCount || 0,
    isPinned: cached?.isPinned,
    publishedAt: item.publishedAt,
    createdAt: item.publishedAt,
  }
}

function dedupeTimelineItems(items: TimelineItem[]) {
  const seen = new Map<string, TimelineItem>()
  for (const item of items) {
    seen.set(`${item.contentType}:${item.id}`, item)
  }
  return Array.from(seen.values())
}

function formatItemDate(iso: string): string {
  return new Intl.DateTimeFormat(locale.value === 'en' ? 'en-US' : 'zh-CN', {
    day: '2-digit',
    month: '2-digit',
  }).format(new Date(iso))
}

function getItemLabel(item: TimelineItem): string {
  if (item.contentType === 'article') {
    return item.title || t('timeline.untitledArticle')
  }

  return getMomentPreviewText(item.content || '')
}

function getItemTail(item: TimelineItem): string {
  if (item.contentType === 'article') {
    const category = item.categoryName || item.categorySlug
    return category ? `${t('timeline.articleLabel')} / ${category}` : t('timeline.articleLabel')
  }

  return t('timeline.momentLabel')
}

function getMomentPreviewText(content: string): string {
  const plain = content
    .replace(/[#>*_`~\-\[\]\(\)!]/g, ' ')
    .replace(/\s+/g, ' ')
    .trim()
  return plain.slice(0, 42) + (plain.length > 42 ? '…' : '')
}

useSeoMeta({ title: t('timeline.title') })
</script>

<style scoped>
.timeline-page {
  padding-top: calc(var(--header-h) + 48px);
  padding-bottom: 80px;
}

.page-hero {
  display: flex;
  flex-direction: column;
  gap: 16px;
  margin-bottom: 40px;
  padding: 4px 0 26px;
  border-bottom: 1px solid color-mix(in srgb, var(--border) 78%, transparent);
}

.hero-eyebrow {
  font-size: 11px;
  font-weight: 700;
  letter-spacing: 0.18em;
  text-transform: uppercase;
  color: var(--accent);
}

.hero-main {
  display: flex;
  align-items: flex-end;
  justify-content: space-between;
  gap: 24px;
  flex-wrap: wrap;
}

.hero-copy {
  max-width: 620px;
}

.page-title {
  font-size: clamp(1.8rem, 4vw, 2.8rem);
  font-weight: 800;
  letter-spacing: -0.02em;
  margin: 0;
}

.page-subtitle {
  margin: 10px 0 0;
  max-width: 560px;
  font-size: 15px;
  line-height: 1.8;
  color: var(--text-subtle);
}

.hero-stats {
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
}

.hero-stat {
  display: inline-flex;
  align-items: center;
  min-height: 34px;
  padding: 0 12px;
  border-radius: 999px;
  border: 1px solid color-mix(in srgb, var(--border) 78%, transparent);
  background: color-mix(in srgb, var(--bg-secondary) 68%, transparent);
  color: var(--text-muted);
  font-size: 12px;
  white-space: nowrap;
}

.hero-stat.muted {
  color: var(--text-subtle);
}

.timeline-groups {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.timeline-group {
  display: flex;
  flex-direction: column;
  gap: 4px;
  animation: fade-up 0.48s var(--ease-spring) both;
}

.group-head {
  position: relative;
  display: grid;
  grid-template-columns: minmax(0, 1fr) auto;
  align-items: center;
  gap: 14px;
  width: 100%;
  padding: 7px 0 7px 24px;
  border: none;
  background: transparent;
  text-align: left;
  cursor: pointer;
}

.group-head--loading {
  position: sticky;
  padding-left: 24px;
}

.group-orbit {
  position: absolute;
  left: 0;
  top: 2px;
  width: 14px;
  height: 14px;
  border-radius: 50%;
  background: color-mix(in srgb, var(--surface) 86%, white);
  border: 2px solid var(--accent);
  box-shadow:
    0 0 0 8px color-mix(in srgb, var(--accent) 8%, transparent),
    0 10px 18px rgba(8, 16, 24, 0.06);
}

.group-kicker {
  font-size: 11px;
  letter-spacing: 0.16em;
  text-transform: uppercase;
  color: var(--text-subtle);
}

.group-title {
  margin: 0;
  font-size: clamp(1.12rem, 2vw, 1.48rem);
  font-weight: 800;
  letter-spacing: -0.03em;
  color: var(--text);
}

.group-total {
  display: inline-flex;
  align-items: center;
  min-height: 24px;
  padding: 0 8px;
  border-radius: 999px;
  background: color-mix(in srgb, var(--bg-secondary) 70%, transparent);
  color: var(--text-subtle);
  font-size: 11px;
}

.group-meta {
  display: inline-flex;
  align-items: center;
  gap: 8px;
}

.group-toggle {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 18px;
  height: 18px;
  border-radius: 999px;
  color: var(--text-subtle);
  font-size: 14px;
}

.group-stream {
  position: relative;
  display: flex;
  flex-direction: column;
  gap: 0;
  padding-left: 24px;
}

.group-stream::before {
  content: '';
  position: absolute;
  left: 5px;
  top: 0;
  bottom: 0;
  width: 1px;
  background: linear-gradient(to bottom, color-mix(in srgb, var(--accent) 26%, transparent), transparent);
}

.tl-row {
  position: relative;
  display: grid;
  grid-template-columns: 62px minmax(0, 1fr) auto;
  align-items: center;
  gap: 12px;
  width: 100%;
  min-height: 30px;
  padding: 2px 0;
  border: none;
  background: transparent;
  text-align: left;
  text-decoration: none;
  color: inherit;
  cursor: pointer;
  transition: transform 180ms cubic-bezier(0.22, 1.18, 0.36, 1), color 180ms ease;
  animation: fade-up 0.34s var(--ease-spring) both;
}

.tl-row::before {
  content: '';
  position: absolute;
  left: -22px;
  top: 50%;
  width: 7px;
  height: 7px;
  transform: translateY(-50%);
  border-radius: 50%;
  background: var(--surface);
  border: 1.5px solid var(--border);
  box-shadow: 0 0 0 4px color-mix(in srgb, var(--surface) 72%, transparent);
  transition: border-color 180ms ease, box-shadow 180ms ease;
}

.tl-row--article::before {
  border-color: var(--accent);
  box-shadow: 0 0 0 4px color-mix(in srgb, var(--accent) 10%, transparent);
}

.tl-row--moment::before {
  background: color-mix(in srgb, var(--bg-secondary) 80%, white);
}

.tl-row:hover {
  transform: translateX(4px);
}

.tl-date {
  font-size: 11px;
  color: var(--text-subtle);
  font-variant-numeric: tabular-nums;
  letter-spacing: 0.04em;
}

.tl-text {
  margin: 0;
  min-width: 0;
  font-size: 14px;
  line-height: 1.2;
  color: var(--text);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.tl-row--article .tl-text {
  font-weight: 600;
}

.tl-row--moment .tl-text {
  color: var(--text-muted);
}

.tl-tail {
  font-size: 12px;
  color: var(--text-subtle);
  white-space: nowrap;
}

.timeline-group--loading .group-head {
  position: static;
}

.skeleton-pill,
.skeleton-line {
  display: block;
  border-radius: 999px;
  background: linear-gradient(90deg, var(--bg-secondary) 25%, var(--surface-hover) 50%, var(--bg-secondary) 75%);
  background-size: 200% 100%;
  animation: shimmer 1.35s infinite;
}

.skeleton-month {
  width: 82px;
  height: 28px;
}

.skeleton-count {
  width: 58px;
  height: 10px;
}

.tl-row--loading {
  cursor: default;
  pointer-events: none;
}

.skeleton-date {
  width: 44px;
  height: 10px;
}

.skeleton-text {
  width: 100%;
  height: 12px;
}

.skeleton-tail {
  width: 36px;
  height: 10px;
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

@media (max-width: 900px) {
  .group-head {
    padding-left: 20px;
  }

  .group-orbit {
    top: 9px;
  }
}

@media (max-width: 640px) {
  .page-hero {
    gap: 14px;
    margin-bottom: 34px;
    padding-bottom: 22px;
  }

  .hero-main {
    align-items: flex-start;
    gap: 18px;
  }

  .page-subtitle {
    font-size: 14px;
  }

  .timeline-groups {
    gap: 10px;
  }

  .group-stream {
    padding-left: 20px;
  }

  .tl-row {
    grid-template-columns: 50px minmax(0, 1fr) auto;
    gap: 8px;
    min-height: 28px;
  }

  .tl-tail {
    font-size: 11px;
  }
}
</style>
