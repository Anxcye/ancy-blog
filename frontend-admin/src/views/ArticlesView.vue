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
            <td>{{ item.title }}</td>
            <td class="mono">{{ item.slug }}</td>
            <td>{{ item.status }}</td>
            <td>{{ item.contentKind }}</td>
            <td>{{ formatDate(item.updatedAt) }}</td>
            <td>
              <RouterLink :to="{ name: 'article-edit', params: { id: item.id } }">
                {{ t('articles.edit') }}
              </RouterLink>
            </td>
          </tr>
          <tr v-if="!loading && rows.length === 0">
            <td colspan="6">{{ t('articles.empty') }}</td>
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

import { listArticles } from '@/api/modules/articles';
import type { Article } from '@/api/types';

const { t } = useI18n();

const loading = ref(false);
const errorText = ref('');
const rows = ref<Article[]>([]);
const total = ref(0);
const page = ref(1);
const pageSize = 12;

const filters = reactive({
  keyword: '',
  status: '',
  contentKind: '',
});

const totalPages = computed(() => Math.max(1, Math.ceil(total.value / pageSize)));

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
    total.value = result.total;
    page.value = nextPage;
  } catch {
    errorText.value = t('common.loadFailed');
  } finally {
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
}
</style>
