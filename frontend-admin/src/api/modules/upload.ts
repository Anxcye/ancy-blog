// File: api/modules/upload.ts
// Purpose: Provide admin image upload API call for rich-text editor and media insertion flows.
// Module: frontend-admin/api/upload, domain gateway layer.
// Related: ArticleEditorView and backend admin upload endpoint.
import { httpClient } from '@/api/http';
import type { ApiEnvelope } from '@/api/types';

interface UploadImageResult {
  key: string;
  url: string;
}

export async function uploadImage(file: File): Promise<UploadImageResult> {
  const formData = new FormData();
  formData.append('file', file);
  const response = await httpClient.post<ApiEnvelope<UploadImageResult>>('/admin/upload/image', formData);
  return response.data.data;
}
