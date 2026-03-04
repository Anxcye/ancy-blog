<!--
File: SystemView.vue
Purpose: Manage integrations and translation jobs/content operations in a unified system workspace.
Module: frontend-admin/views/system, presentation layer.
Related: integrations API module, translations API module, and backend admin system endpoints.
-->
<template>
  <section class="system-page">
    <h1>{{ t('system.title') }}</h1>
    <p class="subtitle">{{ t('system.subtitle') }}</p>

    <div class="mode-tabs">
      <button :class="{ active: tab === 'integrations' }" @click="tab = 'integrations'">{{ t('system.tabIntegrations') }}</button>
      <button :class="{ active: tab === 'translations' }" @click="tab = 'translations'">{{ t('system.tabTranslations') }}</button>
    </div>

    <p v-if="errorText" class="error">{{ errorText }}</p>
    <p v-if="successText" class="success">{{ successText }}</p>

    <template v-if="tab === 'integrations'">
      <article v-for="provider in providers" :key="provider.providerKey" class="panel">
        <header class="panel-header">
          <div>
            <h2>{{ provider.name }}</h2>
            <p>{{ provider.providerKey }} · {{ provider.providerType }}</p>
          </div>
          <label class="switch-row">
            <input v-model="provider.enabled" type="checkbox" />
            <span>{{ provider.enabled ? t('system.enabled') : t('system.disabled') }}</span>
          </label>
        </header>

        <label>
          <span>{{ t('system.configJson') }}</span>
          <textarea v-model="provider.configRaw" rows="6" />
        </label>

        <label>
          <span>{{ t('system.metaJson') }}</span>
          <textarea v-model="provider.metaRaw" rows="4" />
        </label>

        <div class="actions">
          <button :disabled="loading" @click="saveProvider(provider)">{{ t('common.save') }}</button>
          <button :disabled="loading" @click="testProvider(provider.providerKey)">{{ t('system.test') }}</button>
        </div>
      </article>
    </template>

    <template v-else>
      <article class="panel">
        <h2>{{ t('system.translationCreateTitle') }}</h2>
        <div class="grid-2">
          <label>
            <span>sourceType</span>
            <select v-model="jobForm.sourceType">
              <option value="article">article</option>
              <option value="moment">moment</option>
            </select>
          </label>
          <label>
            <span>sourceId</span>
            <input v-model.trim="jobForm.sourceId" type="text" />
          </label>
          <label>
            <span>sourceLocale</span>
            <input v-model.trim="jobForm.sourceLocale" type="text" placeholder="zh-CN" />
          </label>
          <label>
            <span>targetLocale</span>
            <input v-model.trim="jobForm.targetLocale" type="text" placeholder="en-US" />
          </label>
          <label>
            <span>providerKey</span>
            <input v-model.trim="jobForm.providerKey" type="text" placeholder="openai_compatible" />
          </label>
          <label>
            <span>modelName</span>
            <input v-model.trim="jobForm.modelName" type="text" placeholder="gpt-4.1-mini" />
          </label>
          <label>
            <span>maxRetries</span>
            <input v-model.number="jobForm.maxRetries" type="number" min="0" max="10" />
          </label>
          <label>
            <span>publishAt</span>
            <input v-model="jobPublishAtLocal" type="datetime-local" />
          </label>
          <label class="switch-row">
            <input v-model="jobForm.autoPublish" type="checkbox" />
            <span>autoPublish</span>
          </label>
        </div>
        <div class="actions">
          <button :disabled="loading" @click="createJob">{{ t('system.createTranslationJob') }}</button>
        </div>
      </article>

      <article class="panel">
        <header class="panel-header">
          <h2>{{ t('system.translationJobsTitle') }}</h2>
          <div class="toolbar">
            <select v-model="jobStatusFilter" @change="loadJobs()">
              <option value="">all</option>
              <option value="queued">queued</option>
              <option value="running">running</option>
              <option value="succeeded">succeeded</option>
              <option value="failed">failed</option>
            </select>
          </div>
        </header>

        <ul class="list">
          <li v-for="job in jobs" :key="job.id" class="item">
            <div>
              <p class="item-title">{{ job.sourceType }} · {{ job.targetLocale }} · {{ job.status }}</p>
              <p class="item-content">{{ job.id }}</p>
              <p class="item-meta">retry {{ job.retryCount }}/{{ job.maxRetries }} · {{ formatDate(job.updatedAt) }}</p>
            </div>
            <div class="item-actions">
              <button :disabled="loading || job.status !== 'failed'" @click="retryJob(job.id)">{{ t('system.retry') }}</button>
            </div>
          </li>
        </ul>
        <footer class="pager">
          <button :disabled="loading || jobPage <= 1" @click="loadJobs(jobPage - 1)">{{ t('common.prev') }}</button>
          <span>{{ jobPage }} / {{ jobTotalPages }}</span>
          <button :disabled="loading || jobPage >= jobTotalPages" @click="loadJobs(jobPage + 1)">{{ t('common.next') }}</button>
        </footer>
      </article>

      <article class="panel">
        <header class="panel-header">
          <h2>{{ t('system.translationContentsTitle') }}</h2>
          <div class="toolbar grid-mini">
            <select v-model="contentFilter.sourceType" @change="() => loadContents(1)">
              <option value="article">article</option>
              <option value="moment">moment</option>
            </select>
            <input v-model.trim="contentFilter.sourceId" placeholder="sourceId" @keyup.enter="() => loadContents(1)" />
            <input v-model.trim="contentFilter.locale" placeholder="locale" @keyup.enter="() => loadContents(1)" />
            <button :disabled="loading" @click="() => loadContents(1)">{{ t('common.search') }}</button>
          </div>
        </header>

        <ul class="list">
          <li v-for="row in contents" :key="`${row.sourceType}:${row.sourceId}:${row.locale}`" class="item-col">
            <p class="item-title">{{ row.sourceType }} · {{ row.sourceId }} · {{ row.locale }}</p>
            <input v-model="row.title" :placeholder="t('editor.fieldTitle')" />
            <input v-model="row.summary" :placeholder="t('editor.fieldSummary')" />
            <textarea v-model="row.content" rows="6" />
            <div class="grid-mini">
              <select v-model="row.status">
                <option value="draft">draft</option>
                <option value="published">published</option>
              </select>
              <input v-model="row.publishedAt" placeholder="2026-03-04T12:00:00Z" />
              <button :disabled="loading" @click="saveContent(row)">{{ t('common.save') }}</button>
            </div>
          </li>
        </ul>
        <footer class="pager">
          <button :disabled="loading || contentPage <= 1" @click="loadContents(contentPage - 1)">{{ t('common.prev') }}</button>
          <span>{{ contentPage }} / {{ contentTotalPages }}</span>
          <button :disabled="loading || contentPage >= contentTotalPages" @click="loadContents(contentPage + 1)">{{ t('common.next') }}</button>
        </footer>
      </article>
    </template>
  </section>
</template>

<script setup lang="ts">
import { computed, onMounted, reactive, ref } from 'vue';
import { useI18n } from 'vue-i18n';

import { listIntegrations, testIntegration, updateIntegration } from '@/api/modules/integrations';
import { createTranslationJob, listTranslationContents, listTranslationJobs, retryTranslationJob, upsertTranslationContent } from '@/api/modules/translations';
import type { TranslationContent, TranslationJob } from '@/api/types';

type IntegrationViewModel = {
  providerKey: string;
  providerType: string;
  name: string;
  enabled: boolean;
  configRaw: string;
  metaRaw: string;
};

type TabName = 'integrations' | 'translations';

const { t } = useI18n();

const tab = ref<TabName>('integrations');
const loading = ref(false);
const errorText = ref('');
const successText = ref('');
const providers = ref<IntegrationViewModel[]>([]);

const jobStatusFilter = ref('');
const jobs = ref<TranslationJob[]>([]);
const jobPage = ref(1);
const jobTotal = ref(0);
const jobPageSize = 20;
const jobPublishAtLocal = ref('');
const jobForm = reactive({
  sourceType: 'article',
  sourceId: '',
  sourceLocale: 'zh-CN',
  targetLocale: 'en-US',
  providerKey: 'openai_compatible',
  modelName: '',
  maxRetries: 3,
  autoPublish: true,
  publishAt: '',
});

const contentFilter = reactive({
  sourceType: 'article',
  sourceId: '',
  locale: '',
});
const contents = ref<TranslationContent[]>([]);
const contentPage = ref(1);
const contentTotal = ref(0);
const contentPageSize = 20;

const jobTotalPages = computed(() => Math.max(1, Math.ceil(jobTotal.value / jobPageSize)));
const contentTotalPages = computed(() => Math.max(1, Math.ceil(contentTotal.value / contentPageSize)));

function formatDate(value: string): string {
  if (!value) {
    return '-';
  }
  return new Date(value).toLocaleString();
}

function parseJSON(value: string): Record<string, unknown> {
  if (!value.trim()) {
    return {};
  }
  const parsed = JSON.parse(value);
  if (typeof parsed !== 'object' || Array.isArray(parsed) || parsed === null) {
    throw new Error('invalid json object');
  }
  return parsed as Record<string, unknown>;
}

async function loadProviders(): Promise<void> {
  const rows = await listIntegrations();
  providers.value = rows.map((item) => ({
    providerKey: item.providerKey,
    providerType: item.providerType,
    name: item.name,
    enabled: item.enabled,
    configRaw: JSON.stringify(item.configJson ?? {}, null, 2),
    metaRaw: JSON.stringify(item.metaJson ?? {}, null, 2),
  }));
}

async function saveProvider(item: IntegrationViewModel): Promise<void> {
  loading.value = true;
  errorText.value = '';
  successText.value = '';
  try {
    await updateIntegration(item.providerKey, {
      enabled: item.enabled,
      configJson: parseJSON(item.configRaw),
      metaJson: parseJSON(item.metaRaw),
    });
    successText.value = t('common.saveSuccess');
  } catch {
    errorText.value = t('system.invalidJsonOrSaveFailed');
  } finally {
    loading.value = false;
  }
}

async function testProvider(providerKey: string): Promise<void> {
  loading.value = true;
  errorText.value = '';
  successText.value = '';
  try {
    const result = await testIntegration(providerKey);
    successText.value = `${result.ok ? 'OK' : 'FAIL'} · ${result.message}`;
  } catch {
    errorText.value = t('system.testFailed');
  } finally {
    loading.value = false;
  }
}

async function loadJobs(page = 1): Promise<void> {
  loading.value = true;
  errorText.value = '';
  try {
    const result = await listTranslationJobs({
      page,
      pageSize: jobPageSize,
      status: jobStatusFilter.value || undefined,
    });
    jobs.value = result.rows;
    jobTotal.value = result.total;
    jobPage.value = page;
  } catch {
    errorText.value = t('common.loadFailed');
  } finally {
    loading.value = false;
  }
}

async function createJob(): Promise<void> {
  if (!jobForm.sourceId) {
    errorText.value = t('system.translationJobMissingSourceId');
    return;
  }
  loading.value = true;
  errorText.value = '';
  successText.value = '';
  try {
    jobForm.publishAt = jobPublishAtLocal.value ? new Date(jobPublishAtLocal.value).toISOString() : '';
    const id = await createTranslationJob({
      sourceType: jobForm.sourceType,
      sourceId: jobForm.sourceId,
      sourceLocale: jobForm.sourceLocale,
      targetLocale: jobForm.targetLocale,
      providerKey: jobForm.providerKey,
      modelName: jobForm.modelName,
      maxRetries: jobForm.maxRetries,
      autoPublish: jobForm.autoPublish,
      publishAt: jobForm.publishAt || undefined,
    });
    successText.value = `${t('system.translationJobCreated')} ${id}`;
    await loadJobs(1);
  } catch {
    errorText.value = t('common.saveFailed');
    loading.value = false;
  }
}

async function retryJob(id: string): Promise<void> {
  loading.value = true;
  errorText.value = '';
  successText.value = '';
  try {
    await retryTranslationJob(id);
    successText.value = t('system.retrySuccess');
    await loadJobs(jobPage.value);
  } catch {
    errorText.value = t('system.retryFailed');
    loading.value = false;
  }
}

async function loadContents(page = 1): Promise<void> {
  loading.value = true;
  errorText.value = '';
  try {
    const sourceType = contentFilter.sourceType || 'article';
    const result = await listTranslationContents({
      page,
      pageSize: contentPageSize,
      sourceType,
      sourceId: contentFilter.sourceId || undefined,
      locale: contentFilter.locale || undefined,
    });
    contents.value = result.rows;
    contentTotal.value = result.total;
    contentPage.value = page;
  } catch {
    errorText.value = t('common.loadFailed');
  } finally {
    loading.value = false;
  }
}

async function saveContent(row: TranslationContent): Promise<void> {
  loading.value = true;
  errorText.value = '';
  successText.value = '';
  try {
    await upsertTranslationContent({
      sourceType: row.sourceType,
      sourceId: row.sourceId,
      locale: row.locale,
      title: row.title,
      summary: row.summary,
      content: row.content,
      status: row.status,
      publishedAt: row.publishedAt || undefined,
      translatedByJobId: row.translatedByJobId || undefined,
    });
    successText.value = t('common.saveSuccess');
    await loadContents(contentPage.value);
  } catch {
    errorText.value = t('common.saveFailed');
  } finally {
    loading.value = false;
  }
}

onMounted(async () => {
  loading.value = true;
  errorText.value = '';
  try {
    await Promise.all([loadProviders(), loadJobs(1), loadContents(1)]);
  } catch {
    errorText.value = t('common.loadFailed');
  } finally {
    loading.value = false;
  }
});
</script>

<style scoped>
.system-page {
  display: grid;
  gap: 12px;
}

h1,
h2 {
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

.panel-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  gap: 10px;
}

.panel-header p {
  margin: 4px 0 0;
  color: var(--muted);
}

.switch-row {
  display: inline-flex;
  align-items: center;
  gap: 6px;
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

input,
select,
textarea,
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

.actions {
  display: flex;
  gap: 8px;
}

.pager {
  display: flex;
  justify-content: flex-end;
  align-items: center;
  gap: 8px;
}

.toolbar {
  display: flex;
  gap: 8px;
}

.grid-mini {
  display: grid;
  grid-template-columns: repeat(4, minmax(0, 1fr));
  gap: 8px;
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

.item-col {
  border: 1px solid var(--border);
  border-radius: 8px;
  padding: 10px;
  display: grid;
  gap: 8px;
}

.item-title,
.item-content,
.item-meta {
  margin: 0;
}

.item-title {
  font-weight: 600;
}

.item-meta {
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
  .grid-2,
  .grid-mini {
    grid-template-columns: 1fr;
  }

  .panel-header,
  .item {
    flex-direction: column;
  }

  .item-actions,
  .actions,
  .pager {
    width: 100%;
  }

  .item-actions button,
  .actions button,
  .pager button {
    flex: 1;
  }
}
</style>
