<!-- File: app/pages/gallery/index.vue
     Purpose: Gallery stream page with masonry layout, tag filtering, BlurHash placeholders, and floating context label.
     Module: frontend-blog/pages/gallery, presentation layer.
     Related: composables/useApi.ts, pages/gallery/[slug].vue, InfiniteScrollTrigger. -->
<template>
  <div class="gallery-page">
    <div class="container container--wide">

      <!-- Page hero -->
      <div class="page-hero">
        <span class="hero-eyebrow">{{ t('gallery.eyebrow') }}</span>
        <div class="hero-main">
          <div class="hero-copy">
            <h1 class="page-title">{{ t('gallery.title') }}</h1>
            <p class="page-subtitle">{{ t('gallery.subtitle') }}</p>
          </div>
          <div class="hero-stats">
            <span class="hero-stat">{{ t('gallery.total', { n: total }) }}</span>
          </div>
        </div>
      </div>

      <!-- Tag filter pills -->
      <div v-if="galleryTags?.length" class="filter-bar">
        <button
          class="filter-pill"
          :class="{ active: !activeTag }"
          @click="setTag('')"
        >
          {{ t('gallery.allTags') }}
        </button>
        <button
          v-for="tag in galleryTags"
          :key="tag.slug"
          class="filter-pill"
          :class="{ active: activeTag === tag.slug }"
          @click="setTag(tag.slug)"
        >
          {{ tag.name }}
        </button>
      </div>

      <!-- Floating context label -->
      <div v-if="contextLabel" class="gallery-context-label">
        {{ contextLabel }}
      </div>

      <!-- Skeleton loading -->
      <div v-if="pending" class="masonry-grid">
        <div v-for="n in 8" :key="n" class="masonry-item">
          <div class="photo-card skeleton-card" :style="{ paddingBottom: `${60 + (n % 4) * 15}%` }" />
        </div>
      </div>

      <!-- Masonry grid -->
      <div v-else-if="allPhotos.length" ref="masonryRef" class="masonry-grid">
        <NuxtLink
          v-for="(photo, i) in allPhotos"
          :key="photo.id"
          :ref="(el) => setPhotoRef(i, el as any)"
          :to="localePath(`/gallery/${photo.slug}`)"
          class="masonry-item"
          :style="{ animationDelay: `${(i % 12) * 60}ms` }"
          @mouseenter="hoveredIndex = i"
          @mouseleave="hoveredIndex = -1"
        >
          <div class="photo-card" :style="{ paddingBottom: `${(photo.height / photo.width) * 100}%` }">
            <!-- BlurHash placeholder -->
            <canvas
              v-if="photo.placeholderData"
              :ref="(el) => renderBlurHash(el as HTMLCanvasElement, photo.placeholderData!)"
              class="photo-placeholder"
              width="32"
              height="32"
            />

            <!-- Actual image -->
            <img
              :src="photo.displayUrl"
              :alt="photo.title || photo.slug"
              :width="photo.width"
              :height="photo.height"
              class="photo-img"
              loading="lazy"
              @load="($event.target as HTMLImageElement)?.classList.add('loaded')"
            />

            <!-- Hover overlay (desktop) -->
            <div class="photo-overlay">
              <div class="overlay-info">
                <span v-if="photo.locationCity" class="overlay-chip">
                  <svg class="overlay-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M21 10c0 7-9 13-9 13s-9-6-9-13a9 9 0 0 1 18 0z"/><circle cx="12" cy="10" r="3"/></svg>
                  {{ photo.locationCity }}
                </span>
                <span v-if="photo.cameraModel" class="overlay-chip">
                  <svg class="overlay-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M23 19a2 2 0 0 1-2 2H3a2 2 0 0 1-2-2V8a2 2 0 0 1 2-2h4l2-3h6l2 3h4a2 2 0 0 1 2 2z"/><circle cx="12" cy="13" r="4"/></svg>
                  {{ photo.cameraModel }}
                </span>
                <span v-if="photo.tagSlugs?.length" class="overlay-chip overlay-tags">
                  {{ photo.tagSlugs.slice(0, 2).map(s => getTagName(s)).join(', ') }}
                </span>
              </div>
            </div>
          </div>
        </NuxtLink>
      </div>

      <div v-else class="empty-state">
        <p>{{ t('gallery.noPhotos') }}</p>
      </div>

      <!-- Infinite scroll -->
      <InfiniteScrollTrigger
        v-if="allPhotos.length > 0"
        :loading="loadingMore"
        :done="!hasMore"
        @load="loadMore"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import type { GalleryPhotoPublic, GalleryTag } from '~/composables/useApi'

const { t, locale } = useI18n()
const localePath = useLocalePath()
const { listGalleryPhotos, getGalleryTags } = useApi()
const route = useRoute()
const router = useRouter()

// ── Tag filter ────────────────────────────────────────────────────
const activeTag = ref((route.query.tag as string) || '')

// ── Load tags ─────────────────────────────────────────────────────
const { data: galleryTags } = await useAsyncData('gallery-tags', getGalleryTags, {
  default: () => [] as GalleryTag[],
})

// ── Load photos (Infinite Scroll) ─────────────────────────────────
const page = ref(1)
const allPhotos = ref<GalleryPhotoPublic[]>([])
const total = ref(0)
const loadingMore = ref(false)

const { data: initialData, pending } = await useAsyncData(
  'gallery-first-page',
  () => listGalleryPhotos({
    page: 1,
    pageSize: 20,
    tag: activeTag.value || undefined,
  }),
  { watch: [activeTag] }
)

watch(initialData, (newVal) => {
  if (newVal) {
    allPhotos.value = newVal.rows || []
    total.value = newVal.total || 0
    page.value = 1
  }
}, { immediate: true })

const hasMore = computed(() => allPhotos.value.length < total.value)

async function loadMore() {
  if (loadingMore.value || !hasMore.value) return
  loadingMore.value = true
  page.value++
  try {
    const res = await listGalleryPhotos({
      page: page.value,
      pageSize: 20,
      tag: activeTag.value || undefined,
    })
    if (res.rows?.length) {
      allPhotos.value.push(...res.rows)
    }
  } catch {
    page.value--
  } finally {
    loadingMore.value = false
  }
}

// ── Tag filter ────────────────────────────────────────────────────
function setTag(slug: string) {
  activeTag.value = activeTag.value === slug ? '' : slug
  page.value = 1
}

watch(activeTag, () => {
  router.replace({
    query: activeTag.value ? { tag: activeTag.value } : {},
  })
})

watch(() => route.query, (query) => {
  const t = (query.tag as string) || ''
  if (t !== activeTag.value) {
    activeTag.value = t
  }
}, { deep: true })

// ── Floating context label ────────────────────────────────────────
const hoveredIndex = ref(-1)
const photoRefs = ref<Map<number, Element>>(new Map())
const masonryRef = ref<HTMLElement | null>(null)

function setPhotoRef(index: number, el: Element | null) {
  if (el) {
    photoRefs.value.set(index, el)
  }
}

const contextLabel = computed(() => {
  // Collect cities from all currently loaded photos for a simple label
  const visible = allPhotos.value.filter(p => p.locationCity)
  if (!visible.length) return ''

  const cities = [...new Set(visible.map(p => p.locationCity!))]
  if (cities.length <= 2) {
    return cities.join(', ')
  }
  return t('gallery.contextLabelMore', { cities: cities.slice(0, 2).join(', '), n: cities.length - 2 })
})

// ── BlurHash rendering ────────────────────────────────────────────
function renderBlurHash(canvas: HTMLCanvasElement | null, _hash: string) {
  // BlurHash decoding would require a client-side library.
  // For SSR safety, we just set a neutral background; the placeholder canvas
  // acts as a subtle aspect-ratio holder until the real image loads.
  if (!canvas || !import.meta.client) return
  const ctx = canvas.getContext('2d')
  if (!ctx) return
  ctx.fillStyle = '#e5e7eb'
  ctx.fillRect(0, 0, 32, 32)
}

// ── Helpers ───────────────────────────────────────────────────────
function getTagName(slug: string) {
  return galleryTags.value?.find(t => t.slug === slug)?.name ?? slug
}

// ── SEO ──────────────────────────────────────────────────────────
useSeoMeta({ title: t('gallery.title') })
</script>

<style scoped>
.gallery-page {
  padding-top: calc(var(--header-h) + 48px);
  padding-bottom: 80px;
}

/* ── Hero ── */
.page-hero {
  display: flex;
  flex-direction: column;
  gap: 16px;
  margin-bottom: 40px;
  padding: 4px 0 26px;
  border-bottom: 1px solid color-mix(in srgb, var(--border) 78%, transparent);
}

.hero-eyebrow {
  font-size: 11px;
  font-weight: 700;
  letter-spacing: 0.18em;
  text-transform: uppercase;
  color: var(--accent);
}

.hero-main {
  display: flex;
  align-items: flex-end;
  justify-content: space-between;
  gap: 24px;
  flex-wrap: wrap;
}

.hero-copy { max-width: 620px; }

.page-title {
  font-size: clamp(1.8rem, 4vw, 2.8rem);
  font-weight: 800;
  letter-spacing: -0.02em;
  margin: 0;
}

.page-subtitle {
  margin: 10px 0 0;
  font-size: 15px;
  line-height: 1.8;
  color: var(--text-subtle);
}

.hero-stats { display: flex; flex-wrap: wrap; gap: 10px; }

.hero-stat {
  display: inline-flex;
  align-items: center;
  min-height: 34px;
  padding: 0 12px;
  border-radius: 999px;
  border: 1px solid color-mix(in srgb, var(--border) 78%, transparent);
  background: color-mix(in srgb, var(--bg-secondary) 68%, transparent);
  color: var(--text-muted);
  font-size: 12px;
  white-space: nowrap;
}

/* ── Filter pills ── */
.filter-bar {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  margin-bottom: 24px;
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

.filter-pill:hover,
.filter-pill.active {
  background: var(--accent-soft);
  border-color: var(--accent);
  color: var(--accent-text);
  font-weight: 600;
}

/* ── Floating context label ── */
.gallery-context-label {
  position: sticky;
  top: calc(var(--header-h) + 12px);
  z-index: 20;
  display: inline-flex;
  align-items: center;
  padding: 6px 14px;
  border-radius: var(--radius-md);
  background: color-mix(in srgb, var(--surface) 72%, transparent);
  backdrop-filter: blur(16px) saturate(1.4);
  -webkit-backdrop-filter: blur(16px) saturate(1.4);
  border: 1px solid color-mix(in srgb, var(--border) 40%, transparent);
  font-size: 12px;
  color: var(--text-muted);
  margin-bottom: 16px;
  pointer-events: none;
  box-shadow: var(--shadow-sm);
}

/* ── Masonry grid ── */
.masonry-grid {
  columns: 3;
  column-gap: 16px;
}

.masonry-item {
  break-inside: avoid;
  margin-bottom: 16px;
  animation: fadeUp 0.5s var(--ease-smooth) both;
  text-decoration: none;
  display: block;
}

@keyframes fadeUp {
  from {
    opacity: 0;
    transform: translateY(20px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

/* ── Photo card ── */
.photo-card {
  position: relative;
  width: 100%;
  border-radius: var(--radius-md);
  overflow: hidden;
  background: var(--bg-secondary);
  cursor: pointer;
}

.photo-placeholder {
  position: absolute;
  inset: 0;
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.photo-img {
  position: absolute;
  inset: 0;
  width: 100%;
  height: 100%;
  object-fit: cover;
  opacity: 0;
  transition: opacity 0.4s var(--ease-smooth);
}

.photo-img.loaded {
  opacity: 1;
}

/* ── Hover overlay (frosted glass) ── */
.photo-overlay {
  position: absolute;
  inset: 0;
  display: flex;
  align-items: flex-end;
  padding: 12px;
  opacity: 0;
  transition: opacity var(--dur-base) var(--ease-smooth);
  background: linear-gradient(
    to top,
    rgba(0, 0, 0, 0.45) 0%,
    rgba(0, 0, 0, 0.1) 40%,
    transparent 100%
  );
}

.masonry-item:hover .photo-overlay {
  opacity: 1;
}

.overlay-info {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
}

.overlay-chip {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  padding: 3px 10px;
  border-radius: 999px;
  background: rgba(255, 255, 255, 0.18);
  backdrop-filter: blur(10px);
  -webkit-backdrop-filter: blur(10px);
  color: #fff;
  font-size: 11px;
  font-weight: 500;
  white-space: nowrap;
}

.overlay-icon {
  width: 12px;
  height: 12px;
  flex-shrink: 0;
}

.overlay-tags {
  font-style: italic;
  opacity: 0.85;
}

/* ── Skeleton ── */
.skeleton-card {
  background: linear-gradient(
    110deg,
    var(--bg-secondary) 8%,
    color-mix(in srgb, var(--bg-secondary) 60%, var(--surface)) 18%,
    var(--bg-secondary) 33%
  );
  background-size: 200% 100%;
  animation: shimmer 1.4s linear infinite;
  border-radius: var(--radius-md);
}

@keyframes shimmer {
  to { background-position: -200% 0; }
}

/* ── Empty state ── */
.empty-state {
  text-align: center;
  padding: 80px 0;
  color: var(--text-subtle);
}

/* ── Responsive ── */
@media (max-width: 900px) {
  .masonry-grid {
    columns: 2;
    column-gap: 12px;
  }
  .masonry-item {
    margin-bottom: 12px;
  }
}

@media (max-width: 520px) {
  .masonry-grid {
    columns: 1;
    column-gap: 0;
  }
  .masonry-item {
    margin-bottom: 12px;
  }
  /* On mobile, show a subtle bottom bar instead of hover overlay */
  .photo-overlay {
    opacity: 1;
    background: linear-gradient(
      to top,
      rgba(0, 0, 0, 0.35) 0%,
      transparent 60%
    );
  }
}
</style>
