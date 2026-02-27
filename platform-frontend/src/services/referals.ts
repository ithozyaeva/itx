import type { ReferalLink } from '@/models/referals'
import { apiClient } from './api'

export interface ReferalSearchFilters {
  grade?: string
  company?: string
  status?: string
}

function cleanFilters(filters?: ReferalSearchFilters): Record<string, string> {
  if (!filters)
    return {}
  return Object.fromEntries(
    Object.entries(filters).filter(([_, v]) => v !== undefined && v !== ''),
  )
}

export const referalLinkService = {
  search: async (limit: number, offset: number, filters?: ReferalSearchFilters) => {
    return apiClient.get('referals', { searchParams: { limit, offset, ...cleanFilters(filters) } }).json<{ items: ReferalLink[], total: number }>()
  },
  addLink: async (link: Partial<ReferalLink>) => {
    return apiClient.post('referals/add-link', { json: link }).json<ReferalLink>()
  },
  updateLink: async (link: Partial<ReferalLink>) => {
    return apiClient.put('referals/update-link', { json: link }).json<ReferalLink>()
  },

  deleteLink: async (linkId: number) => {
    return apiClient.delete(`referals/delete-link`, { json: { id: linkId } }).json<ReferalLink>()
  },
}
