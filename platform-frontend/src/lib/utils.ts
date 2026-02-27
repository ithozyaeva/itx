import type { ClassValue } from 'clsx'
import { clsx } from 'clsx'
import { twMerge } from 'tailwind-merge'

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs))
}

export function wrapLinks(text: string): string {
  // Регулярное выражение для поиска URL
  const urlRegex = /(https?:\/\/\S+)/g

  return text.replace(urlRegex, '<br /> <a href="$1" target="_blank" rel="noopener noreferrer" class="underline" >$1</a>')
}

export const dateFormatter = new Intl.DateTimeFormat('ru-RU', {
  day: 'numeric',
  month: 'long',
  year: 'numeric',
  hour: '2-digit',
  minute: '2-digit',
  hour12: false,
})
