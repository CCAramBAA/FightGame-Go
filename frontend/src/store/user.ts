import { defineStore } from 'pinia'
import { ref, computed } from 'vue'

export const useUserStore = defineStore('user', () => {
  const token = ref<string>(localStorage.getItem('token') || '')
  const username = ref<string>(localStorage.getItem('username') || '')
  const userId = ref<string>(localStorage.getItem('userId') || '')

  const isLoggedIn = computed(() => !!token.value)

  function setAuth(t: string, name: string, id: string) {
    token.value = t
    username.value = name
    userId.value = id
    localStorage.setItem('token', t)
    localStorage.setItem('username', name)
    localStorage.setItem('userId', id)
  }

  function logout() {
    token.value = ''
    username.value = ''
    userId.value = ''
    localStorage.removeItem('token')
    localStorage.removeItem('username')
    localStorage.removeItem('userId')
  }

  return { token, username, userId, isLoggedIn, setAuth, logout }
})
