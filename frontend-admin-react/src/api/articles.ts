/**
 * File: articles.ts
 * Purpose: Provide typed API call functions for admin article CRUD and batch operations.
 * Module: frontend-admin-react/api, API layer.
 * Related: http client, article types, ArticlesPage, and ArticleEditorPage.
 */

import { httpClient } from '../lib/http';
import type { ApiResponse, PaginatedData } from '../types/api';
import type { Article, ArticleFormValues, ArticleListItem, ArticleListParams } from '../types/article';

export interface Category { id: string; name: string; slug: string }
export interface Tag { id: string; name: string; slug: string }

export async function listCategories(): Promise<Category[]> {
  const res = await httpClient.get<ApiResponse<Category[]>>('/public/categories');
  return res.data.data ?? [];
}

export async function createCategory(payload: { name: string; slug: string }): Promise<Category> {
  const res = await httpClient.post<ApiResponse<Category>>('/admin/categories', payload);
  return res.data.data;
}

export async function deleteCategory(id: string): Promise<void> {
  await httpClient.delete(`/admin/categories/${id}`);
}

export async function listTags(): Promise<Tag[]> {
  const res = await httpClient.get<ApiResponse<Tag[]>>('/public/tags');
  return res.data.data ?? [];
}

export async function createTag(payload: { name: string; slug: string }): Promise<Tag> {
  const res = await httpClient.post<ApiResponse<Tag>>('/admin/tags', payload);
  return res.data.data;
}

export async function deleteTag(id: string): Promise<void> {
  await httpClient.delete(`/admin/tags/${id}`);
}

export async function listArticles(params: ArticleListParams): Promise<PaginatedData<ArticleListItem>> {
  // Strip empty string filters before sending to avoid backend validation errors
  const clean = Object.fromEntries(
    Object.entries(params).filter(([, v]) => v !== '' && v !== undefined),
  );
  const res = await httpClient.get<ApiResponse<PaginatedData<ArticleListItem>>>('/admin/articles', {
    params: clean,
  });
  return res.data.data;
}

export async function getArticle(id: string): Promise<Article> {
  const res = await httpClient.get<ApiResponse<Article>>(`/admin/articles/${id}`);
  return res.data.data;
}

export async function createArticle(payload: ArticleFormValues): Promise<{ id: string }> {
  const res = await httpClient.post<ApiResponse<{ id: string }>>('/admin/articles', payload);
  return res.data.data;
}

export async function updateArticle(id: string, payload: ArticleFormValues): Promise<{ id: string }> {
  const res = await httpClient.put<ApiResponse<{ id: string }>>(`/admin/articles/${id}`, payload);
  return res.data.data;
}

export async function deleteArticle(id: string): Promise<void> {
  await httpClient.delete(`/admin/articles/${id}`);
}

export async function batchStatusArticles(ids: string[], status: string): Promise<{ count: number }> {
  const res = await httpClient.post<ApiResponse<{ count: number }>>('/admin/articles/batch-status', {
    ids,
    status,
  });
  return res.data.data;
}

export async function batchDeleteArticles(ids: string[]): Promise<{ count: number }> {
  const res = await httpClient.post<ApiResponse<{ count: number }>>('/admin/articles/batch-delete', {
    ids,
  });
  return res.data.data;
}

export async function generateAiSummary(
  title: string,
  content: string,
): Promise<{ summary: string; fallbackUsed: boolean }> {
  const res = await httpClient.post<ApiResponse<{ summary: string; fallbackUsed: boolean }>>(
    '/admin/ai/summary',
    { title, content, maxLength: 180 },
  );
  return res.data.data;
}

export async function generateAiSlug(
  title: string,
): Promise<{ slug: string; fallbackUsed: boolean }> {
  const res = await httpClient.post<ApiResponse<{ slug: string; fallbackUsed: boolean }>>(
    '/admin/ai/slug',
    { title },
  );
  return res.data.data;
}
