import { mount } from '@vue/test-utils'
import { describe, expect, it } from 'vitest'
import { Tag } from '../components/tag'

describe('Tag', () => {
  it('renders with terminal rounded-sm style', () => {
    const wrapper = mount(Tag, {
      slots: { default: 'Vue' },
    })
    expect(wrapper.text()).toBe('Vue')
    expect(wrapper.classes().join(' ')).toContain('rounded-sm')
    expect(wrapper.classes().join(' ')).not.toContain('rounded-full')
  })

  it('applies active variant', () => {
    const wrapper = mount(Tag, {
      props: { variant: 'active' },
      slots: { default: 'Active' },
    })
    const classes = wrapper.classes().join(' ')
    expect(classes).toContain('text-accent')
  })

  it('applies default variant', () => {
    const wrapper = mount(Tag, {
      props: { variant: 'default' },
      slots: { default: 'Default' },
    })
    const classes = wrapper.classes().join(' ')
    expect(classes).toContain('text-foreground')
  })

  it('applies disabled state', () => {
    const wrapper = mount(Tag, {
      props: { disabled: true },
      slots: { default: 'Disabled' },
    })
    expect(wrapper.classes().join(' ')).toContain('opacity-50')
  })

  it('uses mono font for terminal style', () => {
    const wrapper = mount(Tag, {
      slots: { default: 'Tag' },
    })
    expect(wrapper.classes().join(' ')).toContain('font-mono')
  })
})
