<!--
File: DashboardView.vue
Purpose: Display admin overview metrics and quick actions for daily operation workflow.
Module: frontend-admin/views/dashboard, presentation layer.
Related: dashboard API module and top-level management routes.
-->
<template>
  <section class="dashboard-page">
    <header class="header">
      <h1>{{ t('dashboard.title') }}</h1>
      <p>{{ t('dashboard.subtitle') }}</p>
    </header>

    <p v-if="errorText" class="error">{{ errorText }}</p>

    <div class="cards">
      <article class="card">
        <h2>{{ t('dashboard.cardArticleTotal') }}</h2>
        <p>{{ metrics.articleTotal }}</p>
      </article>
      <article class="card">
        <h2>{{ t('dashboard.cardArticleDraft') }}</h2>
        <p>{{ metrics.articleDraft }}</p>
      </article>
      <article class="card">
        <h2>{{ t('dashboard.cardArticlePublished') }}</h2>
        <p>{{ metrics.articlePublished }}</p>
      </article>
      <article class="card">
        <h2>{{ t('dashboard.cardMomentTotal') }}</h2>
        <p>{{ metrics.momentTotal }}</p>
      </article>
      <article class="card">
        <h2>{{ t('dashboard.cardCommentPending') }}</h2>
        <p>{{ metrics.commentPending }}</p>
      </article>
      <article class="card">
        <h2>{{ t('dashboard.cardLinkPending') }}</h2>
        <p>{{ metrics.linkPending }}</p>
      </article>
    </div>

    <section class="quick-actions">
      <h2>{{ t('dashboard.quickActions') }}</h2>
      <div class="action-grid">
        <RouterLink :to="{ name: 'article-new' }">{{ t('dashboard.actionNewArticle') }}</RouterLink>
        <RouterLink :to="{ name: 'moments' }">{{ t('dashboard.actionManageMoments') }}</RouterLink>
        <RouterLink :to="{ name: 'interaction' }">{{ t('dashboard.actionReviewInteractions') }}</RouterLink>
        <RouterLink :to="{ name: 'system' }">{{ t('dashboard.actionTranslationCenter') }}</RouterLink>
      </div>
    </section>
  </section>
</template>

<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue';
import { useI18n } from 'vue-i18n';

import { loadDashboardMetrics } from '@/api/modules/dashboard';

const { t } = useI18n();
const errorText = ref('');
const metrics = reactive({
  articleTotal: 0,
  articleDraft: 0,
  articlePublished: 0,
  momentTotal: 0,
  commentPending: 0,
  linkPending: 0,
});

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
  gap: 14px;
}

.header h1 {
  margin: 0;
}

.header p {
  margin: 4px 0 0;
  color: var(--muted);
}

.cards {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 10px;
}

.card {
  border: 1px solid var(--border);
  border-radius: 12px;
  background: var(--surface);
  padding: 14px;
}

.card h2 {
  margin: 0;
  font-size: 14px;
  color: var(--muted);
  font-weight: 500;
}

.card p {
  margin: 8px 0 0;
  font-size: 24px;
  font-weight: 700;
}

.quick-actions {
  border: 1px solid var(--border);
  border-radius: 12px;
  background: var(--surface);
  padding: 14px;
}

.quick-actions h2 {
  margin: 0 0 10px;
}

.action-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 8px;
}

.action-grid a {
  text-decoration: none;
  border: 1px solid var(--border);
  border-radius: 10px;
  padding: 10px;
  color: var(--text);
  transition: all 0.2s ease;
}

.action-grid a:hover {
  border-color: var(--accent);
  background: var(--accent-soft);
  color: var(--accent-hover);
}

.error {
  color: #b64040;
  margin: 0;
}

@media (max-width: 900px) {
  .cards,
  .action-grid {
    grid-template-columns: 1fr;
  }
}
</style>
