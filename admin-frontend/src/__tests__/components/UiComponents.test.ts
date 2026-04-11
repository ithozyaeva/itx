import { describe, expect, it } from 'vitest'
import { mount } from '@vue/test-utils'
import Typography from '@/components/ui/typography/Typography.vue'
import Tag from '@/components/ui/tag/Tag.vue'

describe('Typography', () => {
  it('renders as <span> with body-m classes by default', () => {
    const wrapper = mount(Typography, {
      slots: { default: 'Hello' },
    })
    expect(wrapper.element.tagName).toBe('SPAN')
    expect(wrapper.classes()).toContain('text-base')
  })

  it('renders correct tag with `as` prop', () => {
    const wrapper = mount(Typography, {
      props: { as: 'h1' },
      slots: { default: 'Title' },
    })
    expect(wrapper.element.tagName).toBe('H1')
  })

  it('applies variant classes for h2', () => {
    const wrapper = mount(Typography, {
      props: { variant: 'h2' },
      slots: { default: 'Heading' },
    })
    expect(wrapper.classes()).toContain('text-2xl')
    expect(wrapper.classes()).toContain('font-bold')
    expect(wrapper.classes()).toContain('uppercase')
  })

  it('applies variant classes for h3', () => {
    const wrapper = mount(Typography, {
      props: { variant: 'h3' },
      slots: { default: 'Heading' },
    })
    expect(wrapper.classes()).toContain('text-lg')
    expect(wrapper.classes()).toContain('font-semibold')
  })

  it('applies variant classes for h4', () => {
    const wrapper = mount(Typography, {
      props: { variant: 'h4' },
      slots: { default: 'Heading' },
    })
    expect(wrapper.classes()).toContain('text-base')
    expect(wrapper.classes()).toContain('font-semibold')
  })

  it('renders slot content', () => {
    const wrapper = mount(Typography, {
      slots: { default: 'Slot content here' },
    })
    expect(wrapper.text()).toBe('Slot content here')
  })

  it('passes through class prop', () => {
    const wrapper = mount(Typography, {
      props: { class: 'custom-class' },
      slots: { default: 'Text' },
    })
    expect(wrapper.classes()).toContain('custom-class')
    // Should also retain variant classes
    expect(wrapper.classes()).toContain('text-base')
  })
})

describe('Tag', () => {
  it('renders with default variant styling', () => {
    const wrapper = mount(Tag, {
      slots: { default: 'Tag text' },
    })
    expect(wrapper.classes()).toContain('border-border')
    expect(wrapper.classes()).toContain('text-foreground')
    expect(wrapper.classes()).not.toContain('border-accent/30')
  })

  it('renders with active variant styling', () => {
    const wrapper = mount(Tag, {
      props: { variant: 'active' },
      slots: { default: 'Active tag' },
    })
    expect(wrapper.classes()).toContain('border-accent/30')
    expect(wrapper.classes()).toContain('text-accent')
    expect(wrapper.classes()).not.toContain('border-border')
  })

  it('renders disabled state with opacity-50', () => {
    const wrapper = mount(Tag, {
      props: { disabled: true },
      slots: { default: 'Disabled' },
    })
    expect(wrapper.classes()).toContain('opacity-50')
    expect(wrapper.classes()).toContain('pointer-events-none')
  })

  it('renders slot content', () => {
    const wrapper = mount(Tag, {
      slots: { default: 'My tag' },
    })
    expect(wrapper.text()).toBe('My tag')
  })

  it('renders as custom element with `as` prop', () => {
    const wrapper = mount(Tag, {
      props: { as: 'button' },
      slots: { default: 'Button tag' },
    })
    expect(wrapper.element.tagName).toBe('BUTTON')
  })
})
