<!-- File: app/pages/articles/[slug].vue
     Purpose: Article detail page — TipTap JSON renderer, comment section.
     Module: frontend-blog/pages/articles, presentation layer.
     Related: composables/useApi.ts, components/ArticleContent.vue (inline renderer). -->
<template>
  <div class="article-page">
    <div class="container">

      <!-- Back link -->
      <NuxtLink :to="localePath('/articles')" class="back-link">
        <svg viewBox="0 0 16 16" fill="none" stroke="currentColor" stroke-width="1.5">
          <path d="M10 4L6 8l4 4"/>
        </svg>
        {{ t('nav.articles') }}
      </NuxtLink>

      <!-- ── Article header ── -->
      <header v-if="article" class="article-header">
        <h1 class="article-title">{{ article.title }}</h1>

        <div class="article-meta">
          <span class="meta-item">
            <svg class="meta-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true">
              <rect x="3" y="4" width="18" height="18" rx="2"></rect>
              <line x1="16" y1="2" x2="16" y2="6"></line>
              <line x1="8" y1="2" x2="8" y2="6"></line>
              <line x1="3" y1="10" x2="21" y2="10"></line>
            </svg>
            <time :datetime="article.publishedAt">{{ formatDate(article.publishedAt) }}</time>
          </span>
          <span v-if="showUpdatedMeta" class="meta-item meta-updated">
            <svg class="meta-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true">
              <path d="M21 12a9 9 0 1 1-3.2-6.9"></path>
              <path d="M21 3v6h-6"></path>
            </svg>
            <span class="meta-prefix">{{ t('article.revised') }}</span>
            <time :datetime="article.updatedAt">{{ formatDate(article.updatedAt) }}</time>
          </span>
          <span v-if="article.categorySlug" class="meta-item">
            <svg class="meta-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true">
              <path d="M20.59 13.41 11 3.83a2 2 0 0 0-1.41-.58H4a2 2 0 0 0-2 2v5.59a2 2 0 0 0 .58 1.41l9.59 9.59a2 2 0 0 0 2.83 0l5.59-5.59a2 2 0 0 0 0-2.83Z"></path>
              <circle cx="7.5" cy="7.5" r="1.5"></circle>
            </svg>
            <span>{{ categoryName }}</span>
          </span>
          <span
            v-if="article.aiAssistLevel"
            class="meta-item meta-ai"
          >
            <svg class="meta-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true">
              <path d="M12 2a3 3 0 0 0-3 3v1H8a2 2 0 0 0-2 2v4a6 6 0 0 0 12 0V8a2 2 0 0 0-2-2h-1V5a3 3 0 0 0-3-3Z"></path>
              <path d="M9 18h6"></path>
              <path d="M10 22h4"></path>
            </svg>
            <span>{{ aiAssistLabel(article.aiAssistLevel) }}</span>
            <span class="meta-tooltip" role="tooltip">{{ aiAssistDescription(article.aiAssistLevel) }}</span>
          </span>
        </div>

        <!-- Tags -->
        <div v-if="article.tagSlugs?.length" class="article-tags">
          <span v-for="tag in article.tagSlugs" :key="tag" class="article-tag">#{{ tag }}</span>
        </div>

        <!-- Cover image -->
        <img
          v-if="article.coverImage"
          :src="article.coverImage"
          :alt="article.title"
          class="article-cover"
          loading="eager"
        />
      </header>

      <!-- ── Article body ── -->
      <article v-if="article" class="article-body">
        <TiptapRenderer :content="article.content" />
      </article>

      <!-- ── Comments ── -->
      <section v-if="article?.allowComment && siteSettings?.commentEnabled" class="comments-section">
        <CommentList content-type="article" :content-id="article.id" :require-approval="siteSettings?.commentRequireApproval" />
      </section>

    </div>
  </div>
</template>

<script setup lang="ts">
const { t } = useI18n()
const localePath = useLocalePath()
const route = useRoute()
const { getArticle, getSiteSettings, getCategories } = useApi()

const slug = computed(() => route.params.slug as string)

// ── Fetch article ───────────────────────────────────────────────
const { data: article, error } = await useAsyncData(
  `article-${slug.value}`,
  () => getArticle(slug.value)
)

if (error.value || !article.value) {
  throw createError({ statusCode: 404, message: 'Article not found' })
}

// ── Fetch site settings ────────────────────────
const { data: siteSettings } = await useAsyncData('article-site-settings', getSiteSettings)
const { data: categories } = await useAsyncData('article-categories', getCategories)

// ── Helpers ─────────────────────────────────────────────────────
function formatDate(iso?: string): string {
  if (!iso) return ''
  return new Date(iso).toLocaleDateString('zh-CN', { year: 'numeric', month: 'long', day: 'numeric' })
}

const AI_LEVEL_LABEL_KEY: Record<string, string> = {
  none: 'article.aiLevels.none',
  polish: 'article.aiLevels.polish',
  dictation: 'article.aiLevels.dictation',
  assisted: 'article.aiLevels.assisted',
  generated: 'article.aiLevels.generated',
  translated: 'article.aiLevels.translated',
}

function aiAssistLabel(level: string): string {
  const key = AI_LEVEL_LABEL_KEY[level] || AI_LEVEL_LABEL_KEY.none
  return t(key)
}

const AI_LEVEL_DESCRIPTION_KEY: Record<string, string> = {
  none: 'article.aiDescriptions.none',
  polish: 'article.aiDescriptions.polish',
  dictation: 'article.aiDescriptions.dictation',
  assisted: 'article.aiDescriptions.assisted',
  generated: 'article.aiDescriptions.generated',
  translated: 'article.aiDescriptions.translated',
}

function aiAssistDescription(level?: string): string {
  if (!level) return ''
  const key = AI_LEVEL_DESCRIPTION_KEY[level] || AI_LEVEL_DESCRIPTION_KEY.none
  return t(key)
}

const showUpdatedMeta = computed(() => {
  if (!article.value?.publishedAt || !article.value?.updatedAt) return false
  const publishedAt = new Date(article.value.publishedAt).getTime()
  const updatedAt = new Date(article.value.updatedAt).getTime()
  return updatedAt - publishedAt > 60_000
})

const categoryName = computed(() => {
  const slug = article.value?.categorySlug
  if (!slug) return ''
  return categories.value?.find((item) => item.slug === slug)?.name || slug
})

// ── SEO & JSON-LD ───────────────────────────────────────────────
useArticleSeo(article.value, siteSettings.value || null)
</script>

<style scoped>
.article-page {
  padding-top: calc(var(--header-h) + 40px);
  padding-bottom: 80px;
}

.back-link {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  font-size: 13px;
  color: var(--text-subtle);
  margin-bottom: 32px;
  transition: color var(--dur-fast);
}
.back-link:hover { color: var(--accent-text); }
.back-link svg { width: 16px; height: 16px; }

/* ── Header ── */
.article-header { margin-bottom: 40px; }

.article-title {
  font-family: 'Songti SC', 'SimSun', 'Noto Serif SC', Georgia, serif;
  font-size: clamp(1.6rem, 4vw, 2.2rem);
  font-weight: 800;
  line-height: 1.25;
  margin-bottom: 16px;
  letter-spacing: -0.02em;
}

.article-meta {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  gap: 8px 14px;
  margin-bottom: 16px;
  color: var(--text-subtle);
  font-size: 13px;
}

.meta-item {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  position: relative;
  line-height: 1.4;
}

.meta-item::after {
  content: '·';
  margin-left: 8px;
  color: var(--border-strong, rgba(0, 0, 0, 0.25));
}

.meta-item:last-child::after {
  display: none;
}

.meta-icon {
  width: 14px;
  height: 14px;
  flex-shrink: 0;
  color: var(--text-muted);
}

.meta-prefix {
  color: var(--text-muted);
  font-weight: 500;
}

.meta-ai {
  color: var(--accent-text);
  cursor: help;
}

.meta-tooltip {
  position: absolute;
  left: 0;
  top: calc(100% + 10px);
  z-index: 20;
  width: min(280px, 70vw);
  padding: 10px 12px;
  border-radius: 12px;
  background: rgba(17, 17, 17, 0.94);
  color: #fff;
  font-size: 12px;
  line-height: 1.5;
  opacity: 0;
  transform: translateY(4px);
  pointer-events: none;
  transition: opacity 180ms ease, transform 180ms ease;
  box-shadow: 0 12px 30px rgba(0, 0, 0, 0.18);
}

.meta-tooltip::before {
  content: '';
  position: absolute;
  left: 18px;
  bottom: 100%;
  width: 10px;
  height: 10px;
  background: rgba(17, 17, 17, 0.94);
  transform: rotate(45deg) translateY(6px);
}

.meta-ai:hover .meta-tooltip,
.meta-ai:focus-within .meta-tooltip {
  opacity: 1;
  transform: translateY(0);
}

@media (max-width: 720px) {
  .article-meta {
    gap: 8px;
  }

  .meta-item {
    width: auto;
  }

  .meta-tooltip {
    width: min(260px, calc(100vw - 48px));
  }
}

.article-tags { display: flex; flex-wrap: wrap; gap: 6px; margin-bottom: 28px; }
.article-tag {
  font-size: 12px;
  color: var(--text-subtle);
  background: var(--bg-secondary);
  padding: 3px 10px;
  border-radius: 99px;
  border: 1px solid var(--border);
}

.article-cover {
  width: 100%;
  max-height: 420px;
  object-fit: cover;
  border-radius: var(--radius-lg);
  margin-top: 8px;
  box-shadow: var(--shadow-md);
}

/* ── Article body (rich text) ── */
.article-body {
  font-family: 'Songti SC', 'SimSun', 'Noto Serif SC', Georgia, serif;
  font-size: 1.05rem;
  margin-bottom: 64px;
}

/* ── Comments ── */
.comments-section {
  border-top: 1px solid var(--border);
  padding-top: 48px;
}
</style>
