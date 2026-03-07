import { describe, expect, it } from 'vitest'

describe('MyPoints logic', () => {
  function questProgress(currentCount: number, targetCount: number) {
    return Math.min(100, Math.round((currentCount / targetCount) * 100))
  }

  function formatQuestDeadline(dateStr: string) {
    const date = new Date(dateStr)
    const now = new Date()
    const diffMs = date.getTime() - now.getTime()
    const diffDays = Math.ceil(diffMs / (1000 * 60 * 60 * 24))
    if (diffDays <= 0) return 'Истекло'
    if (diffDays === 1) return 'Остался 1 день'
    if (diffDays <= 7) return `Осталось ${diffDays} дн.`
    return date.toLocaleDateString('ru-RU', { day: '2-digit', month: '2-digit' })
  }

  describe('questProgress', () => {
    it('calculates percentage correctly', () => {
      expect(questProgress(5, 10)).toBe(50)
      expect(questProgress(10, 10)).toBe(100)
      expect(questProgress(0, 10)).toBe(0)
    })

    it('caps at 100%', () => {
      expect(questProgress(15, 10)).toBe(100)
    })

    it('handles fractional progress', () => {
      expect(questProgress(1, 3)).toBe(33)
      expect(questProgress(2, 3)).toBe(67)
    })
  })

  describe('formatQuestDeadline', () => {
    it('returns "Истекло" for past dates', () => {
      const past = new Date(Date.now() - 86400000).toISOString()
      expect(formatQuestDeadline(past)).toBe('Истекло')
    })

    it('returns "Остался 1 день" for tomorrow', () => {
      const tomorrow = new Date(Date.now() + 86400000 * 0.5).toISOString()
      expect(formatQuestDeadline(tomorrow)).toBe('Остался 1 день')
    })

    it('returns days remaining for near future', () => {
      const future = new Date(Date.now() + 86400000 * 4.5).toISOString()
      expect(formatQuestDeadline(future)).toBe('Осталось 5 дн.')
    })
  })
})
