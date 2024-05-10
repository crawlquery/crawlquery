<script lang="ts" setup>
import { ref } from 'vue'
import axios from 'axios'
import { useResultStore } from '@/stores/result';

const resultStore = useResultStore()

const term = ref('')

const search = () => {
    axios.get(`http://localhost:8080/search?q=${term.value}`)
        .then((response: any) => {
            resultStore.setResults(response.data.results)
            console.log(response.data.results)
        })
}
</script>
<template>

    <div>
        <label class="input input-bordered flex items-center gap-2">
            <input type="text" class="grow" placeholder="Search" :model="term" @keyup="search" />
            <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 16 16" fill="currentColor" class="w-4 h-4 opacity-70">
                <path fill-rule="evenodd"
                    d="M9.965 11.026a5 5 0 1 1 1.06-1.06l2.755 2.754a.75.75 0 1 1-1.06 1.06l-2.755-2.754ZM10.5 7a3.5 3.5 0 1 1-7 0 3.5 3.5 0 0 1 7 0Z"
                    clip-rule="evenodd" />
            </svg>
        </label>

        <div v-if="resultStore.results && resultStore.results.length > 0">
            <div v-for="result in resultStore.results" :key="result.id">
                <div class="card">
                    <div class="card-body">
                        <h2 class="card-title">{{ result.title }}</h2>

                        <p class="text-gray-500">{{ result.description }}</p>

                    </div>
                </div>
            </div>
        </div>
    </div>
</template>