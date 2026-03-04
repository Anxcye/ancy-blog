// File: api/modules/moments.ts
// Purpose: Provide admin moment management API calls for list/create/update flows.
// Module: frontend-admin/api/moments, domain gateway layer.
// Related: MomentsView and backend admin moment handlers.
import { httpClient } from '@/api/http';
import type { ApiEnvelope, PageResult } from '@/api/types';

export interface Moment {
  id: string;
  content: string;
  status: 'draft' | 'published' | 'scheduled' | string;
  allowComment: boolean;
  publishedAt: string;
  createdAt: string;
  updatedAt: string;
}

interface IDResponse {
  id: string;
}

interface MomentPayload {
  content: string;
  status: string;
  allowComment: boolean;
  publishedAt?: string;
}

export async function listMoments(params: {
  page: number;
  pageSize: number;
  status?: string;
}): Promise<PageResult<Moment>> {
  const response = await httpClient.get<ApiEnvelope<PageResult<Moment>>>('/admin/moments', { params });
  return response.data.data;
}

export async function createMoment(payload: MomentPayload): Promise<string> {
  const response = await httpClient.post<ApiEnvelope<IDResponse>>('/admin/moments', payload);
  return response.data.data.id;
}

export async function updateMoment(id: string, payload: MomentPayload): Promise<void> {
  await httpClient.put<ApiEnvelope<Moment>>(`/admin/moments/${id}`, payload);
}
