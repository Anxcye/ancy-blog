<template>
  <div class="max-w-3xl mx-auto px-3 pt-10">
    <ArticleViewer :article="article">
      <template #header>
        <h1 class="text-2xl font-bold mb-20">朋友</h1>
        <div class="grid grid-cols-2 md:grid-cols-3">
          <a
            :href="item.address"
            target="_blank"
            v-for="item in linkList"
            :key="item.id"
            class="py-4"
          >
            <div class="flex flex-col items-center justify-center gap-2">
              <a-avatar :src="item.logo" class="w-20 h-20 rounded-3xl" />
              <div class="text-lg font-medium">{{ item.name }}</div>
              <div class="text-sm text-gray text-center">
                {{ item.description }}
              </div>
            </div>
          </a>
        </div>
      </template>
      <template #footer>
        <div class="text-xl font-bold mb-4 mt-20 text-center">添加朋友</div>
        <div class="flex flex-col items-center gap-2 justify-between md:flex-row">
          <div class="flex flex-col gap-2 items-center w-full md:w-1/2">
            <a :href="addLinkParams.address" target="_blank" class="py-4">
              <div class="flex flex-col items-center justify-center gap-2">
                <a-avatar :src="addLinkParams.logo" class="w-20 h-20 rounded-3xl" />
                <div class="text-lg font-medium">{{ addLinkParams.name }}</div>
                <div class="text-sm text-gray text-center">
                  {{ addLinkParams.description }}
                </div>
              </div>
            </a>
            <div class="text-xl font-bold">预览</div>
          </div>
          <a-form :model="addLinkParams" :rules="rules" ref="formRef" class="w-full md:w-1/2">
            <a-form-item label="名称" name="name">
              <a-input v-model:value="addLinkParams.name" placeholder="input placeholder" />
            </a-form-item>
            <a-form-item label="logo" name="logo">
              <a-input v-model:value="addLinkParams.logo" placeholder="input placeholder" />
            </a-form-item>
            <a-form-item label="描述" name="description">
              <a-input v-model:value="addLinkParams.description" placeholder="input placeholder" />
            </a-form-item>
            <a-form-item label="地址" name="address">
              <a-input v-model:value="addLinkParams.address" placeholder="input placeholder" />
            </a-form-item>
            <a-form-item>
              <div class="flex flex-row items-center gap-2 md:justify-end">
                <a-button type="primary" @click="copyInfo">复制本站信息</a-button>
                <a-button type="primary" @click="addLink" :loading="submitLoading">
                  申请添加
                </a-button>
              </div>
            </a-form-item>
          </a-form>
        </div>
      </template>
    </ArticleViewer>
  </div>
</template>

<script setup lang="ts">
import type { ArticleDetailData } from '@/api/article/type'
import { reqLinkAdd, reqLinkGetArticle, reqLinkList } from '@/api/link'
import type { LinkAddParams, LinkListData } from '@/api/link/type'
import { onMounted, ref } from 'vue'
import ArticleViewer from '@/components/ArticleViewer.vue'
import { type FormInstance } from 'ant-design-vue'
import { useBaseInfoStore } from '@/stores/baseInfo'
import { message } from 'ant-design-vue'

const baseInfo = useBaseInfoStore()
const linkList = ref<LinkListData[]>([])
const article = ref<ArticleDetailData>()
const initialLinkParams = {
  name: baseInfo.getName() || '',
  logo: baseInfo.getAvatar() || '',
  description: baseInfo.getPhilosophy() || '',
  address: baseInfo.getAddress() || '',
}
const addLinkParams = ref<LinkAddParams>({
  ...initialLinkParams,
})
const addLinkList = ref<LinkAddParams[]>([])
const rules = ref({
  name: [{ required: true, message: '请输入名称' }],
  logo: [{ required: true, message: '请输入logo' }],
  description: [{ required: true, message: '请输入描述' }],
  address: [{ required: true, message: '请输入地址' }],
})
const formRef = ref<FormInstance>()
const submitLoading = ref(false)

const copyInfo = () => {
  navigator.clipboard.writeText(JSON.stringify(initialLinkParams))
  message.info('已复制，欢迎添加哦')
}

const addLink = async () => {
  submitLoading.value = true
  try {
    if (addLinkParams.value.address === initialLinkParams.address) {
      message.error('至少修改下信息吧~')
      return
    }
    if (addLinkList.value.some((item) => item.address === addLinkParams.value.address)) {
      message.error('已发送电波~ 请等待')
      return
    }

    await formRef.value!.validateFields()
    const res = await reqLinkAdd(addLinkParams.value)
    if (res.code === 200) {
      addLinkList.value.push(addLinkParams.value)
      message.info('已发送电波~ 请等待')
    }
  } finally {
    submitLoading.value = false
  }
}

const getLinkList = async () => {
  const res = await reqLinkList()
  linkList.value = res.data
}

const getArticle = async () => {
  const res = await reqLinkGetArticle()
  article.value = res.data
}

onMounted(async () => {
  await getLinkList()
  await getArticle()
})
</script>

<style scoped></style>
