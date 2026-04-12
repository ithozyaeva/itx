import { mount } from '@vue/test-utils'
import { describe, expect, it } from 'vitest'
import { Button, buttonVariants } from '../components/button'

describe('Button', () => {
  it('renders with default props', () => {
    const wrapper = mount(Button, {
      slots: { default: 'Click me' },
    })
    expect(wrapper.text()).toBe('Click me')
    expect(wrapper.element.tagName).toBe('BUTTON')
  })

  it('applies default variant classes', () => {
    const wrapper = mount(Button)
    const classes = wrapper.classes().join(' ')
    expect(classes).toContain('rounded-sm')
    expect(classes).toContain('bg-accent')
  })

  it('applies destructive variant', () => {
    const wrapper = mount(Button, {
      props: { variant: 'destructive' },
    })
    const classes = wrapper.classes().join(' ')
    expect(classes).toContain('bg-destructive')
  })

  it('applies ghost variant', () => {
    const wrapper = mount(Button, {
      props: { variant: 'ghost' },
    })
    const classes = wrapper.classes().join(' ')
    expect(classes).not.toContain('bg-accent')
  })

  it('applies outline variant', () => {
    const wrapper = mount(Button, {
      props: { variant: 'outline' },
    })
    const classes = wrapper.classes().join(' ')
    expect(classes).toContain('border')
  })

  it('applies size variants', () => {
    const sm = mount(Button, { props: { size: 'sm' } })
    expect(sm.classes().join(' ')).toContain('h-8')

    const lg = mount(Button, { props: { size: 'lg' } })
    expect(lg.classes().join(' ')).toContain('h-10')

    const icon = mount(Button, { props: { size: 'icon' } })
    expect(icon.classes().join(' ')).toContain('w-9')
  })

  it('merges custom class', () => {
    const wrapper = mount(Button, {
      props: { class: 'custom-class' },
    })
    expect(wrapper.classes()).toContain('custom-class')
  })

  it('emits click events', async () => {
    const wrapper = mount(Button)
    await wrapper.trigger('click')
    expect(wrapper.emitted('click')).toHaveLength(1)
  })

  it('uses rounded-sm for terminal style (not rounded-full)', () => {
    const classes = buttonVariants({ variant: 'default', size: 'default' })
    expect(classes).toContain('rounded-sm')
    expect(classes).not.toContain('rounded-full')
    expect(classes).not.toContain('rounded-xl')
  })

  it('has cursor-pointer for interactivity', () => {
    const classes = buttonVariants({ variant: 'default', size: 'default' })
    expect(classes).toContain('cursor-pointer')
  })
})
