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
            />

            <!-- Hover overlay (desktop) -->
            <div class="photo-overlay">
              <div class="overlay-info">
                <div class="overlay-name">{{ formatPhotoName(photo) }}</div>
                <div class="overlay-spec">{{ formatPhotoAssetLine(photo) }}</div>
                <div v-if="formatCameraLine(photo)" class="overlay-camera">{{ formatCameraLine(photo) }}</div>
                <div v-if="formatExposureLine(photo)" class="overlay-exposure">{{ formatExposureLine(photo) }}</div>
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
import { decode } from 'blurhash'
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
function renderBlurHash(canvas: HTMLCanvasElement | null, hash: string) {
  if (!canvas || !import.meta.client || !hash) return
  const ctx = canvas.getContext('2d')
  if (!ctx) return

  try {
    const width = 32
    const height = 32
    const pixels = decode(hash, width, height)
    const imageData = ctx.createImageData(width, height)
    imageData.data.set(pixels)
    ctx.putImageData(imageData, 0, 0)
  } catch {
    ctx.fillStyle = '#e5e7eb'
    ctx.fillRect(0, 0, 32, 32)
  }
}

// ── Helpers ───────────────────────────────────────────────────────
function getTagName(slug: string) {
  return galleryTags.value?.find(t => t.slug === slug)?.name ?? slug
}

function formatPhotoName(photo: GalleryPhotoPublic) {
  if (photo.title?.trim()) {
    return photo.title.trim()
  }
  return formatFallbackPhotoName(photo.slug)
}

function formatFallbackPhotoName(slug: string) {
  const token = slug.replace(/[^a-zA-Z0-9]/g, '').slice(-8).toUpperCase()
  return `IMG${token || '00000000'}`
}

function formatPhotoAssetLine(photo: GalleryPhotoPublic) {
  const type = photo.displayUrl.split('.').pop()?.toUpperCase() || 'JPEG'
  const resolution = photo.width && photo.height ? `${photo.width} × ${photo.height}` : ''
  return [type, resolution, formatFileSize(photo.fileSizeBytes)].filter(Boolean).join(' · ')
}

function formatFileSize(bytes?: number) {
  if (!bytes || bytes <= 0) return ''
  if (bytes < 1024) return `${bytes} B`
  if (bytes < 1024 * 1024) return `${(bytes / 1024).toFixed(1)} KB`
  return `${(bytes / 1024 / 1024).toFixed(2)} MB`
}

function formatCameraLine(photo: GalleryPhotoPublic) {
  return [photo.cameraMake, photo.cameraModel].filter(Boolean).join(' ')
}

function formatExposureLine(photo: GalleryPhotoPublic) {
  return [photo.focalLength, photo.aperture, photo.shutterSpeed, photo.iso ? `ISO ${photo.iso}` : '']
    .filter(Boolean)
    .join(' · ')
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
  column-gap: 4px;
}

.masonry-item {
  break-inside: avoid;
  margin-bottom: 4px;
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
  transition: transform 360ms var(--ease-smooth);
}

.masonry-item:hover .photo-img {
  transform: scale(1.04);
}

/* ── Hover overlay (frosted glass) ── */
.photo-overlay {
  position: absolute;
  inset: 0;
  display: flex;
  align-items: flex-end;
  padding: 16px 18px;
  opacity: 0;
  transition: opacity var(--dur-base) var(--ease-smooth);
  background: linear-gradient(
    to top,
    rgba(0, 0, 0, 0.72) 0%,
    rgba(0, 0, 0, 0.24) 55%,
    transparent 100%
  );
}

.masonry-item:hover .photo-overlay {
  opacity: 1;
}

.overlay-info {
  width: 100%;
  color: rgba(255, 255, 255, 0.92);
}

.overlay-name {
  font-size: 14px;
  font-weight: 800;
  line-height: 1.35;
  letter-spacing: -0.02em;
}

.overlay-spec,
.overlay-camera,
.overlay-exposure {
  margin-top: 5px;
  font-size: 11px;
  line-height: 1.5;
  color: rgba(255, 255, 255, 0.78);
}

.overlay-camera {
  color: #fff;
  font-size: 12px;
  font-weight: 700;
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
    column-gap: 4px;
  }
  .masonry-item {
    margin-bottom: 4px;
  }
}

@media (max-width: 320px) {
  .masonry-grid {
    columns: 1;
    column-gap: 3px;
  }
  .masonry-item {
    margin-bottom: 3px;
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
