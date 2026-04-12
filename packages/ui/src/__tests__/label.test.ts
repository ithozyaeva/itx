import { mount } from '@vue/test-utils'
import { describe, expect, it } from 'vitest'
import { Label } from '../components/label'

describe('Label', () => {
  it('renders a label element', () => {
    const wrapper = mount(Label, {
      slots: { default: 'Name' },
    })
    expect(wrapper.element.tagName).toBe('LABEL')
    expect(wrapper.text()).toBe('Name')
  })

  it('passes for attribute', () => {
    const wrapper = mount(Label, {
      attrs: { for: 'email-input' },
      slots: { default: 'Email' },
    })
    expect(wrapper.attributes('for')).toBe('email-input')
  })
})
