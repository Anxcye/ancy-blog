<template>
  <div>
    <TopGradient height="100%" />
    <div class="text-2xl font-bold text-center mt-40 mb-5">{{ article?.title }}</div>
    <div class="text-sm text-gray text-center mb-5 flex flex-row gap-3 items-center justify-center">
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

    <MdViewer :content="article?.content" />
  </div>
</template>

<script setup lang="ts">
import { reqArticleGetById } from '@/api/article'
import type { ArticleDetailData } from '@/api/article/type'
import { onMounted, ref } from 'vue'
import { useRoute } from 'vue-router'
import MdViewer from '@/components/MdViewer.vue'
import TopGradient from '@/components/TopGradient.vue'
import Icon from '@ant-design/icons-vue'
import { CalendarOutlined, FolderOutlined, EyeOutlined } from '@ant-design/icons-vue'
import timeAgo from '@/utils/timeAgo'

const route = useRoute()
const article = ref<ArticleDetailData>()

const getArticle = async () => {
  const res = await reqArticleGetById(Number(route.params.id))
  article.value = res.data
}
onMounted(async () => {
  await getArticle()
})
</script>

<style scoped lang="scss"></style>
