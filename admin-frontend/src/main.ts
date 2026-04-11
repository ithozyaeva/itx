import { VueQueryPlugin } from '@tanstack/vue-query'
import { createPinia } from 'pinia'

import { createApp } from 'vue'
import App from './App.vue'
import permissionDirective from './directives/permission'
import { queryClient } from './plugins/vueQuery'
import router from './router'
import './assets/main.css'

const app = createApp(App)
app.use(createPinia())
app.use(router)
app.use(VueQueryPlugin, { queryClient })

// Register the permission directive globally
app.directive('permission', permissionDirective)

app.mount('#app')
