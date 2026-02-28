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
