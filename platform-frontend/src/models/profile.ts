export type UserRole = 'UNSUBSCRIBER' | 'SUBSCRIBER' | 'MENTOR' | 'ADMIN' | 'EVENT_MAKER'

export type SubscriptionTierSlug = 'beginner' | 'foreman' | 'master' | 'king'

export interface SubscriptionTier {
  id: number
  slug: SubscriptionTierSlug | string
  name: string
  level: number
}

export interface TelegramUser {
  id: number
  telegramID: number
  tg: string
  birthday: string
  firstName: string
  lastName: string
  bio: string
  grade: string
  company: string
  avatarUrl: string
  roles: UserRole[]
  createdAt?: string
  subscriptionTier?: SubscriptionTier | null
}

export interface Mentor extends TelegramUser {
  id: number
  occupation: string
  experience: string
  profTags: ProfTag[]
  contacts: Contacts[]
  services: Service[]
}

export interface Service {
  id: number
  name: string
  price: number
}

export interface Contacts {
  id: number
  type: number
  link: string
}

export interface ProfTag {
  id: number
  title: string
}

export interface PublicProfile {
  member: TelegramUser
  points: number
  isMentor: boolean
  mentor?: Mentor
}

export interface NotificationSettings {
  id: number
  memberId: number
  muteAll: boolean
  newEvents: boolean
  remindWeek: boolean
  remindDay: boolean
  remindHour: boolean
  eventStart: boolean
  eventUpdates: boolean
  eventCancelled: boolean
}

export const SUBSCRIPTION_LEVELS = [
  'Новичок',
  'Бригадир',
  'Хозяин',
  'Меценат',
  'King',
  'Бизнесмен',
] as const

export type SubscriptionLevel = typeof SUBSCRIPTION_LEVELS[number]

// Маппинг тиров подписки (из subscription_tiers) на UI-уровни.
// beginner = "в комьюнити" без повышенного тира → Новичок.
// MENTOR как отдельная роль — поднимает до Хозяина, если тир ниже.
// ADMIN всегда побеждает и даёт Бизнесмена.
export function getSubscriptionLevel(
  roles: UserRole[],
  tier?: SubscriptionTier | null,
): SubscriptionLevel {
  if (roles.includes('ADMIN'))
    return 'Бизнесмен'
  const slug = tier?.slug
  if (slug === 'king')
    return 'King'
  if (slug === 'master')
    return 'Хозяин'
  if (roles.includes('MENTOR'))
    return 'Хозяин'
  if (slug === 'foreman')
    return 'Бригадир'
  return 'Новичок'
}

export function getSubscriptionLevelIndex(
  roles: UserRole[],
  tier?: SubscriptionTier | null,
): number {
  const level = getSubscriptionLevel(roles, tier)
  return SUBSCRIPTION_LEVELS.indexOf(level)
}

export const CONTACT_TYPES = [
  { value: 1, label: 'Telegram' },
  { value: 2, label: 'Email' },
  { value: 3, label: 'Телефон' },
  { value: 4, label: 'Другое' },
  { value: 5, label: 'LinkedIn' },
  { value: 6, label: 'GitHub' },
  { value: 7, label: 'VK' },
  { value: 8, label: 'Сайт' },
] as const
