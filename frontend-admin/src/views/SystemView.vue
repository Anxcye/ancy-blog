<!--
File: SystemView.vue
Purpose: Manage integrations and translation jobs/content operations with unified Naive UI components.
Module: frontend-admin/views/system, presentation layer.
Related: integrations API module, translations API module, and backend admin system endpoints.
-->
<template>
  <section class="system-page">
    <NCard :bordered="false" class="section-card">
      <div class="topbar">
        <div class="tab-switch">
          <NButton :type="tab === 'integrations' ? 'primary' : 'default'" secondary @click="tab = 'integrations'">{{ t('system.tabIntegrations') }}</NButton>
          <NButton :type="tab === 'translations' ? 'primary' : 'default'" secondary @click="tab = 'translations'">{{ t('system.tabTranslations') }}</NButton>
        </div>
      </div>

      <NAlert v-if="errorText" type="error" :show-icon="false">{{ errorText }}</NAlert>
      <NAlert v-if="successText" type="success" :show-icon="false">{{ successText }}</NAlert>

      <template v-if="tab === 'integrations'">
        <NCard v-for="provider in providers" :key="provider.providerKey" size="small" :bordered="true" class="provider-card">
          <template #header>
            <div class="provider-header">
              <div>
                <strong>{{ provider.name }}</strong>
                <p>{{ provider.providerKey }} · {{ provider.providerType }}</p>
              </div>
              <NSwitch v-model:value="provider.enabled" />
            </div>
          </template>

          <NForm label-placement="top">
            <div v-if="provider.providerKey === 'openai_compatible'" class="grid-3">
              <NFormItem label="base_url">
                <NInput :value="getConfigField(provider, 'base_url')" placeholder="https://api.openai.com/v1" @update:value="(v) => setConfigField(provider, 'base_url', v)" />
              </NFormItem>
              <NFormItem label="api_key">
                <NInput type="password" show-password-on="click" :value="getConfigField(provider, 'api_key')" placeholder="sk-..." @update:value="(v) => setConfigField(provider, 'api_key', v)" />
              </NFormItem>
              <NFormItem label="model">
                <NInput :value="getConfigField(provider, 'model')" placeholder="gpt-4.1-mini" @update:value="(v) => setConfigField(provider, 'model', v)" />
              </NFormItem>
            </div>

            <div v-else-if="provider.providerKey === 'cloudflare_r2'" class="grid-3">
              <NFormItem label="account_id">
                <NInput :value="getConfigField(provider, 'account_id')" @update:value="(v) => setConfigField(provider, 'account_id', v)" />
              </NFormItem>
              <NFormItem label="bucket">
                <NInput :value="getConfigField(provider, 'bucket')" @update:value="(v) => setConfigField(provider, 'bucket', v)" />
              </NFormItem>
              <NFormItem label="public_base_url">
                <NInput :value="getConfigField(provider, 'public_base_url')" placeholder="https://cdn.example.com" @update:value="(v) => setConfigField(provider, 'public_base_url', v)" />
              </NFormItem>
              <NFormItem label="access_key_id">
                <NInput :value="getConfigField(provider, 'access_key_id')" @update:value="(v) => setConfigField(provider, 'access_key_id', v)" />
              </NFormItem>
              <NFormItem label="secret_access_key">
                <NInput type="password" show-password-on="click" :value="getConfigField(provider, 'secret_access_key')" @update:value="(v) => setConfigField(provider, 'secret_access_key', v)" />
              </NFormItem>
            </div>

            <NCollapse>
              <NCollapseItem title="高级 JSON（可选）" name="advanced">
                <NFormItem :label="t('system.configJson')">
                  <NInput v-model:value="provider.configRaw" type="textarea" :autosize="{ minRows: 4, maxRows: 10 }" />
                </NFormItem>
                <NFormItem :label="t('system.metaJson')">
                  <NInput v-model:value="provider.metaRaw" type="textarea" :autosize="{ minRows: 3, maxRows: 8 }" />
                </NFormItem>
              </NCollapseItem>
            </NCollapse>
          </NForm>

          <NSpace>
            <NButton :loading="loading" @click="saveProvider(provider)">{{ t('common.save') }}</NButton>
            <NButton :loading="loading" @click="testProvider(provider.providerKey)">{{ t('system.test') }}</NButton>
          </NSpace>
        </NCard>
      </template>

      <template v-else>
        <NCard size="small" :bordered="true">
          <template #header>{{ t('system.translationCreateTitle') }}</template>
          <NForm label-placement="top">
            <div class="grid-3">
              <NFormItem label="sourceType">
                <NSelect v-model:value="jobForm.sourceType" :options="sourceTypeOptions" />
              </NFormItem>
              <NFormItem label="sourceId">
                <NInput v-model:value="jobForm.sourceId" />
              </NFormItem>
              <NFormItem label="sourceLocale">
                <NInput v-model:value="jobForm.sourceLocale" placeholder="zh-CN" />
              </NFormItem>
              <NFormItem label="targetLocale">
                <NInput v-model:value="jobForm.targetLocale" placeholder="en-US" />
              </NFormItem>
              <NFormItem label="providerKey">
                <NInput v-model:value="jobForm.providerKey" placeholder="openai_compatible" />
              </NFormItem>
              <NFormItem label="modelName">
                <NInput v-model:value="jobForm.modelName" placeholder="gpt-4.1-mini" />
              </NFormItem>
              <NFormItem label="maxRetries">
                <NInputNumber v-model:value="jobForm.maxRetries" :min="0" :max="10" style="width: 100%" />
              </NFormItem>
              <NFormItem label="publishAt">
                <NDatePicker v-model:value="jobPublishAtTs" type="datetime" clearable style="width: 100%" />
              </NFormItem>
              <NFormItem label="autoPublish">
                <NSwitch v-model:value="jobForm.autoPublish" />
              </NFormItem>
            </div>
          </NForm>
          <NButton type="primary" :loading="loading" @click="createJob">{{ t('system.createTranslationJob') }}</NButton>
        </NCard>

        <NCard size="small" :bordered="true">
          <template #header>{{ t('system.translationJobsTitle') }}</template>
          <template #header-extra>
            <NSelect
              v-model:value="jobStatusFilter"
              clearable
              :options="jobStatusOptions"
              style="width: 160px"
              @update:value="() => loadJobs(1)"
            />
          </template>

          <NDataTable remote :loading="loading" :columns="jobColumns" :data="jobs" :pagination="false" :row-key="rowKey" />

          <div class="footer-row">
            <span class="hint">{{ t('articles.total', { total: jobTotal }) }}</span>
            <NPagination :page="jobPage" :page-size="jobPageSize" :item-count="jobTotal" :page-slot="7" @update:page="loadJobs" />
          </div>
        </NCard>

        <NCard size="small" :bordered="true">
          <template #header>{{ t('system.translationContentsTitle') }}</template>
          <NForm inline :show-label="false" class="filters" @submit.prevent="loadContents(1)">
            <NFormItem>
              <NSelect v-model:value="contentFilter.sourceType" :options="sourceTypeOptions" style="width: 130px" />
            </NFormItem>
            <NFormItem>
              <NInput v-model:value="contentFilter.sourceId" placeholder="sourceId" />
            </NFormItem>
            <NFormItem>
              <NInput v-model:value="contentFilter.locale" placeholder="locale" />
            </NFormItem>
            <NFormItem>
              <NButton attr-type="submit">{{ t('common.search') }}</NButton>
            </NFormItem>
          </NForm>

          <NDataTable remote :loading="loading" :columns="contentColumns" :data="contents" :pagination="false" :row-key="contentRowKey" />

          <div class="footer-row">
            <span class="hint">{{ t('articles.total', { total: contentTotal }) }}</span>
            <NPagination :page="contentPage" :page-size="contentPageSize" :item-count="contentTotal" :page-slot="7" @update:page="loadContents" />
          </div>
        </NCard>
      </template>
    </NCard>
  </section>
</template>

<script setup lang="ts">
import { computed, h, onMounted, reactive, ref } from 'vue';
import { useRoute } from 'vue-router';
import { useI18n } from 'vue-i18n';
import type { DataTableColumns } from 'naive-ui';
import { NAlert, NButton, NCard, NCollapse, NCollapseItem, NDataTable, NDatePicker, NForm, NFormItem, NInput, NInputNumber, NPagination, NSelect, NSpace, NSwitch } from 'naive-ui';

import { listIntegrations, testIntegration, updateIntegration } from '@/api/modules/integrations';
import { createTranslationJob, listTranslationContents, listTranslationJobs, retryTranslationJob, upsertTranslationContent } from '@/api/modules/translations';
import type { TranslationContent, TranslationJob } from '@/api/types';

type IntegrationViewModel = {
  providerKey: string;
  providerType: string;
  name: string;
  enabled: boolean;
  configJson: Record<string, unknown>;
  metaJson: Record<string, unknown>;
  configRaw: string;
  metaRaw: string;
};

type TabName = 'integrations' | 'translations';

const { t } = useI18n();
const route = useRoute();

const tab = ref<TabName>('integrations');
const loading = ref(false);
const errorText = ref('');
const successText = ref('');
const providers = ref<IntegrationViewModel[]>([]);

const jobStatusFilter = ref<string | null>(null);
const jobs = ref<TranslationJob[]>([]);
const jobPage = ref(1);
const jobTotal = ref(0);
const jobPageSize = 20;
const jobPublishAtTs = ref<number | null>(null);
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

const sourceTypeOptions = [
  { label: 'article', value: 'article' },
  { label: 'moment', value: 'moment' },
];

const jobStatusOptions = [
  { label: 'queued', value: 'queued' },
  { label: 'running', value: 'running' },
  { label: 'succeeded', value: 'succeeded' },
  { label: 'failed', value: 'failed' },
];

function rowKey(row: TranslationJob): string {
  return row.id;
}

function contentRowKey(row: TranslationContent): string {
  return `${row.sourceType}:${row.sourceId}:${row.locale}`;
}

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

function toPrettyJSON(value: Record<string, unknown>): string {
  return JSON.stringify(value ?? {}, null, 2);
}

function getConfigField(provider: IntegrationViewModel, key: string): string {
  return String(provider.configJson[key] ?? '');
}

function setConfigField(provider: IntegrationViewModel, key: string, value: string): void {
  provider.configJson[key] = value;
  provider.configRaw = toPrettyJSON(provider.configJson);
}

const jobColumns = computed<DataTableColumns<TranslationJob>>(() => [
  {
    title: '任务ID',
    key: 'id',
    ellipsis: {
      tooltip: true,
    },
  },
  {
    title: '状态',
    key: 'status',
    width: 120,
  },
  {
    title: '源/目标',
    key: 'pair',
    width: 220,
    render(row) {
      return `${row.sourceType} -> ${row.targetLocale}`;
    },
  },
  {
    title: '重试',
    key: 'retry',
    width: 120,
    render(row) {
      return `${row.retryCount}/${row.maxRetries}`;
    },
  },
  {
    title: '更新时间',
    key: 'updatedAt',
    width: 180,
    render(row) {
      return formatDate(row.updatedAt);
    },
  },
  {
    title: '操作',
    key: 'actions',
    width: 110,
    render(row) {
      return h(
        NButton,
        {
          size: 'small',
          tertiary: true,
          disabled: loading.value || row.status !== 'failed',
          onClick: () => retryJob(row.id),
        },
        { default: () => t('system.retry') },
      );
    },
  },
]);

const contentColumns = computed<DataTableColumns<TranslationContent>>(() => [
  {
    title: '来源',
    key: 'origin',
    width: 260,
    render(row) {
      return `${row.sourceType} · ${row.sourceId} · ${row.locale}`;
    },
  },
  {
    title: t('editor.fieldTitle'),
    key: 'title',
    render(row) {
      return h(NInput, {
        value: row.title,
        onUpdateValue: (value: string) => {
          row.title = value;
        },
      });
    },
  },
  {
    title: t('editor.fieldSummary'),
    key: 'summary',
    render(row) {
      return h(NInput, {
        value: row.summary,
        onUpdateValue: (value: string) => {
          row.summary = value;
        },
      });
    },
  },
  {
    title: '操作',
    key: 'actions',
    width: 110,
    render(row) {
      return h(
        NButton,
        {
          size: 'small',
          tertiary: true,
          disabled: loading.value,
          onClick: () => saveContent(row),
        },
        { default: () => t('common.save') },
      );
    },
  },
]);

async function loadProviders(): Promise<void> {
  const rows = await listIntegrations();
  providers.value = rows.map((item) => ({
    providerKey: item.providerKey,
    providerType: item.providerType,
    name: item.name,
    enabled: item.enabled,
    configJson: { ...(item.configJson ?? {}) },
    metaJson: { ...(item.metaJson ?? {}) },
    configRaw: toPrettyJSON(item.configJson ?? {}),
    metaRaw: toPrettyJSON(item.metaJson ?? {}),
  }));
}

async function saveProvider(item: IntegrationViewModel): Promise<void> {
  loading.value = true;
  errorText.value = '';
  successText.value = '';
  try {
    item.configJson = parseJSON(item.configRaw);
    item.metaJson = parseJSON(item.metaRaw);
    await updateIntegration(item.providerKey, {
      enabled: item.enabled,
      configJson: item.configJson,
      metaJson: item.metaJson,
    });
    item.configRaw = toPrettyJSON(item.configJson);
    item.metaRaw = toPrettyJSON(item.metaJson);
    const result = await testIntegration(item.providerKey);
    if (result.ok) {
      successText.value = `${t('common.saveSuccess')} ${t('system.autoTestOk')} ${result.message}`;
    } else {
      successText.value = `${t('common.saveSuccess')} ${t('system.autoTestFail')} ${result.message}`;
    }
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
    successText.value = `${result.ok ? t('system.testOk') : t('system.testFail')} · ${result.message}`;
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
    jobForm.publishAt = jobPublishAtTs.value ? new Date(jobPublishAtTs.value).toISOString() : '';
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
  const queryTab = String(route.query.tab || '');
  const querySourceType = String(route.query.sourceType || '');
  const querySourceID = String(route.query.sourceId || '');
  if (queryTab === 'translations') {
    tab.value = 'translations';
  }
  if (querySourceType === 'article' || querySourceType === 'moment') {
    jobForm.sourceType = querySourceType;
    contentFilter.sourceType = querySourceType;
  }
  if (querySourceID) {
    jobForm.sourceId = querySourceID;
    contentFilter.sourceId = querySourceID;
  }

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
}

.section-card {
  border-radius: 14px;
  box-shadow: 0 6px 24px color-mix(in srgb, var(--n-text-color) 8%, transparent);
  display: grid;
  gap: 12px;
}

.topbar {
  display: flex;
  justify-content: flex-start;
  align-items: center;
  gap: 12px;
}

.tab-switch {
  display: flex;
  gap: 8px;
}

.provider-card + .provider-card {
  margin-top: 10px;
}

.provider-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 12px;
}

.provider-header p {
  margin: 4px 0 0;
  color: var(--n-text-color-3);
}

.grid-3 {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 10px;
}

.filters {
  margin-bottom: 10px;
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
  .grid-3,
  .footer-row,
  .topbar {
    grid-template-columns: 1fr;
    flex-direction: column;
    align-items: stretch;
  }
}
</style>
