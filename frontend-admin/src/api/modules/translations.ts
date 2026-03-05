// File: api/modules/translations.ts
// Purpose: Provide translation job and translation content management API calls.
// Module: frontend-admin/api/translations, domain gateway layer.
// Related: SystemView translation tab and backend translation handlers.
import { httpClient } from '@/api/http';
import type { ApiEnvelope, PageResult, TranslationContent, TranslationJob } from '@/api/types';

interface IDResponse {
  id: string;
}

export async function createTranslationJob(payload: {
  sourceType: string;
  sourceId: string;
  sourceLocale: string;
  targetLocale: string;
  providerKey: string;
  modelName: string;
  maxRetries: number;
  autoPublish: boolean;
  publishAt?: string;
}): Promise<string> {
  const response = await httpClient.post<ApiEnvelope<IDResponse>>('/admin/translations/jobs', payload);
  return response.data.data.id;
}

export async function listTranslationJobs(params: {
  page: number;
  pageSize: number;
  status?: string;
  sourceType?: string;
  sourceId?: string;
}): Promise<PageResult<TranslationJob>> {
  const response = await httpClient.get<ApiEnvelope<PageResult<TranslationJob>>>('/admin/translations/jobs', { params });
  return response.data.data;
}

export async function retryTranslationJob(id: string): Promise<void> {
  await httpClient.post<ApiEnvelope<TranslationJob>>(`/admin/translations/jobs/${id}/retry`);
}

export async function listTranslationContents(params: {
  page: number;
  pageSize: number;
  sourceType?: string;
  sourceId?: string;
  locale?: string;
}): Promise<PageResult<TranslationContent>> {
  const response = await httpClient.get<ApiEnvelope<PageResult<TranslationContent>>>('/admin/translations/contents', { params });
  return response.data.data;
}

export async function getTranslationContent(sourceType: string, sourceId: string, locale: string): Promise<TranslationContent> {
  const response = await httpClient.get<ApiEnvelope<TranslationContent>>(`/admin/translations/contents/${sourceType}/${sourceId}/${locale}`);
  return response.data.data;
}

export async function upsertTranslationContent(payload: {
  sourceType: string;
  sourceId: string;
  locale: string;
  title: string;
  summary: string;
  content: string;
  status: string;
  publishedAt?: string;
  translatedByJobId?: string;
}): Promise<TranslationContent> {
  const response = await httpClient.put<ApiEnvelope<TranslationContent>>('/admin/translations/contents', payload);
  return response.data.data;
}
