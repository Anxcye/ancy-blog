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
      <div v-if="pending" class="article-grid">
        <div v-for="n in 6" :key="n" class="skeleton-card">
          <div class="skeleton-line" style="height: 14px; width: 40%;" />
          <div class="skeleton-line" style="height: 22px; width: 85%; margin-top: 8px;" />
          <div class="skeleton-line" style="height: 14px; width: 65%; margin-top: 6px;" />
        </div>
      </div>

      <div v-else-if="articles?.rows?.length" class="article-grid">
        <ArticleCard
          v-for="(article, i) in articles.rows"
          :key="article.id"
          :article="article"
          class="grid-item"
          :style="{ animationDelay: `${i * 50}ms` }"
        />
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
const { listArticles, getCategories, getTags } = useApi()
const route = useRoute()
const router = useRouter()

// ── Filters from URL ──────────────────────────────────────────────
const page = ref(Number(route.query.page) || 1)
const activeCategory = ref((route.query.category as string) || '')
const activeTag = ref((route.query.tag as string) || '')

// ── Load taxonomy ─────────────────────────────────────────────────
const [{ data: categories }, { data: tags }] = await Promise.all([
  useAsyncData('categories', getCategories, { default: () => [] }),
  useAsyncData('tags', getTags, { default: () => [] }),
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
  { watch: [page, activeCategory, activeTag] }
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

/* Article grid — same as homepage */
.article-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 20px;
  margin-bottom: 40px;
}

.grid-item { animation: card-in 0.45s var(--ease-spring) both; }
@keyframes card-in {
  from { opacity: 0; transform: translateY(18px); }
  to   { opacity: 1; transform: none; }
}

/* Skeleton */
.skeleton-card {
  background: var(--surface);
  border: 1px solid var(--border);
  border-radius: var(--radius-lg);
  padding: 20px 24px;
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
  .article-grid { grid-template-columns: 1fr; }
}
</style>
