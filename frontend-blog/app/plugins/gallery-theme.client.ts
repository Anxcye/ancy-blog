// File: gallery-theme.client.ts
// Purpose: Force dark mode on gallery routes and restore the previous theme after leaving gallery pages.
// Module: frontend-blog/plugins, client-side route/theme behavior.
// Related: nuxt.config.ts colorMode settings, layouts/default.vue theme toggle, pages/gallery/*.

export default defineNuxtPlugin(() => {
  const route = useRoute()
  const colorMode = useColorMode()
  const previousGalleryTheme = useState<string | null>('gallery-previous-theme', () => null)

  const isGalleryPath = (path: string) => /^\/(?:en\/)?gallery(?:\/|$)/.test(path)

  watch(
    () => route.path,
    (path) => {
      if (isGalleryPath(path)) {
        if (previousGalleryTheme.value === null) {
          previousGalleryTheme.value = colorMode.preference
        }
        if (colorMode.preference !== 'dark') {
          colorMode.preference = 'dark'
        }
        return
      }

      if (previousGalleryTheme.value !== null) {
        colorMode.preference = previousGalleryTheme.value
        previousGalleryTheme.value = null
      }
    },
    { immediate: true }
  )
})
