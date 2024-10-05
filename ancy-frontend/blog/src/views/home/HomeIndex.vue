<template>
  <TopGradient />
  <SayHi />
  <div class="px-4 py-2">
    <span class="text-xl font-medium">文章</span>
    <a-timeline class="mt-8">
      <a-timeline-item
        :color="colorStore.getPrimaryColor()"
        v-for="item in articleList"
        :key="item.id"
      >
        <div class="flex flex-row items-center justify-between">
          <a class="flex-1" :href="`/article/${item.id}`">{{ item.title }}</a>
          <div class="text-sm text-gray-500">
            {{ timeAgo(new Date(item.createTime)) }}
          </div>
        </div>
      </a-timeline-item>
    </a-timeline>
  </div>
  <HomeFooter />
</template>

<script setup lang="ts">
import { reqArticlePage } from '@/api/article'
import SayHi from './components/SayHi.vue'
import { onMounted, ref } from 'vue'
import type { ArticleListData } from '@/api/article/type'
import { useColorStore } from '@/stores/color'
import timeAgo from '@/utils/timeAgo'
import HomeFooter from './components/HomeFooter.vue'
import TopGradient from '@/components/TopGradient.vue'

const colorStore = useColorStore()
const articleList = ref<ArticleListData[]>([])
const total = ref<number>(0)

const getArticleList = async () => {
  const res = await reqArticlePage({ pageNum: 1, pageSize: 10 })
  articleList.value = res.data.rows
  total.value = res.data.total
}

onMounted(async () => {
  getArticleList()
})
</script>

<style scoped></style>
