import { ref, computed } from 'vue'
import { defineStore } from 'pinia'

export const useAuthStore = defineStore('auth', () => {
  const localStorageToken = localStorage.getItem('token')
  const token = ref<string | null>(
    localStorageToken != '' ? localStorageToken : null
  )

  const setToken = (newToken: string) => {
    token.value = newToken

    // Save the token in localStorage
    localStorage.setItem('token', newToken)
  }

  const logout = () => {
    token.value = null

    // Remove the token from localStorage
    localStorage.removeItem('token')
  }

  return { token, setToken, logout }
})
