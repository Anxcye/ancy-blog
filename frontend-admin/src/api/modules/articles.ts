// File: api/modules/articles.ts
// Purpose: Provide article management API calls for admin list and editor flows.
// Module: frontend-admin/api/articles, domain gateway layer.
// Related: ArticlesView, ArticleEditorView, backend admin article handlers.
import { httpClient } from '@/api/http';
import type { ApiEnvelope, Article, ArticleUpsertPayload, PageResult } from '@/api/types';

export interface ListArticlesParams {
  page: number;
  pageSize: number;
  status?: string;
  contentKind?: string;
  keyword?: string;
}

interface IDResponse {
  id: string;
}

export async function listArticles(params: ListArticlesParams): Promise<PageResult<Article>> {
  const response = await httpClient.get<ApiEnvelope<PageResult<Article>>>('/admin/articles', { params });
  return response.data.data;
}

export async function getArticle(id: string): Promise<Article> {
  const response = await httpClient.get<ApiEnvelope<Article>>(`/admin/articles/${id}`);
  return response.data.data;
}

export async function createArticle(payload: ArticleUpsertPayload): Promise<string> {
  const response = await httpClient.post<ApiEnvelope<IDResponse>>('/admin/articles', payload);
  return response.data.data.id;
}

export async function updateArticle(id: string, payload: ArticleUpsertPayload): Promise<string> {
  const response = await httpClient.put<ApiEnvelope<IDResponse>>(`/admin/articles/${id}`, payload);
  return response.data.data.id;
}

export async function deleteArticle(id: string): Promise<void> {
  await httpClient.delete<ApiEnvelope<boolean>>(`/admin/articles/${id}`);
}

export async function batchUpdateArticleStatus(ids: string[], status: string): Promise<number> {
  const response = await httpClient.post<ApiEnvelope<{ affected: number }>>('/admin/articles/batch-status', { ids, status });
  return response.data.data.affected;
}

export async function batchDeleteArticles(ids: string[]): Promise<number> {
  const response = await httpClient.post<ApiEnvelope<{ affected: number }>>('/admin/articles/batch-delete', { ids });
  return response.data.data.affected;
}
