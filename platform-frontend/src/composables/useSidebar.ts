import type { Component } from 'vue'
import { Award, BookOpen, Calendar, FileText, Folder, Home, MessageSquare, Star, Trophy, Users } from 'lucide-vue-next'
import { ref } from 'vue'

export interface SidebarItem {
  title: string
  path: string
  icon: Component
  indicator?: boolean
}

// Создаем синглтон с состоянием
const state = {
  isOpen: ref(false),
  sidebarItems: ref<SidebarItem[]>([
    {
      title: 'Дом',
      path: '/',
      icon: Home,
    },
    {
      title: 'События',
      path: '/events',
      icon: Calendar,
      indicator: true,
    },
    {
      title: 'Контент',
      path: '/content',
      icon: BookOpen,
    },
    {
      title: 'Участники',
      path: '/mentors',
      icon: Users,
    },
    {
      title: 'Рефералки',
      path: '/referals',
      icon: Folder,
    },
    {
      title: 'Резюме',
      path: '/resumes',
      icon: FileText,
    },
    {
      title: 'Мои отзывы',
      path: '/my-reviews',
      icon: MessageSquare,
    },
    {
      title: 'Мои баллы',
      path: '/points',
      icon: Star,
    },
    {
      title: 'Рейтинг',
      path: '/leaderboard',
      icon: Trophy,
    },
    {
      title: 'Достижения',
      path: '/achievements',
      icon: Award,
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
    sidebarItems: state.sidebarItems,
    toggleSidebar,
  }
}
