// File: contentMarkdown.ts
// Purpose: Provide a shared markdown-it renderer for public content fragments like moments.
// Module: frontend-blog/utils, presentation helper layer.
// Related: app/pages/moments/[[id]].vue and app/components/MomentDetailModal.vue.

import MarkdownIt from 'markdown-it'

const renderer = new MarkdownIt({
  html: false,
  breaks: true,
  linkify: true,
  typographer: true,
})

export function renderContentMarkdown(source: string): string {
  const text = source.trim()
  if (!text) {
    return ''
  }
  return renderer.render(text)
}
