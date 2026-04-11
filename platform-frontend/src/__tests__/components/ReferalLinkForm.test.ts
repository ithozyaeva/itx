import { mount } from '@vue/test-utils'
import { describe, expect, it, vi } from 'vitest'

vi.mock('@/components/ui/typography', () => ({
  Typography: { template: '<div><slot /></div>' },
}))

vi.mock('lucide-vue-next', () => ({
  Loader2: { template: '<span />' },
}))

vi.mock('@/composables/useDictionary', () => ({
  useDictionary: () => ({
    grades: { value: [{ value: 'junior', label: 'Junior' }, { value: 'middle', label: 'Middle' }, { value: 'senior', label: 'Senior' }] },
  }),
}))

vi.mock('@/components/common/ProfTagsInput.vue', () => ({
  default: {
    template: '<div class="prof-tags-input" />',
    props: ['modelValue'],
    emits: ['update:modelValue'],
  },
}))

vi.mock('@/components/ui/button', () => ({
  Button: {
    inheritAttrs: false,
    template: '<button :disabled="$attrs.disabled" @click="$emit(\'click\')"><slot /></button>',
  },
}))

vi.mock('@/components/ui/input', () => ({
  Input: {
    template: '<input :value="modelValue" @input="$emit(\'update:modelValue\', $event.target.value)" />',
    props: ['modelValue', 'type', 'min', 'id'],
    emits: ['update:modelValue'],
  },
}))

vi.mock('@/components/ui/label', () => ({
  Label: { template: '<label><slot /></label>', props: ['for'] },
}))

vi.mock('@/components/ui/select', () => ({
  Select: {
    template: '<div class="select-mock"><slot /></div>',
    props: ['modelValue', 'id'],
    emits: ['update:modelValue'],
  },
  SelectContent: { template: '<div><slot /></div>' },
  SelectItem: { template: '<div><slot /></div>', props: ['value'] },
  SelectTrigger: { template: '<div><slot /></div>' },
  SelectValue: { template: '<span />', props: ['placeholder'] },
}))

import ReferalLinkForm from '@/components/referals/ReferalLinkForm.vue'

describe('ReferalLinkForm', () => {
  it('renders without errors', () => {
    const wrapper = mount(ReferalLinkForm)
    expect(wrapper.exists()).toBe(true)
  })

  it('displays custom title', () => {
    const wrapper = mount(ReferalLinkForm, {
      props: { title: 'Новая ссылка' },
    })
    expect(wrapper.text()).toContain('Новая ссылка')
  })

  it('renders company, grade, tags, vacations and date inputs', () => {
    const wrapper = mount(ReferalLinkForm)
    const inputs = wrapper.findAll('input')
    // company, vacationsCount, expiresAt = 3 inputs
    expect(inputs.length).toBe(3)
    expect(wrapper.find('.select-mock').exists()).toBe(true)
    expect(wrapper.find('.prof-tags-input').exists()).toBe(true)
  })

  it('initializes form with link data', () => {
    const link = {
      company: 'Test Corp',
      grade: 'senior' as const,
      profTags: [{ id: 1, title: 'Go' }],
      vacationsCount: 5,
      expiresAt: '2026-06-01T00:00:00.000Z',
    }
    const wrapper = mount(ReferalLinkForm, {
      props: { link },
    })
    const vm = wrapper.vm as any
    expect(vm.formData.company).toBe('Test Corp')
    expect(vm.formData.grade).toBe('senior')
    expect(vm.formData.vacationsCount).toBe(5)
  })

  it('defaults grade to junior when no link provided', () => {
    const wrapper = mount(ReferalLinkForm)
    const vm = wrapper.vm as any
    expect(vm.formData.grade).toBe('junior')
  })

  it('emits save with form data on save button click', async () => {
    const wrapper = mount(ReferalLinkForm)
    const vm = wrapper.vm as any
    vm.formData.company = 'NewCorp'
    await wrapper.vm.$nextTick()

    const buttons = wrapper.findAll('button')
    const saveBtn = buttons.find(b => b.text().includes('Сохранить'))
    await saveBtn!.trigger('click')

    const emitted = wrapper.emitted('save')
    expect(emitted).toBeDefined()
    expect((emitted![0][0] as any).company).toBe('NewCorp')
  })

  it('emits cancel on cancel button click', async () => {
    const wrapper = mount(ReferalLinkForm)

    const buttons = wrapper.findAll('button')
    const cancelBtn = buttons.find(b => b.text().includes('Отменить'))
    await cancelBtn!.trigger('click')

    expect(wrapper.emitted('cancel')).toBeDefined()
  })

  it('includes expiresAt as ISO string when date is set', async () => {
    const wrapper = mount(ReferalLinkForm)
    const vm = wrapper.vm as any
    vm.expiresAtDate = '2026-12-31'
    await wrapper.vm.$nextTick()

    const saveBtn = wrapper.findAll('button').find(b => b.text().includes('Сохранить'))
    await saveBtn!.trigger('click')

    const emitted = wrapper.emitted('save')
    expect(emitted).toBeDefined()
    const data = emitted![0][0] as any
    expect(data.expiresAt).toContain('2026-12-31')
  })

  it('does not include expiresAt when date is empty', async () => {
    const wrapper = mount(ReferalLinkForm)

    const saveBtn = wrapper.findAll('button').find(b => b.text().includes('Сохранить'))
    await saveBtn!.trigger('click')

    const emitted = wrapper.emitted('save')
    const data = emitted![0][0] as any
    expect(data.expiresAt).toBeUndefined()
  })

  it('disables save button when isSaving is true', () => {
    const wrapper = mount(ReferalLinkForm, {
      props: { isSaving: true },
    })
    const saveBtn = wrapper.findAll('button').find(b => b.text().includes('Сохранить'))
    expect(saveBtn!.attributes('disabled')).toBeDefined()
  })
})
