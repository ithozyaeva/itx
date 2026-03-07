import { mount } from '@vue/test-utils'
import { describe, expect, it, vi } from 'vitest'

const mockUserData = {
  value: {
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
    profTags: [{ id: 1, title: 'Go' }, { id: 2, title: 'Vue' }],
    contacts: [],
    services: [],
  },
}

vi.mock('@/composables/useUser', () => ({
  useUser: () => mockUserData,
}))

vi.mock('@/services/profile', () => ({
  profileService: {
    getAllProfTags: vi.fn().mockResolvedValue([
      { id: 1, title: 'Go' },
      { id: 2, title: 'Vue' },
      { id: 3, title: 'React' },
      { id: 4, title: 'Python' },
    ]),
    updateTags: vi.fn().mockResolvedValue(true),
  },
}))

vi.mock('@/services/errorService', () => ({
  handleError: vi.fn(),
}))

vi.mock('itx-ui-kit', () => ({
  Typography: { template: '<div><slot /></div>' },
}))

vi.mock('lucide-vue-next', () => ({
  Edit: { template: '<span />' },
  Loader2: { template: '<span />' },
}))

// Mock the complex reka-ui based components
vi.mock('@/components/ui/combobox', () => ({
  Combobox: { template: '<div><slot /></div>', props: ['modelValue', 'open', 'ignoreFilter'] },
  ComboboxAnchor: { template: '<div><slot /></div>', props: ['asChild'] },
  ComboboxEmpty: { template: '<div />' },
  ComboboxGroup: { template: '<div><slot /></div>' },
  ComboboxInput: { template: '<input />', props: ['modelValue', 'readonly', 'asChild'] },
  ComboboxItem: { template: '<div><slot /></div>', props: ['value'] },
  ComboboxList: { template: '<div><slot /></div>' },
}))

vi.mock('@/components/ui/tags-input', () => ({
  TagsInput: {
    template: '<div class="tags-input"><slot /></div>',
    props: ['modelValue', 'displayValue', 'convertValue'],
    emits: ['update:modelValue', 'removeTag'],
  },
  TagsInputInput: { template: '<input />', props: ['placeholder', 'readonly'] },
  TagsInputItem: { template: '<span class="tag-item"><slot /></span>', props: ['value'] },
  TagsInputItemDelete: { template: '<button class="tag-delete">x</button>' },
  TagsInputItemText: { template: '<span />' },
}))

vi.mock('@/components/ui/button', () => ({
  Button: {
    template: '<button :disabled="$attrs.disabled" @click="$emit(\'click\')"><slot /></button>',
  },
}))

import ProfTagsForm from '@/components/Profile/ProfTagsForm.vue'
import { profileService } from '@/services/profile'

describe('ProfTagsForm', () => {
  it('renders without errors', () => {
    const wrapper = mount(ProfTagsForm)
    expect(wrapper.exists()).toBe(true)
  })

  it('displays title', () => {
    const wrapper = mount(ProfTagsForm)
    expect(wrapper.text()).toContain('Проф теги')
  })

  it('renders tag items from user data', () => {
    const wrapper = mount(ProfTagsForm)
    const tagItems = wrapper.findAll('.tag-item')
    expect(tagItems.length).toBe(2)
  })

  it('does not show save button in view mode', () => {
    const wrapper = mount(ProfTagsForm)
    const saveBtn = wrapper.findAll('button').find(b => b.text().includes('Сохранить'))
    expect(saveBtn).toBeUndefined()
  })

  it('shows save button after toggling edit mode', async () => {
    const wrapper = mount(ProfTagsForm)
    const editIcon = wrapper.find('.absolute')
    await editIcon.trigger('click')

    const saveBtn = wrapper.findAll('button').find(b => b.text().includes('Сохранить'))
    expect(saveBtn).toBeDefined()
  })

  it('calls profileService.updateTags on submit', async () => {
    const wrapper = mount(ProfTagsForm)
    const editIcon = wrapper.find('.absolute')
    await editIcon.trigger('click')

    const saveBtn = wrapper.findAll('button').find(b => b.text().includes('Сохранить'))
    await saveBtn!.trigger('click')
    await vi.dynamicImportSettled()

    expect(profileService.updateTags).toHaveBeenCalled()
  })

  it('loads all prof tags on mount', async () => {
    mount(ProfTagsForm)
    await vi.dynamicImportSettled()
    expect(profileService.getAllProfTags).toHaveBeenCalled()
  })

  it('shows tag delete buttons only in edit mode', async () => {
    const wrapper = mount(ProfTagsForm)
    // In view mode, no delete buttons
    expect(wrapper.findAll('.tag-delete').length).toBe(0)

    // Switch to edit mode
    const editIcon = wrapper.find('.absolute')
    await editIcon.trigger('click')

    expect(wrapper.findAll('.tag-delete').length).toBe(2)
  })
})
