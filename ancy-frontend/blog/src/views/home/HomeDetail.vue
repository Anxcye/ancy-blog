<template>
  <div>
    <ArticleViewer :article="article">
      <template #header>
        <div class="text-2xl font-bold text-center mt-40 mb-5 md:text-left">
          {{ article?.title }}
        </div>
        <div class="text-sm text-gray text-center md:text-left md:mb-20 md:p-0">
          {{ article?.summary }}
        </div>
      </template>
    </ArticleViewer>
  </div>
</template>

<script setup lang="ts">
import { reqArticleHomeGetById } from '@/api/article'
import type { ArticleDetailData } from '@/api/article/type'
import ArticleViewer from '@/components/ArticleViewer.vue'
import { useBrowserStore } from '@/stores/browser'
import getMeta from '@/utils/meta'
import { useHead } from '@vueuse/head'
import { computed, onMounted, ref, watch } from 'vue'
import { useRoute } from 'vue-router'

const route = useRoute()
const article = ref<ArticleDetailData>()

const getArticleDetail = async (id: number) => {
  const res = await reqArticleHomeGetById(id)
  article.value = res.data
  useBrowserStore().setTitle(article.value?.title ?? '')
}

watch(
  () => route.params.id,
  async (newId) => {
    if (newId) {
      await getArticleDetail(Number(newId))
    }
  },
  { immediate: true },
)

onMounted(async () => {
  await getArticleDetail(Number(route.params.id))
})

useHead({
  meta: getMeta(
    computed(() => article.value?.summary ?? ''),
    computed(() => article.value?.tags.map((item) => item.name) ?? []),
  ),
})
</script>

<style scoped></style>
