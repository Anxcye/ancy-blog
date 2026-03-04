<!--
File: MomentsView.vue
Purpose: Manage moment entries with list filtering, creation, and inline editing actions.
Module: frontend-admin/views/content, presentation layer.
Related: moments API module and backend admin moment endpoints.
-->
<template>
  <section class="moment-page">
    <nav class="subnav">
      <RouterLink :to="{ name: 'articles' }">{{ t('articles.tabArticles') }}</RouterLink>
      <RouterLink :to="{ name: 'moments' }">{{ t('articles.tabMoments') }}</RouterLink>
    </nav>

    <header class="header">
      <div>
        <h1>{{ t('moments.title') }}</h1>
        <p>{{ t('moments.subtitle') }}</p>
      </div>
      <button :disabled="loading" @click="openCreate">{{ t('moments.create') }}</button>
    </header>

    <div class="toolbar">
      <select v-model="statusFilter" @change="loadMoments(1)">
        <option value="">{{ t('moments.statusAll') }}</option>
        <option value="draft">draft</option>
        <option value="published">published</option>
        <option value="scheduled">scheduled</option>
      </select>
    </div>

    <p v-if="errorText" class="error">{{ errorText }}</p>
    <p v-if="successText" class="success">{{ successText }}</p>

    <ul class="list">
      <li v-for="item in rows" :key="item.id" class="item">
        <div>
          <p class="content">{{ item.content }}</p>
          <p class="meta">{{ item.status }} · {{ formatDate(item.updatedAt) }}</p>
        </div>
        <button :disabled="loading" @click="openEdit(item)">{{ t('moments.edit') }}</button>
      </li>
    </ul>

    <div v-if="showForm" class="modal-mask" @click.self="closeForm">
      <article class="modal">
        <h2>{{ editingId ? t('moments.edit') : t('moments.create') }}</h2>
        <label>
          <span>{{ t('moments.fieldContent') }}</span>
          <textarea v-model="form.content" rows="7" />
        </label>
        <div class="grid-2">
          <label>
            <span>{{ t('moments.fieldStatus') }}</span>
            <select v-model="form.status">
              <option value="draft">draft</option>
              <option value="published">published</option>
              <option value="scheduled">scheduled</option>
            </select>
          </label>
          <label>
            <span>{{ t('moments.fieldPublishedAt') }}</span>
            <input v-model="publishedAtLocal" type="datetime-local" />
          </label>
        </div>
        <label class="switch-row">
          <input v-model="form.allowComment" type="checkbox" />
          <span>{{ t('moments.fieldAllowComment') }}</span>
        </label>
        <div class="actions">
          <button :disabled="loading" @click="closeForm">{{ t('common.cancel') }}</button>
          <button :disabled="loading" @click="saveMoment">{{ t('common.save') }}</button>
        </div>
      </article>
    </div>
  </section>
</template>

<script setup lang="ts">
import { onMounted, reactive, ref, watch } from 'vue';
import { useI18n } from 'vue-i18n';

import { createMoment, listMoments, updateMoment, type Moment } from '@/api/modules/moments';

const { t } = useI18n();

const loading = ref(false);
const errorText = ref('');
const successText = ref('');

const rows = ref<Moment[]>([]);
const statusFilter = ref('');

const showForm = ref(false);
const editingId = ref('');
const publishedAtLocal = ref('');
const form = reactive({
  content: '',
  status: 'draft',
  allowComment: true,
  publishedAt: '',
});

watch(publishedAtLocal, (value) => {
  form.publishedAt = value ? new Date(value).toISOString() : '';
});

function formatDate(value: string): string {
  if (!value) {
    return '-';
  }
  return new Date(value).toLocaleString();
}

function toDateLocal(value: string): string {
  if (!value) {
    return '';
  }
  const date = new Date(value);
  const pad = (v: number) => String(v).padStart(2, '0');
  return `${date.getFullYear()}-${pad(date.getMonth() + 1)}-${pad(date.getDate())}T${pad(date.getHours())}:${pad(date.getMinutes())}`;
}

async function loadMoments(page: number): Promise<void> {
  loading.value = true;
  errorText.value = '';
  try {
    const result = await listMoments({ page, pageSize: 30, status: statusFilter.value || undefined });
    rows.value = result.rows;
  } catch {
    errorText.value = t('common.loadFailed');
  } finally {
    loading.value = false;
  }
}

function openCreate(): void {
  editingId.value = '';
  form.content = '';
  form.status = 'draft';
  form.allowComment = true;
  form.publishedAt = '';
  publishedAtLocal.value = '';
  showForm.value = true;
}

function openEdit(item: Moment): void {
  editingId.value = item.id;
  form.content = item.content;
  form.status = item.status;
  form.allowComment = item.allowComment;
  form.publishedAt = item.publishedAt || '';
  publishedAtLocal.value = toDateLocal(form.publishedAt);
  showForm.value = true;
}

function closeForm(): void {
  showForm.value = false;
}

async function saveMoment(): Promise<void> {
  if (!form.content.trim()) {
    errorText.value = t('moments.contentRequired');
    return;
  }
  loading.value = true;
  errorText.value = '';
  successText.value = '';
  try {
    const payload = {
      content: form.content,
      status: form.status,
      allowComment: form.allowComment,
      publishedAt: form.publishedAt || undefined,
    };
    if (editingId.value) {
      await updateMoment(editingId.value, payload);
    } else {
      await createMoment(payload);
    }
    successText.value = t('common.saveSuccess');
    closeForm();
    await loadMoments(1);
  } catch {
    errorText.value = t('common.saveFailed');
    loading.value = false;
  }
}

onMounted(async () => {
  await loadMoments(1);
});
</script>

<style scoped>
.moment-page {
  display: grid;
  gap: 12px;
}

.subnav {
  display: flex;
  gap: 8px;
}

.subnav a {
  text-decoration: none;
  border: 1px solid var(--border);
  border-radius: 999px;
  padding: 6px 10px;
  color: var(--muted);
}

.subnav a.router-link-active {
  border-color: var(--accent);
  color: var(--accent-hover);
  background: var(--accent-soft);
}

.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 10px;
}

h1 {
  margin: 0;
}

p {
  margin: 4px 0 0;
  color: var(--muted);
}

.toolbar {
  display: flex;
  justify-content: flex-end;
}

select,
button,
input,
textarea {
  border: 1px solid var(--border);
  border-radius: 8px;
  padding: 8px 10px;
  font: inherit;
  background: var(--surface);
}

button {
  cursor: pointer;
}

.list {
  list-style: none;
  margin: 0;
  padding: 0;
  display: grid;
  gap: 8px;
}

.item {
  border: 1px solid var(--border);
  border-radius: 8px;
  padding: 10px;
  display: flex;
  justify-content: space-between;
  gap: 10px;
}

.content,
.meta {
  margin: 0;
}

.content {
  color: var(--text);
  white-space: pre-wrap;
}

.meta {
  margin-top: 6px;
  font-size: 13px;
}

.modal-mask {
  position: fixed;
  inset: 0;
  background: rgba(26, 26, 26, 0.3);
  display: grid;
  place-items: center;
  padding: 12px;
  z-index: 30;
}

.modal {
  width: min(740px, 100%);
  border: 1px solid var(--border);
  border-radius: 12px;
  background: var(--surface);
  padding: 14px;
  display: grid;
  gap: 10px;
}

.grid-2 {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 10px;
}

label {
  display: grid;
  gap: 6px;
}

.switch-row {
  display: inline-flex;
  align-items: center;
  gap: 8px;
}

.actions {
  display: flex;
  justify-content: flex-end;
  gap: 8px;
}

.error {
  color: #b64040;
  margin: 0;
}

.success {
  color: var(--accent-hover);
  margin: 0;
}

@media (max-width: 900px) {
  .header,
  .item {
    flex-direction: column;
    align-items: stretch;
  }

  .grid-2 {
    grid-template-columns: 1fr;
  }

  .actions button,
  .header button {
    width: 100%;
  }
}
</style>
