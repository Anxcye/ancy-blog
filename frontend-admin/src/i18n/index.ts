// File: i18n/index.ts
// Purpose: Register admin UI localization resources and default locale.
// Module: frontend-admin/i18n, presentation localization layer.
// Related: locale dictionaries, UI labels, route-visible text.
import { createI18n } from 'vue-i18n';

import zhCN from '@/locales/zh-CN';
import enUS from '@/locales/en-US';

const i18n = createI18n({
  legacy: false,
  locale: 'zh-CN',
  fallbackLocale: 'en-US',
  messages: {
    'zh-CN': zhCN,
    'en-US': enUS,
  },
});

export default i18n;
