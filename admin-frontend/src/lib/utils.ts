import type { ClassValue } from 'clsx'
import { clsx } from 'clsx'
import { twMerge } from 'tailwind-merge'

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs))
}

export function formatDateToInput(date?: string) {
  if (date) {
    return date.slice(0, 10)
  }
  return ''
}
export type DeepPartial<T> = T extends object ? {
  [P in keyof T]?: DeepPartial<T[P]>;
} : T

export function toDatetimeLocal(isoString: string) {
  const date = new Date(isoString)

  const year = date.getFullYear()
  const month = String(date.getMonth() + 1).padStart(2, '0')
  const day = String(date.getDate()).padStart(2, '0')

  const hours = String(date.getHours()).padStart(2, '0')
  const minutes = String(date.getMinutes()).padStart(2, '0')

  // Собираем строку для input[type="datetime-local"]
  return `${year}-${month}-${day}T${hours}:${minutes}`
}

export async function downloadFile(url: string, filename: string) {
  const tgToken = localStorage.getItem('tg_token')
  const response = await fetch(`/api/admin/${url}`, {
    headers: tgToken ? { 'X-Telegram-User-Token': tgToken } : {},
  })
  const blob = await response.blob()
  const link = document.createElement('a')
  link.href = URL.createObjectURL(blob)
  link.download = filename
  link.click()
  URL.revokeObjectURL(link.href)
}

export function cleanParams(params: Record<string, any>) {
  return Object.fromEntries(
    Object.entries(params).filter(([_, v]) => {
      if (Array.isArray(v)) {
        return v.length > 0
      }
      return v !== undefined && v !== null && v !== ''
    }),
  )
}
