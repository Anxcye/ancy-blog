import { defineStore } from 'pinia'
import { useFavicon } from '@vueuse/core'
import { useTitle } from '@vueuse/core'
import { useBaseInfoStore } from './baseInfo'

export const useBrowserStore = defineStore('browser', () => {
  const baseInfoStore = useBaseInfoStore()
  const icon = useFavicon()
  const title = useTitle()

  const setTitle = (newTitle: string) => {
    title.value = newTitle + ' | ' + (baseInfoStore.baseInfo.name ?? 'Ancy')
  }

  const setIcon = () => {
    icon.value = baseInfoStore.baseInfo.avatar
  }

  return { setIcon, setTitle }
})
