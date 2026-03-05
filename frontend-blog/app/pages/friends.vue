<!-- File: app/pages/friends.vue
     Purpose: Friends/Links page showing approved links and an optional pinned article introduction.
     Module: frontend-blog/pages
-->
<template>
  <div class="friends-page">
    <div class="container">
      <div class="page-hero">
        <h1 class="page-title word-bounce">友人帐</h1>
        <p class="page-subtitle word-bounce" style="animation-delay: 100ms">散落平行时空的节点，都在这里建立连接。</p>
      </div>

      <!-- Optional Article Intro -->
      <div v-if="article" class="friends-intro prose" :class="{ dark: isDark }">
        <div class="prose-content" v-html="renderedContent"></div>
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
    </div>
  </div>
</template>

<script setup lang="ts">
import { generateHTML } from '@tiptap/html'
import Document from '@tiptap/extension-document'
import Paragraph from '@tiptap/extension-paragraph'
import Text from '@tiptap/extension-text'
import Heading from '@tiptap/extension-heading'
import Blockquote from '@tiptap/extension-blockquote'
import PrecodeBlock from '~/components/tiptap/PrecodeBlock'
import CustomLink from '~/components/tiptap/CustomLink'
import CustomImage from '~/components/tiptap/CustomImage'

const { getApprovedLinks, getArticle } = useApi()
const colorMode = useColorMode()
const isDark = computed(() => colorMode.value === 'dark')

// Fetch approved links
const { data: links, pending } = await useAsyncData('friends-links', getApprovedLinks, {
  getCachedData: () => undefined
})

// Optionally load 'friends' article for intro text if it exists
const { data: article } = await useAsyncData('friends-intro', async () => {
  try {
    return await getArticle('friends')
  } catch (err: any) {
    if (err.statusCode === 404) return null
    throw err
  }
}, { getCachedData: () => undefined })

const renderedContent = computed(() => {
  if (!article.value?.content) return ''
  try {
    const doc = JSON.parse(article.value.content)
    return generateHTML(doc, [
      Document, Paragraph, Text,
      Heading.configure({ levels: [1, 2, 3, 4, 5, 6] }),
      Blockquote, PrecodeBlock, CustomLink, CustomImage
    ])
  } catch {
    return `<p>${article.value.content}</p>`
  }
})

useSeoMeta({ title: '友人帐 - 友情链接' })
</script>

<style scoped>
.friends-page {
  padding-top: calc(var(--header-h) + 64px);
  padding-bottom: 80px;
}

.page-hero {
  text-align: center;
  margin-bottom: 56px;
}

.page-title {
  font-size: clamp(2rem, 4vw, 3rem);
  font-weight: 800;
  margin-bottom: 12px;
}

.page-subtitle {
  font-size: 1.1rem;
  color: var(--text-muted);
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
  align-items: center;
  gap: 16px;
  padding: 20px;
  background: var(--surface);
  border: 1px solid var(--border);
  border-radius: var(--radius-lg);
  text-decoration: none;
  color: var(--text);
  transition: all var(--dur-base) var(--ease-out);
  opacity: 0;
  animation: fade-up 0.5s var(--ease-spring) forwards;
}

.link-card:hover {
  transform: translateY(-4px);
  border-color: var(--accent);
  box-shadow: var(--shadow-md), 0 0 0 1px var(--accent-soft);
}

.link-avatar {
  width: 56px;
  height: 56px;
  flex-shrink: 0;
  border-radius: 50%;
  overflow: hidden;
  background: var(--bg-secondary);
  border: 1.5px solid var(--border);
  display: grid;
  place-items: center;
  transition: transform var(--dur-base);
}

.link-card:hover .link-avatar {
  transform: scale(1.05) rotate(5deg);
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
  font-size: 16px;
  font-weight: 600;
  margin: 0 0 4px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
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
