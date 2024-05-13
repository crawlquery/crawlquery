<script setup lang="ts">
import { RouterView } from 'vue-router'
import { useRoute, useRouter } from 'vue-router';
import { useAuthStore } from './stores/auth'
const route = useRoute()
const router = useRouter()

const authStore = useAuthStore()
</script>

<template>
  <div class="max-w-3xl mx-auto mt-16">
    <div class="w-1/3 mx-auto">
      <a href="/">
        <img src="/logo.svg" alt="Logo" class="w-full" />
      </a>
    </div>

    <div class="mt-8 navbar border rounded-xl">
      <div class="navbar-center hidden lg:flex">
        <ul class="menu menu-horizontal px-1">
          <li>
            <RouterLink to="/" class="font-semibold"
              :class="route.name == 'search' ? 'text-indigo-300' : 'text-indigo-500 '">
              Search</RouterLink>
          </li>
          <li v-if="authStore.token">
            <RouterLink to="/nodes" class="ml-4 font-semibold"
              :class="route.name == 'nodes' ? 'text-indigo-300' : 'text-indigo-500 '">Nodes</RouterLink>
          </li>
          <li>
            <RouterLink to="/crawl" class="ml-4 font-semibold"
              :class="route.name == 'crawl' ? 'text-indigo-300' : 'text-indigo-500 '">Crawl</RouterLink>
          </li>
          <li>
            <RouterLink to="/query" class="ml-4 font-semibold"
              :class="route.name == 'query' ? 'text-indigo-300' : 'text-indigo-500 '">Query</RouterLink>
          </li>
        </ul>
      </div>
      <div class="grow"></div>
      <div class="navbar-end">
        <button v-if="authStore.token" @click="authStore.logout() && router.push('/')"
          class="btn ml-4 font-semibold text-indigo-500">Logout</button>
        <button v-else @click="router.push('/login')" class="btn ml-4 font-semibold text-indigo-500">Login</button>
      </div>
    </div>
  </div>
  <RouterView />
  <div class="w-full text-center mt-16 text-gray-300">
    <div>crawlquery - all rights reserved.
    </div>
  </div>

</template>