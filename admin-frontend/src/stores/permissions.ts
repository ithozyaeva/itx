import { defineStore } from 'pinia'
import { computed } from 'vue'
import { usePermissions } from '@/composables/usePermissions'

export const usePermissionsStore = defineStore('permissions', () => {
  // подключаем composable с vue-query
  const {
    userPermissions,
    isLoading,
    error,
    hasPermission,
    hasAnyPermission,
    clearPermissions,
    refetch,
  } = usePermissions()

  // computed для удобства
  const isReady = computed(() => !isLoading.value && !error.value)

  return {
    userPermissions,
    isLoading,
    error,
    isReady,
    hasPermission,
    hasAnyPermission,
    clearPermissions,
    refetch,
  }
})
