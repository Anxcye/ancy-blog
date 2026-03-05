<!-- File: app/layouts/default.vue
     Purpose: Site shell — fixed header with nav/theme toggle + footer.
     Module: frontend-blog/layouts, shell layer.
     Related: app.vue, composables/useApi.ts, pages/. -->
<template>
  <div class="site-wrap">
    <!-- ── Header ── -->
    <header class="site-header" :class="{ 'site-header--scrolled': scrolled }">
      <div class="header-inner container container--wide">

        <!-- Left: Avatar / brand -->
        <NuxtLink :to="localePath('/')" class="header-brand" :aria-label="t('nav.home')">
          <div class="header-avatar">
            <img
              v-if="siteSettings?.avatarUrl"
              :src="siteSettings.avatarUrl"
              :alt="siteSettings.siteName"
              width="28" height="28"
            />
            <span v-else class="header-avatar-fallback">A</span>
          </div>
        </NuxtLink>

        <!-- Center: Nav -->
        <nav class="header-nav" aria-label="主导航">
          <NuxtLink
            v-for="(item, i) in navItems"
            :key="item.key"
            :to="localePath(item.to)"
            class="nav-link"
            :style="{ '--nav-i': i }"
          >{{ item.label }}</NuxtLink>
        </nav>

        <!-- Right: Theme + Lang -->
        <div class="header-actions">
          <!-- Language switch -->
          <button class="icon-btn lang-btn" @click="toggleLocale" :title="locale === 'zh' ? 'English' : '中文'">
            <span class="lang-label">{{ locale === 'zh' ? 'EN' : 'ZH' }}</span>
          </button>

          <!-- Theme toggle -->
          <button
            class="icon-btn theme-btn"
            @click="toggleColorMode"
            :aria-label="isDark ? t('meta.lightMode') : t('meta.darkMode')"
          >
            <span class="theme-icon">
              <svg v-if="isDark" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                <circle cx="12" cy="12" r="5"/>
                <path d="M12 1v2M12 21v2M4.22 4.22l1.42 1.42M18.36 18.36l1.42 1.42M1 12h2M21 12h2M4.22 19.78l1.42-1.42M18.36 5.64l1.42-1.42"/>
              </svg>
              <svg v-else viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                <path d="M21 12.79A9 9 0 1 1 11.21 3 7 7 0 0 0 21 12.79z"/>
              </svg>
            </span>
          </button>

          <!-- Mobile menu -->
          <button class="icon-btn menu-btn" @click="mobileOpen = !mobileOpen" aria-label="菜单">
            <span class="menu-icon" :class="{ open: mobileOpen }">
              <span /><span /><span />
            </span>
          </button>
        </div>
      </div>

      <!-- Mobile drawer -->
      <Transition name="mobile-nav">
        <div v-if="mobileOpen" class="mobile-nav" @click="mobileOpen = false">
          <NuxtLink
            v-for="(item, i) in navItems"
            :key="item.key"
            :to="localePath(item.to)"
            class="mobile-nav-link"
            :style="{ animationDelay: `${i * 50}ms` }"
          >{{ item.label }}</NuxtLink>
        </div>
      </Transition>
    </header>

    <!-- ── Page content ── -->
    <main class="site-main">
      <slot />
    </main>

    <!-- ── Footer ── -->
    <footer class="site-footer">
      <div class="container">
        <p class="footer-copy">
          © {{ new Date().getFullYear() }}
          <span class="footer-accent">{{ siteSettings?.siteName || 'Ancy Blog' }}</span>
          · Built with Nuxt
        </p>
      </div>
    </footer>
  </div>
</template>

<script setup lang="ts">
import { useSiteStore } from '~/stores/site'

const { t, locale } = useI18n()
const localePath = useLocalePath()
const colorMode = useColorMode()
const { getSiteSettings } = useApi()
const siteStore = useSiteStore()

// ── Site settings (avatar / site name for header) ─────────────────
const { data: siteSettings } = await useAsyncData('site-settings', getSiteSettings, {
  server: true,
  lazy: false,
})

// ── Nav items ────────────────────────────────────────────────────
const defaultNavItems = computed(() => [
  { key: 'home',     to: '/',         label: t('nav.home') },
  { key: 'articles', to: '/articles', label: t('nav.articles') },
  { key: 'moments',  to: '/moments',  label: t('nav.moments') },
  { key: 'timeline', to: '/timeline', label: t('nav.timeline') },
])
const navItems = computed(() => {
  if (siteStore.navigation.length) {
    return siteStore.navigation.map(n => ({ key: n.id, to: n.targetValue || '/', label: n.name }))
  }
  return defaultNavItems.value
})

// ── Theme ──────────────────────────────────────────────────────
const isDark = computed(() => colorMode.value === 'dark')
function toggleColorMode() {
  colorMode.preference = isDark.value ? 'light' : 'dark'
}

// ── Language ───────────────────────────────────────────────────
const { setLocale } = useI18n()
function toggleLocale() {
  setLocale(locale.value === 'zh' ? 'en' : 'zh')
}

// ── Scroll detection ───────────────────────────────────────────
const scrolled = ref(false)
onMounted(() => {
  const onScroll = () => { scrolled.value = window.scrollY > 20 }
  window.addEventListener('scroll', onScroll, { passive: true })
  onUnmounted(() => window.removeEventListener('scroll', onScroll))
})

// ── Mobile nav ─────────────────────────────────────────────────
const mobileOpen = ref(false)
const route = useRoute()
watch(() => route.path, () => { mobileOpen.value = false })
</script>

<style scoped>
/* ── Site wrap ── */
.site-wrap {
  display: flex;
  flex-direction: column;
  min-height: 100vh;
}

.site-main {
  flex: 1;
}

/* ── Header ── */
.site-header {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  z-index: 100;
  height: var(--header-h);
  background: transparent;
  transition: background var(--dur-base) var(--ease-smooth),
              box-shadow var(--dur-base) var(--ease-smooth),
              backdrop-filter var(--dur-base) var(--ease-smooth);
}

.site-header--scrolled {
  background: rgba(var(--bg-rgb, 249,249,247), 0.85);
  backdrop-filter: blur(16px);
  -webkit-backdrop-filter: blur(16px);
  box-shadow: 0 1px 0 var(--border);
}

.dark .site-header--scrolled {
  background: rgba(15,17,23, 0.85);
}

.header-inner {
  height: 100%;
  display: flex;
  align-items: center;
  gap: 0;
}

/* ── Brand / Avatar ── */
.header-brand {
  display: flex;
  align-items: center;
  flex-shrink: 0;
  transition: opacity var(--dur-fast);
}
.header-brand:hover { opacity: 0.75; }

.header-avatar {
  width: 30px;
  height: 30px;
  border-radius: 50%;
  overflow: hidden;
  background: var(--accent-soft);
  display: grid;
  place-items: center;
  border: 1.5px solid var(--border);
}

.header-avatar img { width: 100%; height: 100%; object-fit: cover; }

.header-avatar-fallback {
  font-size: 13px;
  font-weight: 700;
  color: var(--accent-text);
  text-transform: uppercase;
}

/* ── Nav ── */
.header-nav {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 2px;
}

@keyframes nav-spring {
  0%   { opacity: 0; transform: translateY(-10px) scale(0.9); }
  60%  { opacity: 1; transform: translateY(2px) scale(1.02); }
  100% { opacity: 1; transform: translateY(0) scale(1); }
}

.nav-link {
  padding: 6px 12px;
  border-radius: var(--radius-sm);
  font-size: 14px;
  font-weight: 500;
  color: var(--text-muted);
  transition: color var(--dur-fast), background var(--dur-fast);
  position: relative;
  opacity: 0;
  animation: nav-spring 0.5s var(--ease-spring) forwards;
  animation-delay: calc(var(--nav-i, 0) * 70ms + 100ms);
}

.nav-link:hover {
  color: var(--text);
  background: var(--surface-hover);
}

.nav-link.router-link-active {
  color: var(--text);
  font-weight: 600;
}

.nav-link.router-link-active::after {
  content: '';
  position: absolute;
  bottom: 2px;
  left: 50%;
  transform: translateX(-50%);
  width: 16px;
  height: 2px;
  background: var(--accent);
  border-radius: 2px;
}

/* ── Actions ── */
.header-actions {
  display: flex;
  align-items: center;
  gap: 4px;
  flex-shrink: 0;
}

.icon-btn {
  width: 36px;
  height: 36px;
  border-radius: var(--radius-md);
  display: grid;
  place-items: center;
  color: var(--text-muted);
  transition: color var(--dur-fast), background var(--dur-fast);
  -webkit-tap-highlight-color: transparent;
}

.icon-btn:hover {
  color: var(--text);
  background: var(--surface-hover);
}

.theme-icon svg { width: 16px; height: 16px; }
.lang-label { font-size: 11px; font-weight: 700; letter-spacing: 0.04em; }

/* ── Mobile hamburger ── */
.menu-btn { display: none; }

.menu-icon {
  display: flex;
  flex-direction: column;
  gap: 4px;
  width: 18px;
}

.menu-icon span {
  display: block;
  height: 1.5px;
  background: currentColor;
  border-radius: 2px;
  transition: transform var(--dur-base) var(--ease-smooth), opacity var(--dur-fast);
  transform-origin: center;
}

.menu-icon.open span:nth-child(1) { transform: translateY(5.5px) rotate(45deg); }
.menu-icon.open span:nth-child(2) { opacity: 0; transform: scaleX(0); }
.menu-icon.open span:nth-child(3) { transform: translateY(-5.5px) rotate(-45deg); }

/* ── Mobile nav drawer ── */
.mobile-nav {
  position: absolute;
  top: var(--header-h);
  left: 0;
  right: 0;
  background: var(--surface);
  border-bottom: 1px solid var(--border);
  padding: 8px 20px 16px;
  display: flex;
  flex-direction: column;
  gap: 2px;
  box-shadow: var(--shadow-md);
}

.mobile-nav-link {
  padding: 12px 8px;
  font-size: 15px;
  font-weight: 500;
  color: var(--text-muted);
  border-bottom: 1px solid var(--border);
  transition: color var(--dur-fast);
}

.mobile-nav-link:last-child { border-bottom: none; }
.mobile-nav-link:hover, .mobile-nav-link.router-link-active { color: var(--text); }

/* Mobile nav transition */
.mobile-nav-enter-active, .mobile-nav-leave-active {
  transition: opacity var(--dur-base), transform var(--dur-base) var(--ease-smooth);
}
.mobile-nav-enter-from, .mobile-nav-leave-to {
  opacity: 0;
  transform: translateY(-8px);
}

/* ── Footer ── */
.site-footer {
  padding: 32px 0;
  border-top: 1px solid var(--border);
  margin-top: 80px;
}

.footer-copy {
  font-size: 13px;
  color: var(--text-subtle);
  text-align: center;
}

.footer-accent {
  color: var(--accent-text);
  font-weight: 500;
}

/* ── Responsive ── */
@media (max-width: 640px) {
  .header-nav { display: none; }
  .menu-btn { display: grid; }
  .lang-btn { display: none; }
}
</style>
