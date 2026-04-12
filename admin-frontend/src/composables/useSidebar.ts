import type { Component } from 'vue'
import type { Permission } from '@/types/permissions'
import { computed, ref } from 'vue'
import Award from '~icons/lucide/award'
import BarChart3 from '~icons/lucide/bar-chart-3'
import Calendar from '~icons/lucide/calendar'
import ClipboardList from '~icons/lucide/clipboard-list'
import CreditCard from '~icons/lucide/credit-card'
import Dice5 from '~icons/lucide/dice-5'
import FileText from '~icons/lucide/file-text'
import Gift from '~icons/lucide/gift'
import Home from '~icons/lucide/home'
import Link from '~icons/lucide/link'
import MessageSquare from '~icons/lucide/message-square'
import Star from '~icons/lucide/star'
import Trophy from '~icons/lucide/trophy'
import User from '~icons/lucide/user'
import Users from '~icons/lucide/users'
import { usePermissions } from './usePermissions'

export interface SidebarItem {
  title: string
  path: string
  icon: Component
  requiredPermission?: Permission
}

export interface SidebarGroup {
  label?: string
  items: SidebarItem[]
}

const isCollapsed = ref(false)
const isMobileOpen = ref(false)

export function useSidebar() {
  const { hasPermission } = usePermissions()

  const allSidebarGroups: SidebarGroup[] = [
    {
      label: 'system',
      items: [
        {
          title: 'Дашборд',
          path: '/dashboard',
          icon: Home,
          requiredPermission: 'can_view_admin_panel',
        },
      ],
    },
    {
      label: 'content',
      items: [
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
          title: 'События',
          path: '/events',
          icon: Calendar,
          requiredPermission: 'can_view_admin_events',
        },
        {
          title: 'Резюме',
          path: '/resumes',
          icon: FileText,
          requiredPermission: 'can_view_admin_resumes',
        },
      ],
    },
    {
      label: 'reviews',
      items: [
        {
          title: 'Сообщество',
          path: '/reviews',
          icon: MessageSquare,
          requiredPermission: 'can_view_admin_reviews',
        },
        {
          title: 'Менторы',
          path: '/mentor-reviews',
          icon: MessageSquare,
          requiredPermission: 'can_view_admin_mentors_review',
        },
      ],
    },
    {
      label: 'gamification',
      items: [
        {
          title: 'Баллы',
          path: '/points',
          icon: Star,
          requiredPermission: 'can_view_admin_points',
        },
        {
          title: 'Активность',
          path: '/chat-activity',
          icon: BarChart3,
        },
        {
          title: 'Задания',
          path: '/chat-quests',
          icon: Award,
        },
        {
          title: 'Сезоны',
          path: '/seasons',
          icon: Trophy,
        },
        {
          title: 'Розыгрыши',
          path: '/raffles',
          icon: Gift,
        },
        {
          title: 'Мини-игры',
          path: '/minigames',
          icon: Dice5,
        },
      ],
    },
    {
      label: 'config',
      items: [
        {
          title: 'Подписки',
          path: '/subscriptions',
          icon: CreditCard,
          requiredPermission: 'can_view_admin_subscriptions',
        },
        {
          title: 'Рефералы',
          path: '/referrals',
          icon: Link,
        },
        {
          title: 'Журнал',
          path: '/audit-logs',
          icon: ClipboardList,
          requiredPermission: 'can_view_admin_audit_logs',
        },
      ],
    },
  ]

  const sidebarGroups = computed(() => {
    return allSidebarGroups
      .map(group => ({
        ...group,
        items: group.items.filter((item) => {
          if (!item.requiredPermission)
            return true
          return hasPermission.value(item.requiredPermission)
        }),
      }))
      .filter(group => group.items.length > 0)
  })

  const sidebarItems = computed(() => {
    return sidebarGroups.value.flatMap(g => g.items)
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
    sidebarGroups,
    toggleSidebar,
    toggleMobileSidebar,
    closeMobileSidebar,
  }
}
