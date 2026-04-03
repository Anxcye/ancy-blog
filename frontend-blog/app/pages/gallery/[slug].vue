<!-- File: app/pages/gallery/[slug].vue
     Purpose: Single photo viewer with large image display and detail side panel.
     Module: frontend-blog/pages/gallery, presentation layer.
     Related: composables/useApi.ts, pages/gallery/index.vue. -->
<template>
  <div class="photo-viewer">
    <!-- Close button -->
    <NuxtLink :to="localePath('/gallery')" class="viewer-close" :aria-label="t('gallery.close')">
      <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round">
        <line x1="18" y1="6" x2="6" y2="18" /><line x1="6" y1="6" x2="18" y2="18" />
      </svg>
    </NuxtLink>

    <!-- Mobile details toggle -->
    <button class="viewer-details-toggle" @click="detailOpen = !detailOpen">
      {{ t('gallery.details') }}
    </button>

    <!-- Main image area -->
    <div class="viewer-main" @click="detailOpen = false">
      <div
        v-if="photo"
        class="viewer-stage"
        :style="{ aspectRatio: photo.width && photo.height ? `${photo.width} / ${photo.height}` : undefined }"
      >
        <canvas
          v-if="photo.placeholderData"
          ref="placeholderCanvasRef"
          class="viewer-placeholder"
          width="32"
          height="32"
        />

        <img
          :src="photo.displayUrl"
          :alt="photo.title || photo.slug"
          :width="photo.width"
          :height="photo.height"
          class="viewer-image viewer-image--display"
          @load="handleDisplayLoaded"
        />

        <img
          v-if="largeReady"
          :src="photo.largeUrl"
          :alt="photo.title || photo.slug"
          :width="photo.width"
          :height="photo.height"
          class="viewer-image viewer-image--large loaded"
        />
      </div>

      <div v-if="!displayLoaded" class="viewer-loading">
        <div class="spinner" />
      </div>
    </div>

    <!-- Detail panel -->
    <aside class="viewer-panel" :class="{ open: detailOpen }">
      <div class="panel-scroll" v-if="photo">
        <p class="panel-heading">{{ formatPhotoName() }}</p>
        <p v-if="photo.description" class="panel-desc">{{ photo.description }}</p>

        <div class="detail-section">
          <h2 class="section-title">{{ t('gallery.basicInformation') }}</h2>
          <dl class="detail-list">
            <div class="detail-row">
              <dt>{{ t('gallery.filename') }}</dt>
              <dd>{{ formatPhotoName() }}</dd>
            </div>
            <div class="detail-row">
              <dt>{{ t('gallery.fileType') }}</dt>
              <dd>{{ getFileType() }}</dd>
            </div>
            <div v-if="photo.fileSizeBytes" class="detail-row">
              <dt>{{ t('gallery.fileSize') }}</dt>
              <dd>{{ formatFileSize(photo.fileSizeBytes) }}</dd>
            </div>
            <div class="detail-row">
              <dt>{{ t('gallery.resolution') }}</dt>
              <dd>{{ formatResolution() }}</dd>
            </div>
            <div class="detail-row">
              <dt>{{ t('gallery.pixels') }}</dt>
              <dd>{{ formatMegapixels() }}</dd>
            </div>
            <div v-if="photo.takenAt" class="detail-row">
              <dt>{{ t('gallery.takenAt') }}</dt>
              <dd>{{ formatDate(photo.takenAt) }}</dd>
            </div>
            <div v-if="photo.locationCountry" class="detail-row">
              <dt>{{ t('gallery.country') }}</dt>
              <dd>{{ photo.locationCountry }}</dd>
            </div>
            <div v-if="photo.locationCity || photo.locationName" class="detail-row">
              <dt>{{ t('gallery.city') }}</dt>
              <dd>{{ photo.locationName || photo.locationCity }}</dd>
            </div>
          </dl>
        </div>

        <div v-if="hasExif" class="detail-section">
          <h2 class="section-title">{{ t('gallery.shootingParameters') }}</h2>
          <dl class="detail-list">
            <div v-if="photo.focalLength" class="detail-row">
              <dt>{{ t('gallery.focalLength') }}</dt>
              <dd>{{ photo.focalLength }}</dd>
            </div>
            <div v-if="photo.aperture" class="detail-row">
              <dt>{{ t('gallery.aperture') }}</dt>
              <dd>{{ photo.aperture }}</dd>
            </div>
            <div v-if="photo.shutterSpeed" class="detail-row">
              <dt>{{ t('gallery.shutterSpeed') }}</dt>
              <dd>{{ photo.shutterSpeed }}</dd>
            </div>
            <div v-if="photo.iso" class="detail-row">
              <dt>{{ t('gallery.iso') }}</dt>
              <dd>{{ photo.iso }}</dd>
            </div>
          </dl>
        </div>

        <div v-if="photo.cameraMake || photo.cameraModel || photo.lensModel || photo.focalLength" class="detail-section">
          <h2 class="section-title">{{ t('gallery.equipmentInformation') }}</h2>
          <dl class="detail-list">
            <div v-if="photo.cameraMake || photo.cameraModel" class="detail-row">
              <dt>{{ t('gallery.camera') }}</dt>
              <dd>{{ [photo.cameraMake, photo.cameraModel].filter(Boolean).join(' ') }}</dd>
            </div>
            <div v-if="photo.lensModel" class="detail-row">
              <dt>{{ t('gallery.lens') }}</dt>
              <dd>{{ photo.lensModel }}</dd>
            </div>
            <div v-if="photo.focalLength" class="detail-row">
              <dt>{{ t('gallery.focalLength') }}</dt>
              <dd>{{ photo.focalLength }}</dd>
            </div>
          </dl>
        </div>

        <!-- Tags -->
        <div v-if="photo.tagSlugs?.length" class="panel-tags">
          <NuxtLink
            v-for="slug in photo.tagSlugs"
            :key="slug"
            :to="localePath(`/gallery?tag=${slug}`)"
            class="panel-tag"
          >
            #{{ getTagName(slug) }}
          </NuxtLink>
        </div>
      </div>
    </aside>
  </div>
</template>

<script setup lang="ts">
import { decode } from 'blurhash'
import type { GalleryTag } from '~/composables/useApi'

const { t, locale } = useI18n()
const localePath = useLocalePath()
const route = useRoute()
const { getGalleryPhoto, getGalleryTags } = useApi()

const slug = route.params.slug as string
const detailOpen = ref(false)
const displayLoaded = ref(false)
const largeReady = ref(false)
const placeholderCanvasRef = ref<HTMLCanvasElement | null>(null)

// Fetch photo and tags in parallel
const [{ data: photo }, { data: galleryTags }] = await Promise.all([
  useAsyncData(`gallery-photo-${slug}`, () => getGalleryPhoto(slug)),
  useAsyncData('gallery-tags-viewer', getGalleryTags, { default: () => [] as GalleryTag[] }),
])

watch(
  photo,
  async (currentPhoto) => {
    displayLoaded.value = false
    largeReady.value = false

    if (!currentPhoto) return

    await nextTick()
    renderBlurHash(currentPhoto.placeholderData || '')
    preloadDisplayImage(currentPhoto.displayUrl)
  },
  { immediate: true }
)

const hasExif = computed(() =>
  photo.value?.focalLength || photo.value?.aperture || photo.value?.shutterSpeed || photo.value?.iso
)

function getTagName(slug: string) {
  return galleryTags.value?.find(t => t.slug === slug)?.name ?? slug
}

function formatDate(dateStr?: string) {
  if (!dateStr) return ''
  const d = new Date(dateStr)
  const bcp47 = locale.value === 'zh' ? 'zh-CN' : 'en-US'
  return new Intl.DateTimeFormat(bcp47, {
    year: 'numeric',
    month: 'long',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
  }).format(d)
}

function formatPhotoName() {
  if (!photo.value) return formatFallbackPhotoName(slug)
  if (photo.value.title?.trim()) {
    return photo.value.title.trim()
  }
  return formatFallbackPhotoName(photo.value.slug)
}

function formatFallbackPhotoName(rawSlug: string) {
  const token = rawSlug.replace(/[^a-zA-Z0-9]/g, '').slice(-8).toUpperCase()
  return `IMG${token || '00000000'}`
}

function getFileType() {
  if (!photo.value?.displayUrl) return 'JPEG'
  return photo.value.displayUrl.split('.').pop()?.toUpperCase() || 'JPEG'
}

function formatResolution() {
  if (!photo.value?.width || !photo.value?.height) return '-'
  return `${photo.value.width} × ${photo.value.height}`
}

function formatMegapixels() {
  if (!photo.value?.width || !photo.value?.height) return '-'
  return `${(photo.value.width * photo.value.height / 1_000_000).toFixed(2)} MP`
}

function formatFileSize(bytes?: number) {
  if (!bytes || bytes <= 0) return '-'
  if (bytes < 1024) return `${bytes} B`
  if (bytes < 1024 * 1024) return `${(bytes / 1024).toFixed(1)} KB`
  return `${(bytes / 1024 / 1024).toFixed(2)} MB`
}

function renderBlurHash(hash: string) {
  if (!placeholderCanvasRef.value || !import.meta.client || !hash) return

  const canvas = placeholderCanvasRef.value
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

function handleDisplayLoaded() {
  displayLoaded.value = true
  preloadLargeImage()
}

function preloadDisplayImage(displayURL: string) {
  if (!import.meta.client || !displayURL) return

  const img = new Image()
  img.onload = () => {
    handleDisplayLoaded()
  }
  img.src = displayURL
}

function preloadLargeImage() {
  if (!import.meta.client || !photo.value?.largeUrl) return
  if (photo.value.largeUrl === photo.value.displayUrl) {
    largeReady.value = true
    return
  }

  const img = new Image()
  img.onload = () => {
    largeReady.value = true
  }
  img.src = photo.value.largeUrl
}

// ── SEO ──
useSeoMeta({
  title: formatPhotoName(),
  ogImage: photo.value?.largeUrl,
  description: photo.value?.description,
})

// Close panel on Escape
if (import.meta.client) {
  useEventListener(document, 'keydown', (e: KeyboardEvent) => {
    if (e.key === 'Escape') {
      if (detailOpen.value) {
        detailOpen.value = false
      } else {
        navigateTo(localePath('/gallery'))
      }
    }
  })
}

</script>

<style scoped>
.photo-viewer {
  position: fixed;
  inset: 0;
  z-index: 100;
  display: flex;
  background: var(--bg);
}

/* ── Close button ── */
.viewer-close {
  position: fixed;
  top: 16px;
  right: 16px;
  z-index: 110;
  width: 40px;
  height: 40px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 50%;
  background: color-mix(in srgb, var(--surface) 72%, transparent);
  backdrop-filter: blur(12px);
  -webkit-backdrop-filter: blur(12px);
  border: 1px solid color-mix(in srgb, var(--border) 40%, transparent);
  color: var(--text);
  text-decoration: none;
  transition: all var(--dur-fast);
  cursor: pointer;
}

.viewer-close:hover {
  background: var(--surface);
  box-shadow: var(--shadow-md);
}

.viewer-close svg {
  width: 20px;
  height: 20px;
}

/* ── Mobile details toggle ── */
.viewer-details-toggle {
  display: none;
  position: fixed;
  top: 16px;
  left: 16px;
  z-index: 110;
  padding: 8px 16px;
  border-radius: 999px;
  background: color-mix(in srgb, var(--surface) 72%, transparent);
  backdrop-filter: blur(12px);
  -webkit-backdrop-filter: blur(12px);
  border: 1px solid color-mix(in srgb, var(--border) 40%, transparent);
  color: var(--text);
  font-size: 13px;
  font-weight: 500;
  cursor: pointer;
  transition: all var(--dur-fast);
}

.viewer-details-toggle:hover {
  background: var(--surface);
}

/* ── Main image area ── */
.viewer-main {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 24px;
  overflow: hidden;
  min-width: 0;
}

.viewer-stage {
  position: relative;
  max-width: 100%;
  max-height: 100%;
  border-radius: var(--radius-md);
  overflow: hidden;
  background: var(--bg-secondary);
  box-shadow: var(--shadow-lg);
}

.viewer-placeholder {
  position: absolute;
  inset: 0;
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.viewer-image {
  display: block;
  width: 100%;
  height: auto;
  max-width: 100%;
  max-height: 100%;
  object-fit: contain;
}

.viewer-image--large {
  position: absolute;
  inset: 0;
  width: 100%;
  height: 100%;
}

.viewer-loading {
  display: flex;
  align-items: center;
  justify-content: center;
}

.spinner {
  width: 28px;
  height: 28px;
  border: 2px solid var(--border);
  border-top-color: var(--accent);
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

/* ── Detail panel ── */
.viewer-panel {
  width: 360px;
  flex-shrink: 0;
  overflow-y: auto;
  transition: transform var(--dur-slow) var(--ease-smooth);
}

.panel-scroll {
  padding: 32px 28px 36px;
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.panel-heading {
  margin: 0;
  font-size: 1.2rem;
  font-weight: 800;
  letter-spacing: -0.03em;
  color: var(--text);
}

.panel-desc {
  margin: -14px 0 0;
  font-size: 14px;
  line-height: 1.7;
  color: var(--text-muted);
}

.detail-section {
  padding: 0;
}

.section-title {
  margin: 0 0 12px;
  font-size: 0.82rem;
  font-weight: 800;
  letter-spacing: 0.02em;
  color: var(--text);
}

.detail-list {
  margin: 0;
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.detail-row {
  display: grid;
  grid-template-columns: minmax(96px, 118px) minmax(0, 1fr);
  gap: 12px;
  align-items: baseline;
}

.detail-row dt,
.detail-row dd {
  margin: 0;
}

.detail-row dt {
  font-size: 11px;
  font-weight: 600;
  color: var(--text-subtle);
}

.detail-row dd {
  font-size: 13px;
  font-weight: 500;
  color: var(--text);
  word-break: break-word;
}

/* ── Tags ── */
.panel-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  padding-top: 2px;
}

.panel-tag {
  padding: 0;
  font-size: 12px;
  font-weight: 600;
  color: var(--text-muted);
  text-decoration: none;
  transition: all var(--dur-fast);
}

.panel-tag:hover {
  color: var(--accent-text);
}

/* ── Responsive ── */
@media (max-width: 768px) {
  .photo-viewer {
    flex-direction: column;
  }

  .viewer-details-toggle {
    display: block;
  }

  .viewer-main {
    flex: 1;
    padding: 60px 12px 12px;
  }

  .viewer-panel {
    position: fixed;
    bottom: 0;
    left: 0;
    right: 0;
    width: 100%;
    height: 70vh;
    border-radius: var(--radius-xl) var(--radius-xl) 0 0;
    transform: translateY(100%);
    z-index: 120;
  }

  .viewer-panel.open {
    transform: translateY(0);
  }
}
</style>
