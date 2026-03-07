import { mount } from '@vue/test-utils'
import { describe, expect, it, vi } from 'vitest'

vi.mock('lucide-vue-next', () => ({
  Loader2: { template: '<span />' },
}))

vi.mock('@/components/ui/button', () => ({
  Button: {
    inheritAttrs: false,
    template: '<button :disabled="$attrs.disabled" @click="$emit(\'click\')"><slot /></button>',
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

vi.mock('@/components/ui/textarea', () => ({
  Textarea: {
    template: '<textarea :value="modelValue" @input="$emit(\'update:modelValue\', $event.target.value)" />',
    props: ['modelValue', 'placeholder', 'rows'],
    emits: ['update:modelValue'],
  },
}))

vi.mock('@/services/mentors', () => ({
  mentorsService: {
    addReview: vi.fn().mockResolvedValue(true),
  },
}))

import ReviewForm from '@/components/mentors/ReviewForm.vue'
import { mentorsService } from '@/services/mentors'

const services = [
  { id: 1, name: 'Consulting', price: 100 },
  { id: 2, name: 'Code Review', price: 50 },
]

describe('ReviewForm', () => {
  it('renders without errors', () => {
    const wrapper = mount(ReviewForm, {
      props: { mentorId: 1, services },
    })
    expect(wrapper.exists()).toBe(true)
  })

  it('displays title', () => {
    const wrapper = mount(ReviewForm, {
      props: { mentorId: 1, services },
    })
    expect(wrapper.text()).toContain('Оставить отзыв')
  })

  it('renders select and textarea', () => {
    const wrapper = mount(ReviewForm, {
      props: { mentorId: 1, services },
    })
    expect(wrapper.find('.select-mock').exists()).toBe(true)
    expect(wrapper.find('textarea').exists()).toBe(true)
  })

  it('submit button is disabled when no service or text', () => {
    const wrapper = mount(ReviewForm, {
      props: { mentorId: 1, services },
    })
    const button = wrapper.find('button')
    expect(button.attributes('disabled')).toBeDefined()
  })

  it('calls mentorsService.addReview on submit when data is filled', async () => {
    const wrapper = mount(ReviewForm, {
      props: { mentorId: 5, services },
    })

    // Set the internal state via component vm
    const vm = wrapper.vm as any
    vm.selectedServiceId = 1
    vm.text = 'Great mentor!'
    await wrapper.vm.$nextTick()

    const button = wrapper.find('button')
    await button.trigger('click')
    await vi.dynamicImportSettled()

    expect(mentorsService.addReview).toHaveBeenCalledWith(5, 1, 'Great mentor!')
  })

  it('emits submitted after successful review', async () => {
    const wrapper = mount(ReviewForm, {
      props: { mentorId: 1, services },
    })

    const vm = wrapper.vm as any
    vm.selectedServiceId = 2
    vm.text = 'Nice work'
    await wrapper.vm.$nextTick()

    const button = wrapper.find('button')
    await button.trigger('click')
    await vi.dynamicImportSettled()

    expect(wrapper.emitted('submitted')).toBeDefined()
  })

  it('resets form fields after successful submission', async () => {
    const wrapper = mount(ReviewForm, {
      props: { mentorId: 1, services },
    })

    const vm = wrapper.vm as any
    vm.selectedServiceId = 1
    vm.text = 'Feedback'
    await wrapper.vm.$nextTick()

    const button = wrapper.find('button')
    await button.trigger('click')
    await vi.dynamicImportSettled()

    expect(vm.selectedServiceId).toBe('')
    expect(vm.text).toBe('')
  })

  it('does not call addReview when selectedServiceId is empty', async () => {
    vi.mocked(mentorsService.addReview).mockClear()

    const wrapper = mount(ReviewForm, {
      props: { mentorId: 1, services },
    })

    const vm = wrapper.vm as any
    vm.selectedServiceId = ''
    vm.text = 'Some text'
    await wrapper.vm.$nextTick()

    const button = wrapper.find('button')
    await button.trigger('click')
    await vi.dynamicImportSettled()

    expect(mentorsService.addReview).not.toHaveBeenCalled()
  })

  it('does not call addReview when text is empty', async () => {
    vi.mocked(mentorsService.addReview).mockClear()

    const wrapper = mount(ReviewForm, {
      props: { mentorId: 1, services },
    })

    const vm = wrapper.vm as any
    vm.selectedServiceId = 1
    vm.text = '   '
    await wrapper.vm.$nextTick()

    const button = wrapper.find('button')
    await button.trigger('click')
    await vi.dynamicImportSettled()

    expect(mentorsService.addReview).not.toHaveBeenCalled()
  })
})
