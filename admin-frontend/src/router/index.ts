import { createRouter, createWebHistory } from 'vue-router'
import { checkAuth, ensureAdminAccess, isAuthenticated, logout } from '@/services/authService'

const LoginView = () => import('@/views/LoginView.vue')
const DashboardView = () => import('@/views/DashboardView.vue')
const MentorsView = () => import('@/views/MentorsView.vue')
const MembersView = () => import('@/views/MembersView.vue')
const ReviewsView = () => import('@/views/ReviewsView.vue')
const EventsView = () => import('@/views/EventsView.vue')
const MentorsReviewsView = () => import('@/views/MentorsReviewsView.vue')
const ResumesView = () => import('@/views/ResumesView.vue')
const AuditLogsView = () => import('@/views/AuditLogsView.vue')
const PointsView = () => import('@/views/PointsView.vue')
const CreditsView = () => import('@/views/CreditsView.vue')
const ReferralsView = () => import('@/views/ReferralsView.vue')
const ChatActivityView = () => import('@/views/ChatActivityView.vue')
const ChatQuestsView = () => import('@/views/ChatQuestsView.vue')
const RafflesView = () => import('@/views/RafflesView.vue')
const CasinoView = () => import('@/views/CasinoView.vue')
const DailyTasksView = () => import('@/views/DailyTasksView.vue')
const ChallengesView = () => import('@/views/ChallengesView.vue')
const SubscriptionsView = () => import('@/views/SubscriptionsView.vue')
const FeedbackView = () => import('@/views/FeedbackView.vue')
const ModerationView = () => import('@/views/ModerationView.vue')

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      redirect: () => {
        checkAuth()
        return isAuthenticated.value ? { name: 'dashboard' } : { name: 'login' }
      },
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
    {
      path: '/events',
      name: 'events',
      component: EventsView,
      meta: { requiresAuth: true },
    },
    {
      path: '/resumes',
      name: 'resumes',
      component: ResumesView,
      meta: { requiresAuth: true },
    },
    {
      path: '/audit-logs',
      name: 'audit-logs',
      component: AuditLogsView,
      meta: { requiresAuth: true },
    },
    {
      path: '/points',
      name: 'points',
      component: PointsView,
      meta: { requiresAuth: true },
    },
    {
      path: '/credits',
      name: 'credits',
      component: CreditsView,
      meta: { requiresAuth: true },
    },
    {
      path: '/referrals',
      name: 'referrals',
      component: ReferralsView,
      meta: { requiresAuth: true },
    },
    {
      path: '/chat-activity',
      name: 'chat-activity',
      component: ChatActivityView,
      meta: { requiresAuth: true },
    },
    {
      path: '/chat-quests',
      name: 'chat-quests',
      component: ChatQuestsView,
      meta: { requiresAuth: true },
    },
    {
      path: '/raffles',
      name: 'raffles',
      component: RafflesView,
      meta: { requiresAuth: true },
    },
    {
      path: '/daily-tasks',
      name: 'daily-tasks',
      component: DailyTasksView,
      meta: { requiresAuth: true },
    },
    {
      path: '/challenges',
      name: 'challenges',
      component: ChallengesView,
      meta: { requiresAuth: true },
    },
    {
      path: '/minigames',
      name: 'minigames',
      component: CasinoView,
      meta: { requiresAuth: true },
    },
    {
      path: '/subscriptions',
      name: 'subscriptions',
      component: SubscriptionsView,
      meta: { requiresAuth: true },
    },
    {
      path: '/feedback',
      name: 'feedback',
      component: FeedbackView,
      meta: { requiresAuth: true },
    },
    {
      path: '/moderation',
      name: 'moderation',
      component: ModerationView,
      meta: { requiresAuth: true },
    },
  ],
})

router.beforeEach(async (to, _, next) => {
  checkAuth()

  if (to.meta.requiresAuth && !isAuthenticated.value) {
    next({ name: 'login', query: { redirect: to.fullPath } })
    return
  }

  // Любой Telegram-юзер мог получить tg_token и зайти на /dashboard, даже
  // не имея ни одного админ-пермишена. UI оставался пустым (API отвечал
  // 403 на каждый запрос), но дверь была открыта. Здесь явно проверяем,
  // что пермишены есть; иначе — выкидываем на login с тостом.
  if (to.meta.requiresAuth && isAuthenticated.value) {
    const allowed = await ensureAdminAccess()
    if (!allowed) {
      logout()
      next({ name: 'login' })
      return
    }
  }

  if (to.name === 'login' && isAuthenticated.value) {
    next({ name: 'dashboard' })
    return
  }

  next()
})

export default router
