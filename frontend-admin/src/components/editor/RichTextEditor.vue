<!--
File: RichTextEditor.vue
Purpose: Provide reusable rich-text editor with TipTap toolbar and image upload insertion support.
Module: frontend-admin/components/editor, editor UI layer.
Related: ArticleEditorView and admin upload API module.
-->
<template>
  <div class="editor-root">
    <div class="toolbar">
      <NButton size="small" tertiary :type="isActive('bold') ? 'primary' : 'default'" @click="toggleBold">B</NButton>
      <NButton size="small" tertiary :type="isActive('italic') ? 'primary' : 'default'" @click="toggleItalic">I</NButton>
      <NButton size="small" tertiary :type="isActive('strike') ? 'primary' : 'default'" @click="toggleStrike">S</NButton>
      <NButton size="small" tertiary :type="isActive('heading', { level: 2 }) ? 'primary' : 'default'" @click="toggleHeading">H2</NButton>
      <NButton size="small" tertiary :type="isActive('bulletList') ? 'primary' : 'default'" @click="toggleBulletList">• List</NButton>
      <NButton size="small" tertiary :type="isActive('orderedList') ? 'primary' : 'default'" @click="toggleOrderedList">1. List</NButton>
      <NButton size="small" tertiary :type="isActive('blockquote') ? 'primary' : 'default'" @click="toggleBlockquote">Quote</NButton>
      <NButton size="small" tertiary :type="isActive('codeBlock') ? 'primary' : 'default'" @click="toggleCodeBlock">Code</NButton>
      <NButton size="small" tertiary @click="setLink">Link</NButton>
      <NButton size="small" tertiary :loading="uploading" @click="pickImage">Image</NButton>
      <NButton size="small" tertiary @click="clearContent">Clear</NButton>
    </div>

    <div class="editor-shell">
      <EditorContent :editor="editor" />
    </div>

    <input ref="imageInputRef" class="hidden-input" type="file" accept="image/*" @change="onImageSelected" />
  </div>
</template>

<script setup lang="ts">
import { onBeforeUnmount, ref, watch } from 'vue';
import { EditorContent, useEditor } from '@tiptap/vue-3';
import type { JSONContent } from '@tiptap/core';
import StarterKit from '@tiptap/starter-kit';
import Link from '@tiptap/extension-link';
import Image from '@tiptap/extension-image';
import Placeholder from '@tiptap/extension-placeholder';
import { NButton } from 'naive-ui';

const props = withDefaults(defineProps<{
  modelValue: string;
  placeholder?: string;
  disabled?: boolean;
  uploadImage?: (file: File) => Promise<string>;
}>(), {
  placeholder: '',
  disabled: false,
  uploadImage: undefined,
});

const emit = defineEmits<{
  (e: 'update:modelValue', value: string): void;
  (e: 'upload-error', message: string): void;
}>();

const imageInputRef = ref<HTMLInputElement | null>(null);
const uploading = ref(false);

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
    Placeholder.configure({
      placeholder: props.placeholder,
      showOnlyWhenEditable: false,
    }),
  ],
  onUpdate() {
    emit('update:modelValue', serializeEditor());
  },
});

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

function toggleHeading(): void {
  editor.value?.chain().focus().toggleHeading({ level: 2 }).run();
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

function setLink(): void {
  const target = window.prompt('Enter URL');
  if (target === null) {
    return;
  }
  const trimmed = target.trim();
  if (!trimmed) {
    editor.value?.chain().focus().extendMarkRange('link').unsetLink().run();
    return;
  }
  editor.value?.chain().focus().extendMarkRange('link').setLink({ href: trimmed }).run();
}

function clearContent(): void {
  editor.value?.commands.setContent(emptyDoc());
}

function pickImage(): void {
  imageInputRef.value?.click();
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

.toolbar {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
}

.editor-shell {
  border: 1px solid var(--n-border-color);
  border-radius: 10px;
  background: var(--n-input-color);
  min-height: 320px;
  padding: 10px;
}

.editor-shell :deep(.ProseMirror) {
  min-height: 300px;
  outline: none;
  line-height: 1.7;
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
  border-radius: 8px;
}

.hidden-input {
  display: none;
}
</style>
