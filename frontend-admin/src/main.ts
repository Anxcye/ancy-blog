// File: main.ts
// Purpose: Bootstrap Vue admin application and register core plugins.
// Module: frontend-admin/app, entry layer.
// Related: router, pinia store, i18n setup, App root component.
import { createApp } from 'vue';
import { createPinia } from 'pinia';

import App from './App.vue';
import router from './router';
import i18n from './i18n';
import './styles/tokens.css';

const app = createApp(App);

app.use(createPinia());
app.use(router);
app.use(i18n);
app.mount('#app');
