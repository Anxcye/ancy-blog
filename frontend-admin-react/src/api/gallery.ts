/**
 * Purpose: Provide typed API call functions for admin gallery photo and tag management.
 * Module: frontend-admin-react/api, API layer.
 * Related: http client, gallery types, GalleryPage, and GalleryEditorPage.
 */

import { httpClient } from '../lib/http';
import type { ApiResponse, PaginatedData } from '../types/api';
import type {
  GalleryPhoto,
  GalleryPhotoFormValues,
  GalleryPhotoListParams,
  GalleryTag,
} from '../types/gallery';

// --------------- Photos ---------------

export async function listGalleryPhotos(params: GalleryPhotoListParams): Promise<PaginatedData<GalleryPhoto>> {
  const clean = Object.fromEntries(
    Object.entries(params).filter(([, v]) => v !== '' && v !== undefined),
  );
  const res = await httpClient.get<ApiResponse<PaginatedData<GalleryPhoto>>>('/admin/gallery/photos', {
    params: clean,
  });
  return res.data.data;
}

export async function getGalleryPhoto(id: string): Promise<GalleryPhoto> {
  const res = await httpClient.get<ApiResponse<GalleryPhoto>>(`/admin/gallery/photos/${id}`);
  return res.data.data;
}

export async function createGalleryPhoto(payload: GalleryPhotoFormValues): Promise<{ id: string }> {
  const res = await httpClient.post<ApiResponse<{ id: string }>>('/admin/gallery/photos', payload);
  return res.data.data;
}

export async function updateGalleryPhoto(id: string, payload: GalleryPhotoFormValues): Promise<{ id: string }> {
  const res = await httpClient.put<ApiResponse<{ id: string }>>(`/admin/gallery/photos/${id}`, payload);
  return res.data.data;
}

export async function deleteGalleryPhoto(id: string): Promise<void> {
  await httpClient.delete(`/admin/gallery/photos/${id}`);
}

export async function batchUpdatePhotoStatus(ids: string[], status: string): Promise<{ count: number }> {
  const res = await httpClient.post<ApiResponse<{ count: number }>>('/admin/gallery/photos/batch-status', {
    ids,
    status,
  });
  return res.data.data;
}

export async function uploadGalleryPhoto(file: File): Promise<GalleryPhoto> {
  const formData = new FormData();
  formData.append('file', file);
  const res = await httpClient.post<ApiResponse<GalleryPhoto>>('/admin/gallery/photos/upload', formData, {
    headers: { 'Content-Type': 'multipart/form-data' },
    timeout: 120_000,
  });
  return res.data.data;
}

// --------------- Tags ---------------

export async function listGalleryTags(): Promise<GalleryTag[]> {
  const res = await httpClient.get<ApiResponse<GalleryTag[]>>('/admin/gallery/tags');
  return res.data.data ?? [];
}

export async function createGalleryTag(payload: { name: string; slug: string }): Promise<GalleryTag> {
  const res = await httpClient.post<ApiResponse<GalleryTag>>('/admin/gallery/tags', payload);
  return res.data.data;
}

export async function deleteGalleryTag(id: string): Promise<void> {
  await httpClient.delete(`/admin/gallery/tags/${id}`);
}
