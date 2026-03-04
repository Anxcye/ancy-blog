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

export interface SiteSettings {
  siteName: string;
  avatarUrl: string;
  heroIntroMd: string;
  defaultLocale: string;
}

export interface FooterItem {
  id: string;
  label: string;
  linkType: 'internal' | 'external' | 'none';
  internalArticleSlug: string;
  externalUrl: string;
  rowNum: number;
  orderNum: number;
  enabled: boolean;
}

export interface SocialLink {
  id: string;
  platform: string;
  title: string;
  url: string;
  iconKey: string;
  orderNum: number;
  enabled: boolean;
}

export interface NavItem {
  id: string;
  name: string;
  key: string;
  type: string;
  targetType: string;
  targetValue: string;
  orderNum: number;
  enabled: boolean;
}

export interface IntegrationProvider {
  id: string;
  providerType: string;
  providerKey: string;
  name: string;
  enabled: boolean;
  configJson: Record<string, unknown>;
  metaJson: Record<string, unknown>;
  createdAt: string;
  updatedAt: string;
}
