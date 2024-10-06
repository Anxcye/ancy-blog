<template>
  <SayHi />
  <div class="px-4 py-2">
    <span class="text-xl font-medium">文章</span>
    <TimelineList :list="articleList">
      <template #item="{ item }">
        <router-link class="hover:text-primary" :to="`/article/${item.id}`">
          {{ item.title }}
        </router-link>
      </template>
    </TimelineList>
  </div>
  <HomeFooter />
</template>

<script setup lang="ts">
import { reqArticleRecent } from '@/api/article'
import SayHi from './components/SayHi.vue'
import { onMounted, ref } from 'vue'
import type { ArticleListData } from '@/api/article/type'
import HomeFooter from './components/HomeFooter.vue'
import TimelineList from '@/components/TimelineList.vue'
const articleList = ref<ArticleListData[]>([])

const getArticleList = async () => {
  const res = await reqArticleRecent()
  articleList.value = res.data
}

onMounted(async () => {
  getArticleList()
})
</script>

<style scoped></style>
