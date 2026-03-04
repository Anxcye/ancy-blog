// File: embeds.ts
// Purpose: Define reusable TipTap block extensions for third-party embedded content.
// Module: frontend-admin/components/editor/extensions, editor schema layer.
// Related: RichTextEditor insertion actions and RichTextPreview HTML rendering.
import { Node, mergeAttributes } from '@tiptap/core';

export const XPostEmbed = Node.create({
  name: 'xPostEmbed',
  group: 'block',
  atom: true,
  selectable: true,

  addAttributes() {
    return {
      postId: {
        default: '',
      },
      author: {
        default: '',
      },
    };
  },

  parseHTML() {
    return [{ tag: 'div[data-embed-type="x_post"]' }];
  },

  renderHTML({ HTMLAttributes }) {
    const postId = String(HTMLAttributes.postId || '').trim();
    const author = String(HTMLAttributes.author || '').trim();
    return [
      'div',
      mergeAttributes(HTMLAttributes, {
        'data-embed-type': 'x_post',
        class: 'embed-card embed-card-x',
      }),
      ['div', { class: 'embed-card-title' }, 'X Post'],
      ['div', { class: 'embed-card-main' }, postId ? `Post: ${postId}` : 'Post: (empty)'],
      ['div', { class: 'embed-card-sub' }, author ? `Author: ${author}` : 'Author: -'],
    ];
  },
});

export const TmdbCardEmbed = Node.create({
  name: 'tmdbCardEmbed',
  group: 'block',
  atom: true,
  selectable: true,

  addAttributes() {
    return {
      mediaType: {
        default: 'movie',
      },
      tmdbId: {
        default: '',
      },
      title: {
        default: '',
      },
    };
  },

  parseHTML() {
    return [{ tag: 'div[data-embed-type="tmdb_card"]' }];
  },

  renderHTML({ HTMLAttributes }) {
    const mediaType = String(HTMLAttributes.mediaType || 'movie').trim() || 'movie';
    const tmdbID = String(HTMLAttributes.tmdbId || '').trim();
    const title = String(HTMLAttributes.title || '').trim();
    return [
      'div',
      mergeAttributes(HTMLAttributes, {
        'data-embed-type': 'tmdb_card',
        class: 'embed-card embed-card-tmdb',
      }),
      ['div', { class: 'embed-card-title' }, 'TMDB Card'],
      ['div', { class: 'embed-card-main' }, `${mediaType.toUpperCase()} · ${tmdbID || '(empty id)'}`],
      ['div', { class: 'embed-card-sub' }, title || 'Title: -'],
    ];
  },
});

