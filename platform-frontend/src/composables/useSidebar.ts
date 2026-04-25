import type { Component } from 'vue'
import { Award, Calendar, ClipboardList, Dices, Flame, Gift, Heart, HelpCircle, Home, Star, Trophy, User, Users } from 'lucide-vue-next'
import { ref } from 'vue'

export interface SidebarItem {
  title: string
  path: string
  icon: Component
  indicator?: boolean
  dataOnboarding?: string
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
        { title: 'FAQ', path: '/faq', icon: HelpCircle },
      ],
    },
    {
      label: 'Сообщество',
      items: [
        { title: 'События', path: '/events', icon: Calendar, indicator: true, dataOnboarding: 'events' },
        { title: 'Менторы', path: '/mentors', icon: Users },
        { title: 'Благодарности', path: '/kudos', icon: Heart },
      ],
    },
    {
      label: 'Активность',
      items: [
        { title: 'Мои баллы', path: '/points', icon: Star, dataOnboarding: 'points' },
        { title: 'Рейтинг', path: '/leaderboard', icon: Trophy },
        { title: 'Достижения', path: '/achievements', icon: Award },
        { title: 'Биржа заданий', path: '/tasks', icon: ClipboardList },
        { title: 'Задания в чатах', path: '/quests', icon: Flame },
        { title: 'Розыгрыши', path: '/raffles', icon: Gift },
        { title: 'Мини-игры', path: '/minigames', icon: Dices },
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
