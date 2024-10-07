<template>
  <div class="content">
    <div class="title">阅读</div>
    <div class="summary">记下了{{ total }}条</div>
    <Waterfall :list="readList">
      <template #default="{ item }">
        <div
          class="border-l-4 p-2"
          :style="{
            borderColor: getColor(item.source),
            backgroundColor: getColor(item.source, true),
          }"
        >
          <div>{{ item.content }}</div>
          <TimeTip :time="item.createTime" class="text-sm text-gray" />
          <div class="flex flex-row justify-between gap-2">
            <div v-if="item.addFrom === 1">
              <a-tooltip>
                <template #title>
                  <a :href="anxReaderUrl" target="_blank">
                    摘抄于安读，点击
                    <span class="text-primary">这里</span>
                    了解更多
                  </a>
                </template>
                <a-tag>安读</a-tag>
              </a-tooltip>
            </div>
            <div v-else></div>
            <div class="text-sm text-gray mt-2">
              来自{{ item.author }}{{ item.source && item.author ? '的' : ''
              }}{{ '《' + item.source + '》' }}
            </div>
          </div>
        </div>
      </template>
    </Waterfall>
    <div v-if="loading" class="text-center py-4">正在加载更多...</div>
  </div>
</template>

<script setup lang="ts">
import { reqReadPage } from '@/api/read'
import type { ReadData, ReadPageParam } from '@/api/read/type'
import { useColorStore } from '@/stores/color'
import { onMounted, ref, onUnmounted } from 'vue'
import TimeTip from '@/components/TimeTip.vue'
import { Waterfall } from 'vue-waterfall-plugin-next'
import 'vue-waterfall-plugin-next/dist/style.css'
import { handleScroll } from '@/utils/handleScroll'

const anxReaderUrl = 'https://github.com/anxcye/anx-reader'
const readPageParam = ref<ReadPageParam>({
  pageNum: 1,
  pageSize: 10,
})
const colorStore = useColorStore()
const readList = ref<ReadData[]>([])
const total = ref(0)
const loading = ref(false)

const getReadList = async () => {
  if (loading.value) return
  loading.value = true
  try {
    const res = await reqReadPage(readPageParam.value)
    readList.value = [...readList.value, ...res.data.rows]
    total.value = res.data.total
    readPageParam.value.pageNum++
  } finally {
    loading.value = false
  }
}

const getColor = (source: string, bg = false) => {
  return bg ? colorStore.getColor(source) + '40' : colorStore.getColor(source)
}

const scroll = () => {
  handleScroll(getReadList, loading.value, total.value === readList.value.length)
}

onMounted(async () => {
  await getReadList()
  window.addEventListener('scroll', scroll)
})

onUnmounted(() => {
  window.removeEventListener('scroll', scroll)
})
</script>

<style scoped lang="scss">
.waterfall-list {
  background-color: transparent !important;
}
</style>
