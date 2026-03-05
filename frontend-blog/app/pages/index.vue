<!-- File: app/pages/index.vue
     Purpose: Homepage — full-height hero + recent article grid.
     Module: frontend-blog/pages, presentation layer.
     Related: layouts/default.vue, components/ArticleCard.vue, composables/useApi.ts. -->
<template>
  <div class="home">

    <!-- ═══ Hero ════════════════════════════════════════════════════ -->
    <section class="hero">
      <div class="hero-inner container container--wide">

        <!-- Left: intro + socials -->
        <div class="hero-left" :class="{ visible: heroVisible }">
          <!-- eslint-disable-next-line vue/no-v-html -->
          <div class="hero-intro" v-html="introHtml" />

          <div v-if="siteStore.socialLinks.length" class="hero-socials">
            <a
              v-for="link in siteStore.socialLinks"
              :key="link.id"
              :href="link.url"
              target="_blank"
              rel="noopener noreferrer"
              class="social-link"
              :aria-label="link.title"
            >
              <span class="social-icon">
                <svg v-if="link.platform === 'github'" viewBox="0 0 24 24" fill="currentColor">
                  <path d="M12 2C6.477 2 2 6.484 2 12.017c0 4.425 2.865 8.18 6.839 9.504.5.092.682-.217.682-.483 0-.237-.008-.868-.013-1.703-2.782.605-3.369-1.343-3.369-1.343-.454-1.158-1.11-1.466-1.11-1.466-.908-.62.069-.608.069-.608 1.003.07 1.531 1.032 1.531 1.032.892 1.53 2.341 1.088 2.91.832.092-.647.35-1.088.636-1.338-2.22-.253-4.555-1.113-4.555-4.951 0-1.093.39-1.988 1.029-2.688-.103-.253-.446-1.272.098-2.65 0 0 .84-.27 2.75 1.026A9.564 9.564 0 0 1 12 6.844a9.59 9.59 0 0 1 2.504.337c1.909-1.296 2.747-1.027 2.747-1.027.546 1.379.202 2.398.1 2.651.64.7 1.028 1.595 1.028 2.688 0 3.848-2.339 4.695-4.566 4.943.359.309.678.92.678 1.855 0 1.338-.012 2.419-.012 2.747 0 .268.18.58.688.482A10.02 10.02 0 0 0 22 12.017C22 6.484 17.522 2 12 2z"/>
                </svg>
                <svg v-else-if="link.platform === 'mail'" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                  <path d="M4 4h16c1.1 0 2 .9 2 2v12c0 1.1-.9 2-2 2H4c-1.1 0-2-.9-2-2V6c0-1.1.9-2 2-2z"/><polyline points="22,6 12,13 2,6"/>
                </svg>
                <svg v-else viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="M10 13a5 5 0 0 0 7.54.54l3-3a5 5 0 0 0-7.07-7.07l-1.72 1.71"/>
                  <path d="M14 11a5 5 0 0 0-7.54-.54l-3 3a5 5 0 0 0 7.07 7.07l1.71-1.71"/>
                </svg>
              </span>
              <span class="social-label">{{ link.title }}</span>
            </a>
          </div>
        </div>

        <!-- Right: avatar -->
        <div class="hero-right" :class="{ visible: heroVisible }">
          <div class="hero-avatar-wrap">
            <div class="hero-avatar-ring" />
            <img
              v-if="siteStore.settings?.avatarUrl"
              :src="siteStore.settings.avatarUrl"
              :alt="siteStore.settings?.siteName"
              class="hero-avatar"
              width="180"
              height="180"
            />
            <div v-else class="hero-avatar-placeholder">A</div>
          </div>
        </div>
      </div>

      <!-- Down arrow -->
      <div class="hero-arrow" :class="{ visible: heroVisible }">
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
import { useSiteStore } from '~/stores/site'

const { t } = useI18n()
const localePath = useLocalePath()
const { listArticles } = useApi()
const siteStore = useSiteStore()

// ── Fetch data ───────────────────────────────────────────────────
await siteStore.fetchAll()
const { data: articles, pending } = await useAsyncData('home-articles', () => listArticles({ pageSize: 6 }))

// ── Intro HTML — convert simple markdown line breaks ────────────
const introHtml = computed(() => {
  const raw = siteStore.settings?.heroIntroMd || `Hi, I'm **Ancy** 👋\nI write code and thoughts.`
  return raw
    .replace(/\*\*(.*?)\*\*/g, '<strong>$1</strong>')
    .replace(/\n/g, '<br/>')
})

// ── Hero entrance animation ──────────────────────────────────────
const heroVisible = ref(false)
onMounted(() => {
  const id = requestAnimationFrame(() => { heroVisible.value = true })
  onUnmounted(() => cancelAnimationFrame(id))
})

// ── SEO ─────────────────────────────────────────────────────────
useSeoMeta({
  title: siteStore.settings?.siteName || 'Ancy Blog',
  description: siteStore.settings?.siteDescription || '',
  ogTitle: siteStore.settings?.siteName,
  ogDescription: siteStore.settings?.siteDescription,
})
</script>

<style scoped>
/* ═══ Hero ═══════════════════════════════════════════════════════ */
.hero {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  position: relative;
  padding-top: var(--header-h);
}

.hero-inner {
  display: grid;
  grid-template-columns: 1fr auto;
  gap: 60px;
  align-items: center;
  width: 100%;
}

/* ── Left ── */
.hero-left {
  opacity: 0;
  transform: translateY(24px);
  transition: opacity 0.7s var(--ease-smooth), transform 0.7s var(--ease-spring);
}
.hero-left.visible { opacity: 1; transform: none; }

.hero-intro {
  font-size: clamp(1.5rem, 3.5vw, 2.1rem);
  font-weight: 400;
  line-height: 1.5;
  color: var(--text);
  margin-bottom: 28px;
  font-family: 'Noto Serif SC', serif;
  letter-spacing: 0.01em;
}

.hero-intro :deep(strong) {
  font-weight: 800;
  color: var(--accent-text);
}

.hero-socials {
  display: flex;
  gap: 12px;
  flex-wrap: wrap;
}

.social-link {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 8px 14px;
  border: 1px solid var(--border);
  border-radius: var(--radius-xl);
  font-size: 13px;
  font-weight: 500;
  color: var(--text-muted);
  background: var(--surface);
  transition: border-color var(--dur-base), color var(--dur-base),
              background var(--dur-base), transform var(--dur-base) var(--ease-spring),
              box-shadow var(--dur-base);
}

.social-link:hover {
  border-color: var(--accent);
  color: var(--accent-text);
  background: var(--accent-soft);
  transform: translateY(-2px);
  box-shadow: var(--shadow-sm);
}

.social-icon { display: flex; align-items: center; }
.social-icon svg { width: 16px; height: 16px; }

/* ── Right ── */
.hero-right {
  opacity: 0;
  transform: translateX(24px) scale(0.95);
  transition: opacity 0.7s 0.1s var(--ease-smooth), transform 0.7s 0.1s var(--ease-spring);
}
.hero-right.visible { opacity: 1; transform: none; }

.hero-avatar-wrap {
  position: relative;
  width: 180px;
  height: 180px;
}

.hero-avatar-ring {
  position: absolute;
  inset: -6px;
  border-radius: 50%;
  border: 2px dashed var(--accent);
  opacity: 0.35;
  animation: spin 20s linear infinite;
}

@keyframes spin { to { transform: rotate(360deg); } }

.hero-avatar, .hero-avatar-placeholder {
  width: 180px;
  height: 180px;
  border-radius: 50%;
  object-fit: cover;
  border: 3px solid var(--surface);
  box-shadow: var(--shadow-lg);
}

.hero-avatar-placeholder {
  background: var(--accent-soft);
  display: grid;
  place-items: center;
  font-size: 64px;
  font-weight: 800;
  color: var(--accent-text);
}

/* ── Down arrow ── */
.hero-arrow {
  position: absolute;
  bottom: 32px;
  left: 50%;
  transform: translateX(-50%);
  opacity: 0;
  transition: opacity 1s 0.8s;
  animation: float 2.5s ease-in-out infinite;
}

.hero-arrow.visible { opacity: 0.4; }

.hero-arrow svg { width: 24px; height: 24px; color: var(--text-muted); }

@keyframes float {
  0%, 100% { transform: translateX(-50%) translateY(0); }
  50%       { transform: translateX(-50%) translateY(6px); }
}

/* ═══ Recent articles ════════════════════════════════════════════ */
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
