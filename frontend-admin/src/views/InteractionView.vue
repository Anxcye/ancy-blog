<!--
File: InteractionView.vue
Purpose: Manage comment moderation and friend-link review tasks with unified Naive UI components.
Module: frontend-admin/views/interaction, presentation layer.
Related: interactions API module and backend admin moderation endpoints.
-->
<template>
  <section class="interaction-page">
    <NCard :bordered="false" class="section-card">
      <div class="topbar">
        <div class="tab-switch">
          <NButton :type="tab === 'comments' ? 'primary' : 'default'" secondary @click="switchTab('comments')">{{ t('interaction.commentsTab') }}</NButton>
          <NButton :type="tab === 'links' ? 'primary' : 'default'" secondary @click="switchTab('links')">{{ t('interaction.linksTab') }}</NButton>
        </div>

        <NSelect
          v-if="tab === 'comments'"
          v-model:value="commentStatus"
          :options="commentStatusOptions"
          :placeholder="t('interaction.allStatus')"
          clearable
          style="width: 180px"
          @update:value="() => loadComments(1)"
        />

        <NSelect
          v-else
          v-model:value="linkStatus"
          :options="linkStatusOptions"
          :placeholder="t('interaction.allStatus')"
          clearable
          style="width: 180px"
          @update:value="() => loadLinks(1)"
        />
      </div>

      <NAlert v-if="errorText" type="error" :show-icon="false">{{ errorText }}</NAlert>
      <NAlert v-if="successText" type="success" :show-icon="false">{{ successText }}</NAlert>

      <NDataTable
        v-if="tab === 'comments'"
        remote
        :loading="loading"
        :columns="commentColumns"
        :data="comments"
        :pagination="false"
        :scroll-x="980"
        :row-key="rowKey"
      />

      <NDataTable
        v-else
        remote
        :loading="loading"
        :columns="linkColumns"
        :data="links"
        :pagination="false"
        :scroll-x="900"
        :row-key="rowKey"
      />

      <div class="footer-row">
        <span class="hint">{{ t('articles.total', { total: tab === 'comments' ? commentTotal : linkTotal }) }}</span>
        <NPagination
          v-if="tab === 'comments'"
          :page="commentPage"
          :page-size="commentPageSize"
          :item-count="commentTotal"
          :page-slot="isMobile ? 3 : 7"
          :simple="isMobile"
          @update:page="loadComments"
        />
        <NPagination
          v-else
          :page="linkPage"
          :page-size="linkPageSize"
          :item-count="linkTotal"
          :page-slot="isMobile ? 3 : 7"
          :simple="isMobile"
          @update:page="loadLinks"
        />
      </div>
    </NCard>
  </section>
</template>

<script setup lang="ts">
import { computed, h, onBeforeUnmount, onMounted, ref } from 'vue';
import { useI18n } from 'vue-i18n';
import type { DataTableColumns } from 'naive-ui';
import { NAlert, NButton, NCard, NDataTable, NPagination, NSelect, NSpace } from 'naive-ui';

import { listComments, listLinkSubmissions, reviewLink, updateComment } from '@/api/modules/interactions';
import type { Comment, LinkSubmission } from '@/api/types';

type InteractionTab = 'comments' | 'links';

const { t } = useI18n();

const tab = ref<InteractionTab>('comments');
const loading = ref(false);
const errorText = ref('');
const successText = ref('');
const isMobile = ref(false);

const commentStatus = ref<string | null>(null);
const comments = ref<Comment[]>([]);
const commentPage = ref(1);
const commentTotal = ref(0);
const commentPageSize = 20;

const linkStatus = ref<string | null>(null);
const links = ref<LinkSubmission[]>([]);
const linkPage = ref(1);
const linkTotal = ref(0);
const linkPageSize = 20;

const commentStatusOptions = [
  { label: 'approved', value: 'approved' },
  { label: 'pending', value: 'pending' },
  { label: 'rejected', value: 'rejected' },
];

const linkStatusOptions = [
  { label: 'pending', value: 'pending' },
  { label: 'approved', value: 'approved' },
  { label: 'rejected', value: 'rejected' },
];

function rowKey(row: Comment | LinkSubmission): string {
  return row.id;
}

function syncViewport(): void {
  isMobile.value = window.innerWidth <= 900;
}

function formatDate(value: string): string {
  if (!value) {
    return '-';
  }
  return new Date(value).toLocaleString();
}

const commentColumns = computed<DataTableColumns<Comment>>(() => [
  {
    title: '用户',
    key: 'nickname',
    width: 150,
  },
  {
    title: '状态',
    key: 'status',
    width: 110,
  },
  {
    title: '内容',
    key: 'content',
    ellipsis: {
      tooltip: true,
    },
  },
  {
    title: '来源',
    key: 'meta',
    width: 220,
    render(row) {
      return `${row.ip} · ${formatDate(row.createdAt)}`;
    },
  },
  {
    title: '操作',
    key: 'actions',
    width: 240,
    render(row) {
      return h(NSpace, { wrap: false, size: 8 }, {
        default: () => [
          h(
            NButton,
            {
              size: 'small',
              tertiary: true,
              onClick: () => setCommentStatus(row.id, 'approved', row.isPinned),
            },
            { default: () => 'approve' },
          ),
          h(
            NButton,
            {
              size: 'small',
              tertiary: true,
              onClick: () => setCommentStatus(row.id, 'rejected', row.isPinned),
            },
            { default: () => 'reject' },
          ),
          h(
            NButton,
            {
              size: 'small',
              tertiary: true,
              onClick: () => setCommentStatus(row.id, row.status, row.isPinned === '1' ? '0' : '1'),
            },
            { default: () => (row.isPinned === '1' ? 'unpin' : 'pin') },
          ),
        ],
      });
    },
  },
]);

const linkColumns = computed<DataTableColumns<LinkSubmission>>(() => [
  {
    title: '名称',
    key: 'name',
    width: 180,
  },
  {
    title: '状态',
    key: 'reviewStatus',
    width: 110,
  },
  {
    title: '地址',
    key: 'url',
    ellipsis: {
      tooltip: true,
    },
  },
  {
    title: '来源',
    key: 'meta',
    width: 220,
    render(row) {
      return `${row.contactEmail || '-'} · ${row.submittedIp}`;
    },
  },
  {
    title: '操作',
    key: 'actions',
    width: 160,
    render(row) {
      return h(NSpace, { wrap: false, size: 8 }, {
        default: () => [
          h(
            NButton,
            {
              size: 'small',
              tertiary: true,
              onClick: () => setLinkStatus(row.id, 'approved'),
            },
            { default: () => 'approve' },
          ),
          h(
            NButton,
            {
              size: 'small',
              tertiary: true,
              onClick: () => setLinkStatus(row.id, 'rejected'),
            },
            { default: () => 'reject' },
          ),
        ],
      });
    },
  },
]);

async function switchTab(nextTab: InteractionTab): Promise<void> {
  tab.value = nextTab;
  if (nextTab === 'comments') {
    await loadComments(1);
  } else {
    await loadLinks(1);
  }
}

async function loadComments(page: number): Promise<void> {
  loading.value = true;
  errorText.value = '';
  try {
    const result = await listComments({ page, pageSize: commentPageSize, status: commentStatus.value || undefined });
    comments.value = result.rows;
    commentTotal.value = result.total;
    commentPage.value = page;
  } catch {
    errorText.value = t('common.loadFailed');
  } finally {
    loading.value = false;
  }
}

async function setCommentStatus(id: string, status: string, isPinned: string): Promise<void> {
  loading.value = true;
  errorText.value = '';
  successText.value = '';
  try {
    await updateComment(id, { status, isPinned });
    successText.value = t('common.saveSuccess');
    await loadComments(commentPage.value);
  } catch {
    errorText.value = t('common.saveFailed');
    loading.value = false;
  }
}

async function loadLinks(page: number): Promise<void> {
  loading.value = true;
  errorText.value = '';
  try {
    const result = await listLinkSubmissions({ page, pageSize: linkPageSize, reviewStatus: linkStatus.value || undefined });
    links.value = result.rows;
    linkTotal.value = result.total;
    linkPage.value = page;
  } catch {
    errorText.value = t('common.loadFailed');
  } finally {
    loading.value = false;
  }
}

async function setLinkStatus(id: string, reviewStatus: string): Promise<void> {
  loading.value = true;
  errorText.value = '';
  successText.value = '';
  try {
    await reviewLink(id, { reviewStatus, reviewNote: '', relatedArticleId: '' });
    successText.value = t('common.saveSuccess');
    await loadLinks(linkPage.value);
  } catch {
    errorText.value = t('common.saveFailed');
    loading.value = false;
  }
}

onMounted(async () => {
  syncViewport();
  window.addEventListener('resize', syncViewport);
  await loadComments(1);
});

onBeforeUnmount(() => {
  window.removeEventListener('resize', syncViewport);
});
</script>

<style scoped>
.interaction-page {
  display: grid;
  max-width: 100%;
  overflow-x: hidden;
}

.section-card {
  border-radius: 14px;
  box-shadow: 0 6px 24px color-mix(in srgb, var(--n-text-color) 8%, transparent);
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

.tab-switch {
  display: flex;
  gap: 8px;
}

.footer-row {
  margin-top: 12px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 10px;
}

.hint {
  color: var(--n-text-color-3);
  font-size: 13px;
}

@media (max-width: 900px) {
  .topbar,
  .footer-row {
    flex-direction: column;
    align-items: stretch;
  }
}
</style>
