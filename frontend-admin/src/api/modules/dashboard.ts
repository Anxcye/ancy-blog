// File: api/modules/dashboard.ts
// Purpose: Aggregate lightweight dashboard metrics from existing admin endpoints.
// Module: frontend-admin/api/dashboard, application gateway layer.
// Related: DashboardView and article/moment/comment/link API modules.
import { listArticles } from '@/api/modules/articles';
import { listComments, listLinkSubmissions } from '@/api/modules/interactions';
import { listMoments } from '@/api/modules/moments';

export interface DashboardMetrics {
  articleTotal: number;
  articleDraft: number;
  articlePublished: number;
  momentTotal: number;
  commentPending: number;
  linkPending: number;
}

export async function loadDashboardMetrics(): Promise<DashboardMetrics> {
  const [articleTotal, articleDraft, articlePublished, momentTotal, commentPending, linkPending] = await Promise.all([
    listArticles({ page: 1, pageSize: 1 }),
    listArticles({ page: 1, pageSize: 1, status: 'draft' }),
    listArticles({ page: 1, pageSize: 1, status: 'published' }),
    listMoments({ page: 1, pageSize: 1 }),
    listComments({ page: 1, pageSize: 1, status: 'pending' }),
    listLinkSubmissions({ page: 1, pageSize: 1, reviewStatus: 'pending' }),
  ]);

  return {
    articleTotal: articleTotal.total,
    articleDraft: articleDraft.total,
    articlePublished: articlePublished.total,
    momentTotal: momentTotal.total,
    commentPending: commentPending.total,
    linkPending: linkPending.total,
  };
}
