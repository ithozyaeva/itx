import { flushPromises, mount } from '@vue/test-utils'
import { beforeEach, describe, expect, it, vi } from 'vitest'

const { mockGetMyStats, mockGetLeaderboard } = vi.hoisted(() => ({
  mockGetMyStats: vi.fn(),
  mockGetLeaderboard: vi.fn(),
}))

vi.mock('itx-ui-kit', () => ({
  Typography: { template: '<div><slot /></div>', props: ['variant', 'as'] },
}))

vi.mock('lucide-vue-next', () => ({
  ArrowDown: { template: '<span class="arrow-down" />' },
  ArrowUp: { template: '<span class="arrow-up" />' },
  Calendar: { template: '<span />' },
  CheckCircle: { template: '<span />' },
  ClipboardList: { template: '<span />' },
  Heart: { template: '<span />' },
  Loader2: { template: '<span class="loader" />' },
  MessageSquare: { template: '<span />' },
  Mic: { template: '<span />' },
  Minus: { template: '<span class="arrow-same" />' },
  Share2: { template: '<span />' },
  Star: { template: '<span />' },
  TrendingUp: { template: '<span />' },
  Trophy: { template: '<span />' },
}))

vi.mock('@/components/common/ErrorState.vue', () => ({
  default: { template: '<div class="error-state"><slot /></div>', props: ['message'], emits: ['retry'] },
}))

vi.mock('@/composables/useUser', () => ({
  useUser: () => ({ value: { id: 1, firstName: 'Test', lastName: 'User' } }),
}))

vi.mock('@/services/errorService', () => ({
  handleError: vi.fn(async () => ({ message: 'Error occurred' })),
}))

vi.mock('@/services/profileStats', () => ({
  profileStatsService: {
    getMyStats: mockGetMyStats,
  },
}))

vi.mock('@/services/points', () => ({
  pointsService: {
    getLeaderboard: mockGetLeaderboard,
  },
}))

import MyStats from '@/pages/MyStats.vue'

const sampleStats = {
  eventsAttended: 5,
  eventsHosted: 2,
  reviewsCount: 10,
  referralsCount: 3,
  kudosSent: 7,
  kudosReceived: 12,
  tasksCreated: 4,
  tasksDone: 6,
  pointsBalance: 150,
  memberSince: '2025-01-15T00:00:00Z',
  pointsHistory: [
    { month: '2025-10', total: 80 },
    { month: '2025-11', total: 100 },
    { month: '2025-12', total: 150 },
  ],
}

describe('MyStats page', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  it('shows loading spinner initially', () => {
    mockGetMyStats.mockReturnValue(new Promise(() => {}))
    mockGetLeaderboard.mockReturnValue(new Promise(() => {}))
    const wrapper = mount(MyStats)
    expect(wrapper.find('.loader').exists()).toBe(true)
  })

  it('shows ErrorState when fetch fails', async () => {
    mockGetMyStats.mockRejectedValue(new Error('Network error'))
    mockGetLeaderboard.mockResolvedValue({ items: [] })
    const wrapper = mount(MyStats)
    await flushPromises()
    expect(wrapper.find('.error-state').exists()).toBe(true)
  })

  it('renders stat cards when data loaded', async () => {
    mockGetMyStats.mockResolvedValue(sampleStats)
    mockGetLeaderboard.mockResolvedValue({ items: [] })
    const wrapper = mount(MyStats)
    await flushPromises()
    expect(wrapper.text()).toContain('Баллы')
    expect(wrapper.text()).toContain('150')
    expect(wrapper.text()).toContain('Посещено событий')
    expect(wrapper.text()).toContain('5')
    expect(wrapper.text()).toContain('Проведено событий')
    expect(wrapper.text()).toContain('Отзывов')
    expect(wrapper.text()).toContain('Рефералов')
  })

  it('shows "Участник с" date', async () => {
    mockGetMyStats.mockResolvedValue(sampleStats)
    mockGetLeaderboard.mockResolvedValue({ items: [] })
    const wrapper = mount(MyStats)
    await flushPromises()
    expect(wrapper.text()).toContain('Участник с')
  })

  it('shows days count', async () => {
    mockGetMyStats.mockResolvedValue(sampleStats)
    mockGetLeaderboard.mockResolvedValue({ items: [] })
    const wrapper = mount(MyStats)
    await flushPromises()
    expect(wrapper.text()).toContain('дней в сообществе')
  })

  it('shows leaderboard position when available', async () => {
    mockGetMyStats.mockResolvedValue(sampleStats)
    mockGetLeaderboard.mockResolvedValue({
      items: [
        { memberId: 5, points: 200 },
        { memberId: 1, points: 150 },
        { memberId: 3, points: 100 },
      ],
    })
    const wrapper = mount(MyStats)
    await flushPromises()
    expect(wrapper.text()).toContain('#2')
    expect(wrapper.text()).toContain('Место в рейтинге')
  })

  it('does not show leaderboard position when user not in leaderboard', async () => {
    mockGetMyStats.mockResolvedValue(sampleStats)
    mockGetLeaderboard.mockResolvedValue({
      items: [
        { memberId: 5, points: 200 },
        { memberId: 3, points: 100 },
      ],
    })
    const wrapper = mount(MyStats)
    await flushPromises()
    expect(wrapper.text()).not.toContain('Место в рейтинге')
  })

  it('shows points chart when history exists', async () => {
    mockGetMyStats.mockResolvedValue(sampleStats)
    mockGetLeaderboard.mockResolvedValue({ items: [] })
    const wrapper = mount(MyStats)
    await flushPromises()
    expect(wrapper.text()).toContain('Баллы по месяцам')
    expect(wrapper.find('svg').exists()).toBe(true)
  })

  it('shows trend indicator (up) when points increased', async () => {
    mockGetMyStats.mockResolvedValue(sampleStats)
    mockGetLeaderboard.mockResolvedValue({ items: [] })
    const wrapper = mount(MyStats)
    await flushPromises()
    expect(wrapper.text()).toContain('за месяц')
    expect(wrapper.find('.arrow-up').exists()).toBe(true)
  })

  it('shows trend indicator (down) when points decreased', async () => {
    const statsWithDecrease = {
      ...sampleStats,
      pointsHistory: [
        { month: '2025-10', total: 80 },
        { month: '2025-11', total: 150 },
        { month: '2025-12', total: 100 },
      ],
    }
    mockGetMyStats.mockResolvedValue(statsWithDecrease)
    mockGetLeaderboard.mockResolvedValue({ items: [] })
    const wrapper = mount(MyStats)
    await flushPromises()
    expect(wrapper.find('.arrow-down').exists()).toBe(true)
  })

  it('handles empty pointsHistory gracefully', async () => {
    const statsNoHistory = {
      ...sampleStats,
      pointsHistory: [],
    }
    mockGetMyStats.mockResolvedValue(statsNoHistory)
    mockGetLeaderboard.mockResolvedValue({ items: [] })
    const wrapper = mount(MyStats)
    await flushPromises()
    expect(wrapper.text()).not.toContain('Баллы по месяцам')
    expect(wrapper.find('svg').exists()).toBe(false)
  })
})

describe('MyStats logic', () => {
  function formatMonth(m: string) {
    const [_, month] = m.split('-')
    const months = ['Янв', 'Фев', 'Мар', 'Апр', 'Май', 'Июн', 'Июл', 'Авг', 'Сен', 'Окт', 'Ноя', 'Дек']
    return months[Number.parseInt(month) - 1] || m
  }

  function getMostActiveArea(stats: {
    eventsAttended: number
    eventsHosted: number
    tasksCreated: number
    tasksDone: number
    kudosSent: number
    kudosReceived: number
    reviewsCount: number
    referralsCount: number
  }): string {
    const areas = [
      { label: 'события', value: stats.eventsAttended + stats.eventsHosted },
      { label: 'задания', value: stats.tasksCreated + stats.tasksDone },
      { label: 'благодарности', value: stats.kudosSent + stats.kudosReceived },
      { label: 'отзывы', value: stats.reviewsCount },
      { label: 'рефералы', value: stats.referralsCount },
    ]
    const sorted = [...areas].sort((a, b) => b.value - a.value)
    return sorted[0].value > 0 ? sorted[0].label : 'нет активности'
  }

  function pointsTrend(history: { month: string, total: number }[]) {
    if (!history || history.length < 2)
      return null

    const current = history[history.length - 1]?.total ?? 0
    const previous = history[history.length - 2]?.total ?? 0

    if (previous === 0)
      return null

    const change = ((current - previous) / previous) * 100
    return {
      value: Math.abs(Math.round(change)),
      direction: change > 0 ? 'up' : change < 0 ? 'down' : 'same',
      current,
      previous,
    }
  }

  describe('formatMonth', () => {
    it('formats month 01 as Янв', () => {
      expect(formatMonth('2025-01')).toBe('Янв')
    })

    it('formats month 12 as Дек', () => {
      expect(formatMonth('2025-12')).toBe('Дек')
    })

    it('formats month 06 as Июн', () => {
      expect(formatMonth('2025-06')).toBe('Июн')
    })

    it('returns raw value for invalid month', () => {
      expect(formatMonth('2025-13')).toBe('2025-13')
    })
  })

  describe('getMostActiveArea', () => {
    it('returns events when events are highest', () => {
      const stats = {
        eventsAttended: 10,
        eventsHosted: 5,
        tasksCreated: 1,
        tasksDone: 1,
        kudosSent: 1,
        kudosReceived: 1,
        reviewsCount: 1,
        referralsCount: 1,
      }
      expect(getMostActiveArea(stats)).toBe('события')
    })

    it('returns задания when tasks are highest', () => {
      const stats = {
        eventsAttended: 1,
        eventsHosted: 1,
        tasksCreated: 10,
        tasksDone: 10,
        kudosSent: 1,
        kudosReceived: 1,
        reviewsCount: 1,
        referralsCount: 1,
      }
      expect(getMostActiveArea(stats)).toBe('задания')
    })

    it('returns "нет активности" when all zeros', () => {
      const stats = {
        eventsAttended: 0,
        eventsHosted: 0,
        tasksCreated: 0,
        tasksDone: 0,
        kudosSent: 0,
        kudosReceived: 0,
        reviewsCount: 0,
        referralsCount: 0,
      }
      expect(getMostActiveArea(stats)).toBe('нет активности')
    })
  })

  describe('pointsTrend', () => {
    it('returns null when history has less than 2 entries', () => {
      expect(pointsTrend([{ month: '2025-01', total: 100 }])).toBeNull()
      expect(pointsTrend([])).toBeNull()
    })

    it('returns null when previous value is 0', () => {
      expect(pointsTrend([
        { month: '2025-01', total: 0 },
        { month: '2025-02', total: 100 },
      ])).toBeNull()
    })

    it('returns up direction when points increased', () => {
      const result = pointsTrend([
        { month: '2025-01', total: 100 },
        { month: '2025-02', total: 150 },
      ])
      expect(result).not.toBeNull()
      expect(result!.direction).toBe('up')
      expect(result!.value).toBe(50)
    })

    it('returns down direction when points decreased', () => {
      const result = pointsTrend([
        { month: '2025-01', total: 200 },
        { month: '2025-02', total: 100 },
      ])
      expect(result).not.toBeNull()
      expect(result!.direction).toBe('down')
      expect(result!.value).toBe(50)
    })

    it('returns same direction when points unchanged', () => {
      const result = pointsTrend([
        { month: '2025-01', total: 100 },
        { month: '2025-02', total: 100 },
      ])
      expect(result).not.toBeNull()
      expect(result!.direction).toBe('same')
      expect(result!.value).toBe(0)
    })
  })
})
