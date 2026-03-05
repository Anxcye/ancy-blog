<!-- File: app/app.vue
     Purpose: Root component — injects random accent color on mount, and applies color-mode class.
     Module: frontend-blog/app, root layer.
     Related: nuxt.config.ts (colorMode), assets/css/main.css (accent classes). -->
<template>
  <NuxtLayout>
    <NuxtPage />
  </NuxtLayout>
</template>

<script setup lang="ts">
// ── Accent color palette ──────────────────────────────────────────
const ACCENT_PRESETS = [
  { name: 'teal',   accent: '#2AA889', soft: 'rgba(42,168,137,0.10)',  text: '#1e9170' },
  { name: 'rose',   accent: '#e8738a', soft: 'rgba(232,115,138,0.10)', text: '#d05c75' },
  { name: 'violet', accent: '#8b5cf6', soft: 'rgba(139,92,246,0.10)',  text: '#7c3aed' },
  { name: 'amber',  accent: '#d97706', soft: 'rgba(217,119,6,0.10)',   text: '#b45309' },
  { name: 'sky',    accent: '#0ea5e9', soft: 'rgba(14,165,233,0.10)',  text: '#0284c7' },
]

// Pick a random accent on each page load and inject CSS variables
onMounted(() => {
  const preset = ACCENT_PRESETS[Math.floor(Math.random() * ACCENT_PRESETS.length)]
  const root = document.documentElement
  root.style.setProperty('--accent', preset.accent)
  root.style.setProperty('--accent-soft', preset.soft)
  root.style.setProperty('--accent-text', preset.text)
})

// ── SEO defaults ─────────────────────────────────────────────────
useHead({
  titleTemplate: (title) => title ? `${title} · Ancy Blog` : 'Ancy Blog',
})
</script>
