import { mount } from '@vue/test-utils'
import { describe, expect, it } from 'vitest'
import { Input } from '../components/input'

describe('Input', () => {
  it('renders an input element', () => {
    const wrapper = mount(Input)
    expect(wrapper.element.tagName).toBe('INPUT')
  })

  it('uses terminal rounded-sm style', () => {
    const wrapper = mount(Input)
    const classes = wrapper.classes().join(' ')
    expect(classes).toContain('rounded-sm')
    expect(classes).not.toContain('rounded-xl')
  })

  it('binds v-model', async () => {
    const wrapper = mount(Input, {
      props: {
        modelValue: 'hello',
        'onUpdate:modelValue': (e: string | number) => wrapper.setProps({ modelValue: e }),
      },
    })
    expect((wrapper.element as HTMLInputElement).value).toBe('hello')

    await wrapper.setValue('world')
    expect(wrapper.emitted('update:modelValue')?.[0]).toEqual(['world'])
  })

  it('supports default value', () => {
    const wrapper = mount(Input, {
      props: { defaultValue: 'default' },
    })
    expect((wrapper.element as HTMLInputElement).value).toBe('default')
  })

  it('applies custom class', () => {
    const wrapper = mount(Input, {
      props: { class: 'w-64' },
    })
    expect(wrapper.classes()).toContain('w-64')
  })

  it('has proper focus ring for accessibility', () => {
    const wrapper = mount(Input)
    const classes = wrapper.classes().join(' ')
    expect(classes).toContain('focus-visible:ring-1')
  })
})
