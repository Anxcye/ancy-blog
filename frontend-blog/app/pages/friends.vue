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
            <span class="hero-stat muted">
              {{ siteSettings?.linkSubmissionEnabled ? t('friends.openSubmission') : t('friends.closedSubmission') }}
            </span>
          </div>
        </div>
      </div>

      <TiptapRenderer v-if="article?.content" :content="article.content" class="friends-intro" />

      <section class="directory-shell">
        <div v-if="pending" class="links-directory">
          <div v-for="n in 6" :key="n" class="link-row skeleton-row">
            <div class="link-row-inner">
              <div class="skeleton-avatar"></div>
              <div class="skeleton-copy">
                <div class="skeleton-line line-title"></div>
                <div class="skeleton-line line-desc"></div>
              </div>
              <div class="skeleton-line line-meta"></div>
            </div>
          </div>
        </div>

        <div v-else-if="links?.length" class="links-directory">
          <a
            v-for="(link, i) in links"
            :key="link.id"
            :href="link.url"
            target="_blank"
            rel="noopener noreferrer"
            class="link-row"
            :style="{ animationDelay: `${i * 45 + 160}ms` }"
            @mousemove="handleLinkMove"
            @mouseleave="resetLinkMove"
          >
            <div class="link-row-inner">
              <div class="link-smoke"></div>
              <div class="link-avatar">
                <img v-if="link.avatarUrl" :src="link.avatarUrl" :alt="link.name" loading="lazy" />
                <span v-else class="link-fallback">{{ getInitial(link.name) }}</span>
              </div>
              <div class="link-primary">
                <div class="link-heading">
                  <h3 class="link-name">{{ link.name }}</h3>
                  <span class="link-host">{{ getHost(link.url) }}</span>
                </div>
                <p class="link-desc" :title="link.description">
                  {{ link.description || t('friends.fallbackDescription') }}
                </p>
              </div>
              <span class="link-tail">{{ t('friends.outboundLabel') }}</span>
            </div>
          </a>
        </div>

        <div v-else class="empty-state">
          {{ t('friends.empty') }}
        </div>
      </section>

      <section class="submit-section">
        <div class="submit-shell">
          <div class="submit-head">
            <div>
              <p class="section-kicker">{{ t('friends.submitEyebrow') }}</p>
              <h2 class="section-title">{{ t('friends.submitTitle') }}</h2>
              <p class="submit-note">
                {{ siteSettings?.linkSubmissionEnabled ? t('friends.submitNote') : t('friends.submitClosedNote') }}
              </p>
            </div>
            <button
              v-if="siteSettings?.linkSubmissionEnabled"
              type="button"
              class="submit-trigger"
              :aria-expanded="showSubmission"
              @click="showSubmission = !showSubmission"
            >
              {{ showSubmission ? t('friends.hideSubmission') : t('friends.showSubmission') }}
            </button>
          </div>

          <ul class="submission-rules">
            <li>{{ t('friends.ruleAvailability') }}</li>
            <li>{{ t('friends.ruleOriginality') }}</li>
            <li>{{ t('friends.ruleReciprocal') }}</li>
          </ul>

          <div v-if="siteSettings?.linkSubmissionEnabled && showSubmission" class="submit-container">
            <div class="preview-area">
              <p class="preview-label">{{ t('friends.previewLabel') }}</p>
              <div class="link-row preview-row">
                <div class="link-row-inner">
                  <div class="link-smoke"></div>
                  <div class="link-avatar">
                    <img v-if="form.avatarUrl" :src="form.avatarUrl" :alt="form.name" />
                    <span v-else class="link-fallback">{{ form.name ? getInitial(form.name) : '?' }}</span>
                  </div>
                  <div class="link-primary">
                    <div class="link-heading">
                      <h3 class="link-name">{{ form.name || t('friends.previewName') }}</h3>
                      <span class="link-host">{{ form.url ? getHost(form.url) : t('friends.previewHost') }}</span>
                    </div>
                    <p class="link-desc">{{ form.description || t('friends.previewDescription') }}</p>
                  </div>
                  <span class="link-tail">{{ t('friends.outboundLabel') }}</span>
                </div>
              </div>
            </div>

            <form @submit.prevent="handleSubmit" class="submit-form">
              <input v-model="form.name" type="text" :placeholder="t('friends.formName')" required class="form-input" />
              <input v-model="form.url" type="url" :placeholder="t('friends.formUrl')" required class="form-input" />
              <input v-model="form.avatarUrl" type="url" :placeholder="t('friends.formAvatar')" class="form-input" />
              <textarea v-model="form.description" :placeholder="t('friends.formDescription')" rows="2" class="form-textarea"></textarea>
              <input v-model="form.contactEmail" type="email" :placeholder="t('friends.formEmail')" class="form-input" />
              <button type="submit" :disabled="submitting" class="submit-btn">
                {{ submitting ? t('friends.submitting') : t('friends.submitAction') }}
              </button>
              <p v-if="submitMessage" class="submit-message" :class="{ success: submitSuccess }">{{ submitMessage }}</p>
            </form>
          </div>
        </div>
      </section>

      <section
        v-if="article?.id && article.allowComment && siteSettings?.commentEnabled"
        class="friends-comments"
      >
        <CommentList
          content-type="article"
          :content-id="article.id"
          :require-approval="siteSettings?.commentRequireApproval"
        />
      </section>
    </div>
  </div>
</template>

<script setup lang="ts">
const { t } = useI18n()
const { getApprovedLinks, getArticle, getSiteSettings, submitLink } = useApi()

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

const { data: siteSettings } = await useAsyncData('friends-site-settings', getSiteSettings, {
  getCachedData: () => undefined
})

const form = ref({ name: '', url: '', avatarUrl: '', description: '', contactEmail: '' })
const showSubmission = ref(false)
const submitting = ref(false)
const submitMessage = ref('')
const submitSuccess = ref(false)

function getInitial(name: string) {
  return name.trim().charAt(0).toUpperCase()
}

function getHost(url: string) {
  try {
    return new URL(url).host.replace(/^www\./, '')
  } catch {
    return url
  }
}

function handleLinkMove(event: MouseEvent) {
  const currentTarget = event.currentTarget as HTMLElement | null
  if (!currentTarget) return
  const inner = currentTarget.querySelector<HTMLElement>('.link-row-inner')
  if (!inner) return
  const rect = currentTarget.getBoundingClientRect()
  const offsetX = ((event.clientX - rect.left) / rect.width - 0.5) * 12
  const offsetY = ((event.clientY - rect.top) / rect.height - 0.5) * 8
  inner.style.setProperty('--fx', `${offsetX.toFixed(2)}px`)
  inner.style.setProperty('--fy', `${offsetY.toFixed(2)}px`)
}

function resetLinkMove(event: MouseEvent) {
  const currentTarget = event.currentTarget as HTMLElement | null
  if (!currentTarget) return
  const inner = currentTarget.querySelector<HTMLElement>('.link-row-inner')
  if (!inner) return
  inner.style.setProperty('--fx', '0px')
  inner.style.setProperty('--fy', '0px')
}

async function handleSubmit() {
  if (!siteSettings.value?.linkSubmissionEnabled) {
    submitSuccess.value = false
    submitMessage.value = t('friends.submitDisabled')
    return
  }
  submitting.value = true
  submitMessage.value = ''
  try {
    await submitLink(form.value)
    submitSuccess.value = true
    submitMessage.value = t('friends.submitSuccess')
    form.value = { name: '', url: '', avatarUrl: '', description: '', contactEmail: '' }
  } catch (err: any) {
    submitSuccess.value = false
    submitMessage.value = err.message || t('friends.submitFailure')
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
  margin-bottom: 44px;
  padding: 32px;
  background: var(--surface);
  border: 1px solid var(--border);
  border-radius: var(--radius-lg);
  box-shadow: var(--shadow-sm);
  animation: fade-up 0.6s var(--ease-spring) both;
}

.friends-comments {
  margin-bottom: 54px;
}

.directory-shell {
  display: flex;
  flex-direction: column;
  gap: 0;
}

.section-kicker {
  margin: 0 0 8px;
  font-size: 11px;
  font-weight: 700;
  letter-spacing: 0.16em;
  text-transform: uppercase;
  color: var(--text-muted);
}

.section-title {
  margin: 0;
  font-size: clamp(1.1rem, 2vw, 1.45rem);
  font-weight: 700;
  letter-spacing: -0.01em;
}

.section-note,
.submit-note {
  max-width: 420px;
  margin: 0;
  font-size: 13px;
  line-height: 1.8;
  color: var(--text-subtle);
}

@keyframes fade-up {
  from { opacity: 0; transform: translateY(20px); }
  to { opacity: 1; transform: translateY(0); }
}

.links-directory {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.link-row {
  display: block;
  text-decoration: none;
  color: var(--text);
  opacity: 0;
  animation: fade-up 0.5s var(--ease-spring) forwards;
}

.link-row-inner {
  --fx: 0px;
  --fy: 0px;
  position: relative;
  display: grid;
  grid-template-columns: 56px minmax(0, 1fr) auto;
  align-items: center;
  gap: 16px;
  padding: 16px 18px;
  transform: translate3d(var(--fx), var(--fy), 0);
  transition:
    transform 280ms var(--ease-spring),
    color var(--dur-fast),
    border-color var(--dur-fast);
}

.link-smoke {
  position: absolute;
  inset: 0;
  border-radius: 20px;
  background:
    radial-gradient(circle at 20% 50%, color-mix(in srgb, var(--accent) 9%, transparent), transparent 42%),
    color-mix(in srgb, var(--text) 3%, transparent);
  opacity: 0;
  transition: opacity 220ms ease, transform 280ms var(--ease-spring);
  transform: translate3d(calc(var(--fx) * 0.45), calc(var(--fy) * 0.45), 0);
}

.link-row:hover .link-smoke {
  opacity: 1;
}

.link-row:hover .link-row-inner {
  color: var(--text);
}

.link-primary {
  min-width: 0;
}

.link-heading {
  display: flex;
  align-items: baseline;
  gap: 12px;
  min-width: 0;
}

.link-card,
.preview-card {
  display: block;
}

.link-card {
  opacity: 0;
}

.link-avatar {
  width: 56px;
  height: 56px;
  flex-shrink: 0;
  border-radius: 50%;
  overflow: hidden;
  background: var(--bg-secondary);
  border: 1px solid color-mix(in srgb, var(--border) 84%, transparent);
  display: grid;
  place-items: center;
  transition: transform 280ms var(--ease-spring), border-color var(--dur-fast);
}

.link-avatar img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.link-fallback {
  font-size: 20px;
  font-weight: 700;
  color: var(--text-muted);
}

.link-name {
  font-size: 15px;
  font-weight: 650;
  margin: 0;
  color: var(--text);
}

.link-host {
  min-width: 0;
  font-size: 12px;
  color: var(--text-muted);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.link-desc {
  font-size: 13px;
  color: var(--text-muted);
  margin: 4px 0 0;
  display: -webkit-box;
  -webkit-line-clamp: 1;
  -webkit-box-orient: vertical;
  overflow: hidden;
  line-height: 1.7;
}

.link-tail {
  padding-left: 10px;
  font-size: 12px;
  color: var(--text-subtle);
  white-space: nowrap;
}

/* Skeleton */
.skeleton-row {
  opacity: 1;
  animation: none;
}

.skeleton-row .link-row-inner {
  grid-template-columns: 56px minmax(0, 1fr) 72px;
}

.skeleton-avatar { width: 56px; height: 56px; border-radius: 50%; background: var(--border); }
.skeleton-line { background: var(--border); border-radius: 4px; }
.skeleton-copy {
  display: flex;
  flex-direction: column;
  gap: 10px;
}
.line-title { width: 34%; height: 14px; }
.line-desc { width: 78%; height: 10px; }
.line-meta { width: 56px; height: 12px; }
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
}

.submit-shell {
  padding: 28px 0 0;
  border-top: 1px solid color-mix(in srgb, var(--border) 78%, transparent);
}

.submit-head {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 20px;
  margin-bottom: 18px;
}

.submit-trigger {
  min-height: 42px;
  padding: 0 18px;
  border: 1px solid color-mix(in srgb, var(--border) 84%, transparent);
  border-radius: 999px;
  background: color-mix(in srgb, var(--accent) 10%, transparent);
  color: var(--text);
  font-size: 13px;
  font-weight: 600;
  cursor: pointer;
  transition:
    transform 260ms var(--ease-spring),
    background var(--dur-fast),
    border-color var(--dur-fast);
}

.submit-trigger:hover {
  transform: translateY(-1px);
  border-color: color-mix(in srgb, var(--accent) 28%, var(--border));
}

.submission-rules {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 10px;
  padding: 0;
  margin: 0 0 22px;
  list-style: none;
}

.submission-rules li {
  padding: 10px 12px;
  border-radius: 14px;
  background: color-mix(in srgb, var(--bg-secondary) 66%, transparent);
  font-size: 13px;
  line-height: 1.7;
  color: var(--text-subtle);
}

.submit-container {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 32px;
  padding-top: 10px;
  animation: fade-up 0.45s var(--ease-spring) both;
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

.preview-row {
  pointer-events: none;
  opacity: 1;
  animation: none;
}

.preview-row .link-row-inner {
  background: color-mix(in srgb, var(--bg-secondary) 62%, transparent);
  border-radius: 20px;
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

  .friends-intro {
    padding: 24px;
    margin-bottom: 34px;
  }

  .section-heading,
  .submit-head {
    flex-direction: column;
    align-items: flex-start;
  }

  .submission-rules {
    grid-template-columns: 1fr;
  }

  .submit-container {
    grid-template-columns: 1fr;
  }

  .link-row-inner {
    grid-template-columns: 48px minmax(0, 1fr);
    gap: 14px;
    padding: 14px 0;
  }

  .link-tail {
    grid-column: 2;
    padding-left: 0;
  }

  .link-avatar {
    width: 48px;
    height: 48px;
  }
}
</style>
