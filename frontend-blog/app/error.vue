<!-- File: app/error.vue
     Purpose: Global error boundary for SSR/Client crashes and 404/500 API errors.
     Module: frontend-blog/app
-->
<template>
  <div class="error-page theme-adapt">
    <div class="error-content">
      <div class="error-code">
        {{ error?.statusCode }}
      </div>
      <h1 class="error-title">
        {{ is404 ? '页面好像走丢了' : '出了些小状况' }}
      </h1>
      <p class="error-message">
        {{ is404 ? '您访问的地址可能已被移除，或暂时不可用。' : error?.message || '服务器遇到了一点问题，请稍后再试。' }}
      </p>

      <div class="error-actions">
        <UButton
          size="lg"
          color="gray"
          variant="ghost"
          icon="i-heroicons-arrow-left"
          @click="handleError"
        >
          返回首页
        </UButton>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
const props = defineProps({
  error: Object as () => any
})

const is404 = computed(() => props.error?.statusCode === 404 || props.error?.statusCode === '404')

const handleError = () => clearError({ redirect: '/' })

// Basic meta tags to look good even when layout crashes
useHead({ title: is404.value ? '404 - 未找到页面' : 'Oops... 出错了' })
</script>

<style scoped>
.error-page {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--bg);
  color: var(--text);
  padding: 24px;
  font-family: inherit;
}

.error-content {
  text-align: center;
  max-width: 480px;
  width: 100%;
}

.error-code {
  font-size: 100px;
  font-weight: 800;
  color: var(--accent);
  line-height: 1;
  margin-bottom: 24px;
  opacity: 0.15;
  letter-spacing: -2px;
}

.error-title {
  font-size: 28px;
  font-weight: 700;
  margin-bottom: 12px;
}

.error-message {
  font-size: 15px;
  color: var(--text-muted);
  line-height: 1.6;
  margin-bottom: 32px;
}

.error-actions {
  display: flex;
  justify-content: center;
}
</style>
