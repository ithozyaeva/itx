<script setup lang="ts">
import Footer from '@/components/Footer.vue'
import Header from '@/components/Header.vue'
import FAQSection from '@/sections/FAQSection.vue'
import InviteSection from '@/sections/InviteSection.vue'
import MentorsSection from '@/sections/MentorsSection.vue'
import PromoteSection from '@/sections/PromoteSection.vue'
import ReviewSection from '@/sections/ReviewSection.vue'
import TariffSection from '@/sections/TariffSection.vue'
import WhySection from '@/sections/WhySection.vue'
</script>

<template>
  <div class="flex min-h-screen flex-col">
    <Header />
    <main>
      <PromoteSection />
      <WhySection />
      <MentorsSection />
      <TariffSection />
      <ReviewSection />
      <FAQSection />
      <InviteSection />
    </main>
    <Footer />
  </div>
</template>
