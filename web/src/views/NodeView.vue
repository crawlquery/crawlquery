<script setup lang="ts">
import AddNode from '../components/AddNode.vue'
import NodeList from '../components/NodeList.vue'
import { ref } from 'vue'
import { Node } from '../types'
import axios from 'axios'
import { useAuthStore } from '../stores/auth'

const authStore = useAuthStore()
const err = ref('')
const nodes = ref<Array<Node>>([])

axios.get(`http://localhost:8080/nodes`, {
  headers: {
    Authorization: `Bearer ${authStore.token}`
  }
})
  .then((response: any) => {
    nodes.value = response.data.nodes
  })
  .catch((error: any) => {
    err.value = error.response.data.error
  })
</script>

<template>
  <main>
    <div class="max-w-xl mx-auto mt-24">
      <div v-show="err" role="alert" class="mb-8 alert alert-warning">
        <svg xmlns="http://www.w3.org/2000/svg" class="stroke-current shrink-0 h-6 w-6" fill="none" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
            d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
        </svg>
        <span>{{ err }}!</span>
      </div>
      <NodeList :nodes="nodes" />
      <div class="mt-24">
      </div>
      <AddNode @nodeCreated="nodes.push($event)" />
    </div>
  </main>
</template>
