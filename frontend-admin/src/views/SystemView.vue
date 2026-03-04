<!--
File: SystemView.vue
Purpose: Manage external integration providers and run connectivity tests.
Module: frontend-admin/views/system, presentation layer.
Related: integration API module and backend integration endpoints.
-->
<template>
  <section class="system-page">
    <h1>{{ t('system.title') }}</h1>
    <p class="subtitle">{{ t('system.subtitle') }}</p>

    <p v-if="errorText" class="error">{{ errorText }}</p>
    <p v-if="successText" class="success">{{ successText }}</p>

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
  </section>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue';
import { useI18n } from 'vue-i18n';

import { listIntegrations, testIntegration, updateIntegration } from '@/api/modules/integrations';

type IntegrationViewModel = {
  providerKey: string;
  providerType: string;
  name: string;
  enabled: boolean;
  configRaw: string;
  metaRaw: string;
};

const { t } = useI18n();

const loading = ref(false);
const errorText = ref('');
const successText = ref('');
const providers = ref<IntegrationViewModel[]>([]);

async function loadProviders(): Promise<void> {
  loading.value = true;
  errorText.value = '';
  try {
    const rows = await listIntegrations();
    providers.value = rows.map((item) => ({
      providerKey: item.providerKey,
      providerType: item.providerType,
      name: item.name,
      enabled: item.enabled,
      configRaw: JSON.stringify(item.configJson ?? {}, null, 2),
      metaRaw: JSON.stringify(item.metaJson ?? {}, null, 2),
    }));
  } catch {
    errorText.value = t('common.loadFailed');
  } finally {
    loading.value = false;
  }
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

onMounted(async () => {
  await loadProviders();
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

label {
  display: grid;
  gap: 6px;
}

textarea,
button {
  border: 1px solid var(--border);
  border-radius: 8px;
  padding: 8px 10px;
  font: inherit;
}

.actions {
  display: flex;
  gap: 8px;
}

button {
  cursor: pointer;
  background: var(--surface);
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
  .panel-header {
    flex-direction: column;
  }

  .actions {
    width: 100%;
  }

  .actions button {
    flex: 1;
  }
}
</style>
