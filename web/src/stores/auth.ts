import { ref, computed } from 'vue'
import { defineStore } from 'pinia'

export const useAuthStore = defineStore('auth', () => {
  const token = ref<string | null>(null)

  const setToken = (newToken: string) => {
    token.value = newToken
  }

  const logout = () => {
    token.value = null
  }

  return { token, setToken, logout }
})
