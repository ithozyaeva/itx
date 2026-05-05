import { flushPromises, mount } from '@vue/test-utils'
import { beforeEach, describe, expect, it, vi } from 'vitest'

const { mockGetMyAchievements } = vi.hoisted(() => ({
  mockGetMyAchievements: vi.fn(),
}))

vi.mock('lucide-vue-next', () => {
  // Все иконки, упомянутые в AchievementsPanel.iconMap + сам Award/Loader2/CheckCircle.
  // Vitest требует именованных экспортов (Proxy не подходит).
  const stub = { template: '<span class="icon" />' }
  const names = [
    'Award', 'BookOpen', 'Briefcase', 'CalendarCheck', 'CheckCircle', 'ClipboardList',
    'Crown', 'FileText', 'Flame', 'Footprints', 'Gem', 'GraduationCap', 'HardHat',
    'History', 'ListChecks', 'Loader2', 'Medal', 'MessageSquare', 'MessagesSquare',
    'Mic', 'Package', 'Presentation', 'Share2', 'ShoppingCart', 'Star', 'Swords',
    'Target', 'Trophy', 'UserCheck', 'UserPlus', 'Users', 'Zap',
  ]
  return Object.fromEntries(names.map(n => [n, stub]))
})

vi.mock('@/components/progress', () => ({
  // ProgressBar/TintedIcon подключаются через `@/components/progress` barrel.
  // Заглушки достаточно: рендер не зависит от их внутренностей.
  ProgressBar: { template: '<div class="progress-bar" />', props: ['progress', 'target', 'state'] },
  TintedIcon: { template: '<span class="tinted-icon" />', props: ['icon', 'tone', 'done'] },
}))

vi.mock('@/components/ui/dialog', () => ({
  Dialog: {
    template: '<div class="dialog" :data-open="open"><slot /></div>',
    props: ['open'],
    emits: ['update:open'],
  },
  DialogContent: { template: '<div class="dialog-content"><slot /></div>' },
  DialogDescription: { template: '<div class="dialog-description"><slot /></div>' },
  DialogHeader: { template: '<div class="dialog-header"><slot /></div>' },
  DialogTitle: { template: '<div class="dialog-title"><slot /></div>' },
}))

vi.mock('@/services/errorService', () => ({
  handleError: vi.fn(async () => ({ message: 'Error' })),
}))

vi.mock('@/services/achievements', () => ({
  achievementsService: { getMyAchievements: mockGetMyAchievements },
}))

import AchievementsPanel from '@/components/progress/AchievementsPanel.vue'

const sampleResponse = {
  totalCount: 4,
  unlockedCount: 1,
  items: [
    { id: 'evt-1', title: 'Первый ивент', description: 'Посети событие', icon: 'calendar-check', category: 'events', threshold: 1, progress: 1, unlocked: true },
    { id: 'evt-2', title: 'Десять событий', description: '10 ивентов', icon: 'calendar-check', category: 'events', threshold: 10, progress: 3, unlocked: false },
    { id: 'pts-1', title: 'Сто баллов', description: 'Накопи 100', icon: 'star', category: 'points', threshold: 100, progress: 60, unlocked: false },
    { id: 'soc-1', title: 'Доброжелатель', description: 'Скажи спасибо', icon: 'heart', category: 'social', threshold: 1, progress: 0, unlocked: false },
  ],
}

describe('AchievementsPanel', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  it('shows loader before data resolves', () => {
    mockGetMyAchievements.mockReturnValue(new Promise(() => {}))
    const wrapper = mount(AchievementsPanel)
    expect(wrapper.find('.icon').exists()).toBe(true) // Loader2 stubbed as icon
  })

  it('renders unlocked counter and all achievements by default', async () => {
    mockGetMyAchievements.mockResolvedValue(sampleResponse)
    const wrapper = mount(AchievementsPanel)
    await flushPromises()
    expect(wrapper.text()).toContain('1')
    expect(wrapper.text()).toContain('/ 4')
    expect(wrapper.text()).toContain('Первый ивент')
    expect(wrapper.text()).toContain('Десять событий')
    expect(wrapper.text()).toContain('Сто баллов')
    expect(wrapper.text()).toContain('Доброжелатель')
  })

  it('filters items by category', async () => {
    mockGetMyAchievements.mockResolvedValue(sampleResponse)
    const wrapper = mount(AchievementsPanel)
    await flushPromises()

    const eventsTab = wrapper.findAll('button[role="tab"]').find(b => b.text() === 'События')!
    await eventsTab.trigger('click')
    await flushPromises()

    expect(wrapper.text()).toContain('Первый ивент')
    expect(wrapper.text()).toContain('Десять событий')
    expect(wrapper.text()).not.toContain('Сто баллов')
    expect(wrapper.text()).not.toContain('Доброжелатель')
  })

  it('opens detail dialog when achievement clicked', async () => {
    mockGetMyAchievements.mockResolvedValue(sampleResponse)
    const wrapper = mount(AchievementsPanel)
    await flushPromises()

    const cards = wrapper.findAll('button').filter(b => b.text().includes('Сто баллов'))
    expect(cards.length).toBeGreaterThan(0)
    await cards[0].trigger('click')
    await flushPromises()

    const dialog = wrapper.find('.dialog')
    expect(dialog.attributes('data-open')).toBe('true')
    expect(wrapper.find('.dialog-title').text()).toContain('Сто баллов')
    expect(wrapper.find('.dialog-description').text()).toContain('Накопи 100')
  })

  it('marks selected category tab as aria-selected', async () => {
    mockGetMyAchievements.mockResolvedValue(sampleResponse)
    const wrapper = mount(AchievementsPanel)
    await flushPromises()

    const allTab = wrapper.findAll('button[role="tab"]').find(b => b.text() === 'Все')!
    expect(allTab.attributes('aria-selected')).toBe('true')

    const pointsTab = wrapper.findAll('button[role="tab"]').find(b => b.text() === 'Баллы')!
    await pointsTab.trigger('click')
    await flushPromises()
    expect(pointsTab.attributes('aria-selected')).toBe('true')
    expect(allTab.attributes('aria-selected')).toBe('false')
  })
})
