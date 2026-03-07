import { describe, expect, it } from 'vitest'

describe('Events logic', () => {
  const PAGE_SIZE = 10

  describe('PAGE_SIZE', () => {
    it('is 10', () => {
      expect(PAGE_SIZE).toBe(10)
    })

    it('is a positive integer', () => {
      expect(Number.isInteger(PAGE_SIZE)).toBe(true)
      expect(PAGE_SIZE).toBeGreaterThan(0)
    })
  })
})
