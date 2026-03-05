/**
 * File: tmdb-card-embed.ts
 * Purpose: TipTap extension for rendering TMDB movie/TV card embeds in blog frontend.
 * Module: frontend-blog/components/tiptap-extensions
 * Related: x-post-embed.ts, TiptapRenderer.vue
 */

import { Node } from '@tiptap/core'

export const TmdbCardEmbed = Node.create({
  name: 'tmdbCardEmbed',
  group: 'block',
  atom: true,

  addAttributes() {
    return {
      tmdbId: { default: '' },
      mediaType: { default: 'movie' },
      title: { default: '' },
      overview: { default: '' },
      posterPath: { default: '' },
      releaseDate: { default: '' },
      voteAverage: { default: 0 },
    }
  },

  parseHTML() {
    return [{ tag: 'div[data-embed-type="tmdb_card"]' }]
  },

  renderHTML({ node }) {
    const tmdbId = node.attrs.tmdbId || ''
    const mediaType = node.attrs.mediaType || 'movie'
    const title = node.attrs.title || `TMDB ${mediaType} #${tmdbId}`
    const overview = node.attrs.overview || ''
    const posterPath = node.attrs.posterPath || ''
    const releaseDate = node.attrs.releaseDate || ''
    const voteAverage = node.attrs.voteAverage || 0
    const url = `https://www.themoviedb.org/${mediaType}/${tmdbId}`
    const posterUrl = posterPath ? `https://image.tmdb.org/t/p/w342${posterPath}` : ''

    return [
      'a',
      {
        'data-embed-type': 'tmdb_card',
        'data-tmdb-id': tmdbId,
        'data-media-type': mediaType,
        class: 'tmdb-card-embed',
        href: url,
        target: '_blank',
        rel: 'noopener noreferrer',
      },
      ...(posterUrl ? [['img', { src: posterUrl, alt: title, class: 'tmdb-poster' }]] : []),
      [
        'div',
        { class: 'tmdb-info' },
        ['div', { class: 'tmdb-title' }, title],
        ['div', { class: 'tmdb-meta' }, `${mediaType === 'movie' ? '电影' : '剧集'} · ${releaseDate.slice(0, 4)} · ${voteAverage.toFixed(1)} 分`],
        ...(overview ? [['div', { class: 'tmdb-overview' }, overview]] : []),
      ],
    ]
  },
})
