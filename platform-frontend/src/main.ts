import { VueQueryPlugin } from '@tanstack/vue-query'
import { createPinia } from 'pinia'
import { createApp } from 'vue'
import App from './App.vue'
import { installMockTelegram } from './composables/useTelegramWebApp'
import { queryClient } from './plugins/vueQuery'
import { initWebVitals } from './plugins/webVitals'
import router from './router'
import './index.css'

if (import.meta.env.DEV && import.meta.env.VITE_MOCK_TELEGRAM === 'true')
  installMockTelegram()

const app = createApp(App)
app.use(createPinia())
app.use(router)
app.use(VueQueryPlugin, { queryClient })
app.mount('#app')

initWebVitals()
