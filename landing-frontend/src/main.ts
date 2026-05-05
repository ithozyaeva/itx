import { VueQueryPlugin } from '@tanstack/vue-query'

import { createHead } from '@unhead/vue/client'
import { MasonryWall } from '@yeger/vue-masonry-wall'
import { createApp } from 'vue'
import { initYandexMetrika } from 'yandex-metrika-vue3'
import App from './App.vue'

import { queryClient } from './plugins/vueQuery'
import { initWebVitals } from './plugins/webVitals'
import router from './router'
import '@fontsource-variable/inter/wght.css'
import '@fontsource-variable/jetbrains-mono/wght.css'
import '@fontsource-variable/space-grotesk/wght.css'
import './assets/uikit.css'
import './assets/base.css'

const app = createApp(App)

const head = createHead()
app.use(head)

app.component('MasonryWall', MasonryWall)
app.use(router)
app.use(VueQueryPlugin, { queryClient })
app.mount('#app')

initWebVitals()

// Метрика инициализируется после первого idle, чтобы не конкурировать с LCP
// и не утяжелять INP первой интеракции. Скрипт mc.yandex.ru/metrika/tag.js
// и подписка на route changes ставится только после mount-а.
const metrikaId = import.meta.env.VITE_YANDEX_METRIKA_ID
const metrikaEnabled = import.meta.env.VITE_YANDEX_METRIKA_ENABLED !== 'false'

if (metrikaId && metrikaEnabled) {
  const initMetrika = () => app.use(initYandexMetrika, {
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

  if ('requestIdleCallback' in window) {
    window.requestIdleCallback(initMetrika, { timeout: 4000 })
  }
  else {
    setTimeout(initMetrika, 2000)
  }
}
