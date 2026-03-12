import { mount } from '@vue/test-utils'
import { describe, expect, it, vi } from 'vitest'

vi.mock('lucide-vue-next', () => ({
  AlertTriangle: { template: '<span class="alert-triangle" />' },
  RefreshCw: { template: '<span class="refresh" />' },
  WifiOff: { template: '<span class="wifi-off" />' },
}))

vi.mock('@/components/ui/button', () => ({
  Button: {
    inheritAttrs: false,
    template: '<button @click="$emit(\'click\')"><slot /></button>',
    props: ['variant', 'size'],
  },
}))

import ErrorState from '@/components/common/ErrorState.vue'

describe('ErrorState', () => {
  it('renders default message', () => {
    const wrapper = mount(ErrorState)
    expect(wrapper.text()).toContain('Произошла ошибка при загрузке данных')
  })

  it('renders custom message', () => {
    const wrapper = mount(ErrorState, {
      props: { message: 'Сервер недоступен' },
    })
    expect(wrapper.text()).toContain('Сервер недоступен')
  })

  it('shows WifiOff icon when type is network', () => {
    const wrapper = mount(ErrorState, {
      props: { type: 'network' },
    })
    expect(wrapper.find('.wifi-off').exists()).toBe(true)
    expect(wrapper.find('.alert-triangle').exists()).toBe(false)
  })

  it('shows AlertTriangle icon when type is not network', () => {
    const wrapper = mount(ErrorState, {
      props: { type: 'server' },
    })
    expect(wrapper.find('.alert-triangle').exists()).toBe(true)
    expect(wrapper.find('.wifi-off').exists()).toBe(false)
  })

  it('shows "Нет соединения" title for network type', () => {
    const wrapper = mount(ErrorState, {
      props: { type: 'network' },
    })
    expect(wrapper.text()).toContain('Нет соединения')
  })

  it('shows "Ошибка загрузки" title for non-network type', () => {
    const wrapper = mount(ErrorState)
    expect(wrapper.text()).toContain('Ошибка загрузки')
  })

  it('emits retry when button clicked', async () => {
    const wrapper = mount(ErrorState)
    await wrapper.find('button').trigger('click')
    expect(wrapper.emitted('retry')).toBeDefined()
    expect(wrapper.emitted('retry')!.length).toBe(1)
  })

  it('shows retry button text', () => {
    const wrapper = mount(ErrorState)
    expect(wrapper.text()).toContain('Попробовать снова')
  })
})
