/**
 * File: system.ts
 * Purpose: Define integration provider, translation job, and translation content domain types.
 * Module: frontend-admin-react/types, domain types layer.
 * Related: system API module and SystemPage.
 */

// ── Integration providers ──────────────────────────────
export type ProviderType = 'object_storage' | 'llm';

export interface ProviderMeta {
  lastTestAt?: string;
  lastTestOk?: boolean;
  lastTestMsg?: string;
  latencyMs?: number;
}

export interface IntegrationProvider {
  id: string;
  providerType: ProviderType;
  providerKey: string;
  name: string;
  enabled: boolean;
  configJson: Record<string, unknown>;
  metaJson?: ProviderMeta;
  createdAt: string;
  updatedAt: string;
}

export interface ProviderUpdatePayload {
  enabled: boolean;
  configJson: Record<string, unknown>;
}

export interface ProviderTestResult {
  ok: boolean;
  message: string;
  latencyMs: number;
}

// ── Translation jobs ───────────────────────────────────
export type TranslationJobStatus = 'queued' | 'running' | 'succeeded' | 'failed';
export type TranslationSourceType = 'article' | 'moment';

export interface TranslationJob {
  id: string;
  sourceType: TranslationSourceType;
  sourceId: string;
  sourceLocale: string;
  targetLocale: string;
  providerKey: string;
  modelName: string;
  status: TranslationJobStatus;
  errorMessage?: string;
  retryCount: number;
  maxRetries: number;
  autoPublish: boolean;
  createdAt: string;
  updatedAt: string;
  finishedAt?: string;
}

export interface CreateTranslationJobPayload {
  sourceType: TranslationSourceType;
  sourceId: string;
  sourceLocale: string;
  targetLocale: string;
  providerKey: string;
  modelName: string;
  maxRetries?: number;
  autoPublish?: boolean;
}

export interface TranslationJobListParams {
  page?: number;
  pageSize?: number;
  status?: TranslationJobStatus | '';
  sourceType?: TranslationSourceType | '';
  sourceId?: string;
}

// ── Translation contents ───────────────────────────────
export type TranslationContentStatus = 'draft' | 'published' | 'archived';

export interface TranslationContent {
  id: string;
  sourceType: TranslationSourceType;
  sourceId: string;
  locale: string;
  title?: string;
  summary?: string;
  content: string;
  status: TranslationContentStatus;
  publishedAt?: string;
  translatedByJobId?: string;
  createdAt: string;
  updatedAt: string;
}

export interface TranslationContentListParams {
  sourceType: TranslationSourceType;
  sourceId?: string;
  locale?: string;
  page?: number;
  pageSize?: number;
}

export interface UpdateTranslationContentPayload {
  sourceType: TranslationSourceType;
  sourceId: string;
  locale: string;
  title?: string;
  summary?: string;
  content: string;
  status?: TranslationContentStatus;
  publishedAt?: string;
  translatedByJobId?: string;
}
