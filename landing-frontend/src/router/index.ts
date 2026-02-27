import type { RouteRecordRaw } from 'vue-router'
import { createRouter, createWebHistory } from 'vue-router'
import Home from '../pages/Main.vue'
import Privacy from '../pages/Privacy.vue'

const routes: RouteRecordRaw[] = [
  { path: '/', component: Home },
  { path: '/privacy', component: Privacy },
]

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes,
})

export default router
