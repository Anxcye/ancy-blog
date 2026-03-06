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
          <div class="skeleton-line" style="height: 28px; width: 60%; margin-bottom: 12px;" />
          <div class="skeleton-line" style="height: 16px; width: 30%;" />
        </div>
      </div>

      <div v-else-if="articles?.rows?.length" class="article-list">
        <NuxtLink
          v-for="(article, i) in articles.rows"
          :key="article.id"
          :to="localePath(`/articles/${article.slug}`)"
          class="horizontal-article-item"
          :style="{ animationDelay: `${i * 80}ms` }"
        >
          <h2 class="article-item-title">{{ article.title }}</h2>
          <div class="article-item-meta">
            <time class="meta-date">{{ new Date(article.publishedAt || article.createdAt).toLocaleDateString() }}</time>
          </div>
        </NuxtLink>
      </div>

      <div v-else class="empty-state">
        <p>{{ t('home.noArticles') }}</p>
      </div>

      <!-- Pagination -->
      <div v-if="totalPages > 1" class="pagination">
        <button
          class="page-btn"
          :disabled="page <= 1"
          @click="page--"
        >←</button>

        <span class="page-info">{{ page }} / {{ totalPages }}</span>

        <button
          class="page-btn"
          :disabled="page >= totalPages"
          @click="page++"
        >→</button>
      </div>

    </div>
  </div>
</template>

<script setup lang="ts">
const { t } = useI18n()
const localePath = useLocalePath()
const { listArticles, getCategories, getTags } = useApi()
const route = useRoute()
const router = useRouter()

// ── Filters from URL ──────────────────────────────────────────────
const page = ref(Number(route.query.page) || 1)
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

// ── Load articles (reactive to filters) ───────────────────────────
const { data: articles, pending, refresh } = await useAsyncData(
  'articles-list',
  () => listArticles({
    page: page.value,
    pageSize: 9,
    category: activeCategory.value || undefined,
    tag: activeTag.value || undefined,
  }),
  {
    watch: [page, activeCategory, activeTag],
    getCachedData: () => undefined,  // always re-fetch; prevents stale SSR payload on CSR nav
  }
)

const totalPages = computed(() =>
  articles.value ? Math.ceil(articles.value.total / 9) : 1
)

// ── Sync URL to filters ───────────────────────────────────────────
watch([page, activeCategory, activeTag], () => {
  router.replace({
    query: {
      ...(page.value > 1 ? { page: page.value } : {}),
      ...(activeCategory.value ? { category: activeCategory.value } : {}),
      ...(activeTag.value ? { tag: activeTag.value } : {}),
    }
  })
  window.scrollTo({ top: 0, behavior: 'smooth' })
})

// ── Sync URL changes back to filters ──────────────────────────────
watch(() => route.query, (query) => {
  page.value = Number(query.page) || 1
  activeCategory.value = (query.category as string) || ''
  activeTag.value = (query.tag as string) || ''
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

/* Article List */
.article-list {
  display: flex;
  flex-direction: column;
  gap: 16px;
  margin-bottom: 40px;
}

.horizontal-article-item {
  display: block;
  padding: 24px 0;
  border-bottom: 1px solid var(--border);
  text-decoration: none;
  background: transparent;
  transition: all 0.4s cubic-bezier(0.34, 1.56, 0.64, 1);
  animation: slide-up-spring 0.8s cubic-bezier(0.34, 1.56, 0.64, 1) both;
  cursor: pointer;
  position: relative;
}

.horizontal-article-item:last-child {
  border-bottom: none;
}

.horizontal-article-item:hover {
  transform: translateX(12px);
}

.horizontal-article-item::before {
  content: '';
  position: absolute;
  left: -16px;
  top: 50%;
  transform: translateY(-50%) scale(0);
  width: 4px;
  height: 0;
  background: var(--accent);
  border-radius: 4px;
  transition: all 0.3s cubic-bezier(0.34, 1.56, 0.64, 1);
  opacity: 0;
}

.horizontal-article-item:hover::before {
  transform: translateY(-50%) scale(1);
  height: 40%;
  opacity: 1;
}

.article-item-title {
  font-size: clamp(1.25rem, 2.5vw, 1.75rem);
  font-weight: 700;
  color: var(--text-primary);
  margin-bottom: 8px;
  line-height: 1.3;
  transition: color var(--dur-fast);
}

.horizontal-article-item:hover .article-item-title {
  color: var(--accent);
}

.article-item-meta {
  display: flex;
  gap: 16px;
  font-size: 14px;
  color: var(--text-muted);
}

@keyframes slide-up-spring {
  0% { opacity: 0; transform: translateY(30px) scale(0.98); }
  100% { opacity: 1; transform: translateY(0) scale(1); }
}

/* Skeleton */
.skeleton-article-item {
  padding: 24px 0;
  border-bottom: 1px solid var(--border);
}

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
.pagination {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 20px;
  margin-top: 40px;
}

.page-btn {
  width: 38px;
  height: 38px;
  border-radius: var(--radius-md);
  border: 1px solid var(--border);
  background: var(--surface);
  color: var(--text-muted);
  font-size: 16px;
  display: grid;
  place-items: center;
  transition: all var(--dur-fast);
  cursor: pointer;
}

.page-btn:hover:not(:disabled) {
  border-color: var(--accent);
  color: var(--accent-text);
  background: var(--accent-soft);
}

.page-btn:disabled { opacity: 0.35; cursor: not-allowed; }

.page-info { font-size: 14px; color: var(--text-muted); min-width: 60px; text-align: center; }

/* Empty */
.empty-state { text-align: center; padding: 64px 0; color: var(--text-subtle); }

@media (max-width: 640px) {
  .horizontal-article-item:hover {
    transform: translateX(4px);
  }
}
</style>
