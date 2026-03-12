import type { Component } from 'vue'
import { Award, BarChart3, BookOpen, Calendar, ClipboardList, Dices, FileText, Flame, Folder, Gift, Heart, Home, MessageSquare, Shield, ShoppingBag, Star, Trophy, Users } from 'lucide-vue-next'
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
      ],
    },
    {
      label: 'Сообщество',
      items: [
        { title: 'События', path: '/events', icon: Calendar, indicator: true, dataOnboarding: 'events' },
        { title: 'Контент', path: '/content', icon: BookOpen },
        { title: 'Менторы', path: '/mentors', icon: Users },
        { title: 'Благодарности', path: '/kudos', icon: Heart },
        { title: 'Гильдии', path: '/guilds', icon: Shield },
      ],
    },
    {
      label: 'Активность',
      items: [
        { title: 'Мои баллы', path: '/points', icon: Star, dataOnboarding: 'points' },
        { title: 'Рейтинг', path: '/leaderboard', icon: Trophy },
        { title: 'Сезоны', path: '/seasons', icon: Calendar },
        { title: 'Достижения', path: '/achievements', icon: Award },
        { title: 'Биржа заданий', path: '/tasks', icon: ClipboardList },
        { title: 'Задания в чатах', path: '/quests', icon: Flame },
        { title: 'Розыгрыши', path: '/raffles', icon: Gift },
        { title: 'Казино', path: '/casino', icon: Dices },
      ],
    },
    {
      label: 'Мои разделы',
      items: [
        { title: 'Моя статистика', path: '/my-stats', icon: BarChart3 },
        { title: 'Рефералки', path: '/referals', icon: Folder },
        { title: 'Резюме', path: '/resumes', icon: FileText },
        { title: 'Мои отзывы', path: '/my-reviews', icon: MessageSquare },
        { title: 'Барахолка', path: '/marketplace', icon: ShoppingBag },
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
