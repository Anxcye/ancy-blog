<template>
  <div>
    <ArticleViewer :article="article">
      <template #header>
        <div class="text-2xl font-bold text-center mt-40 mb-5">{{ article?.title }}</div>
        <div
          class="text-sm text-gray text-center mb-5 flex flex-row gap-3 items-center justify-center"
        >
          <span>
            <a-tooltip>
              <template #title>
                {{ article?.createTime || '' }}
              </template>
              <icon :component="CalendarOutlined" />
              {{ timeAgo(new Date(article?.createTime || '')) }}
            </a-tooltip>
          </span>
          <span>
            <icon :component="FolderOutlined" />
            <router-link :to="`/category/${article?.categoryId}`">
              {{ article?.categoryName }}
            </router-link>
          </span>
          <span>
            <icon :component="EyeOutlined" />
            {{ article?.viewCount }}
          </span>
        </div>
        <div class="flex flex-row items-center justify-center mb-5">
          <div v-for="tag in article?.tags" :key="tag.id">
            <a-tag># {{ tag.name }}</a-tag>
          </div>
        </div>
        <div class="text-sm text-gray text-left border border-gray-bg rounded-2xl p-5">
          {{ article?.summary }}
        </div>
      </template>
    </ArticleViewer>
  </div>
</template>

<script setup lang="ts">
import { reqArticleGetById } from '@/api/article'
import type { ArticleDetailData } from '@/api/article/type'
import { onMounted, ref, computed } from 'vue'
import { useRoute } from 'vue-router'
import Icon from '@ant-design/icons-vue'
import { CalendarOutlined, FolderOutlined, EyeOutlined } from '@ant-design/icons-vue'
import timeAgo from '@/utils/timeAgo'
import ArticleViewer from '@/components/ArticleViewer.vue'
import { useBrowserStore } from '@/stores/browser'
import { useHead } from '@vueuse/head'
import getMeta from '@/utils/meta'

const route = useRoute()
const article = ref<ArticleDetailData>()

useHead({
  meta: getMeta(
    computed(() => article.value?.summary ?? ''),
    computed(() => article.value?.tags.map((item) => item.name) ?? []),
  ),
})

const getArticle = async () => {
  const res = await reqArticleGetById(Number(route.params.id))
  article.value = res.data
  useBrowserStore().setTitle(article.value?.title ?? '')
}

onMounted(async () => {
  await getArticle()
})
</script>

<style scoped lang="scss"></style>
