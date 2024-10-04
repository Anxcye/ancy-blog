import { HomeOutlined } from '@ant-design/icons-vue'
import { defineStore } from 'pinia'
import { h } from 'vue'

export const useDockStore = defineStore('dock', () => {
  const items = [
    {
      key: 'home1',
      label: '首页1',
      icon: h(HomeOutlined),
      path: '/',
    },
    {
      key: 'home2',
      label: '首页2',
      icon: h(HomeOutlined),
      path: '/article/3',
    },
    {
      key: 'home3',
      label: '首页3',
      icon: h(HomeOutlined),
      path: '/article/5',
    },
  ]

  return {
    items,
  }
})
