<script setup lang="ts">
import { RouterView } from 'vue-router'
import { useRoute } from 'vue-router';
const route = useRoute()
</script>

<template>
  <div class="max-w-5xl mx-auto mt-16">
    <div class="w-1/3 mx-auto">
      <a href="/">
        <img src="/logo.svg" alt="Logo" class="w-full" />
      </a>
    </div>
  </div>
  <div class="w-full text-center mt-16 text-gray-300">
    <div class="mt-8">
      <RouterLink to="/" class="font-semibold" :class="route.name == 'search' ? 'text-indigo-300' : 'text-indigo-500 '">
        Search</RouterLink>
      <RouterLink to="/crawl" class="ml-4 font-semibold"
        :class="route.name == 'crawl' ? 'text-indigo-300' : 'text-indigo-500 '">Crawl</RouterLink>
      <RouterLink to="/query" class="ml-4 font-semibold"
        :class="route.name == 'query' ? 'text-indigo-300' : 'text-indigo-500 '">Query</RouterLink>
    </div>
  </div>
  <RouterView />
  <div class="w-full text-center mt-16 text-gray-300">
    <div>crawlquery - all rights reserved.
    </div>
  </div>

</template>