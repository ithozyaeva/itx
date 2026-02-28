<script setup lang="ts">
import Header from './Header.vue'
import Sidebar from './Sidebar.vue'
</script>

<template>
  <div class="flex h-screen">
    <Sidebar />

    <div class="flex-1 flex flex-col overflow-hidden min-w-0">
      <Header />

      <main class="flex-1 overflow-y-auto p-4 lg:p-6">
        <slot />
      </main>
    </div>
  </div>
</template>
