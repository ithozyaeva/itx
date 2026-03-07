import { describe, expect, it } from 'vitest'

describe('Casino page logic', () => {
  // diceMultiplier computation
  function calcDiceMultiplier(target: number, direction: 'over' | 'under') {
    const chance = direction === 'over'
      ? (100 - target) / 100
      : target / 100
    if (chance <= 0)
      return 0
    return Math.round((0.97 / chance) * 100) / 100
  }

  describe('diceMultiplier', () => {
    it('calculates multiplier for over 50 (50% chance)', () => {
      expect(calcDiceMultiplier(50, 'over')).toBe(1.94)
    })

    it('calculates multiplier for under 50 (50% chance)', () => {
      expect(calcDiceMultiplier(50, 'under')).toBe(1.94)
    })

    it('calculates higher multiplier for riskier bets', () => {
      const safe = calcDiceMultiplier(25, 'over') // 75% chance
      const risky = calcDiceMultiplier(75, 'over') // 25% chance
      expect(risky).toBeGreaterThan(safe)
    })

    it('returns 0 when chance is 0', () => {
      expect(calcDiceMultiplier(100, 'over')).toBe(0)
    })

    it('calculates correctly for edge targets', () => {
      // over 2 = 98% chance
      const mult = calcDiceMultiplier(2, 'over')
      expect(mult).toBeCloseTo(0.99, 1)

      // over 98 = 2% chance
      const mult2 = calcDiceMultiplier(98, 'over')
      expect(mult2).toBeCloseTo(48.5, 0)
    })
  })

  // gameLabel function
  describe('gameLabel', () => {
    const labels: Record<string, string> = {
      'coin-flip': 'Монетка',
      'dice-roll': 'Кости',
      'wheel': 'Колесо',
    }

    function gameLabel(game: string) {
      return labels[game] ?? game
    }

    it('maps coin-flip', () => {
      expect(gameLabel('coin-flip')).toBe('Монетка')
    })

    it('maps dice-roll', () => {
      expect(gameLabel('dice-roll')).toBe('Кости')
    })

    it('maps wheel', () => {
      expect(gameLabel('wheel')).toBe('Колесо')
    })

    it('falls back to raw game name for unknown game', () => {
      expect(gameLabel('slots')).toBe('slots')
    })
  })

  // formatDate function
  describe('formatDate', () => {
    function formatDate(dateStr: string) {
      const d = new Date(dateStr)
      return d.toLocaleString('ru-RU', { day: '2-digit', month: '2-digit', hour: '2-digit', minute: '2-digit' })
    }

    it('formats a valid ISO date', () => {
      const result = formatDate('2026-03-07T14:30:00Z')
      expect(result).toBeTruthy()
      expect(typeof result).toBe('string')
    })
  })

  // quickBets array
  describe('quickBets', () => {
    const quickBets = [10, 25, 50, 100, 200]

    it('has 5 options', () => {
      expect(quickBets).toHaveLength(5)
    })

    it('starts with min bet 10', () => {
      expect(quickBets[0]).toBe(10)
    })

    it('ends with max bet 200', () => {
      expect(quickBets[quickBets.length - 1]).toBe(200)
    })

    it('is sorted ascending', () => {
      for (let i = 1; i < quickBets.length; i++) {
        expect(quickBets[i]).toBeGreaterThan(quickBets[i - 1])
      }
    })
  })

  // History prepend logic
  describe('history management', () => {
    it('prepends new result and limits to 20', () => {
      const history = Array.from({ length: 20 }, (_, i) => ({ id: i + 1 }))
      const newResult = { id: 100 }

      const updated = [newResult, ...history.slice(0, 19)]

      expect(updated).toHaveLength(20)
      expect(updated[0]).toEqual({ id: 100 })
      expect(updated[19]).toEqual({ id: 19 })
    })

    it('works with empty history', () => {
      const history: { id: number }[] = []
      const newResult = { id: 1 }

      const updated = [newResult, ...history.slice(0, 19)]

      expect(updated).toHaveLength(1)
      expect(updated[0]).toEqual({ id: 1 })
    })
  })

  // Profit display logic
  describe('profit display', () => {
    it('shows positive profit with +', () => {
      const profit = 100
      const display = profit > 0 ? `+${profit}` : `${profit}`
      expect(display).toBe('+100')
    })

    it('shows negative profit without +', () => {
      const profit = -50
      const display = profit > 0 ? `+${profit}` : `${profit}`
      expect(display).toBe('-50')
    })

    it('shows zero profit without +', () => {
      const profit = 0
      const display = profit > 0 ? `+${profit}` : `${profit}`
      expect(display).toBe('0')
    })
  })
})
