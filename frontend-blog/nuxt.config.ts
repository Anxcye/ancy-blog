// File: nuxt.config.ts
// Purpose: Nuxt 4 configuration — modules, i18n, color-mode, API proxy, and font loading.
// Module: frontend-blog, application configuration layer.
// Related: app.vue, layouts/default.vue, composables/useApi.ts.

export default defineNuxtConfig({
  compatibilityDate: '2025-07-15',
  devtools: { enabled: true },

  // ── Directory conventions (Nuxt 4) ──────────────────────────────
  future: { compatibilityVersion: 4 },

  // ── Modules ─────────────────────────────────────────────────────
  modules: [
    '@nuxt/ui',
    '@pinia/nuxt',
    '@vueuse/nuxt',
    '@nuxtjs/i18n',
    '@nuxtjs/google-fonts',
  ],

  // ── Color mode ──────────────────────────────────────────────────
  colorMode: {
    preference: 'system',   // respect OS preference on first visit
    fallback: 'light',
    classSuffix: '',        // adds 'dark' / 'light' class on <html>
    storageKey: 'blog-color-mode',
  },

  // ── i18n ────────────────────────────────────────────────────────
  i18n: {
    strategy: 'prefix_except_default',
    defaultLocale: 'zh',
    locales: [
      { code: 'zh', language: 'zh-CN', name: '中文', file: 'zh.json' },
      { code: 'en', language: 'en-US', name: 'English', file: 'en.json' },
    ],
    langDir: 'locales/',
    detectBrowserLanguage: {
      useCookie: true,
      cookieKey: 'blog-locale',
      redirectOn: 'root',
    },
  },

  // ── Google Fonts ─────────────────────────────────────────────────
  googleFonts: {
    families: {
      'Noto+Serif+SC': [400, 500, 700],
      'Inter': [400, 500, 600],
    },
    display: 'swap',
    preload: true,
  },

  // ── Runtime config ───────────────────────────────────────────────
  runtimeConfig: {
    public: {
      apiBase: process.env.NUXT_PUBLIC_API_BASE || 'http://localhost:8080/api/v1',
    },
  },

  // ── Dev server proxy (avoid CORS in dev) ─────────────────────────
  nitro: {
    devProxy: {
      '/api': {
        target: process.env.NUXT_PUBLIC_API_BASE?.replace('/api/v1', '') || 'http://localhost:8080',
        changeOrigin: true,
      },
    },
  },

  // ── CSS ──────────────────────────────────────────────────────────
  css: ['~/assets/css/main.css'],

  // ── App head ─────────────────────────────────────────────────────
  app: {
    head: {
      htmlAttrs: { lang: 'zh-CN' },
      meta: [
        { charset: 'utf-8' },
        { name: 'viewport', content: 'width=device-width, initial-scale=1' },
        { name: 'theme-color', content: '#ffffff', media: '(prefers-color-scheme: light)' },
        { name: 'theme-color', content: '#0f1117', media: '(prefers-color-scheme: dark)' },
      ],
      link: [
        { rel: 'icon', type: 'image/x-icon', href: '/favicon.ico' },
      ],
    },
    pageTransition: { name: 'page', mode: 'out-in' },
  },

  // ── Vite ─────────────────────────────────────────────────────────
  vite: {
    css: {
      preprocessorOptions: {
        scss: { api: 'modern-compiler' },
      },
    },
  },
})
