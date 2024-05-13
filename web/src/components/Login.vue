<script lang="ts" setup>
import { ref } from 'vue'
import axios from 'axios'
import { useAuthStore } from '../stores/auth'

const authStore = useAuthStore()

const latency = ref('0ms')
const email = ref('')
const password = ref('')
const success = ref('')
const err = ref('')

const signUp = () => {
    const start = Date.now()
    axios.post(`http://localhost:8080/accounts`, {
        email: email.value,
        password: password.value,
    })
        .then((response: any) => {
            err.value = ''
            success.value = "Account created successfully"
            latency.value = `${Date.now() - start}ms`
        })
        .catch((error: any) => {
            err.value = error.response.data.error
        })
}

const logIn = () => {
    const start = Date.now()
    axios.post(`http://localhost:8080/login`, {
        email: email.value,
        password: password.value,
    })
        .then((response: any) => {
            err.value = ''
            success.value = "Logged in successfully"
            authStore.setToken(response.data.token)
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
    <div v-if="success" role="alert" class="mb-8 alert alert-success">
        <svg xmlns="http://www.w3.org/2000/svg" class="stroke-current shrink-0 h-6 w-6" fill="none" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
        </svg>
        <span>{{ success }}</span>
    </div>

    <div v-if="!authStore.token" class="flex gap-4">
        <div class="grow"></div>
        <label class="input input-bordered flex items-center">
            <input type="text" placeholder="Email" v-model="email" />
        </label>
        <label class="input input-bordered flex items-centerw-full">
            <input type="password" placeholder="Password" v-model="password" />
        </label>
        <div class="grow"></div>
    </div>

    <div v-if="!authStore.token" class="flex gap-4">
        <button class="mt-8 btn btn-primary flex-1" style="flex-basis: 0; flex-grow: 1; flex-shrink: 0;"
            @click="signUp">
            Sign Up</button>
        <button class="mt-8 btn btn-primary flex-1" style="flex-basis: 0; flex-grow: 1; flex-shrink: 0;"
            @click="logIn">Log
            In</button>
    </div>
</template>