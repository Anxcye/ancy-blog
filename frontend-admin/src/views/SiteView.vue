<!--
File: SiteView.vue
Purpose: Manage site settings, footer items, social links, and navigation entries.
Module: frontend-admin/views/site, presentation layer.
Related: site API module and backend site/admin endpoints.
-->
<template>
  <section class="site-page">
    <h1>{{ t('site.title') }}</h1>
    <p class="subtitle">{{ t('site.subtitle') }}</p>

    <p v-if="errorText" class="error">{{ errorText }}</p>
    <p v-if="successText" class="success">{{ successText }}</p>

    <div class="panel">
      <h2>{{ t('site.settingsTitle') }}</h2>
      <div class="grid-2">
        <label>
          <span>{{ t('site.siteName') }}</span>
          <input v-model.trim="settingsForm.siteName" type="text" />
        </label>
        <label>
          <span>{{ t('site.defaultLocale') }}</span>
          <input v-model.trim="settingsForm.defaultLocale" type="text" placeholder="zh-CN" />
        </label>
        <label>
          <span>{{ t('site.avatarUrl') }}</span>
          <input v-model.trim="settingsForm.avatarUrl" type="url" />
        </label>
        <label>
          <span>{{ t('site.heroIntroMd') }}</span>
          <textarea v-model="settingsForm.heroIntroMd" rows="4" />
        </label>
      </div>
      <div class="actions">
        <button class="primary" :disabled="loading" @click="saveSettings">{{ t('common.save') }}</button>
      </div>
    </div>

    <div class="panel">
      <h2>{{ t('site.footerTitle') }}</h2>
      <div class="inline-form">
        <input v-model.trim="footerForm.label" :placeholder="t('site.footerLabel')" type="text" />
        <select v-model="footerForm.linkType">
          <option value="none">none</option>
          <option value="internal">internal</option>
          <option value="external">external</option>
        </select>
        <input v-model.trim="footerForm.internalArticleSlug" placeholder="about-me" type="text" />
        <input v-model.trim="footerForm.externalUrl" placeholder="https://example.com" type="url" />
        <input v-model.number="footerForm.rowNum" type="number" min="1" max="3" />
        <input v-model.number="footerForm.orderNum" type="number" min="1" />
        <button :disabled="loading" @click="createFooter">{{ t('common.add') }}</button>
      </div>
      <ul class="list">
        <li v-for="item in footerItems" :key="item.id">
          <span>{{ item.label }} (row {{ item.rowNum }}, #{{ item.orderNum }})</span>
          <code>{{ item.linkType }}</code>
        </li>
      </ul>
    </div>

    <div class="panel">
      <h2>{{ t('site.socialTitle') }}</h2>
      <div class="inline-form">
        <input v-model.trim="socialForm.platform" placeholder="github" type="text" />
        <input v-model.trim="socialForm.title" :placeholder="t('site.socialTitleField')" type="text" />
        <input v-model.trim="socialForm.url" placeholder="https://github.com/" type="url" />
        <input v-model.trim="socialForm.iconKey" placeholder="github" type="text" />
        <input v-model.number="socialForm.orderNum" type="number" min="1" />
        <button :disabled="loading" @click="createSocial">{{ t('common.add') }}</button>
      </div>
      <ul class="list">
        <li v-for="item in socialLinks" :key="item.id">
          <span>{{ item.title }}</span>
          <code>{{ item.url }}</code>
        </li>
      </ul>
    </div>

    <div class="panel">
      <h2>{{ t('site.navTitle') }}</h2>
      <div class="inline-form">
        <input v-model.trim="navForm.name" :placeholder="t('site.navName')" type="text" />
        <input v-model.trim="navForm.key" placeholder="articles" type="text" />
        <input v-model.trim="navForm.type" placeholder="menu" type="text" />
        <input v-model.trim="navForm.targetType" placeholder="route" type="text" />
        <input v-model.trim="navForm.targetValue" placeholder="/archives" type="text" />
        <input v-model.number="navForm.orderNum" type="number" min="1" />
        <button :disabled="loading" @click="createNav">{{ t('common.add') }}</button>
      </div>
      <ul class="list">
        <li v-for="item in navItems" :key="item.id">
          <span>{{ item.name }} ({{ item.key }})</span>
          <code>{{ item.targetType }}: {{ item.targetValue }}</code>
        </li>
      </ul>
    </div>
  </section>
</template>

<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue';
import { useI18n } from 'vue-i18n';

import {
  createFooterItem,
  createNavItem,
  createSocialLink,
  getSiteSettings,
  listFooterItems,
  listNavItems,
  listSocialLinks,
  updateSiteSettings,
} from '@/api/modules/site';
import type { FooterItem, NavItem, SiteSettings, SocialLink } from '@/api/types';

const { t } = useI18n();

const loading = ref(false);
const errorText = ref('');
const successText = ref('');

const settingsForm = reactive<SiteSettings>({
  siteName: '',
  avatarUrl: '',
  heroIntroMd: '',
  defaultLocale: 'zh-CN',
});

const footerForm = reactive<Omit<FooterItem, 'id'>>({
  label: '',
  linkType: 'none',
  internalArticleSlug: '',
  externalUrl: '',
  rowNum: 1,
  orderNum: 1,
  enabled: true,
});

const socialForm = reactive<Omit<SocialLink, 'id'>>({
  platform: '',
  title: '',
  url: '',
  iconKey: '',
  orderNum: 1,
  enabled: true,
});

const navForm = reactive<Omit<NavItem, 'id'>>({
  name: '',
  key: '',
  type: 'menu',
  targetType: 'route',
  targetValue: '',
  orderNum: 1,
  enabled: true,
});

const footerItems = ref<FooterItem[]>([]);
const socialLinks = ref<SocialLink[]>([]);
const navItems = ref<NavItem[]>([]);

async function loadAll(): Promise<void> {
  loading.value = true;
  errorText.value = '';
  try {
    const [settings, footer, socials, navs] = await Promise.all([
      getSiteSettings(),
      listFooterItems(),
      listSocialLinks(),
      listNavItems(),
    ]);
    settingsForm.siteName = settings.siteName;
    settingsForm.avatarUrl = settings.avatarUrl || '';
    settingsForm.heroIntroMd = settings.heroIntroMd || '';
    settingsForm.defaultLocale = settings.defaultLocale || 'zh-CN';

    footerItems.value = footer;
    socialLinks.value = socials;
    navItems.value = navs;
  } catch {
    errorText.value = t('common.loadFailed');
  } finally {
    loading.value = false;
  }
}

async function saveSettings(): Promise<void> {
  loading.value = true;
  errorText.value = '';
  successText.value = '';
  try {
    await updateSiteSettings(settingsForm);
    successText.value = t('common.saveSuccess');
  } catch {
    errorText.value = t('common.saveFailed');
  } finally {
    loading.value = false;
  }
}

async function createFooter(): Promise<void> {
  loading.value = true;
  errorText.value = '';
  successText.value = '';
  try {
    await createFooterItem(footerForm);
    successText.value = t('common.saveSuccess');
    footerForm.label = '';
    footerForm.internalArticleSlug = '';
    footerForm.externalUrl = '';
    await loadAll();
  } catch {
    errorText.value = t('common.saveFailed');
    loading.value = false;
  }
}

async function createSocial(): Promise<void> {
  loading.value = true;
  errorText.value = '';
  successText.value = '';
  try {
    await createSocialLink(socialForm);
    successText.value = t('common.saveSuccess');
    socialForm.platform = '';
    socialForm.title = '';
    socialForm.url = '';
    socialForm.iconKey = '';
    await loadAll();
  } catch {
    errorText.value = t('common.saveFailed');
    loading.value = false;
  }
}

async function createNav(): Promise<void> {
  loading.value = true;
  errorText.value = '';
  successText.value = '';
  try {
    await createNavItem(navForm);
    successText.value = t('common.saveSuccess');
    navForm.name = '';
    navForm.key = '';
    navForm.targetValue = '';
    await loadAll();
  } catch {
    errorText.value = t('common.saveFailed');
    loading.value = false;
  }
}

onMounted(async () => {
  await loadAll();
});
</script>

<style scoped>
.site-page {
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

.grid-2 {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 10px;
}

.inline-form {
  display: grid;
  grid-template-columns: repeat(6, minmax(0, 1fr));
  gap: 8px;
}

label {
  display: grid;
  gap: 6px;
}

input,
textarea,
select,
button {
  border: 1px solid var(--border);
  border-radius: 8px;
  padding: 8px 10px;
  font: inherit;
}

button {
  cursor: pointer;
  background: var(--surface);
}

button.primary {
  background: var(--accent);
  color: #fff;
  border-color: var(--accent);
}

.actions {
  display: flex;
  justify-content: flex-end;
}

.list {
  list-style: none;
  margin: 0;
  padding: 0;
  display: grid;
  gap: 8px;
}

.list li {
  border: 1px solid var(--border);
  border-radius: 8px;
  padding: 8px 10px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 8px;
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
  .grid-2,
  .inline-form {
    grid-template-columns: 1fr;
  }
}
</style>
