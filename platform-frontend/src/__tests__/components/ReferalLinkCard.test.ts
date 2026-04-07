import type { ReferalLink } from '@/models/referals'
import { mount } from '@vue/test-utils'
import { describe, expect, it, vi } from 'vitest'
import { ref } from 'vue'

const mockUserData = ref({
  id: 1,
  telegramID: 123,
  tg: 'currentuser',
  birthday: '',
  firstName: 'Current',
  lastName: 'User',
  bio: '',
  grade: '',
  company: '',
  avatarUrl: '',
  roles: ['SUBSCRIBER'],
})

vi.mock('@/composables/useUser', () => ({
  useUser: () => mockUserData,
}))

vi.mock('@/composables/useDictionary', () => ({
  useDictionary: () => ({
    gradesObject: { value: { junior: 'Junior', middle: 'Middle', senior: 'Senior' } },
    referalLinkStatusesObject: { value: { active: 'Активна', freezed: 'Заморожена' } },
  }),
}))

vi.mock('@/services/referals', () => ({
  referalLinkService: {
    updateLink: vi.fn().mockResolvedValue({ id: 1 }),
    deleteLink: vi.fn().mockResolvedValue(true),
    trackConversion: vi.fn().mockResolvedValue(true),
  },
}))

vi.mock('@/services/errorService', () => ({
  handleError: vi.fn(),
}))

vi.mock('@/lib/utils', () => ({
  dateFormatter: { format: (d: Date) => d.toISOString().slice(0, 10) },
}))

vi.mock('itx-ui-kit', () => ({
  Typography: { template: '<div><slot /></div>' },
}))

vi.mock('lucide-vue-next', () => ({
  Check: { template: '<span />' },
  Loader2: { template: '<span />' },
  Pencil: { template: '<span />' },
  Trash: { template: '<span />' },
}))

vi.mock('@/components/ConfirmDialog.vue', () => ({
  default: {
    template: '<div class="confirm-dialog"><slot name="trigger" /><button class="confirm-btn" @click="$emit(\'confirm\')">confirm</button></div>',
    emits: ['confirm'],
    props: ['title', 'description', 'confirmLabel'],
  },
}))

vi.mock('@/components/referals/ReferalLinkForm.vue', () => ({
  default: {
    template: '<div class="referal-form"><button class="save-btn" @click="$emit(\'save\', {})">save</button><button class="cancel-btn" @click="$emit(\'cancel\')">cancel</button></div>',
    emits: ['save', 'cancel'],
    props: ['link', 'isSaving', 'title'],
  },
}))

vi.mock('@/components/ui/badge', () => ({
  Badge: { template: '<span class="badge"><slot /></span>', props: ['variant'] },
}))

import ReferalLinkCard from '@/components/referals/ReferalLinkCard.vue'
import { referalLinkService } from '@/services/referals'

function createLink(overrides: Partial<ReferalLink> = {}): ReferalLink {
  return {
    id: 1,
    author: {
      id: 1,
      telegramID: 123,
      tg: 'currentuser',
      birthday: '',
      firstName: 'Current',
      lastName: 'User',
      bio: '',
      grade: '',
      company: '',
      avatarUrl: '',
      roles: ['SUBSCRIBER'],
    },
    company: 'ACME Corp',
    grade: 'middle',
    profTags: [{ id: 1, title: 'Go' }],
    status: 'active',
    vacationsCount: 3,
    conversionsCount: 0,
    hasConverted: false,
    updatedAt: '2026-01-15T10:00:00Z',
    ...overrides,
  }
}

describe('ReferalLinkCard', () => {
  it('renders without errors', () => {
    const wrapper = mount(ReferalLinkCard, {
      props: { link: createLink() },
    })
    expect(wrapper.exists()).toBe(true)
  })

  it('displays company name', () => {
    const wrapper = mount(ReferalLinkCard, {
      props: { link: createLink({ company: 'Yandex' }) },
    })
    expect(wrapper.text()).toContain('Yandex')
  })

  it('displays author name', () => {
    const wrapper = mount(ReferalLinkCard, {
      props: { link: createLink() },
    })
    expect(wrapper.text()).toContain('Current')
    expect(wrapper.text()).toContain('User')
  })

  it('shows edit and delete buttons when user is owner', () => {
    const wrapper = mount(ReferalLinkCard, {
      props: { link: createLink({ author: { ...createLink().author, id: 1 } }) },
    })
    // Owner controls should be visible
    expect(wrapper.find('.confirm-dialog').exists()).toBe(true)
  })

  it('hides owner controls when user is not the owner', () => {
    const wrapper = mount(ReferalLinkCard, {
      props: { link: createLink({ author: { ...createLink().author, id: 999, tg: 'otheruser' } }) },
    })
    expect(wrapper.find('.confirm-dialog').exists()).toBe(false)
  })

  it('shows convert button for non-owner on active link', () => {
    const wrapper = mount(ReferalLinkCard, {
      props: {
        link: createLink({
          author: { ...createLink().author, id: 999, tg: 'otheruser' },
          status: 'active',
        }),
      },
    })
    expect(wrapper.text()).toContain('Откликнуться')
  })

  it('does not show convert button for owner', () => {
    const wrapper = mount(ReferalLinkCard, {
      props: { link: createLink() },
    })
    expect(wrapper.text()).not.toContain('Откликнуться')
  })

  it('shows "Открыть чат снова" when already converted', () => {
    const wrapper = mount(ReferalLinkCard, {
      props: {
        link: createLink({
          author: { ...createLink().author, id: 999, tg: 'otheruser' },
          status: 'active',
          hasConverted: true,
        }),
      },
    })
    expect(wrapper.text()).toContain('Открыть чат снова')
  })

  it('emits deleted when delete is confirmed', async () => {
    const wrapper = mount(ReferalLinkCard, {
      props: { link: createLink() },
    })
    const confirmBtn = wrapper.find('.confirm-btn')
    await confirmBtn.trigger('click')
    await vi.dynamicImportSettled()

    expect(referalLinkService.deleteLink).toHaveBeenCalledWith(1)
    expect(wrapper.emitted('deleted')).toBeDefined()
  })

  it('enters editing mode when edit button is clicked', async () => {
    const wrapper = mount(ReferalLinkCard, {
      props: { link: createLink() },
    })

    // Find the pencil button (non-confirm-dialog button)
    const editButtons = wrapper.findAll('button').filter(b => !b.classes().includes('confirm-btn'))
    const editBtn = editButtons[0]
    await editBtn.trigger('click')
    await wrapper.vm.$nextTick()

    expect(wrapper.find('.referal-form').exists()).toBe(true)
  })

  it('exits editing mode on cancel', async () => {
    const wrapper = mount(ReferalLinkCard, {
      props: { link: createLink() },
    })

    const vm = wrapper.vm as any
    vm.isEditing = true
    await wrapper.vm.$nextTick()

    const cancelBtn = wrapper.find('.cancel-btn')
    await cancelBtn.trigger('click')
    await wrapper.vm.$nextTick()

    expect(wrapper.find('.referal-form').exists()).toBe(false)
  })

  it('displays prof tags', () => {
    const wrapper = mount(ReferalLinkCard, {
      props: { link: createLink({ profTags: [{ id: 1, title: 'Go' }, { id: 2, title: 'Vue' }] }) },
    })
    expect(wrapper.text()).toContain('Go, Vue')
  })

  it('displays vacations count', () => {
    const wrapper = mount(ReferalLinkCard, {
      props: { link: createLink({ vacationsCount: 5 }) },
    })
    expect(wrapper.text()).toContain('5')
  })
})
