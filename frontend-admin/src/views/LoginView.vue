<!--
File: LoginView.vue
Purpose: Render admin sign-in form with enterprise-grade validation and feedback.
Module: frontend-admin/views/auth, presentation layer.
Related: auth API module, app store token state, router navigation, global message provider.
-->
<template>
  <section class="login-page">
    <NCard class="login-card" :bordered="false">
      <header class="login-header">
        <h1>{{ t('login.title') }}</h1>
        <p>{{ t('layout.subtitle') }}</p>
      </header>

      <NForm label-placement="top" size="large" @submit.prevent="onSubmit">
        <NFormItem :label="t('login.username')">
          <NInput v-model:value="username" :placeholder="t('login.username')" clearable autocomplete="username" />
        </NFormItem>

        <NFormItem :label="t('login.password')">
          <NInput
            v-model:value="password"
            type="password"
            show-password-on="click"
            :placeholder="t('login.password')"
            autocomplete="current-password"
          />
        </NFormItem>

        <NAlert v-if="errorText" type="error" :show-icon="false">{{ errorText }}</NAlert>

        <NButton type="primary" block attr-type="submit" :loading="submitting">{{ t('login.submit') }}</NButton>
      </NForm>
    </NCard>
  </section>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import { useRouter } from 'vue-router';
import { useI18n } from 'vue-i18n';
import { NAlert, NButton, NCard, NForm, NFormItem, NInput, useMessage } from 'naive-ui';

import { login } from '@/api/modules/auth';
import { useAppStore } from '@/stores/app';

const router = useRouter();
const appStore = useAppStore();
const { t } = useI18n();
const message = useMessage();

const username = ref('');
const password = ref('');
const errorText = ref('');
const submitting = ref(false);

async function onSubmit(): Promise<void> {
  errorText.value = '';
  if (!username.value.trim() || !password.value) {
    errorText.value = t('login.failed');
    return;
  }
  submitting.value = true;
  try {
    const token = await login(username.value.trim(), password.value);
    appStore.setToken(token);
    message.success(t('login.welcome'));
    await router.push({ name: 'dashboard' });
  } catch {
    errorText.value = t('login.failed');
  } finally {
    submitting.value = false;
  }
}
</script>

<style scoped>
.login-page {
  width: min(420px, 100%);
}

.login-card {
  border-radius: 20px;
  box-shadow: 0 22px 48px rgba(17, 23, 31, 0.09);
}

.login-header {
  margin-bottom: 12px;
}

.login-header h1 {
  margin: 0;
  font-size: 26px;
}

.login-header p {
  margin: 8px 0 0;
  color: #6c7780;
}
</style>
