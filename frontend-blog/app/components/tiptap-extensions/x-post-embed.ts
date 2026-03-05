/**
 * File: x-post-embed.ts
 * Purpose: TipTap extension for rendering X (Twitter) post embeds in blog frontend.
 * Module: frontend-blog/components/tiptap-extensions
 * Related: tmdb-card-embed.ts, TiptapRenderer.vue
 */

import { Node } from '@tiptap/core'

export const XPostEmbed = Node.create({
  name: 'xPostEmbed',
  group: 'block',
  atom: true,

  addAttributes() {
    return {
      postId: { default: '' },
      author: { default: '' },
    }
  },

  parseHTML() {
    return [{ tag: 'div[data-embed-type="x_post"]' }]
  },

  renderHTML({ node }) {
    const postId = node.attrs.postId || ''
    const author = node.attrs.author || 'unknown'
    const url = `https://twitter.com/${author}/status/${postId}`
    return [
      'div',
      {
        'data-embed-type': 'x_post',
        'data-post-id': postId,
        'data-author': author,
        class: 'x-post-embed',
      },
      [
        'blockquote',
        { class: 'twitter-tweet', 'data-theme': 'light' },
        [
          'a',
          {
            href: url,
            target: '_blank',
            rel: 'noopener noreferrer',
          },
          '',
        ],
      ],
    ]
  },

  addGlobalAttributes() {
    return [
      {
        types: ['xPostEmbed'],
        attributes: {
          'data-twitter-script': {
            default: null,
            renderHTML: () => {
              if (typeof window !== 'undefined' && !window.twttr) {
                const script = document.createElement('script')
                script.src = 'https://platform.twitter.com/widgets.js'
                script.async = true
                document.body.appendChild(script)
              }
              return {}
            },
          },
        },
      },
    ]
  },
})
