// File: i18n/index.ts
// Purpose: Register admin UI localization resources and initialize locale preference.
// Module: frontend-admin/i18n, presentation localization layer.
// Related: locale dictionaries, app store locale state, UI components using translation keys.
import { createI18n } from 'vue-i18n';

import zhCN from '@/locales/zh-CN';
import enUS from '@/locales/en-US';

const storedLocale = localStorage.getItem('ancy_admin_locale');
const defaultLocale = storedLocale === 'en-US' ? 'en-US' : 'zh-CN';

const i18n = createI18n({
  legacy: false,
  locale: defaultLocale,
  fallbackLocale: 'en-US',
  messages: {
    'zh-CN': zhCN,
    'en-US': enUS,
  },
});

export default i18n;
