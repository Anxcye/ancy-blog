<!-- File: app/pages/friends.vue
     Purpose: Friends/Links page showing approved links and an optional pinned article introduction.
     Module: frontend-blog/pages
-->
<template>
  <div class="friends-page">
    <div class="container">
      <div class="page-hero">
        <span class="hero-eyebrow">{{ t('friends.eyebrow') }}</span>
        <div class="hero-main">
          <div class="hero-copy">
            <h1 class="page-title">{{ t('friends.title') }}</h1>
            <p class="page-subtitle">{{ t('friends.subtitle') }}</p>
          </div>
          <div class="hero-stats">
            <span class="hero-stat">{{ t('friends.total', { n: links?.length || 0 }) }}</span>
            <span class="hero-stat muted">{{ t('friends.openSubmission') }}</span>
          </div>
        </div>
      </div>

      <!-- Skeleton Loading -->
      <div v-if="pending" class="links-grid">
        <div v-for="n in 6" :key="n" class="link-card skeleton-card">
          <div class="skeleton-avatar"></div>
          <div class="skeleton-info">
            <div class="skeleton-line" style="width: 50%; height: 16px;"></div>
            <div class="skeleton-line" style="width: 80%; height: 12px; margin-top: 8px;"></div>
          </div>
        </div>
      </div>

      <!-- Links Grid -->
      <div v-else-if="links?.length" class="links-grid">
        <a
          v-for="(link, i) in links"
          :key="link.id"
          :href="link.url"
          target="_blank"
          rel="noopener noreferrer"
          class="link-card"
          :style="{ animationDelay: `${i * 60 + 200}ms` }"
        >
          <div class="link-avatar">
            <img v-if="link.avatarUrl" :src="link.avatarUrl" :alt="link.name" loading="lazy" />
            <span v-else class="link-fallback">{{ link.name.charAt(0).toUpperCase() }}</span>
          </div>
          <div class="link-info">
            <h3 class="link-name">{{ link.name }}</h3>
            <p class="link-desc" :title="link.description">{{ link.description || '这人很懒，什么都没留下。' }}</p>
          </div>
        </a>
      </div>

      <div v-else class="empty-state">
        目前还没有记录任何宇宙信号...
      </div>

      <!-- Optional Article Intro -->
      <TiptapRenderer v-if="article?.content" :content="article.content" class="friends-intro" />

      <!-- Submission Form -->
      <div class="submit-section">
        <h2 class="submit-title">申请友链</h2>
        <div class="submit-container">
          <!-- Preview -->
          <div class="preview-area">
            <p class="preview-label">预览效果</p>
            <div class="link-card preview-card">
              <div class="link-avatar">
                <img v-if="form.avatarUrl" :src="form.avatarUrl" :alt="form.name" />
                <span v-else class="link-fallback">{{ form.name ? form.name.charAt(0).toUpperCase() : '?' }}</span>
              </div>
              <div class="link-info">
                <h3 class="link-name">{{ form.name || '站点名称' }}</h3>
                <p class="link-desc">{{ form.description || '站点简介' }}</p>
              </div>
            </div>
          </div>

          <!-- Form -->
          <form @submit.prevent="handleSubmit" class="submit-form">
            <input v-model="form.name" type="text" placeholder="站点名称 *" required class="form-input" />
            <input v-model="form.url" type="url" placeholder="站点链接 *" required class="form-input" />
            <input v-model="form.avatarUrl" type="url" placeholder="头像链接" class="form-input" />
            <textarea v-model="form.description" placeholder="站点简介" rows="2" class="form-textarea"></textarea>
            <input v-model="form.contactEmail" type="email" placeholder="联系邮箱" class="form-input" />
            <button type="submit" :disabled="submitting" class="submit-btn">
              {{ submitting ? '提交中...' : '提交申请' }}
            </button>
            <p v-if="submitMessage" class="submit-message" :class="{ success: submitSuccess }">{{ submitMessage }}</p>
          </form>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
const { t } = useI18n()
const { getApprovedLinks, getArticle, submitLink } = useApi()

const { data: links, pending } = await useAsyncData('friends-links', getApprovedLinks, {
  getCachedData: () => undefined
})

const { data: article } = await useAsyncData('friends-intro', async () => {
  try {
    return await getArticle('friends')
  } catch (err: any) {
    if (err.statusCode === 404) return null
    throw err
  }
}, { getCachedData: () => undefined })

const form = ref({ name: '', url: '', avatarUrl: '', description: '', contactEmail: '' })
const submitting = ref(false)
const submitMessage = ref('')
const submitSuccess = ref(false)

async function handleSubmit() {
  submitting.value = true
  submitMessage.value = ''
  try {
    await submitLink(form.value)
    submitSuccess.value = true
    submitMessage.value = '提交成功！等待审核通过后即可显示。'
    form.value = { name: '', url: '', avatarUrl: '', description: '', contactEmail: '' }
  } catch (err: any) {
    submitSuccess.value = false
    submitMessage.value = err.message || '提交失败，请稍后重试。'
  } finally {
    submitting.value = false
  }
}

useSeoMeta({ title: () => `${t('friends.title')} - ${t('nav.links')}` })
</script>

<style scoped>
.friends-page {
  padding-top: calc(var(--header-h) + 64px);
  padding-bottom: 80px;
}

.page-hero {
  display: flex;
  flex-direction: column;
  gap: 16px;
  margin-bottom: 44px;
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

.hero-copy {
  max-width: 620px;
}

.page-title {
  font-size: clamp(1.8rem, 4vw, 2.8rem);
  font-weight: 800;
  letter-spacing: -0.02em;
  margin: 0;
}

.page-subtitle {
  margin: 10px 0 0;
  max-width: 560px;
  font-size: 15px;
  line-height: 1.8;
  color: var(--text-subtle);
}

.hero-stats {
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
}

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

.hero-stat.muted {
  color: var(--text-subtle);
}

.friends-intro {
  margin-bottom: 56px;
  padding: 32px;
  background: var(--surface);
  border: 1px solid var(--border);
  border-radius: var(--radius-lg);
  box-shadow: var(--shadow-sm);
  animation: fade-up 0.6s var(--ease-spring) both;
}

@keyframes fade-up {
  from { opacity: 0; transform: translateY(20px); }
  to { opacity: 1; transform: translateY(0); }
}

.links-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
  gap: 24px;
}

.link-card {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 12px;
  padding: 20px;
  text-decoration: none;
  color: var(--text);
  transition: all var(--dur-base) var(--ease-out);
  opacity: 0;
  animation: fade-up 0.5s var(--ease-spring) forwards;
  text-align: center;
}

.link-card:hover {
  transform: translateY(-4px);
}

.link-avatar {
  width: 64px;
  height: 64px;
  flex-shrink: 0;
  border-radius: 50%;
  overflow: hidden;
  background: var(--bg-secondary);
  border: 2px solid var(--border);
  display: grid;
  place-items: center;
  transition: transform var(--dur-base);
}

.link-card:hover .link-avatar {
  transform: scale(1.1);
  border-color: var(--accent);
}

.link-avatar img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.link-fallback {
  font-size: 24px;
  font-weight: 700;
  color: var(--text-muted);
}

.link-info {
  flex: 1;
  min-width: 0;
}

.link-name {
  font-size: 15px;
  font-weight: 600;
  margin: 0 0 6px;
  color: var(--text);
  transition: color var(--dur-fast);
}

.link-card:hover .link-name {
  color: var(--accent-text);
}

.link-desc {
  font-size: 13px;
  color: var(--text-muted);
  margin: 0;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
  line-height: 1.5;
}

/* Skeleton */
.skeleton-card { animation: none; opacity: 1; border-color: transparent; }
.skeleton-avatar { width: 56px; height: 56px; border-radius: 50%; background: var(--border); }
.skeleton-line { background: var(--border); border-radius: 4px; }
.skeleton-avatar, .skeleton-line {
  background: linear-gradient(90deg, var(--bg-secondary) 25%, var(--surface-hover) 50%, var(--bg-secondary) 75%);
  background-size: 200% 100%;
  animation: shimmer 1.5s infinite;
}

@keyframes shimmer {
  from { background-position: 200% 0; }
  to { background-position: -200% 0; }
}

.empty-state {
  text-align: center;
  padding: 80px 20px;
  color: var(--text-muted);
  font-size: 15px;
}

/* Submission Form */
.submit-section {
  margin-top: 80px;
  padding: 40px;
  border-top: 1px solid var(--border);
}

.submit-title {
  font-size: 1.5rem;
  font-weight: 700;
  margin: 0 0 24px;
  text-align: center;
}

.submit-container {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 32px;
  max-width: 900px;
  margin: 0 auto;
}

.preview-area {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.preview-label {
  font-size: 13px;
  color: var(--text-muted);
  font-weight: 500;
}

.preview-card {
  pointer-events: none;
  opacity: 1;
  animation: none;
}

.submit-form {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.form-input, .form-textarea {
  width: 100%;
  padding: 10px 14px;
  background: var(--bg);
  border: 1px solid var(--border);
  border-radius: var(--radius-md);
  color: var(--text);
  font-size: 14px;
  transition: border-color var(--dur-fast);
}

.form-input:focus, .form-textarea:focus {
  outline: none;
  border-color: var(--accent);
}

.form-textarea {
  resize: vertical;
  font-family: inherit;
}

.submit-btn {
  padding: 10px 20px;
  background: var(--accent);
  color: white;
  border: none;
  border-radius: var(--radius-md);
  font-size: 14px;
  font-weight: 600;
  cursor: pointer;
  transition: all var(--dur-fast);
}

.submit-btn:hover:not(:disabled) {
  background: var(--accent-hover);
}

.submit-btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.submit-message {
  text-align: center;
  padding: 10px;
  border-radius: var(--radius-md);
  font-size: 13px;
  background: var(--surface-hover);
  color: var(--text-muted);
}

.submit-message.success {
  background: var(--accent-soft);
  color: var(--accent-text);
}

@media (max-width: 768px) {
  .page-hero {
    gap: 14px;
    margin-bottom: 34px;
    padding-bottom: 22px;
  }

  .hero-main {
    align-items: flex-start;
    gap: 18px;
  }

  .page-subtitle {
    font-size: 14px;
  }

  .submit-container {
    grid-template-columns: 1fr;
  }
  .submit-section {
    padding: 24px;
  }
}
</style>
