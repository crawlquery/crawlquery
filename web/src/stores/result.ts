import { ref, computed } from 'vue'
import { defineStore } from 'pinia'
import type { Result } from '@/types/types'

export const useResultStore = defineStore('result', () => {
  const results = ref<Result[]>([])

  const setResults = (newResults: Result[]) => {
    results.value = newResults
  }

  return { results, setResults }
})
