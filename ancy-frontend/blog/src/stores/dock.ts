import { reqArticleGetFrontList } from '@/api/article'
import type { ArticleListData } from '@/api/article/type'
import { reqCategoryList } from '@/api/category'
import type { CategoryListData } from '@/api/category/type'
import {
  BookOutlined,
  HomeOutlined,
  CalendarOutlined,
  LinkOutlined,
  MoreOutlined,
  UserOutlined,
  ProjectOutlined,
  ReadOutlined,
} from '@ant-design/icons-vue'
import { defineStore } from 'pinia'
import { computed, h, onMounted, ref } from 'vue'

export const useDockStore = defineStore('dock', () => {
  const homeArticles = ref<ArticleListData[]>([])
  const category = ref<CategoryListData[]>([])
  const homeItems = computed(() => {
    return homeArticles.value.map((item) => {
      return {
        key: item.id,
        label: item.title,
        icon: h(HomeOutlined),
        path: `/home/${item.id}`,
        group: 'home',
      }
    })
  })
  const categoryItems = computed(() => {
    return category.value.map((item) => {
      return {
        key: item.id,
        label: item.name,
        icon: h(HomeOutlined),
        path: `/category/${item.id}`,
        group: 'article',
      }
    })
  })

  const items = [
    {
      key: 'home',
      label: '首页',
      icon: h(HomeOutlined),
      path: '/',
      group: 'home',
      children: homeItems,
    },
    {
      key: 'article',
      label: '文章',
      icon: h(BookOutlined),
      path: '/article',
      group: 'article',
      children: categoryItems,
    },
    {
      key: 'timeline',
      label: '回溯',
      icon: h(CalendarOutlined),
      path: '/timeline',
      group: 'timeline',
    },
    {
      key: 'note',
      label: '日志',
      icon: h(CalendarOutlined),
      path: '/note',
      group: 'note',
    },
    {
      key: 'link',
      label: '友链',
      icon: h(LinkOutlined),
      path: '/link',
      group: 'link',
    },
    {
      key: 'more',
      label: '更多',
      icon: h(MoreOutlined),
      path: '/more',
      group: 'more',
      children: [
        {
          key: 'project',
          label: '项目',
          icon: h(ProjectOutlined),
          path: '/project',
          group: 'more',
        },
        {
          key: 'read',
          label: '阅读',
          icon: h(ReadOutlined),
          path: '/read',
          group: 'more',
        },
      ],
    },
  ]

  const reqHomeArticles = async () => {
    const res = await reqArticleGetFrontList()
    homeArticles.value = res.data
  }

  const reqCategory = async () => {
    const res = await reqCategoryList()
    category.value = res.data
  }

  onMounted(async () => {
    await reqHomeArticles()
    await reqCategory()
  })

  return {
    items,
    homeItems,
    reqHomeArticles,
  }
})
