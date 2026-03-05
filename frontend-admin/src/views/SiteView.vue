<!--
File: SiteView.vue
Purpose: Manage site settings, footer items, social links, nav entries, and content slots with unified Naive UI components.
Module: frontend-admin/views/site, presentation layer.
Related: site API module and backend site/admin endpoints.
-->
<template>
  <section class="site-page">
    <NCard :bordered="false" class="section-card">
      <NAlert v-if="errorText" type="error" :show-icon="false">{{ errorText }}</NAlert>
      <NAlert v-if="successText" type="success" :show-icon="false">{{ successText }}</NAlert>

      <NCard size="small" :bordered="true">
        <template #header>{{ t('site.settingsTitle') }}</template>
        <NForm label-placement="top">
          <div class="grid-2">
            <NFormItem :label="t('site.siteName')">
              <NInput v-model:value="settingsForm.siteName" />
            </NFormItem>
            <NFormItem :label="t('site.defaultLocale')">
              <NInput v-model:value="settingsForm.defaultLocale" placeholder="zh-CN" />
            </NFormItem>
            <NFormItem :label="t('site.avatarUrl')">
              <NInput v-model:value="settingsForm.avatarUrl" />
            </NFormItem>
            <NFormItem :label="t('site.heroIntroMd')">
              <NInput v-model:value="settingsForm.heroIntroMd" type="textarea" :autosize="{ minRows: 4, maxRows: 10 }" />
            </NFormItem>
          </div>
        </NForm>
        <NButton type="primary" :loading="loading" @click="saveSettings">{{ t('common.save') }}</NButton>
      </NCard>

      <NCard size="small" :bordered="true">
        <template #header>{{ t('site.footerTitle') }}</template>
        <NForm label-placement="top">
          <div class="grid-3">
            <NFormItem :label="t('site.footerLabel')">
              <NInput v-model:value="footerForm.label" />
            </NFormItem>
            <NFormItem label="linkType">
              <NSelect v-model:value="footerForm.linkType" :options="footerLinkTypeOptions" />
            </NFormItem>
            <NFormItem label="internalSlug">
              <NInput v-model:value="footerForm.internalArticleSlug" />
            </NFormItem>
            <NFormItem label="externalUrl">
              <NInput v-model:value="footerForm.externalUrl" />
            </NFormItem>
            <NFormItem label="rowNum">
              <NInputNumber v-model:value="footerForm.rowNum" :min="1" :max="3" style="width: 100%" />
            </NFormItem>
            <NFormItem label="orderNum">
              <NInputNumber v-model:value="footerForm.orderNum" :min="1" style="width: 100%" />
            </NFormItem>
          </div>
        </NForm>
        <NSpace>
          <NButton :loading="loading" @click="submitFooter">{{ editingFooterId ? t('common.update') : t('common.add') }}</NButton>
          <NButton v-if="editingFooterId" :loading="loading" tertiary @click="cancelFooterEdit">{{ t('common.cancel') }}</NButton>
        </NSpace>

        <NDataTable remote :loading="loading" :columns="footerColumns" :data="footerItems" :pagination="false" :row-key="rowKey" class="table-block" />
      </NCard>

      <NCard size="small" :bordered="true">
        <template #header>{{ t('site.socialTitle') }}</template>
        <NForm label-placement="top">
          <div class="grid-3">
            <NFormItem label="platform">
              <NInput v-model:value="socialForm.platform" />
            </NFormItem>
            <NFormItem :label="t('site.socialTitleField')">
              <NInput v-model:value="socialForm.title" />
            </NFormItem>
            <NFormItem label="url">
              <NInput v-model:value="socialForm.url" />
            </NFormItem>
            <NFormItem label="iconKey">
              <NInput v-model:value="socialForm.iconKey" />
            </NFormItem>
            <NFormItem label="orderNum">
              <NInputNumber v-model:value="socialForm.orderNum" :min="1" style="width: 100%" />
            </NFormItem>
          </div>
        </NForm>
        <NSpace>
          <NButton :loading="loading" @click="submitSocial">{{ editingSocialId ? t('common.update') : t('common.add') }}</NButton>
          <NButton v-if="editingSocialId" :loading="loading" tertiary @click="cancelSocialEdit">{{ t('common.cancel') }}</NButton>
        </NSpace>

        <NDataTable remote :loading="loading" :columns="socialColumns" :data="socialLinks" :pagination="false" :row-key="rowKey" class="table-block" />
      </NCard>

      <NCard size="small" :bordered="true">
        <template #header>{{ t('site.navTitle') }}</template>
        <NForm label-placement="top">
          <div class="grid-3">
            <NFormItem label="parentId (optional)">
              <NSelect v-model:value="navForm.parentId" :options="parentNavOptions" clearable />
            </NFormItem>
            <NFormItem :label="t('site.navName')">
              <NInput v-model:value="navForm.name" />
            </NFormItem>
            <NFormItem label="key">
              <NInput v-model:value="navForm.key" />
            </NFormItem>
            <NFormItem label="type">
              <NInput v-model:value="navForm.type" />
            </NFormItem>
            <NFormItem label="targetType">
              <NInput v-model:value="navForm.targetType" />
            </NFormItem>
            <NFormItem label="targetValue">
              <NInput v-model:value="navForm.targetValue" />
            </NFormItem>
            <NFormItem label="orderNum">
              <NInputNumber v-model:value="navForm.orderNum" :min="1" style="width: 100%" />
            </NFormItem>
          </div>
        </NForm>
        <NSpace>
          <NButton :loading="loading" @click="submitNav">{{ editingNavId ? t('common.update') : t('common.add') }}</NButton>
          <NButton v-if="editingNavId" :loading="loading" tertiary @click="cancelNavEdit">{{ t('common.cancel') }}</NButton>
        </NSpace>

        <NDataTable remote :loading="loading" :columns="navColumns" :data="navItems" :pagination="false" :row-key="rowKey" class="table-block" />
      </NCard>

      <NCard size="small" :bordered="true">
        <template #header>{{ t('site.slotTitle') }}</template>
        <NForm label-placement="top">
          <div class="grid-3">
            <NFormItem label="slotKey">
              <NInput v-model:value="slotForm.slotKey" />
            </NFormItem>
            <NFormItem :label="t('site.slotName')">
              <NInput v-model:value="slotForm.name" />
            </NFormItem>
            <NFormItem :label="t('site.slotDescription')">
              <NInput v-model:value="slotForm.description" />
            </NFormItem>
          </div>
        </NForm>
        <NButton :loading="loading" @click="submitSlot">{{ t('common.add') }}</NButton>

        <NSpace class="slot-tools" vertical>
          <NSelect v-model:value="selectedSlotKey" :options="slotOptions" style="max-width: 360px" @update:value="onSlotChange" />

          <div v-if="selectedSlotKey" class="slot-inline">
            <NSelect v-model:value="slotItemForm.contentType" :options="sourceTypeOptions" style="width: 140px" />
            <NInput v-model:value="slotItemForm.contentId" placeholder="content id (uuid)" />
            <NInputNumber v-model:value="slotItemForm.orderNum" :min="1" style="width: 140px" />
            <NButton :loading="loading" @click="submitSlotItem">{{ t('common.add') }}</NButton>
          </div>
        </NSpace>

        <NDataTable remote :loading="loading" :columns="slotItemColumns" :data="slotItems" :pagination="false" :row-key="rowKey" class="table-block" />
      </NCard>
    </NCard>
  </section>
</template>

<script setup lang="ts">
import { computed, h, onMounted, reactive, ref } from 'vue';
import { useI18n } from 'vue-i18n';
import type { DataTableColumns } from 'naive-ui';
import { NAlert, NButton, NCard, NDataTable, NForm, NFormItem, NInput, NInputNumber, NSelect, NSpace } from 'naive-ui';

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
  parentId: '',
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

const footerLinkTypeOptions = [
  { label: 'none', value: 'none' },
  { label: 'internal', value: 'internal' },
  { label: 'external', value: 'external' },
];

const sourceTypeOptions = [
  { label: 'article', value: 'article' },
  { label: 'moment', value: 'moment' },
];

const slotOptions = computed(() => slots.value.map((slot) => ({
  label: `${slot.slotKey} · ${slot.name}`,
  value: slot.slotKey,
})));

const parentNavOptions = computed(() => {
  const topLevel = navItems.value.filter(n => !n.parentId);
  return topLevel.map(n => ({ label: n.name, value: n.id }));
});

function rowKey(row: { id: string }): string {
  return row.id;
}

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
  navForm.parentId = '';
  navForm.name = '';
  navForm.key = '';
  navForm.type = 'menu';
  navForm.targetType = 'route';
  navForm.targetValue = '';
  navForm.orderNum = 1;
  navForm.enabled = true;
}

const footerColumns = computed<DataTableColumns<FooterItem>>(() => [
  { title: t('site.footerLabel'), key: 'label' },
  { title: 'linkType', key: 'linkType', width: 120 },
  {
    title: 'row/order',
    key: 'order',
    width: 120,
    render(row) {
      return `${row.rowNum}/${row.orderNum}`;
    },
  },
  {
    title: t('articles.colAction'),
    key: 'actions',
    width: 170,
    render(row) {
      return h(NSpace, { wrap: false, size: 8 }, {
        default: () => [
          h(NButton, { size: 'small', tertiary: true, onClick: () => startFooterEdit(row) }, { default: () => t('common.edit') }),
          h(NButton, { size: 'small', tertiary: true, type: 'error', onClick: () => removeFooter(row.id) }, { default: () => t('common.delete') }),
        ],
      });
    },
  },
]);

const socialColumns = computed<DataTableColumns<SocialLink>>(() => [
  { title: 'platform', key: 'platform', width: 120 },
  { title: t('site.socialTitleField'), key: 'title', width: 150 },
  { title: 'url', key: 'url' },
  {
    title: t('articles.colAction'),
    key: 'actions',
    width: 170,
    render(row) {
      return h(NSpace, { wrap: false, size: 8 }, {
        default: () => [
          h(NButton, { size: 'small', tertiary: true, onClick: () => startSocialEdit(row) }, { default: () => t('common.edit') }),
          h(NButton, { size: 'small', tertiary: true, type: 'error', onClick: () => removeSocial(row.id) }, { default: () => t('common.delete') }),
        ],
      });
    },
  },
]);

const navColumns = computed<DataTableColumns<NavItem>>(() => [
  { title: t('site.navName'), key: 'name', width: 150 },
  { title: 'key', key: 'key', width: 120 },
  {
    title: 'target',
    key: 'target',
    render(row) {
      return `${row.targetType}:${row.targetValue}`;
    },
  },
  {
    title: t('articles.colAction'),
    key: 'actions',
    width: 170,
    render(row) {
      return h(NSpace, { wrap: false, size: 8 }, {
        default: () => [
          h(NButton, { size: 'small', tertiary: true, onClick: () => startNavEdit(row) }, { default: () => t('common.edit') }),
          h(NButton, { size: 'small', tertiary: true, type: 'error', onClick: () => removeNav(row.id) }, { default: () => t('common.delete') }),
        ],
      });
    },
  },
]);

const slotItemColumns = computed<DataTableColumns<SlotItem>>(() => [
  {
    title: 'content',
    key: 'content',
    render(row) {
      return `${row.contentType} · ${row.contentId}`;
    },
  },
  { title: 'order', key: 'orderNum', width: 100 },
  {
    title: t('articles.colAction'),
    key: 'actions',
    width: 110,
    render(row) {
      return h(NButton, {
        size: 'small',
        tertiary: true,
        type: 'error',
        onClick: () => removeSlotItem(row.id),
      }, { default: () => t('common.delete') });
    },
  },
]);

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
  navForm.parentId = item.parentId || '';
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
  if (!selectedSlotKey.value) {
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
}

.section-card {
  border-radius: 14px;
  box-shadow: 0 6px 24px color-mix(in srgb, var(--n-text-color) 8%, transparent);
  display: grid;
  gap: 12px;
}

.grid-2 {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 10px;
}

.grid-3 {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 10px;
}

.table-block {
  margin-top: 10px;
}

.slot-tools {
  margin-top: 10px;
}

.slot-inline {
  display: grid;
  grid-template-columns: 140px 1fr 140px auto;
  gap: 10px;
  align-items: center;
}

@media (max-width: 900px) {
  .grid-2,
  .grid-3,
  .slot-inline {
    grid-template-columns: 1fr;
  }
}
</style>
