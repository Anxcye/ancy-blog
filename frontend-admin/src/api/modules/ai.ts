// File: api/modules/ai.ts
// Purpose: Provide AI-assist API calls for summary generation and slug suggestion.
// Module: frontend-admin/api/ai, domain gateway layer.
// Related: ArticleEditorView and backend admin AI handlers.
import { httpClient } from '@/api/http';
import type { ApiEnvelope } from '@/api/types';

interface SummaryData {
  summary: string;
  fallbackUsed: boolean;
}

interface SlugData {
  slug: string;
  fallbackUsed: boolean;
}

export async function generateSummary(payload: {
  title: string;
  content: string;
  providerKey?: string;
  modelName?: string;
  maxLength?: number;
}): Promise<SummaryData> {
  const response = await httpClient.post<ApiEnvelope<SummaryData>>('/admin/ai/summary', payload);
  return response.data.data;
}

export async function suggestSlug(payload: {
  title: string;
  providerKey?: string;
  modelName?: string;
}): Promise<SlugData> {
  const response = await httpClient.post<ApiEnvelope<SlugData>>('/admin/ai/slug', payload);
  return response.data.data;
}
