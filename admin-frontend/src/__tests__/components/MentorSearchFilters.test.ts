import { describe, expect, it, vi } from 'vitest'

vi.mock('@/components/ui/button', () => ({
  Button: { template: '<button><slot /></button>', props: ['variant', 'size'] },
}))

vi.mock('@/components/ui/input', () => ({
  Input: { template: '<input />', props: ['id', 'modelValue', 'placeholder'] },
}))

vi.mock('@/components/ui/label', () => ({
  Label: { template: '<label><slot /></label>', props: ['for'] },
}))

vi.mock('./ui/card/Card.vue', () => ({
  default: { template: '<div><slot /></div>' },
}))

import { mount } from '@vue/test-utils'
import MentorSearchFilters from '@/components/MentorSearchFilters.vue'

describe('MentorSearchFilters', () => {
  it('renders the component', () => {
    const wrapper = mount(MentorSearchFilters)
    expect(wrapper.exists()).toBe(true)
  })

  it('renders labels for filter fields', () => {
    const wrapper = mount(MentorSearchFilters)
    expect(wrapper.text()).toContain('Имя / Username')
    expect(wrapper.text()).toContain('Профессиональный тег')
  })

  it('renders apply and reset buttons', () => {
    const wrapper = mount(MentorSearchFilters)
    expect(wrapper.text()).toContain('Применить')
    expect(wrapper.text()).toContain('Сбросить')
  })

  it('emits apply with default filters when apply is clicked', async () => {
    const wrapper = mount(MentorSearchFilters)
    const applyButton = wrapper.findAll('button').find(b => b.text() === 'Применить')
    await applyButton!.trigger('click')
    expect(wrapper.emitted('apply')).toHaveLength(1)
    expect(wrapper.emitted('apply')![0]).toEqual([{ name: '', tag: '' }])
  })

  it('emits apply with reset filters when reset is clicked', async () => {
    const wrapper = mount(MentorSearchFilters)
    const resetButton = wrapper.findAll('button').find(b => b.text() === 'Сбросить')
    await resetButton!.trigger('click')
    expect(wrapper.emitted('apply')).toHaveLength(1)
    expect(wrapper.emitted('apply')![0]).toEqual([{ name: '', tag: '' }])
  })
})
