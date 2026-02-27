import type { NavigationGuardNext, RouteLocationNormalized, RouteRecordRaw } from 'vue-router'
import Events from '@/pages/Events.vue'
import ReferalLinks from '@/pages/ReferalLinks.vue'
import Resumes from '@/pages/Resumes.vue'
import Home from '@/pages/User.vue'
import { createRouter, createWebHistory } from 'vue-router'
import { useMainStore } from '../store'

const routes: RouteRecordRaw[] = [
  { path: '/me', component: Home, name: 'profile' },
  { path: '/events', component: Events, name: 'events' },
  { path: '/referals', component: ReferalLinks, name: 'referals' },
  { path: '/resumes', component: Resumes, name: 'resumes' },
  { path: '/', redirect: { name: 'events' } },
]

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes,
})

router.beforeEach((to: RouteLocationNormalized, from: RouteLocationNormalized, next: NavigationGuardNext) => {
  const store = useMainStore()

  if (!store.user) {
    store.initFromLocalStorage()
  }

  if (to.meta.requiresAuth && !store.user) {
    next('/login')
  }
  else {
    next()
  }
})

export default router
