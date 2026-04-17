<script lang="ts" setup>
import { computed } from 'vue'
import StructuredData from '@/components/StructuredData.vue'
import { usePageMeta } from '@/composables/useMeta'
import { useRevealObserver } from '@/composables/useRevealObserver'
import CalendarSection from '@/sections/CalendarSection.vue'
import FAQSection from '@/sections/FAQSection.vue'
import MentorsSection from '@/sections/MentorsSection.vue'
import PromoteSection from '@/sections/PromoteSection.vue'
import TariffSection from '@/sections/TariffSection.vue'
import WhySection from '@/sections/WhySection.vue'

useRevealObserver()

usePageMeta({
  title: 'IT-ХОЗЯЕВА - закрытое сообщество IT-специалистов с менторами',
  description: 'Закрытое сообщество IT-специалистов с базой менторов. Помогаем с карьерой, поиском работы и собесами. Регулярные встречи, тренинги и менторство от опытных IT-профессионалов. Присоединяйтесь к сообществу!',
  url: 'https://ithozyaeva.ru',
  type: 'website',
})

const structuredData = computed(() => [
  {
    '@context': 'https://schema.org',
    '@type': 'WebSite',
    'name': 'IT-ХОЗЯЕВА',
    'url': 'https://ithozyaeva.ru',
    'description': 'IT-сообщество на стыке технологий и искусственного интеллекта: менторство, вайбкодинг, нетворкинг',
  },
  {
    '@context': 'https://schema.org',
    '@type': 'Service',
    'serviceType': 'IT-сообщество с AI и менторством',
    'provider': {
      '@type': 'Organization',
      'name': 'IT-ХОЗЯЕВА',
    },
    'areaServed': {
      '@type': 'Country',
      'name': 'Россия',
    },
    'description': 'IT-комьюнити: менторство, вайбкодинг (vibe coding), нетворкинг, AI-практики, подготовка к собеседованиям',
    'offers': {
      '@type': 'Offer',
      'priceCurrency': 'RUB',
      'price': '520',
      'description': 'Тарифы от 520₽ в месяц',
    },
  },
])
</script>

<template>
  <StructuredData :data="structuredData" />
  <PromoteSection />
  <WhySection />
  <MentorsSection />
  <CalendarSection />
  <TariffSection />
  <FAQSection />
</template>
