<!-- File: app/layouts/default.vue
     Purpose: Site shell — fixed header with nav/theme toggle + footer.
     Module: frontend-blog/layouts, shell layer.
     Related: app.vue, composables/useApi.ts, pages/. -->
<template>
  <div class="site-wrap">
    <!-- ── Header ── -->
    <header class="site-header" :class="{ 'site-header--scrolled': scrolled }">
      <div class="header-inner container container--wide">

        <!-- Mobile menu button (first in DOM for left position) -->
        <button class="icon-btn menu-btn" @click="mobileOpen = !mobileOpen" aria-label="菜单">
          <span class="menu-icon" :class="{ open: mobileOpen }">
            <span /><span /><span />
          </span>
        </button>

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
          <div
            v-for="(item, i) in navItems"
            :key="item.key"
            class="nav-item-wrap"
            :style="{ '--nav-i': i }"
          >
            <a v-if="item.isExternal" :href="item.to" target="_blank" class="nav-link">{{ item.label }}</a>
            <NuxtLink v-else :to="localePath(item.to)" class="nav-link">{{ item.label }}</NuxtLink>

            <!-- Dropdown -->
            <div v-if="item.children && item.children.length" class="nav-dropdown">
              <template v-for="child in item.children" :key="child.key">
                <a v-if="child.isExternal" :href="child.to" target="_blank" class="dropdown-link">{{ child.label }}</a>
                <NuxtLink v-else :to="localePath(child.to)" class="dropdown-link">{{ child.label }}</NuxtLink>
              </template>
            </div>
          </div>
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
        </div>
      </div>

      <!-- Mobile drawer -->
      <Transition name="mobile-nav">
        <div v-if="mobileOpen" class="mobile-nav-overlay" @click="mobileOpen = false">
          <div class="mobile-nav" @click.stop>
            <template v-for="item in navItems" :key="item.key">
              <!-- Primary nav item -->
              <div class="mobile-nav-section">
                <a v-if="item.isExternal" :href="item.to" target="_blank" class="mobile-nav-primary" @click="mobileOpen = false">
                  {{ item.label }}
                </a>
                <NuxtLink v-else :to="localePath(item.to)" class="mobile-nav-primary" @click="mobileOpen = false">
                  {{ item.label }}
                </NuxtLink>

                <!-- Secondary nav items (children) -->
                <div v-if="item.children?.length" class="mobile-nav-secondary">
                  <template v-for="child in item.children" :key="child.key">
                    <a v-if="child.isExternal" :href="child.to" target="_blank" class="mobile-nav-link" @click="mobileOpen = false">
                      {{ child.label }}
                    </a>
                    <NuxtLink v-else :to="localePath(child.to)" class="mobile-nav-link" @click="mobileOpen = false">
                      {{ child.label }}
                    </NuxtLink>
                  </template>
                </div>
              </div>
            </template>
          </div>
        </div>
      </Transition>
    </header>

    <!-- ── Page content ── -->
    <main class="site-main">
      <slot />
    </main>

    <!-- ── Footer ── -->
    <footer class="site-footer">
      <div class="container footer-container">
        <!-- Left: Footer rows -->
        <div class="footer-left">
          <div v-if="footerRows.length" class="footer-rows">
            <div v-for="row in footerRows" :key="row.rowNum" class="footer-row">
              <template v-for="item in row.items" :key="item.id">
                <a v-if="item.linkType === 'external'" :href="item.externalUrl" target="_blank" class="footer-link">
                  {{ item.label }}
                </a>
                <NuxtLink v-else-if="item.linkType === 'internal'" :to="localePath(`/articles/${item.internalArticleSlug}`)" class="footer-link">
                  {{ item.label }}
                </NuxtLink>
                <span v-else class="footer-text">{{ item.label }}</span>
              </template>
            </div>
          </div>

          <!-- Copyright -->
          <p class="footer-copy">
            © {{ new Date().getFullYear() }}
            <a href="https://github.com/anxcye/ancy-blog" target="_blank" class="footer-accent">{{ siteSettings?.siteName || 'Ancy Blog' }}</a>
          </p>
        </div>

        <!-- Right: Social links -->
        <div v-if="siteStore.socialLinks?.length" class="footer-right">
          <a
            v-for="link in siteStore.socialLinks"
            :key="link.id"
            :href="link.url"
            target="_blank"
            class="footer-social"
            :title="link.title"
          >
            {{ link.title }}
          </a>
        </div>
      </div>
    </footer>
  </div>
</template>

<script setup lang="ts">
import { useSiteStore } from '~/stores/site'

const { t, locale } = useI18n()
const localePath = useLocalePath()
const colorMode = useColorMode()
const { getSiteSettings, getCategories } = useApi()
const siteStore = useSiteStore()

// ── Site settings (avatar / site name for header) ─────────────────
const { data: siteSettings } = await useAsyncData('site-settings', getSiteSettings, {
  server: true,
  lazy: false,
})

// ── Categories for dynamic dropdown ────────────────────────────────
const { data: categories } = await useAsyncData('categories', getCategories, {
  server: true,
  lazy: false,
})

// ── Load site navigation ───────────────────────────────────────────
await siteStore.fetchAll()

// ── Nav items ────────────────────────────────────────────────────
function resolveTarget(n: any) {
  switch (n.targetType) {
    case 'external': return n.targetValue || '#';
    case 'article': return `/articles/${n.targetValue}`;
    case 'category': return `/articles/category/${n.targetValue}`;
    case 'route': default: return n.targetValue || '/';
  }
}

function mapNav(n: any): any {
  const label = getNavLabel(n.key, n.name)

  let children = n.children?.length ? n.children.map(mapNav) : undefined

  // Auto-inject categories: if key is 'articles' and no children, show category dropdown
  if (n.key === 'articles' && !children && categories.value?.length) {
    children = categories.value.map((cat: any) => ({
      key: `cat-${cat.slug}`,
      label: cat.name,
      to: `/articles?category=${cat.slug}`,
      isExternal: false,
    }))
  }

  // Or if targetType is 'category' with no targetValue and no children, inject all categories
  if (n.targetType === 'category' && !n.targetValue && !children && categories.value?.length) {
    children = categories.value.map((cat: any) => ({
      key: `cat-${cat.slug}`,
      label: cat.name,
      to: `/articles?category=${cat.slug}`,
      isExternal: false,
    }))
  }

  return {
    key: n.id || n.key,
    label,
    to: resolveTarget(n),
    isExternal: n.targetType === 'external',
    children
  }
}

const defaultNavItems = computed(() => [
  { key: 'home',     to: '/',         label: getNavLabel('home', t('nav.home')), isExternal: false },
  { key: 'articles', to: '/articles', label: getNavLabel('articles', t('nav.articles')), isExternal: false },
  { key: 'moments',  to: '/moments',  label: getNavLabel('moments', t('nav.moments')), isExternal: false },
  { key: 'timeline', to: '/timeline', label: getNavLabel('timeline', t('nav.timeline')), isExternal: false },
  { key: 'links',    to: '/friends',  label: getNavLabel('links', t('nav.links')), isExternal: false },
])

const navItems = computed(() => {
  if (siteStore.navigation?.length) {
    return siteStore.navigation.map(mapNav)
  }
  return defaultNavItems.value
})

const footerRows = computed(() => {
  const footer = siteStore.footer
  if (!footer || !Object.keys(footer).length) return []

  return [1, 2, 3]
    .map(rowNum => ({
      rowNum,
      items: footer[rowNum] || []
    }))
    .filter(row => row.items.length > 0)
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

function getNavLabel(key: string, fallback: string) {
  const displayKey = `navDisplay.${key}`
  if (t(displayKey) !== displayKey) {
    return t(displayKey)
  }

  const navKey = `nav.${key}`
  if (t(navKey) !== navKey) {
    return t(navKey)
  }

  return fallback
}
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

.nav-item-wrap {
  position: relative;
  display: inline-block;
}

.nav-link {
  display: block;
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

.nav-link:hover, .nav-link.router-link-active {
  color: var(--text);
  background: var(--surface-hover);
}

/* Nav dropdown styling */
.nav-dropdown {
  position: absolute;
  top: 100%;
  left: 50%;
  transform: translateX(-50%) translateY(4px);
  min-width: 120px;
  background: rgba(var(--bg-rgb, 249,249,247), 0.75);
  backdrop-filter: blur(20px);
  -webkit-backdrop-filter: blur(20px);
  border: 1px solid rgba(var(--border-rgb, 0,0,0), 0.08);
  border-radius: 8px;
  box-shadow: 0 4px 16px rgba(0,0,0,0.08), 0 1px 4px rgba(0,0,0,0.04);
  padding: 0;
  display: flex;
  flex-direction: column;
  gap: 0;
  opacity: 0;
  visibility: hidden;
  pointer-events: none;
  transition: opacity 0.2s var(--ease-out), transform 0.2s var(--ease-out), visibility 0.2s;
  z-index: 10;
  overflow: hidden;
}

.dark .nav-dropdown {
  background: rgba(15,17,23, 0.85);
  border-color: rgba(255,255,255,0.08);
  box-shadow: 0 4px 20px rgba(0,0,0,0.4), 0 1px 4px rgba(0,0,0,0.2);
}

.nav-item-wrap:hover .nav-dropdown {
  opacity: 1;
  visibility: visible;
  pointer-events: auto;
  transform: translateX(-50%) translateY(0);
}

.dropdown-link {
  display: block;
  padding: 10px 16px;
  font-size: 14px;
  font-weight: 500;
  color: var(--text-muted);
  transition: background 0.15s var(--ease-out), color 0.15s var(--ease-out);
  white-space: nowrap;
  text-align: center;
  text-decoration: none;
}

.dropdown-link:hover {
  background: var(--accent-soft);
  color: var(--accent-text);
}

.dropdown-link.router-link-active:hover {
  background: var(--accent-soft);
  color: var(--accent-text);
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
.mobile-nav-overlay {
  position: fixed;
  top: var(--header-h);
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0,0,0,0.3);
  backdrop-filter: blur(4px);
  z-index: 99;
}

.mobile-nav {
  background: var(--surface);
  border-bottom: 1px solid var(--border);
  padding: 16px 20px;
  box-shadow: var(--shadow-md);
  max-height: calc(100vh - var(--header-h));
  overflow-y: auto;
}

.mobile-nav-section {
  margin-bottom: 8px;
  opacity: 0;
  animation: mobile-nav-spring 0.5s var(--ease-spring) forwards;
}

.mobile-nav-section:nth-child(1) { animation-delay: 0.05s; }
.mobile-nav-section:nth-child(2) { animation-delay: 0.1s; }
.mobile-nav-section:nth-child(3) { animation-delay: 0.15s; }
.mobile-nav-section:nth-child(4) { animation-delay: 0.2s; }
.mobile-nav-section:nth-child(5) { animation-delay: 0.25s; }

@keyframes mobile-nav-spring {
  0% {
    opacity: 0;
    transform: translateY(-20px) scale(0.95);
  }
  60% {
    opacity: 1;
    transform: translateY(4px) scale(1.02);
  }
  100% {
    opacity: 1;
    transform: translateY(0) scale(1);
  }
}

.mobile-nav-primary {
  display: block;
  padding: 12px 0;
  font-size: 16px;
  font-weight: 600;
  color: var(--text);
  text-decoration: none;
  border-bottom: 1px solid var(--border);
}

.mobile-nav-secondary {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  margin-top: 12px;
  padding-bottom: 12px;
}

.mobile-nav-link {
  display: inline-block;
  padding: 8px 16px;
  font-size: 13px;
  font-weight: 500;
  color: var(--text-muted);
  background: var(--bg-secondary);
  border: 1px solid var(--border);
  border-radius: 20px;
  text-decoration: none;
  transition: all var(--dur-fast);
}

.mobile-nav-link:hover, .mobile-nav-link.router-link-active {
  color: var(--accent-text);
  background: var(--accent-soft);
  border-color: var(--accent);
}

/* Mobile nav transition */
.mobile-nav-enter-active, .mobile-nav-leave-active {
  transition: opacity var(--dur-base);
}
.mobile-nav-enter-from, .mobile-nav-leave-to {
  opacity: 0;
}

/* ── Footer ── */
.site-footer {
  padding: 48px 0 32px;
  margin-top: 80px;
  background: var(--accent-soft);
  border-top: 1px solid var(--accent);
}

.footer-container {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  gap: 48px;
}

.footer-left {
  flex: 1;
}

.footer-rows {
  display: flex;
  flex-direction: column;
  gap: 12px;
  margin-bottom: 24px;
}

.footer-row {
  display: flex;
  flex-wrap: wrap;
  gap: 20px;
  font-size: 14px;
}

.footer-row:first-child {
  font-size: 16px;
  font-weight: 500;
}

.footer-link {
  color: var(--text-muted);
  text-decoration: none;
  transition: color var(--dur-fast);
}

.footer-link:hover {
  color: var(--accent-text);
}

.footer-text {
  color: var(--text-subtle);
}

.footer-copy {
  font-size: 13px;
  color: var(--text-subtle);
}

.footer-accent {
  color: var(--accent-text);
  font-weight: 600;
  text-decoration: none;
  transition: opacity var(--dur-fast);
}

.footer-accent:hover {
  opacity: 0.8;
}

.footer-right {
  display: flex;
  flex-direction: column;
  gap: 8px;
  align-items: flex-end;
}

.footer-social {
  font-size: 13px;
  color: var(--text-muted);
  text-decoration: none;
  transition: color var(--dur-fast);
}

.footer-social:hover {
  color: var(--accent-text);
}

@media (max-width: 640px) {
  .footer-container {
    flex-direction: column;
    gap: 32px;
  }

  .footer-right {
    align-items: flex-start;
  }
}

/* ── Responsive ── */
@media (max-width: 640px) {
  .header-inner {
    justify-content: space-between;
  }

  /* Left: menu button */
  .menu-btn {
    display: grid;
  }

  /* Center: avatar */
  .header-brand {
    position: absolute;
    left: 50%;
    transform: translateX(-50%);
  }

  /* Right: lang + theme buttons */
  .header-actions {
    justify-content: flex-end;
  }

  .header-nav { display: none; }
}
</style>
