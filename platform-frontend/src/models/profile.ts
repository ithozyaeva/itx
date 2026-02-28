export type UserRole = 'UNSUBSCRIBER' | 'SUBSCRIBER' | 'MENTOR' | 'ADMIN' | 'EVENT_MAKER'

export interface TelegramUser {
  id: number
  telegramID: number
  tg: string
  birthday: string
  firstName: string
  lastName: string
  bio: string
  avatarUrl: string
  roles: UserRole[]
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

export const SUBSCRIPTION_LEVELS = [
  'Новичок',
  'Бригадир',
  'Хозяин',
  'Меценат',
  'King',
  'Бизнесмен',
] as const

export type SubscriptionLevel = typeof SUBSCRIPTION_LEVELS[number]

export function getSubscriptionLevel(roles: UserRole[]): SubscriptionLevel {
  if (roles.includes('ADMIN'))
    return 'Бизнесмен'
  if (roles.includes('MENTOR'))
    return 'Хозяин'
  if (roles.includes('SUBSCRIBER'))
    return 'Бригадир'
  return 'Новичок'
}

export function getSubscriptionLevelIndex(roles: UserRole[]): number {
  const level = getSubscriptionLevel(roles)
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
