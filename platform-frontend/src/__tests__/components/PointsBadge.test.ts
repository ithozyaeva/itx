import { mount } from '@vue/test-utils'
import { describe, expect, it } from 'vitest'
import PointsBadge from '@/components/progress/PointsBadge.vue'

describe('PointsBadge', () => {
  it('shows "+N" when not earned', () => {
    const w = mount(PointsBadge, { props: { amount: 25 } })
    expect(w.text()).toBe('+25')
  })

  it('shows just N when earned', () => {
    const w = mount(PointsBadge, { props: { amount: 25, earned: true } })
    expect(w.text()).toBe('25')
  })

  it('uses yellow when not earned, green when earned', () => {
    const notEarned = mount(PointsBadge, { props: { amount: 1 } })
    expect(notEarned.classes().some(c => c.includes('yellow'))).toBe(true)
    const earned = mount(PointsBadge, { props: { amount: 1, earned: true } })
    expect(earned.classes().some(c => c.includes('green'))).toBe(true)
  })
})
