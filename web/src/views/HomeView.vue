<script setup lang="ts">
import Search from '../components/Search.vue'
</script>

<template>
  <main>
    <div class="max-w-5xl mx-auto mt-16">
      <div class="md:w-1/3 mx-auto">
        <a href="/">
          <img src="/logo.svg" alt="Logo" class="w-full" />
        </a>
      </div>
      <div class="max-w-xl mx-auto mt-24">
        <Search />
      </div>
    </div>

    <div class="w-full text-center mt-16 text-gray-300">
      crawlquery 2024, all rights reserved.
    </div>
  </main>
</template>
