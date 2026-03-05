// File: app/composables/useArticleSeo.ts
// Purpose: Enterprise SEO hook to generate JSON-LD structured data and OpenGraph tags for articles.

import type { ArticleDetail, SiteSettings } from '~/composables/useApi'

export function useArticleSeo(article: ArticleDetail, site: SiteSettings | null) {
    const config = useRuntimeConfig()
    const route = useRoute()

    // Full canonical URL
    const url = `${config.public.apiBase.replace('/api/v1', '')}${route.path}`

    // Schema.org JSON-LD for Articles
    const jsonLd = {
        '@context': 'https://schema.org',
        '@type': 'BlogPosting',
        headline: article.title,
        description: article.summary,
        image: article.coverImage || site?.ogImageUrl || '',
        datePublished: article.publishedAt || article.createdAt,
        dateModified: article.updatedAt || article.createdAt,
        author: {
            '@type': 'Person',
            name: site?.siteName || 'Author',
            url: config.public.apiBase.replace('/api/v1', ''),
        },
        publisher: {
            '@type': 'Organization',
            name: site?.siteName || 'Blog',
            logo: {
                '@type': 'ImageObject',
                url: site?.avatarUrl || '',
            }
        },
        mainEntityOfPage: {
            '@type': 'WebPage',
            '@id': url
        }
    }

    useHead({
        script: [
            {
                type: 'application/ld+json',
                innerHTML: JSON.stringify(jsonLd)
            }
        ]
    })

    // Open Graph / Twitter Cards
    useSeoMeta({
        title: article.title,
        description: article.summary,
        ogTitle: article.title,
        ogDescription: article.summary,
        ogImage: article.coverImage || site?.ogImageUrl,
        ogUrl: url,
        ogType: 'article',
        twitterCard: 'summary_large_image',
        twitterTitle: article.title,
        twitterDescription: article.summary,
        twitterImage: article.coverImage || site?.ogImageUrl,
    })
}
