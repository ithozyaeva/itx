import type { TelegramUser } from '@/models/profile'

export type AIMaterialContentType = 'prompt' | 'link' | 'agent'
export type AIMaterialKind = 'prompt' | 'skill' | 'library' | 'tutorial' | 'agent'
export type AIMaterialSort = 'new' | 'popular'

export interface AIMaterial {
  id: number
  authorId: number
  author?: TelegramUser
  title: string
  summary: string
  contentType: AIMaterialContentType
  materialKind: AIMaterialKind
  promptBody: string
  externalUrl: string
  agentConfig: string
  likesCount: number
  bookmarksCount: number
  commentsCount: number
  isHidden: boolean
  createdAt: string
  updatedAt: string
  tags: string[]
  liked: boolean
  bookmarked: boolean
}

export interface AIMaterialSearchResponse {
  items: AIMaterial[]
  total: number
}

export interface ToggleLikeResponse {
  liked: boolean
  likesCount: number
}

export interface ToggleBookmarkResponse {
  bookmarked: boolean
  bookmarksCount: number
}

export interface CreateAIMaterialRequest {
  title: string
  summary: string
  contentType: AIMaterialContentType
  materialKind: AIMaterialKind
  promptBody?: string
  externalUrl?: string
  agentConfig?: string
  tags: string[]
}

export interface AIMaterialFilters {
  kind?: AIMaterialKind | ''
  tag?: string
  q?: string
  sort?: AIMaterialSort
  mine?: boolean
  bookmarked?: boolean
  limit?: number
  offset?: number
}

export const AI_MATERIAL_KIND_LABELS: Record<AIMaterialKind, string> = {
  prompt: 'Промт',
  skill: 'Скилл',
  library: 'Библиотека',
  tutorial: 'Туториал',
  agent: 'Агент',
}

export const AI_MATERIAL_KIND_OPTIONS: { value: AIMaterialKind, label: string }[] = [
  { value: 'prompt', label: AI_MATERIAL_KIND_LABELS.prompt },
  { value: 'skill', label: AI_MATERIAL_KIND_LABELS.skill },
  { value: 'library', label: AI_MATERIAL_KIND_LABELS.library },
  { value: 'tutorial', label: AI_MATERIAL_KIND_LABELS.tutorial },
  { value: 'agent', label: AI_MATERIAL_KIND_LABELS.agent },
]

export const AI_MATERIAL_CONTENT_TYPE_LABELS: Record<AIMaterialContentType, string> = {
  prompt: 'Текст промта',
  link: 'Ссылка',
  agent: 'Конфиг агента',
}

export const AI_MATERIAL_CONTENT_TYPE_OPTIONS: { value: AIMaterialContentType, label: string, description: string }[] = [
  { value: 'prompt', label: 'Текст промта', description: 'Промт или инструкция, которую можно скопировать одной кнопкой' },
  { value: 'link', label: 'Ссылка', description: 'Внешний ресурс — Notion, GitHub, Telegraph и т.д.' },
  { value: 'agent', label: 'Конфиг агента', description: 'JSON/YAML конфиг готового AI-агента' },
]

export const AI_MATERIAL_LIMITS = {
  titleMin: 3,
  titleMax: 120,
  summaryMin: 30,
  summaryMax: 500,
  promptMax: 50_000,
  agentMax: 50_000,
  urlMax: 2048,
  tagsMax: 5,
  tagLenMax: 40,
} as const
