// File: api/types.ts
// Purpose: Define shared API envelope and pagination types for admin HTTP modules.
// Module: frontend-admin/api, transport typing layer.
// Related: api/http client, module request functions, view-level data mapping.
export interface ApiEnvelope<T> {
  code: string;
  message: string;
  data: T;
}

export interface PageResult<T> {
  total: number;
  rows: T[];
}

export interface Article {
  id: string;
  title: string;
  slug: string;
  contentKind: 'post' | 'page';
  summary: string;
  content: string;
  status: 'draft' | 'published' | 'scheduled';
  visibility: 'public' | 'private' | 'unlisted';
  allowComment: boolean;
  originType: 'original' | 'repost' | 'translated';
  sourceUrl: string;
  aiAssistLevel: string;
  coverImage: string;
  categorySlug: string;
  tagSlugs: string[];
  publishedAt: string;
  createdAt: string;
  updatedAt: string;
}

export interface ArticleUpsertPayload {
  title: string;
  slug: string;
  contentKind: 'post' | 'page';
  summary: string;
  content: string;
  status: 'draft' | 'published' | 'scheduled';
  visibility: 'public' | 'private' | 'unlisted';
  allowComment: boolean;
  originType: 'original' | 'repost' | 'translated';
  sourceUrl: string;
  aiAssistLevel: string;
  coverImage: string;
  categorySlug: string;
  tagSlugs: string[];
  publishedAt?: string;
}
