<script setup lang="ts">
import CookiesBanner from '@/components/CookiesBanner.vue'
import Footer from '@/components/Footer.vue'
import Header from '@/components/Header.vue'
import { useScrollProgress } from '@/composables/useScrollProgress'

const { progress } = useScrollProgress()
</script>

<template>
  <div class="flex min-h-screen flex-col">
    <!-- Scroll progress bar -->
    <div
      class="fixed top-0 left-0 h-[2px] bg-accent z-[60] origin-left transition-none"
      :style="{ transform: `scaleX(${progress})` }"
    />
    <Header />
    <main>
      <router-view />
    </main>
    <Footer />
    <CookiesBanner />
  </div>
</template>
