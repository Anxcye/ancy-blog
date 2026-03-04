<!--
File: InteractionView.vue
Purpose: Manage comment moderation and friend-link review tasks in one workspace.
Module: frontend-admin/views/interaction, presentation layer.
Related: interactions API module and backend admin moderation endpoints.
-->
<template>
  <section class="interaction-page">
    <h1>{{ t('interaction.title') }}</h1>
    <p class="subtitle">{{ t('interaction.subtitle') }}</p>

    <div class="mode-tabs">
      <button :class="{ active: tab === 'comments' }" @click="switchTab('comments')">{{ t('interaction.commentsTab') }}</button>
      <button :class="{ active: tab === 'links' }" @click="switchTab('links')">{{ t('interaction.linksTab') }}</button>
    </div>

    <p v-if="errorText" class="error">{{ errorText }}</p>
    <p v-if="successText" class="success">{{ successText }}</p>

    <section v-if="tab === 'comments'" class="panel">
      <div class="toolbar">
        <select v-model="commentStatus" @change="loadComments(1)">
          <option value="">{{ t('interaction.allStatus') }}</option>
          <option value="approved">approved</option>
          <option value="pending">pending</option>
          <option value="rejected">rejected</option>
        </select>
      </div>

      <ul class="list">
        <li v-for="item in comments" :key="item.id" class="item">
          <div>
            <p class="item-title">{{ item.nickname }} · {{ item.status }}</p>
            <p class="item-content">{{ item.content }}</p>
            <p class="item-meta">{{ item.ip }} · {{ formatDate(item.createdAt) }}</p>
          </div>
          <div class="item-actions">
            <button @click="setCommentStatus(item.id, 'approved', item.isPinned)">approve</button>
            <button @click="setCommentStatus(item.id, 'rejected', item.isPinned)">reject</button>
            <button @click="setCommentStatus(item.id, item.status, item.isPinned === '1' ? '0' : '1')">
              {{ item.isPinned === '1' ? 'unpin' : 'pin' }}
            </button>
          </div>
        </li>
      </ul>
    </section>

    <section v-else class="panel">
      <div class="toolbar">
        <select v-model="linkStatus" @change="loadLinks(1)">
          <option value="">{{ t('interaction.allStatus') }}</option>
          <option value="pending">pending</option>
          <option value="approved">approved</option>
          <option value="rejected">rejected</option>
        </select>
      </div>

      <ul class="list">
        <li v-for="item in links" :key="item.id" class="item">
          <div>
            <p class="item-title">{{ item.name }} · {{ item.reviewStatus }}</p>
            <p class="item-content"><a :href="item.url" target="_blank" rel="noreferrer">{{ item.url }}</a></p>
            <p class="item-meta">{{ item.contactEmail || '-' }} · {{ item.submittedIp }}</p>
          </div>
          <div class="item-actions">
            <button @click="setLinkStatus(item.id, 'approved')">approve</button>
            <button @click="setLinkStatus(item.id, 'rejected')">reject</button>
          </div>
        </li>
      </ul>
    </section>
  </section>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue';
import { useI18n } from 'vue-i18n';

import { listComments, listLinkSubmissions, reviewLink, updateComment } from '@/api/modules/interactions';
import type { Comment, LinkSubmission } from '@/api/types';

type InteractionTab = 'comments' | 'links';

const { t } = useI18n();

const tab = ref<InteractionTab>('comments');
const loading = ref(false);
const errorText = ref('');
const successText = ref('');

const commentStatus = ref('');
const comments = ref<Comment[]>([]);

const linkStatus = ref('');
const links = ref<LinkSubmission[]>([]);

function formatDate(value: string): string {
  if (!value) {
    return '-';
  }
  return new Date(value).toLocaleString();
}

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
    const result = await listComments({ page, pageSize: 20, status: commentStatus.value || undefined });
    comments.value = result.rows;
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
    await loadComments(1);
  } catch {
    errorText.value = t('common.saveFailed');
    loading.value = false;
  }
}

async function loadLinks(page: number): Promise<void> {
  loading.value = true;
  errorText.value = '';
  try {
    const result = await listLinkSubmissions({ page, pageSize: 20, reviewStatus: linkStatus.value || undefined });
    links.value = result.rows;
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
    await loadLinks(1);
  } catch {
    errorText.value = t('common.saveFailed');
    loading.value = false;
  }
}

onMounted(async () => {
  await loadComments(1);
});
</script>

<style scoped>
.interaction-page {
  display: grid;
  gap: 12px;
}

h1 {
  margin: 0;
}

.subtitle {
  margin: -4px 0 0;
  color: var(--muted);
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
  padding: 8px 12px;
  background: var(--surface);
  cursor: pointer;
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
  border: 1px solid var(--border);
  border-radius: 12px;
  background: var(--surface);
  padding: 14px;
  display: grid;
  gap: 10px;
}

.toolbar {
  display: flex;
  justify-content: flex-end;
}

select,
button {
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

.item-title,
.item-content,
.item-meta {
  margin: 0;
}

.item-title {
  font-weight: 600;
}

.item-content {
  margin-top: 6px;
  color: var(--text);
  white-space: pre-wrap;
}

.item-meta {
  margin-top: 6px;
  color: var(--muted);
  font-size: 13px;
}

.item-actions {
  display: flex;
  gap: 8px;
  align-items: flex-start;
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
  .item {
    flex-direction: column;
  }

  .item-actions {
    width: 100%;
  }

  .item-actions button {
    flex: 1;
  }
}
</style>
