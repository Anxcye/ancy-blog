<!--
File: ArticleEditorView.vue
Purpose: Provide article create/edit form with content, metadata, and preview modes.
Module: frontend-admin/views/content, editor presentation layer.
Related: article APIs, route params, translation and publish workflows.
-->
<template>
  <section class="editor-page">
    <header class="editor-header">
      <h1>{{ isEdit ? t('editor.editTitle') : t('editor.createTitle') }}</h1>
      <div class="header-actions">
        <button :disabled="saving" @click="saveAs('draft')">{{ t('editor.saveDraft') }}</button>
        <button class="primary" :disabled="saving" @click="saveAs('published')">{{ t('editor.publish') }}</button>
      </div>
    </header>

    <p v-if="errorText" class="error">{{ errorText }}</p>
    <p v-if="successText" class="success">{{ successText }}</p>

    <div class="mode-tabs">
      <button :class="{ active: mode === 'content' }" @click="mode = 'content'">{{ t('editor.tabContent') }}</button>
      <button :class="{ active: mode === 'meta' }" @click="mode = 'meta'">{{ t('editor.tabMeta') }}</button>
      <button :class="{ active: mode === 'preview' }" @click="mode = 'preview'">{{ t('editor.tabPreview') }}</button>
    </div>

    <div v-if="mode === 'content'" class="panel">
      <label>
        <span>{{ t('editor.fieldTitle') }}</span>
        <input v-model.trim="form.title" type="text" required />
      </label>
      <label>
        <span>{{ t('editor.fieldSlug') }}</span>
        <div class="inline-row">
          <input v-model.trim="form.slug" type="text" required />
          <button type="button" :disabled="saving || aiLoading" @click="onSuggestSlug">
            {{ t('editor.aiSuggestSlug') }}
          </button>
        </div>
      </label>
      <label>
        <span>{{ t('editor.fieldSummary') }}</span>
        <div class="inline-col">
          <textarea v-model="form.summary" rows="3" />
          <button type="button" :disabled="saving || aiLoading" @click="onGenerateSummary">
            {{ t('editor.aiGenerateSummary') }}
          </button>
        </div>
      </label>
      <label>
        <span>{{ t('editor.fieldContent') }}</span>
        <textarea v-model="form.content" rows="18" />
      </label>
    </div>

    <div v-else-if="mode === 'meta'" class="panel grid-2">
      <label>
        <span>{{ t('editor.fieldKind') }}</span>
        <select v-model="form.contentKind">
          <option value="post">post</option>
          <option value="page">page</option>
        </select>
      </label>
      <label>
        <span>{{ t('editor.fieldStatus') }}</span>
        <select v-model="form.status">
          <option value="draft">draft</option>
          <option value="published">published</option>
          <option value="scheduled">scheduled</option>
        </select>
      </label>
      <label>
        <span>{{ t('editor.fieldVisibility') }}</span>
        <select v-model="form.visibility">
          <option value="public">public</option>
          <option value="unlisted">unlisted</option>
          <option value="private">private</option>
        </select>
      </label>
      <label>
        <span>{{ t('editor.fieldOriginType') }}</span>
        <select v-model="form.originType">
          <option value="original">original</option>
          <option value="repost">repost</option>
          <option value="translated">translated</option>
        </select>
      </label>
      <label>
        <span>{{ t('editor.fieldAIAssistLevel') }}</span>
        <input v-model.trim="form.aiAssistLevel" type="text" placeholder="none" />
      </label>
      <label>
        <span>{{ t('editor.fieldAIProviderKey') }}</span>
        <input v-model.trim="aiOptions.providerKey" type="text" placeholder="openai_compatible" />
      </label>
      <label>
        <span>{{ t('editor.fieldAIModelName') }}</span>
        <input v-model.trim="aiOptions.modelName" type="text" placeholder="gpt-4.1-mini" />
      </label>
      <label>
        <span>{{ t('editor.fieldCategory') }}</span>
        <input v-model.trim="form.categorySlug" type="text" placeholder="tech" />
      </label>
      <label>
        <span>{{ t('editor.fieldTags') }}</span>
        <input v-model.trim="tagInput" type="text" placeholder="go,backend" />
      </label>
      <label>
        <span>{{ t('editor.fieldSourceUrl') }}</span>
        <input v-model.trim="form.sourceUrl" type="url" />
      </label>
      <label>
        <span>{{ t('editor.fieldCoverImage') }}</span>
        <input v-model.trim="form.coverImage" type="url" />
      </label>
      <label>
        <span>{{ t('editor.fieldPublishedAt') }}</span>
        <input v-model="publishedAtLocal" type="datetime-local" />
      </label>
      <label class="switch-row">
        <input v-model="form.allowComment" type="checkbox" />
        <span>{{ t('editor.fieldAllowComment') }}</span>
      </label>
    </div>

    <div v-else class="panel preview">
      <h2>{{ form.title || '-' }}</h2>
      <p class="preview-meta">/{{ form.slug || '-' }}</p>
      <p class="preview-summary">{{ form.summary || '-' }}</p>
      <pre>{{ form.content || '-' }}</pre>
    </div>
  </section>
</template>

<script setup lang="ts">
import { computed, onMounted, reactive, ref, watch } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { useI18n } from 'vue-i18n';

import { createArticle, getArticle, updateArticle } from '@/api/modules/articles';
import { generateSummary, suggestSlug } from '@/api/modules/ai';
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
const publishedAtLocal = ref('');
const tagInput = ref('');
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

const isEdit = computed(() => Boolean(route.params.id));

watch(tagInput, (value) => {
  form.tagSlugs = value
    .split(',')
    .map((item) => item.trim())
    .filter((item) => item.length > 0);
});

watch(publishedAtLocal, (value) => {
  form.publishedAt = value ? new Date(value).toISOString() : '';
});

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
    publishedAtLocal.value = form.publishedAt ? formatToDateTimeLocal(form.publishedAt) : '';
  } catch {
    errorText.value = t('common.loadFailed');
  }
}

function formatToDateTimeLocal(value: string): string {
  const date = new Date(value);
  const pad = (part: number) => String(part).padStart(2, '0');
  return `${date.getFullYear()}-${pad(date.getMonth() + 1)}-${pad(date.getDate())}T${pad(date.getHours())}:${pad(date.getMinutes())}`;
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
      await router.replace({ name: 'article-edit', params: { id } });
    }
    successText.value = t('editor.saveSuccess');
  } catch {
    errorText.value = t('editor.saveFailed');
  } finally {
    saving.value = false;
  }
}

async function onGenerateSummary(): Promise<void> {
  if (!form.title && !form.content) {
    errorText.value = t('editor.requiredHint');
    return;
  }
  aiLoading.value = true;
  errorText.value = '';
  successText.value = '';
  try {
    const result = await generateSummary({
      title: form.title,
      content: form.content,
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
  await loadArticle();
});
</script>

<style scoped>
.editor-page {
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

button {
  border: 1px solid var(--border);
  border-radius: 8px;
  padding: 8px 12px;
  background: var(--surface);
  cursor: pointer;
}

.inline-row {
  display: flex;
  gap: 8px;
}

.inline-row input {
  flex: 1;
}

.inline-col {
  display: grid;
  gap: 8px;
}

button.primary {
  background: var(--accent);
  color: #fff;
}

button:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.error {
  margin: 0;
  color: #b64040;
}

.success {
  margin: 0;
  color: var(--accent-hover);
}

.mode-tabs {
  display: inline-flex;
  border: 1px solid var(--border);
  border-radius: 10px;
  overflow: hidden;
}

.mode-tabs button {
  border: 0;
  border-right: 1px solid var(--border);
  border-radius: 0;
}

.mode-tabs button:last-child {
  border-right: 0;
}

.mode-tabs button.active {
  background: var(--accent-soft);
  color: var(--accent-hover);
  font-weight: 600;
}

.panel {
  padding: 16px;
  border: 1px solid var(--border);
  border-radius: 12px;
  background: var(--surface);
  display: grid;
  gap: 12px;
}

.grid-2 {
  grid-template-columns: repeat(2, minmax(0, 1fr));
}

label {
  display: grid;
  gap: 6px;
}

input,
select,
textarea {
  border: 1px solid var(--border);
  border-radius: 8px;
  padding: 8px 10px;
  font: inherit;
  width: 100%;
}

textarea {
  resize: vertical;
}

.switch-row {
  display: flex;
  align-items: center;
  gap: 8px;
}

.switch-row input {
  width: auto;
}

.preview h2 {
  margin: 0;
}

.preview-meta {
  margin: 4px 0;
  color: var(--muted);
  font-family: ui-monospace, SFMono-Regular, Menlo, monospace;
}

.preview-summary {
  margin: 6px 0 10px;
  color: var(--muted);
}

pre {
  white-space: pre-wrap;
  margin: 0;
  font: inherit;
  line-height: 1.7;
}

@media (max-width: 900px) {
  .editor-header {
    align-items: stretch;
    flex-direction: column;
  }

  .header-actions {
    width: 100%;
  }

  .header-actions button {
    flex: 1;
  }

  .grid-2 {
    grid-template-columns: 1fr;
  }
}
</style>
