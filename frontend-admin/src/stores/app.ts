// File: stores/app.ts
// Purpose: Hold admin session token, color theme mode, and locale preference state.
// Module: frontend-admin/stores, state management layer.
// Related: auth API module, router auth guard, i18n runtime, app shell controls.
import { computed, ref } from 'vue';
import { defineStore } from 'pinia';

const TOKEN_KEY = 'ancy_admin_token';
const THEME_KEY = 'ancy_admin_theme';
const LOCALE_KEY = 'ancy_admin_locale';

type ThemeMode = 'light' | 'dark';
type LocaleCode = 'zh-CN' | 'en-US';

function readThemeMode(): ThemeMode {
  const mode = localStorage.getItem(THEME_KEY);
  return mode === 'dark' ? 'dark' : 'light';
}

function readLocaleCode(): LocaleCode {
  const locale = localStorage.getItem(LOCALE_KEY);
  return locale === 'en-US' ? 'en-US' : 'zh-CN';
}

export const useAppStore = defineStore('app', () => {
  const token = ref<string | null>(localStorage.getItem(TOKEN_KEY));
  const themeMode = ref<ThemeMode>(readThemeMode());
  const locale = ref<LocaleCode>(readLocaleCode());

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

  function toggleThemeMode(): void {
    setThemeMode(themeMode.value === 'dark' ? 'light' : 'dark');
  }

  function setLocale(nextLocale: LocaleCode): void {
    locale.value = nextLocale;
    localStorage.setItem(LOCALE_KEY, nextLocale);
  }

  return {
    token,
    isAuthenticated,
    themeMode,
    locale,
    setToken,
    clearToken,
    setThemeMode,
    toggleThemeMode,
    setLocale,
  };
});
