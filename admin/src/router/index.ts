import { checkAuth, isAuthenticated } from '@/services/authService'
import { createRouter, createWebHistory } from 'vue-router'

const LoginView = () => import('@/views/LoginView.vue')
const DashboardView = () => import('@/views/DashboardView.vue')
const MentorsView = () => import('@/views/MentorsView.vue')
const MembersView = () => import('@/views/MembersView.vue')
const ReviewsView = () => import('@/views/ReviewsView.vue')
const MentorsReviewsView = () => import('@/views/MentorsReviewsView.vue')

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      redirect: to => {
        checkAuth()
        return isAuthenticated.value ? { name: 'dashboard' } : { name: 'login' }
      }
    },
    {
      path: '/login',
      name: 'login',
      component: LoginView,
    },
    {
      path: '/dashboard',
      name: 'dashboard',
      component: DashboardView,
      meta: { requiresAuth: true },
    },
    {
      path: '/mentors',
      name: 'mentors',
      component: MentorsView,
      meta: { requiresAuth: true },
    },
    {
      path: '/members',
      name: 'members',
      component: MembersView,
      meta: { requiresAuth: true },
    },
    {
      path: '/reviews',
      name: 'reviews',
      component: ReviewsView,
      meta: { requiresAuth: true },
    },
    {
      path: '/mentor-reviews',
      name: 'mentor-reviews',
      component: MentorsReviewsView,
      meta: { requiresAuth: true },
    },
  ],
})

router.beforeEach((to, _, next) => {
  checkAuth()

  if (to.meta.requiresAuth && !isAuthenticated.value) {
    next({ name: 'login', query: { redirect: to.fullPath } })
  }
  else if (to.name === 'login' && isAuthenticated.value) {
    next({ name: 'dashboard' })
  }
  else {
    next()
  }
})

export default router
