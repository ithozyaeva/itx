import type { Comment, CommentEntityType, CommentsResponse, ToggleCommentLikeResponse } from '@/models/comment'
import { COMMENT_ENTITY_SEGMENT } from '@/models/comment'
import { apiClient } from './api'

// Универсальный сервис для работы с комментариями всех сущностей.
// Бэкенд использует один CommentService с visibility-checker'ами per-entity,
// фронт — один shared компонент Comments.vue, дёргающий эти методы.
export const commentsService = {
  async list(entityType: CommentEntityType, entityId: number, limit = 20, offset = 0) {
    const params = new URLSearchParams()
    params.set('limit', limit.toString())
    params.set('offset', offset.toString())
    return apiClient
      .get(`${COMMENT_ENTITY_SEGMENT[entityType]}/${entityId}/comments`, { searchParams: params })
      .json<CommentsResponse>()
  },

  async create(entityType: CommentEntityType, entityId: number, body: string) {
    return apiClient
      .post(`${COMMENT_ENTITY_SEGMENT[entityType]}/${entityId}/comments`, { json: { body } })
      .json<Comment>()
  },

  async update(commentId: number, body: string) {
    return apiClient.patch(`comments/${commentId}`, { json: { body } }).json<Comment>()
  },

  async remove(commentId: number) {
    return apiClient.delete(`comments/${commentId}`)
  },

  async toggleLike(commentId: number) {
    return apiClient.post(`comments/${commentId}/like`).json<ToggleCommentLikeResponse>()
  },
}
