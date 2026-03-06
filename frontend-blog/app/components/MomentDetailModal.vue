<!-- File: app/components/MomentDetailModal.vue
     Purpose: Render the URL-driven moment detail modal with focused scrolling and comment thread.
     Module: frontend-blog/components, presentation layer.
     Related: app/pages/moments/[[id]].vue and components/CommentList.vue. -->
<template>
  <Teleport to="body">
    <Transition name="moment-modal">
      <div v-if="open" class="moment-modal" @click.self="$emit('close')">
        <div ref="dialogRef" class="moment-dialog" role="dialog" aria-modal="true" :aria-label="t('moments.detailTitle')" tabindex="-1">
          <button class="close-btn" type="button" :aria-label="t('moments.closeDetail')" @click="$emit('close')">
            ×
          </button>

          <div v-if="moment" class="dialog-body">
            <div class="dialog-surface">
              <div class="dialog-meta">
                <span class="meta-pill">{{ formatDate(moment.publishedAt || moment.createdAt) }}</span>
                <span class="meta-pill">{{ t('moments.commentCount', { n: commentCount }) }}</span>
              </div>

              <div class="dialog-content markdown-body" v-html="renderMomentContent(moment.content)"></div>

              <div class="dialog-nav" :class="{ 'single-side': !previousMoment || !nextMoment }">
                <button
                  class="nav-btn"
                  type="button"
                  :disabled="!previousMoment"
                  @click="$emit('prev')"
                >
                  <span class="nav-arrow" aria-hidden="true">←</span>
                  <span class="nav-label">{{ t('moments.previous') }}</span>
                </button>
                <button
                  class="nav-btn align-right"
                  type="button"
                  :disabled="!nextMoment"
                  @click="$emit('next')"
                >
                  <span class="nav-label">{{ t('moments.next') }}</span>
                  <span class="nav-arrow" aria-hidden="true">→</span>
                </button>
              </div>

              <div v-if="commentEnabled" class="dialog-comments">
                <CommentList
                  :key="moment.id"
                  content-type="moment"
                  :content-id="moment.id"
                  :require-approval="requireApproval"
                  :show-loading-skeleton="false"
                  @count-change="$emit('countChange', $event)"
                />
              </div>
            </div>
          </div>

          <div v-else-if="loading" class="dialog-loading">
            <div class="spinner"></div>
            <span>{{ t('comment.loading') }}</span>
          </div>

          <div v-else class="dialog-empty">
            {{ t('moments.detailMissing') }}
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<script setup lang="ts">
import { computed, nextTick, ref, watch, onBeforeUnmount } from 'vue'
import { useI18n } from 'vue-i18n'
import CommentList from '~/components/CommentList.vue'
import type { Moment } from '~/composables/useApi'
import { renderContentMarkdown } from '~/utils/contentMarkdown'

const props = defineProps<{
  open: boolean
  moment?: Moment | null
  loading?: boolean
  commentEnabled?: boolean
  requireApproval?: boolean
  previousMoment?: Moment | null
  nextMoment?: Moment | null
}>()

defineEmits<{
  (e: 'close'): void
  (e: 'countChange', total: number): void
  (e: 'prev'): void
  (e: 'next'): void
}>()

const { t, locale } = useI18n()
const commentCount = computed(() => props.moment?.commentCount || 0)
const dialogRef = ref<HTMLElement | null>(null)
let previousBodyOverflow = ''

watch(() => props.open, async (open) => {
  if (import.meta.server) return
  if (open) {
    previousBodyOverflow = document.body.style.overflow
    document.body.style.overflow = 'hidden'
    await nextTick()
    if (dialogRef.value) {
      dialogRef.value.scrollTop = 0
    }
    dialogRef.value?.focus()
    return
  }
  document.body.style.overflow = previousBodyOverflow
}, { immediate: true })

onBeforeUnmount(() => {
  if (import.meta.client) {
    document.body.style.overflow = previousBodyOverflow
  }
})

function formatDate(iso: string): string {
  return new Intl.DateTimeFormat(locale.value === 'en' ? 'en-US' : 'zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
  }).format(new Date(iso))
}

function renderMomentContent(content: string): string {
  return renderContentMarkdown(content)
}
</script>

<style scoped>
.moment-modal {
  position: fixed;
  inset: 0;
  z-index: 70;
  display: flex;
  align-items: flex-start;
  justify-content: center;
  padding: 72px 24px 24px;
  background: color-mix(in srgb, #081018 38%, transparent);
  backdrop-filter: blur(8px);
}

.moment-dialog {
  position: relative;
  width: min(760px, 100%);
  max-height: calc(100vh - 96px);
  overflow: auto;
  overscroll-behavior: contain;
  padding: 0;
  border: none;
  background: transparent;
  box-shadow: none;
  outline: none;
}

.close-btn {
  position: absolute;
  top: 14px;
  right: 14px;
  width: 40px;
  height: 40px;
  border: none;
  border-radius: 999px;
  background: color-mix(in srgb, var(--bg-secondary) 82%, transparent);
  color: var(--text-muted);
  font-size: 22px;
  cursor: pointer;
}

.dialog-meta {
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
  margin-bottom: 18px;
}

.dialog-surface {
  padding: 22px 22px 26px;
  border: 1px solid color-mix(in srgb, var(--border) 88%, transparent);
  border-radius: 22px;
  background: #ffffff;
  box-shadow: 0 10px 24px rgba(8, 16, 24, 0.06);
}

.meta-pill {
  display: inline-flex;
  align-items: center;
  min-height: 30px;
  padding: 0 12px;
  border-radius: 999px;
  background: color-mix(in srgb, var(--accent) 7%, var(--bg-secondary));
  color: var(--text-subtle);
  font-size: 12px;
  letter-spacing: 0.04em;
}

.dialog-content {
  font-size: 16px;
  line-height: 1.95;
  color: var(--text);
  word-break: break-word;
}

.dialog-content :deep(p),
.dialog-content :deep(ul),
.dialog-content :deep(ol),
.dialog-content :deep(blockquote),
.dialog-content :deep(pre) {
  margin: 0 0 14px;
}

.dialog-content :deep(p:last-child),
.dialog-content :deep(ul:last-child),
.dialog-content :deep(ol:last-child),
.dialog-content :deep(blockquote:last-child),
.dialog-content :deep(pre:last-child) {
  margin-bottom: 0;
}

.dialog-content :deep(ul),
.dialog-content :deep(ol) {
  padding-left: 22px;
}

.dialog-content :deep(blockquote) {
  padding-left: 14px;
  border-left: 3px solid var(--border);
  color: var(--text-muted);
}

.dialog-content :deep(pre) {
  overflow-x: auto;
  padding: 14px 16px;
  border-radius: 16px;
  background: #f6f7f8;
}

.dialog-content :deep(code) {
  font-family: 'Fira Code', monospace;
  font-size: 0.92em;
}

.dialog-content :deep(pre code) {
  background: transparent;
}

.dialog-nav {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  justify-content: space-between;
  gap: 10px 16px;
  margin-top: 28px;
}

.dialog-nav.single-side {
  justify-content: flex-start;
}

.nav-btn {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  min-height: auto;
  padding: 0;
  border: none;
  background: transparent;
  color: var(--accent-text);
  text-align: left;
  cursor: pointer;
  transition:
    transform 220ms cubic-bezier(0.22, 1.18, 0.36, 1),
    color 180ms ease,
    opacity 180ms ease;
}

.nav-btn:hover:not(:disabled) {
  transform: translateY(-1px);
  color: var(--accent);
}

.nav-btn:disabled {
  cursor: default;
  opacity: 0.34;
}

.align-right {
  text-align: right;
  margin-left: auto;
}

.nav-label {
  font-size: 13px;
  line-height: 1.4;
  font-weight: 600;
}

.nav-arrow {
  font-size: 14px;
  line-height: 1;
}

.dialog-comments {
  margin-top: 28px;
}

.dialog-loading,
.dialog-empty {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 10px;
  min-height: 240px;
  color: var(--text-subtle);
}

.spinner {
  width: 22px;
  height: 22px;
  border: 2px solid var(--border);
  border-top-color: var(--accent);
  border-radius: 50%;
  animation: moment-modal-spin 0.8s linear infinite;
}

.moment-modal-enter-active,
.moment-modal-leave-active {
  transition: opacity 220ms ease;
}

.moment-modal-enter-active .moment-dialog,
.moment-modal-leave-active .moment-dialog {
  transition: transform 340ms cubic-bezier(0.22, 1.18, 0.36, 1), opacity 220ms ease;
}

.moment-modal-enter-from,
.moment-modal-leave-to {
  opacity: 0;
}

.moment-modal-enter-from .moment-dialog,
.moment-modal-leave-to .moment-dialog {
  opacity: 0;
  transform: translateY(28px) scale(0.98);
}

@keyframes moment-modal-spin {
  to {
    transform: rotate(360deg);
  }
}

@media (max-width: 640px) {
  .moment-modal {
    padding: 56px 12px 12px;
  }

  .moment-dialog {
    max-height: calc(100vh - 68px);
  }

  .dialog-surface {
    padding: 18px 16px 22px;
    border-radius: 18px;
  }

  .dialog-nav {
    justify-content: flex-start;
  }
}
</style>
