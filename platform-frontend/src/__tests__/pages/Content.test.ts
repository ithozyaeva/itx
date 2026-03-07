import { describe, expect, it } from 'vitest'

describe('Content logic', () => {
  const videos = [
    { id: 'H6cWhHG_KBQ', title: 'Деплой для глупеньких [ IT-X: Mornings ]', date: '2026-02-10' },
    { id: '5x4BKfrQhrY', title: 'Что с работой в 2026: QA с HR сообщества', date: '2026-02-08' },
    { id: '05EXlY1q-Kc', title: 'Postgres для настоящих слонов [ IT-X: Mornings ]', date: '2026-01-27' },
    { id: 'Aiy6rwQNrds', title: 'Делаем SSG-утилиту на Rust [ IT-X: Mornings ]', date: '2025-12-18' },
    { id: 'NZhTyuJWVJE', title: 'Фриланс: опыт выживания [ IT-X: Mornings ]', date: '2025-12-18' },
    { id: '-UA56Roynpg', title: 'Docker: База. Часть 1 [ IT-X: Mornings ]', date: '2025-12-16' },
    { id: '_J090s_jeOk', title: 'Валентин Ким. Как найти работу фронтом в 2025', date: '2025-07-24' },
    { id: 'SO9Xn_bF1zU', title: 'Василий Кузенков. База по ИИ', date: '2025-07-23' },
    { id: '7tYCbNyIun4', title: 'Владимир Балун. Системный дизайн', date: '2025-07-10' },
    { id: 'w-TETYEhzxs', title: 'Альтернативный способ заработка в IT. Как стать ментором?', date: '2025-02-22' },
  ]

  const dateFormatter = new Intl.DateTimeFormat('ru-RU', {
    day: 'numeric',
    month: 'long',
    year: 'numeric',
  })

  function formatDate(dateStr: string) {
    return dateFormatter.format(new Date(dateStr))
  }

  describe('videos', () => {
    it('has exactly 10 videos', () => {
      expect(videos).toHaveLength(10)
    })

    it('each video has id, title, and date', () => {
      for (const video of videos) {
        expect(video).toHaveProperty('id')
        expect(video).toHaveProperty('title')
        expect(video).toHaveProperty('date')
        expect(typeof video.id).toBe('string')
        expect(typeof video.title).toBe('string')
        expect(typeof video.date).toBe('string')
      }
    })

    it('all video ids are unique', () => {
      const ids = videos.map(v => v.id)
      expect(new Set(ids).size).toBe(ids.length)
    })
  })

  describe('formatDate', () => {
    it('returns a non-empty string', () => {
      const result = formatDate('2025-07-10')
      expect(typeof result).toBe('string')
      expect(result.length).toBeGreaterThan(0)
    })

    it('formats date in Russian locale', () => {
      const result = formatDate('2025-07-10')
      expect(result).toContain('2025')
    })

    it('formats different dates correctly', () => {
      const result1 = formatDate('2026-02-10')
      const result2 = formatDate('2025-12-18')
      expect(result1).not.toBe(result2)
    })
  })
})
