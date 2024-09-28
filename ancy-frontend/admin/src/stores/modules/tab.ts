import router from '@/router'
import { defineStore, type StoreDefinition } from 'pinia'
import { ref } from 'vue'
import type { RouteLocationNormalized } from 'vue-router'

export const useTabStore: StoreDefinition = defineStore('tab', () => {
  const historyTabs = ref<RouteLocationNormalized[]>([])
  const currentTab = ref<RouteLocationNormalized | null>(null)
  const cacheTabs = ref<string[]>([])

  const addHistoryTab = (tab: RouteLocationNormalized) => {
    if (historyTabs.value.some((item) => item.path === tab.path)) {
      return
    }
    historyTabs.value.push(tab)
  }

  const removeHistoryTab = (tab: RouteLocationNormalized) => {
    historyTabs.value = historyTabs.value.filter(
      (item) => item.path !== tab.path,
    )
    if (historyTabs.value.length === 0) {
      router.push('/ancy')
    }

    if (currentTab.value!.path === tab.path) {
      currentTab.value = historyTabs.value[historyTabs.value.length - 1]
      router.push(currentTab.value.path)
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
