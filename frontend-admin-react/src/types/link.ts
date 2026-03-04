/**
 * File: link.ts
 * Purpose: Define friend-link domain types for admin review workflow.
 * Module: frontend-admin-react/types, domain types layer.
 * Related: links API module and InteractionPage.
 */

export type ReviewStatus = 'pending' | 'approved' | 'rejected';

export interface Link {
  id: string;
  name: string;
  url: string;
  avatarUrl?: string;
  description?: string;
  contactEmail?: string;
  reviewStatus: ReviewStatus;
  reviewNote?: string;
  submittedIp?: string;
  relatedArticleId?: string;
  approvedAt?: string;
  createdAt: string;
  updatedAt: string;
}

export interface LinkListParams {
  page?: number;
  pageSize?: number;
  reviewStatus?: ReviewStatus | '';
}

export interface LinkReviewPayload {
  reviewStatus: ReviewStatus;
  reviewNote?: string;
  relatedArticleId?: string;
}
