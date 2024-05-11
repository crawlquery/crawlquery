<script setup lang="ts">
import { RouterView } from 'vue-router'
</script>

<template>
  <div class="max-w-5xl mx-auto mt-16">
    <div class="md:w-1/3 mx-auto">
      <a href="/">
        <img src="/logo.svg" alt="Logo" class="w-full" />
      </a>
    </div>
  </div>
  <RouterView />
  <div class="w-full text-center mt-16 text-gray-300">
    <div>CrawlQuery 2024, all rights reserved.
    </div>
    <div class="mt-8">
      <RouterLink to="/crawl" class="text-indigo-500 font-semibold">Crawl Page</RouterLink>
    </div>
  </div>
</template>