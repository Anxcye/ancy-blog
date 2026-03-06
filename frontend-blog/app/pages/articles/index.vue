<!-- File: app/pages/articles/index.vue
     Purpose: Article list page with category/tag filter pills and pagination.
     Module: frontend-blog/pages/articles, presentation layer.
     Related: components/ArticleCard.vue, composables/useApi.ts. -->
<template>
  <div class="articles-page">
    <div class="container">

      <!-- Page title -->
      <div class="page-hero">
        <h1 class="page-title">{{ t('nav.articles') }}</h1>
        <p class="page-subtitle">{{ t('home.recentArticles') }}</p>
      </div>

      <!-- Filter pills -->
      <div class="filter-bar">
        <button
          class="filter-pill"
          :class="{ active: !activeCategory && !activeTag }"
          @click="clearFilters"
        >
          全部
        </button>

        <template v-if="categories?.length">
          <button
            v-for="cat in categories"
            :key="cat.slug"
            class="filter-pill"
            :class="{ active: activeCategory === cat.slug }"
            @click="setCategory(cat.slug)"
          >
            {{ cat.name }}
          </button>
        </template>
      </div>

      <!-- Tag cloud -->
      <div v-if="tags?.length" class="tag-cloud">
        <button
          v-for="tag in tags"
          :key="tag.slug"
          class="tag-pill"
          :class="{ active: activeTag === tag.slug }"
          @click="setTag(tag.slug)"
        >
          #{{ tag.name }}
        </button>
      </div>

      <!-- Articles -->
      <div v-if="pending" class="article-list">
        <div v-for="n in 6" :key="n" class="skeleton-article-item">
          <div class="skeleton-inner">
            <div class="skeleton-line" style="height: 32px; width: 55%; margin-bottom: 14px;" />
            <div style="display:flex; gap:8px;">
              <div class="skeleton-line" style="height: 20px; width: 60px;" />
              <div class="skeleton-line" style="height: 20px; width: 80px;" />
              <div class="skeleton-line" style="height: 20px; width: 50px;" />
            </div>
          </div>
        </div>
      </div>

      <div v-else-if="allArticles.length" class="article-list">
        <NuxtLink
          v-for="(article, i) in allArticles"
          :key="article.id"
          :to="localePath(`/articles/${article.slug}`)"
          class="article-item"
          :style="{ animationDelay: `${(i % 9) * 70}ms` }"
        >
          <!-- Left: index number (global index based on i) -->
          <span class="article-index">{{ String(i + 1).padStart(2, '0') }}</span>

          <!-- Center: title + meta -->
          <div class="article-body">
            <h2 class="article-title">{{ article.title }}</h2>

            <div class="article-meta">
              <!-- Date -->
              <span class="meta-item meta-date">
                <svg class="meta-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect x="3" y="4" width="18" height="18" rx="2" ry="2"></rect><line x1="16" y1="2" x2="16" y2="6"></line><line x1="8" y1="2" x2="8" y2="6"></line><line x1="3" y1="10" x2="21" y2="10"></line></svg>
                <time>{{ formatDate(article.publishedAt || article.createdAt) }}</time>
              </span>

              <!-- Category -->
              <span v-if="article.categorySlug" class="meta-item meta-category" @click.prevent="setCategory(article.categorySlug)">
                <svg class="meta-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M22 19a2 2 0 0 1-2 2H4a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h5l2 3h9a2 2 0 0 1 2 2z"></path></svg>
                <span>{{ getCategoryName(article.categorySlug) }}</span>
              </span>

              <!-- Tags (max 3) -->
              <template v-if="article.tagSlugs?.length">
                <span
                  v-for="slug in article.tagSlugs.slice(0, 3)"
                  :key="slug"
                  class="meta-item meta-tag"
                  @click.prevent="setTag(slug)"
                >
                  <svg class="meta-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><line x1="4" y1="9" x2="20" y2="9"></line><line x1="4" y1="15" x2="20" y2="15"></line><line x1="10" y1="3" x2="8" y2="21"></line><line x1="16" y1="3" x2="14" y2="21"></line></svg>
                  <span>{{ getTagName(slug) }}</span>
                </span>
              </template>
              
              <!-- Views -->
              <span class="meta-item meta-views">
                <svg class="meta-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z"></path><circle cx="12" cy="12" r="3"></circle></svg>
                <span>{{ article.viewCount }}</span>
              </span>
            </div>
          </div>

          <!-- Right: arrow indicator -->
          <span class="article-arrow">→</span>
        </NuxtLink>
      </div>

      <div v-else class="empty-state">
        <p>{{ t('home.noArticles') }}</p>
      </div>

      <!-- Infinite Scroll Trigger -->
      <InfiniteScrollTrigger
        v-if="allArticles.length > 0"
        :loading="loadingMore"
        :done="!hasMore"
        @load="loadMore"
      />

    </div>
  </div>
</template>

<script setup lang="ts">
const { t, locale } = useI18n()
const localePath = useLocalePath()
const { listArticles, getCategories, getTags } = useApi()
const route = useRoute()
const router = useRouter()

// ── Filters from URL ──────────────────────────────────────────────
const activeCategory = ref((route.query.category as string) || '')
const activeTag = ref((route.query.tag as string) || '')

// ── Load taxonomy ─────────────────────────────────────────────────
const [{ data: categories }, { data: tags }] = await Promise.all([
  useAsyncData('article-page-categories', getCategories, {
    default: () => [],
  }),
  useAsyncData('article-page-tags', getTags, {
    default: () => [],
  }),
])

// ── Load articles (Infinite Scroll) ───────────────────────────────
const page = ref(1)
const allArticles = ref<any[]>([])
const total = ref(0)
const loadingMore = ref(false)

const fetchParams = computed(() => ({
  page: page.value,
  pageSize: 9,
  category: activeCategory.value || undefined,
  tag: activeTag.value || undefined,
}))

// SSR Initial load (force page 1)
const { data: initialData, pending } = await useAsyncData(
  'articles-list',
  () => listArticles({ ...fetchParams.value, page: 1 }),
  { getCachedData: () => undefined }
)

// Sync SSR data to local refs
if (initialData.value) {
  allArticles.value = initialData.value.rows || []
  total.value = initialData.value.total || 0
}

const hasMore = computed(() => allArticles.value.length < total.value)

async function loadMore() {
  if (loadingMore.value || !hasMore.value) return
  loadingMore.value = true
  page.value++
  try {
    const res = await listArticles(fetchParams.value)
    if (res.rows?.length) {
      allArticles.value.push(...res.rows)
    }
  } catch (err) {
    page.value--
    console.error(err)
  } finally {
    loadingMore.value = false
  }
}

// ── Sync URL to filters & Reset List ──────────────────────────────
watch([activeCategory, activeTag], async () => {
  page.value = 1
  pending.value = true
  try {
    const res = await listArticles(fetchParams.value)
    allArticles.value = res.rows || []
    total.value = res.total || 0
    
    router.replace({
      query: {
        ...(activeCategory.value ? { category: activeCategory.value } : {}),
        ...(activeTag.value ? { tag: activeTag.value } : {}),
      }
    })
  } finally {
    pending.value = false
  }
})

// ── Sync URL changes back to filters ──────────────────────────────
watch(() => route.query, (query) => {
  if (query.category !== activeCategory.value || query.tag !== activeTag.value) {
    activeCategory.value = (query.category as string) || ''
    activeTag.value = (query.tag as string) || ''
  }
}, { deep: true })

function setCategory(slug: string) {
  activeCategory.value = activeCategory.value === slug ? '' : slug
  activeTag.value = ''
  page.value = 1
}

function setTag(slug: string) {
  activeTag.value = activeTag.value === slug ? '' : slug
  activeCategory.value = ''
  page.value = 1
}

function clearFilters() {
  activeCategory.value = ''
  activeTag.value = ''
  page.value = 1
}

// ── Lookup helpers ──────────────────────────────────────────────
function getCategoryName(slug: string) {
  return categories.value?.find(c => c.slug === slug)?.name ?? slug
}
function getTagName(slug: string) {
  return tags.value?.find(t => t.slug === slug)?.name ?? slug
}

// ── Date formatter (i18n-aware) ────────────────────────────────
function formatDate(dateStr: string) {
  const d = new Date(dateStr)
  // Map nuxt/i18n locale codes to BCP 47 locale tags
  const bcp47 = locale.value === 'zh' ? 'zh-CN' : 'en-US'
  return new Intl.DateTimeFormat(bcp47, {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    weekday: 'long',
  }).format(d)
}

// ── SEO ──────────────────────────────────────────────────────────
useSeoMeta({ title: t('nav.articles') })
</script>

<style scoped>
.articles-page {
  padding-top: calc(var(--header-h) + 48px);
  padding-bottom: 80px;
}

/* Page hero */
.page-hero { margin-bottom: 36px; }
.page-title {
  font-size: clamp(1.5rem, 3vw, 2rem);
  font-weight: 800;
  margin-bottom: 6px;
}
.page-subtitle { font-size: 14px; color: var(--text-muted); }

/* Filters */
.filter-bar {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  margin-bottom: 12px;
}

.filter-pill {
  padding: 6px 16px;
  border-radius: 99px;
  border: 1px solid var(--border);
  font-size: 13px;
  font-weight: 500;
  color: var(--text-muted);
  background: var(--surface);
  transition: all var(--dur-fast);
  cursor: pointer;
}

.filter-pill:hover, .filter-pill.active {
  background: var(--accent-soft);
  border-color: var(--accent);
  color: var(--accent-text);
  font-weight: 600;
}

/* Tags */
.tag-cloud {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
  margin-bottom: 32px;
}

.tag-pill {
  padding: 4px 12px;
  border-radius: 99px;
  border: 1px solid transparent;
  font-size: 12px;
  color: var(--text-subtle);
  background: var(--bg-secondary);
  transition: all var(--dur-fast);
  cursor: pointer;
}

.tag-pill:hover { color: var(--text-muted); border-color: var(--border); }
.tag-pill.active { color: var(--accent-text); background: var(--accent-soft); border-color: var(--accent); }

/* ── Article List ─────────────────────────────────────────────── */
.article-list {
  display: flex;
  flex-direction: column;
  margin-bottom: 40px;
}

.article-item {
  display: flex;
  align-items: center;
  gap: 20px;
  padding: 22px 0;
  text-decoration: none;
  background: transparent;
  border-radius: 0;
  position: relative;
  animation: slide-up-spring 0.65s cubic-bezier(0.34, 1.4, 0.64, 1) both;
  transition: transform 0.35s cubic-bezier(0.34, 1.5, 0.64, 1);
  cursor: pointer;
  overflow: hidden;
}

.article-item:hover {
  transform: translateX(10px);
}

/* Index number on the left */
.article-index {
  font-size: 11px;
  font-weight: 700;
  letter-spacing: 0.05em;
  color: var(--text-subtle);
  opacity: 0.45;
  min-width: 26px;
  font-variant-numeric: tabular-nums;
  transition: opacity 0.25s, color 0.25s;
  flex-shrink: 0;
  margin-top: 2px;
  align-self: flex-start;
}

.article-item:hover .article-index {
  opacity: 1;
  color: var(--accent);
}

/* Title + meta block */
.article-body {
  flex: 1;
  min-width: 0;
}

.article-title {
  font-size: clamp(1.1rem, 2.2vw, 1.5rem);
  font-weight: 700;
  color: var(--text-primary);
  margin-bottom: 10px;
  line-height: 1.35;
  transition: color 0.2s ease;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.article-item:hover .article-title {
  color: var(--accent);
}

/* Meta row */
.article-meta {
  display: flex;
  align-items: center;
  flex-wrap: wrap;
  gap: 14px;
  font-size: 13px;
  color: var(--text-muted);
}

.meta-item {
  display: inline-flex;
  align-items: center;
  gap: 5px;
  transition: color 0.2s;
}

.meta-icon {
  width: 13.5px;
  height: 13.5px;
  opacity: 0.65;
}

/* Clear capsule styling, just interactive text */
.meta-category, .meta-tag {
  cursor: pointer;
}

.meta-category:hover, .meta-tag:hover {
  color: var(--accent);
}

.meta-date {
  font-variant-numeric: tabular-nums;
  letter-spacing: 0.01em;
}

.meta-views {
  font-variant-numeric: tabular-nums;
}

/* Arrow on right */
.article-arrow {
  color: var(--text-subtle);
  opacity: 0;
  font-size: 16px;
  flex-shrink: 0;
  transform: translateX(-6px);
  transition: opacity 0.25s, transform 0.3s cubic-bezier(0.34, 1.5, 0.64, 1), color 0.2s;
}
.article-item:hover .article-arrow {
  opacity: 1;
  transform: translateX(0);
  color: var(--accent);
}

/* Spring in animation */
@keyframes slide-up-spring {
  0%   { opacity: 0; transform: translateY(28px) scale(0.97); }
  100% { opacity: 1; transform: translateY(0)   scale(1); }
}

/* -- Skeleton ----------------------------------------------------------- */
.skeleton-article-item {
  padding: 22px 0;
}
.skeleton-inner { padding-left: 46px; }

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

/* Empty */
.empty-state { text-align: center; padding: 64px 0; color: var(--text-subtle); }

@media (max-width: 640px) {
  .article-item { gap: 12px; }
  .article-item:hover { transform: translateX(4px); }
  .article-title { font-size: 1rem; white-space: normal; }
}
</style>
