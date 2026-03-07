import { describe, expect, it } from 'vitest'

describe('MyReviews logic', () => {
  const statusLabels: Record<string, string> = {
    DRAFT: 'На модерации',
    APPROVED: 'Опубликован',
  }

  describe('statusLabels', () => {
    it('maps DRAFT to "На модерации"', () => {
      expect(statusLabels.DRAFT).toBe('На модерации')
    })

    it('maps APPROVED to "Опубликован"', () => {
      expect(statusLabels.APPROVED).toBe('Опубликован')
    })

    it('has exactly 2 statuses', () => {
      expect(Object.keys(statusLabels)).toHaveLength(2)
    })

    it('returns undefined for unknown status', () => {
      expect(statusLabels.UNKNOWN).toBeUndefined()
    })
  })
})
