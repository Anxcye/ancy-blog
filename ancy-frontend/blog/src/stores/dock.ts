import { HomeOutlined } from '@ant-design/icons-vue'
import { defineStore } from 'pinia'
import { h, ref } from 'vue'

export const useDockStore = defineStore('dock', () => {
  const items = [
    {
      key: 'home',
      label: '首页',
      icon: h(HomeOutlined),
      path: '/',
    },
    {
      key: 'article',
      label: '文章',
      icon: h(HomeOutlined),
      path: '/article',
    },
    {
      key: 'timeline',
      label: '回溯',
      icon: h(HomeOutlined),
      path: '/timeline',
    },
  ]

  const homeArticles = ref([])

  return {
    items,
  }
})
