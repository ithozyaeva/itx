<script setup lang="ts">
import Footer from './Footer.vue'
import Header from './Header.vue'
import Sidebar from './Sidebar.vue'
</script>

<template>
  <div class="flex h-screen">
    <Sidebar />
    <div class="flex-1 flex flex-col overflow-auto">
      <Header />
      <main class="flex-grow w-full max-w-[1440px] mx-auto p-1">
        <slot />
      </main>
      <Footer />
    </div>
  </div>
</template>
