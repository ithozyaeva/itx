import type { CommunityEvent } from '@/models/event'
import { mount } from '@vue/test-utils'
import { describe, expect, it, vi } from 'vitest'

vi.mock('lucide-vue-next', () => ({
  ChevronLeft: { template: '<span class="chevron-left" />' },
  ChevronRight: { template: '<span class="chevron-right" />' },
}))

vi.mock('@/components/events/EventCard.vue', () => ({
  default: {
    template: '<div class="event-card">{{ event.title }}</div>',
    props: ['event'],
  },
}))

import CalendarView from '@/components/events/CalendarView.vue'

const monthNames = [
  'Январь', 'Февраль', 'Март', 'Апрель', 'Май', 'Июнь',
  'Июль', 'Август', 'Сентябрь', 'Октябрь', 'Ноябрь', 'Декабрь',
]

function createEvent(overrides: Partial<CommunityEvent> = {}): CommunityEvent {
  return {
    id: 1,
    title: 'Test Event',
    description: 'Description',
    date: new Date().toISOString(),
    timezone: 'Europe/Moscow',
    placeType: 'ONLINE',
    place: 'https://zoom.us/j/123',
    customPlaceType: '',
    eventType: 'ONLINE',
    open: true,
    videoLink: '',
    isRepeating: false,
    recordingUrl: '',
    maxParticipants: 0,
    exclusiveChatId: null,
    hosts: [],
    members: [],
    eventTags: [],
    ...overrides,
  }
}

describe('CalendarView', () => {
  it('renders month name and year', () => {
    const wrapper = mount(CalendarView, {
      props: { events: [] },
    })
    const today = new Date()
    expect(wrapper.text()).toContain(monthNames[today.getMonth()])
    expect(wrapper.text()).toContain(String(today.getFullYear()))
  })

  it('renders 7 weekday headers', () => {
    const wrapper = mount(CalendarView, {
      props: { events: [] },
    })
    const weekDays = ['Пн', 'Вт', 'Ср', 'Чт', 'Пт', 'Сб', 'Вс']
    for (const day of weekDays) {
      expect(wrapper.text()).toContain(day)
    }
  })

  it('renders 42 day cells (6 weeks)', () => {
    const wrapper = mount(CalendarView, {
      props: { events: [] },
    })
    const grid = wrapper.find('.grid.grid-cols-7.gap-px.bg-border')
    const buttons = grid.findAll('button')
    expect(buttons.length).toBe(42)
  })

  it('highlights today', () => {
    const wrapper = mount(CalendarView, {
      props: { events: [] },
    })
    const todayEl = wrapper.find('.bg-primary.text-primary-foreground')
    expect(todayEl.exists()).toBe(true)
  })

  it('shows event dots for dates with events', () => {
    const today = new Date()
    const event = createEvent({ date: today.toISOString() })
    const wrapper = mount(CalendarView, {
      props: { events: [event] },
    })
    const dots = wrapper.findAll('.rounded-full.bg-primary')
    expect(dots.length).toBeGreaterThan(0)
  })

  it('clicking a day selects it (ring-2 class)', async () => {
    const wrapper = mount(CalendarView, {
      props: { events: [] },
    })
    const grid = wrapper.find('.grid.grid-cols-7.gap-px.bg-border')
    const dayButton = grid.findAll('button')[10]
    await dayButton.trigger('click')
    expect(dayButton.classes()).toContain('ring-2')
  })

  it('clicking selected day deselects it', async () => {
    const wrapper = mount(CalendarView, {
      props: { events: [] },
    })
    const grid = wrapper.find('.grid.grid-cols-7.gap-px.bg-border')
    const dayButton = grid.findAll('button')[10]
    await dayButton.trigger('click')
    expect(dayButton.classes()).toContain('ring-2')
    await dayButton.trigger('click')
    expect(dayButton.classes()).not.toContain('ring-2')
  })

  it('previous/next month buttons navigate', async () => {
    const wrapper = mount(CalendarView, {
      props: { events: [] },
    })
    const today = new Date()
    const currentMonth = today.getMonth()
    const prevMonthIndex = currentMonth === 0 ? 11 : currentMonth - 1

    const prevBtn = wrapper.find('button[aria-label="Предыдущий месяц"]')
    await prevBtn.trigger('click')
    expect(wrapper.text()).toContain(monthNames[prevMonthIndex])

    const nextBtn = wrapper.find('button[aria-label="Следующий месяц"]')
    await nextBtn.trigger('click')
    expect(wrapper.text()).toContain(monthNames[currentMonth])
  })

  it('"Сегодня" button resets to current month', async () => {
    const wrapper = mount(CalendarView, {
      props: { events: [] },
    })
    const today = new Date()

    const prevBtn = wrapper.find('button[aria-label="Предыдущий месяц"]')
    await prevBtn.trigger('click')
    await prevBtn.trigger('click')

    const todayBtn = wrapper.findAll('button').find(b => b.text() === 'Сегодня')
    expect(todayBtn).toBeDefined()
    await todayBtn!.trigger('click')
    expect(wrapper.text()).toContain(monthNames[today.getMonth()])
    expect(wrapper.text()).toContain(String(today.getFullYear()))
  })

  it('shows EventCard list when day with events is selected', async () => {
    const today = new Date()
    const event = createEvent({ id: 1, title: 'My Event', date: today.toISOString() })
    const wrapper = mount(CalendarView, {
      props: { events: [event] },
    })

    // Find the today button (it has the bg-primary span)
    const grid = wrapper.find('.grid.grid-cols-7.gap-px.bg-border')
    const buttons = grid.findAll('button')
    const todayButton = buttons.find(b => b.find('.bg-primary.text-primary-foreground').exists())
    expect(todayButton).toBeDefined()
    await todayButton!.trigger('click')

    expect(wrapper.find('.event-card').exists()).toBe(true)
    expect(wrapper.text()).toContain('My Event')
  })

  it('shows "Нет событий в этот день" for selected empty day', async () => {
    const wrapper = mount(CalendarView, {
      props: { events: [] },
    })
    const grid = wrapper.find('.grid.grid-cols-7.gap-px.bg-border')
    const dayButton = grid.findAll('button')[10]
    await dayButton.trigger('click')
    expect(wrapper.text()).toContain('Нет событий в этот день')
  })

  it('navigation buttons have aria-labels', () => {
    const wrapper = mount(CalendarView, {
      props: { events: [] },
    })
    expect(wrapper.find('button[aria-label="Предыдущий месяц"]').exists()).toBe(true)
    expect(wrapper.find('button[aria-label="Следующий месяц"]').exists()).toBe(true)
  })
})
