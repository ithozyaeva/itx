import type { Component } from 'vue'
import { Calendar, FileText, Folder, MessageSquare, Trophy, User, Users } from 'lucide-vue-next'
import { ref } from 'vue'

export interface SidebarItem {
  title: string
  path: string
  icon: Component
}

// Создаем синглтон с состоянием
const state = {
  isOpen: ref(false),
  sidebarItems: ref<SidebarItem[]>([
    {
      title: 'Профиль',
      path: '/me',
      icon: User,
    },
    {
      title: 'Календарь событий',
      path: '/events',
      icon: Calendar,
    },
    {
      title: 'Менторы',
      path: '/mentors',
      icon: Users,
    },
    {
      title: 'Таблица рефералов',
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
      title: 'Рейтинг',
      path: '/leaderboard',
      icon: Trophy,
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
