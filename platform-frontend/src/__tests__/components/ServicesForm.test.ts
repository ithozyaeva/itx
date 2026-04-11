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
  occupation: 'Dev',
  experience: '5 years',
  profTags: [],
  contacts: [],
  services: [
    { id: 1, name: 'Consulting', price: 100 },
    { id: 2, name: 'Code Review', price: 50 },
  ],
})

vi.mock('@/composables/useUser', () => ({
  useUser: () => mockUserData,
}))

vi.mock('@/services/profile', () => ({
  profileService: {
    updateServices: vi.fn().mockResolvedValue(true),
  },
}))

vi.mock('@/services/errorService', () => ({
  handleError: vi.fn(),
}))

vi.mock('@/components/ui/typography', () => ({
  Typography: { template: '<div><slot /></div>' },
}))

vi.mock('lucide-vue-next', () => ({
  Edit: { template: '<span />' },
  Loader2: { template: '<span />' },
  Plus: { template: '<span />' },
  Trash2: { template: '<span />' },
}))

vi.mock('@/components/ui/button', () => ({
  Button: {
    inheritAttrs: false,
    template: '<button :class="$attrs.class" :disabled="$attrs.disabled" @click="$emit(\'click\')"><slot /></button>',
  },
}))

vi.mock('@/components/ui/input', () => ({
  Input: {
    template: '<input :value="modelValue" :readonly="readonly" @input="$emit(\'update:modelValue\', $event.target.value)" />',
    props: ['modelValue', 'placeholder', 'readonly', 'type'],
    emits: ['update:modelValue'],
  },
}))

vi.mock('@/components/ui/textarea', () => ({
  Textarea: {
    template: '<textarea :value="modelValue" :readonly="readonly" @input="$emit(\'update:modelValue\', $event.target.value)" />',
    props: ['modelValue', 'placeholder', 'readonly'],
    emits: ['update:modelValue'],
  },
}))

vi.mock('@/components/ui/table', () => ({
  Table: { template: '<table><slot /></table>' },
  TableBody: { template: '<tbody><slot /></tbody>' },
  TableCell: { template: '<td><slot /></td>' },
  TableHead: { template: '<th><slot /></th>' },
  TableHeader: { template: '<thead><slot /></thead>' },
  TableRow: { template: '<tr><slot /></tr>' },
}))

import ServicesForm from '@/components/Profile/ServicesForm.vue'
import { profileService } from '@/services/profile'

describe('ServicesForm', () => {
  it('renders without errors', () => {
    const wrapper = mount(ServicesForm)
    expect(wrapper.exists()).toBe(true)
  })

  it('displays title', () => {
    const wrapper = mount(ServicesForm)
    expect(wrapper.text()).toContain('Услуги')
  })

  it('renders service rows from user data', () => {
    const wrapper = mount(ServicesForm)
    // 1 header row + 2 data rows = 3 total
    const rows = wrapper.findAll('tr')
    expect(rows.length).toBe(3)
  })

  it('does not show add/save buttons in view mode', () => {
    const wrapper = mount(ServicesForm)
    const addBtn = wrapper.findAll('button').find(b => b.text().includes('Добавить'))
    const saveBtn = wrapper.findAll('button').find(b => b.text().includes('Сохранить'))
    expect(addBtn).toBeUndefined()
    expect(saveBtn).toBeUndefined()
  })

  it('shows add and save buttons in edit mode', async () => {
    const wrapper = mount(ServicesForm)
    const editIcon = wrapper.find('.absolute')
    await editIcon.trigger('click')

    const addBtn = wrapper.findAll('button').find(b => b.text().includes('Добавить'))
    const saveBtn = wrapper.findAll('button').find(b => b.text().includes('Сохранить'))
    expect(addBtn).toBeDefined()
    expect(saveBtn).toBeDefined()
  })

  it('adds a new service row when add button is clicked', async () => {
    const wrapper = mount(ServicesForm)
    const editIcon = wrapper.find('.absolute')
    await editIcon.trigger('click')
    await wrapper.vm.$nextTick()

    const rowsBefore = wrapper.findAll('tbody tr').length
    const addBtn = wrapper.findAll('button').find(b => b.text().includes('Добавить'))
    await addBtn!.trigger('click')
    await wrapper.vm.$nextTick()

    const rowsAfter = wrapper.findAll('tbody tr').length
    expect(rowsAfter).toBeGreaterThan(rowsBefore)
  })

  it('removes a service when delete button is clicked in edit mode', async () => {
    const wrapper = mount(ServicesForm)
    const editIcon = wrapper.find('.absolute')
    await editIcon.trigger('click')
    await wrapper.vm.$nextTick()

    const rowsBefore = wrapper.findAll('tbody tr').length
    expect(rowsBefore).toBeGreaterThan(0)

    const deleteBtn = wrapper.find('button.text-destructive')
    expect(deleteBtn.exists()).toBe(true)
    await deleteBtn.trigger('click')
    await wrapper.vm.$nextTick()

    const rowsAfter = wrapper.findAll('tbody tr').length
    expect(rowsAfter).toBeLessThan(rowsBefore)
  })

  it('calls profileService.updateServices on submit', async () => {
    const wrapper = mount(ServicesForm)
    const editIcon = wrapper.find('.absolute')
    await editIcon.trigger('click')

    const saveBtn = wrapper.findAll('button').find(b => b.text().includes('Сохранить'))
    await saveBtn!.trigger('click')
    await vi.dynamicImportSettled()

    expect(profileService.updateServices).toHaveBeenCalled()
  })

  it('shows "Нет услуг" when no services and not in edit mode', () => {
    const originalServices = mockUserData.value.services
    mockUserData.value = { ...mockUserData.value, services: [] }

    const wrapper = mount(ServicesForm)
    expect(wrapper.text()).toContain('Нет услуг')

    mockUserData.value = { ...mockUserData.value, services: originalServices }
  })

  it('exits edit mode after successful submit', async () => {
    const wrapper = mount(ServicesForm)
    const editIcon = wrapper.find('.absolute')
    await editIcon.trigger('click')

    // Verify we are in edit mode
    let saveBtn = wrapper.findAll('button').find(b => b.text().includes('Сохранить'))
    expect(saveBtn).toBeDefined()

    await saveBtn!.trigger('click')
    await vi.dynamicImportSettled()

    // After submit, save button should no longer be visible
    saveBtn = wrapper.findAll('button').find(b => b.text().includes('Сохранить'))
    expect(saveBtn).toBeUndefined()
  })
})
