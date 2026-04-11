import type { CommunityEvent } from '@/models/event'
import { mount } from '@vue/test-utils'
import { describe, expect, it, vi } from 'vitest'

vi.mock('@/components/ui/typography', () => ({
  Typography: { template: '<div><slot /></div>' },
}))

vi.mock('lucide-vue-next', () => ({
  Calendar: { template: '<span />' },
  ChevronDown: { template: '<span />' },
  Crown: { template: '<span />' },
  Loader2: { template: '<span />' },
  MapPin: { template: '<span />' },
}))

vi.mock('@/composables/useUser', () => ({
  useUser: () => ({ value: { id: 99, firstName: 'Test', lastName: 'User', roles: ['SUBSCRIBER'] } }),
}))

vi.mock('@/composables/useDictionary', () => ({
  useDictionary: () => ({
    placeTypesObject: { value: { ONLINE: 'Онлайн', OFFLINE: 'Офлайн', HYBRID: 'Гибрид' } },
    placeTypes: { value: [] },
  }),
}))

vi.mock('@/composables/useGoogleCalendar', () => ({
  useGoogleCalendar: () => ({
    openInGoogleCalendar: vi.fn(),
  }),
}))

vi.mock('@/services/events', () => ({
  eventsService: {
    applyEvent: vi.fn(),
    declineEvent: vi.fn(),
  },
}))

vi.mock('@/services/errorService', () => ({
  handleError: vi.fn(),
}))

vi.mock('@/components/ConfirmDialog.vue', () => ({
  default: {
    template: '<div><slot name="trigger" /></div>',
    props: ['title', 'description', 'confirmLabel'],
  },
}))

import EventCard from '@/components/events/EventCard.vue'

function createEvent(overrides: Partial<CommunityEvent> = {}): CommunityEvent {
  return {
    id: 1,
    title: 'Test Event',
    description: 'Event description here',
    date: new Date(Date.now() + 86400000).toISOString(), // tomorrow
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
    hosts: [{ id: 1, telegramID: 1, tg: 'host', birthday: '', firstName: 'Host', lastName: 'User', bio: '', grade: '', company: '', avatarUrl: '', roles: ['SUBSCRIBER'] }],
    members: [],
    eventTags: [{ id: 1, name: 'Go' }],
    ...overrides,
  }
}

describe('EventCard', () => {
  const globalConfig = {
    stubs: {
      Button: { template: '<button :disabled="$attrs.disabled"><slot /></button>' },
    },
  }

  it('renders event title and description', () => {
    const wrapper = mount(EventCard, {
      props: { event: createEvent() },
      global: globalConfig,
    })
    expect(wrapper.text()).toContain('Test Event')
    expect(wrapper.text()).toContain('Event description here')
  })

  it('renders event tags', () => {
    const wrapper = mount(EventCard, {
      props: { event: createEvent({ eventTags: [{ id: 1, name: 'Go' }, { id: 2, name: 'Vue' }] }) },
      global: globalConfig,
    })
    expect(wrapper.text()).toContain('Go')
    expect(wrapper.text()).toContain('Vue')
  })

  it('displays hosts', () => {
    const wrapper = mount(EventCard, {
      props: { event: createEvent() },
      global: globalConfig,
    })
    expect(wrapper.text()).toContain('Host User')
  })

  it('shows "Участвую!" button for non-member non-host', () => {
    const wrapper = mount(EventCard, {
      props: { event: createEvent() },
      global: globalConfig,
    })
    const buttons = wrapper.findAll('button')
    const applyBtn = buttons.find(b => b.text().includes('Участвую!'))
    expect(applyBtn).toBeDefined()
  })

  it('shows "Отменить участие" button when user is a member', () => {
    const event = createEvent({
      members: [{ id: 99, telegramID: 99, tg: 'test', birthday: '', firstName: 'Test', lastName: 'User', bio: '', grade: '', company: '', avatarUrl: '', roles: ['SUBSCRIBER'] }],
    })
    const wrapper = mount(EventCard, {
      props: { event },
      global: globalConfig,
    })
    expect(wrapper.text()).toContain('Отменить участие')
  })

  it('does not show action buttons for past events', () => {
    const event = createEvent({ date: new Date(Date.now() - 86400000).toISOString() })
    const wrapper = mount(EventCard, {
      props: { event },
      global: globalConfig,
    })
    const applyBtn = wrapper.findAll('button').find(b => b.text().includes('Участвую!'))
    expect(applyBtn).toBeUndefined()
  })

  it('shows "Мест нет" when event is full', () => {
    const event = createEvent({
      maxParticipants: 1,
      members: [{ id: 50, telegramID: 50, tg: 'other', birthday: '', firstName: 'Other', lastName: 'User', bio: '', grade: '', company: '', avatarUrl: '', roles: ['SUBSCRIBER'] }],
    })
    const wrapper = mount(EventCard, {
      props: { event },
      global: globalConfig,
    })
    const btn = wrapper.findAll('button').find(b => b.text().includes('Мест нет'))
    expect(btn).toBeDefined()
    expect(btn!.attributes('disabled')).toBeDefined()
  })

  it('shows exclusive badge when event has exclusiveChatId', () => {
    const event = createEvent({ exclusiveChatId: 123, exclusiveChatTitle: 'VIP Chat' })
    const wrapper = mount(EventCard, {
      props: { event },
      global: globalConfig,
    })
    expect(wrapper.text()).toContain('VIP Chat')
  })

  it('shows member count', async () => {
    const event = createEvent({
      members: [
        { id: 10, telegramID: 10, tg: 'a', birthday: '', firstName: 'Alice', lastName: 'A', bio: '', grade: '', company: '', avatarUrl: '', roles: ['SUBSCRIBER'] },
        { id: 11, telegramID: 11, tg: 'b', birthday: '', firstName: 'Bob', lastName: 'B', bio: '', grade: '', company: '', avatarUrl: '', roles: ['SUBSCRIBER'] },
      ],
    })
    const wrapper = mount(EventCard, {
      props: { event },
      global: globalConfig,
    })
    expect(wrapper.text()).toContain('Участники (2)')
  })

  it('shows members with max participants when set', () => {
    const event = createEvent({ maxParticipants: 10 })
    const wrapper = mount(EventCard, {
      props: { event },
      global: globalConfig,
    })
    expect(wrapper.text()).toContain('Участники (0/10)')
  })

  it('shows timezone', () => {
    const wrapper = mount(EventCard, {
      props: { event: createEvent({ timezone: 'Europe/Moscow' }) },
      global: globalConfig,
    })
    expect(wrapper.text()).toContain('Europe/Moscow')
  })

  it('shows repeat info for repeating events', () => {
    const event = createEvent({
      isRepeating: true,
      repeatPeriod: 'WEEKLY',
      repeatInterval: 1,
    })
    const wrapper = mount(EventCard, {
      props: { event },
      global: globalConfig,
    })
    expect(wrapper.text()).toContain('Повторяется')
  })

  it('shows online place as link', () => {
    const event = createEvent({ placeType: 'ONLINE', place: 'https://zoom.us/j/123' })
    const wrapper = mount(EventCard, {
      props: { event },
      global: globalConfig,
    })
    const link = wrapper.find('a[href="https://zoom.us/j/123"]')
    expect(link.exists()).toBe(true)
  })

  it('shows offline place as text', () => {
    const event = createEvent({ placeType: 'OFFLINE', place: 'Moscow, Red Square 1' })
    const wrapper = mount(EventCard, {
      props: { event },
      global: globalConfig,
    })
    expect(wrapper.text()).toContain('Moscow, Red Square 1')
  })

  it('shows video link when present', () => {
    const event = createEvent({ videoLink: 'https://youtube.com/live' })
    const wrapper = mount(EventCard, {
      props: { event },
      global: globalConfig,
    })
    const link = wrapper.find('a[href="https://youtube.com/live"]')
    expect(link.exists()).toBe(true)
  })

  it('shows recording url when present', () => {
    const event = createEvent({ recordingUrl: 'https://youtube.com/recording' })
    const wrapper = mount(EventCard, {
      props: { event },
      global: globalConfig,
    })
    const link = wrapper.find('a[href="https://youtube.com/recording"]')
    expect(link.exists()).toBe(true)
  })
})
