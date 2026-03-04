<!--
File: RichTextEditor.vue
Purpose: Provide reusable rich-text editor with TipTap toolbar, link tools, and image upload insertion support.
Module: frontend-admin/components/editor, editor UI layer.
Related: ArticleEditorView, RichTextPreview, and admin upload API module.
-->
<template>
  <div class="editor-root">
    <div class="toolbar-shell">
      <div class="toolbar-group">
        <NTooltip trigger="hover"><template #trigger><NButton size="small" quaternary class="tool-btn" :class="{ 'is-active': isActive('bold') }" :disabled="isReadonly" @click="toggleBold">B</NButton></template><span>Bold</span></NTooltip>
        <NTooltip trigger="hover"><template #trigger><NButton size="small" quaternary class="tool-btn" :class="{ 'is-active': isActive('italic') }" :disabled="isReadonly" @click="toggleItalic">I</NButton></template><span>Italic</span></NTooltip>
        <NTooltip trigger="hover"><template #trigger><NButton size="small" quaternary class="tool-btn" :class="{ 'is-active': isActive('strike') }" :disabled="isReadonly" @click="toggleStrike">S</NButton></template><span>Strike</span></NTooltip>
        <NTooltip trigger="hover"><template #trigger><NButton size="small" quaternary class="tool-btn code-btn" :class="{ 'is-active': isActive('code') }" :disabled="isReadonly" @click="toggleInlineCode">&lt;/&gt;</NButton></template><span>Inline Code</span></NTooltip>
      </div>

      <div class="toolbar-group">
        <NTooltip trigger="hover"><template #trigger><NButton size="small" quaternary class="tool-btn" :class="{ 'is-active': isActive('heading', { level: 2 }) }" :disabled="isReadonly" @click="setHeading(2)">H2</NButton></template><span>Heading 2</span></NTooltip>
        <NTooltip trigger="hover"><template #trigger><NButton size="small" quaternary class="tool-btn" :class="{ 'is-active': isActive('heading', { level: 3 }) }" :disabled="isReadonly" @click="setHeading(3)">H3</NButton></template><span>Heading 3</span></NTooltip>
        <NTooltip trigger="hover"><template #trigger><NButton size="small" quaternary class="tool-btn" :class="{ 'is-active': isActive('bulletList') }" :disabled="isReadonly" @click="toggleBulletList">• List</NButton></template><span>Bullet List</span></NTooltip>
        <NTooltip trigger="hover"><template #trigger><NButton size="small" quaternary class="tool-btn" :class="{ 'is-active': isActive('orderedList') }" :disabled="isReadonly" @click="toggleOrderedList">1. List</NButton></template><span>Ordered List</span></NTooltip>
        <NTooltip trigger="hover"><template #trigger><NButton size="small" quaternary class="tool-btn" :class="{ 'is-active': isActive('blockquote') }" :disabled="isReadonly" @click="toggleBlockquote">Quote</NButton></template><span>Blockquote</span></NTooltip>
        <NTooltip trigger="hover"><template #trigger><NButton size="small" quaternary class="tool-btn" :class="{ 'is-active': isActive('codeBlock') }" :disabled="isReadonly" @click="toggleCodeBlock">Code</NButton></template><span>Code Block</span></NTooltip>
      </div>

      <div class="toolbar-group">
        <NPopover trigger="click" :show="linkPopoverVisible" @update:show="onLinkPopoverVisibleChange">
          <template #trigger>
            <NButton size="small" quaternary class="tool-btn" :class="{ 'is-active': isActive('link') }" :disabled="isReadonly" @click="prepareLinkInput">Link</NButton>
          </template>
          <div class="link-popover">
            <NInput v-model:value="linkInput" placeholder="https://example.com" size="small" />
            <div class="link-actions">
              <NButton size="small" type="primary" :disabled="isReadonly" @click="applyLink">Apply</NButton>
              <NButton size="small" tertiary :disabled="isReadonly" @click="removeLink">Remove</NButton>
            </div>
          </div>
        </NPopover>
        <NTooltip trigger="hover"><template #trigger><NButton size="small" quaternary class="tool-btn" :loading="uploading" :disabled="isReadonly" @click="pickImage">Image</NButton></template><span>Upload Image</span></NTooltip>
        <NTooltip trigger="hover"><template #trigger><NButton size="small" quaternary class="tool-btn" :disabled="isReadonly" @click="insertHorizontalRule">Divider</NButton></template><span>Horizontal Rule</span></NTooltip>
      </div>

      <div class="toolbar-group">
        <NTooltip trigger="hover"><template #trigger><NButton size="small" quaternary class="tool-btn" :disabled="isReadonly" @click="insertXPostEmbed">X Post</NButton></template><span>Insert X Post Block</span></NTooltip>
        <NTooltip trigger="hover"><template #trigger><NButton size="small" quaternary class="tool-btn" :disabled="isReadonly" @click="insertTmdbCardEmbed">TMDB</NButton></template><span>Insert TMDB Card Block</span></NTooltip>
      </div>

      <div class="toolbar-group">
        <NTooltip trigger="hover"><template #trigger><NButton size="small" quaternary class="tool-btn" :disabled="isReadonly || !canUndo" @click="undo">Undo</NButton></template><span>Undo</span></NTooltip>
        <NTooltip trigger="hover"><template #trigger><NButton size="small" quaternary class="tool-btn" :disabled="isReadonly || !canRedo" @click="redo">Redo</NButton></template><span>Redo</span></NTooltip>
        <NTooltip trigger="hover"><template #trigger><NButton size="small" quaternary class="tool-btn clear-btn" :disabled="isReadonly" @click="clearContent">Clear</NButton></template><span>Clear Content</span></NTooltip>
      </div>
    </div>

    <div class="editor-shell">
      <EditorContent :editor="editor" />
    </div>

    <input ref="imageInputRef" class="hidden-input" type="file" accept="image/*" @change="onImageSelected" />
  </div>
</template>

<script setup lang="ts">
import { computed, onBeforeUnmount, ref, watch } from 'vue';
import { EditorContent, useEditor } from '@tiptap/vue-3';
import type { JSONContent } from '@tiptap/core';
import StarterKit from '@tiptap/starter-kit';
import Link from '@tiptap/extension-link';
import Image from '@tiptap/extension-image';
import Placeholder from '@tiptap/extension-placeholder';
import { NButton, NInput, NPopover, NTooltip } from 'naive-ui';
import { TmdbCardEmbed, XPostEmbed } from '@/components/editor/extensions/embeds';

const props = withDefaults(
  defineProps<{
    modelValue: string;
    placeholder?: string;
    disabled?: boolean;
    uploadImage?: (file: File) => Promise<string>;
  }>(),
  {
    placeholder: '',
    disabled: false,
    uploadImage: undefined,
  },
);

const emit = defineEmits<{
  (e: 'update:modelValue', value: string): void;
  (e: 'upload-error', message: string): void;
}>();

const imageInputRef = ref<HTMLInputElement | null>(null);
const uploading = ref(false);
const linkPopoverVisible = ref(false);
const linkInput = ref('');

function emptyDoc(): JSONContent {
  return { type: 'doc', content: [{ type: 'paragraph' }] };
}

function toDocContent(raw: string): JSONContent {
  if (!raw || !raw.trim()) {
    return emptyDoc();
  }
  try {
    const parsed = JSON.parse(raw) as JSONContent;
    if (typeof parsed === 'object' && parsed !== null && parsed.type === 'doc') {
      return parsed;
    }
  } catch {
    // Fallback to paragraph text content below.
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

function serializeEditor(): string {
  if (!editor.value) {
    return JSON.stringify(emptyDoc());
  }
  return JSON.stringify(editor.value.getJSON());
}

const editor = useEditor({
  content: toDocContent(props.modelValue),
  editable: !props.disabled,
  extensions: [
    StarterKit,
    Link.configure({ openOnClick: false }),
    Image,
    XPostEmbed,
    TmdbCardEmbed,
    Placeholder.configure({
      placeholder: props.placeholder,
      showOnlyWhenEditable: false,
    }),
  ],
  onUpdate() {
    emit('update:modelValue', serializeEditor());
  },
});

const isReadonly = computed(() => props.disabled || !editor.value);
const canUndo = computed(() => editor.value?.can().chain().focus().undo().run() ?? false);
const canRedo = computed(() => editor.value?.can().chain().focus().redo().run() ?? false);

watch(
  () => props.modelValue,
  (value) => {
    if (!editor.value) {
      return;
    }
    const current = serializeEditor();
    const next = JSON.stringify(toDocContent(value));
    if (current !== next) {
      editor.value.commands.setContent(toDocContent(value), { emitUpdate: false });
    }
  },
);

watch(
  () => props.disabled,
  (value) => {
    editor.value?.setEditable(!value);
  },
);

function isActive(name: string, attrs?: Record<string, unknown>): boolean {
  return editor.value?.isActive(name, attrs) ?? false;
}

function toggleBold(): void {
  editor.value?.chain().focus().toggleBold().run();
}

function toggleItalic(): void {
  editor.value?.chain().focus().toggleItalic().run();
}

function toggleStrike(): void {
  editor.value?.chain().focus().toggleStrike().run();
}

function toggleInlineCode(): void {
  editor.value?.chain().focus().toggleCode().run();
}

function setHeading(level: 2 | 3): void {
  editor.value?.chain().focus().toggleHeading({ level }).run();
}

function toggleBulletList(): void {
  editor.value?.chain().focus().toggleBulletList().run();
}

function toggleOrderedList(): void {
  editor.value?.chain().focus().toggleOrderedList().run();
}

function toggleBlockquote(): void {
  editor.value?.chain().focus().toggleBlockquote().run();
}

function toggleCodeBlock(): void {
  editor.value?.chain().focus().toggleCodeBlock().run();
}

function prepareLinkInput(): void {
  const href = editor.value?.getAttributes('link').href;
  linkInput.value = typeof href === 'string' ? href : '';
}

function onLinkPopoverVisibleChange(value: boolean): void {
  linkPopoverVisible.value = value;
}

function applyLink(): void {
  const trimmed = linkInput.value.trim();
  if (!trimmed) {
    removeLink();
    return;
  }
  editor.value?.chain().focus().extendMarkRange('link').setLink({ href: trimmed }).run();
  linkPopoverVisible.value = false;
}

function removeLink(): void {
  editor.value?.chain().focus().extendMarkRange('link').unsetLink().run();
  linkPopoverVisible.value = false;
}

function clearContent(): void {
  editor.value?.commands.setContent(emptyDoc());
}

function pickImage(): void {
  imageInputRef.value?.click();
}

function insertHorizontalRule(): void {
  editor.value?.chain().focus().setHorizontalRule().run();
}

function insertXPostEmbed(): void {
  const postID = window.prompt('X post id');
  if (postID === null) {
    return;
  }
  const author = window.prompt('X author (optional)') || '';
  editor.value
    ?.chain()
    .focus()
    .insertContent({
      type: 'xPostEmbed',
      attrs: {
        postId: postID.trim(),
        author: author.trim(),
      },
    })
    .run();
}

function insertTmdbCardEmbed(): void {
  const tmdbID = window.prompt('TMDB id');
  if (tmdbID === null) {
    return;
  }
  const mediaTypeInput = (window.prompt('TMDB type: movie or tv') || 'movie').trim().toLowerCase();
  const mediaType = mediaTypeInput === 'tv' ? 'tv' : 'movie';
  const title = window.prompt('Title (optional)') || '';
  editor.value
    ?.chain()
    .focus()
    .insertContent({
      type: 'tmdbCardEmbed',
      attrs: {
        mediaType,
        tmdbId: tmdbID.trim(),
        title: title.trim(),
      },
    })
    .run();
}

function undo(): void {
  editor.value?.chain().focus().undo().run();
}

function redo(): void {
  editor.value?.chain().focus().redo().run();
}

async function onImageSelected(event: Event): Promise<void> {
  const input = event.target as HTMLInputElement;
  const file = input.files?.[0];
  input.value = '';
  if (!file) {
    return;
  }
  if (!props.uploadImage) {
    emit('upload-error', 'Image uploader is not configured.');
    return;
  }

  uploading.value = true;
  try {
    const imageURL = await props.uploadImage(file);
    editor.value?.chain().focus().setImage({ src: imageURL }).run();
  } catch (error) {
    const message = error instanceof Error ? error.message : 'Image upload failed.';
    emit('upload-error', message);
  } finally {
    uploading.value = false;
  }
}

onBeforeUnmount(() => {
  editor.value?.destroy();
});
</script>

<style scoped>
.editor-root {
  display: grid;
  gap: 10px;
}

.toolbar-shell {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  align-items: center;
  padding: 8px;
  border: 1px solid var(--n-border-color);
  border-radius: 12px;
  background: color-mix(in srgb, var(--n-card-color) 92%, var(--n-primary-color) 8%);
  position: sticky;
  top: 8px;
  z-index: 12;
  backdrop-filter: blur(6px);
}

.toolbar-group {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 4px;
  border-radius: 10px;
  background: var(--n-color);
  border: 1px solid var(--n-border-color);
}

.tool-btn {
  min-width: 44px;
  font-weight: 600;
  border-radius: 9px;
  border: 1px solid color-mix(in srgb, var(--n-border-color) 75%, transparent);
  background: color-mix(in srgb, var(--n-card-color) 92%, var(--n-primary-color) 8%);
  transition: transform 140ms ease, border-color 140ms ease, background-color 140ms ease, color 140ms ease;
}

.tool-btn:hover {
  transform: translateY(-1px);
  border-color: color-mix(in srgb, var(--n-primary-color) 48%, var(--n-border-color));
  background: color-mix(in srgb, var(--n-card-color) 85%, var(--n-primary-color) 15%);
}

.tool-btn:active {
  transform: translateY(0);
}

.tool-btn.is-active {
  color: #fff;
  border-color: color-mix(in srgb, var(--n-primary-color) 85%, transparent);
  background: linear-gradient(
    135deg,
    color-mix(in srgb, var(--n-primary-color) 78%, #ffffff 22%) 0%,
    color-mix(in srgb, var(--n-primary-color) 92%, #000000 8%) 100%
  );
}

.code-btn {
  font-family: ui-monospace, SFMono-Regular, Menlo, monospace;
}

.clear-btn {
  min-width: 56px;
}

.link-popover {
  width: 260px;
  display: grid;
  gap: 8px;
}

.link-actions {
  display: flex;
  justify-content: flex-end;
  gap: 6px;
}

.editor-shell {
  border: 1px solid var(--n-border-color);
  border-radius: 12px;
  background: var(--n-input-color);
  min-height: 380px;
  padding: 14px;
  box-shadow: inset 0 1px 0 color-mix(in srgb, var(--n-primary-color) 10%, transparent);
}

.editor-shell :deep(.ProseMirror) {
  min-height: 348px;
  outline: none;
  line-height: 1.75;
  font-size: 15px;
  color: var(--n-text-color);
}

.editor-shell :deep(.ProseMirror h1),
.editor-shell :deep(.ProseMirror h2),
.editor-shell :deep(.ProseMirror h3) {
  margin: 1.1em 0 0.6em;
  line-height: 1.3;
}

.editor-shell :deep(.ProseMirror p) {
  margin: 0.75em 0;
}

.editor-shell :deep(.ProseMirror blockquote) {
  margin: 0.9em 0;
  padding: 0.2em 1em;
  border-left: 3px solid var(--n-primary-color);
  color: var(--n-text-color-2);
}

.editor-shell :deep(.ProseMirror pre) {
  margin: 0.9em 0;
  overflow: auto;
  padding: 10px;
  border-radius: 10px;
  background: var(--n-code-color);
}

.editor-shell :deep(.ProseMirror code) {
  font-family: ui-monospace, SFMono-Regular, Menlo, monospace;
}

.editor-shell :deep(.ProseMirror p.is-editor-empty:first-child::before) {
  content: attr(data-placeholder);
  color: var(--n-placeholder-color);
  pointer-events: none;
  float: left;
  height: 0;
}

.editor-shell :deep(.ProseMirror img) {
  max-width: 100%;
  height: auto;
  border-radius: 10px;
  box-shadow: 0 6px 20px color-mix(in srgb, var(--n-primary-color) 18%, transparent);
}

.editor-shell :deep(.ProseMirror .embed-card) {
  margin: 0.9em 0;
  border: 1px solid color-mix(in srgb, var(--n-primary-color) 30%, var(--n-border-color));
  border-radius: 12px;
  padding: 10px 12px;
  background: color-mix(in srgb, var(--n-card-color) 90%, var(--n-primary-color) 10%);
}

.editor-shell :deep(.ProseMirror .embed-card-title) {
  font-size: 12px;
  text-transform: uppercase;
  letter-spacing: 0.08em;
  color: var(--n-text-color-3);
}

.editor-shell :deep(.ProseMirror .embed-card-main) {
  margin-top: 6px;
  font-size: 14px;
  font-weight: 600;
  color: var(--n-text-color);
}

.editor-shell :deep(.ProseMirror .embed-card-sub) {
  margin-top: 4px;
  font-size: 13px;
  color: var(--n-text-color-2);
}

.hidden-input {
  display: none;
}

@media (max-width: 720px) {
  .toolbar-shell {
    position: static;
  }

  .editor-shell {
    min-height: 320px;
    padding: 10px;
  }

  .editor-shell :deep(.ProseMirror) {
    min-height: 300px;
    font-size: 14px;
  }
}
</style>
