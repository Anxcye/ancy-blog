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
        <button :disabled="loading" @click="submitFooter">
          {{ editingFooterId ? t('common.update') : t('common.add') }}
        </button>
        <button v-if="editingFooterId" :disabled="loading" @click="cancelFooterEdit">{{ t('common.cancel') }}</button>
      </div>
      <ul class="list">
        <li v-for="item in footerItems" :key="item.id">
          <span>{{ item.label }} (row {{ item.rowNum }}, #{{ item.orderNum }})</span>
          <div class="row-actions">
            <code>{{ item.linkType }}</code>
            <button :disabled="loading" @click="startFooterEdit(item)">{{ t('common.edit') }}</button>
            <button :disabled="loading" class="danger" @click="removeFooter(item.id)">{{ t('common.delete') }}</button>
          </div>
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
        <button :disabled="loading" @click="submitSocial">
          {{ editingSocialId ? t('common.update') : t('common.add') }}
        </button>
        <button v-if="editingSocialId" :disabled="loading" @click="cancelSocialEdit">{{ t('common.cancel') }}</button>
      </div>
      <ul class="list">
        <li v-for="item in socialLinks" :key="item.id">
          <span>{{ item.title }}</span>
          <div class="row-actions">
            <code>{{ item.url }}</code>
            <button :disabled="loading" @click="startSocialEdit(item)">{{ t('common.edit') }}</button>
            <button :disabled="loading" class="danger" @click="removeSocial(item.id)">{{ t('common.delete') }}</button>
          </div>
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
        <button :disabled="loading" @click="submitNav">
          {{ editingNavId ? t('common.update') : t('common.add') }}
        </button>
        <button v-if="editingNavId" :disabled="loading" @click="cancelNavEdit">{{ t('common.cancel') }}</button>
      </div>
      <ul class="list">
        <li v-for="item in navItems" :key="item.id">
          <span>{{ item.name }} ({{ item.key }})</span>
          <div class="row-actions">
            <code>{{ item.targetType }}: {{ item.targetValue }}</code>
            <button :disabled="loading" @click="startNavEdit(item)">{{ t('common.edit') }}</button>
            <button :disabled="loading" class="danger" @click="removeNav(item.id)">{{ t('common.delete') }}</button>
          </div>
        </li>
      </ul>
    </div>

    <div class="panel">
      <h2>{{ t('site.slotTitle') }}</h2>
      <div class="inline-form">
        <input v-model.trim="slotForm.slotKey" placeholder="home_featured" type="text" />
        <input v-model.trim="slotForm.name" :placeholder="t('site.slotName')" type="text" />
        <input v-model.trim="slotForm.description" :placeholder="t('site.slotDescription')" type="text" />
        <button :disabled="loading" @click="submitSlot">{{ t('common.add') }}</button>
      </div>

      <div class="slot-tools">
        <select v-model="selectedSlotKey" @change="onSlotChange">
          <option v-for="slot in slots" :key="slot.id" :value="slot.slotKey">{{ slot.slotKey }} · {{ slot.name }}</option>
        </select>
      </div>

      <div v-if="selectedSlotKey" class="inline-form">
        <select v-model="slotItemForm.contentType">
          <option value="article">article</option>
          <option value="moment">moment</option>
        </select>
        <input v-model.trim="slotItemForm.contentId" placeholder="content id (uuid)" type="text" />
        <input v-model.number="slotItemForm.orderNum" type="number" min="1" />
        <button :disabled="loading" @click="submitSlotItem">{{ t('common.add') }}</button>
      </div>

      <ul class="list">
        <li v-for="item in slotItems" :key="item.id">
          <span>{{ item.contentType }} · {{ item.contentId }} · #{{ item.orderNum }}</span>
          <div class="row-actions">
            <code>{{ item.enabled ? 'enabled' : 'disabled' }}</code>
            <button :disabled="loading" class="danger" @click="removeSlotItem(item.id)">{{ t('common.delete') }}</button>
          </div>
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
  createSlot,
  createSlotItem,
  createSocialLink,
  deleteFooterItem,
  deleteNavItem,
  deleteSlotItem,
  deleteSocialLink,
  getSiteSettings,
  listFooterItems,
  listNavItems,
  listSlotItems,
  listSlots,
  listSocialLinks,
  updateFooterItem,
  updateNavItem,
  updateSiteSettings,
  updateSocialLink,
} from '@/api/modules/site';
import type { ContentSlot, FooterItem, NavItem, SiteSettings, SlotItem, SocialLink } from '@/api/types';

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

const editingFooterId = ref('');
const footerForm = reactive<Omit<FooterItem, 'id'>>({
  label: '',
  linkType: 'none',
  internalArticleSlug: '',
  externalUrl: '',
  rowNum: 1,
  orderNum: 1,
  enabled: true,
});

const editingSocialId = ref('');
const socialForm = reactive<Omit<SocialLink, 'id'>>({
  platform: '',
  title: '',
  url: '',
  iconKey: '',
  orderNum: 1,
  enabled: true,
});

const editingNavId = ref('');
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
const slots = ref<ContentSlot[]>([]);
const slotItems = ref<SlotItem[]>([]);
const selectedSlotKey = ref('');

const slotForm = reactive({
  slotKey: '',
  name: '',
  description: '',
  enabled: true,
});

const slotItemForm = reactive({
  contentType: 'article',
  contentId: '',
  orderNum: 1,
  enabled: true,
});

function setSuccess(message: string): void {
  successText.value = message;
  errorText.value = '';
}

function resetFooterForm(): void {
  editingFooterId.value = '';
  footerForm.label = '';
  footerForm.linkType = 'none';
  footerForm.internalArticleSlug = '';
  footerForm.externalUrl = '';
  footerForm.rowNum = 1;
  footerForm.orderNum = 1;
  footerForm.enabled = true;
}

function resetSocialForm(): void {
  editingSocialId.value = '';
  socialForm.platform = '';
  socialForm.title = '';
  socialForm.url = '';
  socialForm.iconKey = '';
  socialForm.orderNum = 1;
  socialForm.enabled = true;
}

function resetNavForm(): void {
  editingNavId.value = '';
  navForm.name = '';
  navForm.key = '';
  navForm.type = 'menu';
  navForm.targetType = 'route';
  navForm.targetValue = '';
  navForm.orderNum = 1;
  navForm.enabled = true;
}

async function loadAll(): Promise<void> {
  loading.value = true;
  errorText.value = '';
  try {
    const [settings, footer, socials, navs, slotRows] = await Promise.all([
      getSiteSettings(),
      listFooterItems(),
      listSocialLinks(),
      listNavItems(),
      listSlots(),
    ]);
    settingsForm.siteName = settings.siteName;
    settingsForm.avatarUrl = settings.avatarUrl || '';
    settingsForm.heroIntroMd = settings.heroIntroMd || '';
    settingsForm.defaultLocale = settings.defaultLocale || 'zh-CN';
    footerItems.value = footer;
    socialLinks.value = socials;
    navItems.value = navs;
    slots.value = slotRows;
    if (!selectedSlotKey.value && slotRows.length > 0) {
      selectedSlotKey.value = slotRows[0].slotKey;
    }
    if (selectedSlotKey.value) {
      await loadSlotItems(selectedSlotKey.value);
    } else {
      slotItems.value = [];
    }
  } catch {
    errorText.value = t('common.loadFailed');
  } finally {
    loading.value = false;
  }
}

async function loadSlotItems(slotKey: string): Promise<void> {
  slotItems.value = await listSlotItems(slotKey);
}

async function onSlotChange(): Promise<void> {
  if (!selectedSlotKey.value) {
    slotItems.value = [];
    return;
  }
  loading.value = true;
  errorText.value = '';
  try {
    await loadSlotItems(selectedSlotKey.value);
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
    setSuccess(t('common.saveSuccess'));
  } catch {
    errorText.value = t('common.saveFailed');
  } finally {
    loading.value = false;
  }
}

function startFooterEdit(item: FooterItem): void {
  editingFooterId.value = item.id;
  footerForm.label = item.label;
  footerForm.linkType = item.linkType;
  footerForm.internalArticleSlug = item.internalArticleSlug || '';
  footerForm.externalUrl = item.externalUrl || '';
  footerForm.rowNum = item.rowNum;
  footerForm.orderNum = item.orderNum;
  footerForm.enabled = item.enabled;
}

function cancelFooterEdit(): void {
  resetFooterForm();
}

async function submitFooter(): Promise<void> {
  loading.value = true;
  errorText.value = '';
  successText.value = '';
  try {
    if (editingFooterId.value) {
      await updateFooterItem(editingFooterId.value, footerForm);
    } else {
      await createFooterItem(footerForm);
    }
    setSuccess(t('common.saveSuccess'));
    resetFooterForm();
    await loadAll();
  } catch {
    errorText.value = t('common.saveFailed');
    loading.value = false;
  }
}

async function removeFooter(id: string): Promise<void> {
  if (!window.confirm(t('common.confirmDelete'))) {
    return;
  }
  loading.value = true;
  errorText.value = '';
  successText.value = '';
  try {
    await deleteFooterItem(id);
    setSuccess(t('common.deleteSuccess'));
    if (editingFooterId.value === id) {
      resetFooterForm();
    }
    await loadAll();
  } catch {
    errorText.value = t('common.deleteFailed');
    loading.value = false;
  }
}

function startSocialEdit(item: SocialLink): void {
  editingSocialId.value = item.id;
  socialForm.platform = item.platform;
  socialForm.title = item.title;
  socialForm.url = item.url;
  socialForm.iconKey = item.iconKey || '';
  socialForm.orderNum = item.orderNum;
  socialForm.enabled = item.enabled;
}

function cancelSocialEdit(): void {
  resetSocialForm();
}

async function submitSocial(): Promise<void> {
  loading.value = true;
  errorText.value = '';
  successText.value = '';
  try {
    if (editingSocialId.value) {
      await updateSocialLink(editingSocialId.value, socialForm);
    } else {
      await createSocialLink(socialForm);
    }
    setSuccess(t('common.saveSuccess'));
    resetSocialForm();
    await loadAll();
  } catch {
    errorText.value = t('common.saveFailed');
    loading.value = false;
  }
}

async function removeSocial(id: string): Promise<void> {
  if (!window.confirm(t('common.confirmDelete'))) {
    return;
  }
  loading.value = true;
  errorText.value = '';
  successText.value = '';
  try {
    await deleteSocialLink(id);
    setSuccess(t('common.deleteSuccess'));
    if (editingSocialId.value === id) {
      resetSocialForm();
    }
    await loadAll();
  } catch {
    errorText.value = t('common.deleteFailed');
    loading.value = false;
  }
}

function startNavEdit(item: NavItem): void {
  editingNavId.value = item.id;
  navForm.name = item.name;
  navForm.key = item.key;
  navForm.type = item.type;
  navForm.targetType = item.targetType;
  navForm.targetValue = item.targetValue || '';
  navForm.orderNum = item.orderNum;
  navForm.enabled = item.enabled;
}

function cancelNavEdit(): void {
  resetNavForm();
}

async function submitNav(): Promise<void> {
  loading.value = true;
  errorText.value = '';
  successText.value = '';
  try {
    if (editingNavId.value) {
      await updateNavItem(editingNavId.value, navForm);
    } else {
      await createNavItem(navForm);
    }
    setSuccess(t('common.saveSuccess'));
    resetNavForm();
    await loadAll();
  } catch {
    errorText.value = t('common.saveFailed');
    loading.value = false;
  }
}

async function removeNav(id: string): Promise<void> {
  if (!window.confirm(t('common.confirmDelete'))) {
    return;
  }
  loading.value = true;
  errorText.value = '';
  successText.value = '';
  try {
    await deleteNavItem(id);
    setSuccess(t('common.deleteSuccess'));
    if (editingNavId.value === id) {
      resetNavForm();
    }
    await loadAll();
  } catch {
    errorText.value = t('common.deleteFailed');
    loading.value = false;
  }
}

async function submitSlot(): Promise<void> {
  if (!slotForm.slotKey.trim() || !slotForm.name.trim()) {
    errorText.value = t('site.slotRequired');
    return;
  }
  loading.value = true;
  errorText.value = '';
  successText.value = '';
  try {
    await createSlot(slotForm);
    setSuccess(t('common.saveSuccess'));
    slotForm.slotKey = '';
    slotForm.name = '';
    slotForm.description = '';
    await loadAll();
  } catch {
    errorText.value = t('common.saveFailed');
    loading.value = false;
  }
}

async function submitSlotItem(): Promise<void> {
  if (!selectedSlotKey.value || !slotItemForm.contentId.trim()) {
    errorText.value = t('site.slotItemRequired');
    return;
  }
  loading.value = true;
  errorText.value = '';
  successText.value = '';
  try {
    await createSlotItem(selectedSlotKey.value, slotItemForm);
    setSuccess(t('common.saveSuccess'));
    slotItemForm.contentId = '';
    slotItemForm.orderNum = 1;
    await loadSlotItems(selectedSlotKey.value);
  } catch {
    errorText.value = t('common.saveFailed');
  } finally {
    loading.value = false;
  }
}

async function removeSlotItem(id: string): Promise<void> {
  if (!selectedSlotKey.value || !window.confirm(t('common.confirmDelete'))) {
    return;
  }
  loading.value = true;
  errorText.value = '';
  successText.value = '';
  try {
    await deleteSlotItem(selectedSlotKey.value, id);
    setSuccess(t('common.deleteSuccess'));
    await loadSlotItems(selectedSlotKey.value);
  } catch {
    errorText.value = t('common.deleteFailed');
  } finally {
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
  grid-template-columns: repeat(8, minmax(0, 1fr));
  gap: 8px;
}

.slot-tools {
  display: flex;
  justify-content: flex-start;
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

button.danger {
  border-color: #e8b9b9;
  color: #b64040;
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

.row-actions {
  display: flex;
  align-items: center;
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

  .list li {
    align-items: flex-start;
    flex-direction: column;
  }
}
</style>
