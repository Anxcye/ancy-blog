<!-- File: app/error.vue
     Purpose: Custom boundary for 404s and server errors.
     Module: frontend-blog/app, error boundary.
-->
<template>
  <div class="error-page" :class="{ dark: isDark }">
    <div class="error-content">
      <h1 class="error-code word-bounce">{{ error?.statusCode || 500 }}</h1>
      <p class="error-msg word-bounce" style="animation-delay: 100ms">
        {{ is404 ? '这里空空如也，或许像一些曾经的念想。' : '抱歉，服务器的脑回路短路了。' }}
      </p>
      
      <p v-if="error?.message && !is404" class="error-detail word-bounce" style="animation-delay: 200ms">
        {{ error.message }}
      </p>
      
      <button class="error-btn word-bounce" style="animation-delay: 300ms" @click="handleError">
        回到起点
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { NuxtError } from '#app'

const props = defineProps({
  error: Object as () => NuxtError
})

const colorMode = useColorMode()
const isDark = computed(() => colorMode.value === 'dark')

const is404 = computed(() => props.error?.statusCode === 404 || props.error?.statusCode === '404')
const handleError = () => clearError({ redirect: '/' })

useHead({ title: is404.value ? '404 - 迷失了方向' : 'Oops... 故障发生' })
</script>

<style scoped>
.error-page {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--bg, #f9f9f7);
  color: var(--text, #1c1e21);
  text-align: center;
  padding: 24px;
}

.error-content {
  max-width: 480px;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 20px;
}

.error-code {
  font-size: clamp(6rem, 20vw, 10rem);
  font-weight: 800;
  line-height: 1;
  background: linear-gradient(135deg, var(--text) 0%, var(--accent-soft, #14b8a6) 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  margin: 0;
}

.error-msg {
  font-size: 1.15rem;
  font-weight: 500;
  color: var(--text-muted, #64748b);
  margin: 0;
  line-height: 1.6;
}

.error-detail {
  font-size: 0.875rem;
  color: #ef4444;
  background: rgba(239, 68, 68, 0.1);
  padding: 8px 16px;
  border-radius: var(--radius-sm);
  font-family: monospace;
  word-break: break-all;
}

.error-btn {
  margin-top: 24px;
  padding: 14px 32px;
  font-size: 15px;
  font-weight: 700;
  color: var(--bg);
  background: var(--text);
  border: none;
  border-radius: 99px;
  cursor: pointer;
  transition: transform 0.2s cubic-bezier(0.34, 1.56, 0.64, 1), opacity 0.2s, background 0.2s, color 0.2s;
}

.error-btn:hover {
  transform: scale(1.05) translateY(-2px);
  background: var(--accent);
  color: var(--accent-text);
}

.error-btn:active {
  transform: scale(0.95);
}

.word-bounce {
  display: inline-block;
  opacity: 0;
  animation: word-spring 0.8s cubic-bezier(0.34, 1.56, 0.64, 1) forwards;
}

@keyframes word-spring {
  0% { opacity: 0; transform: translateY(20px) scale(0.8); }
  60% { opacity: 1; transform: translateY(-4px) scale(1.05); }
  100% { opacity: 1; transform: translateY(0) scale(1); }
}
</style>
