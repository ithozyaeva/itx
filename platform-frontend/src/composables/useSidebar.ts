import type { Component } from 'vue'
import type { SubscriptionTierSlug } from '@/models/profile'
import { Calendar, ClipboardList, Crown, Dices, Gift, HelpCircle, Home, Sparkles, Sprout, User, Users } from 'lucide-vue-next'
import { ref } from 'vue'

export interface SidebarItem {
  title: string
  path: string
  icon: Component
  indicator?: boolean
  dataOnboarding?: string
  // requiresSubscription — пункт скрывается у UNSUBSCRIBER. Совпадает с
  // meta.requiresSubscription в роутере, чтобы UI и guard не расходились.
  requiresSubscription?: boolean
  // requiresMinTier — пункт скрывается у тех, чей tier ниже указанного.
  // Совпадает с meta.requiresMinTier в роутере (например, 'master' для
  // премиум-разделов).
  requiresMinTier?: SubscriptionTierSlug
  // visibleFor — 'unsubscribed' значит виден только без подписки (например,
  // пункт «Тарифы»). Без флага — виден всем авторизованным.
  visibleFor?: 'unsubscribed'
}

export interface SidebarGroup {
  label?: string
  items: SidebarItem[]
}

// Создаем синглтон с состоянием
const state = {
  isOpen: ref(false),
  sidebarGroups: ref<SidebarGroup[]>([
    {
      items: [
        { title: 'Главная', path: '/', icon: Home },
        { title: 'Мой профиль', path: '/me', icon: User, dataOnboarding: 'profile' },
        { title: 'Тарифы', path: '/tariffs', icon: Crown, visibleFor: 'unsubscribed' },
        { title: 'FAQ', path: '/faq', icon: HelpCircle },
      ],
    },
    {
      label: 'Сообщество',
      items: [
        { title: 'События', path: '/events', icon: Calendar, indicator: true, dataOnboarding: 'events', requiresSubscription: true },
        { title: 'Менторы', path: '/mentors', icon: Users },
      ],
    },
    {
      label: 'Знания',
      items: [
        { title: 'AI-материалы', path: '/ai-materials', icon: Sparkles, requiresSubscription: true },
      ],
    },
    {
      label: 'Активность',
      items: [
        { title: 'Прогресс', path: '/progress', icon: Sprout, dataOnboarding: 'points', requiresSubscription: true },
        { title: 'Биржа заданий', path: '/tasks', icon: ClipboardList, requiresSubscription: true },
      ],
    },
    {
      label: 'Бонусы',
      items: [
        { title: 'Розыгрыши', path: '/raffles', icon: Gift, requiresSubscription: true },
        { title: 'Мини-игры', path: '/minigames', icon: Dices, requiresSubscription: true },
      ],
    },
  ]),
}

// Функция для управления состоянием
function toggleSidebar() {
  state.isOpen.value = !state.isOpen.value
}

// Экспортируем composable
export function useSidebar() {
  return {
    isOpen: state.isOpen,
    sidebarGroups: state.sidebarGroups,
    toggleSidebar,
  }
}
