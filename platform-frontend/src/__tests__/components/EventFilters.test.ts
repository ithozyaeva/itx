import { mount } from '@vue/test-utils'
import { describe, expect, it, vi } from 'vitest'

vi.mock('@/composables/useDictionary', () => ({
  useDictionary: () => ({
    placeTypes: { value: [{ value: 'ONLINE', label: 'Онлайн' }, { value: 'OFFLINE', label: 'Офлайн' }] },
  }),
}))

vi.mock('@/components/ui/input', () => ({
  Input: {
    template: '<input :value="modelValue" @input="$emit(\'update:modelValue\', $event.target.value)" />',
    props: ['modelValue'],
    emits: ['update:modelValue'],
  },
}))

vi.mock('@/components/ui/select', () => ({
  Select: {
    template: '<div class="select-mock"><slot /></div>',
    props: ['modelValue'],
    emits: ['update:modelValue'],
  },
  SelectContent: { template: '<div><slot /></div>' },
  SelectItem: { template: '<div><slot /></div>', props: ['value'] },
  SelectTrigger: { template: '<div><slot /></div>' },
  SelectValue: { template: '<span />', props: ['placeholder'] },
}))

import EventFilters from '@/components/events/EventFilters.vue'

describe('EventFilters', () => {
  it('renders without errors', () => {
    const wrapper = mount(EventFilters)
    expect(wrapper.exists()).toBe(true)
  })

  it('renders input for title search', () => {
    const wrapper = mount(EventFilters)
    const input = wrapper.find('input')
    expect(input.exists()).toBe(true)
  })

  it('emits change event when title input changes', async () => {
    const wrapper = mount(EventFilters)
    const input = wrapper.find('input')
    await input.setValue('meetup')

    const emitted = wrapper.emitted('change')
    expect(emitted).toBeDefined()
    expect(emitted!.length).toBeGreaterThan(0)

    const lastEmit = emitted![emitted!.length - 1][0] as any
    expect(lastEmit.title).toBe('meetup')
  })

  it('renders select for place type', () => {
    const wrapper = mount(EventFilters)
    expect(wrapper.find('.select-mock').exists()).toBe(true)
  })
})
