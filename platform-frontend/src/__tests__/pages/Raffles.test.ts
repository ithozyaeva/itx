import { describe, expect, it } from 'vitest'

describe('Raffles logic', () => {
  function timeLeft(endsAt: string) {
    const diff = new Date(endsAt).getTime() - Date.now()
    if (diff <= 0) return 'Завершён'
    const hours = Math.floor(diff / 3600000)
    const days = Math.floor(hours / 24)
    if (days > 0) return `${days} дн. ${hours % 24} ч.`
    return `${hours} ч.`
  }

  interface RaffleItem {
    winnerFirstName?: string
    winnerLastName?: string
  }

  function winnerName(r: RaffleItem) {
    return [r.winnerFirstName, r.winnerLastName].filter(Boolean).join(' ')
  }

  function getTicketCount(ticketCounts: Record<number, number>, id: number) {
    return ticketCounts[id] ?? 1
  }

  describe('timeLeft', () => {
    it('returns "Завершён" for past date', () => {
      const past = new Date(Date.now() - 86400000).toISOString()
      expect(timeLeft(past)).toBe('Завершён')
    })

    it('returns hours only when less than a day', () => {
      const fiveHoursLater = new Date(Date.now() + 5 * 3600000 + 60000).toISOString()
      expect(timeLeft(fiveHoursLater)).toBe('5 ч.')
    })

    it('returns days and hours when more than a day', () => {
      const future = new Date(Date.now() + 27 * 3600000 + 60000).toISOString()
      expect(timeLeft(future)).toBe('1 дн. 3 ч.')
    })

    it('returns "Завершён" for exactly now', () => {
      const now = new Date(Date.now() - 1).toISOString()
      expect(timeLeft(now)).toBe('Завершён')
    })

    it('returns "0 ч." when less than one hour remains', () => {
      const soon = new Date(Date.now() + 30 * 60000).toISOString()
      expect(timeLeft(soon)).toBe('0 ч.')
    })
  })

  describe('winnerName', () => {
    it('joins first and last name', () => {
      expect(winnerName({ winnerFirstName: 'Иван', winnerLastName: 'Иванов' })).toBe('Иван Иванов')
    })

    it('handles only first name', () => {
      expect(winnerName({ winnerFirstName: 'Иван', winnerLastName: '' })).toBe('Иван')
    })

    it('handles only last name', () => {
      expect(winnerName({ winnerFirstName: '', winnerLastName: 'Иванов' })).toBe('Иванов')
    })

    it('handles neither name', () => {
      expect(winnerName({ winnerFirstName: '', winnerLastName: '' })).toBe('')
    })

    it('handles undefined names', () => {
      expect(winnerName({})).toBe('')
    })
  })

  describe('getTicketCount', () => {
    it('returns count from map when present', () => {
      expect(getTicketCount({ 1: 5, 2: 3 }, 1)).toBe(5)
    })

    it('defaults to 1 when id not in map', () => {
      expect(getTicketCount({ 1: 5 }, 99)).toBe(1)
    })

    it('defaults to 1 for empty map', () => {
      expect(getTicketCount({}, 1)).toBe(1)
    })
  })
})
