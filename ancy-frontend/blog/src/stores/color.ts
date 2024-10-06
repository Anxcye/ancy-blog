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
  const isDarkMode = ref(false)

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
    updateTheme()
  }

  const toggleDarkMode = () => {
    isDarkMode.value = !isDarkMode.value
    updateTheme()
  }

  const updateTheme = () => {
    if (isDarkMode.value) {
      document.documentElement.classList.add('dark-mode')
    } else {
      document.documentElement.classList.remove('dark-mode')
    }
  }

  onMounted(() => {
    setPrimaryColor(getPrimaryColor())
    // 可以根据用户的系统偏好设置初始主题
    isDarkMode.value = window.matchMedia('(prefers-color-scheme: dark)').matches
    updateTheme()
  })

  return { currentColor, getPrimaryColor, resetPrimaryColor, isDarkMode, toggleDarkMode }
})
