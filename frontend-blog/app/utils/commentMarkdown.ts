// File: commentMarkdown.ts
// Purpose: Provide a shared markdown-it renderer for comment display and preview.
// Module: frontend-blog/utils, presentation helper layer.
// Related: app/components/CommentForm.vue and app/components/CommentItem.vue.

import MarkdownIt from 'markdown-it'

const renderer = new MarkdownIt({
  html: false,
  breaks: true,
  linkify: true,
  typographer: true,
})

export function renderCommentMarkdown(source: string): string {
  const text = source.trim()
  if (!text) {
    return ''
  }
  return renderer.render(text)
}
