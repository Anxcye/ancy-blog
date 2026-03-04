/**
 * File: links.ts
 * Purpose: Provide typed API call functions for admin friend-link review workflow.
 * Module: frontend-admin-react/api, API layer.
 * Related: http client, link types, and InteractionPage.
 */

import { httpClient } from '../lib/http';
import type { ApiResponse, PaginatedData } from '../types/api';
import type { Link, LinkListParams, LinkReviewPayload } from '../types/link';

export async function listLinks(params: LinkListParams): Promise<PaginatedData<Link>> {
  const clean = Object.fromEntries(
    Object.entries(params).filter(([, v]) => v !== '' && v !== undefined),
  );
  const res = await httpClient.get<ApiResponse<PaginatedData<Link>>>('/admin/links', {
    params: clean,
  });
  return res.data.data;
}

export async function reviewLink(id: string, payload: LinkReviewPayload): Promise<boolean> {
  const res = await httpClient.patch<ApiResponse<boolean>>(
    `/admin/links/${id}/review`,
    payload,
  );
  return res.data.data;
}
