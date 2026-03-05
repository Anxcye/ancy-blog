<!-- File: components/TiptapRenderer.vue
     Purpose: Render TipTap JSON content using read-only editor.
     Module: frontend-blog/components
-->
<template>
  <ClientOnly>
    <EditorContent v-if="editor" :editor="editor" />
  </ClientOnly>
</template>

<script setup lang="ts">
import { useEditor, EditorContent } from '@tiptap/vue-3'
import StarterKit from '@tiptap/starter-kit'
import Link from '@tiptap/extension-link'
import Image from '@tiptap/extension-image'
import TextAlign from '@tiptap/extension-text-align'
import { TextStyle } from '@tiptap/extension-text-style'
import { Color } from '@tiptap/extension-color'
import Highlight from '@tiptap/extension-highlight'
import TaskList from '@tiptap/extension-task-list'
import TaskItem from '@tiptap/extension-task-item'
import Underline from '@tiptap/extension-underline'
import Subscript from '@tiptap/extension-subscript'
import Superscript from '@tiptap/extension-superscript'
import Typography from '@tiptap/extension-typography'
import { XPostEmbed } from './tiptap-extensions/x-post-embed'
import { TmdbCardEmbed } from './tiptap-extensions/tmdb-card-embed'

const props = defineProps<{ content: string }>()

const editor = useEditor({
  editable: false,
  content: props.content ? JSON.parse(props.content) : null,
  extensions: [
    StarterKit,
    Link.configure({ openOnClick: false }),
    Image,
    TextAlign.configure({ types: ['heading', 'paragraph'] }),
    TextStyle,
    Color,
    Highlight.configure({ multicolor: true }),
    TaskList,
    TaskItem.configure({ nested: true }),
    Underline,
    Subscript,
    Superscript,
    Typography,
    XPostEmbed,
    TmdbCardEmbed,
  ],
})

// Load Twitter widgets script
onMounted(() => {
  if (typeof window !== 'undefined') {
    const w = window as any
    if (!w.twttr) {
      const script = document.createElement('script')
      script.src = 'https://platform.twitter.com/widgets.js'
      script.async = true
      script.onload = () => {
        if (w.twttr?.widgets) {
          w.twttr.widgets.load()
        }
      }
      document.body.appendChild(script)
    } else if (w.twttr?.widgets) {
      setTimeout(() => w.twttr.widgets.load(), 100)
    }
  }
})

watch(() => props.content, (newContent) => {
  if (editor.value && newContent) {
    editor.value.commands.setContent(JSON.parse(newContent))
  }
})

onBeforeUnmount(() => {
  editor.value?.destroy()
})
</script>

<style>
/* Import TipTap base styles */
.tiptap {
  outline: none;
}

.tiptap p {
  margin: 1em 0;
}

.tiptap h1,
.tiptap h2,
.tiptap h3,
.tiptap h4 {
  font-weight: 600;
  margin-top: 2em;
  margin-bottom: 0.5em;
}

.tiptap h1 { font-size: 1.875rem; }
.tiptap h2 { font-size: 1.5rem; }
.tiptap h3 { font-size: 1.25rem; }
.tiptap h4 { font-size: 1.1rem; }

.tiptap ul,
.tiptap ol {
  margin: 1.5em 0;
  padding-left: 1.5em;
}

.tiptap ul {
  list-style: disc;
}

.tiptap ol {
  list-style: decimal;
}

.tiptap li {
  margin: 0.5em 0;
}

.tiptap ul[data-type="taskList"] {
  list-style: none;
  padding-left: 0;
}

.tiptap ul[data-type="taskList"] li {
  display: flex;
  gap: 0.5em;
  margin: 0.5em 0;
}

.tiptap ul[data-type="taskList"] li label {
  display: flex;
  flex-shrink: 0;
  padding-top: 0.1em;
}

.tiptap ul[data-type="taskList"] li label input[type="checkbox"] {
  width: 1em;
  height: 1em;
  margin: 0;
  cursor: pointer;
  accent-color: var(--accent);
}

.tiptap ul[data-type="taskList"] li label span {
  display: none;
}

.tiptap ul[data-type="taskList"] li > div {
  flex: 1;
}

.tiptap ul[data-type="taskList"] li > div p {
  margin: 0;
  line-height: 1.6;
}

.tiptap ul[data-type="taskList"] li[data-checked="true"] > div p {
  text-decoration: line-through;
  opacity: 0.6;
  color: var(--text-muted);
}

.tiptap blockquote {
  margin: 1.5em 0;
  padding-left: 1.2em;
  border-left: 3px solid var(--accent);
  color: var(--text-muted);
  font-style: italic;
}

.tiptap pre {
  margin: 1.5em 0;
  padding: 1em;
  background: var(--bg-secondary);
  border: 1px solid var(--border);
  border-radius: var(--radius-md);
  overflow-x: auto;
}

.tiptap code {
  font-family: 'Fira Code', monospace;
  font-size: 0.9em;
  padding: 0.2em 0.4em;
  background: var(--bg-secondary);
  border-radius: 4px;
}

.tiptap pre code {
  padding: 0;
  background: none;
}

.tiptap img {
  max-width: 100%;
  height: auto;
  border-radius: var(--radius-md);
}

.tiptap a {
  color: var(--accent-text);
  text-decoration: underline;
}

.tiptap hr {
  margin: 2em 0;
  border: none;
  border-top: 1px solid var(--border);
}

.tiptap strong {
  font-weight: 600;
}

.tiptap mark {
  padding: 0.1em 0.2em;
  border-radius: 3px;
  /* Background color comes from inline style attribute */
}

.tiptap u {
  text-decoration: underline;
}

.tiptap s {
  text-decoration: line-through;
}

.tiptap sub {
  vertical-align: sub;
  font-size: 0.75em;
}

.tiptap sup {
  vertical-align: super;
  font-size: 0.75em;
}

/* X Post Embed */
.tiptap .x-post-embed {
  margin: 2em 0;
}

/* TMDB Card Embed */
.tiptap .tmdb-card-embed {
  display: block;
  margin: 2em 0;
  border: 1px solid var(--border);
  border-radius: var(--radius-lg);
  overflow: hidden;
  background: var(--surface);
  text-decoration: none;
  transition: all 0.2s;
}

.tiptap .tmdb-card-embed:hover {
  border-color: var(--accent);
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.08);
  transform: translateY(-2px);
}

.tiptap .tmdb-card-content {
  display: flex;
  align-items: stretch;
  position: relative;
}

.tiptap .tmdb-poster {
  width: 80px;
  height: 120px;
  object-fit: cover;
  flex-shrink: 0;
  display: block;
  vertical-align: top;
  mask-image: linear-gradient(to right, rgba(0,0,0,1) 0%, rgba(0,0,0,1) 70%, rgba(0,0,0,0.8) 85%, rgba(0,0,0,0) 100%);
  -webkit-mask-image: linear-gradient(to right, rgba(0,0,0,1) 0%, rgba(0,0,0,1) 70%, rgba(0,0,0,0.8) 85%, rgba(0,0,0,0) 100%);
}

.tiptap .tmdb-info {
  flex: 1;
  padding: 14px 20px;
  display: flex;
  flex-direction: column;
  gap: 6px;
  justify-content: center;
}

.tiptap .tmdb-title {
  font-weight: 600;
  font-size: 1.1em;
  line-height: 1.3;
  color: var(--text);
}

.tiptap .tmdb-meta {
  font-size: 0.85em;
  color: var(--text-muted);
}

.tiptap .tmdb-overview {
  font-size: 0.85em;
  line-height: 1.5;
  color: var(--text-subtle);
  display: -webkit-box;
  -webkit-line-clamp: 2;
  line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}
</style>
