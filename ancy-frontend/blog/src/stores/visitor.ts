import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useVisitorStore = defineStore('visitor', () => {
  const user = ref<{
    nickname?: string
    email?: string
    avatar?: string
  }>(JSON.parse(localStorage.getItem('user') || '{}'))

  const liked = ref<Set<number>>(new Set(JSON.parse(localStorage.getItem('liked') || '[]')))

  const getUserInfo = () => {
    return user.value
  }

  const setUserInfo = (info: { nickname?: string; email?: string; avatar?: string }) => {
    user.value = info
    localStorage.setItem('user', JSON.stringify(user.value))
  }

  const checkLiked = (id: number) => {
    return liked.value.has(id)
  }

  const addLiked = (id: number) => {
    liked.value.add(id)
    localStorage.setItem('liked', JSON.stringify(Array.from(liked.value)))
  }

  const removeLiked = (id: number) => {
    liked.value.delete(id)
    localStorage.setItem('liked', JSON.stringify(Array.from(liked.value)))
  }

  return {
    user,
    getUserInfo,
    setUserInfo,
    checkLiked,
    addLiked,
    removeLiked,
  }
})
