/**
 * File: dashboard.ts
 * Purpose: Provide aggregated data fetchers for the admin dashboard overview.
 * Module: frontend-admin-react/api, API layer.
 * Related: http client, articles/moments/comments API modules, DashboardPage.
 */

import { httpClient } from '../lib/http';
import type { ApiResponse, PaginatedData } from '../types/api';
import type { ArticleListItem } from '../types/article';
import type { Moment } from '../types/moment';
import type { Comment } from '../types/comment';

export interface DashboardStats {
    totalArticles: number;
    publishedArticles: number;
    draftArticles: number;
    totalMoments: number;
    pendingComments: number;
    totalComments: number;
}

/** Fetch aggregated dashboard statistics by parallel API calls. */
export async function fetchDashboardStats(): Promise<DashboardStats> {
    const [articlesAll, articlesDraft, momentsAll, commentsPending, commentsAll] = await Promise.all([
        httpClient
            .get<ApiResponse<PaginatedData<ArticleListItem>>>('/admin/articles', {
                params: { page: 1, pageSize: 1, contentKind: 'post' },
            })
            .then((r) => r.data.data?.total ?? 0)
            .catch(() => 0),

        httpClient
            .get<ApiResponse<PaginatedData<ArticleListItem>>>('/admin/articles', {
                params: { page: 1, pageSize: 1, status: 'draft' },
            })
            .then((r) => r.data.data?.total ?? 0)
            .catch(() => 0),

        httpClient
            .get<ApiResponse<PaginatedData<Moment>>>('/admin/moments', {
                params: { page: 1, pageSize: 1 },
            })
            .then((r) => r.data.data?.total ?? 0)
            .catch(() => 0),

        httpClient
            .get<ApiResponse<PaginatedData<Comment>>>('/admin/comments', {
                params: { page: 1, pageSize: 1, status: 'pending' },
            })
            .then((r) => r.data.data?.total ?? 0)
            .catch(() => 0),

        httpClient
            .get<ApiResponse<PaginatedData<Comment>>>('/admin/comments', {
                params: { page: 1, pageSize: 1 },
            })
            .then((r) => r.data.data?.total ?? 0)
            .catch(() => 0),
    ]);

    return {
        totalArticles: articlesAll,
        publishedArticles: articlesAll - articlesDraft,
        draftArticles: articlesDraft,
        totalMoments: momentsAll,
        pendingComments: commentsPending,
        totalComments: commentsAll,
    };
}

/** Fetch recent articles for the dashboard activity feed. */
export async function fetchRecentArticles(): Promise<ArticleListItem[]> {
    const res = await httpClient.get<ApiResponse<PaginatedData<ArticleListItem>>>('/admin/articles', {
        params: { page: 1, pageSize: 5 },
    });
    return res.data.data?.rows ?? [];
}

/** Fetch recent comments for the dashboard activity feed. */
export async function fetchRecentComments(): Promise<Comment[]> {
    const res = await httpClient.get<ApiResponse<PaginatedData<Comment>>>('/admin/comments', {
        params: { page: 1, pageSize: 5 },
    });
    return res.data.data?.rows ?? [];
}
