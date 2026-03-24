import type { CommunityEvent } from '@/models/event'

function advanceDate(date: Date, period: string, interval: number): Date {
  const next = new Date(date)
  switch (period) {
    case 'DAILY':
      next.setDate(next.getDate() + interval)
      break
    case 'WEEKLY':
      next.setDate(next.getDate() + interval * 7)
      break
    case 'MONTHLY':
      next.setMonth(next.getMonth() + interval)
      break
    case 'YEARLY':
      next.setFullYear(next.getFullYear() + interval)
      break
  }
  return next
}

/**
 * Возвращает дату следующего (ближайшего будущего) вхождения события.
 * Для неповторяющихся событий — исходная дата.
 */
export function getNextOccurrenceDate(event: CommunityEvent, now: Date = new Date()): Date {
  if (!event.isRepeating || !event.repeatPeriod) {
    return new Date(event.date)
  }

  const interval = event.repeatInterval ?? 1
  const start = new Date(event.date)

  if (start >= now) {
    return start
  }

  switch (event.repeatPeriod) {
    case 'DAILY': {
      const daysSince = Math.floor((now.getTime() - start.getTime()) / (24 * 60 * 60 * 1000))
      const nextDays = (Math.floor(daysSince / interval) + 1) * interval
      const next = new Date(start)
      next.setDate(next.getDate() + nextDays)
      return next
    }
    case 'WEEKLY': {
      const weeksSince = Math.floor((now.getTime() - start.getTime()) / (7 * 24 * 60 * 60 * 1000))
      const nextWeeks = (Math.floor(weeksSince / interval) + 1) * interval
      const next = new Date(start)
      next.setDate(next.getDate() + nextWeeks * 7)
      return next
    }
    case 'MONTHLY': {
      let current = new Date(start)
      while (current < now) {
        current = advanceDate(current, 'MONTHLY', interval)
      }
      return current
    }
    case 'YEARLY': {
      let current = new Date(start)
      while (current < now) {
        current = advanceDate(current, 'YEARLY', interval)
      }
      return current
    }
    default:
      return start
  }
}

/**
 * Возвращает все вхождения повторяющегося события в заданном месяце.
 * Для неповторяющихся событий — массив с одной датой, если она попадает в месяц.
 */
export function getOccurrencesInMonth(event: CommunityEvent, year: number, month: number): Date[] {
  const monthStart = new Date(year, month, 1)
  const monthEnd = new Date(year, month + 1, 0, 23, 59, 59, 999)

  if (!event.isRepeating || !event.repeatPeriod) {
    const date = new Date(event.date)
    if (date >= monthStart && date <= monthEnd) {
      return [date]
    }
    return []
  }

  if (event.repeatEndDate && new Date(event.repeatEndDate) < monthStart) {
    return []
  }

  const interval = event.repeatInterval ?? 1
  let current = new Date(event.date)

  if (current > monthEnd) {
    return []
  }

  // Перематываем до первого вхождения >= начала месяца
  while (current < monthStart) {
    current = advanceDate(current, event.repeatPeriod, interval)
  }

  const occurrences: Date[] = []
  while (current <= monthEnd) {
    if (event.repeatEndDate && current > new Date(event.repeatEndDate)) {
      break
    }
    occurrences.push(new Date(current))
    current = advanceDate(current, event.repeatPeriod, interval)
  }

  return occurrences
}
