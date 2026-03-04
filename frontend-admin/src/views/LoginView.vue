<!--
File: LoginView.vue
Purpose: Render admin login form and execute authentication flow.
Module: frontend-admin/views/auth, presentation layer.
Related: auth API module, app store token state, router navigation.
-->
<template>
  <section class="login-card">
    <h1>{{ t('login.title') }}</h1>
    <form @submit.prevent="onSubmit">
      <label>
        <span>{{ t('login.username') }}</span>
        <input v-model.trim="username" type="text" required autocomplete="username" />
      </label>
      <label>
        <span>{{ t('login.password') }}</span>
        <input v-model="password" type="password" required autocomplete="current-password" />
      </label>
      <button :disabled="submitting" type="submit">
        {{ t('login.submit') }}
      </button>
      <p v-if="errorText" class="error">{{ errorText }}</p>
    </form>
  </section>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import { useRouter } from 'vue-router';
import { useI18n } from 'vue-i18n';

import { login } from '@/api/modules/auth';
import { useAppStore } from '@/stores/app';

const router = useRouter();
const appStore = useAppStore();
const { t } = useI18n();

const username = ref('');
const password = ref('');
const errorText = ref('');
const submitting = ref(false);

async function onSubmit(): Promise<void> {
  errorText.value = '';
  submitting.value = true;

  try {
    const token = await login(username.value, password.value);
    appStore.setToken(token);
    await router.push({ name: 'dashboard' });
  } catch {
    errorText.value = t('login.failed');
  } finally {
    submitting.value = false;
  }
}
</script>

<style scoped>
.login-card {
  width: min(420px, 100%);
  margin: 64px auto;
  padding: 24px;
  border: 1px solid var(--border);
  border-radius: 12px;
  background: var(--surface);
}

h1 {
  margin: 0 0 20px;
}

form {
  display: grid;
  gap: 14px;
}

label {
  display: grid;
  gap: 6px;
}

input {
  width: 100%;
  padding: 10px;
  border-radius: 8px;
  border: 1px solid var(--border);
}

button {
  padding: 10px;
  border: 0;
  border-radius: 8px;
  background: var(--accent);
  color: #fff;
  cursor: pointer;
}

button:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.error {
  margin: 0;
  color: #b64040;
}
</style>
