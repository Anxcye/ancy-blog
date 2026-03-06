<!-- File: app/components/InfiniteScrollTrigger.vue
     Purpose: A reusable intersection-observer component to trigger loading more items.
     Module: components, presentation layer. -->
<template>
  <div ref="target" class="infinite-scroll-trigger">
    <slot name="loading" v-if="loading">
      <div class="spinner"></div>
    </slot>
    <slot name="done" v-else-if="done">
      <span class="no-more">{{ doneText || '—— 没有更多内容了 ——' }}</span>
    </slot>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useIntersectionObserver } from '@vueuse/core'

const props = defineProps<{
  loading: boolean
  done: boolean
  doneText?: string
}>()

const emit = defineEmits<{
  (e: 'load'): void
}>()

const target = ref<HTMLElement | null>(null)

useIntersectionObserver(target, ([entry]) => {
  if (entry?.isIntersecting && !props.loading && !props.done) {
    emit('load')
  }
}, {
  rootMargin: '0px', // trigger fetch exactly when it hits the viewport
})
</script>

<style scoped>
.infinite-scroll-trigger {
  display: flex;
  justify-content: center;
  align-items: center;
  padding: 32px 0;
  min-height: 40px;
  color: var(--text-subtle);
  font-size: 13px;
  clear: both;
}

.no-more {
  opacity: 0.6;
  letter-spacing: 0.05em;
}

.spinner {
  width: 22px;
  height: 22px;
  border: 2px solid var(--border);
  border-top-color: var(--accent);
  border-radius: 50%;
  animation: is-spin 0.8s linear infinite;
}

@keyframes is-spin {
  to { transform: rotate(360deg); }
}
</style>
