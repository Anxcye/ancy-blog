// File: app/composables/useApi.ts
// Purpose: Typed API client wrapping $fetch with base URL and locale injection.
// Module: frontend-blog/composables, data layer.
// Related: nuxt.config.ts (runtimeConfig), all page components fetching data.

export interface ApiResponse<T> {
    code: string
    message: string
    data: T
}

export interface PaginatedData<T> {
    rows: T[]
    total: number
    page: number
    pageSize: number
}

// ── Article types ──────────────────────────────────────────────────
export type ArticleStatus = 'draft' | 'published' | 'scheduled' | 'archived'
export type ContentKind = 'post' | 'page'
export type CommentContentType = 'article' | 'moment'

export interface ArticleCard {
    id: string
    title: string
    slug: string
    summary: string
    contentKind: ContentKind
    status: ArticleStatus
    isPinned: boolean
    isFeatured: boolean
    coverImage?: string
    publishedAt?: string
    createdAt: string
    updatedAt: string
    categorySlug?: string
    tagSlugs?: string[]
    viewCount: number
}

export interface ArticleDetail extends ArticleCard {
    content: string // TipTap JSON string
    allowComment: boolean
    visibility: string
    originType: string
    sourceUrl?: string
    aiAssistLevel?: string
}

// ── Moment types ───────────────────────────────────────────────────
export interface Moment {
    id: string
    content: string
    status: string
    allowComment: boolean
    commentCount: number
    isPinned?: boolean
    publishedAt?: string
    createdAt: string
}

// ── Comment types ──────────────────────────────────────────────────
export interface Comment {
    id: string
    articleId?: string
    contentType: CommentContentType
    contentId: string
    parentId?: string
    rootId?: string
    toCommentId?: string
    content: string
    status: string
    isPinned: boolean
    likeCount: number
    replyCount: number
    nickname: string
    isAuthor: boolean
    website?: string
    avatarUrl?: string
    toCommentNickname?: string
    createdAt: string
}

export interface CommentThread extends Comment {
    children: CommentThread[]
}

export interface CommentCreatePayload {
    articleId?: string
    contentType: CommentContentType
    contentId: string
    parentId?: string
    rootId?: string
    toCommentId?: string
    content: string
    nickname: string
    email?: string
    website?: string
    avatarUrl?: string
}

export interface LinkSubmissionPayload {
    name: string
    url: string
    avatarUrl?: string
    description?: string
    contactEmail?: string
}

// ── Category / Tag ─────────────────────────────────────────────────
export interface Category { id: string; name: string; slug: string }
export interface Tag { id: string; name: string; slug: string }

// ── Site settings ──────────────────────────────────────────────────
export interface SiteSettings {
    siteName: string
    avatarUrl?: string
    faviconUrl?: string
    heroIntroMd?: string
    defaultLocale: string
    commentEnabled: boolean
    commentRequireApproval: boolean
    linkSubmissionEnabled: boolean
    siteDescription?: string
    seoKeywords?: string
    ogImageUrl?: string
}

// ── Timeline ───────────────────────────────────────────────────────
export interface TimelineItem {
    contentType: 'article' | 'moment'
    id: string
    title?: string
    slug?: string
    summary?: string
    categorySlug?: string
    categoryName?: string
    content?: string
    publishedAt: string
}

// ── Composable ────────────────────────────────────────────────────
export function useApi() {
    const config = useRuntimeConfig()
    const { locale } = useI18n()

    /**
   * Core fetch wrapper — prepends base URL, injects locale header.
   * Also intercepts system/business errors and triggers Nuxt error boundaries automatically.
   */
    async function apiFetch<T>(path: string, opts?: Parameters<typeof $fetch>[1]): Promise<T> {
        // SSR requires absolute path to connect to backend directly;
        // CSR points to current origin (`/api/v1`) to leverage Nuxt proxy and avoid CORS
        const baseURL = import.meta.client ? '/api/v1' : config.public.apiBase

        try {
            const res = await $fetch<ApiResponse<T>>(path, {
                baseURL,
                headers: {
                    'Accept-Language': locale.value === 'en' ? 'en-US' : 'zh-CN',
                    ...(opts?.headers as Record<string, string> ?? {}),
                },
                onResponseError({ response }) {
                    // HTTP 级别异常（500/404 等硬件隔离的系统异常）
                    throw createError({
                        statusCode: response.status,
                        message: response._data?.message || `Network error (${response.status})`,
                        fatal: true, // Fatal triggers full SSR boundary crash so `error.vue` handles it
                    })
                },
                ...opts,
            })

            // App 业务级异常（后端自己组装的 code / message）
            if (!res || !['SUCCESS', '0', 'ok', ''].includes((res.code || '').toLowerCase())) {
                throw createError({
                    statusCode: 400,
                    message: res?.message || 'Server returned invalid state',
                    fatal: true,
                })
            }

            return res.data
        } catch (err: any) {
            // Re-throw the structured nuxt error to boundary 
            throw err
        }
    }

    // ── Public APIs ────────────────────────────────────────────────
    return {
        /** Fetch site settings */
        getSiteSettings: () =>
            apiFetch<SiteSettings>('/public/site/settings'),

        /** Fetch paginated article list */
        listArticles: (params?: { page?: number; pageSize?: number; category?: string; tag?: string }) =>
            apiFetch<PaginatedData<ArticleCard>>('/public/articles', { params }),

        /** Fetch article detail by slug */
        getArticle: (slug: string) =>
            apiFetch<ArticleDetail>(`/public/articles/${slug}`, {
                params: locale.value === 'en' ? { locale: 'en-US' } : undefined,
            }),

        /** Fetch categories */
        getCategories: () =>
            apiFetch<Category[]>('/public/categories'),

        /** Fetch tags */
        getTags: () =>
            apiFetch<Tag[]>('/public/tags'),

        /** Fetch paginated moments */
        listMoments: (params?: { page?: number; pageSize?: number }) =>
            apiFetch<PaginatedData<Moment>>('/public/moments', { params }),

        /** Fetch moment detail */
        getMoment: (id: string) =>
            apiFetch<Moment>(`/public/moments/${id}`, {
                params: locale.value === 'en' ? { locale: 'en-US' } : undefined,
            }),

        /** Fetch mixed timeline */
        listTimeline: (params?: { page?: number; pageSize?: number }) =>
            apiFetch<PaginatedData<TimelineItem>>('/public/timeline', { params }),

        /** Fetch comments for an article */
        listComments: (contentType: CommentContentType, contentId: string, params?: { page?: number; pageSize?: number }) =>
            apiFetch<PaginatedData<CommentThread>>(`/public/comments/content/${contentType}/${contentId}`, { params }),

        /** Fetch total approved comment count */
        getCommentTotal: (contentType: CommentContentType, contentId: string) =>
            apiFetch<number>(`/public/comments/content/${contentType}/${contentId}/total`),

        /** Fetch children for a root comment */
        listCommentChildren: (commentId: string, params?: { page?: number; pageSize?: number }) =>
            apiFetch<PaginatedData<Comment>>(`/public/comments/${commentId}/children`, { params }),

        /** Submit a comment */
        createComment: (payload: CommentCreatePayload) =>
            apiFetch<{ id: string }>('/public/comments', { method: 'POST', body: payload }),

        /** Fetch navigation items */
        getNav: () =>
            apiFetch<Array<{ id: string; name: string; key: string; targetType: string; targetValue?: string; children?: any[] }>>('/public/site/nav'),

        /** Fetch social links */
        getSocialLinks: () =>
            apiFetch<Array<{ id: string; platform: string; title: string; url: string; iconKey?: string }>>('/public/site/social-links'),

        /** Fetch footer items */
        getFooter: () =>
            apiFetch<Record<string, Array<{ id: string; label: string; linkType: string; internalArticleSlug?: string; externalUrl?: string; rowNum: number; orderNum: number }>>>('/public/site/footer'),

        /** Fetch approved friends links */
        getApprovedLinks: () =>
            apiFetch<Array<{ id: string; name: string; url: string; avatarUrl?: string; description?: string }>>('/public/links'),

        /** Submit a friend link */
        submitLink: (payload: LinkSubmissionPayload) =>
            apiFetch<{ id: string }>('/public/links/submissions', { method: 'POST', body: payload }),
    }
}
