<script lang="ts" setup>
import { ref } from 'vue'
import axios from 'axios'
import { useAuthStore } from '../stores/auth'

const emit = defineEmits(['nodeCreated'])
const authStore = useAuthStore()

const latency = ref('0ms')
const hostname = ref('')
const port = ref('')
const success = ref(false)
const err = ref('')
const addNode = () => {
    const start = Date.now()
    axios.post(`http://localhost:8080/nodes`, {
        hostname: hostname.value,
        port: port.value,
    }, {
        headers: {
            Authorization: `Bearer ${authStore.token}`
        }
    })
        .then((response: any) => {
            success.value = true
            latency.value = `${Date.now() - start}ms`
            emit('nodeCreated', response.data.node)
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
    <div v-if="success" role="alert" class="mb-8 alert alert-success">
        <svg xmlns="http://www.w3.org/2000/svg" class="stroke-current shrink-0 h-6 w-6" fill="none" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
        </svg>
        <span>Node added successfully!</span>
    </div>
    <div class="flex gap-4">
        <label class="input input-bordered flex items-center gap-2 w-full">
            <input type="text" class="grow" placeholder="hostname e.g node1.cluster.com" v-model="hostname" />
            <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5"
                stroke="currentColor" class="w-6 h-6">
                <path stroke-linecap="round" stroke-linejoin="round"
                    d="M5.25 14.25h13.5m-13.5 0a3 3 0 0 1-3-3m3 3a3 3 0 1 0 0 6h13.5a3 3 0 1 0 0-6m-16.5-3a3 3 0 0 1 3-3h13.5a3 3 0 0 1 3 3m-19.5 0a4.5 4.5 0 0 1 .9-2.7L5.737 5.1a3.375 3.375 0 0 1 2.7-1.35h7.126c1.062 0 2.062.5 2.7 1.35l2.587 3.45a4.5 4.5 0 0 1 .9 2.7m0 0a3 3 0 0 1-3 3m0 3h.008v.008h-.008v-.008Zm0-6h.008v.008h-.008v-.008Zm-3 6h.008v.008h-.008v-.008Zm0-6h.008v.008h-.008v-.008Z" />
            </svg>
        </label>
        <input type="number" class="input input-bordered w-1/4" placeholder="Port" v-model="port" />
    </div>
    <button class="mt-8 btn btn-primary w-full" @click="addNode">Add</button>
</template>