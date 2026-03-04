/**
 * File: moments.ts
 * Purpose: Provide typed API call functions for admin moment CRUD and batch operations.
 * Module: frontend-admin-react/api, API layer.
 * Related: http client, moment types, and MomentsPage.
 */

import { httpClient } from '../lib/http';
import type { ApiResponse, PaginatedData } from '../types/api';
import type { Moment, MomentFormValues, MomentListParams } from '../types/moment';

export async function listMoments(params: MomentListParams): Promise<PaginatedData<Moment>> {
  const clean = Object.fromEntries(
    Object.entries(params).filter(([, v]) => v !== '' && v !== undefined),
  );
  const res = await httpClient.get<ApiResponse<PaginatedData<Moment>>>('/admin/moments', {
    params: clean,
  });
  return res.data.data;
}

export async function createMoment(payload: MomentFormValues): Promise<{ id: string }> {
  const res = await httpClient.post<ApiResponse<{ id: string }>>('/admin/moments', payload);
  return res.data.data;
}

export async function updateMoment(id: string, payload: MomentFormValues): Promise<Moment> {
  const res = await httpClient.put<ApiResponse<Moment>>(`/admin/moments/${id}`, payload);
  return res.data.data;
}

export async function deleteMoment(id: string): Promise<void> {
  await httpClient.delete(`/admin/moments/${id}`);
}

export async function batchStatusMoments(ids: string[], status: string): Promise<{ count: number }> {
  const res = await httpClient.post<ApiResponse<{ count: number }>>('/admin/moments/batch-status', {
    ids,
    status,
  });
  return res.data.data;
}

export async function batchDeleteMoments(ids: string[]): Promise<{ count: number }> {
  const res = await httpClient.post<ApiResponse<{ count: number }>>('/admin/moments/batch-delete', {
    ids,
  });
  return res.data.data;
}
