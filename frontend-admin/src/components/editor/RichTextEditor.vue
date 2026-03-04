<!--
File: RichTextEditor.vue
Purpose: Provide reusable rich-text editor with component-library-driven toolbar and extensible embed insertion.
Module: frontend-admin/components/editor, editor UI layer.
Related: ArticleEditorView, RichTextPreview, upload API module, and embed extensions.
-->
<template>
  <div class="editor-root">
    <div class="toolbar-shell">
      <div class="toolbar-row">
        <NButtonGroup size="small">
          <NTooltip trigger="hover"><template #trigger><NButton class="tool-btn" :class="{ 'is-active': isActive('bold') }" :disabled="isReadonly" @click="toggleBold"><template #icon><NIcon><TextOutline /></NIcon></template>B</NButton></template><span>Bold</span></NTooltip>
          <NTooltip trigger="hover"><template #trigger><NButton class="tool-btn" :class="{ 'is-active': isActive('italic') }" :disabled="isReadonly" @click="toggleItalic"><template #icon><NIcon><TextOutline /></NIcon></template>I</NButton></template><span>Italic</span></NTooltip>
          <NTooltip trigger="hover"><template #trigger><NButton class="tool-btn" :class="{ 'is-active': isActive('strike') }" :disabled="isReadonly" @click="toggleStrike"><template #icon><NIcon><TextOutline /></NIcon></template>S</NButton></template><span>Strike</span></NTooltip>
          <NTooltip trigger="hover"><template #trigger><NButton class="tool-btn" :class="{ 'is-active': isActive('code') }" :disabled="isReadonly" @click="toggleInlineCode"><template #icon><NIcon><CodeSlashOutline /></NIcon></template></NButton></template><span>Inline Code</span></NTooltip>
        </NButtonGroup>

        <NDropdown trigger="click" :options="blockMenuOptions" @select="onSelectBlockMenu">
          <NButton size="small" class="menu-btn" :disabled="isReadonly">
            <template #icon><NIcon><ListOutline /></NIcon></template>
            Block
            <NIcon class="menu-caret"><ChevronDownOutline /></NIcon>
          </NButton>
        </NDropdown>

        <NPopover trigger="click" :show="linkPopoverVisible" @update:show="onLinkPopoverVisibleChange">
          <template #trigger>
            <NButton size="small" class="tool-btn" :class="{ 'is-active': isActive('link') }" :disabled="isReadonly" @click="prepareLinkInput">
              <template #icon><NIcon><LinkOutline /></NIcon></template>
            </NButton>
          </template>
          <div class="link-popover">
            <NInput v-model:value="linkInput" placeholder="https://example.com" size="small" />
            <div class="link-actions">
              <NButton size="small" type="primary" :disabled="isReadonly" @click="applyLink">Apply</NButton>
              <NButton size="small" tertiary :disabled="isReadonly" @click="removeLink">Remove</NButton>
            </div>
          </div>
        </NPopover>

        <NTooltip trigger="hover"><template #trigger><NButton size="small" class="tool-btn" :loading="uploading" :disabled="isReadonly" @click="pickImage"><template #icon><NIcon><ImageOutline /></NIcon></template></NButton></template><span>Upload Image</span></NTooltip>

        <NDropdown trigger="click" :options="embedMenuOptions" @select="onSelectEmbedMenu">
          <NButton size="small" class="menu-btn" :disabled="isReadonly">
            <template #icon><NIcon><AddCircleOutline /></NIcon></template>
            Embed
            <NIcon class="menu-caret"><ChevronDownOutline /></NIcon>
          </NButton>
        </NDropdown>

        <div class="toolbar-spacer" />

        <NButtonGroup size="small">
          <NTooltip trigger="hover"><template #trigger><NButton class="tool-btn" :disabled="isReadonly || !canUndo" @click="undo"><template #icon><NIcon><ArrowUndoOutline /></NIcon></template></NButton></template><span>Undo</span></NTooltip>
          <NTooltip trigger="hover"><template #trigger><NButton class="tool-btn" :disabled="isReadonly || !canRedo" @click="redo"><template #icon><NIcon><ArrowRedoOutline /></NIcon></template></NButton></template><span>Redo</span></NTooltip>
          <NTooltip trigger="hover"><template #trigger><NButton class="tool-btn clear-btn" :disabled="isReadonly" @click="clearContent"><template #icon><NIcon><TrashOutline /></NIcon></template></NButton></template><span>Clear Content</span></NTooltip>
        </NButtonGroup>
      </div>
    </div>

    <div class="editor-shell">
      <EditorContent :editor="editor" />
    </div>

    <input ref="imageInputRef" class="hidden-input" type="file" accept="image/*" @change="onImageSelected" />

    <NModal v-model:show="embedModal.show" preset="card" :title="embedModalTitle" style="width: min(520px, 96vw)">
      <NForm label-placement="top">
        <template v-if="embedModal.type === 'x_post'">
          <NFormItem label="Post ID">
            <NInput v-model:value="embedModal.xPost.postId" placeholder="e.g. 1860000000000000000" />
          </NFormItem>
          <NFormItem label="Author (optional)">
            <NInput v-model:value="embedModal.xPost.author" placeholder="@username" />
          </NFormItem>
        </template>

        <template v-else-if="embedModal.type === 'tmdb_card'">
          <NFormItem label="TMDB Type">
            <NSelect v-model:value="embedModal.tmdb.mediaType" :options="tmdbTypeOptions" />
          </NFormItem>
          <NFormItem label="TMDB ID">
            <NInput v-model:value="embedModal.tmdb.tmdbId" placeholder="e.g. 603" />
          </NFormItem>
          <NFormItem label="Title (optional)">
            <NInput v-model:value="embedModal.tmdb.title" placeholder="The Matrix" />
          </NFormItem>
        </template>
      </NForm>

      <template #footer>
        <div class="modal-actions">
          <NButton @click="embedModal.show = false">Cancel</NButton>
          <NButton type="primary" :disabled="isReadonly" @click="confirmEmbedInsert">Insert</NButton>
        </div>
      </template>
    </NModal>
  </div>
</template>

<script setup lang="ts">
import { computed, h, onBeforeUnmount, reactive, ref, watch } from 'vue';
import { EditorContent, useEditor } from '@tiptap/vue-3';
import type { JSONContent } from '@tiptap/core';
import StarterKit from '@tiptap/starter-kit';
import Link from '@tiptap/extension-link';
import Image from '@tiptap/extension-image';
import Placeholder from '@tiptap/extension-placeholder';
import type { DropdownOption } from 'naive-ui';
import { NButton, NButtonGroup, NDropdown, NForm, NFormItem, NIcon, NInput, NModal, NPopover, NSelect, NTooltip } from 'naive-ui';
import {
  AddCircleOutline,
  ArrowRedoOutline,
  ArrowUndoOutline,
  ChevronDownOutline,
  CodeSlashOutline,
  ImageOutline,
  LinkOutline,
  ListOutline,
  LogoTwitter,
  FilmOutline,
  RemoveOutline,
  TextOutline,
  TrashOutline,
} from '@vicons/ionicons5';

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

const embedModal = reactive({
  show: false,
  type: '' as '' | 'x_post' | 'tmdb_card',
  xPost: {
    postId: '',
    author: '',
  },
  tmdb: {
    mediaType: 'movie',
    tmdbId: '',
    title: '',
  },
});

const tmdbTypeOptions = [
  { label: 'Movie', value: 'movie' },
  { label: 'TV', value: 'tv' },
];

const embedModalTitle = computed(() => {
  if (embedModal.type === 'x_post') {
    return 'Insert X Post';
  }
  if (embedModal.type === 'tmdb_card') {
    return 'Insert TMDB Card';
  }
  return 'Insert Embed';
});

function renderIcon(icon: any) {
  return () => h(NIcon, null, { default: () => h(icon) });
}

const blockMenuOptions: DropdownOption[] = [
  { label: 'Paragraph', key: 'paragraph', icon: renderIcon(TextOutline) },
  { label: 'Heading 2', key: 'h2', icon: renderIcon(TextOutline) },
  { label: 'Heading 3', key: 'h3', icon: renderIcon(TextOutline) },
  { type: 'divider', key: 'd1' },
  { label: 'Bullet List', key: 'bullet', icon: renderIcon(ListOutline) },
  { label: 'Ordered List', key: 'ordered', icon: renderIcon(ListOutline) },
  { label: 'Blockquote', key: 'quote', icon: renderIcon(LogoTwitter) },
  { label: 'Code Block', key: 'codeblock', icon: renderIcon(CodeSlashOutline) },
  { type: 'divider', key: 'd2' },
  { label: 'Divider', key: 'divider', icon: renderIcon(RemoveOutline) },
];

const embedMenuOptions: DropdownOption[] = [
  { label: 'X Post', key: 'x_post', icon: renderIcon(LogoTwitter) },
  { label: 'TMDB Card', key: 'tmdb_card', icon: renderIcon(FilmOutline) },
];

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

function onSelectBlockMenu(key: string): void {
  if (!editor.value) {
    return;
  }
  const chain = editor.value.chain().focus();
  switch (key) {
    case 'paragraph':
      chain.setParagraph().run();
      break;
    case 'h2':
      chain.toggleHeading({ level: 2 }).run();
      break;
    case 'h3':
      chain.toggleHeading({ level: 3 }).run();
      break;
    case 'bullet':
      chain.toggleBulletList().run();
      break;
    case 'ordered':
      chain.toggleOrderedList().run();
      break;
    case 'quote':
      chain.toggleBlockquote().run();
      break;
    case 'codeblock':
      chain.toggleCodeBlock().run();
      break;
    case 'divider':
      chain.setHorizontalRule().run();
      break;
    default:
      break;
  }
}

function onSelectEmbedMenu(key: string): void {
  if (key === 'x_post') {
    embedModal.type = 'x_post';
    embedModal.show = true;
    return;
  }
  if (key === 'tmdb_card') {
    embedModal.type = 'tmdb_card';
    embedModal.show = true;
  }
}

function confirmEmbedInsert(): void {
  if (!editor.value || !embedModal.type) {
    return;
  }

  if (embedModal.type === 'x_post') {
    const postID = embedModal.xPost.postId.trim();
    if (!postID) {
      emit('upload-error', 'Post ID is required.');
      return;
    }
    editor.value
      .chain()
      .focus()
      .insertContent({
        type: 'xPostEmbed',
        attrs: {
          postId: postID,
          author: embedModal.xPost.author.trim(),
        },
      })
      .run();
  }

  if (embedModal.type === 'tmdb_card') {
    const tmdbID = embedModal.tmdb.tmdbId.trim();
    if (!tmdbID) {
      emit('upload-error', 'TMDB ID is required.');
      return;
    }
    editor.value
      .chain()
      .focus()
      .insertContent({
        type: 'tmdbCardEmbed',
        attrs: {
          mediaType: embedModal.tmdb.mediaType === 'tv' ? 'tv' : 'movie',
          tmdbId: tmdbID,
          title: embedModal.tmdb.title.trim(),
        },
      })
      .run();
  }

  embedModal.show = false;
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
  border: 1px solid var(--n-border-color);
  border-radius: 12px;
  background: color-mix(in srgb, var(--n-card-color) 92%, var(--n-primary-color) 8%);
  position: sticky;
  top: 8px;
  z-index: 12;
  backdrop-filter: blur(6px);
  padding: 8px;
}

.toolbar-row {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  gap: 8px;
}

.toolbar-spacer {
  flex: 1;
}

.tool-btn,
.menu-btn {
  border-radius: 10px;
}

.menu-caret {
  margin-left: 4px;
  opacity: 0.7;
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

.clear-btn {
  color: var(--n-error-color);
}

.link-popover {
  width: 260px;
  display: grid;
  gap: 8px;
}

.link-actions,
.modal-actions {
  display: flex;
  justify-content: flex-end;
  gap: 8px;
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

  .toolbar-row {
    gap: 6px;
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
