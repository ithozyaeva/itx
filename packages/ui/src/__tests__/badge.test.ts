import { mount } from '@vue/test-utils'
import { describe, expect, it } from 'vitest'
import { Badge } from '../components/badge'

describe('Badge', () => {
  it('renders with default variant', () => {
    const wrapper = mount(Badge, {
      slots: { default: 'New' },
    })
    expect(wrapper.text()).toBe('New')
  })

  it('applies variant classes', () => {
    const wrapper = mount(Badge, {
      props: { variant: 'destructive' },
      slots: { default: 'Error' },
    })
    expect(wrapper.classes().join(' ')).toContain('destructive')
  })

  it('applies outline variant', () => {
    const wrapper = mount(Badge, {
      props: { variant: 'outline' },
      slots: { default: 'Tag' },
    })
    expect(wrapper.classes().join(' ')).toContain('text-foreground')
  })

  it('merges custom class', () => {
    const wrapper = mount(Badge, {
      props: { class: 'ml-2' },
      slots: { default: 'x' },
    })
    expect(wrapper.classes()).toContain('ml-2')
  })
})
