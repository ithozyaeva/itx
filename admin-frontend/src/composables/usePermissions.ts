import type { Permission } from '@/types/permissions'
import { useQuery, useQueryClient } from '@tanstack/vue-query'

import { computed } from 'vue'
import api from '@/lib/api'

export function usePermissions() {
  const queryClient = useQueryClient()

  const {
    data: userPermissions,
    isLoading,
    error,
    refetch,
  } = useQuery({
    queryKey: ['permissions'],
    queryFn: async () => {
      return await api.get('me/permissions').json<Permission[]>()
    },
    staleTime: Infinity,
    gcTime: Infinity,
    retry: false,
    refetchOnWindowFocus: false,
  })

  const hasPermission = computed(() => {
    return (permission: Permission) => {
      return !!userPermissions.value?.includes(permission)
    }
  })
  const hasAnyPermission = computed(() => (permissions: Permission[]): boolean => {
    return permissions.some(permission => hasPermission.value(permission))
  })

  const clearPermissions = () => {
    queryClient.invalidateQueries({ queryKey: ['permissions'] })
  }

  return {
    userPermissions,
    isLoading,
    error,
    refetch,
    hasPermission,
    hasAnyPermission,
    clearPermissions,
  }
}
