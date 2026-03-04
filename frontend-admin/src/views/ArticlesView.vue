<!--
File: ArticlesView.vue
Purpose: Render article management table with filters, batch operations, and editor navigation.
Module: frontend-admin/views/content, presentation layer.
Related: articles API module, translation center routing, dashboard content workflow.
-->
<template>
  <section class="articles-page">
    <NCard :bordered="false" class="section-card">
      <div class="topbar">
        <NTabs type="segment" :value="'articles'">
          <NTabPane name="articles" :tab="t('articles.tabArticles')" />
          <NTabPane name="moments" :tab="t('articles.tabMoments')" @click="router.push({ name: 'moments' })" />
        </NTabs>

        <NButton type="primary" @click="router.push({ name: 'article-new' })">{{ t('articles.create') }}</NButton>
      </div>

      <NForm inline :show-label="false" class="filters" @submit.prevent="reload(1)">
        <NFormItem>
          <NInput v-model:value="filters.keyword" :placeholder="t('articles.filterKeyword')" clearable />
        </NFormItem>
        <NFormItem>
          <NSelect v-model:value="filters.status" :options="statusOptions" :placeholder="t('articles.filterStatusAll')" clearable style="width: 150px" />
        </NFormItem>
        <NFormItem>
          <NSelect v-model:value="filters.contentKind" :options="kindOptions" :placeholder="t('articles.filterKindAll')" clearable style="width: 150px" />
        </NFormItem>
        <NFormItem>
          <NButton attr-type="submit">{{ t('common.search') }}</NButton>
        </NFormItem>
      </NForm>

      <div class="batch-row">
        <span class="hint">{{ t('articles.selected', { count: selectedRowKeys.length }) }}</span>
        <NSpace>
          <NButton :disabled="selectedRowKeys.length === 0 || loading" @click="applyStatus('draft')">{{ t('articles.toDraft') }}</NButton>
          <NButton :disabled="selectedRowKeys.length === 0 || loading" @click="applyStatus('published')">{{ t('articles.toPublished') }}</NButton>
          <NButton type="error" tertiary :disabled="selectedRowKeys.length === 0 || loading" @click="removeSelected">{{ t('common.delete') }}</NButton>
        </NSpace>
      </div>

      <NAlert v-if="errorText" type="error" :show-icon="false">{{ errorText }}</NAlert>

      <NDataTable
        remote
        class="table"
        :loading="loading"
        :columns="columns"
        :data="rows"
        :pagination="false"
        :row-key="rowKey"
        :checked-row-keys="selectedRowKeys"
        @update:checked-row-keys="handleCheckedRows"
      />

      <div class="footer-row">
        <span class="hint">{{ t('articles.total', { total }) }}</span>
        <NPagination
          :page="page"
          :page-size="pageSize"
          :item-count="total"
          :page-slot="7"
          @update:page="reload"
        />
      </div>
    </NCard>
  </section>
</template>

<script setup lang="ts">
import { computed, h, onMounted, reactive, ref } from 'vue';
import { RouterLink, useRouter } from 'vue-router';
import { useI18n } from 'vue-i18n';
import type { DataTableColumns } from 'naive-ui';
import { NAlert, NButton, NCard, NDataTable, NForm, NFormItem, NInput, NPagination, NSelect, NSpace, NTabPane, NTabs, NTag, useDialog, useMessage } from 'naive-ui';

import { batchDeleteArticles, batchUpdateArticleStatus, deleteArticle, listArticles } from '@/api/modules/articles';
import type { Article } from '@/api/types';

const { t } = useI18n();
const router = useRouter();
const dialog = useDialog();
const message = useMessage();

const loading = ref(false);
const errorText = ref('');
const rows = ref<Article[]>([]);
const selectedRowKeys = ref<Array<string | number>>([]);
const total = ref(0);
const page = ref(1);
const pageSize = 12;

const filters = reactive({
  keyword: null as string | null,
  status: null as string | null,
  contentKind: null as string | null,
});

const statusOptions = computed(() => [
  { label: t('articles.statusDraft'), value: 'draft' },
  { label: t('articles.statusPublished'), value: 'published' },
  { label: t('articles.statusScheduled'), value: 'scheduled' },
]);

const kindOptions = [
  { label: 'Post', value: 'post' },
  { label: 'Page', value: 'page' },
];

function rowKey(row: Article): string {
  return row.id;
}

function handleCheckedRows(keys: Array<string | number>): void {
  selectedRowKeys.value = keys;
}

function formatDate(value: string): string {
  if (!value) {
    return '-';
  }
  return new Date(value).toLocaleString();
}

const columns = computed<DataTableColumns<Article>>(() => [
  {
    type: 'selection',
  },
  {
    title: t('articles.colTitle'),
    key: 'title',
    ellipsis: {
      tooltip: true,
    },
  },
  {
    title: t('articles.colSlug'),
    key: 'slug',
    ellipsis: {
      tooltip: true,
    },
  },
  {
    title: t('articles.colStatus'),
    key: 'status',
    render(row) {
      const typeMap: Record<string, 'default' | 'success' | 'warning'> = {
        draft: 'default',
        published: 'success',
        scheduled: 'warning',
      };
      return h(NTag, { type: typeMap[row.status] || 'default', round: true }, { default: () => row.status });
    },
  },
  {
    title: t('articles.colKind'),
    key: 'contentKind',
  },
  {
    title: t('articles.colUpdatedAt'),
    key: 'updatedAt',
    render(row) {
      return formatDate(row.updatedAt);
    },
  },
  {
    title: t('articles.colAction'),
    key: 'actions',
    width: 220,
    render(row) {
      return h(NSpace, { wrap: false, size: 8 }, {
        default: () => [
          h(
            RouterLink,
            { to: { name: 'article-edit', params: { id: row.id } } },
            {
              default: () => h(NButton, { size: 'small', tertiary: true }, { default: () => t('articles.edit') }),
            },
          ),
          h(
            RouterLink,
            { to: { name: 'system', query: { tab: 'translations', sourceType: 'article', sourceId: row.id } } },
            {
              default: () => h(NButton, { size: 'small', tertiary: true }, { default: () => t('articles.toTranslate') }),
            },
          ),
          h(
            NButton,
            {
              size: 'small',
              tertiary: true,
              type: 'error',
              onClick: () => {
                removeOne(row.id);
              },
            },
            { default: () => t('common.delete') },
          ),
        ],
      });
    },
  },
]);

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
    selectedRowKeys.value = [];
  } catch {
    errorText.value = t('common.loadFailed');
  } finally {
    loading.value = false;
  }
}

async function applyStatus(status: 'draft' | 'published'): Promise<void> {
  if (selectedRowKeys.value.length === 0) {
    return;
  }
  loading.value = true;
  errorText.value = '';
  try {
    await batchUpdateArticleStatus(selectedRowKeys.value.map((item) => String(item)), status);
    message.success(t('common.saveSuccess'));
    await reload(page.value);
  } catch {
    errorText.value = t('common.saveFailed');
    loading.value = false;
  }
}

function removeSelected(): void {
  if (selectedRowKeys.value.length === 0) {
    return;
  }
  dialog.warning({
    title: t('common.delete'),
    content: t('common.confirmDelete'),
    positiveText: t('common.delete'),
    negativeText: t('common.cancel'),
    onPositiveClick: async () => {
      loading.value = true;
      errorText.value = '';
      try {
        await batchDeleteArticles(selectedRowKeys.value.map((item) => String(item)));
        message.success(t('common.deleteSuccess'));
        await reload(page.value);
      } catch {
        errorText.value = t('common.deleteFailed');
        loading.value = false;
      }
    },
  });
}

function removeOne(id: string): void {
  dialog.warning({
    title: t('common.delete'),
    content: t('common.confirmDelete'),
    positiveText: t('common.delete'),
    negativeText: t('common.cancel'),
    onPositiveClick: async () => {
      loading.value = true;
      errorText.value = '';
      try {
        await deleteArticle(id);
        message.success(t('common.deleteSuccess'));
        await reload(page.value);
      } catch {
        errorText.value = t('common.deleteFailed');
        loading.value = false;
      }
    },
  });
}

onMounted(async () => {
  await reload(1);
});
</script>

<style scoped>
.articles-page {
  display: grid;
}

.section-card {
  border-radius: 14px;
  box-shadow: 0 6px 24px rgba(15, 31, 36, 0.05);
}

.topbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 12px;
  margin-bottom: 12px;
}

.filters {
  margin-bottom: 12px;
}

.batch-row,
.footer-row {
  margin: 10px 0 12px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 10px;
}

.hint {
  color: #6e7a84;
  font-size: 13px;
}

.table {
  overflow: hidden;
}

@media (max-width: 900px) {
  .topbar,
  .batch-row,
  .footer-row {
    flex-direction: column;
    align-items: stretch;
  }
}
</style>
