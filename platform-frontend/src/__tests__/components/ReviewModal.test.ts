import { mount } from '@vue/test-utils'
import { beforeEach, describe, expect, it, vi } from 'vitest'

vi.mock('itx-ui-kit', () => ({
  Typography: { template: '<div><slot /></div>' },
}))

vi.mock('lucide-vue-next', () => ({
  X: { template: '<span />' },
}))

import ReviewModal from '@/components/ReviewModal.vue'

describe('ReviewModal', () => {
  it('does not render content when isOpen is false', () => {
    const wrapper = mount(ReviewModal, {
      props: { isOpen: false },
    })
    expect(wrapper.find('textarea').exists()).toBe(false)
  })

  it('renders content when isOpen is true', () => {
    const wrapper = mount(ReviewModal, {
      props: { isOpen: true },
    })
    expect(wrapper.find('textarea').exists()).toBe(true)
    expect(wrapper.text()).toContain('Добавить отзыв')
  })

  it('emits close when close button is clicked', async () => {
    const wrapper = mount(ReviewModal, {
      props: { isOpen: true },
    })
    // The close button has the X icon
    const closeButton = wrapper.find('button.absolute')
    await closeButton.trigger('click')
    expect(wrapper.emitted('close')).toHaveLength(1)
  })

  it('emits close when cancel button is clicked', async () => {
    const wrapper = mount(ReviewModal, {
      props: { isOpen: true },
    })
    const buttons = wrapper.findAll('button')
    const cancelButton = buttons.find(b => b.text().includes('Отменить'))
    expect(cancelButton).toBeDefined()
    await cancelButton!.trigger('click')
    expect(wrapper.emitted('close')).toHaveLength(1)
  })

  it('save button is disabled when textarea is empty', () => {
    const wrapper = mount(ReviewModal, {
      props: { isOpen: true },
    })
    const buttons = wrapper.findAll('button')
    const saveButton = buttons.find(b => b.text().includes('Сохранить'))
    expect(saveButton).toBeDefined()
    expect(saveButton!.attributes('disabled')).toBeDefined()
  })

  it('emits save with review text when save is clicked', async () => {
    const wrapper = mount(ReviewModal, {
      props: { isOpen: true },
    })
    const textarea = wrapper.find('textarea')
    await textarea.setValue('Great service!')

    const buttons = wrapper.findAll('button')
    const saveButton = buttons.find(b => b.text().includes('Сохранить'))
    await saveButton!.trigger('click')

    expect(wrapper.emitted('save')).toHaveLength(1)
    expect(wrapper.emitted('save')![0]).toEqual(['Great service!'])
  })

  it('clears textarea after save', async () => {
    const wrapper = mount(ReviewModal, {
      props: { isOpen: true },
    })
    const textarea = wrapper.find('textarea')
    await textarea.setValue('Test review')

    const buttons = wrapper.findAll('button')
    const saveButton = buttons.find(b => b.text().includes('Сохранить'))
    await saveButton!.trigger('click')

    expect((textarea.element as HTMLTextAreaElement).value).toBe('')
  })

  it('emits close on backdrop click', async () => {
    const wrapper = mount(ReviewModal, {
      props: { isOpen: true },
    })
    // Click on the backdrop (outermost div inside Transition)
    const backdrop = wrapper.find('.fixed')
    await backdrop.trigger('click')
    expect(wrapper.emitted('close')).toHaveLength(1)
  })
})
