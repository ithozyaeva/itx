import { mount } from '@vue/test-utils'
import { describe, expect, it } from 'vitest'

import FormField from '@/components/common/FormField.vue'

describe('FormField', () => {
  it('renders label text', () => {
    const wrapper = mount(FormField, {
      props: { label: 'Имя' },
    })
    expect(wrapper.find('label').text()).toContain('Имя')
  })

  it('renders slot content', () => {
    const wrapper = mount(FormField, {
      props: { label: 'Имя' },
      slots: { default: '<input type="text" class="test-input" />' },
    })
    expect(wrapper.find('.test-input').exists()).toBe(true)
  })

  it('shows error message when error prop provided', () => {
    const wrapper = mount(FormField, {
      props: { label: 'Имя', error: 'Обязательное поле' },
    })
    expect(wrapper.text()).toContain('Обязательное поле')
  })

  it('hides error message when error prop not provided', () => {
    const wrapper = mount(FormField, {
      props: { label: 'Имя' },
    })
    expect(wrapper.find('p').exists()).toBe(false)
  })

  it('error message has role="alert"', () => {
    const wrapper = mount(FormField, {
      props: { label: 'Имя', error: 'Ошибка' },
    })
    expect(wrapper.find('p[role="alert"]').exists()).toBe(true)
  })

  it('shows asterisk when required is true', () => {
    const wrapper = mount(FormField, {
      props: { label: 'Имя', required: true },
    })
    expect(wrapper.find('label').text()).toContain('*')
  })

  it('hides asterisk when required is false', () => {
    const wrapper = mount(FormField, {
      props: { label: 'Имя', required: false },
    })
    expect(wrapper.find('label').text()).not.toContain('*')
  })

  it('sets htmlFor on label when provided', () => {
    const wrapper = mount(FormField, {
      props: { label: 'Имя', htmlFor: 'name-input' },
    })
    expect(wrapper.find('label').attributes('for')).toBe('name-input')
  })
})
