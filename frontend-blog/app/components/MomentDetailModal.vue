<!-- File: app/components/MomentDetailModal.vue
     Purpose: Render a URL-driven moment detail modal with comments.
     Module: frontend-blog/components, presentation layer.
     Related: app/pages/moments/[[id]].vue and components/CommentList.vue. -->
<template>
  <Teleport to="body">
    <Transition name="moment-modal">
      <div v-if="open" class="moment-modal" @click.self="$emit('close')">
        <div class="moment-dialog" role="dialog" aria-modal="true" :aria-label="t('moments.detailTitle')">
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
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import CommentList from '~/components/CommentList.vue'
import type { Moment } from '~/composables/useApi'

const props = defineProps<{
  open: boolean
  moment?: Moment | null
  loading?: boolean
  requireApproval?: boolean
}>()

defineEmits<{
  (e: 'close'): void
  (e: 'countChange', total: number): void
}>()

const { t, locale } = useI18n()

const commentCount = computed(() => props.moment?.commentCount || 0)

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
  align-items: flex-end;
  justify-content: center;
  padding: 24px;
  background: color-mix(in srgb, #081018 38%, transparent);
  backdrop-filter: blur(8px);
}

.moment-dialog {
  position: relative;
  width: min(760px, 100%);
  max-height: min(88vh, 920px);
  overflow: auto;
  padding: 28px;
  border-radius: 28px;
  background: color-mix(in srgb, var(--bg-primary) 94%, white);
  box-shadow: 0 24px 80px color-mix(in srgb, #081018 18%, transparent);
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
    padding: 12px;
  }

  .moment-dialog {
    padding: 22px 18px;
    border-radius: 24px;
  }
}
</style>
