import { defineStore } from 'pinia'
import { onMounted, ref } from 'vue'

const colorScheme = [
  '#66ccff',
  '#00b96b',
  '#ff9900',
  '#ff4d4f',
  '#9933ff',
  '#ff6633',
  '#33cc33',
  '#ff3333',
  '#3333ff',
]

export const useColorStore = defineStore('color', () => {
  const currentColor = ref<string | null>(null)

  const getPrimaryColor = () => {
    if (currentColor.value === null) {
      currentColor.value = colorScheme[Math.floor(Math.random() * colorScheme.length)]
    }
    return currentColor.value
  }

  const resetPrimaryColor = () => {
    currentColor.value = null
  }

  const setPrimaryColor = (color: string) => {
    document.documentElement.style.setProperty('--primary-color', color)
  }

  onMounted(() => {
    setPrimaryColor(getPrimaryColor())
  })

  return { currentColor, getPrimaryColor, resetPrimaryColor }
})
