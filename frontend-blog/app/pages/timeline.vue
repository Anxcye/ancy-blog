<!-- File: app/pages/timeline.vue
     Purpose: Mixed timeline of articles and moments in chronological order.
     Module: frontend-blog/pages, presentation layer.
     Related: composables/useApi.ts. -->
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
            <span class="hero-stat">{{ t('timeline.total', { n: items?.total || items?.rows?.length || 0 }) }}</span>
          </div>
        </div>
      </div>

      <div v-if="pending" class="timeline">
        <div v-for="n in 8" :key="n" class="tl-item">
          <div class="tl-dot" />
          <div class="tl-body">
            <div class="skeleton-line" style="height: 12px; width: 30%;" />
            <div class="skeleton-line" style="height: 18px; width: 75%; margin-top: 6px;" />
          </div>
        </div>
      </div>

      <div v-else-if="items?.rows?.length" class="timeline">
        <div
          v-for="(item, i) in items.rows"
          :key="item.id"
          class="tl-item"
          :class="`tl-item--${item.type}`"
          :style="{ animationDelay: `${i * 50}ms` }"
        >
          <div class="tl-dot">
            <!-- Article icon -->
            <svg v-if="item.type === 'article'" viewBox="0 0 16 16" fill="none" stroke="currentColor" stroke-width="1.5">
              <rect x="2" y="2" width="12" height="12" rx="2"/>
              <path d="M5 6h6M5 9h4"/>
            </svg>
            <!-- Moment icon -->
            <svg v-else viewBox="0 0 16 16" fill="none" stroke="currentColor" stroke-width="1.5">
              <circle cx="8" cy="8" r="3"/>
            </svg>
          </div>

          <div class="tl-body">
            <time class="tl-date">{{ formatDate(item.publishedAt) }}</time>
            <span class="tl-type-badge" :class="`badge--${item.type}`">
              {{ item.type === 'article' ? '文章' : '瞬间' }}
            </span>

            <!-- Article: click to navigate -->
            <template v-if="item.type === 'article' && item.slug">
              <NuxtLink :to="localePath(`/articles/${item.slug}`)" class="tl-title">
                {{ item.title }}
              </NuxtLink>
              <p v-if="item.summary" class="tl-summary">{{ item.summary }}</p>
            </template>

            <!-- Moment: inline content -->
            <template v-else>
              <p class="tl-moment-content">{{ item.content }}</p>
            </template>
          </div>
        </div>

        <!-- Pagination -->
        <div v-if="totalPages > 1" class="pagination">
          <button class="page-btn" :disabled="page <= 1" @click="page--">←</button>
          <span class="page-info">{{ page }} / {{ totalPages }}</span>
          <button class="page-btn" :disabled="page >= totalPages" @click="page++">→</button>
        </div>
      </div>

      <div v-else class="empty-state">{{ t('timeline.empty') }}</div>
    </div>
  </div>
</template>

<script setup lang="ts">
const { t } = useI18n()
const localePath = useLocalePath()
const { listTimeline } = useApi()

const PAGE_SIZE = 20
const page = ref(1)

const { data: items, pending } = await useAsyncData(
  'timeline-list',
  () => listTimeline({ page: page.value, pageSize: PAGE_SIZE }),
  {
    watch: [page],
    getCachedData: () => undefined,
  }
)

const totalPages = computed(() =>
  items.value ? Math.ceil(items.value.total / PAGE_SIZE) : 1
)

function formatDate(iso: string): string {
  return new Date(iso).toLocaleDateString('zh-CN', {
    year: 'numeric', month: 'numeric', day: 'numeric',
  })
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

/* Timeline */
.timeline { position: relative; padding-left: 36px; }
.timeline::before {
  content: '';
  position: absolute;
  left: 11px;
  top: 16px;
  bottom: 0;
  width: 1px;
  background: linear-gradient(to bottom, var(--accent), transparent);
  opacity: 0.25;
}

.tl-item {
  position: relative;
  display: flex;
  gap: 20px;
  margin-bottom: 36px;
  animation: fade-up 0.4s var(--ease-spring) both;
}

@keyframes fade-up {
  from { opacity: 0; transform: translateY(14px); }
  to   { opacity: 1; transform: none; }
}

.tl-dot {
  position: absolute;
  left: -29px;
  top: 2px;
  width: 22px;
  height: 22px;
  border-radius: 50%;
  background: var(--surface);
  border: 1.5px solid var(--border);
  display: grid;
  place-items: center;
  color: var(--text-subtle);
  transition: border-color var(--dur-fast), color var(--dur-fast), box-shadow var(--dur-fast);
}

.tl-dot svg { width: 11px; height: 11px; }

.tl-item--article .tl-dot {
  border-color: var(--accent);
  color: var(--accent-text);
  box-shadow: 0 0 0 4px var(--accent-soft);
}

.tl-body { flex: 1; min-width: 0; }

.tl-date { font-size: 11px; color: var(--text-subtle); display: block; margin-bottom: 4px; }

.tl-type-badge {
  font-size: 10px;
  font-weight: 700;
  text-transform: uppercase;
  letter-spacing: 0.05em;
  padding: 2px 8px;
  border-radius: 99px;
  margin-bottom: 6px;
  display: inline-block;
}

.badge--article { background: var(--accent-soft); color: var(--accent-text); }
.badge--moment  { background: var(--bg-secondary); color: var(--text-subtle); }

.tl-title {
  display: block;
  font-size: 15px;
  font-weight: 700;
  color: var(--text);
  transition: color var(--dur-fast);
  margin-bottom: 4px;
}
.tl-title:hover { color: var(--accent-text); }

.tl-summary {
  font-size: 13px;
  color: var(--text-muted);
  line-height: 1.6;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.tl-moment-content {
  font-size: 14px;
  color: var(--text-muted);
  line-height: 1.65;
  white-space: pre-wrap;
  word-break: break-word;
}

/* Skeleton */
.skeleton-line {
  background: linear-gradient(90deg, var(--bg-secondary) 25%, var(--surface-hover) 50%, var(--bg-secondary) 75%);
  background-size: 200% 100%;
  border-radius: var(--radius-sm);
  animation: shimmer 1.4s infinite;
}
@keyframes shimmer {
  from { background-position: 200% 0; }
  to   { background-position: -200% 0; }
}

/* Pagination */
.pagination { display: flex; align-items: center; justify-content: center; gap: 20px; margin-top: 20px; }
.page-btn {
  width: 38px; height: 38px; border-radius: var(--radius-md);
  border: 1px solid var(--border); background: var(--surface);
  color: var(--text-muted); font-size: 16px; display: grid; place-items: center;
  transition: all var(--dur-fast); cursor: pointer;
}
.page-btn:hover:not(:disabled) { border-color: var(--accent); color: var(--accent-text); background: var(--accent-soft); }
.page-btn:disabled { opacity: 0.35; cursor: not-allowed; }
.page-info { font-size: 14px; color: var(--text-muted); min-width: 60px; text-align: center; }

.empty-state { text-align: center; padding: 64px 0; color: var(--text-subtle); font-size: 15px; }

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
}
</style>
