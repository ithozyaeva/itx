import type { Component } from 'vue'
import type { Permission } from '@/types/permissions'
import { computed, ref } from 'vue'
import ClipboardList from '~icons/lucide/clipboard-list'
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

// Global state so Header and Sidebar share the same refs
const isCollapsed = ref(false)
const isMobileOpen = ref(false)

export function useSidebar() {
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
    {
      title: 'Журнал действий',
      path: '/audit-logs',
      icon: ClipboardList,
      requiredPermission: 'can_view_admin_audit_logs',
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

  const toggleMobileSidebar = () => {
    isMobileOpen.value = !isMobileOpen.value
  }

  const closeMobileSidebar = () => {
    isMobileOpen.value = false
  }

  return {
    isCollapsed,
    isMobileOpen,
    sidebarItems,
    toggleSidebar,
    toggleMobileSidebar,
    closeMobileSidebar,
  }
}
