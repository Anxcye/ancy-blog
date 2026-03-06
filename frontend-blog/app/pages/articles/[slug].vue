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
        <div v-if="article.categorySlug" class="article-category">{{ article.categorySlug }}</div>
        <h1 class="article-title">{{ article.title }}</h1>

        <div class="article-meta">
          <time class="meta-date" :datetime="article.publishedAt">
            {{ t('article.publishedAt') }} {{ formatDate(article.publishedAt) }}
          </time>
          <span v-if="article.aiAssistLevel && article.aiAssistLevel !== 'none'" class="meta-ai">
            🤖 {{ aiAssistLabel(article.aiAssistLevel) }}
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
const { getArticle, getSiteSettings } = useApi()

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

// ── Helpers ─────────────────────────────────────────────────────
function formatDate(iso?: string): string {
  if (!iso) return ''
  return new Date(iso).toLocaleDateString('zh-CN', { year: 'numeric', month: 'long', day: 'numeric' })
}

const AI_LEVEL_LABEL: Record<string, string> = {
  polish: '文字润色', dictation: '语音速记', assisted: 'AI 辅助', generated: 'AI 生成', translated: 'AI 翻译'
}
function aiAssistLabel(level: string): string { return AI_LEVEL_LABEL[level] || level }

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

.article-category {
  font-size: 11px;
  font-weight: 700;
  color: var(--accent-text);
  text-transform: uppercase;
  letter-spacing: 0.06em;
  margin-bottom: 12px;
}

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
  align-items: center;
  gap: 16px;
  margin-bottom: 16px;
}

.meta-date { font-size: 13px; color: var(--text-subtle); }
.meta-ai {
  font-size: 11px;
  padding: 3px 10px;
  border-radius: 99px;
  background: var(--accent-soft);
  color: var(--accent-text);
  font-weight: 600;
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
