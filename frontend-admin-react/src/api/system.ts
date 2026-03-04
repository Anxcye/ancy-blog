/**
 * File: system.ts
 * Purpose: Provide typed API call functions for integration providers and translation management.
 * Module: frontend-admin-react/api, API layer.
 * Related: http client, system types, and SystemPage.
 */

import { httpClient } from '../lib/http';
import type { ApiResponse, PaginatedData } from '../types/api';
import type {
  CreateTranslationJobPayload,
  IntegrationProvider,
  ProviderTestResult,
  ProviderType,
  ProviderUpdatePayload,
  TranslationContent,
  TranslationContentListParams,
  TranslationJob,
  TranslationJobListParams,
  UpdateTranslationContentPayload,
} from '../types/system';

// ── Integration providers ──────────────────────────────

export async function listProviders(providerType?: ProviderType): Promise<IntegrationProvider[]> {
  const res = await httpClient.get<ApiResponse<IntegrationProvider[]>>('/admin/integrations', {
    params: providerType ? { providerType } : undefined,
  });
  return res.data.data;
}

export async function updateProvider(
  providerKey: string,
  payload: ProviderUpdatePayload,
): Promise<boolean> {
  const res = await httpClient.put<ApiResponse<boolean>>(
    `/admin/integrations/${providerKey}`,
    payload,
  );
  return res.data.data;
}

export async function testProvider(providerKey: string): Promise<ProviderTestResult> {
  const res = await httpClient.post<ApiResponse<ProviderTestResult>>(
    `/admin/integrations/${providerKey}/test`,
  );
  return res.data.data;
}

// ── Translation jobs ───────────────────────────────────

export async function createTranslationJob(
  payload: CreateTranslationJobPayload,
): Promise<{ id: string }> {
  const res = await httpClient.post<ApiResponse<{ id: string }>>(
    '/admin/translations/jobs',
    payload,
  );
  return res.data.data;
}

export async function listTranslationJobs(
  params: TranslationJobListParams,
): Promise<PaginatedData<TranslationJob>> {
  const clean = Object.fromEntries(
    Object.entries(params).filter(([, v]) => v !== '' && v !== undefined),
  );
  const res = await httpClient.get<ApiResponse<PaginatedData<TranslationJob>>>(
    '/admin/translations/jobs',
    { params: clean },
  );
  return res.data.data;
}

export async function retryTranslationJob(id: string): Promise<TranslationJob> {
  const res = await httpClient.post<ApiResponse<TranslationJob>>(
    `/admin/translations/jobs/${id}/retry`,
  );
  return res.data.data;
}

// ── Translation contents ───────────────────────────────

export async function listTranslationContents(
  params: TranslationContentListParams,
): Promise<PaginatedData<TranslationContent>> {
  const clean = Object.fromEntries(
    Object.entries(params).filter(([, v]) => v !== '' && v !== undefined),
  );
  const res = await httpClient.get<ApiResponse<PaginatedData<TranslationContent>>>(
    '/admin/translations/contents',
    { params: clean },
  );
  return res.data.data;
}

export async function getTranslationContent(
  sourceType: string,
  sourceId: string,
  locale: string,
): Promise<TranslationContent> {
  const res = await httpClient.get<ApiResponse<TranslationContent>>(
    `/admin/translations/contents/${sourceType}/${sourceId}/${locale}`,
  );
  return res.data.data;
}

export async function updateTranslationContent(
  payload: UpdateTranslationContentPayload,
): Promise<TranslationContent> {
  const res = await httpClient.put<ApiResponse<TranslationContent>>(
    '/admin/translations/contents',
    payload,
  );
  return res.data.data;
}
