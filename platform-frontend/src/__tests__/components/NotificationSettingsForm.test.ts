import { mount } from '@vue/test-utils'
import { describe, expect, it, vi } from 'vitest'

vi.mock('@/services/notificationSettings', () => ({
  notificationSettingsService: {
    get: vi.fn().mockResolvedValue({
      muteAll: false,
      newEvents: true,
      remindWeek: true,
      remindDay: true,
      remindHour: true,
      eventStart: true,
      eventUpdates: true,
      eventCancelled: true,
    }),
    update: vi.fn().mockResolvedValue({
      muteAll: false,
      newEvents: true,
      remindWeek: true,
      remindDay: true,
      remindHour: true,
      eventStart: true,
      eventUpdates: true,
      eventCancelled: true,
    }),
  },
}))

vi.mock('@/services/errorService', () => ({
  handleError: vi.fn(),
}))

vi.mock('itx-ui-kit', () => ({
  Typography: { template: '<div><slot /></div>' },
}))

vi.mock('lucide-vue-next', () => ({
  Bell: { template: '<span />' },
  Loader2: { template: '<span />' },
}))

import NotificationSettingsForm from '@/components/Profile/NotificationSettingsForm.vue'

describe('NotificationSettingsForm', () => {
  const globalConfig = {
    stubs: {
      Button: { template: '<button :disabled="$attrs.disabled" @click="$emit(\'click\')"><slot /></button>' },
    },
  }

  it('renders without errors', () => {
    const wrapper = mount(NotificationSettingsForm, { global: globalConfig })
    expect(wrapper.exists()).toBe(true)
  })

  it('shows title', () => {
    const wrapper = mount(NotificationSettingsForm, { global: globalConfig })
    expect(wrapper.text()).toContain('Уведомления в Telegram')
  })

  it('shows loading spinner initially', () => {
    const wrapper = mount(NotificationSettingsForm, { global: globalConfig })
    // isLoading starts true, so settings toggles should not be visible yet
    expect(wrapper.text()).not.toContain('Отключить все уведомления')
  })

  it('shows settings after loading', async () => {
    const wrapper = mount(NotificationSettingsForm, { global: globalConfig })
    await vi.dynamicImportSettled()
    await wrapper.vm.$nextTick()
    await wrapper.vm.$nextTick()

    expect(wrapper.text()).toContain('Отключить все уведомления')
    expect(wrapper.text()).toContain('Новые события')
    expect(wrapper.text()).toContain('Напоминание за неделю')
  })

  it('shows all toggle items', async () => {
    const wrapper = mount(NotificationSettingsForm, { global: globalConfig })
    await vi.dynamicImportSettled()
    await wrapper.vm.$nextTick()
    await wrapper.vm.$nextTick()

    const expectedLabels = [
      'Новые события',
      'Напоминание за неделю',
      'Напоминание за день',
      'Напоминание за час',
      'Начало события',
      'Изменения событий',
      'Отмена событий',
    ]

    for (const label of expectedLabels) {
      expect(wrapper.text()).toContain(label)
    }
  })

  it('does not show save button initially (no changes)', async () => {
    const wrapper = mount(NotificationSettingsForm, { global: globalConfig })
    await vi.dynamicImportSettled()
    await wrapper.vm.$nextTick()
    await wrapper.vm.$nextTick()

    const saveBtn = wrapper.findAll('button').find(b => b.text().includes('Сохранить'))
    expect(saveBtn).toBeUndefined()
  })

  it('shows save button after toggling a setting', async () => {
    const wrapper = mount(NotificationSettingsForm, { global: globalConfig })
    await vi.dynamicImportSettled()
    await wrapper.vm.$nextTick()
    await wrapper.vm.$nextTick()

    // Find and click the mute all row
    const muteAllRow = wrapper.find('.cursor-pointer')
    await muteAllRow.trigger('click')

    const saveBtn = wrapper.findAll('button').find(b => b.text().includes('Сохранить'))
    expect(saveBtn).toBeDefined()
  })

  it('toggle switches have correct aria-checked state', async () => {
    const wrapper = mount(NotificationSettingsForm, { global: globalConfig })
    await vi.dynamicImportSettled()
    await wrapper.vm.$nextTick()
    await wrapper.vm.$nextTick()

    const switches = wrapper.findAll('button[role="switch"]')
    // First is mute all (false)
    expect(switches[0].attributes('aria-checked')).toBe('false')
    // Second is newEvents (true)
    expect(switches[1].attributes('aria-checked')).toBe('true')
  })
})
