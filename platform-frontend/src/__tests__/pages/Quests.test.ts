import { describe, expect, it } from 'vitest'

describe('Quests logic', () => {
  const filters = [
    { key: 'all', label: 'Все' },
    { key: 'active', label: 'Активные' },
    { key: 'completed', label: 'Выполненные' },
  ]

  interface Quest {
    currentCount: number
    targetCount: number
    questType: string
    completed: boolean
    endsAt: string
  }

  function questProgress(quest: { currentCount: number, targetCount: number }) {
    return Math.min(100, Math.round((quest.currentCount / quest.targetCount) * 100))
  }

  function questDeadlineDays(dateStr: string) {
    const date = new Date(dateStr)
    const now = new Date()
    return Math.ceil((date.getTime() - now.getTime()) / (1000 * 60 * 60 * 24))
  }

  function formatQuestDeadline(dateStr: string) {
    const days = questDeadlineDays(dateStr)
    if (days <= 0)
      return 'Истекает'
    if (days === 1)
      return '1 день'
    if (days <= 7)
      return `${days} дн.`
    return new Date(dateStr).toLocaleDateString('ru-RU', { day: '2-digit', month: '2-digit' })
  }

  function questTypeLabel(quest: { questType: string, currentCount: number, targetCount: number }) {
    if (quest.questType === 'daily_streak')
      return `${quest.currentCount} / ${quest.targetCount} дней подряд`
    return `${quest.currentCount} / ${quest.targetCount} сообщений`
  }

  function filteredQuests(quests: Quest[], activeFilter: 'all' | 'active' | 'completed') {
    if (activeFilter === 'all')
      return quests
    if (activeFilter === 'active')
      return quests.filter(q => !q.completed)
    return quests.filter(q => q.completed)
  }

  describe('filters', () => {
    it('has 3 filter options', () => {
      expect(filters).toHaveLength(3)
    })

    it('has correct keys', () => {
      expect(filters.map(f => f.key)).toEqual(['all', 'active', 'completed'])
    })
  })

  describe('questProgress', () => {
    it('calculates percentage correctly', () => {
      expect(questProgress({ currentCount: 5, targetCount: 10 })).toBe(50)
      expect(questProgress({ currentCount: 10, targetCount: 10 })).toBe(100)
      expect(questProgress({ currentCount: 0, targetCount: 10 })).toBe(0)
    })

    it('caps at 100%', () => {
      expect(questProgress({ currentCount: 15, targetCount: 10 })).toBe(100)
    })

    it('handles fractional progress', () => {
      expect(questProgress({ currentCount: 1, targetCount: 3 })).toBe(33)
      expect(questProgress({ currentCount: 2, targetCount: 3 })).toBe(67)
    })
  })

  describe('questDeadlineDays', () => {
    it('returns positive days for future dates', () => {
      const future = new Date(Date.now() + 86400000 * 5).toISOString()
      expect(questDeadlineDays(future)).toBeGreaterThan(0)
    })

    it('returns negative or zero for past dates', () => {
      const past = new Date(Date.now() - 86400000 * 2).toISOString()
      expect(questDeadlineDays(past)).toBeLessThanOrEqual(0)
    })
  })

  describe('formatQuestDeadline', () => {
    it('returns "Истекает" for past dates', () => {
      const past = new Date(Date.now() - 86400000).toISOString()
      expect(formatQuestDeadline(past)).toBe('Истекает')
    })

    it('returns "1 день" for 1 day remaining', () => {
      const tomorrow = new Date(Date.now() + 86400000 * 0.5).toISOString()
      expect(formatQuestDeadline(tomorrow)).toBe('1 день')
    })

    it('returns days format for 2-7 days', () => {
      const future = new Date(Date.now() + 86400000 * 4.5).toISOString()
      expect(formatQuestDeadline(future)).toBe('5 дн.')
    })

    it('returns formatted date for more than 7 days', () => {
      const farFuture = new Date(Date.now() + 86400000 * 30).toISOString()
      const result = formatQuestDeadline(farFuture)
      expect(result).toMatch(/^\d{2}\.\d{2}$/)
    })
  })

  describe('questTypeLabel', () => {
    it('returns daily_streak label', () => {
      expect(questTypeLabel({ questType: 'daily_streak', currentCount: 3, targetCount: 7 }))
        .toBe('3 / 7 дней подряд')
    })

    it('returns messages label for other types', () => {
      expect(questTypeLabel({ questType: 'messages', currentCount: 5, targetCount: 20 }))
        .toBe('5 / 20 сообщений')
    })
  })

  describe('filteredQuests', () => {
    const quests: Quest[] = [
      { currentCount: 5, targetCount: 10, questType: 'messages', completed: false, endsAt: '' },
      { currentCount: 10, targetCount: 10, questType: 'daily_streak', completed: true, endsAt: '' },
      { currentCount: 3, targetCount: 7, questType: 'messages', completed: false, endsAt: '' },
    ]

    it('returns all quests for "all" filter', () => {
      expect(filteredQuests(quests, 'all')).toHaveLength(3)
    })

    it('returns only active quests for "active" filter', () => {
      const result = filteredQuests(quests, 'active')
      expect(result).toHaveLength(2)
      expect(result.every(q => !q.completed)).toBe(true)
    })

    it('returns only completed quests for "completed" filter', () => {
      const result = filteredQuests(quests, 'completed')
      expect(result).toHaveLength(1)
      expect(result.every(q => q.completed)).toBe(true)
    })

    it('returns empty array when no quests match filter', () => {
      const activeOnly: Quest[] = [
        { currentCount: 1, targetCount: 5, questType: 'messages', completed: false, endsAt: '' },
      ]
      expect(filteredQuests(activeOnly, 'completed')).toHaveLength(0)
    })
  })
})
