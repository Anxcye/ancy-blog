/**
 * File: moment.ts
 * Purpose: Define moment domain types for admin short-form content management.
 * Module: frontend-admin-react/types, domain types layer.
 * Related: moments API module and MomentsPage.
 */

export type MomentStatus = 'draft' | 'published' | 'archived';

export interface Moment {
  id: string;
  content: string;
  status: MomentStatus;
  isPinned: boolean;
  allowComment: boolean;
  publishedAt?: string;
  createdAt: string;
  updatedAt: string;
}

export interface MomentListParams {
  page?: number;
  pageSize?: number;
  status?: MomentStatus | '';
}

export interface MomentFormValues {
  content: string;
  status: MomentStatus;
}
