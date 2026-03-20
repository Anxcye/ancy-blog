/**
 * File: comment.ts
 * Purpose: Define comment domain types for admin interaction moderation.
 * Module: frontend-admin-react/types, domain types layer.
 * Related: comments API module and InteractionPage.
 */

export type CommentStatus = 'pending' | 'approved' | 'rejected' | 'spam' | 'deleted';

export interface Comment {
  id: string;
  articleId: string;
  contentType: 'article' | 'moment';
  contentId: string;
  parentId?: string;
  rootId?: string;
  content: string;
  status: CommentStatus;
  isPinned: boolean;
  likeCount: number;
  replyCount: number;
  nickname: string;
  email?: string;
  website?: string;
  avatarUrl?: string;
  source: string;
  ip: string;
  userAgent?: string;
  riskScore: number;
  approvedAt?: string;
  approvedBy?: string;
  toCommentId?: string;
  toCommentNickname?: string;
  createdAt: string;
  updatedAt: string;
}

export interface CommentListParams {
  page?: number;
  pageSize?: number;
  status?: CommentStatus | '';
}

export interface CommentUpdatePayload {
  status?: CommentStatus;
  isPinned?: boolean;
}

export interface CommentReplyPayload {
  content: string;
}
