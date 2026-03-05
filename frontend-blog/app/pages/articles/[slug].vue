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
        <h2 class="comments-title">
          {{ t('article.comments') }}
          <span v-if="commentTotal" class="comments-count">{{ commentTotal.total }}</span>
        </h2>

        <!-- Comment list -->
        <div v-if="comments?.rows?.length" class="comment-list">
          <div
            v-for="comment in comments.rows"
            :key="comment.id"
            class="comment-item"
          >
            <div class="comment-avatar">
              <img v-if="comment.avatarUrl" :src="comment.avatarUrl" :alt="comment.nickname" />
              <span v-else>{{ comment.nickname.charAt(0).toUpperCase() }}</span>
            </div>
            <div class="comment-content">
              <div class="comment-header">
                <a
                  v-if="comment.website"
                  :href="comment.website"
                  target="_blank"
                  rel="noopener noreferrer nofollow"
                  class="comment-name"
                >{{ comment.nickname }}</a>
                <span v-else class="comment-name">{{ comment.nickname }}</span>
                <time class="comment-date">{{ formatDateShort(comment.createdAt) }}</time>
              </div>
              <p class="comment-text">{{ comment.content }}</p>
            </div>
          </div>
        </div>

        <p v-else class="comments-empty">{{ t('article.noComments') }}</p>

        <!-- Comment form -->
        <form class="comment-form" @submit.prevent="submitComment">
          <h3 class="form-title">{{ t('article.leaveComment') }}</h3>

          <div class="form-row">
            <div class="form-field">
              <label>{{ t('comment.nickname') }} *</label>
              <input v-model="commentForm.nickname" type="text" required :placeholder="t('comment.nickname')" />
            </div>
            <div class="form-field">
              <label>{{ t('comment.email') }}</label>
              <input v-model="commentForm.email" type="email" :placeholder="t('comment.email')" />
            </div>
          </div>

          <div class="form-field">
            <label>{{ t('comment.website') }}</label>
            <input v-model="commentForm.website" type="url" :placeholder="t('comment.website')" />
          </div>

          <div class="form-field">
            <label>{{ t('article.comments') }} *</label>
            <textarea
              v-model="commentForm.content"
              required
              rows="4"
              :placeholder="t('comment.content')"
            />
          </div>

          <div class="form-actions">
            <button type="submit" class="submit-btn" :disabled="submitting">
              {{ submitting ? t('comment.submitting') : t('comment.submit') }}
            </button>
          </div>
        </form>
      </section>

    </div>
  </div>
</template>

<script setup lang="ts">
const { t } = useI18n()
const localePath = useLocalePath()
const route = useRoute()
const { getArticle, listComments, getCommentTotal, createComment, getSiteSettings } = useApi()

const slug = computed(() => route.params.slug as string)

// ── Fetch article ───────────────────────────────────────────────
const { data: article, error } = await useAsyncData(
  `article-${slug.value}`,
  () => getArticle(slug.value)
)

if (error.value || !article.value) {
  throw createError({ statusCode: 404, message: 'Article not found' })
}

// ── Fetch site settings, comments, total ────────────────────────
const [{ data: siteSettings }, { data: comments }, { data: commentTotal }] = await Promise.all([
  useAsyncData('article-site-settings', getSiteSettings),
  useAsyncData(`comments-${article.value.id}`, () => listComments(article.value!.id, { pageSize: 50 })),
  useAsyncData(`comment-total-${article.value.id}`, () => getCommentTotal(article.value!.id)),
])

// ── Helpers ─────────────────────────────────────────────────────
function formatDate(iso?: string): string {
  if (!iso) return ''
  return new Date(iso).toLocaleDateString('zh-CN', { year: 'numeric', month: 'long', day: 'numeric' })
}

function formatDateShort(iso: string): string {
  return new Date(iso).toLocaleDateString('zh-CN', { month: 'numeric', day: 'numeric' })
}

const AI_LEVEL_LABEL: Record<string, string> = {
  polish: '文字润色', dictation: '语音速记', assisted: 'AI 辅助', generated: 'AI 生成', translated: 'AI 翻译'
}
function aiAssistLabel(level: string): string { return AI_LEVEL_LABEL[level] || level }

// ── Comment form ─────────────────────────────────────────────────
const toast = useToast()
const commentForm = reactive({ nickname: '', email: '', website: '', content: '' })
const submitting = ref(false)

const doSubmitComment = async () => {
  if (!article.value || !commentForm.content.trim()) return
  submitting.value = true
  
  try {
    await createComment({
      articleId: article.value.id,
      content: commentForm.content,
      nickname: commentForm.nickname,
      email: commentForm.email || undefined,
      website: commentForm.website || undefined,
    })
    toast.add({ title: t('comment.success'), color: 'green', icon: 'i-heroicons-check-circle' })
    commentForm.content = ''
  } catch {
    toast.add({ title: '提交失败，请稍后重试', color: 'red', icon: 'i-heroicons-x-circle' })
  } finally {
    submitting.value = false
  }
}

// Enterprise Debounce: prevents spamming the submit button
const submitComment = useDebounceFn(doSubmitComment, 500)

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

.comments-title {
  font-size: 1.2rem;
  font-weight: 700;
  margin-bottom: 28px;
  display: flex;
  align-items: center;
  gap: 10px;
}

.comments-count {
  font-size: 13px;
  font-weight: 500;
  color: var(--text-subtle);
  background: var(--bg-secondary);
  padding: 2px 10px;
  border-radius: 99px;
}

.comment-list { display: flex; flex-direction: column; gap: 20px; margin-bottom: 40px; }

.comment-item {
  display: flex;
  gap: 14px;
}

.comment-avatar {
  width: 36px;
  height: 36px;
  border-radius: 50%;
  background: var(--accent-soft);
  flex-shrink: 0;
  overflow: hidden;
  display: grid;
  place-items: center;
  font-weight: 700;
  font-size: 15px;
  color: var(--accent-text);
  border: 1.5px solid var(--border);
}

.comment-avatar img { width: 100%; height: 100%; object-fit: cover; }

.comment-content { flex: 1; min-width: 0; }

.comment-header {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-bottom: 6px;
}

.comment-name {
  font-size: 14px;
  font-weight: 600;
  color: var(--text);
}

a.comment-name { color: var(--accent-text); }
a.comment-name:hover { text-decoration: underline; }

.comment-date { font-size: 12px; color: var(--text-subtle); }

.comment-text {
  font-size: 14px;
  line-height: 1.65;
  color: var(--text-muted);
  white-space: pre-wrap;
  word-break: break-word;
}

.comments-empty { color: var(--text-subtle); font-size: 14px; margin-bottom: 40px; }

/* ── Comment form ── */
.comment-form {
  background: var(--surface);
  border: 1px solid var(--border);
  border-radius: var(--radius-lg);
  padding: 28px;
}

.form-title {
  font-size: 1rem;
  font-weight: 700;
  margin-bottom: 20px;
}

.form-row {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 16px;
}

.form-field {
  display: flex;
  flex-direction: column;
  gap: 6px;
  margin-bottom: 16px;
}

.form-field label {
  font-size: 12px;
  font-weight: 600;
  color: var(--text-muted);
  text-transform: uppercase;
  letter-spacing: 0.04em;
}

.form-field input,
.form-field textarea {
  border: 1px solid var(--border);
  border-radius: var(--radius-md);
  padding: 10px 14px;
  font-size: 14px;
  font-family: inherit;
  color: var(--text);
  background: var(--bg);
  transition: border-color var(--dur-fast), box-shadow var(--dur-fast);
  resize: vertical;
}

.form-field input:focus,
.form-field textarea:focus {
  outline: none;
  border-color: var(--accent);
  box-shadow: 0 0 0 3px var(--accent-soft);
}

.form-actions { display: flex; align-items: center; gap: 16px; margin-top: 4px; }

.submit-btn {
  padding: 10px 24px;
  background: var(--accent);
  color: white;
  border-radius: var(--radius-md);
  font-size: 14px;
  font-weight: 600;
  transition: opacity var(--dur-fast), transform var(--dur-fast) var(--ease-spring);
}

.submit-btn:hover:not(:disabled) { opacity: 0.88; transform: translateY(-1px); }
.submit-btn:disabled { opacity: 0.55; cursor: not-allowed; }

.submit-result { font-size: 13px; }
.submit-result.success { color: #16a34a; }
.submit-result.error { color: #dc2626; }

@media (max-width: 640px) {
  .form-row { grid-template-columns: 1fr; }
}
</style>
