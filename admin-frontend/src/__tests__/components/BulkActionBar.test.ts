import { describe, expect, it, vi } from 'vitest'

vi.mock('@/components/ConfirmDialog.vue', () => ({
  default: {
    template: '<div><slot name="trigger" /></div>',
    props: ['title', 'description', 'confirmLabel', 'variant'],
    emits: ['confirm'],
  },
}))

vi.mock('@/components/ui/button', () => ({
  Button: { template: '<button><slot /></button>', props: ['variant', 'size', 'disabled'] },
}))

import { mount } from '@vue/test-utils'
import BulkActionBar from '@/components/BulkActionBar.vue'

describe('BulkActionBar', () => {
  const defaultActions = [
    { label: 'Удалить', variant: 'destructive' as const, handler: vi.fn() },
    { label: 'Архивировать', handler: vi.fn() },
  ]

  it('does not render when count is 0', () => {
    const wrapper = mount(BulkActionBar, {
      props: { count: 0, actions: defaultActions },
    })
    expect(wrapper.find('div').exists()).toBe(false)
  })

  it('renders when count > 0', () => {
    const wrapper = mount(BulkActionBar, {
      props: { count: 3, actions: defaultActions },
    })
    expect(wrapper.text()).toContain('Выбрано: 3')
  })

  it('renders action buttons for each action', () => {
    const wrapper = mount(BulkActionBar, {
      props: { count: 2, actions: defaultActions },
    })
    expect(wrapper.text()).toContain('Удалить')
    expect(wrapper.text()).toContain('Архивировать')
  })

  it('renders clear selection button', () => {
    const wrapper = mount(BulkActionBar, {
      props: { count: 1, actions: defaultActions },
    })
    expect(wrapper.text()).toContain('Снять выбор')
  })

  it('emits clear when clear button is clicked', async () => {
    const wrapper = mount(BulkActionBar, {
      props: { count: 1, actions: defaultActions },
    })
    const buttons = wrapper.findAll('button')
    const clearButton = buttons.find(b => b.text() === 'Снять выбор')
    expect(clearButton).toBeTruthy()
    await clearButton!.trigger('click')
    expect(wrapper.emitted('clear')).toHaveLength(1)
  })
})
