import { QueryClient } from '@tanstack/vue-query'

// Create a client
export const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      // With this enabled, queries will refetch on window focus by default
      refetchOnWindowFocus: false,
      // With this enabled, queries will refetch on reconnect by default
      refetchOnReconnect: true,
      // With this enabled, queries will refetch on mount if the data is stale
      refetchOnMount: true,
      // Set a default stale time of 5 minutes
      staleTime: 5 * 60 * 1000,
      // Set a default garbage collection time of 10 minutes (formerly cacheTime)
      gcTime: 10 * 60 * 1000,
      // Don't retry failed queries by default
      retry: false,
    },
  },
})
