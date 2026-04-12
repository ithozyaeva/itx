import { mount } from '@vue/test-utils'
import { describe, expect, it } from 'vitest'
import { Skeleton } from '../components/skeleton'

describe('Skeleton', () => {
  it('renders with animate-pulse', () => {
    const wrapper = mount(Skeleton)
    expect(wrapper.classes().join(' ')).toContain('animate-pulse')
  })

  it('applies muted background', () => {
    const wrapper = mount(Skeleton)
    expect(wrapper.classes().join(' ')).toContain('bg-muted')
  })

  it('merges custom class for sizing', () => {
    const wrapper = mount(Skeleton, {
      props: { class: 'h-4 w-32' },
    })
    expect(wrapper.classes()).toContain('h-4')
    expect(wrapper.classes()).toContain('w-32')
  })
})
