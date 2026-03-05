<!--
File: DashboardView.vue
Purpose: Display operational metrics and action shortcuts for admin daily workflow.
Module: frontend-admin/views/dashboard, presentation layer.
Related: dashboard API module, content/interaction routes, global loading state.
-->
<template>
  <section class="dashboard-page">
    <NAlert v-if="errorText" type="error" :show-icon="false">{{ errorText }}</NAlert>

    <NGrid cols="1 s:2 m:3" responsive="screen" :x-gap="12" :y-gap="12">
      <NGridItem v-for="item in statCards" :key="item.key">
        <NCard :bordered="false" class="metric-card">
          <NStatistic :label="item.label" :value="item.value" />
        </NCard>
      </NGridItem>
    </NGrid>

    <NCard :bordered="false" class="actions-card">
      <template #header>{{ t('dashboard.quickActions') }}</template>
      <div class="action-grid">
        <NButton type="primary" secondary @click="router.push({ name: 'article-new' })">{{ t('dashboard.actionNewArticle') }}</NButton>
        <NButton secondary @click="router.push({ name: 'moments' })">{{ t('dashboard.actionManageMoments') }}</NButton>
        <NButton secondary @click="router.push({ name: 'interaction' })">{{ t('dashboard.actionReviewInteractions') }}</NButton>
        <NButton secondary @click="router.push({ name: 'system' })">{{ t('dashboard.actionTranslationCenter') }}</NButton>
      </div>
    </NCard>
  </section>
</template>

<script setup lang="ts">
import { computed, onMounted, reactive, ref } from 'vue';
import { useRouter } from 'vue-router';
import { useI18n } from 'vue-i18n';
import { NAlert, NButton, NCard, NGrid, NGridItem, NStatistic } from 'naive-ui';

import { loadDashboardMetrics } from '@/api/modules/dashboard';

const { t } = useI18n();
const router = useRouter();
const errorText = ref('');
const metrics = reactive({
  articleTotal: 0,
  articleDraft: 0,
  articlePublished: 0,
  momentTotal: 0,
  commentPending: 0,
  linkPending: 0,
});

const statCards = computed(() => [
  { key: 'articleTotal', label: t('dashboard.cardArticleTotal'), value: metrics.articleTotal },
  { key: 'articleDraft', label: t('dashboard.cardArticleDraft'), value: metrics.articleDraft },
  { key: 'articlePublished', label: t('dashboard.cardArticlePublished'), value: metrics.articlePublished },
  { key: 'momentTotal', label: t('dashboard.cardMomentTotal'), value: metrics.momentTotal },
  { key: 'commentPending', label: t('dashboard.cardCommentPending'), value: metrics.commentPending },
  { key: 'linkPending', label: t('dashboard.cardLinkPending'), value: metrics.linkPending },
]);

onMounted(async () => {
  errorText.value = '';
  try {
    const data = await loadDashboardMetrics();
    Object.assign(metrics, data);
  } catch {
    errorText.value = t('common.loadFailed');
  }
});
</script>

<style scoped>
.dashboard-page {
  display: grid;
  gap: 12px;
}

.metric-card,
.actions-card {
  border-radius: 14px;
  box-shadow: 0 6px 24px rgba(15, 31, 36, 0.05);
}

.action-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 10px;
}

@media (max-width: 900px) {
  .action-grid {
    grid-template-columns: 1fr;
  }
}
</style>
