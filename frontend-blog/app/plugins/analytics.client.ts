// File: analytics.client.ts
// Purpose: Report public page analytics events from the browser after route entry.
// Module: frontend-blog/plugins/analytics, client instrumentation layer.
// Related: public analytics ingest API and admin analytics reporting views.

interface AnalyticsPayload {
  eventId: string
  eventType: 'page_view' | 'page_ping'
  occurredAt: string
  visitorId: string
  sessionId: string
  path: string
  routeName: string
  pageTitle: string
  referrer: string
  contentType: string
  contentId: string
  contentSlug: string
  locale: string
  screenWidth: number
  screenHeight: number
  viewportWidth: number
  viewportHeight: number
  timezone: string
}

const VISITOR_KEY = 'ancy_blog_visitor_id'
const SESSION_KEY = 'ancy_blog_session_id'
const SESSION_AT_KEY = 'ancy_blog_session_last_seen_at'
const SESSION_TTL_MS = 30 * 60 * 1000
const PING_INTERVAL_MS = 15_000

function randomId(prefix: string): string {
  if (typeof crypto !== 'undefined' && typeof crypto.randomUUID === 'function') {
    return `${prefix}_${crypto.randomUUID()}`
  }
  return `${prefix}_${Math.random().toString(36).slice(2)}${Date.now().toString(36)}`
}

function getVisitorId(): string {
  const current = window.localStorage.getItem(VISITOR_KEY)
  if (current) return current
  const next = randomId('visitor')
  window.localStorage.setItem(VISITOR_KEY, next)
  return next
}

function getSessionId(now = Date.now()): string {
  const current = window.sessionStorage.getItem(SESSION_KEY)
  const lastSeenRaw = window.sessionStorage.getItem(SESSION_AT_KEY)
  const lastSeen = lastSeenRaw ? Number(lastSeenRaw) : 0
  if (current && lastSeen > 0 && now-lastSeen < SESSION_TTL_MS) {
    window.sessionStorage.setItem(SESSION_AT_KEY, String(now))
    return current
  }
  const next = randomId('session')
  window.sessionStorage.setItem(SESSION_KEY, next)
  window.sessionStorage.setItem(SESSION_AT_KEY, String(now))
  return next
}

function normalizedPath(pathname: string): string {
  if (!pathname) return '/'
  return pathname.startsWith('/') ? pathname : `/${pathname}`
}

function inferContent(path: string): Pick<AnalyticsPayload, 'contentType' | 'contentId' | 'contentSlug'> {
  const normalized = path.replace(/^\/en(?=\/|$)/, '') || '/'
  const parts = normalized.split('/').filter(Boolean)
  if (parts[0] === 'articles' && parts[1]) {
    return { contentType: 'article', contentId: parts[1], contentSlug: parts[1] }
  }
  if (parts[0] === 'moments' && parts[1]) {
    return { contentType: 'moment', contentId: parts[1], contentSlug: '' }
  }
  return { contentType: 'site', contentId: normalized, contentSlug: '' }
}

function sendAnalytics(events: AnalyticsPayload[]): void {
  const body = JSON.stringify({ events })
  const url = '/api/v1/public/analytics/events'
  if (typeof navigator.sendBeacon === 'function') {
    const blob = new Blob([body], { type: 'application/json' })
    navigator.sendBeacon(url, blob)
    return
  }
  void fetch(url, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body,
    keepalive: true,
  }).catch(() => undefined)
}

export default defineNuxtPlugin((nuxtApp) => {
  let previousReferrer = document.referrer || ''
  let currentPath = ''
  let pingTimer: ReturnType<typeof window.setInterval> | null = null

  function buildPayload(eventType: 'page_view' | 'page_ping'): AnalyticsPayload {
    const path = normalizedPath(window.location.pathname)
    const content = inferContent(path)
    return {
      eventId: randomId('evt'),
      eventType,
      occurredAt: new Date().toISOString(),
      visitorId: getVisitorId(),
      sessionId: getSessionId(),
      path,
      routeName: String(nuxtApp.$router.currentRoute.value.name ?? ''),
      pageTitle: document.title,
      referrer: previousReferrer,
      contentType: content.contentType,
      contentId: content.contentId,
      contentSlug: content.contentSlug,
      locale: document.documentElement.lang || navigator.language || '',
      screenWidth: window.screen?.width ?? 0,
      screenHeight: window.screen?.height ?? 0,
      viewportWidth: window.innerWidth,
      viewportHeight: window.innerHeight,
      timezone: Intl.DateTimeFormat().resolvedOptions().timeZone || '',
    }
  }

  function restartPing(): void {
    if (pingTimer) {
      window.clearInterval(pingTimer)
      pingTimer = null
    }
    if (document.visibilityState !== 'visible') return
    pingTimer = window.setInterval(() => {
      sendAnalytics([buildPayload('page_ping')])
    }, PING_INTERVAL_MS)
  }

  function trackPageView(): void {
    const payload = buildPayload('page_view')
    sendAnalytics([payload])
    currentPath = payload.path
    previousReferrer = `${window.location.origin}${payload.path}`
    restartPing()
  }

  nuxtApp.hook('page:finish', () => {
    window.setTimeout(() => {
      trackPageView()
    }, 0)
  })

  document.addEventListener('visibilitychange', () => {
    if (document.visibilityState === 'visible') {
      getSessionId()
      restartPing()
      return
    }
    if (pingTimer) {
      window.clearInterval(pingTimer)
      pingTimer = null
    }
  })
})
