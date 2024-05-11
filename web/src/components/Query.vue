<script lang="ts" setup>
import { ref } from 'vue'
import axios from 'axios'

const latency = ref('0ms')
const statement = ref('')
const success = ref(false)
const query = () => {
    const start = Date.now()
    axios.post(`http://localhost:8080/query?s=${statement.value}`)
        .then((response: any) => {
            success.value = true
            latency.value = `${Date.now() - start}ms`
        })
}
</script>
<template>

    <div class="flex gap-4">
        <textarea type="text" class="w-full focus:outline-none textarea textarea-bordered textarea-lg w-full"
            placeholder="select * from pages where title = 'home';" rows="4" v-model="statement"
            @keyup.enter="query"> </textarea>

        <button class="btn btn-primary" @click="query">Query</button>
    </div>
</template>