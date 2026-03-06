/**
 * File: site.ts
 * Purpose: Define site management domain types for settings, social, footer, and nav.
 * Module: frontend-admin-react/types, domain types layer.
 * Related: site API module and SitePage.
 */

// ── Site settings ──────────────────────────────
export interface SiteSettings {
  siteName: string;
  avatarUrl?: string;
  heroIntroMd?: string;
  defaultLocale: string;
  commentEnabled: boolean;
  commentRequireApproval: boolean;
  linkSubmissionEnabled: boolean;
  siteDescription?: string;
  seoKeywords?: string;
  ogImageUrl?: string;
  updatedAt?: string;
}

export type SiteSettingsPayload = Omit<SiteSettings, 'updatedAt'>;

// ── Social links ───────────────────────────────
export type SocialPlatform = 'github' | 'mail' | 'x' | 'linkedin' | 'custom';

export interface SocialLink {
  id: string;
  platform: SocialPlatform;
  title: string;
  url: string;
  iconKey?: string;
  orderNum: number;
  enabled: boolean;
  createdAt: string;
  updatedAt: string;
}

export interface SocialLinkFormValues {
  platform: SocialPlatform;
  title: string;
  url: string;
  iconKey?: string;
  orderNum?: number;
  enabled?: boolean;
}

// ── Footer items ───────────────────────────────
export type LinkType = 'none' | 'internal' | 'external';

export interface FooterItem {
  id: string;
  label: string;
  linkType: LinkType;
  internalArticleSlug?: string;
  externalUrl?: string;
  rowNum: 1 | 2 | 3;
  orderNum: number;
  enabled: boolean;
  createdAt: string;
  updatedAt: string;
}

export interface FooterItemFormValues {
  label: string;
  linkType: LinkType;
  internalArticleSlug?: string;
  externalUrl?: string;
  rowNum: 1 | 2 | 3;
  orderNum?: number;
  enabled?: boolean;
}

// ── Nav items ──────────────────────────────────
export type NavItemType = 'menu' | 'dropdown' | 'link';
export type NavTargetType = 'route' | 'category' | 'article' | 'external';

export interface NavItem {
  id: string;
  parentId?: string;
  name: string;
  key: string;
  type: NavItemType;
  targetType: NavTargetType;
  targetValue?: string;
  orderNum: number;
  enabled: boolean;
  children?: NavItem[];
  createdAt?: string;
  updatedAt?: string;
}

export interface NavItemFormValues {
  parentId?: string;
  name: string;
  key: string;
  type: NavItemType;
  targetType: NavTargetType;
  targetValue?: string;
  orderNum?: number;
  enabled?: boolean;
}

// ── Content slots ──────────────────────────────
export interface ContentSlot {
  id: string;
  slotKey: string;
  name: string;
  description?: string;
  enabled: boolean;
  createdAt: string;
  updatedAt: string;
}
