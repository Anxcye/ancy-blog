<template>
  <div class="max-w-3xl mx-auto px-3">
    <ArticleViewer :article="projectDetail" v-if="projectDetail" :showComment="false">
      <template #header>
        <div class="flex flex-col gap-2 items-center md:items-start">
          <div class="text-2xl font-medium">{{ projectDetail?.title }}</div>
          <div class="flex flex-row items-center gap-2">
            <a-tag v-if="projectDetail?.type === '0'" class="bg-green-500 text-white">Active</a-tag>
            <a-tag v-if="projectDetail?.type === '1'">Archived</a-tag>
          </div>
          <div class="text-sm text-gray">{{ projectDetail?.summary }}</div>
          <div v-if="projectDetail?.srcUrl" class="mt-10">
            <span class="text-xl font-bold">开源地址：</span>
            <a :href="projectDetail?.srcUrl" target="_blank">
              {{ projectDetail?.srcUrl }}
            </a>
          </div>
          <div v-if="projectDetail?.displayUrl">
            <span class="text-xl font-bold">展示地址：</span>
            <a :href="projectDetail?.displayUrl" target="_blank">
              {{ projectDetail?.displayUrl }}
            </a>
          </div>
        </div>
      </template>
    </ArticleViewer>
  </div>
</template>

<script setup lang="ts">
import { useRoute } from 'vue-router'
import { reqProjectDetail } from '@/api/project'
import type { ProjectDetailData } from '@/api/project/type'
import { computed, onMounted, ref } from 'vue'
import ArticleViewer from '@/components/ArticleViewer.vue'
import { useBrowserStore } from '@/stores/browser'
import { useHead } from '@vueuse/head'
import getMeta from '@/utils/meta'

const route = useRoute()
const projectDetail = ref<ProjectDetailData>()

useHead({
  meta: getMeta(
    computed(() => projectDetail.value?.summary ?? ''),
    computed(() => [
      '项目',
      ...(projectDetail.value?.srcUrl
        ? ['开源', 'Github', projectDetail.value?.srcUrl.split('/').pop() ?? '']
        : []),
      ...(projectDetail.value?.displayUrl ? ['展示', projectDetail.value?.displayUrl] : []),
    ]),
  ),
})

const getProjectDetail = async () => {
  const res = await reqProjectDetail(Number(route.params.id))
  projectDetail.value = res.data
  useBrowserStore().setTitle('项目' + projectDetail.value?.title ?? '')
}

onMounted(async () => {
  await getProjectDetail()
})
</script>

<style scoped></style>
