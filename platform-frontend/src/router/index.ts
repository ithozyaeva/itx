import type { RouteRecordRaw } from 'vue-router'
import { createRouter, createWebHistory } from 'vue-router'
import Achievements from '@/pages/Achievements.vue'
import AutoApplyBot from '@/pages/AutoApplyBot.vue'
import Casino from '@/pages/Casino.vue'

import Dashboard from '@/pages/Dashboard.vue'
import Events from '@/pages/Events.vue'
import Kudos from '@/pages/Kudos.vue'
import Leaderboard from '@/pages/Leaderboard.vue'
import Marketplace from '@/pages/Marketplace.vue'
import MemberProfile from '@/pages/MemberProfile.vue'
import MentorProfile from '@/pages/MentorProfile.vue'
import Mentors from '@/pages/Mentors.vue'
import MyPoints from '@/pages/MyPoints.vue'
import MyReviews from '@/pages/MyReviews.vue'
import MyStats from '@/pages/MyStats.vue'

import Quests from '@/pages/Quests.vue'
import Raffles from '@/pages/Raffles.vue'
import ReferalLinks from '@/pages/ReferalLinks.vue'
import Resumes from '@/pages/Resumes.vue'
import Seasons from '@/pages/Seasons.vue'
import TaskExchange from '@/pages/TaskExchange.vue'
import Home from '@/pages/User.vue'

const routes: RouteRecordRaw[] = [
  { path: '/', component: Dashboard, name: 'dashboard' },
  { path: '/me', component: Home, name: 'profile', meta: { breadcrumb: [{ label: 'Мой профиль' }] } },
  { path: '/events', component: Events, name: 'events', meta: { breadcrumb: [{ label: 'События' }] } },
  { path: '/content', redirect: '/events?tab=content' },
  { path: '/members/:id', component: MemberProfile, name: 'memberProfile', meta: { breadcrumb: [{ label: 'Рейтинг', to: '/leaderboard' }, { label: 'Профиль участника' }] } },
  { path: '/mentors', component: Mentors, name: 'mentors', meta: { breadcrumb: [{ label: 'Менторы' }] } },
  { path: '/mentors/:id', component: MentorProfile, name: 'mentorProfile', meta: { breadcrumb: [{ label: 'Менторы', to: '/mentors' }, { label: 'Профиль ментора' }] } },
  { path: '/referals', component: ReferalLinks, name: 'referals', meta: { breadcrumb: [{ label: 'Рефералы' }] } },
  { path: '/resumes', component: Resumes, name: 'resumes', meta: { breadcrumb: [{ label: 'Резюме' }] } },
  { path: '/my-reviews', component: MyReviews, name: 'myReviews', meta: { breadcrumb: [{ label: 'Мои отзывы' }] } },
  { path: '/points', component: MyPoints, name: 'myPoints', meta: { breadcrumb: [{ label: 'Мои баллы' }] } },
  { path: '/leaderboard', component: Leaderboard, name: 'leaderboard', meta: { breadcrumb: [{ label: 'Рейтинг' }] } },
  { path: '/achievements', component: Achievements, name: 'achievements', meta: { breadcrumb: [{ label: 'Достижения' }] } },
  { path: '/marketplace', component: Marketplace, name: 'marketplace', meta: { breadcrumb: [{ label: 'Барахолка' }] } },
  { path: '/tasks', component: TaskExchange, name: 'taskExchange', meta: { breadcrumb: [{ label: 'Биржа заданий' }] } },
  { path: '/quests', component: Quests, name: 'quests', meta: { breadcrumb: [{ label: 'Квесты' }] } },
  { path: '/auto-apply', component: AutoApplyBot, name: 'autoApplyBot', meta: { breadcrumb: [{ label: 'Автоотклики' }] } },
  { path: '/kudos', component: Kudos, name: 'kudos', meta: { breadcrumb: [{ label: 'Благодарности' }] } },
  { path: '/seasons', component: Seasons, name: 'seasons', meta: { breadcrumb: [{ label: 'Сезоны' }] } },
  { path: '/raffles', component: Raffles, name: 'raffles', meta: { breadcrumb: [{ label: 'Розыгрыши' }] } },
  { path: '/minigames', component: Casino, name: 'minigames', meta: { breadcrumb: [{ label: 'Мини-игры' }] } },
  { path: '/my-stats', component: MyStats, name: 'myStats', meta: { breadcrumb: [{ label: 'Моя статистика' }] } },
  { path: '/notifications', redirect: '/me' },
  { path: '/faq', name: 'faq', component: () => import('@/pages/FAQ.vue'), meta: { breadcrumb: [{ label: 'FAQ' }] } },
  { path: '/:pathMatch(.*)*', redirect: '/' },
]

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes,
})

export default router
