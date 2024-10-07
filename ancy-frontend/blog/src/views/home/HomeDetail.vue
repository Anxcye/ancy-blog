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
import { onMounted, ref, watch } from 'vue'
import { useRoute } from 'vue-router'

const route = useRoute()
const article = ref<ArticleDetailData>()

const getArticleDetail = async (id: number) => {
  const res = await reqArticleHomeGetById(id)
  article.value = res.data
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
</script>

<style scoped></style>
