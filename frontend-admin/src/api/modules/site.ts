// File: api/modules/site.ts
// Purpose: Provide site orchestration API calls for settings, footer, social links, and nav items.
// Module: frontend-admin/api/site, domain gateway layer.
// Related: SiteView and backend admin/public site handlers.
import { httpClient } from '@/api/http';
import type { ApiEnvelope, FooterItem, NavItem, SiteSettings, SocialLink } from '@/api/types';

interface IDResponse {
  id: string;
}

export async function getSiteSettings(): Promise<SiteSettings> {
  const response = await httpClient.get<ApiEnvelope<SiteSettings>>('/public/site/settings');
  return response.data.data;
}

export async function updateSiteSettings(payload: SiteSettings): Promise<SiteSettings> {
  const response = await httpClient.put<ApiEnvelope<SiteSettings>>('/admin/site/settings', payload);
  return response.data.data;
}

export async function listFooterItems(): Promise<FooterItem[]> {
  const response = await httpClient.get<ApiEnvelope<Record<string, FooterItem[]>>>('/public/site/footer');
  const grouped = response.data.data || {};
  return Object.values(grouped)
    .flat()
    .sort((a, b) => (a.rowNum - b.rowNum) || (a.orderNum - b.orderNum));
}

export async function createFooterItem(payload: Omit<FooterItem, 'id'>): Promise<string> {
  const response = await httpClient.post<ApiEnvelope<IDResponse>>('/admin/site/footer-items', payload);
  return response.data.data.id;
}

export async function listSocialLinks(): Promise<SocialLink[]> {
  const response = await httpClient.get<ApiEnvelope<SocialLink[]>>('/public/site/social-links');
  return response.data.data;
}

export async function createSocialLink(payload: Omit<SocialLink, 'id'>): Promise<string> {
  const response = await httpClient.post<ApiEnvelope<IDResponse>>('/admin/site/social-links', payload);
  return response.data.data.id;
}

export async function listNavItems(): Promise<NavItem[]> {
  const response = await httpClient.get<ApiEnvelope<NavItem[]>>('/public/site/nav');
  return response.data.data;
}

export async function createNavItem(payload: Omit<NavItem, 'id'>): Promise<string> {
  const response = await httpClient.post<ApiEnvelope<IDResponse>>('/admin/site/nav-items', payload);
  return response.data.data.id;
}
