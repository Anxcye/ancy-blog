/**
 * File: markdown.ts
 * Purpose: Provide a shared markdown-it renderer for admin preview surfaces.
 * Module: frontend-admin-react/utils, presentation helper layer.
 * Related: pages/content/MomentsPage and public blog markdown rendering.
 */

import MarkdownIt from 'markdown-it';

const renderer = new MarkdownIt({
  html: false,
  breaks: true,
  linkify: true,
  typographer: true,
});

export function renderMarkdown(source: string): string {
  const text = source.trim();
  if (!text) {
    return '';
  }
  return renderer.render(text);
}
