import { mount } from '@vue/test-utils'
import { describe, expect, it, vi } from 'vitest'
import { ref } from 'vue'

const mockRoute = ref({
  path: '/',
  meta: {} as Record<string, unknown>,
})

const mockPush = vi.fn()

vi.mock('vue-router', () => ({
  useRoute: () => mockRoute.value,
  useRouter: () => ({ push: mockPush }),
}))

vi.mock('@/composables/useBreadcrumb', () => ({
  useBreadcrumb: () => ({
    dynamicLabel: ref(null),
    clearLabel: vi.fn(),
  }),
}))

vi.mock('lucide-vue-next', () => ({
  ChevronRight: { template: '<span class="chevron" />' },
  Home: { template: '<span class="home-icon" />' },
}))

import Breadcrumbs from '@/components/common/Breadcrumbs.vue'

describe('Breadcrumbs', () => {
  it('renders nothing when no breadcrumb meta', () => {
    mockRoute.value = { path: '/', meta: {} }
    const wrapper = mount(Breadcrumbs)
    expect(wrapper.find('nav').exists()).toBe(false)
  })

  it('renders nav element when breadcrumb meta exists', () => {
    mockRoute.value = {
      path: '/events',
      meta: { breadcrumb: [{ label: 'События' }] },
    }
    const wrapper = mount(Breadcrumbs)
    expect(wrapper.find('nav').exists()).toBe(true)
  })

  it('shows home button', () => {
    mockRoute.value = {
      path: '/events',
      meta: { breadcrumb: [{ label: 'События' }] },
    }
    const wrapper = mount(Breadcrumbs)
    expect(wrapper.find('.home-icon').exists()).toBe(true)
  })

  it('renders breadcrumb labels', () => {
    mockRoute.value = {
      path: '/events/1',
      meta: { breadcrumb: [{ label: 'События', to: '/events' }, { label: 'Детали' }] },
    }
    const wrapper = mount(Breadcrumbs)
    expect(wrapper.text()).toContain('События')
    expect(wrapper.text()).toContain('Детали')
  })

  it('renders clickable links for non-last items with to', async () => {
    mockRoute.value = {
      path: '/events/1',
      meta: { breadcrumb: [{ label: 'События', to: '/events' }, { label: 'Детали' }] },
    }
    const wrapper = mount(Breadcrumbs)
    const buttons = wrapper.findAll('button')
    // First button is home, second is "События" link
    const eventButton = buttons.find(b => b.text().includes('События'))
    expect(eventButton).toBeDefined()
    await eventButton!.trigger('click')
    expect(mockPush).toHaveBeenCalledWith('/events')
  })

  it('last item is not a link (rendered as span)', () => {
    mockRoute.value = {
      path: '/events/1',
      meta: { breadcrumb: [{ label: 'События', to: '/events' }, { label: 'Детали' }] },
    }
    const wrapper = mount(Breadcrumbs)
    const span = wrapper.find('span.text-foreground\\/70')
    expect(span.exists()).toBe(true)
    expect(span.text()).toBe('Детали')
  })

  it('has aria-label="Навигация" on nav', () => {
    mockRoute.value = {
      path: '/events',
      meta: { breadcrumb: [{ label: 'События' }] },
    }
    const wrapper = mount(Breadcrumbs)
    expect(wrapper.find('nav').attributes('aria-label')).toBe('Навигация')
  })

  it('home button has aria-label="Главная"', () => {
    mockRoute.value = {
      path: '/events',
      meta: { breadcrumb: [{ label: 'События' }] },
    }
    const wrapper = mount(Breadcrumbs)
    const homeBtn = wrapper.find('button[aria-label="Главная"]')
    expect(homeBtn.exists()).toBe(true)
  })
})
