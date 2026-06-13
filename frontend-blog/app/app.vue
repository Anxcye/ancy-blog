<!-- File: app/app.vue
     Purpose: Root component — provides app shell, page layout, and global SEO defaults.
     Module: frontend-blog/app, root layer.
     Related: nuxt.config.ts, layouts/default.vue, assets/css/main.css. -->
<template>
  <UApp>
    <NuxtLayout>
      <NuxtPage />
    </NuxtLayout>
  </UApp>
</template>

<script setup lang="ts">
const { getSiteSettings } = useApi()
const accentClasses = ['accent-teal', 'accent-rose', 'accent-violet', 'accent-amber', 'accent-sky']
const accentClass = useState('site-accent-class', () => accentClasses[Math.floor(Math.random() * accentClasses.length)])

// ── SEO defaults ─────────────────────────────────────────────────
const { data: siteSettings } = await useAsyncData('global-site-settings', getSiteSettings)

useHead({
  htmlAttrs: {
    class: accentClass.value,
  },
  titleTemplate: (title) => {
    const siteName = siteSettings.value?.siteName || 'Ancy Blog'
    return title ? `${title} · ${siteName}` : siteName
  },
  link: [
    {
      rel: 'icon',
      type: 'image/x-icon',
      href: siteSettings.value?.faviconUrl || '/favicon.ico',
    },
  ],
})
</script>
