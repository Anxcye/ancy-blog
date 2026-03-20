/**
 * File: analytics.ts
 * Purpose: Define admin analytics domain types for overview, pages, and raw visit records.
 * Module: frontend-admin-react/types/analytics, type layer.
 * Related: analytics API module and AnalyticsPage.
 */

export interface AnalyticsPathStat {
  path: string;
  contentType?: string;
  contentId?: string;
  contentSlug?: string;
  pageViews: number;
  uniqueVisitors: number;
  uniqueIPs: number;
  lastVisitedAt: string;
}

export interface AnalyticsReferrerStat {
  referrerHost: string;
  visits: number;
}

export interface AnalyticsDeviceStat {
  deviceType: string;
  visits: number;
}

export interface AnalyticsDailyStat {
  date: string;
  pageViews: number;
  uniqueVisitors: number;
  uniqueIPs: number;
}

export interface AnalyticsOverview {
  rangeStart: string;
  rangeEnd: string;
  pageViews: number;
  uniqueVisitors: number;
  uniqueIPs: number;
  uniqueSessions: number;
  topPaths: AnalyticsPathStat[];
  topReferrers: AnalyticsReferrerStat[];
  deviceBreakdown: AnalyticsDeviceStat[];
  daily: AnalyticsDailyStat[];
}

export interface AnalyticsVisit {
  id: string;
  eventId: string;
  eventType: string;
  occurredAt: string;
  receivedAt: string;
  visitorId: string;
  sessionId: string;
  path: string;
  routeName?: string;
  pageTitle?: string;
  referrer?: string;
  referrerHost?: string;
  contentType?: string;
  contentId?: string;
  contentSlug?: string;
  locale?: string;
  screenWidth?: number;
  screenHeight?: number;
  viewportWidth?: number;
  viewportHeight?: number;
  timezone?: string;
  ip: string;
  userAgent?: string;
  deviceType?: string;
  browserName?: string;
  osName?: string;
  isBot: boolean;
}

export interface AnalyticsPageListParams {
  page?: number;
  pageSize?: number;
  days?: number;
  path?: string;
  contentType?: string;
}

export interface AnalyticsVisitListParams extends AnalyticsPageListParams {
  eventType?: string;
  visitorId?: string;
  sessionId?: string;
  ip?: string;
  deviceType?: string;
  browserName?: string;
  osName?: string;
  isBot?: string;
}
