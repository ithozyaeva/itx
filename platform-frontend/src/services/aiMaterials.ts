import type {
  AIMaterial,
  AIMaterialComment,
  AIMaterialCommentsResponse,
  AIMaterialFilters,
  AIMaterialSearchResponse,
  CreateAIMaterialRequest,
  ToggleBookmarkResponse,
  ToggleLikeResponse,
} from '@/models/aiMaterial'
import { apiClient } from './api'

function buildSearchParams(filters: AIMaterialFilters): URLSearchParams {
  const params = new URLSearchParams()
  if (filters.kind)
    params.set('kind', filters.kind)
  if (filters.tag)
    params.set('tag', filters.tag)
  if (filters.q)
    params.set('q', filters.q)
  if (filters.sort)
    params.set('sort', filters.sort)
  if (filters.mine)
    params.set('mine', 'true')
  if (filters.bookmarked)
    params.set('bookmarked', 'true')
  if (filters.limit != null)
    params.set('limit', filters.limit.toString())
  if (filters.offset != null)
    params.set('offset', filters.offset.toString())
  return params
}

export const aiMaterialsService = {
  async search(filters: AIMaterialFilters = {}) {
    return apiClient
      .get('ai-materials', { searchParams: buildSearchParams(filters) })
      .json<AIMaterialSearchResponse>()
  },

  async getById(id: number) {
    return apiClient.get(`ai-materials/${id}`).json<AIMaterial>()
  },

  async create(data: CreateAIMaterialRequest) {
    return apiClient.post('ai-materials', { json: data }).json<AIMaterial>()
  },

  async update(id: number, data: CreateAIMaterialRequest) {
    return apiClient.put(`ai-materials/${id}`, { json: data }).json<AIMaterial>()
  },

  async remove(id: number) {
    return apiClient.delete(`ai-materials/${id}`)
  },

  async setHidden(id: number, hidden: boolean) {
    return apiClient.post(`ai-materials/${id}/hidden`, { json: { hidden } })
  },

  async topTags(q?: string, limit = 20) {
    const params = new URLSearchParams()
    if (q)
      params.set('q', q)
    params.set('limit', limit.toString())
    return apiClient.get('ai-materials/tags', { searchParams: params }).json<{ tags: string[] }>()
  },

  async toggleLike(id: number) {
    return apiClient.post(`ai-materials/${id}/like`).json<ToggleLikeResponse>()
  },

  async toggleBookmark(id: number) {
    return apiClient.post(`ai-materials/${id}/bookmark`).json<ToggleBookmarkResponse>()
  },

  async listComments(id: number, limit = 20, offset = 0) {
    const params = new URLSearchParams()
    params.set('limit', limit.toString())
    params.set('offset', offset.toString())
    return apiClient.get(`ai-materials/${id}/comments`, { searchParams: params }).json<AIMaterialCommentsResponse>()
  },

  async createComment(id: number, body: string) {
    return apiClient.post(`ai-materials/${id}/comments`, { json: { body } }).json<AIMaterialComment>()
  },

  async updateComment(commentId: number, body: string) {
    return apiClient.patch(`ai-material-comments/${commentId}`, { json: { body } }).json<AIMaterialComment>()
  },

  async deleteComment(commentId: number) {
    return apiClient.delete(`ai-material-comments/${commentId}`)
  },

  async toggleCommentLike(commentId: number) {
    return apiClient.post(`ai-material-comments/${commentId}/like`).json<ToggleLikeResponse>()
  },
}
