<script lang="ts" setup>
import { computed } from 'vue'
import StructuredData from '@/components/StructuredData.vue'
import { usePageMeta } from '@/composables/useMeta'
import CalendarSection from '@/sections/CalendarSection.vue'
import FAQSection from '@/sections/FAQSection.vue'
import MentorsSection from '@/sections/MentorsSection.vue'
import PromoteSection from '@/sections/PromoteSection.vue'
import TariffSection from '@/sections/TariffSection.vue'
import WhySection from '@/sections/WhySection.vue'

usePageMeta({
  title: 'IT-ХОЗЯЕВА - закрытое сообщество IT-специалистов с менторами',
  description: 'Закрытое сообщество IT-специалистов с базой менторов. Помогаем с карьерой, поиском работы и собесами. Регулярные встречи, тренинги и менторство от опытных IT-профессионалов. Присоединяйтесь к сообществу!',
  keywords: 'IT сообщество, IT менторы, менторство, программисты, разработчики, IT специалисты, закрытое сообщество, карьера в IT, IT профессионалы, IT хозяева, сообщество разработчиков, IT сообщество России, база менторов',
  url: 'https://ithozyaeva.ru',
  type: 'website',
})

const structuredData = computed(() => [
  {
    '@context': 'https://schema.org',
    '@type': 'Organization',
    'name': 'IT-ХОЗЯЕВА',
    'description': 'Закрытое сообщество IT-специалистов, в котором ценят опыт, юмор и профессионализм',
    'url': 'https://ithozyaeva.ru',
    'logo': 'https://ithozyaeva.ru/favicon.svg',
  },
  {
    '@context': 'https://schema.org',
    '@type': 'WebSite',
    'name': 'IT-ХОЗЯЕВА',
    'url': 'https://ithozyaeva.ru',
    'description': 'Закрытое сообщество IT-специалистов',
  },
  {
    '@context': 'https://schema.org',
    '@type': 'Service',
    'serviceType': 'IT Сообщество',
    'provider': {
      '@type': 'Organization',
      'name': 'IT-ХОЗЯЕВА',
    },
    'areaServed': {
      '@type': 'Country',
      'name': 'Россия',
    },
    'description': 'Закрытое сообщество IT-специалистов с менторством, встречами, обсуждениями и помощью в карьере',
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
