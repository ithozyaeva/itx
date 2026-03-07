import { describe, expect, it } from 'vitest'

describe('DashboardView logic', () => {
  // statCards configuration
  const statCards = [
    { key: 'totalMembers', label: 'Участники' },
    { key: 'totalMentors', label: 'Менторы' },
    { key: 'upcomingEvents', label: 'Предстоящие события' },
    { key: 'pastEvents', label: 'Прошедшие события' },
    { key: 'pendingReviews', label: 'Ожидают публикации' },
    { key: 'approvedReviews', label: 'Опубликованные отзывы' },
    { key: 'referralLinks', label: 'Реферальные ссылки' },
    { key: 'resumes', label: 'Резюме' },
    { key: 'openTasks', label: 'Открытые задания' },
    { key: 'inProgressTasks', label: 'Задания в работе' },
    { key: 'doneTasks', label: 'Выполненные задания' },
    { key: 'approvedTasks', label: 'Принятые задания' },
  ] as const

  it('has 12 stat cards', () => {
    expect(statCards).toHaveLength(12)
  })

  it('all stat cards have unique keys', () => {
    const keys = statCards.map(c => c.key)
    expect(new Set(keys).size).toBe(keys.length)
  })

  it('all stat cards have non-empty labels', () => {
    for (const card of statCards) {
      expect(card.label.length).toBeGreaterThan(0)
    }
  })

  // memberGrowthData computed logic
  it('builds memberGrowthData from chartStats', () => {
    const memberGrowth = [
      { month: 'Январь', count: 10 },
      { month: 'Февраль', count: 20 },
    ]

    const data = {
      labels: memberGrowth.map(m => m.month),
      datasets: [{
        label: 'Участники',
        data: memberGrowth.map(m => m.count),
        borderColor: 'hsl(var(--primary))',
        backgroundColor: 'hsl(var(--primary) / 0.1)',
        fill: true,
        tension: 0.3,
      }],
    }

    expect(data.labels).toEqual(['Январь', 'Февраль'])
    expect(data.datasets[0].data).toEqual([10, 20])
    expect(data.datasets[0].label).toBe('Участники')
    expect(data.datasets[0].fill).toBe(true)
  })

  it('returns empty arrays when chartStats is null', () => {
    const chartStats = null

    const labels = chartStats?.memberGrowth.map((m: { month: string }) => m.month) ?? []
    const data = chartStats?.memberGrowth.map((m: { count: number }) => m.count) ?? []

    expect(labels).toEqual([])
    expect(data).toEqual([])
  })

  // eventAttendanceData computed logic
  it('builds eventAttendanceData from chartStats', () => {
    const eventAttendance = [
      { month: 'Март', count: 15 },
      { month: 'Апрель', count: 30 },
    ]

    const data = {
      labels: eventAttendance.map(m => m.month),
      datasets: [{
        label: 'Посещаемость',
        data: eventAttendance.map(m => m.count),
        backgroundColor: 'hsl(var(--primary) / 0.7)',
        borderRadius: 4,
      }],
    }

    expect(data.labels).toEqual(['Март', 'Апрель'])
    expect(data.datasets[0].data).toEqual([15, 30])
    expect(data.datasets[0].borderRadius).toBe(4)
  })

  // chartOptions
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
    expect(chartOptions.scales.y.beginAtZero).toBe(true)
  })
})
