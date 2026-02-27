import { VueQueryPlugin } from '@tanstack/vue-query'

import { createHead } from '@unhead/vue/client'
import MasonryWall from '@yeger/vue-masonry-wall'
import { createApp } from 'vue'
import { initYandexMetrika } from 'yandex-metrika-vue3'
import App from './App.vue'

import { queryClient } from './plugins/vueQuery'
import router from './router'
import 'itx-ui-kit/style.css'
import './assets/base.css'

const app = createApp(App)

const head = createHead()
app.use(head)

const metrikaId = import.meta.env.VITE_YANDEX_METRIKA_ID
const metrikaEnabled = import.meta.env.VITE_YANDEX_METRIKA_ENABLED !== 'false'

if (metrikaId && metrikaEnabled) {
  app.use(initYandexMetrika, {
    id: metrikaId,
    router,
    env: import.meta.env.MODE === 'development' ? 'production' : import.meta.env.MODE,
    scriptSrc: 'https://mc.yandex.ru/metrika/tag.js',
    options: {
      clickmap: true,
      trackLinks: true,
      accurateTrackBounce: true,
      webvisor: false,
    },
  })
}

app.use(MasonryWall)
app.use(router)
app.use(VueQueryPlugin, { queryClient })
app.mount('#app')
