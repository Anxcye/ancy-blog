<!-- File: app/components/ArticleCard.vue
     Purpose: Reusable article card for list and home page.
     Module: frontend-blog/components, presentation layer.
     Related: pages/index.vue, pages/articles/index.vue. -->
<template>
  <article class="article-card" :class="{ 'article-card--featured': featured }">
    <!-- Cover image -->
    <NuxtLink v-if="article.coverImage && featured" :to="articlePath" class="card-cover-link">
      <img :src="article.coverImage" :alt="article.title" class="card-cover" loading="lazy" />
    </NuxtLink>

    <div class="card-body">
      <!-- Meta: category + date -->
      <div class="card-meta">
        <span v-if="article.categorySlug" class="card-category">
          {{ article.categorySlug }}
        </span>
        <time class="card-date" :datetime="article.publishedAt">
          {{ formatDate(article.publishedAt || article.createdAt) }}
        </time>
        <span v-if="article.isPinned" class="card-pin" title="置顶">📌</span>
      </div>

      <!-- Title -->
      <NuxtLink :to="articlePath" class="card-title-link">
        <h2 class="card-title">{{ article.title }}</h2>
      </NuxtLink>

      <!-- Summary -->
      <p v-if="article.summary" class="card-summary">{{ article.summary }}</p>

      <!-- Footer: tags + read more -->
      <div class="card-footer">
        <div class="card-tags">
          <span v-for="tag in (article.tagSlugs || []).slice(0, 3)" :key="tag" class="card-tag">
            #{{ tag }}
          </span>
        </div>
        <NuxtLink :to="articlePath" class="card-readmore">
          {{ t('home.readMore') }}
          <svg viewBox="0 0 16 16" fill="none" stroke="currentColor" stroke-width="1.5">
            <path d="M3 8h10M9 4l4 4-4 4"/>
          </svg>
        </NuxtLink>
      </div>
    </div>
  </article>
</template>

<script setup lang="ts">
import type { ArticleCard } from '~/composables/useApi'

const props = defineProps<{
  article: ArticleCard
  featured?: boolean
}>()

const { t } = useI18n()
const localePath = useLocalePath()

const articlePath = computed(() =>
  localePath(`/articles/${props.article.slug}`)
)

function formatDate(iso?: string): string {
  if (!iso) return ''
  return new Date(iso).toLocaleDateString('zh-CN', {
    year: 'numeric', month: 'long', day: 'numeric'
  })
}
</script>

<style scoped>
.article-card {
  background: var(--surface);
  border: 1px solid var(--border);
  border-radius: var(--radius-lg);
  overflow: hidden;
  transition: border-color var(--dur-base), box-shadow var(--dur-base), transform var(--dur-base) var(--ease-spring);
}

.article-card:hover {
  border-color: var(--border-strong);
  box-shadow: var(--shadow-md);
  transform: translateY(-2px);
}

.card-cover-link { display: block; overflow: hidden; }
.card-cover {
  width: 100%;
  height: 220px;
  object-fit: cover;
  transition: transform var(--dur-slow) var(--ease-smooth);
}
.article-card:hover .card-cover { transform: scale(1.03); }

.card-body { padding: 20px 24px 20px; }

.card-meta {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-bottom: 10px;
  font-size: 12px;
}

.card-category {
  color: var(--accent-text);
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.04em;
  font-size: 11px;
}

.card-date { color: var(--text-subtle); }
.card-pin { font-size: 12px; }

.card-title-link { display: block; }
.card-title {
  font-size: 1.15rem;
  font-weight: 700;
  color: var(--text);
  line-height: 1.4;
  margin-bottom: 8px;
  transition: color var(--dur-fast);
}
.article-card:hover .card-title { color: var(--accent-text); }

.card-summary {
  font-size: 14px;
  color: var(--text-muted);
  line-height: 1.6;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
  margin-bottom: 16px;
}

.card-footer {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 8px;
}

.card-tags {
  display: flex;
  gap: 6px;
  flex-wrap: wrap;
}

.card-tag {
  font-size: 11px;
  color: var(--text-subtle);
  background: var(--bg-secondary);
  padding: 2px 8px;
  border-radius: 99px;
  border: 1px solid var(--border);
}

.card-readmore {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 13px;
  font-weight: 600;
  color: var(--accent-text);
  white-space: nowrap;
  transition: gap var(--dur-fast);
}

.card-readmore:hover { gap: 8px; }

.card-readmore svg {
  width: 14px;
  height: 14px;
  transition: transform var(--dur-fast);
}

.card-readmore:hover svg { transform: translateX(2px); }

/* Featured variant */
.article-card--featured {
  grid-column: 1 / -1;
}

.article-card--featured .card-title { font-size: 1.35rem; }
</style>
