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
  parentId?: string;
  name: string;
  key: string;
  type: string;
  targetType: string;
  targetValue: string;
  orderNum: number;
  enabled: boolean;
  children?: NavItem[];
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

export interface ContentSlot {
  id: string;
  slotKey: string;
  name: string;
  description: string;
  enabled: boolean;
}

export interface SlotItem {
  id: string;
  slotKey: string;
  contentType: 'article' | 'moment' | string;
  contentId: string;
  orderNum: number;
  enabled: boolean;
}

export interface Comment {
  id: string;
  articleId: string;
  content: string;
  status: 'approved' | 'pending' | 'rejected' | string;
  isPinned: string;
  nickname: string;
  email: string;
  website: string;
  avatarUrl: string;
  source: string;
  ip: string;
  userAgent: string;
  createdAt: string;
  updatedAt: string;
}

export interface LinkSubmission {
  id: string;
  name: string;
  url: string;
  avatarUrl: string;
  description: string;
  contactEmail: string;
  reviewStatus: 'pending' | 'approved' | 'rejected' | string;
  reviewNote: string;
  relatedArticleId: string;
  submittedIp: string;
  submittedUserAgent: string;
  createdAt: string;
  updatedAt: string;
}

export interface TranslationJob {
  id: string;
  sourceType: 'article' | 'moment' | string;
  sourceId: string;
  sourceLocale: string;
  targetLocale: string;
  providerKey: string;
  modelName: string;
  status: 'queued' | 'running' | 'succeeded' | 'failed' | string;
  errorMessage: string;
  resultText: string;
  requestedBy: string;
  retryCount: number;
  maxRetries: number;
  nextRetryAt: string;
  autoPublish: boolean;
  publishAt: string;
  createdAt: string;
  updatedAt: string;
  finishedAt: string;
}

export interface TranslationContent {
  sourceType: 'article' | 'moment' | string;
  sourceId: string;
  locale: string;
  title: string;
  summary: string;
  content: string;
  status: 'draft' | 'published' | string;
  publishedAt: string;
  translatedByJobId: string;
  createdAt: string;
  updatedAt: string;
}
