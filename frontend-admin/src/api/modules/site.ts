// File: api/modules/site.ts
// Purpose: Provide site orchestration API calls for settings, footer, social links, and nav items.
// Module: frontend-admin/api/site, domain gateway layer.
// Related: SiteView and backend admin/public site handlers.
import { httpClient } from '@/api/http';
import type { ApiEnvelope, ContentSlot, FooterItem, NavItem, SiteSettings, SlotItem, SocialLink } from '@/api/types';

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

export async function updateFooterItem(id: string, payload: Omit<FooterItem, 'id'>): Promise<void> {
  await httpClient.put<ApiEnvelope<FooterItem>>(`/admin/site/footer-items/${id}`, payload);
}

export async function deleteFooterItem(id: string): Promise<void> {
  await httpClient.delete<ApiEnvelope<boolean>>(`/admin/site/footer-items/${id}`);
}

export async function listSocialLinks(): Promise<SocialLink[]> {
  const response = await httpClient.get<ApiEnvelope<SocialLink[]>>('/public/site/social-links');
  return response.data.data;
}

export async function createSocialLink(payload: Omit<SocialLink, 'id'>): Promise<string> {
  const response = await httpClient.post<ApiEnvelope<IDResponse>>('/admin/site/social-links', payload);
  return response.data.data.id;
}

export async function updateSocialLink(id: string, payload: Omit<SocialLink, 'id'>): Promise<void> {
  await httpClient.put<ApiEnvelope<SocialLink>>(`/admin/site/social-links/${id}`, payload);
}

export async function deleteSocialLink(id: string): Promise<void> {
  await httpClient.delete<ApiEnvelope<boolean>>(`/admin/site/social-links/${id}`);
}

export async function listNavItems(): Promise<NavItem[]> {
  const response = await httpClient.get<ApiEnvelope<NavItem[]>>('/public/site/nav');
  return response.data.data;
}

export async function createNavItem(payload: Omit<NavItem, 'id'>): Promise<string> {
  const response = await httpClient.post<ApiEnvelope<IDResponse>>('/admin/site/nav-items', payload);
  return response.data.data.id;
}

export async function updateNavItem(id: string, payload: Omit<NavItem, 'id'>): Promise<void> {
  await httpClient.put<ApiEnvelope<NavItem>>(`/admin/site/nav-items/${id}`, payload);
}

export async function deleteNavItem(id: string): Promise<void> {
  await httpClient.delete<ApiEnvelope<boolean>>(`/admin/site/nav-items/${id}`);
}

export async function listSlots(): Promise<ContentSlot[]> {
  const response = await httpClient.get<ApiEnvelope<ContentSlot[]>>('/admin/site/slots');
  return response.data.data;
}

export async function createSlot(payload: {
  slotKey: string;
  name: string;
  description: string;
  enabled: boolean;
}): Promise<string> {
  const response = await httpClient.post<ApiEnvelope<IDResponse>>('/admin/site/slots', payload);
  return response.data.data.id;
}

export async function listSlotItems(slotKey: string): Promise<SlotItem[]> {
  const response = await httpClient.get<ApiEnvelope<SlotItem[]>>(`/admin/site/slots/${slotKey}/items`);
  return response.data.data;
}

export async function createSlotItem(slotKey: string, payload: {
  contentType: string;
  contentId: string;
  orderNum: number;
  enabled: boolean;
}): Promise<string> {
  const response = await httpClient.post<ApiEnvelope<IDResponse>>(`/admin/site/slots/${slotKey}/items`, payload);
  return response.data.data.id;
}

export async function deleteSlotItem(slotKey: string, id: string): Promise<void> {
  await httpClient.delete<ApiEnvelope<boolean>>(`/admin/site/slots/${slotKey}/items/${id}`);
}
