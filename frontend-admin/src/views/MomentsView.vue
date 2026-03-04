<!--
File: MomentsView.vue
Purpose: Render moment management table with filters, batch operations, and create/edit modal form.
Module: frontend-admin/views/content, presentation layer.
Related: moments API module, article list route, backend admin moment endpoints.
-->
<template>
  <section class="moments-page">
    <NCard :bordered="false" class="section-card">
      <div class="topbar">
        <div class="content-switch">
          <NButton secondary @click="router.push({ name: 'articles' })">{{ t('articles.tabArticles') }}</NButton>
          <NButton type="primary" secondary>{{ t('articles.tabMoments') }}</NButton>
        </div>

        <NButton type="primary" @click="openCreate">{{ t('moments.create') }}</NButton>
      </div>

      <NForm inline :show-label="false" class="filters" @submit.prevent="reload(1)">
        <NFormItem>
          <NSelect v-model:value="statusFilter" :options="statusOptions" :placeholder="t('moments.statusAll')" clearable style="width: 170px" />
        </NFormItem>
        <NFormItem>
          <NButton attr-type="submit">{{ t('common.search') }}</NButton>
        </NFormItem>
      </NForm>

      <div class="batch-row">
        <span class="hint">{{ t('moments.selected', { count: selectedRowKeys.length }) }}</span>
        <NSpace>
          <NButton :disabled="selectedRowKeys.length === 0 || loading" @click="applyStatus('draft')">{{ t('moments.toDraft') }}</NButton>
          <NButton :disabled="selectedRowKeys.length === 0 || loading" @click="applyStatus('published')">{{ t('moments.toPublished') }}</NButton>
          <NButton type="error" tertiary :disabled="selectedRowKeys.length === 0 || loading" @click="removeSelected">{{ t('common.delete') }}</NButton>
        </NSpace>
      </div>

      <NAlert v-if="errorText" type="error" :show-icon="false">{{ errorText }}</NAlert>

      <NDataTable
        remote
        :loading="loading"
        :columns="columns"
        :data="rows"
        :pagination="false"
        :scroll-x="920"
        :row-key="rowKey"
        :checked-row-keys="selectedRowKeys"
        @update:checked-row-keys="handleCheckedRows"
      />

      <div class="footer-row">
        <span class="hint">{{ t('articles.total', { total }) }}</span>
        <NPagination :page="page" :page-size="pageSize" :item-count="total" :page-slot="isMobile ? 3 : 7" :simple="isMobile" @update:page="reload" />
      </div>
    </NCard>

    <NModal v-model:show="showForm" preset="card" :title="editingId ? t('moments.edit') : t('moments.create')" style="width: min(760px, 96vw)">
      <NForm label-placement="top" class="modal-form">
        <NFormItem :label="t('moments.fieldContent')" required>
          <NInput v-model:value="form.content" type="textarea" :autosize="{ minRows: 7, maxRows: 14 }" />
        </NFormItem>

        <div class="modal-grid">
          <NFormItem :label="t('moments.fieldStatus')">
            <NSelect v-model:value="form.status" :options="statusOptions" />
          </NFormItem>
          <NFormItem :label="t('moments.fieldPublishedAt')">
            <NDatePicker v-model:value="publishedAtTs" type="datetime" clearable style="width: 100%" />
          </NFormItem>
        </div>

        <NFormItem>
          <NCheckbox v-model:checked="form.allowComment">{{ t('moments.fieldAllowComment') }}</NCheckbox>
        </NFormItem>
      </NForm>

      <template #footer>
        <NSpace justify="end">
          <NButton @click="showForm = false">{{ t('common.cancel') }}</NButton>
          <NButton type="primary" :loading="loading" @click="saveMoment">{{ t('common.save') }}</NButton>
        </NSpace>
      </template>
    </NModal>
  </section>
</template>

<script setup lang="ts">
import { computed, h, onBeforeUnmount, onMounted, reactive, ref, watch } from 'vue';
import { useRouter } from 'vue-router';
import { useI18n } from 'vue-i18n';
import type { DataTableColumns } from 'naive-ui';
import { NAlert, NButton, NCard, NCheckbox, NDataTable, NDatePicker, NForm, NFormItem, NInput, NModal, NPagination, NSelect, NSpace, NTag, useDialog, useMessage } from 'naive-ui';

import { batchDeleteMoments, batchUpdateMomentStatus, createMoment, deleteMoment, listMoments, type Moment, updateMoment } from '@/api/modules/moments';

const { t } = useI18n();
const router = useRouter();
const dialog = useDialog();
const message = useMessage();

const loading = ref(false);
const errorText = ref('');

const rows = ref<Moment[]>([]);
const total = ref(0);
const page = ref(1);
const pageSize = 20;
const isMobile = ref(false);
const statusFilter = ref<string | null>(null);
const selectedRowKeys = ref<Array<string | number>>([]);

const showForm = ref(false);
const editingId = ref('');
const publishedAtTs = ref<number | null>(null);
const form = reactive({
  content: '',
  status: 'draft',
  allowComment: true,
  publishedAt: '',
});

watch(publishedAtTs, (value) => {
  form.publishedAt = value ? new Date(value).toISOString() : '';
});

const statusOptions = computed(() => [
  { label: 'draft', value: 'draft' },
  { label: 'published', value: 'published' },
  { label: 'scheduled', value: 'scheduled' },
]);

function rowKey(row: Moment): string {
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

const columns = computed<DataTableColumns<Moment>>(() => [
  { type: 'selection' },
  {
    title: 'ID',
    key: 'id',
    width: 220,
    ellipsis: {
      tooltip: true,
    },
  },
  {
    title: t('moments.fieldContent'),
    key: 'content',
    ellipsis: {
      tooltip: true,
    },
  },
  {
    title: t('moments.fieldStatus'),
    key: 'status',
    width: 120,
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
    title: t('articles.colUpdatedAt'),
    key: 'updatedAt',
    width: 180,
    render(row) {
      return formatDate(row.updatedAt);
    },
  },
  {
    title: t('articles.colAction'),
    key: 'actions',
    width: 170,
    render(row) {
      return h(NSpace, { wrap: false, size: 8 }, {
        default: () => [
          h(
            NButton,
            {
              size: 'small',
              tertiary: true,
              onClick: () => openEdit(row),
            },
            { default: () => t('moments.edit') },
          ),
          h(
            NButton,
            {
              size: 'small',
              tertiary: true,
              type: 'error',
              onClick: () => removeOne(row.id),
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
    const result = await listMoments({
      page: nextPage,
      pageSize,
      status: statusFilter.value || undefined,
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

function syncViewport(): void {
  isMobile.value = window.innerWidth <= 900;
}

function openCreate(): void {
  editingId.value = '';
  form.content = '';
  form.status = 'draft';
  form.allowComment = true;
  form.publishedAt = '';
  publishedAtTs.value = null;
  showForm.value = true;
}

function openEdit(item: Moment): void {
  editingId.value = item.id;
  form.content = item.content;
  form.status = item.status;
  form.allowComment = item.allowComment;
  form.publishedAt = item.publishedAt || '';
  publishedAtTs.value = form.publishedAt ? new Date(form.publishedAt).getTime() : null;
  showForm.value = true;
}

async function saveMoment(): Promise<void> {
  if (!form.content.trim()) {
    errorText.value = t('moments.contentRequired');
    return;
  }
  loading.value = true;
  errorText.value = '';
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
    message.success(t('common.saveSuccess'));
    showForm.value = false;
    await reload(1);
  } catch {
    errorText.value = t('common.saveFailed');
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
    await batchUpdateMomentStatus(selectedRowKeys.value.map((item) => String(item)), status);
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
        await batchDeleteMoments(selectedRowKeys.value.map((item) => String(item)));
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
        await deleteMoment(id);
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
  syncViewport();
  window.addEventListener('resize', syncViewport);
  await reload(1);
});

onBeforeUnmount(() => {
  window.removeEventListener('resize', syncViewport);
});
</script>

<style scoped>
.moments-page {
  display: grid;
  max-width: 100%;
  overflow-x: hidden;
}

.section-card {
  border-radius: 14px;
  box-shadow: 0 6px 24px rgba(15, 31, 36, 0.05);
  max-width: 100%;
  overflow-x: hidden;
}

.topbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 12px;
  margin-bottom: 12px;
}

.content-switch {
  display: flex;
  gap: 8px;
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
  color: var(--n-text-color-3);
  font-size: 13px;
}

.modal-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 12px;
}

@media (max-width: 900px) {
  .topbar,
  .batch-row,
  .footer-row {
    flex-direction: column;
    align-items: stretch;
  }

  .modal-grid {
    grid-template-columns: 1fr;
  }
}
</style>
