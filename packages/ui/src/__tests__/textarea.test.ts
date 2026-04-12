import { mount } from '@vue/test-utils'
import { describe, expect, it } from 'vitest'
import { Textarea } from '../components/textarea'

describe('Textarea', () => {
  it('renders a textarea element', () => {
    const wrapper = mount(Textarea)
    expect(wrapper.element.tagName).toBe('TEXTAREA')
  })

  it('uses terminal rounded-sm style', () => {
    const wrapper = mount(Textarea)
    const classes = wrapper.classes().join(' ')
    expect(classes).toContain('rounded-sm')
    expect(classes).not.toContain('rounded-xl')
  })

  it('binds v-model', async () => {
    const wrapper = mount(Textarea, {
      props: {
        modelValue: 'text',
        'onUpdate:modelValue': (e: string | number) => wrapper.setProps({ modelValue: e }),
      },
    })
    expect((wrapper.element as HTMLTextAreaElement).value).toBe('text')

    await wrapper.setValue('new text')
    expect(wrapper.emitted('update:modelValue')?.[0]).toEqual(['new text'])
  })
})
