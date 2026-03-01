import type { ReferalLink } from '@/models/referals'
import { apiClient } from './api'

export interface ReferalSearchFilters {
  grade?: string
  company?: string
  status?: string
  profTagIds?: number[]
}

function cleanFilters(filters?: ReferalSearchFilters): Record<string, string> {
  if (!filters)
    return {}
  const result: Record<string, string> = {}
  if (filters.grade)
    result.grade = filters.grade
  if (filters.company)
    result.company = filters.company
  if (filters.status)
    result.status = filters.status
  if (filters.profTagIds && filters.profTagIds.length > 0)
    result.profTagIds = filters.profTagIds.join(',')
  return result
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

  trackConversion: async (referralLinkId: number) => {
    await apiClient.post('referals/track-conversion', { json: { referralLinkId } })
  },
}
