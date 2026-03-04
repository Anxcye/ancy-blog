/**
 * File: article.ts
 * Purpose: Define article domain types used by the admin content management module.
 * Module: frontend-admin-react/types, domain types layer.
 * Related: articles API module, ArticlesPage, and ArticleEditorPage.
 */

export type ArticleStatus = 'draft' | 'published' | 'scheduled' | 'archived';
export type ContentKind = 'post' | 'page';
export type Visibility = 'public' | 'unlisted' | 'private';
export type OriginType = 'original' | 'repost' | 'translation';

export interface ArticleListItem {
  id: string;
  title: string;
  slug: string;
  summary: string;
  contentKind: ContentKind;
  status: ArticleStatus;
  isPinned: boolean;
  coverImage?: string;
  publishedAt?: string;
  createdAt: string;
  updatedAt: string;
}

export interface Article {
  id: string;
  title: string;
  slug: string;
  summary: string;
  content: string;
  contentKind: ContentKind;
  status: ArticleStatus;
  visibility: Visibility;
  allowComment: boolean;
  isPinned: boolean;
  coverImage?: string;
  originType: OriginType;
  sourceUrl?: string;
  aiAssistLevel?: string;
  publishedAt?: string;
  createdAt: string;
  updatedAt: string;
}

export interface ArticleListParams {
  page?: number;
  pageSize?: number;
  status?: ArticleStatus | '';
  contentKind?: ContentKind | '';
  keyword?: string;
}

export interface ArticleFormValues {
  title: string;
  slug: string;
  summary?: string;
  content?: string;
  contentKind: ContentKind;
  status: ArticleStatus;
  visibility?: Visibility;
  allowComment?: boolean;
  coverImage?: string;
  originType?: OriginType;
  sourceUrl?: string;
  publishedAt?: string;
}
