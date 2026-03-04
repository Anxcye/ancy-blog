// File: stores/app.ts
// Purpose: Hold admin session token and theme mode state.
// Module: frontend-admin/stores, state management layer.
// Related: auth API module, router auth guard, app shell controls.
import { computed, ref } from 'vue';
import { defineStore } from 'pinia';

const TOKEN_KEY = 'ancy_admin_token';
const THEME_KEY = 'ancy_admin_theme';

type ThemeMode = 'light' | 'dark';

function readThemeMode(): ThemeMode {
  const mode = localStorage.getItem(THEME_KEY);
  return mode === 'dark' ? 'dark' : 'light';
}

export const useAppStore = defineStore('app', () => {
  const token = ref<string | null>(localStorage.getItem(TOKEN_KEY));
  const themeMode = ref<ThemeMode>(readThemeMode());

  const isAuthenticated = computed(() => Boolean(token.value));

  function setToken(value: string): void {
    token.value = value;
    localStorage.setItem(TOKEN_KEY, value);
  }

  function clearToken(): void {
    token.value = null;
    localStorage.removeItem(TOKEN_KEY);
  }

  function setThemeMode(mode: ThemeMode): void {
    themeMode.value = mode;
    localStorage.setItem(THEME_KEY, mode);
  }

  return {
    token,
    isAuthenticated,
    themeMode,
    setToken,
    clearToken,
    setThemeMode,
  };
});
