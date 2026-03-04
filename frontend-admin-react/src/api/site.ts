/**
 * File: site.ts
 * Purpose: Provide typed API call functions for admin site settings and content management.
 * Module: frontend-admin-react/api, API layer.
 * Related: http client, site types, and SitePage.
 *
 * NOTE: Read endpoints for footer/social/nav use public routes as admin-specific GET list
 *       endpoints are not yet defined in the API contract. Add ADM-SITE-GET-* when available.
 */

import { httpClient } from '../lib/http';
import type { ApiResponse } from '../types/api';
import type {
  ContentSlot,
  FooterItem,
  FooterItemFormValues,
  NavItem,
  NavItemFormValues,
  SiteSettings,
  SiteSettingsPayload,
  SocialLink,
  SocialLinkFormValues,
} from '../types/site';

// ── Settings ──────────────────────────────────
export async function getSiteSettings(): Promise<SiteSettings> {
  const res = await httpClient.get<ApiResponse<SiteSettings>>('/public/site/settings');
  return res.data.data;
}

export async function updateSiteSettings(payload: SiteSettingsPayload): Promise<boolean> {
  const res = await httpClient.put<ApiResponse<boolean>>('/admin/site/settings', payload);
  return res.data.data;
}

// ── Social links ──────────────────────────────
export async function listSocialLinks(): Promise<SocialLink[]> {
  const res = await httpClient.get<ApiResponse<SocialLink[]>>('/public/site/social-links');
  return res.data.data;
}

export async function createSocialLink(payload: SocialLinkFormValues): Promise<{ id: string }> {
  const res = await httpClient.post<ApiResponse<{ id: string }>>(
    '/admin/site/social-links',
    payload,
  );
  return res.data.data;
}

export async function updateSocialLink(
  id: string,
  payload: SocialLinkFormValues,
): Promise<boolean> {
  const res = await httpClient.put<ApiResponse<boolean>>(
    `/admin/site/social-links/${id}`,
    payload,
  );
  return res.data.data;
}

export async function deleteSocialLink(id: string): Promise<boolean> {
  const res = await httpClient.delete<ApiResponse<boolean>>(`/admin/site/social-links/${id}`);
  return res.data.data;
}

// ── Footer items ──────────────────────────────
export async function listFooterItems(): Promise<FooterItem[]> {
  // Public endpoint returns items grouped by rowNum; flatten all rows
  const res = await httpClient.get<ApiResponse<Record<string, FooterItem[]>>>(
    '/public/site/footer',
  );
  return Object.values(res.data.data).flat();
}

export async function createFooterItem(payload: FooterItemFormValues): Promise<{ id: string }> {
  const res = await httpClient.post<ApiResponse<{ id: string }>>(
    '/admin/site/footer-items',
    payload,
  );
  return res.data.data;
}

export async function updateFooterItem(
  id: string,
  payload: FooterItemFormValues,
): Promise<boolean> {
  const res = await httpClient.put<ApiResponse<boolean>>(
    `/admin/site/footer-items/${id}`,
    payload,
  );
  return res.data.data;
}

export async function deleteFooterItem(id: string): Promise<boolean> {
  const res = await httpClient.delete<ApiResponse<boolean>>(`/admin/site/footer-items/${id}`);
  return res.data.data;
}

// ── Nav items ─────────────────────────────────
export async function listNavItems(): Promise<NavItem[]> {
  const res = await httpClient.get<ApiResponse<NavItem[]>>('/public/site/nav');
  return res.data.data;
}

export async function createNavItem(payload: NavItemFormValues): Promise<{ id: string }> {
  const res = await httpClient.post<ApiResponse<{ id: string }>>(
    '/admin/site/nav-items',
    payload,
  );
  return res.data.data;
}

export async function updateNavItem(id: string, payload: NavItemFormValues): Promise<boolean> {
  const res = await httpClient.put<ApiResponse<boolean>>(
    `/admin/site/nav-items/${id}`,
    payload,
  );
  return res.data.data;
}

export async function deleteNavItem(id: string): Promise<boolean> {
  const res = await httpClient.delete<ApiResponse<boolean>>(`/admin/site/nav-items/${id}`);
  return res.data.data;
}

// ── Translation policy ────────────────────────
export interface TranslationPolicy {
  enabled: boolean;
  targetLocales: string[];
  providerKey: string;
  autoPublish: boolean;
}

export async function getTranslationPolicy(): Promise<TranslationPolicy> {
  const res = await httpClient.get<ApiResponse<TranslationPolicy>>(
    '/admin/site/translation-policy',
  );
  return res.data.data;
}

export async function updateTranslationPolicy(payload: TranslationPolicy): Promise<TranslationPolicy> {
  const res = await httpClient.put<ApiResponse<TranslationPolicy>>(
    '/admin/site/translation-policy',
    payload,
  );
  return res.data.data;
}

// ── Content slots ─────────────────────────────
export async function listSlots(): Promise<ContentSlot[]> {
  const res = await httpClient.get<ApiResponse<ContentSlot[]>>('/admin/site/slots');
  return res.data.data;
}
