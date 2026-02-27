import type { Component } from 'vue'
import type { Permission } from '@/types/permissions'
import { computed, ref } from 'vue'
import FileText from '~icons/lucide/file-text'
import Home from '~icons/lucide/home'
import MessageSquare from '~icons/lucide/message-square'
import User from '~icons/lucide/user'
import Users from '~icons/lucide/users'
import { usePermissions } from './usePermissions'

export interface SidebarItem {
  title: string
  path: string
  icon: Component
  requiredPermission?: Permission
}

export function useSidebar() {
  const isCollapsed = ref(false)
  const { hasPermission } = usePermissions()

  const allSidebarItems = ref<SidebarItem[]>([
    {
      title: 'Дашборд',
      path: '/dashboard',
      icon: Home,
      requiredPermission: 'can_view_admin_panel',
    },
    {
      title: 'Менторы',
      path: '/mentors',
      icon: Users,
      requiredPermission: 'can_view_admin_mentors',
    },
    {
      title: 'Участники',
      path: '/members',
      icon: User,
      requiredPermission: 'can_view_admin_members',
    },
    {
      title: 'Отзывы на сообщество',
      path: '/reviews',
      icon: MessageSquare,
      requiredPermission: 'can_view_admin_reviews',
    },
    {
      title: 'Отзывы на менторов',
      path: '/mentor-reviews',
      icon: MessageSquare,
      requiredPermission: 'can_view_admin_mentors_review',
    },
    {
      title: 'События',
      path: '/events',
      icon: Users,
      requiredPermission: 'can_view_admin_events',
    },
    {
      title: 'Резюме',
      path: '/resumes',
      icon: FileText,
      requiredPermission: 'can_view_admin_resumes',
    },
  ])

  const sidebarItems = computed(() => {
    return allSidebarItems.value.filter((item) => {
      if (!item.requiredPermission) {
        return true
      }
      return hasPermission.value(item.requiredPermission)
    })
  })

  const toggleSidebar = () => {
    isCollapsed.value = !isCollapsed.value
  }

  return {
    isCollapsed,
    sidebarItems,
    toggleSidebar,
  }
}
