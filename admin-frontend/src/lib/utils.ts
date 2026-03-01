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

/**
 * Парсит строку таймзоны вида "UTC", "UTC+3", "UTC-5" в смещение в минутах.
 */
export function parseTimezoneOffsetMinutes(tz: string): number {
  if (!tz || tz === 'UTC')
    return 0
  const match = tz.match(/^UTC([+-])(\d+)$/)
  if (!match)
    return 0
  const sign = match[1] === '+' ? 1 : -1
  return sign * Number(match[2]) * 60
}

/**
 * Конвертирует ISO строку (UTC) в значение для datetime-local в указанной таймзоне.
 */
export function toDatetimeLocal(isoString: string, timezone = 'UTC') {
  const date = new Date(isoString)
  const offsetMs = parseTimezoneOffsetMinutes(timezone) * 60 * 1000
  const local = new Date(date.getTime() + offsetMs)

  const year = local.getUTCFullYear()
  const month = String(local.getUTCMonth() + 1).padStart(2, '0')
  const day = String(local.getUTCDate()).padStart(2, '0')
  const hours = String(local.getUTCHours()).padStart(2, '0')
  const minutes = String(local.getUTCMinutes()).padStart(2, '0')

  return `${year}-${month}-${day}T${hours}:${minutes}`
}

/**
 * Конвертирует значение datetime-local в ISO строку (UTC), учитывая таймзону.
 * datetime-local показывает время в выбранной таймзоне, нужно вычесть смещение чтобы получить UTC.
 */
export function datetimeLocalToISO(datetimeLocal: string, timezone: string): string {
  const offsetMs = parseTimezoneOffsetMinutes(timezone) * 60 * 1000
  const asUtc = new Date(`${datetimeLocal}Z`)
  const utc = new Date(asUtc.getTime() - offsetMs)
  return utc.toISOString()
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
