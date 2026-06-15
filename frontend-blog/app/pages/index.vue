<!-- File: app/pages/index.vue
     Purpose: Homepage with atmospheric hero, recent writing, moments, and gallery preview.
     Module: frontend-blog/pages, presentation layer.
     Related: layouts/default.vue, components/ArticleCard.vue, composables/useApi.ts. -->
<template>
  <div class="home">
    <section class="home-hero">
      <div class="atelier-field" aria-hidden="true">
        <span class="paper-mark paper-mark--one" />
        <span class="paper-mark paper-mark--two" />
        <span class="paper-mark paper-mark--three" />
      </div>

      <div class="hero-inner container">
        <div class="hero-copy">
          <h1 class="hero-title">
            <span>Hi, I'm</span>
            <span class="hero-avatar-wrap">
              <img
                v-if="siteSettings?.avatarUrl"
                :src="siteSettings.avatarUrl"
                :alt="siteSettings?.siteName"
                class="hero-avatar"
              />
              <span v-else class="hero-avatar-fallback">{{ siteInitial }}</span>
            </span>
            <strong>{{ displayName }}</strong>
          </h1>
          <p class="hero-slogan">
            {{ heroSlogan }}
          </p>
          <p v-if="heroDescription" class="hero-description">
            {{ heroDescription }}
          </p>
        </div>

        <div v-if="heroQuote" class="hero-quote">
          <span class="quote-mark">「</span>
          <span>{{ heroQuote }}</span>
          <span class="quote-mark">」</span>
        </div>

        <div v-if="socialLinks.length" class="hero-socials" :aria-label="t('home.socialLinks')">
          <a
            v-for="link in socialLinks"
            :key="link.id"
            :href="link.url"
            target="_blank"
            rel="noreferrer"
            class="social-link"
            :title="link.title"
          >
            <img v-if="link.iconKey" :src="link.iconKey" :alt="link.title" class="social-icon-img" />
            <span v-else>{{ link.title.slice(0, 1) }}</span>
          </a>
        </div>
      </div>
    </section>

    <section class="home-overview">
      <div class="container container--wide overview-grid">
        <section class="writing-panel">
          <div class="section-header editorial-header">
            <div>
              <p class="section-eyebrow">{{ t('home.recentWritingEyebrow') }}</p>
              <h2 class="section-title">{{ t('home.recentArticles') }}</h2>
            </div>
            <NuxtLink :to="localePath('/articles')" class="section-more">
              {{ t('home.viewAll') }}
              <svg viewBox="0 0 16 16" fill="none" stroke="currentColor" stroke-width="1.5">
                <path d="M3 8h10M9 4l4 4-4 4" />
              </svg>
            </NuxtLink>
          </div>

          <div v-if="pending" class="writing-list">
            <div v-for="n in 5" :key="n" class="writing-skeleton" />
          </div>

          <ol v-else-if="articles?.rows?.length" class="writing-list">
            <li v-for="(article, i) in articles.rows.slice(0, 5)" :key="article.id" class="writing-item">
              <NuxtLink :to="localePath(`/articles/${article.slug}`)" class="writing-link">
                <span class="writing-index">{{ String(i + 1).padStart(2, '0') }}</span>
                <span class="writing-main">
                  <span class="writing-title">{{ article.title }}</span>
                  <span class="writing-meta">
                    {{ formatArticleMeta(article) }}
                  </span>
                </span>
                <span class="writing-arrow">→</span>
              </NuxtLink>
            </li>
          </ol>

          <div v-else class="empty-state">
            <p>{{ t('home.noArticles') }}</p>
          </div>
        </section>

        <aside class="side-panel">
          <section class="musing-card">
            <div class="section-header compact">
              <div>
                <p class="section-eyebrow">{{ t('home.musingsEyebrow') }}</p>
                <h2 class="section-title">{{ t('home.musingsTitle') }}</h2>
              </div>
              <NuxtLink :to="localePath('/moments')" class="section-more">
                {{ t('home.viewAll') }}
              </NuxtLink>
            </div>
            <blockquote v-if="latestMoment" class="moment-preview">
              {{ normalizeMoment(latestMoment.content) }}
            </blockquote>
            <p v-else class="side-empty">{{ t('moments.empty') }}</p>
          </section>

          <section class="gallery-strip">
            <div class="section-header compact">
              <div>
                <p class="section-eyebrow">{{ t('home.galleryEyebrow') }}</p>
                <h2 class="section-title">{{ t('home.galleryTitle') }}</h2>
              </div>
              <NuxtLink :to="localePath('/gallery')" class="section-more">
                {{ t('home.viewAll') }}
              </NuxtLink>
            </div>

            <div v-if="photos?.rows?.length" class="photo-row">
              <NuxtLink
                v-for="photo in photos.rows.slice(0, 3)"
                :key="photo.id"
                :to="localePath(`/gallery/${photo.slug}`)"
                class="photo-thumb"
              >
                <img :src="photo.displayUrl" :alt="photo.title || photo.slug" loading="lazy" />
              </NuxtLink>
            </div>
            <p v-else class="side-empty">{{ t('gallery.noPhotos') }}</p>
          </section>
        </aside>
      </div>
    </section>
  </div>
</template>

<script setup lang="ts">
import type { ArticleCard, Moment } from '~/composables/useApi'

const { t, locale } = useI18n()
const localePath = useLocalePath()
const { listArticles, listMoments, listGalleryPhotos, getSiteSettings, getSocialLinks } = useApi()

const [
  { data: siteSettings },
  { data: socialLinks },
  { data: articles, pending },
  { data: moments },
  { data: photos },
] = await Promise.all([
  useAsyncData('site-settings', getSiteSettings),
  useAsyncData('social-links', getSocialLinks, { default: () => [] }),
  useAsyncData('home-articles', () => listArticles({ pageSize: 5 })),
  useAsyncData('home-moments', () => listMoments({ pageSize: 1 })),
  useAsyncData('home-gallery', () => listGalleryPhotos({ pageSize: 3 })),
])

const heroSlogan = computed(() => siteSettings.value?.heroIntroMd?.trim() || t('home.heroSubtitle'))
const heroDescription = computed(() => siteSettings.value?.siteDescription?.trim() || t('home.heroDescription'))
const heroQuoteRandom = useState('home-hero-quote-random', () => Math.random())
const heroQuote = computed(() => {
  const targetLocale = locale.value === 'en' ? 'en-US' : 'zh-CN'
  const quotes = (siteSettings.value?.heroQuotes || [])
    .filter((quote) => quote.locale === targetLocale && quote.text.trim())
    .map((quote) => quote.text.trim())

  if (!quotes.length) return ''
  return quotes[Math.floor(heroQuoteRandom.value * quotes.length) % quotes.length]
})
const displayName = computed(() => {
  const raw = siteSettings.value?.siteName?.trim() || 'Ancy'
  const parts = raw.split(/\s+/)
  return parts.length > 1 ? parts[parts.length - 1] : raw
})
const siteInitial = computed(() => (siteSettings.value?.siteName || 'A').slice(0, 1).toUpperCase())
const latestMoment = computed(() => moments.value?.rows?.[0] || null)

function formatArticleMeta(article: ArticleCard) {
  const type = article.contentKind === 'page' ? t('home.pageLabel') : t('home.postLabel')
  const date = formatDate(article.publishedAt || article.createdAt)
  const category = article.categorySlug ? ` · ${article.categorySlug}` : ''
  return `${type} · ${date}${category}`
}

function formatDate(iso?: string): string {
  if (!iso) return ''
  return new Date(iso).toLocaleDateString('zh-CN', {
    month: 'short',
    day: 'numeric',
  })
}

function normalizeMoment(content: Moment['content']) {
  return content
    .replace(/<[^>]+>/g, '')
    .replace(/[#>*_`~-]/g, '')
    .replace(/\s+/g, ' ')
    .trim()
    .slice(0, 96)
}

useSeoMeta({
  title: siteSettings.value?.siteName || 'Ancy Blog',
  description: siteSettings.value?.siteDescription || '',
  ogTitle: siteSettings.value?.siteName,
  ogDescription: siteSettings.value?.siteDescription,
})
</script>

<style scoped>
.home {
  position: relative;
  overflow: hidden;
}

.home::before,
.home::after {
  content: '';
  position: fixed;
  pointer-events: none;
  border-radius: 999px;
  filter: blur(18px);
  opacity: 0.5;
  z-index: 0;
}

.home::before {
  width: 420px;
  height: 420px;
  left: 8%;
  top: 58dvh;
  background: radial-gradient(circle, rgba(199, 120, 135, 0.12), transparent 64%);
}

.home::after {
  width: 520px;
  height: 520px;
  right: 5%;
  top: 70dvh;
  background: radial-gradient(circle, rgba(128, 165, 194, 0.13), transparent 68%);
}

.home-hero {
  min-height: 92dvh;
  position: relative;
  z-index: 1;
  display: flex;
  align-items: center;
  padding: 116px 0 120px;
}

.atelier-field {
  position: fixed;
  inset: 0;
  pointer-events: none;
  z-index: 0;
}

.paper-mark {
  position: fixed;
  width: clamp(180px, 18vw, 360px);
  height: clamp(180px, 18vw, 360px);
  border-radius: 999px;
  background:
    radial-gradient(circle, rgba(199, 120, 135, 0.12), transparent 58%),
    radial-gradient(circle at 70% 70%, rgba(127, 159, 180, 0.10), transparent 62%);
  filter: blur(18px);
  opacity: 0.58;
  animation: paper-breathe 9s ease-in-out infinite;
}

.paper-mark--one {
  left: 18%;
  bottom: 8dvh;
}

.paper-mark--two {
  right: 20%;
  bottom: 14dvh;
  animation-delay: -3s;
}

.paper-mark--three {
  right: -3%;
  top: 20%;
  opacity: 0.34;
  animation-delay: -5s;
}

@keyframes paper-breathe {
  0%, 100% { transform: translate3d(0, 0, 0) scale(1); }
  50% { transform: translate3d(0, -10px, 0) scale(1.04); }
}

.hero-inner {
  position: relative;
  z-index: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  text-align: center;
  transform: translateY(10dvh);
}

.hero-avatar-wrap {
  width: clamp(56px, 7.2vw, 84px);
  height: clamp(56px, 7.2vw, 84px);
  display: inline-grid;
  place-items: center;
  flex: 0 0 auto;
  border-radius: 50%;
  background: color-mix(in srgb, var(--surface) 74%, transparent);
  box-shadow:
    0 0 0 clamp(6px, 0.78vw, 10px) rgba(255, 255, 255, 0.36),
    0 16px 54px rgba(88, 96, 110, 0.15);
  vertical-align: middle;
  animation: avatar-rise 720ms var(--ease-smooth) both;
}

.hero-avatar,
.hero-avatar-fallback {
  width: calc(100% - 14px);
  height: calc(100% - 14px);
  border-radius: 50%;
}

.hero-avatar {
  object-fit: cover;
}

.hero-avatar-fallback {
  display: grid;
  place-items: center;
  background: var(--accent-soft);
  color: var(--accent-text);
  font-size: clamp(1.35rem, 2.4vw, 2rem);
  font-family: var(--font-display);
}

@keyframes avatar-rise {
  from { opacity: 0; transform: translateY(18px) scale(0.96); }
  to { opacity: 1; transform: none; }
}

.section-eyebrow {
  font-family: var(--font-display);
  color: var(--text-subtle);
  letter-spacing: 0;
  text-transform: uppercase;
}

.hero-title {
  display: flex;
  align-items: center;
  justify-content: center;
  flex-wrap: wrap;
  gap: clamp(10px, 1.2vw, 18px);
  font-size: clamp(2.55rem, 5.6vw, 5.35rem);
  font-family: var(--font-display);
  font-weight: 650;
  line-height: 0.98;
  color: var(--text);
  letter-spacing: 0;
}

.hero-title span {
  color: var(--text-muted);
  font-weight: 520;
}

.hero-title strong {
  color: var(--accent-ink);
  font-weight: 760;
  font-style: italic;
}

.hero-slogan {
  max-width: 860px;
  margin: 20px auto 0;
  font-size: clamp(1.42rem, 3.25vw, 3.1rem);
  line-height: 1.18;
  font-weight: 520;
  text-wrap: balance;
  color: var(--text);
}

.hero-description {
  max-width: 760px;
  margin: 18px auto 0;
  color: var(--text-muted);
  font-size: clamp(0.95rem, 1.4vw, 1.08rem);
  letter-spacing: 0;
  text-transform: uppercase;
}

.hero-quote {
  margin-top: 176px;
  color: var(--text-subtle);
  font-family: var(--font-serif);
  font-size: 0.98rem;
  font-style: italic;
}

.quote-mark {
  color: color-mix(in srgb, var(--accent) 55%, var(--text-subtle));
}

.hero-socials {
  display: flex;
  justify-content: center;
  flex-wrap: wrap;
  gap: 12px;
  margin-top: clamp(148px, 20dvh, 232px);
}

.hero-quote + .hero-socials {
  margin-top: 52px;
}

.social-link {
  width: 34px;
  height: 34px;
  display: grid;
  place-items: center;
  border-radius: 50%;
  border: 1px solid var(--border);
  background: color-mix(in srgb, var(--surface) 70%, transparent);
  color: var(--text-muted);
  box-shadow: var(--shadow-sm);
  transition: transform var(--dur-base) var(--ease-spring), border-color var(--dur-base), color var(--dur-base);
  overflow: hidden;
}

.social-link:hover {
  transform: translateY(-4px);
  color: var(--accent-text);
  border-color: color-mix(in srgb, var(--accent) 38%, var(--border));
}

.social-icon-img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  filter: saturate(0.68) contrast(0.92);
  opacity: 0.78;
  transition: filter var(--dur-base), opacity var(--dur-base);
}

.social-link:hover .social-icon-img {
  filter: none;
  opacity: 1;
}

.home-overview {
  position: relative;
  z-index: 1;
  padding: 4px 0 90px;
}

.overview-grid {
  display: grid;
  grid-template-columns: minmax(0, 1.24fr) minmax(240px, 0.68fr);
  gap: clamp(38px, 8vw, 118px);
  align-items: start;
}

.section-header {
  display: flex;
  align-items: flex-end;
  justify-content: space-between;
  gap: 24px;
  margin-bottom: 18px;
}

.section-header.compact {
  align-items: flex-start;
  margin-bottom: 18px;
}

.section-eyebrow {
  margin-bottom: 5px;
  font-size: 0.66rem;
  color: var(--text-subtle);
}

.section-title {
  font-family: var(--font-serif);
  font-size: clamp(1.36rem, 2.2vw, 1.82rem);
  font-weight: 500;
  line-height: 1.2;
  color: var(--text);
}

.section-more {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  color: var(--accent-text);
  font-size: 0.78rem;
  font-weight: 600;
  white-space: nowrap;
}

.section-more svg {
  width: 14px;
  height: 14px;
}

.writing-list {
  display: flex;
  flex-direction: column;
  list-style: none;
}

.writing-item {
  border-top: 1px solid var(--border);
}

.writing-item:last-child {
  border-bottom: 1px solid var(--border);
}

.writing-link {
  display: grid;
  grid-template-columns: 42px minmax(0, 1fr) 20px;
  gap: 14px;
  align-items: center;
  padding: 14px 0;
  transition: color var(--dur-base), transform var(--dur-base) var(--ease-spring);
}

.writing-link:hover {
  color: var(--accent-text);
  transform: translateX(6px);
}

.writing-index {
  font-family: var(--font-display);
  color: color-mix(in srgb, var(--accent) 70%, var(--text-subtle));
  font-size: 0.82rem;
  font-style: italic;
}

.writing-main {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.writing-title {
  color: var(--text);
  font-family: var(--font-serif);
  font-size: clamp(0.94rem, 1.55vw, 1.12rem);
  font-weight: 560;
  line-height: 1.55;
}

.writing-link:hover .writing-title {
  color: var(--accent-text);
}

.writing-meta {
  color: var(--text-subtle);
  font-size: 0.74rem;
}

.writing-arrow {
  color: var(--text-subtle);
  transition: transform var(--dur-base);
}

.writing-link:hover .writing-arrow {
  transform: translateX(5px);
}

.writing-skeleton {
  height: 74px;
  border-top: 1px solid var(--border);
  background: linear-gradient(90deg, transparent, color-mix(in srgb, var(--surface-hover) 70%, transparent), transparent);
  background-size: 220% 100%;
  animation: shimmer 1.4s infinite;
}

@keyframes shimmer {
  from { background-position: 220% 0; }
  to { background-position: -220% 0; }
}

.side-panel {
  display: flex;
  flex-direction: column;
  gap: 34px;
}

.musing-card,
.gallery-strip {
  position: relative;
  padding: 0 0 0 24px;
  border-left: 1px solid var(--border);
  background: transparent;
  box-shadow: none;
}

.moment-preview {
  color: var(--text-muted);
  font-family: var(--font-serif);
  font-size: 0.92rem;
  line-height: 1.9;
}

.photo-row {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 6px;
}

.photo-thumb {
  aspect-ratio: 4 / 5;
  overflow: hidden;
  border-radius: var(--radius-sm);
  background: var(--bg-secondary);
}

.photo-thumb img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  filter: saturate(0.82);
  transition: transform var(--dur-slow) var(--ease-smooth), filter var(--dur-base);
}

.photo-thumb:hover img {
  transform: scale(1.05);
  filter: saturate(1.04);
}

.empty-state,
.side-empty {
  color: var(--text-subtle);
  font-size: 0.82rem;
  line-height: 1.8;
}

@media (max-width: 820px) {
  .home-hero {
    min-height: 84dvh;
    padding: 96px 0 76px;
  }

  .hero-inner {
    transform: translateY(5dvh);
  }

  .hero-avatar-wrap {
    width: 56px;
    height: 56px;
  }

  .hero-avatar,
  .hero-avatar-fallback {
    width: 46px;
    height: 46px;
  }

  .hero-description {
    text-transform: none;
    letter-spacing: 0;
  }

  .hero-quote {
    margin-top: 76px;
  }

  .hero-socials {
    margin-top: 116px;
  }

  .hero-quote + .hero-socials {
    margin-top: 42px;
  }

  .overview-grid {
    grid-template-columns: 1fr;
  }

  .writing-link {
    grid-template-columns: 42px minmax(0, 1fr) 20px;
    gap: 12px;
  }

  .musing-card,
  .gallery-strip {
    padding: 22px;
  }
}

@media (prefers-reduced-motion: reduce) {
  .paper-mark,
  .hero-avatar-wrap,
  .writing-skeleton {
    animation: none;
  }
}
</style>
