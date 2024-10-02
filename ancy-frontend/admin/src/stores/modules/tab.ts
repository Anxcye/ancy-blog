import router from '@/router'
import { defineStore, type StoreDefinition } from 'pinia'
import { ref } from 'vue'
import type { RouteRecordNormalized } from 'vue-router'

export const useTabStore: StoreDefinition = defineStore('tab', () => {
  const historyTabs = ref<RouteRecordNormalized[]>([])
  const currentTab = ref<string | null>(null)
  const cacheTabs = ref<string[]>([])

  const addHistoryTab = (tab: string) => {
    if (historyTabs.value.some((item) => item.path === tab)) {
      return
    }
    const route = router.getRoutes().find((item) => item.path === tab)
    if (route) {
      historyTabs.value.push(route)
    }
  }

  const removeHistoryTab = (tab: string) => {
    historyTabs.value = historyTabs.value.filter((item) => item.path !== tab)
    if (historyTabs.value.length === 0) {
      router.push('/')
      currentTab.value = '/'
    }

    if (currentTab.value === tab) {
      currentTab.value = historyTabs.value[historyTabs.value.length - 1].path
      router.push(currentTab.value)
    }
  }

  return {
    historyTabs,
    currentTab,
    cacheTabs,
    addHistoryTab,
    removeHistoryTab,
  }
})
