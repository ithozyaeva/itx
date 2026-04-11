import { mount } from '@vue/test-utils'
import { describe, expect, it } from 'vitest'
import Typography from '@/components/ui/typography/Typography.vue'

describe('Typography', () => {
  it('renders as <span> with body-m classes by default', () => {
    const wrapper = mount(Typography)
    expect(wrapper.element.tagName).toBe('SPAN')
    expect(wrapper.classes()).toContain('text-base')
  })

  it('renders correct HTML tag when as prop is set', () => {
    const wrapper = mount(Typography, { props: { as: 'h1' } })
    expect(wrapper.element.tagName).toBe('H1')
  })

  it('applies correct classes for h1 variant', () => {
    const wrapper = mount(Typography, { props: { variant: 'h1' } })
    expect(wrapper.classes()).toContain('text-3xl')
    expect(wrapper.classes()).toContain('font-bold')
    expect(wrapper.classes()).toContain('tracking-tight')
    expect(wrapper.classes()).toContain('uppercase')
  })

  it('applies correct classes for h2 variant', () => {
    const wrapper = mount(Typography, { props: { variant: 'h2' } })
    expect(wrapper.classes()).toContain('text-2xl')
    expect(wrapper.classes()).toContain('font-bold')
    expect(wrapper.classes()).toContain('uppercase')
  })

  it('applies correct classes for h3 variant', () => {
    const wrapper = mount(Typography, { props: { variant: 'h3' } })
    expect(wrapper.classes()).toContain('text-lg')
    expect(wrapper.classes()).toContain('font-semibold')
  })

  it('applies correct classes for h4 variant', () => {
    const wrapper = mount(Typography, { props: { variant: 'h4' } })
    expect(wrapper.classes()).toContain('text-base')
    expect(wrapper.classes()).toContain('font-semibold')
  })

  it('applies correct classes for body-m variant', () => {
    const wrapper = mount(Typography, { props: { variant: 'body-m' } })
    expect(wrapper.classes()).toContain('text-base')
  })

  it('applies correct classes for date variant', () => {
    const wrapper = mount(Typography, { props: { variant: 'date' } })
    expect(wrapper.classes()).toContain('text-sm')
    expect(wrapper.classes()).toContain('text-muted-foreground')
  })

  it('renders slot content', () => {
    const wrapper = mount(Typography, {
      slots: { default: 'Hello World' },
    })
    expect(wrapper.text()).toBe('Hello World')
  })

  it('passes through additional class prop', () => {
    const wrapper = mount(Typography, {
      props: { class: 'mt-4' },
    })
    expect(wrapper.classes()).toContain('mt-4')
  })

  it('combines variant classes with custom class', () => {
    const wrapper = mount(Typography, {
      props: { variant: 'h1', class: 'text-red-500' },
    })
    expect(wrapper.classes()).toContain('text-3xl')
    expect(wrapper.classes()).toContain('font-bold')
    expect(wrapper.classes()).toContain('text-red-500')
  })
})
