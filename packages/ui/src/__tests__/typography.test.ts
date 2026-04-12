import { mount } from '@vue/test-utils'
import { describe, expect, it } from 'vitest'
import { Typography } from '../components/typography'

describe('Typography', () => {
  it('renders h1 variant', () => {
    const wrapper = mount(Typography, {
      props: { variant: 'h1', as: 'h1' },
      slots: { default: 'Heading' },
    })
    expect(wrapper.element.tagName).toBe('H1')
    expect(wrapper.text()).toBe('Heading')
  })

  it('renders h2 variant', () => {
    const wrapper = mount(Typography, {
      props: { variant: 'h2', as: 'h2' },
      slots: { default: 'Sub' },
    })
    expect(wrapper.element.tagName).toBe('H2')
  })

  it('renders p variant', () => {
    const wrapper = mount(Typography, {
      props: { variant: 'p', as: 'p' },
      slots: { default: 'Text' },
    })
    expect(wrapper.element.tagName).toBe('P')
  })

  it('renders as custom tag', () => {
    const wrapper = mount(Typography, {
      props: { variant: 'h3', as: 'span' },
      slots: { default: 'Custom' },
    })
    expect(wrapper.element.tagName).toBe('SPAN')
  })

  it('merges custom class', () => {
    const wrapper = mount(Typography, {
      props: { variant: 'h1', class: 'text-accent' },
      slots: { default: 'Styled' },
    })
    expect(wrapper.classes()).toContain('text-accent')
  })
})
