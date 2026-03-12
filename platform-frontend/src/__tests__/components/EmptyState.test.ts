import { mount } from '@vue/test-utils'
import { describe, expect, it, vi } from 'vitest'

vi.mock('@/components/ui/button', () => ({
  Button: {
    inheritAttrs: false,
    template: '<button @click="$emit(\'click\')"><slot /></button>',
  },
}))

import EmptyState from '@/components/common/EmptyState.vue'

const MockIcon = { template: '<span class="mock-icon" />' }

describe('EmptyState', () => {
  it('renders title text', () => {
    const wrapper = mount(EmptyState, {
      props: { icon: MockIcon, title: 'Нет данных' },
    })
    expect(wrapper.text()).toContain('Нет данных')
  })

  it('renders description when provided', () => {
    const wrapper = mount(EmptyState, {
      props: { icon: MockIcon, title: 'Нет данных', description: 'Попробуйте позже' },
    })
    expect(wrapper.text()).toContain('Попробуйте позже')
  })

  it('hides description when not provided', () => {
    const wrapper = mount(EmptyState, {
      props: { icon: MockIcon, title: 'Нет данных' },
    })
    expect(wrapper.find('p').exists()).toBe(false)
  })

  it('renders action button when actionLabel provided', () => {
    const wrapper = mount(EmptyState, {
      props: { icon: MockIcon, title: 'Нет данных', actionLabel: 'Создать' },
    })
    expect(wrapper.find('button').exists()).toBe(true)
    expect(wrapper.text()).toContain('Создать')
  })

  it('hides action button when actionLabel not provided', () => {
    const wrapper = mount(EmptyState, {
      props: { icon: MockIcon, title: 'Нет данных' },
    })
    expect(wrapper.find('button').exists()).toBe(false)
  })

  it('emits action when button clicked', async () => {
    const wrapper = mount(EmptyState, {
      props: { icon: MockIcon, title: 'Нет данных', actionLabel: 'Создать' },
    })
    await wrapper.find('button').trigger('click')
    expect(wrapper.emitted('action')).toBeDefined()
    expect(wrapper.emitted('action')!.length).toBe(1)
  })

  it('renders icon component', () => {
    const wrapper = mount(EmptyState, {
      props: { icon: MockIcon, title: 'Нет данных' },
    })
    expect(wrapper.find('.mock-icon').exists()).toBe(true)
  })
})
