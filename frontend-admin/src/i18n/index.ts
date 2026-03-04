// File: i18n/index.ts
// Purpose: Register admin UI localization resources with fixed Chinese locale.
// Module: frontend-admin/i18n, presentation localization layer.
// Related: locale dictionaries and UI components using translation keys.
import { createI18n } from 'vue-i18n';

import zhCN from '@/locales/zh-CN';

const i18n = createI18n({
  legacy: false,
  locale: 'zh-CN',
  fallbackLocale: 'zh-CN',
  messages: {
    'zh-CN': zhCN,
  },
});

export default i18n;
