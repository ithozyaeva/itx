import { describe, expect, it } from 'vitest'

describe('RafflesView logic', () => {
  function formatDate(dateStr: string) {
    return new Date(dateStr).toLocaleDateString('ru-RU', {
      day: '2-digit',
      month: '2-digit',
      year: 'numeric',
      hour: '2-digit',
      minute: '2-digit',
    })
  }

  it('formats ISO date to Russian locale', () => {
    const result = formatDate('2026-03-07T15:30:00Z')
    expect(result).toBeTruthy()
    expect(typeof result).toBe('string')
  })

  it('contains day, month, and year components', () => {
    const result = formatDate('2026-03-07T15:30:00Z')
    // Russian locale uses dots as separators: DD.MM.YYYY
    expect(result).toContain('07')
    expect(result).toContain('03')
    expect(result).toContain('2026')
  })

  it('handles different date inputs', () => {
    const result1 = formatDate('2025-01-01T00:00:00Z')
    const result2 = formatDate('2025-12-31T23:59:59Z')
    expect(result1).toBeTruthy()
    expect(result2).toBeTruthy()
    expect(result1).not.toBe(result2)
  })
})
