<script lang="ts" setup>
import { ref } from 'vue'
import axios from 'axios'

const latency = ref('0ms')
const url = ref('')
const success = ref(false)
const crawl = () => {
    const start = Date.now()
    axios.post(`http://localhost:8080/crawl?url=${url.value}`)
        .then((response: any) => {
            success.value = true
            latency.value = `${Date.now() - start}ms`
        })
}
</script>
<template>

    <div class="flex gap-4">
        <label class="input input-bordered flex items-center gap-2 w-full">
            <input type="text" class="grow" placeholder="Enter a URL" v-model="url" @keyup.enter="crawl" />
            <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5"
                stroke="currentColor" class="w-6 h-6">
                <path stroke-linecap="round" stroke-linejoin="round"
                    d="M15.042 21.672 13.684 16.6m0 0-2.51 2.225.569-9.47 5.227 7.917-3.286-.672Zm-7.518-.267A8.25 8.25 0 1 1 20.25 10.5M8.288 14.212A5.25 5.25 0 1 1 17.25 10.5" />
            </svg>
        </label>

        <button class="btn btn-primary" @click="crawl">Crawl</button>
    </div>
</template>