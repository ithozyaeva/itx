import { mount } from '@vue/test-utils'
import { describe, expect, it } from 'vitest'
import ProgressBar from '@/components/progress/ProgressBar.vue'

describe('ProgressBar', () => {
  it('exposes ARIA attributes', () => {
    const w = mount(ProgressBar, { props: { progress: 3, target: 10 } })
    const el = w.find('[role="progressbar"]')
    expect(el.exists()).toBe(true)
    expect(el.attributes('aria-valuenow')).toBe('3')
    expect(el.attributes('aria-valuemax')).toBe('10')
    expect(el.attributes('aria-valuemin')).toBe('0')
  })

  function fill(w: ReturnType<typeof mount>) {
    return w.find('[role="progressbar"] > div')
  }

  it('renders fill width as percentage of target', () => {
    const w = mount(ProgressBar, { props: { progress: 5, target: 10 } })
    expect(fill(w).attributes('style')).toContain('width: 50%')
  })

  it('caps fill at 100%', () => {
    const w = mount(ProgressBar, { props: { progress: 50, target: 10 } })
    expect(fill(w).attributes('style')).toContain('width: 100%')
  })

  it('uses green fill when state=done', () => {
    const w = mount(ProgressBar, { props: { progress: 5, target: 5, state: 'done' } })
    expect(fill(w).classes().some(c => c.includes('green'))).toBe(true)
  })

  it('uses yellow fill when state=awarded', () => {
    const w = mount(ProgressBar, { props: { progress: 5, target: 5, state: 'awarded' } })
    expect(fill(w).classes().some(c => c.includes('yellow'))).toBe(true)
  })

  it('falls back to label "Прогресс N из M" by default', () => {
    const w = mount(ProgressBar, { props: { progress: 2, target: 8 } })
    expect(w.find('[role="progressbar"]').attributes('aria-label')).toBe('Прогресс 2 из 8')
  })

  it('honors custom label prop', () => {
    const w = mount(ProgressBar, { props: { progress: 2, target: 8, label: 'My label' } })
    expect(w.find('[role="progressbar"]').attributes('aria-label')).toBe('My label')
  })
})
