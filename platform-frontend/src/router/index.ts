import type { NavigationGuardNext, RouteLocationNormalized, RouteRecordRaw } from 'vue-router'
import { createRouter, createWebHistory } from 'vue-router'
import Achievements from '@/pages/Achievements.vue'
import AutoApplyBot from '@/pages/AutoApplyBot.vue'
import Casino from '@/pages/Casino.vue'
import Content from '@/pages/Content.vue'
import Dashboard from '@/pages/Dashboard.vue'
import Events from '@/pages/Events.vue'
import Guilds from '@/pages/Guilds.vue'
import Kudos from '@/pages/Kudos.vue'
import Leaderboard from '@/pages/Leaderboard.vue'
import Marketplace from '@/pages/Marketplace.vue'
import MemberProfile from '@/pages/MemberProfile.vue'
import MentorProfile from '@/pages/MentorProfile.vue'
import Mentors from '@/pages/Mentors.vue'
import MyPoints from '@/pages/MyPoints.vue'
import MyReviews from '@/pages/MyReviews.vue'
import MyStats from '@/pages/MyStats.vue'
import NotificationSettings from '@/pages/NotificationSettings.vue'
import Quests from '@/pages/Quests.vue'
import Raffles from '@/pages/Raffles.vue'
import ReferalLinks from '@/pages/ReferalLinks.vue'
import Resumes from '@/pages/Resumes.vue'
import Seasons from '@/pages/Seasons.vue'
import TaskExchange from '@/pages/TaskExchange.vue'
import Home from '@/pages/User.vue'
import { useMainStore } from '../store'

const routes: RouteRecordRaw[] = [
  { path: '/', component: Dashboard, name: 'dashboard', meta: { requiresAuth: true } },
  { path: '/me', component: Home, name: 'profile', meta: { requiresAuth: true, breadcrumb: [{ label: 'Мой профиль' }] } },
  { path: '/events', component: Events, name: 'events', meta: { requiresAuth: true, breadcrumb: [{ label: 'События' }] } },
  { path: '/content', component: Content, name: 'content', meta: { requiresAuth: true, breadcrumb: [{ label: 'Контент' }] } },
  { path: '/members/:id', component: MemberProfile, name: 'memberProfile', meta: { requiresAuth: true, breadcrumb: [{ label: 'Профиль участника' }] } },
  { path: '/mentors', component: Mentors, name: 'mentors', meta: { requiresAuth: true, breadcrumb: [{ label: 'Менторы' }] } },
  { path: '/mentors/:id', component: MentorProfile, name: 'mentorProfile', meta: { requiresAuth: true, breadcrumb: [{ label: 'Менторы', to: '/mentors' }, { label: 'Профиль ментора' }] } },
  { path: '/referals', component: ReferalLinks, name: 'referals', meta: { requiresAuth: true, breadcrumb: [{ label: 'Рефералы' }] } },
  { path: '/resumes', component: Resumes, name: 'resumes', meta: { requiresAuth: true, breadcrumb: [{ label: 'Резюме' }] } },
  { path: '/my-reviews', component: MyReviews, name: 'myReviews', meta: { requiresAuth: true, breadcrumb: [{ label: 'Мои отзывы' }] } },
  { path: '/points', component: MyPoints, name: 'myPoints', meta: { requiresAuth: true, breadcrumb: [{ label: 'Мои баллы' }] } },
  { path: '/leaderboard', component: Leaderboard, name: 'leaderboard', meta: { requiresAuth: true, breadcrumb: [{ label: 'Рейтинг' }] } },
  { path: '/achievements', component: Achievements, name: 'achievements', meta: { requiresAuth: true, breadcrumb: [{ label: 'Достижения' }] } },
  { path: '/marketplace', component: Marketplace, name: 'marketplace', meta: { requiresAuth: true, breadcrumb: [{ label: 'Барахолка' }] } },
  { path: '/tasks', component: TaskExchange, name: 'taskExchange', meta: { requiresAuth: true, breadcrumb: [{ label: 'Биржа заданий' }] } },
  { path: '/quests', component: Quests, name: 'quests', meta: { requiresAuth: true, breadcrumb: [{ label: 'Квесты' }] } },
  { path: '/auto-apply', component: AutoApplyBot, name: 'autoApplyBot', meta: { requiresAuth: true, breadcrumb: [{ label: 'Автоотклики' }] } },
  { path: '/kudos', component: Kudos, name: 'kudos', meta: { requiresAuth: true, breadcrumb: [{ label: 'Благодарности' }] } },
  { path: '/seasons', component: Seasons, name: 'seasons', meta: { requiresAuth: true, breadcrumb: [{ label: 'Сезоны' }] } },
  { path: '/raffles', component: Raffles, name: 'raffles', meta: { requiresAuth: true, breadcrumb: [{ label: 'Розыгрыши' }] } },
  { path: '/casino', component: Casino, name: 'casino', meta: { requiresAuth: true, breadcrumb: [{ label: 'Казино' }] } },
  { path: '/guilds', component: Guilds, name: 'guilds', meta: { requiresAuth: true, breadcrumb: [{ label: 'Гильдии' }] } },
  { path: '/my-stats', component: MyStats, name: 'myStats', meta: { requiresAuth: true, breadcrumb: [{ label: 'Моя статистика' }] } },
  { path: '/notifications', component: NotificationSettings, name: 'notifications', meta: { requiresAuth: true, breadcrumb: [{ label: 'Уведомления' }] } },
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

  if (!to.meta.requiresAuth || store.user) {
    next()
  }
  else {
    next(false)
  }
})

export default router
