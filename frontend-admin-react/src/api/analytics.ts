/**
 * File: analytics.ts
 * Purpose: Provide admin analytics API helpers for overview, page aggregates, and raw visits.
 * Module: frontend-admin-react/api/analytics, API layer.
 * Related: shared http client, analytics types, and AnalyticsPage.
 */

import { httpClient } from '../lib/http';
import type { ApiResponse, PaginatedData } from '../types/api';
import type {
  AnalyticsOverview,
  AnalyticsPageListParams,
  AnalyticsPathStat,
  AnalyticsVisit,
  AnalyticsVisitListParams,
} from '../types/analytics';

export async function getAnalyticsOverview(days = 7): Promise<AnalyticsOverview> {
  const res = await httpClient.get<ApiResponse<AnalyticsOverview>>('/admin/analytics/overview', {
    params: { days },
  });
  return res.data.data;
}

export async function listAnalyticsPages(params: AnalyticsPageListParams): Promise<PaginatedData<AnalyticsPathStat>> {
  const res = await httpClient.get<ApiResponse<PaginatedData<AnalyticsPathStat>>>('/admin/analytics/pages', {
    params,
  });
  return res.data.data;
}

export async function listAnalyticsVisits(params: AnalyticsVisitListParams): Promise<PaginatedData<AnalyticsVisit>> {
  const res = await httpClient.get<ApiResponse<PaginatedData<AnalyticsVisit>>>('/admin/analytics/visits', {
    params,
  });
  return res.data.data;
}
