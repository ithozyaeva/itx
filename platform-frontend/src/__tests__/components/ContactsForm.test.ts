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
  contacts: [
    { id: 1, type: 1, link: 'https://t.me/johndoe' },
    { id: 2, type: 2, link: 'john@example.com' },
  ],
  services: [],
})

vi.mock('@/composables/useUser', () => ({
  useUser: () => mockUserData,
}))

vi.mock('@/services/profile', () => ({
  profileService: {
    updateContacts: vi.fn().mockResolvedValue(true),
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
    props: ['modelValue', 'placeholder', 'readonly'],
    emits: ['update:modelValue'],
  },
}))

vi.mock('@/components/ui/select', () => ({
  Select: { template: '<div><slot /></div>', props: ['modelValue', 'disabled'] },
  SelectContent: { template: '<div><slot /></div>' },
  SelectItem: { template: '<div><slot /></div>', props: ['value'] },
  SelectTrigger: { template: '<div><slot /></div>' },
  SelectValue: { template: '<span />', props: ['placeholder'] },
}))

vi.mock('@/components/ui/table', () => ({
  Table: { template: '<table><slot /></table>' },
  TableBody: { template: '<tbody><slot /></tbody>' },
  TableCell: { template: '<td><slot /></td>' },
  TableHead: { template: '<th><slot /></th>' },
  TableHeader: { template: '<thead><slot /></thead>' },
  TableRow: { template: '<tr><slot /></tr>' },
}))

import ContactsForm from '@/components/Profile/ContactsForm.vue'
import { profileService } from '@/services/profile'

describe('ContactsForm', () => {
  it('renders without errors', () => {
    const wrapper = mount(ContactsForm)
    expect(wrapper.exists()).toBe(true)
  })

  it('displays title', () => {
    const wrapper = mount(ContactsForm)
    expect(wrapper.text()).toContain('Контакты')
  })

  it('renders contact rows from user data', () => {
    const wrapper = mount(ContactsForm)
    // 1 header row + 2 data rows = 3 total
    const rows = wrapper.findAll('tr')
    expect(rows.length).toBe(3)
  })

  it('does not show add/save buttons in view mode', () => {
    const wrapper = mount(ContactsForm)
    const addBtn = wrapper.findAll('button').find(b => b.text().includes('Добавить'))
    const saveBtn = wrapper.findAll('button').find(b => b.text().includes('Сохранить'))
    expect(addBtn).toBeUndefined()
    expect(saveBtn).toBeUndefined()
  })

  it('shows edit, add and save buttons in edit mode', async () => {
    const wrapper = mount(ContactsForm)
    const editIcon = wrapper.find('.absolute')
    await editIcon.trigger('click')

    const addBtn = wrapper.findAll('button').find(b => b.text().includes('Добавить'))
    const saveBtn = wrapper.findAll('button').find(b => b.text().includes('Сохранить'))
    expect(addBtn).toBeDefined()
    expect(saveBtn).toBeDefined()
  })

  it('adds a new contact row when add button is clicked', async () => {
    const wrapper = mount(ContactsForm)
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

  it('removes a contact when delete button is clicked in edit mode', async () => {
    const wrapper = mount(ContactsForm)
    const editIcon = wrapper.find('.absolute')
    await editIcon.trigger('click')
    await wrapper.vm.$nextTick()

    const rowsBefore = wrapper.findAll('tbody tr').length
    expect(rowsBefore).toBeGreaterThan(0)

    // Find a delete button (button with class text-destructive)
    const deleteBtn = wrapper.find('button.text-destructive')
    expect(deleteBtn.exists()).toBe(true)
    await deleteBtn.trigger('click')
    await wrapper.vm.$nextTick()

    const rowsAfter = wrapper.findAll('tbody tr').length
    expect(rowsAfter).toBeLessThan(rowsBefore)
  })

  it('calls profileService.updateContacts on submit', async () => {
    const wrapper = mount(ContactsForm)
    const editIcon = wrapper.find('.absolute')
    await editIcon.trigger('click')

    const saveBtn = wrapper.findAll('button').find(b => b.text().includes('Сохранить'))
    await saveBtn!.trigger('click')
    await vi.dynamicImportSettled()

    expect(profileService.updateContacts).toHaveBeenCalled()
  })

  it('shows "Нет контактов" when no contacts and not in edit mode', () => {
    const originalContacts = mockUserData.value.contacts
    mockUserData.value = { ...mockUserData.value, contacts: [] }

    const wrapper = mount(ContactsForm)
    expect(wrapper.text()).toContain('Нет контактов')

    mockUserData.value = { ...mockUserData.value, contacts: originalContacts }
  })
})
