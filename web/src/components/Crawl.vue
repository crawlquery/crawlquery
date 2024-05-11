<script lang="ts" setup>
import { ref } from 'vue'
import axios from 'axios'

const latency = ref('0ms')
const url = ref('')
const success = ref(false)
const err = ref('')
const crawl = () => {
    const start = Date.now()
    axios.post(`http://localhost:8080/crawl?url=${url.value}`)
        .then((response: any) => {
            success.value = true
            latency.value = `${Date.now() - start}ms`
        })
        .catch((error: any) => {
            err.value = error.response.data.error
        })
}
</script>
<template>
    <div v-show="err" role="alert" class="mb-8 alert alert-warning">
        <svg xmlns="http://www.w3.org/2000/svg" class="stroke-current shrink-0 h-6 w-6" fill="none" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
        </svg>
        <span>{{ err }}!</span>
    </div>
    <div class="flex gap-4">
        <label class="input input-bordered flex items-center gap-2 w-full">
            <input type="text" class="grow" placeholder="Enter a URL to crawl" v-model="url" @keyup.enter="crawl" />
            <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5"
                stroke="currentColor" class="w-6 h-6">
                <path stroke-linecap="round" stroke-linejoin="round"
                    d="M15.042 21.672 13.684 16.6m0 0-2.51 2.225.569-9.47 5.227 7.917-3.286-.672Zm-7.518-.267A8.25 8.25 0 1 1 20.25 10.5M8.288 14.212A5.25 5.25 0 1 1 17.25 10.5" />
            </svg>
        </label>

        <button class="btn btn-primary" @click="crawl">Crawl</button>
    </div>
</template>