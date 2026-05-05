import type { RouteRecordRaw } from 'vue-router'
import type { SubscriptionTierSlug } from '@/models/profile'
import { createRouter, createWebHistory } from 'vue-router'
import { hasMinTier, isUserSubscribed, useUserLevel } from '@/composables/useUser'
// Главные точки входа — Dashboard и /me — оставлены статически,
// чтобы первый paint не ждал лишний chunk. Остальные страницы lazy.
import Dashboard from '@/pages/Dashboard.vue'
import Home from '@/pages/User.vue'

declare module 'vue-router' {
  // requiresMinTier — гейт по минимальному уровню тира. Используется
  // для премиум-разделов (например, AI-материалы — только master+).
  // Имя slug совпадает с TIER_SLUG_LEVELS в useUser.
  interface RouteMeta {
    requiresMinTier?: SubscriptionTierSlug
  }
}

// requiresSubscription: true — UNSUBSCRIBER редиректится на главную.
// Главная сама показывает teaser «Открой полный доступ» с кнопкой на /tariffs,
// в отличие от лобового приземления на витрину тарифов из любого
// гейтнутого роута.
// Открытые для UNSUBSCRIBER (преимущественно «прогрев»): /, /me, /mentors,
// /mentors/:id, /members/:id, /faq, /tariffs.
const routes: RouteRecordRaw[] = [
  { path: '/', component: Dashboard, name: 'dashboard' },
  { path: '/me', component: Home, name: 'profile', meta: { breadcrumb: [{ label: 'Мой профиль' }] } },
  { path: '/tariffs', name: 'tariffs', component: () => import('@/pages/Tariffs.vue'), meta: { breadcrumb: [{ label: 'Тарифы' }] } },
  { path: '/events', component: () => import('@/pages/Events.vue'), name: 'events', meta: { breadcrumb: [{ label: 'События' }], requiresSubscription: true } },
  { path: '/content', redirect: '/events?tab=content' },
  { path: '/members/:id', component: () => import('@/pages/MemberProfile.vue'), name: 'memberProfile', meta: { breadcrumb: [{ label: 'Рейтинг', to: '/progress?tab=leaderboard' }, { label: 'Профиль участника' }] } },
  { path: '/mentors', component: () => import('@/pages/Mentors.vue'), name: 'mentors', meta: { breadcrumb: [{ label: 'Менторы' }] } },
  { path: '/mentors/:id', component: () => import('@/pages/MentorProfile.vue'), name: 'mentorProfile', meta: { breadcrumb: [{ label: 'Менторы', to: '/mentors' }, { label: 'Профиль ментора' }] } },
  { path: '/referals', component: () => import('@/pages/ReferalLinks.vue'), name: 'referals', meta: { breadcrumb: [{ label: 'Рефералы' }], requiresSubscription: true } },
  { path: '/resumes', component: () => import('@/pages/Resumes.vue'), name: 'resumes', meta: { breadcrumb: [{ label: 'Резюме' }], requiresSubscription: true } },
  { path: '/my-reviews', component: () => import('@/pages/MyReviews.vue'), name: 'myReviews', meta: { breadcrumb: [{ label: 'Мои отзывы' }] } },
  { path: '/progress', component: () => import('@/pages/Progress.vue'), name: 'progress', meta: { breadcrumb: [{ label: 'Прогресс' }], requiresSubscription: true } },
  // Старые URL дейликов/квестов/челленджей/баллов сохраняем как редиректы,
  // чтобы не ломать ссылки из бота, e-mail-а и закладок пользователей.
  { path: '/points', redirect: '/progress?tab=history' },
  { path: '/dailies', redirect: '/progress?tab=today' },
  { path: '/quests', redirect: '/progress?tab=period&kind=chats' },
  { path: '/challenges', redirect: '/progress?tab=period' },
  { path: '/leaderboard', redirect: '/progress?tab=leaderboard' },
  { path: '/achievements', redirect: '/progress?tab=achievements' },
  { path: '/marketplace', component: () => import('@/pages/Marketplace.vue'), name: 'marketplace', meta: { breadcrumb: [{ label: 'Барахолка' }], requiresSubscription: true } },
  { path: '/ai-materials', component: () => import('@/pages/AIMaterials.vue'), name: 'aiMaterials', meta: { breadcrumb: [{ label: 'AI-материалы' }], requiresSubscription: true } },
  { path: '/ai-materials/:id', component: () => import('@/pages/AIMaterialDetail.vue'), name: 'aiMaterialDetail', meta: { breadcrumb: [{ label: 'AI-материалы', to: '/ai-materials' }, { label: 'Материал' }], requiresSubscription: true } },
  { path: '/tasks', component: () => import('@/pages/TaskExchange.vue'), name: 'taskExchange', meta: { breadcrumb: [{ label: 'Биржа заданий' }], requiresSubscription: true } },
  { path: '/auto-apply', component: () => import('@/pages/AutoApplyBot.vue'), name: 'autoApplyBot', meta: { breadcrumb: [{ label: 'Автоотклики' }], requiresSubscription: true } },
  { path: '/kudos', redirect: '/progress?tab=kudos' },
  { path: '/raffles', component: () => import('@/pages/Raffles.vue'), name: 'raffles', meta: { breadcrumb: [{ label: 'Розыгрыши' }], requiresSubscription: true } },
  { path: '/minigames', component: () => import('@/pages/Casino.vue'), name: 'minigames', meta: { breadcrumb: [{ label: 'Мини-игры' }], requiresSubscription: true } },
  { path: '/my-stats', redirect: '/progress?tab=stats' },
  { path: '/notifications', redirect: '/me' },
  { path: '/faq', name: 'faq', component: () => import('@/pages/FAQ.vue'), meta: { breadcrumb: [{ label: 'FAQ' }] } },
  { path: '/:pathMatch(.*)*', redirect: '/' },
]

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes,
})

router.beforeEach((to) => {
  const subscribed = isUserSubscribed().value
  // UNSUBSCRIBER на гейтнутом роуте → дашборд. Дашборд сам рендерит
  // subscription-teaser, оттуда юзер сам идёт на /tariffs.
  if (to.meta.requiresSubscription && !subscribed) {
    return { name: 'dashboard' }
  }
  // requiresMinTier — гейт по уровню подписки выше базового (например,
  // master+ для AI-материалов). Не-достаточный тир кидаем на дашборд
  // (там подписчик увидит обычный набор; мотивации к апгрейду пока нет
  // отдельного экрана — это решим, когда раздел станет основным).
  if (to.meta.requiresMinTier && !hasMinTier(to.meta.requiresMinTier).value) {
    return { name: 'dashboard' }
  }
  // /tariffs — витрина для UNSUBSCRIBER. Кому показывать нечего:
  //   - подписчики любого тира (есть subscriptionTier);
  //   - ADMIN/MENTOR — у них levelIndex > 0 даже без оплаченного тира,
  //     иначе админ/ментор приземляется на витрину тарифов из ссылок
  //     бота / прямого URL и видит «купи подписку», что нелепо.
  const { levelIndex } = useUserLevel()
  if (to.name === 'tariffs' && (subscribed || levelIndex.value > 0)) {
    return { name: 'dashboard' }
  }
  return true
})

export default router
