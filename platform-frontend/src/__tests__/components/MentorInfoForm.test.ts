import { mount } from '@vue/test-utils'
import { describe, expect, it, vi } from 'vitest'
import { ref } from 'vue'

const mockUserData = ref({
  id: 1,
  telegramID: 123,
  tg: 'testuser',
  birthday: '',
  firstName: 'John',
  lastName: 'Doe',
  bio: '',
  grade: '',
  company: '',
  avatarUrl: '',
  roles: ['MENTOR'],
  occupation: 'Backend Developer',
  experience: '5 years of Go',
  profTags: [],
  contacts: [],
  services: [],
})

vi.mock('@/composables/useUser', () => ({
  useUser: () => mockUserData,
}))

vi.mock('@/services/profile', () => ({
  profileService: {
    updateMentorInfo: vi.fn().mockResolvedValue(true),
  },
}))

vi.mock('lucide-vue-next', () => ({
  Edit: { template: '<span />' },
}))

vi.mock('@/components/ui/button', () => ({
  Button: {
    inheritAttrs: false,
    template: '<button :disabled="$attrs.disabled" @click="$emit(\'click\')"><slot /></button>',
  },
}))

vi.mock('@/components/ui/input', () => ({
  Input: {
    template: '<input :value="modelValue" :readonly="readonly" @input="$emit(\'update:modelValue\', $event.target.value)" />',
    props: ['modelValue', 'placeholder', 'readonly', 'id', 'maxLength'],
    emits: ['update:modelValue'],
  },
}))

vi.mock('@/components/ui/label', () => ({
  Label: { template: '<label><slot /></label>', props: ['for'] },
}))

vi.mock('@/components/ui/textarea', () => ({
  Textarea: {
    template: '<textarea :value="modelValue" :readonly="readonly" @input="$emit(\'update:modelValue\', $event.target.value)" />',
    props: ['modelValue', 'placeholder', 'readonly', 'rows', 'id', 'maxLength'],
    emits: ['update:modelValue'],
  },
}))

import MentorInfoForm from '@/components/Profile/MentorInfoForm.vue'
import { profileService } from '@/services/profile'

describe('MentorInfoForm', () => {
  it('renders without errors', () => {
    const wrapper = mount(MentorInfoForm)
    expect(wrapper.exists()).toBe(true)
  })

  it('displays occupation and experience labels', () => {
    const wrapper = mount(MentorInfoForm)
    expect(wrapper.text()).toContain('Специализация')
    expect(wrapper.text()).toContain('Опыт работы')
  })

  it('shows user occupation in input', () => {
    const wrapper = mount(MentorInfoForm)
    const input = wrapper.find('input')
    expect(input.element.value).toBe('Backend Developer')
  })

  it('shows user experience in textarea', () => {
    const wrapper = mount(MentorInfoForm)
    const textarea = wrapper.find('textarea')
    expect(textarea.element.value).toBe('5 years of Go')
  })

  it('does not show save button in view mode', () => {
    const wrapper = mount(MentorInfoForm)
    const saveBtn = wrapper.findAll('button').find(b => b.text().includes('Сохранить'))
    expect(saveBtn).toBeUndefined()
  })

  it('shows save button in edit mode', async () => {
    const wrapper = mount(MentorInfoForm)
    const editIcon = wrapper.find('.absolute')
    await editIcon.trigger('click')

    const saveBtn = wrapper.findAll('button').find(b => b.text().includes('Сохранить'))
    expect(saveBtn).toBeDefined()
  })

  it('calls profileService.updateMentorInfo on submit', async () => {
    const wrapper = mount(MentorInfoForm)
    const editIcon = wrapper.find('.absolute')
    await editIcon.trigger('click')

    const saveBtn = wrapper.findAll('button').find(b => b.text().includes('Сохранить'))
    await saveBtn!.trigger('click')
    await vi.dynamicImportSettled()

    expect(profileService.updateMentorInfo).toHaveBeenCalledWith({
      occupation: 'Backend Developer',
      experience: '5 years of Go',
    })
  })

  it('exits edit mode after submit', async () => {
    const wrapper = mount(MentorInfoForm)
    const editIcon = wrapper.find('.absolute')
    await editIcon.trigger('click')

    let saveBtn = wrapper.findAll('button').find(b => b.text().includes('Сохранить'))
    expect(saveBtn).toBeDefined()

    await saveBtn!.trigger('click')
    await vi.dynamicImportSettled()

    saveBtn = wrapper.findAll('button').find(b => b.text().includes('Сохранить'))
    expect(saveBtn).toBeUndefined()
  })

  it('toggles edit mode on edit icon click', async () => {
    const wrapper = mount(MentorInfoForm)
    const editIcon = wrapper.find('.absolute')

    await editIcon.trigger('click')
    let saveBtn = wrapper.findAll('button').find(b => b.text().includes('Сохранить'))
    expect(saveBtn).toBeDefined()

    await editIcon.trigger('click')
    saveBtn = wrapper.findAll('button').find(b => b.text().includes('Сохранить'))
    expect(saveBtn).toBeUndefined()
  })

  it('inputs are readonly when not in edit mode', () => {
    const wrapper = mount(MentorInfoForm)
    const input = wrapper.find('input')
    expect(input.attributes('readonly')).toBeDefined()
  })
})
