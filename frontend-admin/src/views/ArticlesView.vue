<!--
File: ArticlesView.vue
Purpose: Render admin article list with filters and entry points to editor actions.
Module: frontend-admin/views/content, presentation layer.
Related: articles API module, ArticleEditorView route, dashboard workflow.
-->
<template>
  <section class="panel">
    <nav class="subnav">
      <RouterLink :to="{ name: 'articles' }">{{ t('articles.tabArticles') }}</RouterLink>
      <RouterLink :to="{ name: 'moments' }">{{ t('articles.tabMoments') }}</RouterLink>
    </nav>

    <header class="panel-header">
      <div>
        <h1>{{ t('articles.title') }}</h1>
        <p>{{ t('articles.subtitle') }}</p>
      </div>
      <RouterLink class="create-btn" :to="{ name: 'article-new' }">
        {{ t('articles.create') }}
      </RouterLink>
    </header>

    <div class="bulk-actions">
      <label class="check-all">
        <input :checked="allChecked" type="checkbox" @change="toggleAll(($event.target as HTMLInputElement).checked)" />
        <span>{{ t('articles.selected', { count: selectedIds.length }) }}</span>
      </label>
      <div class="bulk-buttons">
        <button :disabled="selectedIds.length === 0 || loading" @click="applyStatus('draft')">{{ t('articles.toDraft') }}</button>
        <button :disabled="selectedIds.length === 0 || loading" @click="applyStatus('published')">{{ t('articles.toPublished') }}</button>
        <button :disabled="selectedIds.length === 0 || loading" class="danger" @click="removeSelected">{{ t('common.delete') }}</button>
      </div>
    </div>

    <form class="filters" @submit.prevent="reload(1)">
      <input
        v-model.trim="filters.keyword"
        :placeholder="t('articles.filterKeyword')"
        type="text"
      />
      <select v-model="filters.status">
        <option value="">{{ t('articles.filterStatusAll') }}</option>
        <option value="draft">{{ t('articles.statusDraft') }}</option>
        <option value="published">{{ t('articles.statusPublished') }}</option>
        <option value="scheduled">{{ t('articles.statusScheduled') }}</option>
      </select>
      <select v-model="filters.contentKind">
        <option value="">{{ t('articles.filterKindAll') }}</option>
        <option value="post">Post</option>
        <option value="page">Page</option>
      </select>
      <button type="submit">{{ t('common.search') }}</button>
    </form>

    <p v-if="errorText" class="error">{{ errorText }}</p>

    <div class="table-wrap">
      <table>
        <thead>
          <tr>
            <th></th>
            <th>{{ t('articles.colTitle') }}</th>
            <th>{{ t('articles.colSlug') }}</th>
            <th>{{ t('articles.colStatus') }}</th>
            <th>{{ t('articles.colKind') }}</th>
            <th>{{ t('articles.colUpdatedAt') }}</th>
            <th>{{ t('articles.colAction') }}</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="item in rows" :key="item.id">
            <td>
              <input :checked="selectedIds.includes(item.id)" type="checkbox" @change="toggleOne(item.id, ($event.target as HTMLInputElement).checked)" />
            </td>
            <td>{{ item.title }}</td>
            <td class="mono">{{ item.slug }}</td>
            <td>{{ item.status }}</td>
            <td>{{ item.contentKind }}</td>
            <td>{{ formatDate(item.updatedAt) }}</td>
            <td>
              <RouterLink :to="{ name: 'article-edit', params: { id: item.id } }">
                {{ t('articles.edit') }}
              </RouterLink>
              <button class="link-btn danger" :disabled="loading" @click="removeOne(item.id)">{{ t('common.delete') }}</button>
            </td>
          </tr>
          <tr v-if="!loading && rows.length === 0">
            <td colspan="7">{{ t('articles.empty') }}</td>
          </tr>
        </tbody>
      </table>
    </div>

    <footer class="pager">
      <span>{{ t('articles.total', { total }) }}</span>
      <div class="pager-actions">
        <button :disabled="page <= 1 || loading" @click="reload(page - 1)">{{ t('common.prev') }}</button>
        <span>{{ page }}</span>
        <button :disabled="page >= totalPages || loading" @click="reload(page + 1)">{{ t('common.next') }}</button>
      </div>
    </footer>
  </section>
</template>

<script setup lang="ts">
import { computed, onMounted, reactive, ref } from 'vue';
import { useI18n } from 'vue-i18n';

import { batchDeleteArticles, batchUpdateArticleStatus, deleteArticle, listArticles } from '@/api/modules/articles';
import type { Article } from '@/api/types';

const { t } = useI18n();

const loading = ref(false);
const errorText = ref('');
const rows = ref<Article[]>([]);
const selectedIds = ref<string[]>([]);
const total = ref(0);
const page = ref(1);
const pageSize = 12;

const filters = reactive({
  keyword: '',
  status: '',
  contentKind: '',
});

const totalPages = computed(() => Math.max(1, Math.ceil(total.value / pageSize)));
const allChecked = computed(() => rows.value.length > 0 && rows.value.every((row) => selectedIds.value.includes(row.id)));

function formatDate(value: string): string {
  if (!value) {
    return '-';
  }
  return new Date(value).toLocaleString();
}

async function reload(nextPage = page.value): Promise<void> {
  loading.value = true;
  errorText.value = '';
  try {
    const result = await listArticles({
      page: nextPage,
      pageSize,
      status: filters.status || undefined,
      contentKind: filters.contentKind || undefined,
      keyword: filters.keyword || undefined,
    });
    rows.value = result.rows;
    selectedIds.value = [];
    total.value = result.total;
    page.value = nextPage;
  } catch {
    errorText.value = t('common.loadFailed');
  } finally {
    loading.value = false;
  }
}

function toggleOne(id: string, checked: boolean): void {
  if (checked) {
    if (!selectedIds.value.includes(id)) {
      selectedIds.value = [...selectedIds.value, id];
    }
    return;
  }
  selectedIds.value = selectedIds.value.filter((item) => item !== id);
}

function toggleAll(checked: boolean): void {
  if (checked) {
    selectedIds.value = rows.value.map((item) => item.id);
    return;
  }
  selectedIds.value = [];
}

async function applyStatus(status: 'draft' | 'published'): Promise<void> {
  loading.value = true;
  errorText.value = '';
  try {
    await batchUpdateArticleStatus(selectedIds.value, status);
    await reload(page.value);
  } catch {
    errorText.value = t('common.saveFailed');
    loading.value = false;
  }
}

async function removeSelected(): Promise<void> {
  if (!window.confirm(t('common.confirmDelete'))) {
    return;
  }
  loading.value = true;
  errorText.value = '';
  try {
    await batchDeleteArticles(selectedIds.value);
    await reload(page.value);
  } catch {
    errorText.value = t('common.deleteFailed');
    loading.value = false;
  }
}

async function removeOne(id: string): Promise<void> {
  if (!window.confirm(t('common.confirmDelete'))) {
    return;
  }
  loading.value = true;
  errorText.value = '';
  try {
    await deleteArticle(id);
    await reload(page.value);
  } catch {
    errorText.value = t('common.deleteFailed');
    loading.value = false;
  }
}

onMounted(async () => {
  await reload(1);
});
</script>

<style scoped>
.panel {
  padding: 20px;
  border: 1px solid var(--border);
  border-radius: 12px;
  background: var(--surface);
}

.subnav {
  display: flex;
  gap: 8px;
  margin-bottom: 12px;
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

.panel-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  margin-bottom: 14px;
}

.bulk-actions {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 10px;
  margin-bottom: 10px;
}

.check-all {
  display: inline-flex;
  align-items: center;
  gap: 8px;
}

.bulk-buttons {
  display: flex;
  gap: 8px;
}

h1 {
  margin: 0;
}

p {
  margin: 4px 0 0;
  color: var(--muted);
}

.create-btn {
  text-decoration: none;
  background: var(--accent);
  color: #fff;
  border-radius: 8px;
  padding: 8px 12px;
}

.filters {
  display: grid;
  grid-template-columns: 1.2fr 0.8fr 0.8fr auto;
  gap: 8px;
  margin-bottom: 14px;
}

input,
select,
button {
  border: 1px solid var(--border);
  border-radius: 8px;
  padding: 8px 10px;
  font: inherit;
}

button.danger {
  border-color: #e8b9b9;
  color: #b64040;
}

button {
  background: var(--surface);
  cursor: pointer;
}

.error {
  color: #b64040;
  margin-bottom: 10px;
}

.table-wrap {
  overflow: auto;
}

table {
  width: 100%;
  border-collapse: collapse;
}

th,
td {
  text-align: left;
  padding: 10px 8px;
  border-bottom: 1px solid var(--border);
}

th {
  color: var(--muted);
  font-weight: 600;
}

.mono {
  font-family: ui-monospace, SFMono-Regular, Menlo, monospace;
}

.link-btn {
  border: 0;
  background: transparent;
  padding: 0;
  margin-left: 8px;
  cursor: pointer;
}

.pager {
  margin-top: 12px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 10px;
}

.pager-actions {
  display: flex;
  align-items: center;
  gap: 8px;
}

@media (max-width: 900px) {
  .panel {
    padding: 14px;
  }

  .filters {
    grid-template-columns: 1fr;
  }

  .panel-header {
    align-items: flex-start;
    flex-direction: column;
  }

  .create-btn {
    width: 100%;
    text-align: center;
  }

  .bulk-actions {
    flex-direction: column;
    align-items: stretch;
  }

  .bulk-buttons {
    width: 100%;
  }

  .bulk-buttons button {
    flex: 1;
  }
}
</style>
