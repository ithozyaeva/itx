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

vi.mock('@/components/ui/select', () => ({
  Select: { template: '<div><slot /></div>', props: ['id', 'modelValue', 'multiple'] },
  SelectContent: { template: '<div><slot /></div>' },
  SelectItem: { template: '<div><slot /></div>', props: ['value'] },
  SelectTrigger: { template: '<div><slot /></div>' },
  SelectValue: { template: '<span />', props: ['placeholder'] },
}))

vi.mock('./ui/card/Card.vue', () => ({
  default: { template: '<div><slot /></div>' },
}))

vi.mock('@/composables/useDictionary', () => ({
  useDictionary: () => ({
    memberRoles: { value: [
      { value: 'ADMIN', label: 'Админ' },
      { value: 'SUBSCRIBER', label: 'Подписчик' },
    ] },
  }),
}))

import { mount } from '@vue/test-utils'
import MemberSearchFilters from '@/components/MemberSearchFilters.vue'

describe('MemberSearchFilters', () => {
  it('renders the component', () => {
    const wrapper = mount(MemberSearchFilters)
    expect(wrapper.exists()).toBe(true)
  })

  it('renders apply and reset buttons', () => {
    const wrapper = mount(MemberSearchFilters)
    expect(wrapper.text()).toContain('Применить')
    expect(wrapper.text()).toContain('Сбросить')
  })

  it('renders labels for filter fields', () => {
    const wrapper = mount(MemberSearchFilters)
    expect(wrapper.text()).toContain('TG username')
    expect(wrapper.text()).toContain('Роли')
  })

  it('emits apply with filters when apply is clicked', async () => {
    const wrapper = mount(MemberSearchFilters)
    const applyButton = wrapper.findAll('button').find(b => b.text() === 'Применить')
    await applyButton!.trigger('click')
    expect(wrapper.emitted('apply')).toHaveLength(1)
    expect(wrapper.emitted('apply')![0]).toEqual([{ username: '', roles: [] }])
  })

  it('emits apply with reset filters when reset is clicked', async () => {
    const wrapper = mount(MemberSearchFilters)
    const resetButton = wrapper.findAll('button').find(b => b.text() === 'Сбросить')
    await resetButton!.trigger('click')
    expect(wrapper.emitted('apply')).toHaveLength(1)
    expect(wrapper.emitted('apply')![0]).toEqual([{ username: '', roles: [] }])
  })
})
