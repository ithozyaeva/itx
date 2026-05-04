import { formatShortDate } from './utils'

export function progressPct(progress: number, target: number): number {
  if (!target || target <= 0)
    return 0
  return Math.min(100, Math.round((progress / target) * 100))
}

export function daysUntil(dateStr: string): number {
  const diffMs = new Date(dateStr).getTime() - Date.now()
  return Math.ceil(diffMs / 86400000)
}

export function formatDeadline(dateStr: string): string {
  const days = daysUntil(dateStr)
  if (days <= 0)
    return 'Истекает'
  if (days === 1)
    return '1 день'
  if (days <= 7)
    return `${days} дн.`
  return formatShortDate(dateStr)
}

export type DeadlineUrgency = 'expired' | 'critical' | 'warning' | 'normal'

export function deadlineUrgency(dateStr: string): DeadlineUrgency {
  const days = daysUntil(dateStr)
  if (days <= 0)
    return 'expired'
  if (days <= 1)
    return 'critical'
  if (days <= 3)
    return 'warning'
  return 'normal'
}
