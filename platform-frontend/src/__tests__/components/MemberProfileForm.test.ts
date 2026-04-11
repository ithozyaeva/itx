import { mount } from '@vue/test-utils'
import { describe, expect, it, vi } from 'vitest'
import { ref } from 'vue'

const mockUser = ref({
  id: 1,
  telegramID: 123,
  tg: 'testuser',
  birthday: '1990-05-15',
  firstName: 'John',
  lastName: 'Doe',
  bio: 'Hello world',
  grade: 'Senior',
  company: 'ACME',
  avatarUrl: 'https://example.com/avatar.jpg',
  roles: ['SUBSCRIBER'],
})

vi.mock('@/composables/useUser', () => ({
  useUser: () => mockUser,
}))

vi.mock('@/services/profile', () => ({
  profileService: {
    updateMe: vi.fn(),
    uploadAvatar: vi.fn(),
  },
}))

vi.mock('@/services/points', () => ({
  pointsService: {
    getMyPoints: vi.fn().mockResolvedValue({ balance: 42 }),
  },
}))

vi.mock('@/services/errorService', () => ({
  handleError: vi.fn(),
}))

vi.mock('@/components/ui/typography', () => ({
  Typography: { template: '<div><slot /></div>' },
}))

vi.mock('lucide-vue-next', () => ({
  Camera: { template: '<span />' },
  Edit: { template: '<span />' },
  Loader2: { template: '<span />' },
  Star: { template: '<span />' },
}))

import MemberProfileForm from '@/components/Profile/MemberProfileForm.vue'
import { profileService } from '@/services/profile'

describe('MemberProfileForm', () => {
  const globalConfig = {
    stubs: {
      Button: { template: '<button @click="$emit(\'click\')"><slot /></button>' },
    },
  }

  it('renders without errors', () => {
    const wrapper = mount(MemberProfileForm, { global: globalConfig })
    expect(wrapper.exists()).toBe(true)
  })

  it('displays user name in view mode', () => {
    const wrapper = mount(MemberProfileForm, { global: globalConfig })
    expect(wrapper.text()).toContain('John')
    expect(wrapper.text()).toContain('Doe')
  })

  it('displays telegram handle in view mode', () => {
    const wrapper = mount(MemberProfileForm, { global: globalConfig })
    expect(wrapper.text()).toContain('testuser')
  })

  it('displays grade and company in view mode', () => {
    const wrapper = mount(MemberProfileForm, { global: globalConfig })
    expect(wrapper.text()).toContain('Senior')
    expect(wrapper.text()).toContain('ACME')
  })

  it('displays bio in view mode', () => {
    const wrapper = mount(MemberProfileForm, { global: globalConfig })
    expect(wrapper.text()).toContain('Hello world')
  })

  it('switches to edit mode when edit icon is clicked', async () => {
    const wrapper = mount(MemberProfileForm, { global: globalConfig })
    const editIcon = wrapper.find('.absolute.right-0')
    await editIcon.trigger('click')

    const inputs = wrapper.findAll('input[type="text"]')
    expect(inputs.length).toBeGreaterThanOrEqual(2)
  })

  it('shows form inputs in edit mode with correct values', async () => {
    const wrapper = mount(MemberProfileForm, { global: globalConfig })
    const editIcon = wrapper.find('.absolute.right-0')
    await editIcon.trigger('click')

    const firstNameInput = wrapper.find('input[placeholder="Имя"]')
    const lastNameInput = wrapper.find('input[placeholder="Фамилия"]')
    expect((firstNameInput.element as HTMLInputElement).value).toBe('John')
    expect((lastNameInput.element as HTMLInputElement).value).toBe('Doe')
  })

  it('shows grade and company inputs in edit mode', async () => {
    const wrapper = mount(MemberProfileForm, { global: globalConfig })
    const editIcon = wrapper.find('.absolute.right-0')
    await editIcon.trigger('click')

    const gradeInput = wrapper.find('input[placeholder="Грейд (Junior, Middle, Senior...)"]')
    const companyInput = wrapper.find('input[placeholder="Место работы"]')
    expect(gradeInput.exists()).toBe(true)
    expect(companyInput.exists()).toBe(true)
    expect((gradeInput.element as HTMLInputElement).value).toBe('Senior')
    expect((companyInput.element as HTMLInputElement).value).toBe('ACME')
  })

  it('calls profileService.updateMe on submit', async () => {
    const wrapper = mount(MemberProfileForm, { global: globalConfig })
    const editIcon = wrapper.find('.absolute.right-0')
    await editIcon.trigger('click')

    const saveBtn = wrapper.findAll('button').find(b => b.text().includes('Сохранить'))
    expect(saveBtn).toBeDefined()
    await saveBtn!.trigger('click')

    expect(profileService.updateMe).toHaveBeenCalled()
  })

  it('exits edit mode after submit', async () => {
    const wrapper = mount(MemberProfileForm, { global: globalConfig })
    const editIcon = wrapper.find('.absolute.right-0')
    await editIcon.trigger('click')

    const saveBtn = wrapper.findAll('button').find(b => b.text().includes('Сохранить'))
    await saveBtn!.trigger('click')

    const nameInput = wrapper.find('input[placeholder="Имя"]')
    expect(nameInput.exists()).toBe(false)
  })

  it('renders avatar image', () => {
    const wrapper = mount(MemberProfileForm, { global: globalConfig })
    const img = wrapper.find('img')
    expect(img.exists()).toBe(true)
    expect(img.attributes('src')).toBe('https://example.com/avatar.jpg')
  })

  it('shows file upload in edit mode', async () => {
    const wrapper = mount(MemberProfileForm, { global: globalConfig })
    const editIcon = wrapper.find('.absolute.right-0')
    await editIcon.trigger('click')

    const fileInput = wrapper.find('input[type="file"]')
    expect(fileInput.exists()).toBe(true)
  })

  it('shows points balance after mount', async () => {
    const wrapper = mount(MemberProfileForm, { global: globalConfig })
    await vi.dynamicImportSettled()
    await wrapper.vm.$nextTick()
    expect(wrapper.text()).toContain('42')
  })
})
