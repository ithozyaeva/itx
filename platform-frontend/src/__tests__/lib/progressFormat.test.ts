import { describe, expect, it } from 'vitest'
import { daysUntil, deadlineUrgency, formatDeadline, progressPct } from '@/lib/progressFormat'

describe('progressFormat', () => {
  describe('progressPct', () => {
    it('calculates percentage', () => {
      expect(progressPct(5, 10)).toBe(50)
      expect(progressPct(0, 10)).toBe(0)
      expect(progressPct(10, 10)).toBe(100)
    })

    it('caps at 100', () => {
      expect(progressPct(15, 10)).toBe(100)
    })

    it('rounds to nearest integer', () => {
      expect(progressPct(1, 3)).toBe(33)
      expect(progressPct(2, 3)).toBe(67)
    })

    it('returns 0 for non-positive target', () => {
      expect(progressPct(5, 0)).toBe(0)
      expect(progressPct(5, -1)).toBe(0)
    })
  })

  describe('daysUntil', () => {
    it('returns positive for future', () => {
      const future = new Date(Date.now() + 86400000 * 5).toISOString()
      expect(daysUntil(future)).toBeGreaterThan(0)
    })

    it('returns non-positive for past', () => {
      const past = new Date(Date.now() - 86400000 * 2).toISOString()
      expect(daysUntil(past)).toBeLessThanOrEqual(0)
    })
  })

  describe('formatDeadline', () => {
    it('returns "Истекает" for past dates', () => {
      const past = new Date(Date.now() - 86400000).toISOString()
      expect(formatDeadline(past)).toBe('Истекает')
    })

    it('returns "1 день" when ~1 day remains', () => {
      const tomorrow = new Date(Date.now() + 86400000 * 0.5).toISOString()
      expect(formatDeadline(tomorrow)).toBe('1 день')
    })

    it('returns "N дн." for 2-7 days remaining', () => {
      const future = new Date(Date.now() + 86400000 * 4.5).toISOString()
      expect(formatDeadline(future)).toBe('5 дн.')
    })

    it('returns formatted date for >7 days', () => {
      const farFuture = new Date(Date.now() + 86400000 * 30).toISOString()
      const result = formatDeadline(farFuture)
      // formatShortDate возвращает локализованную строку «N <месяц> YYYY»
      expect(result).toMatch(/\d{4}/)
    })
  })

  describe('deadlineUrgency', () => {
    it('returns expired for past', () => {
      expect(deadlineUrgency(new Date(Date.now() - 86400000).toISOString())).toBe('expired')
    })

    it('returns critical for ≤1 day', () => {
      expect(deadlineUrgency(new Date(Date.now() + 86400000 * 0.5).toISOString())).toBe('critical')
    })

    it('returns warning for 2-3 days', () => {
      expect(deadlineUrgency(new Date(Date.now() + 86400000 * 2.5).toISOString())).toBe('warning')
    })

    it('returns normal for >3 days', () => {
      expect(deadlineUrgency(new Date(Date.now() + 86400000 * 10).toISOString())).toBe('normal')
    })
  })
})
