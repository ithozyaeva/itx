<script setup lang="ts">
import Breadcrumbs from '@/components/common/Breadcrumbs.vue'
import Footer from './Footer.vue'
import Header from './Header.vue'
import Sidebar from './Sidebar.vue'
</script>

<template>
  <div class="flex h-screen">
    <Sidebar />
    <div class="flex-1 flex flex-col overflow-y-auto overflow-x-hidden">
      <Header />
      <main class="flex-grow w-full max-w-[1440px] mx-auto p-1">
        <Breadcrumbs />
        <slot />
      </main>
      <Footer />
    </div>
  </div>
</template>
