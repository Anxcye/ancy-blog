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
      <img
        v-if="photo"
        :src="photo.largeUrl"
        :alt="photo.title || photo.slug"
        :width="photo.width"
        :height="photo.height"
        class="viewer-image"
        @load="imageLoaded = true"
      />
      <div v-if="!imageLoaded" class="viewer-loading">
        <div class="spinner" />
      </div>
    </div>

    <!-- Detail panel -->
    <aside class="viewer-panel" :class="{ open: detailOpen }">
      <div class="panel-scroll" v-if="photo">
        <!-- Title -->
        <h1 v-if="photo.title" class="panel-title">{{ photo.title }}</h1>

        <!-- Description -->
        <p v-if="photo.description" class="panel-desc">{{ photo.description }}</p>

        <!-- Metadata rows -->
        <div class="panel-meta">
          <div v-if="photo.locationName || photo.locationCity" class="meta-row">
            <svg class="meta-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M21 10c0 7-9 13-9 13s-9-6-9-13a9 9 0 0 1 18 0z"/><circle cx="12" cy="10" r="3"/></svg>
            <div class="meta-content">
              <span class="meta-label">{{ t('gallery.location') }}</span>
              <span class="meta-value">{{ photo.locationName || photo.locationCity }}</span>
            </div>
          </div>

          <div v-if="photo.takenAt" class="meta-row">
            <svg class="meta-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><rect x="3" y="4" width="18" height="18" rx="2" ry="2"/><line x1="16" y1="2" x2="16" y2="6"/><line x1="8" y1="2" x2="8" y2="6"/><line x1="3" y1="10" x2="21" y2="10"/></svg>
            <div class="meta-content">
              <span class="meta-label">{{ t('gallery.takenAt') }}</span>
              <span class="meta-value">{{ formatDate(photo.takenAt) }}</span>
            </div>
          </div>

          <div v-if="photo.cameraMake || photo.cameraModel" class="meta-row">
            <svg class="meta-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M23 19a2 2 0 0 1-2 2H3a2 2 0 0 1-2-2V8a2 2 0 0 1 2-2h4l2-3h6l2 3h4a2 2 0 0 1 2 2z"/><circle cx="12" cy="13" r="4"/></svg>
            <div class="meta-content">
              <span class="meta-label">{{ t('gallery.camera') }}</span>
              <span class="meta-value">{{ [photo.cameraMake, photo.cameraModel].filter(Boolean).join(' ') }}</span>
            </div>
          </div>

          <div v-if="photo.lensModel" class="meta-row">
            <svg class="meta-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="12" cy="12" r="10"/><circle cx="12" cy="12" r="3"/></svg>
            <div class="meta-content">
              <span class="meta-label">{{ t('gallery.lens') }}</span>
              <span class="meta-value">{{ photo.lensModel }}</span>
            </div>
          </div>
        </div>

        <!-- EXIF specs row -->
        <div v-if="hasExif" class="exif-grid">
          <div v-if="photo.focalLength" class="exif-item">
            <span class="exif-label">{{ t('gallery.focalLength') }}</span>
            <span class="exif-value">{{ photo.focalLength }}</span>
          </div>
          <div v-if="photo.aperture" class="exif-item">
            <span class="exif-label">{{ t('gallery.aperture') }}</span>
            <span class="exif-value">{{ photo.aperture }}</span>
          </div>
          <div v-if="photo.shutterSpeed" class="exif-item">
            <span class="exif-label">{{ t('gallery.shutterSpeed') }}</span>
            <span class="exif-value">{{ photo.shutterSpeed }}</span>
          </div>
          <div v-if="photo.iso" class="exif-item">
            <span class="exif-label">{{ t('gallery.iso') }}</span>
            <span class="exif-value">{{ photo.iso }}</span>
          </div>
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
import type { GalleryTag } from '~/composables/useApi'

const { t, locale } = useI18n()
const localePath = useLocalePath()
const route = useRoute()
const { getGalleryPhoto, getGalleryTags } = useApi()

const slug = route.params.slug as string
const detailOpen = ref(false)
const imageLoaded = ref(false)

// Fetch photo and tags in parallel
const [{ data: photo }, { data: galleryTags }] = await Promise.all([
  useAsyncData(`gallery-photo-${slug}`, () => getGalleryPhoto(slug)),
  useAsyncData('gallery-tags-viewer', getGalleryTags, { default: () => [] as GalleryTag[] }),
])

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

// ── SEO ──
useSeoMeta({
  title: photo.value?.title || slug,
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

.viewer-image {
  max-width: 100%;
  max-height: 100%;
  object-fit: contain;
  border-radius: var(--radius-md);
  box-shadow: var(--shadow-lg);
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
  border-left: 1px solid color-mix(in srgb, var(--border) 60%, transparent);
  background: color-mix(in srgb, var(--surface) 88%, transparent);
  backdrop-filter: blur(20px) saturate(1.4);
  -webkit-backdrop-filter: blur(20px) saturate(1.4);
  overflow-y: auto;
  transition: transform var(--dur-slow) var(--ease-smooth);
}

.panel-scroll {
  padding: 32px 24px;
}

.panel-title {
  font-size: 1.3rem;
  font-weight: 700;
  margin: 0 0 12px;
  letter-spacing: -0.01em;
}

.panel-desc {
  font-size: 14px;
  line-height: 1.7;
  color: var(--text-muted);
  margin: 0 0 24px;
}

/* ── Metadata rows ── */
.panel-meta {
  display: flex;
  flex-direction: column;
  gap: 16px;
  margin-bottom: 24px;
}

.meta-row {
  display: flex;
  align-items: flex-start;
  gap: 12px;
}

.meta-icon {
  width: 18px;
  height: 18px;
  flex-shrink: 0;
  margin-top: 2px;
  color: var(--text-subtle);
}

.meta-content {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.meta-label {
  font-size: 11px;
  font-weight: 600;
  letter-spacing: 0.04em;
  text-transform: uppercase;
  color: var(--text-subtle);
}

.meta-value {
  font-size: 14px;
  color: var(--text);
}

/* ── EXIF grid ── */
.exif-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 12px;
  padding: 16px;
  border-radius: var(--radius-md);
  background: color-mix(in srgb, var(--bg-secondary) 60%, transparent);
  border: 1px solid color-mix(in srgb, var(--border) 40%, transparent);
  margin-bottom: 24px;
}

.exif-item {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.exif-label {
  font-size: 10px;
  font-weight: 600;
  letter-spacing: 0.05em;
  text-transform: uppercase;
  color: var(--text-subtle);
}

.exif-value {
  font-size: 14px;
  font-weight: 600;
  color: var(--text);
}

/* ── Tags ── */
.panel-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.panel-tag {
  padding: 4px 12px;
  border-radius: 999px;
  border: 1px solid var(--border);
  font-size: 12px;
  font-weight: 500;
  color: var(--text-muted);
  text-decoration: none;
  transition: all var(--dur-fast);
}

.panel-tag:hover {
  background: var(--accent-soft);
  border-color: var(--accent);
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
    border-left: none;
    border-top: 1px solid color-mix(in srgb, var(--border) 60%, transparent);
    border-radius: var(--radius-xl) var(--radius-xl) 0 0;
    transform: translateY(100%);
    z-index: 120;
  }

  .viewer-panel.open {
    transform: translateY(0);
  }
}
</style>
