import { flushPromises, mount } from '@vue/test-utils'
import { beforeEach, describe, expect, it, vi } from 'vitest'

const { mockGetLeaderboard, mockUserValue } = vi.hoisted(() => ({
  mockGetLeaderboard: vi.fn(),
  mockUserValue: { id: 1, firstName: 'Test', lastName: 'User' } as { id: number, firstName: string, lastName: string } | null,
}))

vi.mock('lucide-vue-next', () => ({
  Loader2: { template: '<span class="loader" />' },
  Trophy: { template: '<span class="trophy" />' },
}))

vi.mock('@/components/ui/button', () => ({
  Button: {
    template: '<button class="btn" @click="$emit(\'click\')"><slot /></button>',
    props: ['variant'],
    emits: ['click'],
  },
}))

vi.mock('@/components/common/ErrorState.vue', () => ({
  default: {
    template: '<div class="error-state"><slot />{{ message }}</div>',
    props: ['message'],
    emits: ['retry'],
  },
}))

vi.mock('@/composables/useUser', () => ({
  useUser: () => ({ value: mockUserValue }),
}))

vi.mock('@/services/errorService', () => ({
  handleError: vi.fn(async () => ({ message: 'Network error' })),
}))

vi.mock('@/services/points', () => ({
  pointsService: {
    getLeaderboard: mockGetLeaderboard,
  },
}))

import LeaderboardPanel from '@/components/progress/LeaderboardPanel.vue'

const stubRouterLink = {
  template: '<a class="router-link" :data-to="to"><slot /></a>',
  props: ['to'],
}

function makeEntries(n: number, currentId = 1) {
  return Array.from({ length: n }, (_, i) => ({
    memberId: i === 0 ? currentId : i + 100,
    firstName: `User${i}`,
    lastName: `Last${i}`,
    tg: `user${i}`,
    avatarUrl: '',
    total: 1000 - i * 10,
  }))
}

describe('LeaderboardPanel', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  it('shows loader initially', () => {
    mockGetLeaderboard.mockReturnValue(new Promise(() => {}))
    const wrapper = mount(LeaderboardPanel, {
      global: { stubs: { RouterLink: stubRouterLink } },
    })
    expect(wrapper.find('.loader').exists()).toBe(true)
  })

  it('shows ErrorState when fetch fails', async () => {
    mockGetLeaderboard.mockRejectedValue(new Error('Boom'))
    const wrapper = mount(LeaderboardPanel, {
      global: { stubs: { RouterLink: stubRouterLink } },
    })
    await flushPromises()
    expect(wrapper.find('.error-state').exists()).toBe(true)
  })

  it('shows empty hint when no entries', async () => {
    mockGetLeaderboard.mockResolvedValue({ items: [] })
    const wrapper = mount(LeaderboardPanel, {
      global: { stubs: { RouterLink: stubRouterLink } },
    })
    await flushPromises()
    expect(wrapper.text()).toContain('Пока нет данных о баллах')
  })

  it('renders entries and highlights current user', async () => {
    mockGetLeaderboard.mockResolvedValue({ items: makeEntries(5, 1) })
    const wrapper = mount(LeaderboardPanel, {
      global: { stubs: { RouterLink: stubRouterLink } },
    })
    await flushPromises()
    const links = wrapper.findAll('.router-link')
    expect(links.length).toBeGreaterThanOrEqual(5)
    expect(wrapper.text()).toContain('User0')
    expect(wrapper.text()).toContain('1000')
  })

  it('shows "Your position" banner when current user is below the visible page', async () => {
    // Текущий пользователь на 60-й позиции, дефолтный PAGE_SIZE=50 — баннер виден.
    const entries = makeEntries(80, 999)
    entries[59].memberId = 1 // current user at rank 60
    mockGetLeaderboard.mockResolvedValue({ items: entries })
    const wrapper = mount(LeaderboardPanel, {
      global: { stubs: { RouterLink: stubRouterLink } },
    })
    await flushPromises()
    expect(wrapper.text()).toContain('Ваша позиция')
  })

  it('does not show "Your position" banner when current user is on the visible page', async () => {
    const entries = makeEntries(20, 1)
    mockGetLeaderboard.mockResolvedValue({ items: entries })
    const wrapper = mount(LeaderboardPanel, {
      global: { stubs: { RouterLink: stubRouterLink } },
    })
    await flushPromises()
    expect(wrapper.text()).not.toContain('Ваша позиция')
  })

  it('"Show more" button appears only when there are more than PAGE_SIZE entries', async () => {
    mockGetLeaderboard.mockResolvedValue({ items: makeEntries(60, 1) })
    const wrapper = mount(LeaderboardPanel, {
      global: { stubs: { RouterLink: stubRouterLink } },
    })
    await flushPromises()
    const showMoreBtn = wrapper.findAll('.btn').find(b => b.text().includes('Показать ещё'))
    expect(showMoreBtn).toBeDefined()

    await showMoreBtn!.trigger('click')
    await flushPromises()
    // После клика подгрузилось ещё 50, теперь видно >= 60 ссылок (без баннера).
    expect(wrapper.findAll('.router-link').length).toBeGreaterThanOrEqual(60)
  })
})
