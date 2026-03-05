// File: api/modules/interactions.ts
// Purpose: Provide comment moderation and friend-link review API calls for admin interaction workflows.
// Module: frontend-admin/api/interactions, domain gateway layer.
// Related: InteractionView and backend admin comment/link handlers.
import { httpClient } from '@/api/http';
import type { ApiEnvelope, Comment, LinkSubmission, PageResult } from '@/api/types';

export async function listComments(params: {
  page: number;
  pageSize: number;
  status?: string;
}): Promise<PageResult<Comment>> {
  const response = await httpClient.get<ApiEnvelope<PageResult<Comment>>>('/admin/comments', { params });
  return response.data.data;
}

export async function updateComment(
  id: string,
  payload: {
    status: string;
    isPinned: string;
  },
): Promise<void> {
  await httpClient.put<ApiEnvelope<Comment>>(`/admin/comments/${id}`, payload);
}

export async function listLinkSubmissions(params: {
  page: number;
  pageSize: number;
  reviewStatus?: string;
}): Promise<PageResult<LinkSubmission>> {
  const response = await httpClient.get<ApiEnvelope<PageResult<LinkSubmission>>>('/admin/links', { params });
  return response.data.data;
}

export async function reviewLink(
  id: string,
  payload: {
    reviewStatus: string;
    reviewNote: string;
    relatedArticleId: string;
  },
): Promise<void> {
  await httpClient.patch<ApiEnvelope<boolean>>(`/admin/links/${id}/review`, payload);
}
