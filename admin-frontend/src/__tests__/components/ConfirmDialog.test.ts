import { describe, expect, it, vi } from 'vitest'

vi.mock('@/components/ui/alert-dialog', () => ({
  AlertDialog: { template: '<div><slot /></div>' },
  AlertDialogAction: { template: '<div><slot /></div>' },
  AlertDialogCancel: { template: '<div><slot /></div>' },
  AlertDialogContent: { template: '<div><slot /></div>' },
  AlertDialogDescription: { template: '<div><slot /></div>' },
  AlertDialogFooter: { template: '<div><slot /></div>' },
  AlertDialogHeader: { template: '<div><slot /></div>' },
  AlertDialogTitle: { template: '<div><slot /></div>' },
  AlertDialogTrigger: { template: '<div><slot /></div>' },
}))

vi.mock('@/components/ui/button', () => ({
  Button: { template: '<button><slot /></button>', props: ['variant', 'size', 'disabled'] },
}))

import { mount } from '@vue/test-utils'
import ConfirmDialog from '@/components/ConfirmDialog.vue'

describe('ConfirmDialog', () => {
  it('renders with default props', () => {
    const wrapper = mount(ConfirmDialog)
    expect(wrapper.exists()).toBe(true)
    expect(wrapper.text()).toContain('Вы уверены?')
    expect(wrapper.text()).toContain('Это действие нельзя отменить.')
    expect(wrapper.text()).toContain('Подтвердить')
    expect(wrapper.text()).toContain('Отмена')
  })

  it('renders with custom props', () => {
    const wrapper = mount(ConfirmDialog, {
      props: {
        title: 'Удалить?',
        description: 'Элемент будет удалён.',
        confirmLabel: 'Удалить',
        variant: 'default',
      },
    })
    expect(wrapper.text()).toContain('Удалить?')
    expect(wrapper.text()).toContain('Элемент будет удалён.')
    expect(wrapper.text()).toContain('Удалить')
  })

  it('emits confirm when confirm button is clicked', async () => {
    const wrapper = mount(ConfirmDialog)
    const buttons = wrapper.findAll('button')
    const confirmButton = buttons.find(b => b.text() === 'Подтвердить')
    expect(confirmButton).toBeTruthy()
    await confirmButton!.trigger('click')
    expect(wrapper.emitted('confirm')).toHaveLength(1)
  })

  it('renders trigger slot content', () => {
    const wrapper = mount(ConfirmDialog, {
      slots: {
        trigger: '<span class="trigger-content">Open</span>',
      },
    })
    expect(wrapper.find('.trigger-content').exists()).toBe(true)
  })
})
