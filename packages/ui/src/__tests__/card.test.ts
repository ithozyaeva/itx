import { mount } from '@vue/test-utils'
import { describe, expect, it } from 'vitest'
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from '../components/card'

describe('Card', () => {
  it('renders with terminal-card class', () => {
    const wrapper = mount(Card, {
      slots: { default: 'Content' },
    })
    const classes = wrapper.classes().join(' ')
    expect(classes).toContain('terminal-card')
    expect(classes).toContain('rounded-sm')
    expect(classes).not.toContain('rounded-xl')
    expect(classes).not.toContain('rounded-3xl')
  })

  it('renders slot content', () => {
    const wrapper = mount(Card, {
      slots: { default: '<p>Hello</p>' },
    })
    expect(wrapper.html()).toContain('<p>Hello</p>')
  })

  it('merges custom class', () => {
    const wrapper = mount(Card, {
      props: { class: 'extra' },
      slots: { default: 'x' },
    })
    expect(wrapper.classes()).toContain('extra')
    expect(wrapper.classes()).toContain('terminal-card')
  })
})

describe('CardHeader', () => {
  it('renders with correct spacing', () => {
    const wrapper = mount(CardHeader, {
      slots: { default: 'Header' },
    })
    expect(wrapper.classes().join(' ')).toContain('flex')
    expect(wrapper.text()).toBe('Header')
  })
})

describe('CardTitle', () => {
  it('renders as h3 by default', () => {
    const wrapper = mount(CardTitle, {
      slots: { default: 'Title' },
    })
    expect(wrapper.element.tagName).toBe('H3')
    expect(wrapper.text()).toBe('Title')
  })
})

describe('CardContent', () => {
  it('renders slot', () => {
    const wrapper = mount(CardContent, {
      slots: { default: 'Body' },
    })
    expect(wrapper.text()).toBe('Body')
  })
})

describe('CardDescription', () => {
  it('renders with muted styling', () => {
    const wrapper = mount(CardDescription, {
      slots: { default: 'Description' },
    })
    expect(wrapper.classes().join(' ')).toContain('text-muted-foreground')
  })
})

describe('CardFooter', () => {
  it('renders with flex layout', () => {
    const wrapper = mount(CardFooter, {
      slots: { default: 'Footer' },
    })
    expect(wrapper.classes().join(' ')).toContain('flex')
  })
})
