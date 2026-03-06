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
            <div class="dialog-meta">
              <span class="meta-pill">{{ formatDate(moment.publishedAt || moment.createdAt) }}</span>
              <span class="meta-pill">{{ t('moments.commentCount', { n: commentCount }) }}</span>
            </div>

            <div class="dialog-content">
              <p>{{ moment.content }}</p>
            </div>

            <div class="dialog-nav" :class="{ 'single-side': !previousMoment || !nextMoment }">
              <button
                class="nav-btn"
                type="button"
                :disabled="!previousMoment"
                @click="$emit('prev')"
              >
                <span class="nav-label">{{ t('moments.previous') }}</span>
                <span v-if="previousMoment" class="nav-text">{{ previousMoment.content }}</span>
              </button>
              <button
                class="nav-btn align-right"
                type="button"
                :disabled="!nextMoment"
                @click="$emit('next')"
              >
                <span class="nav-label">{{ t('moments.next') }}</span>
                <span v-if="nextMoment" class="nav-text">{{ nextMoment.content }}</span>
              </button>
            </div>

            <div v-if="moment.allowComment" class="dialog-comments">
              <CommentList
                content-type="moment"
                :content-id="moment.id"
                :require-approval="requireApproval"
                @count-change="$emit('countChange', $event)"
              />
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

const props = defineProps<{
  open: boolean
  moment?: Moment | null
  loading?: boolean
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
  padding: 28px;
  border-radius: 28px;
  background: color-mix(in srgb, var(--bg-primary) 94%, white);
  box-shadow: 0 24px 80px color-mix(in srgb, #081018 18%, transparent);
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
  white-space: pre-wrap;
  word-break: break-word;
}

.dialog-nav {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 12px;
  margin-top: 28px;
}

.dialog-nav.single-side {
  grid-template-columns: 1fr;
}

.nav-btn {
  display: flex;
  flex-direction: column;
  gap: 6px;
  min-height: 72px;
  padding: 14px 16px;
  border: 1px solid var(--border);
  border-radius: 18px;
  background: color-mix(in srgb, var(--bg-secondary) 44%, transparent);
  color: var(--text);
  text-align: left;
  cursor: pointer;
  transition:
    transform 260ms cubic-bezier(0.22, 1.18, 0.36, 1),
    border-color 180ms ease,
    background 180ms ease;
}

.nav-btn:hover:not(:disabled) {
  transform: translateY(-2px);
  border-color: color-mix(in srgb, var(--accent) 35%, var(--border));
}

.nav-btn:disabled {
  cursor: default;
  opacity: 0.48;
}

.align-right {
  text-align: right;
}

.nav-label {
  font-size: 12px;
  letter-spacing: 0.08em;
  text-transform: uppercase;
  color: var(--text-subtle);
}

.nav-text {
  display: -webkit-box;
  overflow: hidden;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  line-height: 1.7;
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
    padding: 22px 18px;
    border-radius: 24px;
  }

  .dialog-nav {
    grid-template-columns: 1fr;
  }
}
</style>
