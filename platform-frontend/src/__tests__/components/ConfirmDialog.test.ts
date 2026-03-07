import { mount } from '@vue/test-utils'
import { describe, expect, it, vi } from 'vitest'

vi.mock('@/components/ui/alert-dialog', () => ({
  AlertDialog: { template: '<div class="alert-dialog"><slot /></div>' },
  AlertDialogAction: { template: '<div class="alert-dialog-action"><slot /></div>' },
  AlertDialogCancel: { template: '<div class="alert-dialog-cancel"><slot /></div>' },
  AlertDialogContent: { template: '<div class="alert-dialog-content"><slot /></div>' },
  AlertDialogDescription: { template: '<div class="alert-dialog-description"><slot /></div>' },
  AlertDialogFooter: { template: '<div class="alert-dialog-footer"><slot /></div>' },
  AlertDialogHeader: { template: '<div class="alert-dialog-header"><slot /></div>' },
  AlertDialogTitle: { template: '<div class="alert-dialog-title"><slot /></div>' },
  AlertDialogTrigger: { template: '<div class="alert-dialog-trigger"><slot /></div>' },
}))

vi.mock('@/components/ui/button', () => ({
  Button: {
    inheritAttrs: false,
    template: '<button :class="$attrs.class" @click="$emit(\'click\')"><slot /></button>',
    props: ['variant'],
  },
}))

import ConfirmDialog from '@/components/ConfirmDialog.vue'

describe('ConfirmDialog', () => {
  it('renders without errors', () => {
    const wrapper = mount(ConfirmDialog)
    expect(wrapper.exists()).toBe(true)
  })

  it('displays default title', () => {
    const wrapper = mount(ConfirmDialog)
    expect(wrapper.text()).toContain('Вы уверены?')
  })

  it('displays default description', () => {
    const wrapper = mount(ConfirmDialog)
    expect(wrapper.text()).toContain('Это действие нельзя отменить.')
  })

  it('displays default confirm label', () => {
    const wrapper = mount(ConfirmDialog)
    expect(wrapper.text()).toContain('Подтвердить')
  })

  it('displays custom title', () => {
    const wrapper = mount(ConfirmDialog, {
      props: { title: 'Удалить запись?' },
    })
    expect(wrapper.text()).toContain('Удалить запись?')
  })

  it('displays custom description', () => {
    const wrapper = mount(ConfirmDialog, {
      props: { description: 'Запись будет удалена навсегда.' },
    })
    expect(wrapper.text()).toContain('Запись будет удалена навсегда.')
  })

  it('displays custom confirm label', () => {
    const wrapper = mount(ConfirmDialog, {
      props: { confirmLabel: 'Удалить' },
    })
    expect(wrapper.text()).toContain('Удалить')
  })

  it('displays cancel button', () => {
    const wrapper = mount(ConfirmDialog)
    expect(wrapper.text()).toContain('Отмена')
  })

  it('emits confirm when confirm button is clicked', async () => {
    const wrapper = mount(ConfirmDialog)
    const actionArea = wrapper.find('.alert-dialog-action')
    const confirmBtn = actionArea.find('button')
    await confirmBtn.trigger('click')

    expect(wrapper.emitted('confirm')).toBeDefined()
    expect(wrapper.emitted('confirm')!.length).toBe(1)
  })

  it('renders trigger slot content', () => {
    const wrapper = mount(ConfirmDialog, {
      slots: {
        trigger: '<span class="custom-trigger">Open</span>',
      },
    })
    expect(wrapper.find('.custom-trigger').exists()).toBe(true)
    expect(wrapper.text()).toContain('Open')
  })
})
