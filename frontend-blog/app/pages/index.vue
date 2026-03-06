<!-- File: app/pages/index.vue
     Purpose: Homepage — full-height hero + recent article grid.
     Module: frontend-blog/pages, presentation layer.
     Related: layouts/default.vue, components/ArticleCard.vue, composables/useApi.ts. -->
<template>
  <div class="home">

    <!-- ═══ Hero ════════════════════════════════════════════════════ -->
    <section class="hero">
      <div class="hero-inner container">
        <!-- Main title with embedded avatar -->
        <h1 class="hero-title">
          <span class="title-line">
            <span v-for="(char, i) in 'Hi, I\'m'.split('')" :key="`hi-${i}`" class="char char-muted" :style="{ animationDelay: `${i * 50}ms` }">{{ char === ' ' ? '\u00A0' : char }}</span>
            <img
              v-if="siteSettings?.avatarUrl"
              :src="siteSettings.avatarUrl"
              :alt="siteSettings?.siteName"
              class="title-avatar"
              :style="{ animationDelay: '400ms' }"
            />
            <span v-for="(char, i) in (siteSettings?.siteName || 'Ancy').split('')" :key="`name-${i}`" class="char char-name" :style="{ animationDelay: `${(i + 9) * 50}ms` }">{{ char }}</span>
          </span>
          <span class="title-line">
            <span v-for="(char, i) in (siteSettings?.heroIntroMd || t('home.heroSubtitle')).split('')" :key="`intro-${i}`" class="char char-gradient" :style="{ animationDelay: `${i * 50}ms` }">{{ char === ' ' ? '\u00A0' : char }}</span>
          </span>
        </h1>

        <!-- Floating social links -->
        <div v-if="socialLinks.length" class="hero-socials">
          <a
            v-for="link in socialLinks"
            :key="link.id"
            :href="link.url"
            target="_blank"
            class="social-icon"
            :title="link.title"
          >
            <img v-if="link.iconKey" :src="link.iconKey" :alt="link.title" class="social-icon-img" />
          </a>
        </div>
      </div>

      <!-- Down arrow -->
      <div class="hero-arrow">
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round">
          <path d="M12 5v14M5 12l7 7 7-7"/>
        </svg>
      </div>
    </section>

    <!-- ═══ Recent articles ══════════════════════════════════════════ -->
    <section class="recent-section">
      <div class="container">
        <div class="section-header">
          <h2 class="section-title">{{ t('home.recentArticles') }}</h2>
          <NuxtLink :to="localePath('/articles')" class="section-more">
            {{ t('home.viewAll') }}
            <svg viewBox="0 0 16 16" fill="none" stroke="currentColor" stroke-width="1.5">
              <path d="M3 8h10M9 4l4 4-4 4"/>
            </svg>
          </NuxtLink>
        </div>

        <!-- Loading skeletons -->
        <div v-if="pending" class="article-grid">
          <div v-for="n in 6" :key="n" class="skeleton-card">
            <div class="skeleton-line" style="height: 14px; width: 40%;" />
            <div class="skeleton-line" style="height: 22px; width: 85%; margin-top: 8px;" />
            <div class="skeleton-line" style="height: 14px; width: 65%; margin-top: 6px;" />
          </div>
        </div>

        <!-- Articles -->
        <div v-else-if="articles?.rows?.length" class="article-grid">
          <ArticleCard
            v-for="(article, i) in articles.rows"
            :key="article.id"
            :article="article"
            :featured="i === 0 && !!article.coverImage"
            class="grid-item"
            :style="{ animationDelay: `${i * 60}ms` }"
          />
        </div>

        <!-- Empty -->
        <div v-else class="empty-state">
          <p>{{ t('home.noArticles') }}</p>
        </div>
      </div>
    </section>
  </div>
</template>

<script setup lang="ts">
const { t } = useI18n()
const localePath = useLocalePath()
const { listArticles, getSiteSettings, getSocialLinks } = useApi()

// ── Fetch data ───────────────────────────────────────────────────
const [{ data: siteSettings }, { data: socialLinks }, { data: articles, pending }] = await Promise.all([
  useAsyncData('site-settings', getSiteSettings),
  useAsyncData('social-links', getSocialLinks, { default: () => [] }),
  useAsyncData('home-articles', () => listArticles({ pageSize: 6 })),
])

// ── SEO ────────────────────────────────────────────────────────────
useSeoMeta({
  title: siteSettings.value?.siteName || 'Ancy Blog',
  description: siteSettings.value?.siteDescription || '',
  ogTitle: siteSettings.value?.siteName,
  ogDescription: siteSettings.value?.siteDescription,
})
</script>

<style scoped>
/* ═══ Hero ═══════════════════════════════════════════════════════ */
.hero {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 120px 0 80px;
  position: relative;
}

.hero-inner {
  text-align: center;
  position: relative;
}

/* ── Title ── */
.hero-title {
  font-size: clamp(3rem, 8vw, 5rem);
  font-weight: 800;
  line-height: 1.1;
  margin: 0;
}

.char {
  display: inline-block;
  opacity: 0;
  transform: translateY(20px) scale(0.8);
  animation: char-spring 0.6s cubic-bezier(0.34, 1.56, 0.64, 1) forwards;
}

.char-muted { color: var(--text-muted); font-weight: 600; }
.char-name  { color: var(--text); }
.char-gradient {
  background: linear-gradient(135deg, var(--accent), var(--accent-text));
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

@keyframes char-spring {
  0% { opacity: 0; transform: translateY(20px) scale(0.8); }
  60% { opacity: 1; transform: translateY(-4px) scale(1.05); }
  100% { opacity: 1; transform: translateY(0) scale(1); }
}

.title-line {
  display: block;
  margin-bottom: 1em;
}

.title-avatar {
  display: inline-block;
  width: 0.9em;
  height: 0.9em;
  border-radius: 50%;
  object-fit: cover;
  vertical-align: -0.15em;
  margin: 0 0.15em;
  border: 3px solid var(--accent);
  box-shadow: 0 4px 20px rgba(var(--accent-rgb, 31,143,138), 0.3);
  opacity: 0;
  animation: char-spring 0.6s cubic-bezier(0.34, 1.56, 0.64, 1) forwards;
}

.title-gradient {
  background: linear-gradient(135deg, var(--accent), var(--accent-text));
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

/* ── Floating socials ── */
.hero-socials {
  display: flex;
  gap: 16px;
  justify-content: center;
  margin-top: 120px;
}

.social-icon {
  width: 44px;
  height: 44px;
  display: grid;
  place-items: center;
  color: var(--text-muted);
  transition: all 0.3s var(--ease-out);
  opacity: 0;
  animation: char-spring 0.6s cubic-bezier(0.34, 1.56, 0.64, 1) forwards;
  animation-delay: 1s;
}

.social-icon:hover {
  color: var(--accent-text);
  transform: translateY(-4px);
}

.social-icon svg {
  width: 20px;
  height: 20px;
}

.social-icon-img {
  width: 42px;
  height: 42px;
  object-fit: contain;
}

@media (max-width: 640px) {
  .hero {
    min-height: 100vh;
    padding: 80px 0 60px;
  }

  .title-avatar {
    width: 0.8em;
    height: 0.8em;
  }
}

/* ── Down arrow ── */
.hero-arrow {
  position: absolute;
  bottom: 32px;
  left: 50%;
  transform: translateX(-50%);
  opacity: 0;
  animation: arrow-fade 1s 1.2s forwards, float 2.5s 1.2s ease-in-out infinite;
}

@keyframes arrow-fade {
  to { opacity: 0.4; }
}

.hero-arrow svg { width: 24px; height: 24px; color: var(--text-muted); }

@keyframes float {
  0%, 100% { transform: translateX(-50%) translateY(0); }
  50%       { transform: translateX(-50%) translateY(6px); }
}

/* ═══ Recent articles ═════════════════════════════════════════════ */
.recent-section { padding: 64px 0 80px; }

.section-header {
  display: flex;
  align-items: baseline;
  justify-content: space-between;
  margin-bottom: 32px;
}

.section-title {
  font-size: 1.25rem;
  font-weight: 700;
  position: relative;
}

.section-title::after {
  content: '';
  display: block;
  width: 32px;
  height: 3px;
  background: var(--accent);
  border-radius: 2px;
  margin-top: 6px;
}

.section-more {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 13px;
  font-weight: 600;
  color: var(--accent-text);
  transition: gap var(--dur-fast);
}

.section-more:hover { gap: 8px; }
.section-more svg { width: 14px; height: 14px; }

/* ── Article grid ── */
.article-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 20px;
}

.grid-item {
  animation: card-in 0.5s var(--ease-spring) both;
}

@keyframes card-in {
  from { opacity: 0; transform: translateY(20px); }
  to   { opacity: 1; transform: none; }
}

/* ── Skeleton ── */
.skeleton-card {
  background: var(--surface);
  border: 1px solid var(--border);
  border-radius: var(--radius-lg);
  padding: 20px 24px;
}

.skeleton-line {
  background: linear-gradient(
    90deg,
    var(--bg-secondary) 25%,
    var(--surface-hover) 50%,
    var(--bg-secondary) 75%
  );
  background-size: 200% 100%;
  border-radius: var(--radius-sm);
  animation: shimmer 1.4s infinite;
}

@keyframes shimmer {
  from { background-position: 200% 0; }
  to   { background-position: -200% 0; }
}

/* ── Empty ── */
.empty-state {
  text-align: center;
  padding: 64px 0;
  color: var(--text-subtle);
  font-size: 15px;
}

/* ── Mobile ── */
@media (max-width: 640px) {
  .hero-inner {
    grid-template-columns: 1fr;
    text-align: center;
    gap: 32px;
  }

  .hero-right {
    order: -1;
    display: flex;
    justify-content: center;
  }

  .hero-avatar-wrap { width: 120px; height: 120px; }
  .hero-avatar, .hero-avatar-placeholder { width: 120px; height: 120px; }
  .hero-avatar-placeholder { font-size: 42px; }

  .hero-socials { justify-content: center; }

  .article-grid { grid-template-columns: 1fr; }
}
</style>
