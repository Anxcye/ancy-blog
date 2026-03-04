// File: stores/app.ts
// Purpose: Hold admin session token state and basic app-level flags.
// Module: frontend-admin/stores, state management layer.
// Related: auth API module, router auth guard, login view.
import { computed, ref } from 'vue';
import { defineStore } from 'pinia';

const TOKEN_KEY = 'ancy_admin_token';

export const useAppStore = defineStore('app', () => {
  const token = ref<string | null>(localStorage.getItem(TOKEN_KEY));

  const isAuthenticated = computed(() => Boolean(token.value));

  function setToken(value: string): void {
    token.value = value;
    localStorage.setItem(TOKEN_KEY, value);
  }

  function clearToken(): void {
    token.value = null;
    localStorage.removeItem(TOKEN_KEY);
  }

  return {
    token,
    isAuthenticated,
    setToken,
    clearToken,
  };
});
