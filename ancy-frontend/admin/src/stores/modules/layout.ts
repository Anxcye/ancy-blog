import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useLayoutStore = defineStore('layout', () => {
  const refresh = ref(false)

  const setRefresh = () => {
    refresh.value = !refresh.value
  }

  return {
    refresh,
    setRefresh,
  }
})
