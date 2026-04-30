import type { TelegramUser } from '@/models/profile'

export type CommentEntityType = 'ai_material' | 'event'

export interface Comment {
  id: number
  entityType: CommentEntityType
  entityId: number
  authorId: number
  author?: TelegramUser
  body: string
  likesCount: number
  liked: boolean
  isHidden: boolean
  createdAt: string
  updatedAt: string
}

export interface CommentsResponse {
  items: Comment[]
  total: number
}

export interface ToggleCommentLikeResponse {
  liked: boolean
  likesCount: number
}

export const COMMENT_MAX_LEN = 4_000

// Маппинг entity_type → URL-сегмент. Бэкенд монтирует list/create на
// `/<entity-segment>/:id/comments`, поэтому фронту нужно знать соответствие.
export const COMMENT_ENTITY_SEGMENT: Record<CommentEntityType, string> = {
  ai_material: 'ai-materials',
  event: 'events',
}
