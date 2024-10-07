<template>
  <div class="max-w-3xl mx-auto px-3">
    <h1 class="text-2xl font-medium" v-if="route.params.id">
      分类-{{ articleList[0]?.categoryName }}
    </h1>
    <h1 class="text-2xl font-medium" v-else>全部文章</h1>
    <div class="text-sm text-gray">{{ total }} 篇文章</div>

    <TimelineList :list="articleList" :total="total">
      <template #item="{ item }">
        <router-link class="hover:text-primary" :to="`/article/${item.id}`">
          <div
            class="flex flex-col hover:bg-primary-bg-1 p-2 rounded-lg hover:shadow-md hover:scale-105 transition-all"
          >
            <div class="flex flex-row items-center justify-between">
              <div flex-1>
                {{ item.title }}
                <div class="text-sm text-gray mt-2 mb-4">{{ item.summary }}</div>
              </div>
              <div class="text-sm text-gray flex flex-row items-center gap-2">
                <a-tag v-if="item.isTop === '1'" class="bg-primary text-white">置顶</a-tag>
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
import { onMounted, onUnmounted, ref, watch } from 'vue'
import { useRoute } from 'vue-router'
import TimelineList from '@/components/TimelineList.vue'
import { EyeOutlined, FolderOutlined } from '@ant-design/icons-vue'
import Icon from '@ant-design/icons-vue'
import { handleScroll } from '@/utils/handleScroll'

const articleList = ref<ArticleListData[]>([])
const total = ref<number>(0)
const route = useRoute()
const loading = ref(false)
const pageParam = ref({
  pageNum: 1,
  pageSize: 10,
  categoryId: route.params.id ? Number(route.params.id) : undefined,
})

const getArticleList = async (replace = false) => {
  if (loading.value) return
  loading.value = true
  try {
    const res = await reqArticlePage(pageParam.value)
    if (replace) {
      articleList.value = res.data.rows
    } else {
      articleList.value = [...articleList.value, ...res.data.rows]
    }
    total.value = res.data.total
    pageParam.value.pageNum++
  } finally {
    loading.value = false
  }
}

watch(
  () => route.params.id,
  async () => {
    pageParam.value.pageNum = 1
    pageParam.value.categoryId = route.params.id ? Number(route.params.id) : undefined
    await getArticleList(true)
  },
)

const scroll = () => {
  handleScroll(getArticleList, loading.value, total.value === articleList.value.length)
}

onMounted(async () => {
  await getArticleList()
  window.addEventListener('scroll', scroll)
})

onUnmounted(() => {
  window.removeEventListener('scroll', scroll)
})
</script>

<style scoped lang="scss"></style>
