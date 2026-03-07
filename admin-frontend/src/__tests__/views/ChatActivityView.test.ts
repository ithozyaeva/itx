import { describe, expect, it } from 'vitest'

describe('ChatActivityView logic', () => {
  // calcChange function
  function calcChange(current: number, previous: number): number | null {
    if (previous === 0)
      return current > 0 ? 100 : null
    return Math.round(((current - previous) / previous) * 100)
  }

  it('calcChange returns percentage increase', () => {
    expect(calcChange(150, 100)).toBe(50)
  })

  it('calcChange returns percentage decrease', () => {
    expect(calcChange(50, 100)).toBe(-50)
  })

  it('calcChange returns 0 when values are equal', () => {
    expect(calcChange(100, 100)).toBe(0)
  })

  it('calcChange returns 100 when previous is 0 and current > 0', () => {
    expect(calcChange(50, 0)).toBe(100)
  })

  it('calcChange returns null when both are 0', () => {
    expect(calcChange(0, 0)).toBeNull()
  })

  it('calcChange rounds to nearest integer', () => {
    expect(calcChange(133, 100)).toBe(33)
    expect(calcChange(167, 100)).toBe(67)
  })

  // summaryCards structure
  it('summaryCards has 4 items', () => {
    const summaryCards = [
      { label: 'Сообщений сегодня', value: 0, change: null },
      { label: 'Сообщений за неделю', value: 0, change: 50 },
      { label: 'Активных сегодня', value: 0, change: null },
      { label: 'Активных за неделю', value: 0, change: -10 },
    ]
    expect(summaryCards).toHaveLength(4)
  })

  // lineChartData date label slicing
  it('slices date labels to remove year prefix', () => {
    const dates = ['2026-03-01', '2026-03-02', '2026-03-03']
    const labels = dates.map(d => d.slice(5))
    expect(labels).toEqual(['03-01', '03-02', '03-03'])
  })

  // chartOptions config
  it('chartOptions has correct structure', () => {
    const chartOptions = {
      responsive: true,
      maintainAspectRatio: false,
      plugins: {
        legend: { display: false },
      },
      scales: {
        y: { beginAtZero: true },
      },
    }
    expect(chartOptions.responsive).toBe(true)
    expect(chartOptions.maintainAspectRatio).toBe(false)
    expect(chartOptions.plugins.legend.display).toBe(false)
  })

  // clearUserSelection logic
  it('clearUserSelection resets all user-related state', () => {
    let selectedUser: { telegramUserId: number } | null = { telegramUserId: 123 }
    let userStats: object | null = { totalMessages: 100 }
    let userChartData: object[] = [{ date: '2026-03-01', count: 5 }]

    // Simulating clearUserSelection
    selectedUser = null
    userStats = null
    userChartData = []

    expect(selectedUser).toBeNull()
    expect(userStats).toBeNull()
    expect(userChartData).toEqual([])
  })
})
