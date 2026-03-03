import type { Component } from 'vue'
import { Award, BookOpen, Calendar, ClipboardList, FileText, Folder, Home, MessageSquare, ShoppingBag, Star, Trophy, Users } from 'lucide-vue-next'
import { ref } from 'vue'

export interface SidebarItem {
  title: string
  path: string
  icon: Component
  indicator?: boolean
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
        { title: 'События', path: '/events', icon: Calendar, indicator: true },
        { title: 'Контент', path: '/content', icon: BookOpen },
        { title: 'Участники', path: '/mentors', icon: Users },
      ],
    },
    {
      label: 'Активность',
      items: [
        { title: 'Мои баллы', path: '/points', icon: Star },
        { title: 'Рейтинг', path: '/leaderboard', icon: Trophy },
        { title: 'Достижения', path: '/achievements', icon: Award },
        { title: 'Биржа заданий', path: '/tasks', icon: ClipboardList },
      ],
    },
    {
      label: 'Мои разделы',
      items: [
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
