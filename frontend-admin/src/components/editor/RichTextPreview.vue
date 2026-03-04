<!--
File: RichTextPreview.vue
Purpose: Render TipTap JSON content into HTML for read-only preview in admin editor.
Module: frontend-admin/components/editor, editor preview layer.
Related: ArticleEditorView and RichTextEditor content JSON schema.
-->
<template>
  <div class="preview-render" v-html="html" />
</template>

<script setup lang="ts">
import { computed } from 'vue';
import type { JSONContent } from '@tiptap/core';
import { generateHTML } from '@tiptap/html';
import StarterKit from '@tiptap/starter-kit';
import Link from '@tiptap/extension-link';
import Image from '@tiptap/extension-image';
import { TmdbCardEmbed, XPostEmbed } from '@/components/editor/extensions/embeds';

const props = defineProps<{
  content: string;
}>();

function toDocContent(raw: string): JSONContent {
  if (!raw || !raw.trim()) {
    return { type: 'doc', content: [{ type: 'paragraph' }] };
  }
  try {
    const parsed = JSON.parse(raw) as JSONContent;
    if (typeof parsed === 'object' && parsed !== null && parsed.type === 'doc') {
      return parsed;
    }
  } catch {
    // Fallback to plain text paragraph.
  }
  return {
    type: 'doc',
    content: [
      {
        type: 'paragraph',
        content: [{ type: 'text', text: raw }],
      },
    ],
  };
}

const html = computed(() => {
  const doc = toDocContent(props.content);
  return generateHTML(doc, [StarterKit, Link, Image, XPostEmbed, TmdbCardEmbed]);
});
</script>

<style scoped>
.preview-render {
  line-height: 1.75;
  color: var(--n-text-color);
}

.preview-render :deep(h1),
.preview-render :deep(h2),
.preview-render :deep(h3) {
  margin: 1.1em 0 0.6em;
}

.preview-render :deep(p) {
  margin: 0.7em 0;
}

.preview-render :deep(code) {
  font-family: ui-monospace, SFMono-Regular, Menlo, monospace;
  background: color-mix(in srgb, var(--n-primary-color) 14%, transparent);
  border-radius: 6px;
  padding: 2px 6px;
}

.preview-render :deep(pre) {
  overflow: auto;
  padding: 10px;
  border-radius: 10px;
  background: var(--n-code-color);
}

.preview-render :deep(blockquote) {
  margin: 0.9em 0;
  padding: 0.2em 1em;
  border-left: 3px solid var(--n-primary-color);
  color: var(--n-text-color-2);
}

.preview-render :deep(a) {
  color: var(--n-primary-color);
}

.preview-render :deep(img) {
  max-width: 100%;
  height: auto;
  border-radius: 8px;
}

.preview-render :deep(.embed-card) {
  margin: 0.9em 0;
  border: 1px solid color-mix(in srgb, var(--n-primary-color) 30%, var(--n-border-color));
  border-radius: 12px;
  padding: 10px 12px;
  background: color-mix(in srgb, var(--n-card-color) 90%, var(--n-primary-color) 10%);
}

.preview-render :deep(.embed-card-title) {
  font-size: 12px;
  text-transform: uppercase;
  letter-spacing: 0.08em;
  color: var(--n-text-color-3);
}

.preview-render :deep(.embed-card-main) {
  margin-top: 6px;
  font-size: 14px;
  font-weight: 600;
  color: var(--n-text-color);
}

.preview-render :deep(.embed-card-sub) {
  margin-top: 4px;
  font-size: 13px;
  color: var(--n-text-color-2);
}
</style>
