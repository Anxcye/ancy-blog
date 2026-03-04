<!--
File: ArticleEditorView.vue
Purpose: Provide article create/edit form with content, metadata, and preview modes.
Module: frontend-admin/views/content, editor presentation layer.
Related: article APIs, route params, translation and publish workflows.
-->
<template>
  <section class="editor-page">
    <NCard :bordered="false" class="section-card">
      <header class="editor-header">
        <h1>{{ isEdit ? t('editor.editTitle') : t('editor.createTitle') }}</h1>
        <div class="header-actions">
          <NButton :loading="saving" @click="saveAs('draft')">{{ t('editor.saveDraft') }}</NButton>
          <NButton type="primary" :loading="saving" @click="saveAs('published')">{{ t('editor.publish') }}</NButton>
        </div>
      </header>

      <NAlert v-if="errorText" type="error" :show-icon="false">{{ errorText }}</NAlert>
      <NAlert v-if="successText" type="success" :show-icon="false">{{ successText }}</NAlert>

      <div class="mode-tabs">
        <NButton :type="mode === 'content' ? 'primary' : 'default'" secondary @click="mode = 'content'">{{ t('editor.tabContent') }}</NButton>
        <NButton :type="mode === 'meta' ? 'primary' : 'default'" secondary @click="mode = 'meta'">{{ t('editor.tabMeta') }}</NButton>
        <NButton :type="mode === 'preview' ? 'primary' : 'default'" secondary @click="mode = 'preview'">{{ t('editor.tabPreview') }}</NButton>
      </div>

      <NForm v-if="mode === 'content'" label-placement="top" class="panel">
        <NFormItem :label="t('editor.fieldTitle')" required>
          <NInput v-model:value="form.title" />
        </NFormItem>

        <NFormItem :label="t('editor.fieldSlug')" required>
          <div class="inline-row">
            <NInput v-model:value="form.slug" />
            <NButton :loading="aiLoading" :disabled="saving || aiLoading" @click="onSuggestSlug">{{ t('editor.aiSuggestSlug') }}</NButton>
          </div>
        </NFormItem>

        <NFormItem :label="t('editor.fieldSummary')">
          <div class="inline-col">
            <NInput v-model:value="form.summary" type="textarea" :autosize="{ minRows: 3, maxRows: 8 }" />
            <NButton :loading="aiLoading" :disabled="saving || aiLoading" @click="onGenerateSummary">{{ t('editor.aiGenerateSummary') }}</NButton>
          </div>
        </NFormItem>

        <NFormItem :label="t('editor.fieldContent')">
          <RichTextEditor
            v-model="form.content"
            :disabled="saving || aiLoading"
            placeholder="Start writing your article..."
            :upload-image="uploadEditorImage"
            @upload-error="onEditorUploadError"
          />
        </NFormItem>
      </NForm>

      <NForm v-else-if="mode === 'meta'" label-placement="top" class="panel">
        <div class="grid-2">
          <NFormItem :label="t('editor.fieldKind')">
            <NSelect v-model:value="form.contentKind" :options="kindOptions" />
          </NFormItem>
          <NFormItem :label="t('editor.fieldStatus')">
            <NSelect v-model:value="form.status" :options="statusOptions" />
          </NFormItem>
          <NFormItem :label="t('editor.fieldVisibility')">
            <NSelect v-model:value="form.visibility" :options="visibilityOptions" />
          </NFormItem>
          <NFormItem :label="t('editor.fieldOriginType')">
            <NSelect v-model:value="form.originType" :options="originOptions" />
          </NFormItem>
          <NFormItem :label="t('editor.fieldAIAssistLevel')">
            <NInput v-model:value="form.aiAssistLevel" placeholder="none" />
          </NFormItem>
          <NFormItem :label="t('editor.fieldAIProviderKey')">
            <NInput v-model:value="aiOptions.providerKey" placeholder="openai_compatible" />
          </NFormItem>
          <NFormItem :label="t('editor.fieldAIModelName')">
            <NInput v-model:value="aiOptions.modelName" placeholder="gpt-4.1-mini" />
          </NFormItem>
          <NFormItem :label="t('editor.fieldCategory')">
            <NInput v-model:value="form.categorySlug" placeholder="tech" />
          </NFormItem>
          <NFormItem :label="t('editor.fieldTags')">
            <NInput v-model:value="tagInput" placeholder="go,backend" />
          </NFormItem>
          <NFormItem :label="t('editor.fieldSourceUrl')">
            <NInput v-model:value="form.sourceUrl" />
          </NFormItem>
          <NFormItem :label="t('editor.fieldCoverImage')">
            <NInput v-model:value="form.coverImage" />
          </NFormItem>
          <NFormItem :label="t('editor.fieldPublishedAt')">
            <NDatePicker v-model:value="publishedAtTs" type="datetime" clearable style="width: 100%" />
          </NFormItem>
          <NFormItem>
            <NCheckbox v-model:checked="form.allowComment">{{ t('editor.fieldAllowComment') }}</NCheckbox>
          </NFormItem>
        </div>
      </NForm>

      <NCard v-else size="small" :bordered="true" class="preview">
        <h2>{{ form.title || '-' }}</h2>
        <p class="preview-meta">/{{ form.slug || '-' }}</p>
        <p class="preview-summary">{{ form.summary || '-' }}</p>
        <RichTextPreview :content="form.content" />
      </NCard>
    </NCard>
  </section>
</template>

<script setup lang="ts">
import { computed, onBeforeUnmount, onMounted, reactive, ref, watch } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { useI18n } from 'vue-i18n';
import { NAlert, NButton, NCard, NCheckbox, NDatePicker, NForm, NFormItem, NInput, NSelect } from 'naive-ui';

import RichTextEditor from '@/components/editor/RichTextEditor.vue';
import RichTextPreview from '@/components/editor/RichTextPreview.vue';
import { createArticle, getArticle, updateArticle } from '@/api/modules/articles';
import { generateSummary, suggestSlug } from '@/api/modules/ai';
import { uploadImage } from '@/api/modules/upload';
import type { ArticleUpsertPayload } from '@/api/types';

type EditorMode = 'content' | 'meta' | 'preview';

type SaveStatus = 'draft' | 'published';

const route = useRoute();
const router = useRouter();
const { t } = useI18n();

const mode = ref<EditorMode>('content');
const saving = ref(false);
const aiLoading = ref(false);
const errorText = ref('');
const successText = ref('');
const publishedAtTs = ref<number | null>(null);
const tagInput = ref('');
let autosaveTimer: ReturnType<typeof setTimeout> | null = null;
const aiOptions = reactive({
  providerKey: 'openai_compatible',
  modelName: '',
});

const form = reactive<ArticleUpsertPayload>({
  title: '',
  slug: '',
  contentKind: 'post',
  summary: '',
  content: '',
  status: 'draft',
  visibility: 'public',
  allowComment: true,
  originType: 'original',
  sourceUrl: '',
  aiAssistLevel: 'none',
  coverImage: '',
  categorySlug: '',
  tagSlugs: [],
  publishedAt: '',
});

const kindOptions = [
  { label: 'post', value: 'post' },
  { label: 'page', value: 'page' },
];

const statusOptions = [
  { label: 'draft', value: 'draft' },
  { label: 'published', value: 'published' },
  { label: 'scheduled', value: 'scheduled' },
];

const visibilityOptions = [
  { label: 'public', value: 'public' },
  { label: 'unlisted', value: 'unlisted' },
  { label: 'private', value: 'private' },
];

const originOptions = [
  { label: 'original', value: 'original' },
  { label: 'repost', value: 'repost' },
  { label: 'translated', value: 'translated' },
];

const isEdit = computed(() => Boolean(route.params.id));
const draftStorageKey = computed(() => {
  const id = String(route.params.id || 'new');
  return `ancy_admin_editor_draft:${id}`;
});

watch(tagInput, (value) => {
  form.tagSlugs = value
    .split(',')
    .map((item) => item.trim())
    .filter((item) => item.length > 0);
});

watch(publishedAtTs, (value) => {
  form.publishedAt = value ? new Date(value).toISOString() : '';
});

watch(
  () => ({
    title: form.title,
    slug: form.slug,
    contentKind: form.contentKind,
    summary: form.summary,
    content: form.content,
    status: form.status,
    visibility: form.visibility,
    allowComment: form.allowComment,
    originType: form.originType,
    sourceUrl: form.sourceUrl,
    aiAssistLevel: form.aiAssistLevel,
    coverImage: form.coverImage,
    categorySlug: form.categorySlug,
    tagInput: tagInput.value,
    publishedAt: form.publishedAt,
    aiProviderKey: aiOptions.providerKey,
    aiModelName: aiOptions.modelName,
  }),
  () => {
    if (autosaveTimer) {
      clearTimeout(autosaveTimer);
    }
    autosaveTimer = setTimeout(() => {
      const payload = {
        form: { ...form },
        tagInput: tagInput.value,
        aiOptions: { ...aiOptions },
      };
      localStorage.setItem(draftStorageKey.value, JSON.stringify(payload));
    }, 500);
  },
  { deep: true },
);

async function loadArticle(): Promise<void> {
  if (!isEdit.value) {
    return;
  }
  const id = String(route.params.id);
  try {
    const article = await getArticle(id);
    form.title = article.title;
    form.slug = article.slug;
    form.contentKind = article.contentKind;
    form.summary = article.summary;
    form.content = article.content;
    form.status = article.status;
    form.visibility = article.visibility;
    form.allowComment = article.allowComment;
    form.originType = article.originType;
    form.sourceUrl = article.sourceUrl || '';
    form.aiAssistLevel = article.aiAssistLevel || 'none';
    form.coverImage = article.coverImage || '';
    form.categorySlug = article.categorySlug || '';
    form.tagSlugs = article.tagSlugs || [];
    form.publishedAt = article.publishedAt || '';

    tagInput.value = form.tagSlugs.join(',');
    publishedAtTs.value = form.publishedAt ? new Date(form.publishedAt).getTime() : null;
  } catch {
    errorText.value = t('common.loadFailed');
  }
}

function restoreDraftIfNeeded(): void {
  if (isEdit.value) {
    return;
  }
  const raw = localStorage.getItem(draftStorageKey.value);
  if (!raw) {
    return;
  }
  try {
    const parsed = JSON.parse(raw) as {
      form?: Partial<ArticleUpsertPayload>;
      tagInput?: string;
      aiOptions?: Partial<typeof aiOptions>;
    };
    if (parsed.form) {
      Object.assign(form, parsed.form);
    }
    if (typeof parsed.tagInput === 'string') {
      tagInput.value = parsed.tagInput;
    }
    if (parsed.aiOptions) {
      Object.assign(aiOptions, parsed.aiOptions);
    }
    publishedAtTs.value = form.publishedAt ? new Date(form.publishedAt).getTime() : null;
  } catch {
    localStorage.removeItem(draftStorageKey.value);
  }
}

function toPayload(status: SaveStatus): ArticleUpsertPayload {
  const publishedAt = form.publishedAt && form.publishedAt.trim() !== '' ? form.publishedAt : undefined;
  return {
    ...form,
    status,
    publishedAt,
    aiAssistLevel: form.aiAssistLevel || 'none',
  };
}

async function saveAs(status: SaveStatus): Promise<void> {
  if (!form.title || !form.slug) {
    errorText.value = t('editor.requiredHint');
    return;
  }
  saving.value = true;
  errorText.value = '';
  successText.value = '';

  try {
    const payload = toPayload(status);
    if (isEdit.value) {
      const id = String(route.params.id);
      await updateArticle(id, payload);
    } else {
      const id = await createArticle(payload);
      localStorage.removeItem(draftStorageKey.value);
      await router.replace({ name: 'article-edit', params: { id } });
    }
    if (isEdit.value) {
      localStorage.removeItem(draftStorageKey.value);
    }
    successText.value = t('editor.saveSuccess');
  } catch {
    errorText.value = t('editor.saveFailed');
  } finally {
    saving.value = false;
  }
}

async function onGenerateSummary(): Promise<void> {
  const textContent = extractTextFromContent(form.content);
  if (!form.title && !textContent) {
    errorText.value = t('editor.requiredHint');
    return;
  }
  aiLoading.value = true;
  errorText.value = '';
  successText.value = '';
  try {
    const result = await generateSummary({
      title: form.title,
      content: textContent,
      providerKey: aiOptions.providerKey || undefined,
      modelName: aiOptions.modelName || undefined,
      maxLength: 180,
    });
    form.summary = result.summary;
    successText.value = t('editor.aiSummaryDone');
  } catch {
    errorText.value = t('editor.aiSummaryFailed');
  } finally {
    aiLoading.value = false;
  }
}

function extractTextFromContent(raw: string): string {
  if (!raw || !raw.trim()) {
    return '';
  }
  try {
    const parsed = JSON.parse(raw) as Record<string, unknown>;
    const parts: string[] = [];
    const walk = (node: unknown): void => {
      if (!node || typeof node !== 'object') {
        return;
      }
      const map = node as Record<string, unknown>;
      if (typeof map.text === 'string') {
        parts.push(map.text);
      }
      if (Array.isArray(map.content)) {
        map.content.forEach(walk);
      }
    };
    walk(parsed);
    return parts.join(' ').replace(/\s+/g, ' ').trim();
  } catch {
    return raw;
  }
}

async function uploadEditorImage(file: File): Promise<string> {
  const result = await uploadImage(file);
  successText.value = 'Image uploaded successfully.';
  return result.url;
}

function onEditorUploadError(message: string): void {
  errorText.value = message || t('common.saveFailed');
}

async function onSuggestSlug(): Promise<void> {
  if (!form.title) {
    errorText.value = t('editor.requiredHint');
    return;
  }
  aiLoading.value = true;
  errorText.value = '';
  successText.value = '';
  try {
    const result = await suggestSlug({
      title: form.title,
      providerKey: aiOptions.providerKey || undefined,
      modelName: aiOptions.modelName || undefined,
    });
    form.slug = result.slug;
    successText.value = t('editor.aiSlugDone');
  } catch {
    errorText.value = t('editor.aiSlugFailed');
  } finally {
    aiLoading.value = false;
  }
}

onMounted(async () => {
  restoreDraftIfNeeded();
  await loadArticle();
});

onBeforeUnmount(() => {
  if (autosaveTimer) {
    clearTimeout(autosaveTimer);
    autosaveTimer = null;
  }
});
</script>

<style scoped>
.editor-page {
  display: grid;
}

.section-card {
  border-radius: 14px;
  box-shadow: 0 6px 24px color-mix(in srgb, var(--n-text-color) 8%, transparent);
  display: grid;
  gap: 12px;
}

.editor-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 10px;
}

h1 {
  margin: 0;
}

.header-actions {
  display: flex;
  gap: 8px;
}

.mode-tabs {
  display: flex;
  gap: 8px;
}

.panel {
  display: grid;
  gap: 12px;
}

.inline-row {
  display: flex;
  gap: 8px;
}

.inline-row :deep(.n-input) {
  flex: 1;
}

.inline-col {
  display: grid;
  gap: 8px;
}

.grid-2 {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 10px;
}

.preview h2 {
  margin: 0;
}

.preview-meta {
  margin: 4px 0;
  color: var(--n-text-color-3);
  font-family: ui-monospace, SFMono-Regular, Menlo, monospace;
}

.preview-summary {
  margin: 6px 0 10px;
  color: var(--n-text-color-3);
}

@media (max-width: 900px) {
  .editor-header {
    align-items: stretch;
    flex-direction: column;
  }

  .header-actions {
    width: 100%;
  }

  .header-actions :deep(button) {
    flex: 1;
  }

  .inline-row,
  .grid-2 {
    grid-template-columns: 1fr;
    flex-direction: column;
  }
}
</style>
