import type { ClassValue } from 'clsx'
import { clsx } from 'clsx'
import { twMerge } from 'tailwind-merge'

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs))
}

const RE_AMP = /&/g
const RE_LT = /</g
const RE_GT = />/g
const RE_QUOT = /"/g
const RE_APOS = /'/g

function escapeHtml(str: string): string {
  return str
    .replace(RE_AMP, '&amp;')
    .replace(RE_LT, '&lt;')
    .replace(RE_GT, '&gt;')
    .replace(RE_QUOT, '&quot;')
    .replace(RE_APOS, '&#039;')
}

const RE_URL = /(https?:\/\/\S+)/g

// data-external-link — маркер для глобального делегата клика в App.vue:
// внутри Telegram Mini App клик уходит в openLink() → openTelegramLink/
// openLink Telegram-клиента, и юзер не вываливается во внешний браузер.
// В обычном браузере делегат не вмешивается — отрабатывает обычный
// target="_blank".
export function wrapLinks(text: string): string {
  const escaped = escapeHtml(text)

  return escaped.replace(RE_URL, '<br /> <a href="$1" target="_blank" rel="noopener noreferrer" class="underline" data-external-link >$1</a>')
}

export const dateFormatter = new Intl.DateTimeFormat('ru-RU', {
  day: 'numeric',
  month: 'long',
  year: 'numeric',
  hour: '2-digit',
  minute: '2-digit',
  hour12: false,
})

export const shortDateFormatter = new Intl.DateTimeFormat('ru-RU', {
  day: 'numeric',
  month: 'long',
  year: 'numeric',
})

export function displayName(firstName?: string | null, lastName?: string | null): string {
  return [firstName, lastName].filter(Boolean).join(' ') || '—'
}

export function formatShortDate(date: string | Date): string {
  const d = typeof date === 'string' ? new Date(date) : date
  if (Number.isNaN(d.getTime()))
    return ''
  return shortDateFormatter.format(d)
}
