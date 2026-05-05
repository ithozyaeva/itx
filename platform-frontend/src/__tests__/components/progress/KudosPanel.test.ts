import { flushPromises, mount } from '@vue/test-utils'
import { beforeEach, describe, expect, it, vi } from 'vitest'

const { mockGetRecent, mockSend, mockApiGet, mockToast } = vi.hoisted(() => ({
  mockGetRecent: vi.fn(),
  mockSend: vi.fn(),
  mockApiGet: vi.fn(),
  mockToast: vi.fn(),
}))

vi.mock('lucide-vue-next', () => ({
  Heart: { template: '<span class="heart" />' },
  Loader2: { template: '<span class="loader" />' },
  Send: { template: '<span class="send" />' },
}))

vi.mock('@/components/common/EmptyState.vue', () => ({
  default: {
    template: '<div class="empty-state">{{ title }}<button class="empty-action" @click="$emit(\'action\')">{{ actionLabel }}</button></div>',
    props: ['icon', 'title', 'description', 'actionLabel'],
    emits: ['action'],
  },
}))

vi.mock('@/components/common/ErrorState.vue', () => ({
  default: {
    template: '<div class="error-state">{{ message }}</div>',
    props: ['message'],
    emits: ['retry'],
  },
}))

vi.mock('@/components/ui/button', () => ({
  Button: {
    template: '<button class="btn" :disabled="disabled" @click="$emit(\'click\')"><slot /></button>',
    props: ['variant', 'disabled'],
    emits: ['click'],
  },
}))

vi.mock('@/components/ui/dialog', () => ({
  Dialog: {
    template: '<div class="dialog" :data-open="open"><slot /></div>',
    props: ['open'],
    emits: ['update:open'],
  },
  DialogScrollContent: { template: '<div class="dialog-content"><slot /></div>' },
  DialogHeader: { template: '<div class="dialog-header"><slot /></div>' },
  DialogTitle: { template: '<div class="dialog-title"><slot /></div>' },
  DialogFooter: { template: '<div class="dialog-footer"><slot /></div>' },
}))

vi.mock('@/components/ui/select', () => ({
  Select: {
    name: 'Select',
    template: '<div class="select" :data-value="modelValue"><slot /></div>',
    props: ['modelValue'],
    emits: ['update:modelValue'],
  },
  SelectTrigger: { name: 'SelectTrigger', template: '<div class="select-trigger"><slot /></div>' },
  SelectValue: { name: 'SelectValue', template: '<div class="select-value"><slot /></div>', props: ['placeholder'] },
  SelectContent: { name: 'SelectContent', template: '<div class="select-content"><slot /></div>' },
  SelectItem: { name: 'SelectItem', template: '<div class="select-item" :data-value="value"><slot /></div>', props: ['value'] },
}))

vi.mock('@/components/ui/toast', () => ({
  useToast: () => ({ toast: mockToast }),
}))

vi.mock('@/composables/useSSE', () => ({
  useSSE: vi.fn(),
}))

vi.mock('@/composables/useUser', () => ({
  useUser: () => ({ value: { id: 1, firstName: 'Me', lastName: 'Self' } }),
}))

vi.mock('@/services/errorService', () => ({
  handleError: vi.fn(async () => ({ message: 'Boom' })),
}))

vi.mock('@/services/api', () => ({
  apiClient: {
    get: () => ({ json: mockApiGet }),
  },
}))

vi.mock('@/services/kudos', () => ({
  kudosService: {
    getRecent: mockGetRecent,
    send: mockSend,
  },
}))

import KudosPanel from '@/components/progress/KudosPanel.vue'

const stubRouterLink = {
  template: '<a class="router-link" :data-to="to"><slot /></a>',
  props: ['to'],
}

const sampleItems = [
  {
    id: 1,
    fromId: 10,
    fromFirstName: 'Алиса',
    fromLastName: 'А',
    fromAvatarUrl: '',
    toId: 1,
    toFirstName: 'Я',
    toLastName: 'Сам',
    toAvatarUrl: '',
    message: 'Спасибо за помощь',
    createdAt: new Date(Date.now() - 5 * 60 * 1000).toISOString(),
  },
]

describe('KudosPanel', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  it('shows loader initially', () => {
    mockGetRecent.mockReturnValue(new Promise(() => {}))
    const wrapper = mount(KudosPanel, {
      global: { stubs: { RouterLink: stubRouterLink, 'router-link': stubRouterLink } },
    })
    expect(wrapper.find('.loader').exists()).toBe(true)
  })

  it('shows ErrorState when fetch fails', async () => {
    mockGetRecent.mockRejectedValue(new Error('Network'))
    const wrapper = mount(KudosPanel, {
      global: { stubs: { RouterLink: stubRouterLink, 'router-link': stubRouterLink } },
    })
    await flushPromises()
    expect(wrapper.find('.error-state').exists()).toBe(true)
  })

  it('shows empty state when feed is empty', async () => {
    mockGetRecent.mockResolvedValue({ items: [], total: 0 })
    const wrapper = mount(KudosPanel, {
      global: { stubs: { RouterLink: stubRouterLink, 'router-link': stubRouterLink } },
    })
    await flushPromises()
    expect(wrapper.find('.empty-state').exists()).toBe(true)
  })

  it('renders kudos feed', async () => {
    mockGetRecent.mockResolvedValue({ items: sampleItems, total: 1 })
    const wrapper = mount(KudosPanel, {
      global: { stubs: { RouterLink: stubRouterLink, 'router-link': stubRouterLink } },
    })
    await flushPromises()
    expect(wrapper.text()).toContain('Спасибо за помощь')
    expect(wrapper.text()).toContain('Алиса')
  })

  it('opens dialog when "Поблагодарить" button is clicked', async () => {
    mockGetRecent.mockResolvedValue({ items: sampleItems, total: 1 })
    mockApiGet.mockResolvedValue({ items: [{ id: 5, firstName: 'Бо', lastName: 'Б', tg: 'bo' }] })
    const wrapper = mount(KudosPanel, {
      global: { stubs: { RouterLink: stubRouterLink, 'router-link': stubRouterLink } },
    })
    await flushPromises()

    const triggers = wrapper.findAll('button').filter(b => b.text().includes('Поблагодарить'))
    expect(triggers.length).toBeGreaterThan(0)
    await triggers[0].trigger('click')
    await flushPromises()

    expect(wrapper.find('.dialog').attributes('data-open')).toBe('true')
    expect(mockApiGet).toHaveBeenCalled()
  })

  it('disables submit button when message or recipient is empty', async () => {
    mockGetRecent.mockResolvedValue({ items: sampleItems, total: 1 })
    mockApiGet.mockResolvedValue({ items: [] })
    const wrapper = mount(KudosPanel, {
      global: { stubs: { RouterLink: stubRouterLink, 'router-link': stubRouterLink } },
    })
    await flushPromises()
    const trigger = wrapper.findAll('button').find(b => b.text().includes('Поблагодарить'))!
    await trigger.trigger('click')
    await flushPromises()

    const submitBtn = wrapper.find('button[type="submit"]')
    expect(submitBtn.exists()).toBe(true)
    expect(submitBtn.attributes('disabled')).toBeDefined()
  })

  it('calls kudosService.send and refetches on submit', async () => {
    mockGetRecent.mockResolvedValue({ items: sampleItems, total: 1 })
    mockApiGet.mockResolvedValue({ items: [{ id: 5, firstName: 'Бо', lastName: 'Б', tg: 'bo' }] })
    mockSend.mockResolvedValue({})
    const wrapper = mount(KudosPanel, {
      global: { stubs: { RouterLink: stubRouterLink, 'router-link': stubRouterLink } },
    })
    await flushPromises()

    const trigger = wrapper.findAll('button').find(b => b.text().includes('Поблагодарить'))!
    await trigger.trigger('click')
    await flushPromises()

    // Заполняем форму через установку refs изнутри: проще симулировать через v-model.
    // Симулируем выбор получателя: вызываем update:modelValue на Select.
    const selectStub = wrapper.findComponent({ name: 'Select' })
    expect(selectStub.exists()).toBe(true)
    selectStub.vm.$emit('update:modelValue', '5')
    await flushPromises()

    const textarea = wrapper.find('textarea')
    expect(textarea.exists()).toBe(true)
    await textarea.setValue('Спасибо!')

    mockGetRecent.mockClear()
    mockGetRecent.mockResolvedValue({ items: sampleItems, total: 1 })

    const form = wrapper.find('form')
    await form.trigger('submit.prevent')
    await flushPromises()

    expect(mockSend).toHaveBeenCalledWith(5, 'Спасибо!')
    expect(mockToast).toHaveBeenCalled()
    // refetch после отправки
    expect(mockGetRecent).toHaveBeenCalled()
  })
})
