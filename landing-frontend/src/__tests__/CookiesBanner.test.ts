import { mount } from '@vue/test-utils'
import { ref } from 'vue'
import { beforeEach, describe, expect, it, vi } from 'vitest'
import CookiesBanner from '@/components/CookiesBanner.vue'

const reachGoalMock = vi.fn()

vi.mock('yandex-metrika-vue3', () => ({
  useYandexMetrika: () => ({ reachGoal: reachGoalMock }),
}))

vi.mock('@vueuse/core', () => ({
  useLocalStorage: (key: string, initialValue: unknown) => {
    const stored = localStorage.getItem(key)
    return ref(
      stored !== null ? JSON.parse(stored) : initialValue,
    )
  },
}))

vi.mock('itx-ui-kit', () => ({
  Button: {
    name: 'Button',
    template: '<button><slot /></button>',
    props: ['as', 'variant'],
  },
  Typography: {
    name: 'Typography',
    template: '<span><slot /></span>',
    props: ['variant', 'as'],
  },
}))

describe('CookiesBanner', () => {
  beforeEach(() => {
    localStorage.clear()
    reachGoalMock.mockClear()
  })

  it('is visible when cookies have not been accepted', () => {
    const wrapper = mount(CookiesBanner)
    expect(wrapper.find('.fixed').exists()).toBe(true)
  })

  it('is hidden when cookies were previously accepted', () => {
    localStorage.setItem('cookies_accepted', 'true')
    const wrapper = mount(CookiesBanner)
    expect(wrapper.find('.fixed').exists()).toBe(false)
  })

  it('hides the banner and fires analytics when accept button is clicked', async () => {
    const wrapper = mount(CookiesBanner)

    expect(wrapper.find('.fixed').exists()).toBe(true)

    await wrapper.find('button').trigger('click')
    await wrapper.vm.$nextTick()

    expect(reachGoalMock).toHaveBeenCalledWith('cookies_banner_accept')
    expect(wrapper.find('.fixed').exists()).toBe(false)
  })

  it('renders the accept button text', () => {
    const wrapper = mount(CookiesBanner)
    expect(wrapper.find('button').text()).toContain('Понятно')
  })
})
