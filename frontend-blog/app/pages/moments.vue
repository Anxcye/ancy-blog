<!-- File: app/pages/moments.vue
     Purpose: Moments (short content) waterfall feed.
     Module: frontend-blog/pages, presentation layer.
     Related: composables/useApi.ts. -->
<template>
  <div class="moments-page">
    <div class="container">
      <div class="page-hero">
        <h1 class="page-title">{{ t('moments.title') }}</h1>
      </div>

      <!-- Skeleton -->
      <div v-if="pending" class="moments-feed">
        <div v-for="n in 5" :key="n" class="moment-item moment-skeleton">
          <div class="skeleton-dot" />
          <div class="skeleton-content">
            <div class="skeleton-line" style="height: 14px; width: 90%;" />
            <div class="skeleton-line" style="height: 14px; width: 70%; margin-top: 6px;" />
            <div class="skeleton-line" style="height: 11px; width: 30%; margin-top: 10px;" />
          </div>
        </div>
      </div>

      <!-- Feed -->
      <div v-else-if="moments?.rows?.length" class="moments-feed">
        <div
          v-for="(moment, i) in moments.rows"
          :key="moment.id"
          class="moment-item"
          :style="{ animationDelay: `${i * 60}ms` }"
        >
          <div class="moment-dot" />
          <div class="moment-body">
            <p class="moment-content">{{ moment.content }}</p>
            <time class="moment-date">{{ formatDate(moment.publishedAt || moment.createdAt) }}</time>
          </div>
        </div>

        <!-- Pagination -->
        <div v-if="totalPages > 1" class="pagination">
          <button class="page-btn" :disabled="page <= 1" @click="page--">←</button>
          <span class="page-info">{{ page }} / {{ totalPages }}</span>
          <button class="page-btn" :disabled="page >= totalPages" @click="page++">→</button>
        </div>
      </div>

      <div v-else class="empty-state">{{ t('moments.empty') }}</div>
    </div>
  </div>
</template>

<script setup lang="ts">
const { t } = useI18n()
const { listMoments } = useApi()

const PAGE_SIZE = 20
const page = ref(1)

const { data: moments, pending } = await useAsyncData(
  'moments-list',
  () => listMoments({ page: page.value, pageSize: PAGE_SIZE }),
  {
    watch: [page],
    getCachedData: () => undefined,
  }
)

const totalPages = computed(() =>
  moments.value ? Math.ceil(moments.value.total / PAGE_SIZE) : 1
)

function formatDate(iso: string): string {
  return new Date(iso).toLocaleDateString('zh-CN', {
    year: 'numeric', month: 'numeric', day: 'numeric',
    hour: '2-digit', minute: '2-digit',
  })
}

useSeoMeta({ title: t('moments.title') })
</script>

<style scoped>
.moments-page {
  padding-top: calc(var(--header-h) + 48px);
  padding-bottom: 80px;
}

.page-hero { margin-bottom: 40px; }
.page-title { font-size: clamp(1.5rem, 3vw, 2rem); font-weight: 800; }

/* Feed timeline */
.moments-feed {
  position: relative;
  padding-left: 24px;
}

.moments-feed::before {
  content: '';
  position: absolute;
  left: 5px;
  top: 10px;
  bottom: 0;
  width: 1px;
  background: linear-gradient(to bottom, var(--accent), transparent);
  opacity: 0.3;
}

.moment-item {
  position: relative;
  display: flex;
  gap: 20px;
  margin-bottom: 32px;
  animation: fade-up 0.4s var(--ease-spring) both;
}

@keyframes fade-up {
  from { opacity: 0; transform: translateY(16px); }
  to   { opacity: 1; transform: none; }
}

.moment-dot {
  position: absolute;
  left: -19px;
  top: 6px;
  width: 9px;
  height: 9px;
  border-radius: 50%;
  background: var(--accent);
  box-shadow: 0 0 0 3px var(--accent-soft);
  flex-shrink: 0;
}

.moment-body { flex: 1; }

.moment-content {
  font-size: 15px;
  line-height: 1.75;
  color: var(--text);
  white-space: pre-wrap;
  word-break: break-word;
  margin-bottom: 8px;
}

.moment-date { font-size: 12px; color: var(--text-subtle); }

/* Skeleton */
.moment-skeleton { opacity: 0.6; }
.skeleton-dot {
  position: absolute;
  left: -19px;
  top: 6px;
  width: 9px;
  height: 9px;
  border-radius: 50%;
  background: var(--border-strong);
}
.skeleton-content { flex: 1; }
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
</style>
