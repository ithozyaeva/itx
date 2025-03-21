import type { Component } from 'vue'
import { ref } from 'vue'
import Home from '~icons/lucide/home'
import MessageSquare from '~icons/lucide/message-square'
import User from '~icons/lucide/user'
import Users from '~icons/lucide/users'

export interface SidebarItem {
  title: string
  path: string
  icon: Component
}

export function useSidebar() {
  const isCollapsed = ref(false)

  const sidebarItems = ref<SidebarItem[]>([
    {
      title: 'Дашборд',
      path: '/dashboard',
      icon: Home,
    },
    {
      title: 'Менторы',
      path: '/mentors',
      icon: Users,
    },
    {
      title: 'Участники',
      path: '/members',
      icon: User,
    },
    {
      title: 'Отзывы на сообщество',
      path: '/reviews',
      icon: MessageSquare,
    },
    {
      title: 'Отзывы на менторов',
      path: '/mentor-reviews',
      icon: MessageSquare,
    },
  ])

  const toggleSidebar = () => {
    isCollapsed.value = !isCollapsed.value
  }

  return {
    isCollapsed,
    sidebarItems,
    toggleSidebar,
  }
}
