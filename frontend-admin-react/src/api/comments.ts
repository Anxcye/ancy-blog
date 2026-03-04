/**
 * File: comments.ts
 * Purpose: Provide typed API call functions for admin comment moderation.
 * Module: frontend-admin-react/api, API layer.
 * Related: http client, comment types, and InteractionPage.
 */

import { httpClient } from '../lib/http';
import type { ApiResponse, PaginatedData } from '../types/api';
import type { Comment, CommentListParams, CommentUpdatePayload } from '../types/comment';

export async function listComments(params: CommentListParams): Promise<PaginatedData<Comment>> {
  const clean = Object.fromEntries(
    Object.entries(params).filter(([, v]) => v !== '' && v !== undefined),
  );
  const res = await httpClient.get<ApiResponse<PaginatedData<Comment>>>('/admin/comments', {
    params: clean,
  });
  return res.data.data;
}

export async function updateComment(id: string, payload: CommentUpdatePayload): Promise<Comment> {
  const res = await httpClient.put<ApiResponse<Comment>>(`/admin/comments/${id}`, payload);
  return res.data.data;
}
