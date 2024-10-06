<template>
  <div class="max-w-3xl mx-auto px-3">
    <h1 class="text-2xl font-medium" v-if="route.params.id">
      分类-{{ articleList[0]?.categoryName }}
    </h1>
    <h1 class="text-2xl font-medium" v-else>全部文章</h1>
    <div class="text-sm text-gray">{{ total }} 篇文章</div>

    <TimelineList :list="articleList">
      <template #item="{ item }">
        <router-link class="hover:text-primary" :to="`/article/${item.id}`">
          <div
            class="flex flex-col hover:bg-primary-bg-1 p-2 rounded-lg hover:shadow-md hover:scale-105 transition"
          >
            <div class="flex flex-row items-center justify-between">
              <div flex-1>
                {{ item.title }}
                <div class="text-sm text-gray mt-2 mb-4">{{ item.summary }}</div>
              </div>
              <div class="text-sm text-gray flex flex-row items-center gap-2">
                <a-tag v-if="item.isTop === '1'" class="bg-primary text-white">置顶</a-tag>
                <div>{{ timeAgo(new Date(item.createTime)) }}</div>
              </div>
            </div>
            <div class="flex flex-row gap-2 justify-between">
              <div class="flex flex-row gap-2">
                <div v-if="!route.params.id">
                  <icon :component="FolderOutlined" />
                  {{ item.categoryName }}
                </div>
                <div class="flex flex-row gap-2">
                  <span>
                    <icon :component="EyeOutlined" />
                    {{ item.viewCount }}
                  </span>
                </div>
              </div>
              <div v-if="item.tags.length > 0">
                <a-tag v-for="tag in item.tags" :key="tag.id"># {{ tag.name }}</a-tag>
              </div>
            </div>
          </div>
        </router-link>
      </template>
    </TimelineList>
  </div>
</template>

<script setup lang="ts">
import { reqArticlePage } from '@/api/article'
import type { ArticleListData } from '@/api/article/type'
import { onMounted, ref, watch } from 'vue'
import { useRoute } from 'vue-router'
import TimelineList from '@/components/TimelineList.vue'
import timeAgo from '@/utils/timeAgo'
import { EyeOutlined, FolderOutlined } from '@ant-design/icons-vue'
import Icon from '@ant-design/icons-vue'

const articleList = ref<ArticleListData[]>([])
const total = ref<number>(0)
const route = useRoute()

const getArticleList = async () => {
  const res = await reqArticlePage({
    pageNum: 1,
    pageSize: 10,
    categoryId: route.params.id ? Number(route.params.id) : undefined,
  })
  articleList.value = res.data.rows
  total.value = res.data.total
}

watch(
  () => route.params.id,
  async () => {
    await getArticleList()
  },
)

onMounted(async () => {
  await getArticleList()
})
</script>

<style scoped lang="scss"></style>
